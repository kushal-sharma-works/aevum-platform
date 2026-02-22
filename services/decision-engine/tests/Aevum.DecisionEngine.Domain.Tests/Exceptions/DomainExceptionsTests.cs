using FluentAssertions;
using Aevum.DecisionEngine.Domain.Enums;
using Aevum.DecisionEngine.Domain.Exceptions;
using Aevum.DecisionEngine.Domain.Models;

namespace Aevum.DecisionEngine.Domain.Tests.Exceptions;

public sealed class DomainExceptionsTests
{
    [Fact]
    public void RuleNotFoundException_ShouldExposeProperties()
    {
        var ex = new RuleNotFoundException("rule-1", 2);

        ex.RuleId.Should().Be("rule-1");
        ex.Version.Should().Be(2);
        ex.Message.Should().Contain("rule-1");
    }

    [Fact]
    public void EvaluationException_WithRuleContext_ShouldExposeFields()
    {
        var ex = new EvaluationException("rule-1", "amount", "invalid number");

        ex.RuleId.Should().Be("rule-1");
        ex.FieldName.Should().Be("amount");
        ex.Message.Should().Contain("amount");
    }

    [Fact]
    public void InvalidRuleException_ShouldContainValidationErrors()
    {
        var ex = new InvalidRuleException("bad rule")
        {
            ValidationErrors = ["name required", "conditions required"]
        };

        ex.ValidationErrors.Should().HaveCount(2);
        ex.ErrorCode.Should().Be("DOMAIN_ERROR");
    }

    [Fact]
    public void EvaluationContextAndResult_ShouldStoreValues()
    {
        var context = new EvaluationContext
        {
            Data = new Dictionary<string, object> { ["score"] = 90 },
            RequestId = "req-1",
            Timestamp = DateTimeOffset.UtcNow
        };

        var result = new EvaluationResult
        {
            IsMatch = true,
            MatchedConditions = ["score > 80"],
            ActionsToExecute =
            [
                new RuleAction
                {
                    Type = ActionType.StoreDecision,
                    Parameters = new Dictionary<string, object> { ["approved"] = true },
                    Order = 1
                }
            ],
            Status = DecisionStatus.Approved,
            DeterministicHash = "hash"
        };

        context.Data.Should().ContainKey("score");
        result.IsMatch.Should().BeTrue();
        result.ActionsToExecute.Should().HaveCount(1);
    }
}
