using Aevum.DecisionEngine.Application.DTOs;
using Aevum.DecisionEngine.Domain.Models;

namespace Aevum.DecisionEngine.Application.Mapping;

public static class DecisionMapper
{
    public static EvaluationContext ToEvaluationContext(this EvaluateDecisionRequest request)
    {
        return new EvaluationContext
        {
            Data = request.Context,
            RequestId = request.RequestId,
            Timestamp = DateTimeOffset.UtcNow, // Will be replaced by TimeProvider in service layer
            Metadata = request.Metadata ?? new Dictionary<string, string>()
        };
    }

    public static DecisionResponse ToResponse(this Decision decision)
    {
        return new DecisionResponse
        {
            Id = decision.Id,
            RuleId = decision.RuleId,
            RuleVersion = decision.RuleVersion,
            RequestId = decision.RequestId,
            Status = decision.Status,
            InputContext = decision.InputContext.ToDictionary(kvp => kvp.Key, kvp => kvp.Value),
            MatchedConditions = decision.MatchedConditions,
            ExecutedActions = decision.ExecutedActions.Select(a => a.ToDto()).ToList(),
            EvaluatedAt = decision.EvaluatedAt,
            DeterministicHash = decision.DeterministicHash,
            ErrorMessage = decision.ErrorMessage,
            OutputData = decision.OutputData.ToDictionary(kvp => kvp.Key, kvp => kvp.Value),
            EvaluationDurationMs = decision.EvaluationDurationMs
        };
    }
}
