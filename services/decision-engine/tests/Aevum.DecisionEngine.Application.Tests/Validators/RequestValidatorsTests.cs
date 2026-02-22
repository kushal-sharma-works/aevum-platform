using FluentAssertions;
using Aevum.DecisionEngine.Application.DTOs;
using Aevum.DecisionEngine.Application.Validators;
using Aevum.DecisionEngine.Domain.Enums;

namespace Aevum.DecisionEngine.Application.Tests.Validators;

public sealed class RequestValidatorsTests
{
    [Fact]
    public void CreateRuleRequestValidator_ShouldRejectEmptyName()
    {
        var validator = new CreateRuleRequestValidator();
        var request = new CreateRuleRequest
        {
            Name = string.Empty,
            Conditions = [new RuleConditionDto { Field = "amount", Operator = ComparisonOperator.GreaterThan, Value = 10 }],
            Actions = [new RuleActionDto { Type = ActionType.StoreDecision, Parameters = new Dictionary<string, object> { ["k"] = "v" }, Order = 1 }],
            Priority = 10
        };

        var result = validator.Validate(request);

        result.IsValid.Should().BeFalse();
        result.Errors.Should().Contain(e => e.ErrorMessage.Contains("Rule name is required"));
    }

    [Fact]
    public void CreateRuleRequestValidator_ShouldRejectTooDeepNestedConditions()
    {
        var validator = new CreateRuleRequestValidator();

        RuleConditionDto MakeNested(int depth)
        {
            if (depth == 0)
            {
                return new RuleConditionDto { Field = "amount", Operator = ComparisonOperator.GreaterThan, Value = 10 };
            }

            return new RuleConditionDto
            {
                Field = $"level_{depth}",
                Operator = ComparisonOperator.GreaterThan,
                Value = depth,
                NestedConditions = [MakeNested(depth - 1)]
            };
        }

        var request = new CreateRuleRequest
        {
            Name = "deep",
            Conditions = [MakeNested(6)],
            Actions = [new RuleActionDto { Type = ActionType.StoreDecision, Parameters = new Dictionary<string, object> { ["k"] = "v" }, Order = 1 }],
            Priority = 10
        };

        var result = validator.Validate(request);

        result.IsValid.Should().BeFalse();
        result.Errors.Should().Contain(e => e.ErrorMessage.Contains("Nested conditions exceed maximum depth"));
    }

    [Fact]
    public void EvaluateDecisionRequestValidator_ShouldAcceptValidRequest()
    {
        var validator = new EvaluateDecisionRequestValidator();
        var request = new EvaluateDecisionRequest
        {
            RuleId = "rule-1",
            RuleVersion = 1,
            Context = new Dictionary<string, object> { ["score"] = 85 },
            RequestId = "req-1"
        };

        var result = validator.Validate(request);

        result.IsValid.Should().BeTrue();
    }

    [Fact]
    public void EvaluateDecisionRequestValidator_ShouldRejectInvalidVersion()
    {
        var validator = new EvaluateDecisionRequestValidator();
        var request = new EvaluateDecisionRequest
        {
            RuleId = "rule-1",
            RuleVersion = 0,
            Context = new Dictionary<string, object> { ["score"] = 85 },
            RequestId = "req-1"
        };

        var result = validator.Validate(request);

        result.IsValid.Should().BeFalse();
        result.Errors.Should().Contain(e => e.ErrorMessage.Contains("RuleVersion must be positive"));
    }
}
