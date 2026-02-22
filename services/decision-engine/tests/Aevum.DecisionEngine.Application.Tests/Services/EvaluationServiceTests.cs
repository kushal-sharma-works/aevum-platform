using FluentAssertions;
using NSubstitute;
using Aevum.DecisionEngine.Application.Services;
using Aevum.DecisionEngine.Domain.Enums;
using Aevum.DecisionEngine.Domain.Interfaces;
using Aevum.DecisionEngine.Domain.Models;

namespace Aevum.DecisionEngine.Application.Tests.Services;

public sealed class EvaluationServiceTests
{
    private readonly IDeterministicEvaluator _evaluator;
    private readonly IDecisionRepository _decisionRepository;
    private readonly IEventTimelineClient _eventTimelineClient;
    private readonly TimeProvider _timeProvider;
    private readonly EvaluationService _service;

    public EvaluationServiceTests()
    {
        _evaluator = Substitute.For<IDeterministicEvaluator>();
        _decisionRepository = Substitute.For<IDecisionRepository>();
        _eventTimelineClient = Substitute.For<IEventTimelineClient>();
        _timeProvider = TimeProvider.System;

        _service = new EvaluationService(
            _evaluator,
            _decisionRepository,
            _eventTimelineClient,
            _timeProvider);
    }

    [Fact]
    public async Task EvaluateAsync_ShouldCreateDecision()
    {
        // Arrange
        var rule = CreateTestRule();
        var context = CreateTestContext();
        var hash = "testhash123";

        _evaluator.ComputeHash(rule, context).Returns(hash);
        _decisionRepository.GetByRequestIdAsync(context.RequestId, Arg.Any<CancellationToken>())
            .Returns((Decision?)null);
        _evaluator.Evaluate(rule, context).Returns(new EvaluationResult
        {
            IsMatch = true,
            MatchedConditions = ["amount > 100"],
            ActionsToExecute = rule.Actions.ToList(),
            Status = DecisionStatus.Approved,
            DeterministicHash = hash
        });

        _decisionRepository.CreateAsync(Arg.Any<Decision>(), Arg.Any<CancellationToken>())
            .Returns(callInfo => callInfo.Arg<Decision>());

        // Act
        var result = await _service.EvaluateAsync(rule, context, CancellationToken.None);

        // Assert
        result.Should().NotBeNull();
        result.Status.Should().Be(DecisionStatus.Approved);
        result.DeterministicHash.Should().Be(hash);
        await _decisionRepository.Received(1).CreateAsync(Arg.Any<Decision>(), Arg.Any<CancellationToken>());
    }

    [Fact]
    public async Task EvaluateAsync_WithExistingRequestId_ShouldReturnCachedDecision()
    {
        // Arrange
        var rule = CreateTestRule();
        var context = CreateTestContext();
        var hash = "testhash123";
        var existingDecision = CreateTestDecision();

        _evaluator.ComputeHash(rule, context).Returns(hash);
        _decisionRepository.GetByRequestIdAsync(context.RequestId, Arg.Any<CancellationToken>())
            .Returns(existingDecision);

        // Act
        var result = await _service.EvaluateAsync(rule, context, CancellationToken.None);

        // Assert
        result.Should().Be(existingDecision);
        _evaluator.DidNotReceive().Evaluate(Arg.Any<Rule>(), Arg.Any<EvaluationContext>());
        await _decisionRepository.DidNotReceive().CreateAsync(Arg.Any<Decision>(), Arg.Any<CancellationToken>());
    }

    [Fact]
    public async Task EvaluateAsync_WhenCreateFailsButDecisionExists_ShouldReturnExistingDecision()
    {
        // Arrange
        var rule = CreateTestRule();
        var context = CreateTestContext();
        var hash = "testhash123";
        var existingDecision = CreateTestDecision();

        _decisionRepository.GetByRequestIdAsync(context.RequestId, Arg.Any<CancellationToken>())
            .Returns((Decision?)null, existingDecision);
        _evaluator.ComputeHash(rule, context).Returns(hash);
        _evaluator.Evaluate(rule, context).Returns(new EvaluationResult
        {
            IsMatch = true,
            MatchedConditions = ["amount > 100"],
            ActionsToExecute = rule.Actions.ToList(),
            Status = DecisionStatus.Approved,
            DeterministicHash = hash
        });

        _decisionRepository.CreateAsync(Arg.Any<Decision>(), Arg.Any<CancellationToken>())
            .Returns<Task<Decision>>(_ => throw new InvalidOperationException("duplicate key"));

        // Act
        var result = await _service.EvaluateAsync(rule, context, CancellationToken.None);

        // Assert
        result.Should().Be(existingDecision);
    }

    private static Rule CreateTestRule()
    {
        return new Rule
        {
            Id = "rule-1",
            Name = "Test Rule",
            Conditions =
            [
                new RuleCondition
                {
                    Field = "amount",
                    Operator = ComparisonOperator.GreaterThan,
                    Value = 100
                }
            ],
            Actions =
            [
                new RuleAction
                {
                    Type = ActionType.StoreDecision,
                    Parameters = new Dictionary<string, object> { ["key"] = "value" },
                    Order = 1
                }
            ],
            Status = RuleStatus.Active,
            Version = 1,
            Priority = 10,
            CreatedAt = DateTimeOffset.UtcNow,
            UpdatedAt = DateTimeOffset.UtcNow
        };
    }

    private static EvaluationContext CreateTestContext()
    {
        return new EvaluationContext
        {
            Data = new Dictionary<string, object> { ["amount"] = 150 },
            RequestId = "req-123",
            Timestamp = DateTimeOffset.UtcNow
        };
    }

    private static Decision CreateTestDecision()
    {
        return new Decision
        {
            Id = "decision-1",
            RuleId = "rule-1",
            RuleVersion = 1,
            RequestId = "req-123",
            Status = DecisionStatus.Approved,
            InputContext = new Dictionary<string, object> { ["amount"] = 150 },
            MatchedConditions = ["amount > 100"],
            ExecutedActions = [],
            EvaluatedAt = DateTimeOffset.UtcNow,
            DeterministicHash = "testhash123",
            EvaluationDurationMs = 10
        };
    }
}
