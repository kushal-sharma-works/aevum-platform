using System.Diagnostics;
using Aevum.DecisionEngine.Domain.Enums;
using Aevum.DecisionEngine.Domain.Interfaces;
using Aevum.DecisionEngine.Domain.Models;

namespace Aevum.DecisionEngine.Application.Services;

public sealed class EvaluationService(
    IDeterministicEvaluator evaluator,
    IDecisionRepository decisionRepository,
    IEventTimelineClient eventTimelineClient,
    TimeProvider timeProvider)
{
    private readonly IDeterministicEvaluator _evaluator = evaluator;
    private readonly IDecisionRepository _decisionRepository = decisionRepository;
    private readonly IEventTimelineClient _eventTimelineClient = eventTimelineClient;
    private readonly TimeProvider _timeProvider = timeProvider;

    public async Task<Decision> EvaluateAsync(
        Rule rule,
        EvaluationContext context,
        CancellationToken cancellationToken = default)
    {
        var stopwatch = Stopwatch.StartNew();
        var deterministicHash = _evaluator.ComputeHash(rule, context);

        // Check for idempotency
        var existingDecision = await _decisionRepository.ExistsByHashAsync(deterministicHash, cancellationToken);
        if (existingDecision)
        {
            var existing = (await _decisionRepository.GetByRequestIdAsync(context.RequestId, cancellationToken))!;
            return existing;
        }

        var result = _evaluator.Evaluate(rule, context);
        stopwatch.Stop();

        var decision = new Decision
        {
            Id = Guid.NewGuid().ToString(),
            RuleId = rule.Id,
            RuleVersion = rule.Version,
            RequestId = context.RequestId,
            Status = result.Status,
            InputContext = context.Data,
            MatchedConditions = result.MatchedConditions,
            ExecutedActions = result.ActionsToExecute,
            EvaluatedAt = _timeProvider.GetUtcNow(),
            DeterministicHash = deterministicHash,
            ErrorMessage = result.ErrorMessage,
            OutputData = result.OutputData,
            EvaluationDurationMs = stopwatch.ElapsedMilliseconds
        };

        var savedDecision = await _decisionRepository.CreateAsync(decision, cancellationToken);

        // Ingest event to timeline (fire and forget with error handling)
        _ = Task.Run(async () =>
        {
            try
            {
                await _eventTimelineClient.IngestEventAsync(
                    streamId: $"decision-{decision.Id}",
                    eventType: "decision.evaluated",
                    data: new
                    {
                        decision.Id,
                        decision.RuleId,
                        decision.RuleVersion,
                        decision.Status,
                        decision.EvaluatedAt,
                        decision.DeterministicHash
                    },
                    CancellationToken.None);
            }
            catch
            {
                // Silently fail - event ingestion is not critical
            }
        }, CancellationToken.None);

        return savedDecision;
    }
}
