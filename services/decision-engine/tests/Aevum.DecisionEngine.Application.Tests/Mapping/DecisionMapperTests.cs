using FluentAssertions;
using Aevum.DecisionEngine.Application.DTOs;
using Aevum.DecisionEngine.Application.Mapping;
using Aevum.DecisionEngine.Domain.Enums;
using Aevum.DecisionEngine.Domain.Models;

namespace Aevum.DecisionEngine.Application.Tests.Mapping;

public sealed class DecisionMapperTests
{
    [Fact]
    public void EvaluateDecisionRequest_ToEvaluationContext_ShouldMapWithMetadataFallback()
    {
        var request = new EvaluateDecisionRequest
        {
            RuleId = "rule-1",
            Context = new Dictionary<string, object> { ["amount"] = 123 },
            RequestId = "req-1",
            Metadata = null
        };

        var context = request.ToEvaluationContext();

        context.RequestId.Should().Be("req-1");
        context.Data.Should().ContainKey("amount");
        context.Metadata.Should().BeEmpty();
        context.Timestamp.Should().BeCloseTo(DateTimeOffset.UtcNow, TimeSpan.FromSeconds(2));
    }

    [Fact]
    public void Decision_ToResponse_ShouldMapCollections()
    {
        var decision = new Decision
        {
            Id = "dec-1",
            RuleId = "rule-1",
            RuleVersion = 3,
            RequestId = "req-1",
            Status = DecisionStatus.Approved,
            InputContext = new Dictionary<string, object> { ["amount"] = 100 },
            MatchedConditions = ["amount > 50"],
            ExecutedActions =
            [
                new RuleAction
                {
                    Type = ActionType.StoreDecision,
                    Parameters = new Dictionary<string, object> { ["key"] = "value" },
                    Order = 1
                }
            ],
            EvaluatedAt = DateTimeOffset.UtcNow,
            DeterministicHash = "hash",
            OutputData = new Dictionary<string, object> { ["approved"] = true },
            EvaluationDurationMs = 12
        };

        var response = decision.ToResponse();

        response.Id.Should().Be("dec-1");
        response.InputContext.Should().ContainKey("amount");
        response.ExecutedActions.Should().HaveCount(1);
        response.OutputData.Should().ContainKey("approved");
    }
}
