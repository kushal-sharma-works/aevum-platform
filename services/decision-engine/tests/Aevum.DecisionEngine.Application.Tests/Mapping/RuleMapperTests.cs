using FluentAssertions;
using Aevum.DecisionEngine.Application.DTOs;
using Aevum.DecisionEngine.Application.Mapping;
using Aevum.DecisionEngine.Domain.Enums;
using Aevum.DecisionEngine.Domain.Models;

namespace Aevum.DecisionEngine.Application.Tests.Mapping;

public sealed class RuleMapperTests
{
    [Fact]
    public void CreateRuleRequest_ToDomain_ShouldMapAndApplyDefaults()
    {
        var request = new CreateRuleRequest
        {
            Name = "rule-1",
            Description = "desc",
            Conditions =
            [
                new RuleConditionDto
                {
                    Field = "amount",
                    Operator = ComparisonOperator.GreaterThan,
                    Value = 100,
                    NestedConditions =
                    [
                        new RuleConditionDto
                        {
                            Field = "country",
                            Operator = ComparisonOperator.Equals,
                            Value = "DE"
                        }
                    ]
                }
            ],
            Actions =
            [
                new RuleActionDto
                {
                    Type = ActionType.StoreDecision,
                    Parameters = new Dictionary<string, object> { ["key"] = "value" },
                    Order = 1,
                    Description = "action"
                }
            ],
            Priority = 7,
            CreatedBy = "tester",
            Metadata = null
        };

        var result = request.ToDomain();

        result.Name.Should().Be("rule-1");
        result.Status.Should().Be(RuleStatus.Draft);
        result.Version.Should().Be(1);
        result.Metadata.Should().BeEmpty();
        result.Conditions[0].NestedConditions.Should().HaveCount(1);
    }

    [Fact]
    public void UpdateRuleRequest_ToDomain_ShouldMapStatusFallback()
    {
        var request = new UpdateRuleRequest
        {
            Name = "rule-1-updated",
            Conditions =
            [
                new RuleConditionDto
                {
                    Field = "amount",
                    Operator = ComparisonOperator.LessThan,
                    Value = 200
                }
            ],
            Actions =
            [
                new RuleActionDto
                {
                    Type = ActionType.SendNotification,
                    Parameters = new Dictionary<string, object>(),
                    Order = 0
                }
            ],
            Priority = 3,
            Status = null,
            Metadata = null
        };

        var result = request.ToDomain("rule-id");

        result.Id.Should().Be("rule-id");
        result.Status.Should().Be(RuleStatus.Draft);
        result.Version.Should().Be(0);
        result.Metadata.Should().BeEmpty();
    }

    [Fact]
    public void Rule_ToResponse_ShouldMapAllFields()
    {
        var now = DateTimeOffset.UtcNow;
        var rule = new Rule
        {
            Id = "rule-id",
            Name = "rule",
            Conditions =
            [
                new RuleCondition
                {
                    Field = "score",
                    Operator = ComparisonOperator.GreaterThanOrEqual,
                    Value = 80
                }
            ],
            Actions =
            [
                new RuleAction
                {
                    Type = ActionType.LogEvent,
                    Parameters = new Dictionary<string, object> { ["message"] = "ok" },
                    Order = 2
                }
            ],
            Status = RuleStatus.Active,
            Version = 2,
            Priority = 9,
            CreatedAt = now,
            UpdatedAt = now,
            Metadata = new Dictionary<string, string> { ["env"] = "test" }
        };

        var response = rule.ToResponse();

        response.Id.Should().Be("rule-id");
        response.Status.Should().Be(RuleStatus.Active);
        response.Conditions.Should().HaveCount(1);
        response.Actions.Should().HaveCount(1);
        response.Metadata.Should().ContainKey("env");
    }
}
