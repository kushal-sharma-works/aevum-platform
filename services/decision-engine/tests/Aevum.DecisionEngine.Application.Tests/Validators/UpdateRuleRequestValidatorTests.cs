using FluentAssertions;
using Aevum.DecisionEngine.Application.DTOs;
using Aevum.DecisionEngine.Application.Validators;
using Aevum.DecisionEngine.Domain.Enums;

namespace Aevum.DecisionEngine.Application.Tests.Validators;

public sealed class UpdateRuleRequestValidatorTests
{
    [Fact]
    public void Validate_ValidRequest_Succeeds()
    {
        var validator = new UpdateRuleRequestValidator();
        var request = new UpdateRuleRequest
        {
            Name = "rule",
            Description = "desc",
            Conditions =
            [
                new RuleConditionDto
                {
                    Field = "amount",
                    Operator = ComparisonOperator.GreaterThan,
                    Value = 10
                }
            ],
            Actions =
            [
                new RuleActionDto
                {
                    Type = ActionType.StoreDecision,
                    Parameters = new Dictionary<string, object> { ["k"] = "v" },
                    Order = 1
                }
            ],
            Priority = 5,
            Status = RuleStatus.Active,
            EffectiveFrom = DateTimeOffset.UtcNow,
            EffectiveUntil = DateTimeOffset.UtcNow.AddHours(1)
        };

        var result = validator.Validate(request);

        result.IsValid.Should().BeTrue();
    }

    [Fact]
    public void Validate_InvalidDateRange_Fails()
    {
        var validator = new UpdateRuleRequestValidator();
        var now = DateTimeOffset.UtcNow;
        var request = new UpdateRuleRequest
        {
            Name = "rule",
            Conditions = [new RuleConditionDto { Field = "x", Operator = ComparisonOperator.Equals, Value = 1 }],
            Actions = [new RuleActionDto { Type = ActionType.LogEvent, Parameters = new Dictionary<string, object>(), Order = 0 }],
            Priority = 1,
            EffectiveFrom = now,
            EffectiveUntil = now.AddMinutes(-1)
        };

        var result = validator.Validate(request);

        result.IsValid.Should().BeFalse();
        result.Errors.Should().Contain(e => e.ErrorMessage.Contains("EffectiveUntil must be after EffectiveFrom"));
    }
}
