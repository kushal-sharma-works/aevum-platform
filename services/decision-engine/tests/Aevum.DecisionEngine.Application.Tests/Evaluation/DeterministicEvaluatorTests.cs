using FluentAssertions;
using Aevum.DecisionEngine.Application.Evaluation;
using Aevum.DecisionEngine.Domain.Enums;
using Aevum.DecisionEngine.Domain.Models;

namespace Aevum.DecisionEngine.Application.Tests.Evaluation;

public sealed class DeterministicEvaluatorTests
{
    private readonly DeterministicEvaluator _evaluator;
    private readonly TimeProvider _timeProvider;

    public DeterministicEvaluatorTests()
    {
        _timeProvider = TimeProvider.System;
        _evaluator = new DeterministicEvaluator(_timeProvider);
    }

    [Fact]
    public void ComputeHash_WithSameInputs_ShouldReturnSameHash()
    {
        // Arrange
        var rule = CreateTestRule();
        var context = CreateTestContext();

        // Act
        var hash1 = _evaluator.ComputeHash(rule, context);
        var hash2 = _evaluator.ComputeHash(rule, context);

        // Assert
        hash1.Should().Be(hash2);
    }

    [Fact]
    public void ComputeHash_WithDifferentInputs_ShouldReturnDifferentHash()
    {
        // Arrange
        var rule = CreateTestRule();
        var context1 = CreateTestContext();
        var context2 = context1 with
        {
            Data = new Dictionary<string, object> { ["amount"] = 200 }
        };

        // Act
        var hash1 = _evaluator.ComputeHash(rule, context1);
        var hash2 = _evaluator.ComputeHash(rule, context2);

        // Assert
        hash1.Should().NotBe(hash2);
    }

    [Fact]
    public void Evaluate_EqualsOperator_ShouldMatch()
    {
        // Arrange
        var rule = CreateRuleWithCondition(ComparisonOperator.Equals, "status", "active");
        var context = new EvaluationContext
        {
            Data = new Dictionary<string, object> { ["status"] = "active" },
            RequestId = "req-1",
            Timestamp = DateTimeOffset.UtcNow
        };

        // Act
        var result = _evaluator.Evaluate(rule, context);

        // Assert
        result.IsMatch.Should().BeTrue();
        result.Status.Should().Be(DecisionStatus.Approved);
    }

    [Fact]
    public void Evaluate_GreaterThanOperator_ShouldMatch()
    {
        // Arrange
        var rule = CreateRuleWithCondition(ComparisonOperator.GreaterThan, "amount", 100);
        var context = new EvaluationContext
        {
            Data = new Dictionary<string, object> { ["amount"] = 150 },
            RequestId = "req-1",
            Timestamp = DateTimeOffset.UtcNow
        };

        // Act
        var result = _evaluator.Evaluate(rule, context);

        // Assert
        result.IsMatch.Should().BeTrue();
    }

    [Fact]
    public void Evaluate_ContainsOperator_ShouldMatch()
    {
        // Arrange
        var rule = CreateRuleWithCondition(ComparisonOperator.Contains, "tags", "premium");
        var context = new EvaluationContext
        {
            Data = new Dictionary<string, object> { ["tags"] = "premium,vip,gold" },
            RequestId = "req-1",
            Timestamp = DateTimeOffset.UtcNow
        };

        // Act
        var result = _evaluator.Evaluate(rule, context);

        // Assert
        result.IsMatch.Should().BeTrue();
    }

    [Fact]
    public void Evaluate_RegexOperator_ShouldMatch()
    {
        // Arrange
        var rule = CreateRuleWithCondition(ComparisonOperator.Regex, "email", @"^[\w\.-]+@[\w\.-]+\.\w+$");
        var context = new EvaluationContext
        {
            Data = new Dictionary<string, object> { ["email"] = "test@example.com" },
            RequestId = "req-1",
            Timestamp = DateTimeOffset.UtcNow
        };

        // Act
        var result = _evaluator.Evaluate(rule, context);

        // Assert
        result.IsMatch.Should().BeTrue();
    }

    [Fact]
    public void Evaluate_MultipleConditionsWithAnd_ShouldMatch()
    {
        // Arrange
        var rule = new Rule
        {
            Id = "rule-1",
            Name = "Multi Condition Rule",
            Conditions =
            [
                new RuleCondition
                {
                    Field = "amount",
                    Operator = ComparisonOperator.GreaterThan,
                    Value = 100,
                    LogicalOperator = LogicalOperator.And
                },
                new RuleCondition
                {
                    Field = "status",
                    Operator = ComparisonOperator.Equals,
                    Value = "active"
                }
            ],
            Actions = [CreateTestAction()],
            Status = RuleStatus.Active,
            Version = 1,
            Priority = 10,
            CreatedAt = DateTimeOffset.UtcNow,
            UpdatedAt = DateTimeOffset.UtcNow
        };

        var context = new EvaluationContext
        {
            Data = new Dictionary<string, object>
            {
                ["amount"] = 150,
                ["status"] = "active"
            },
            RequestId = "req-1",
            Timestamp = DateTimeOffset.UtcNow
        };

        // Act
        var result = _evaluator.Evaluate(rule, context);

        // Assert
        result.IsMatch.Should().BeTrue();
    }

    [Theory]
    [InlineData(50, false)]
    [InlineData(100, false)]
    [InlineData(101, true)]
    [InlineData(1000, true)]
    public void Evaluate_NumericComparison_ShouldEvaluateCorrectly(int amount, bool expected)
    {
        // Arrange
        var rule = CreateRuleWithCondition(ComparisonOperator.GreaterThan, "amount", 100);
        var context = new EvaluationContext
        {
            Data = new Dictionary<string, object> { ["amount"] = amount },
            RequestId = "req-1",
            Timestamp = DateTimeOffset.UtcNow
        };

        // Act
        var result = _evaluator.Evaluate(rule, context);

        // Assert
        result.IsMatch.Should().Be(expected);
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
            Actions = [CreateTestAction()],
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
            Timestamp = new DateTimeOffset(2024, 1, 1, 0, 0, 0, TimeSpan.Zero)
        };
    }

    private static Rule CreateRuleWithCondition(ComparisonOperator op, string field, object value)
    {
        return new Rule
        {
            Id = "rule-1",
            Name = "Test Rule",
            Conditions =
            [
                new RuleCondition
                {
                    Field = field,
                    Operator = op,
                    Value = value
                }
            ],
            Actions = [CreateTestAction()],
            Status = RuleStatus.Active,
            Version = 1,
            Priority = 10,
            CreatedAt = DateTimeOffset.UtcNow,
            UpdatedAt = DateTimeOffset.UtcNow
        };
    }

    private static RuleAction CreateTestAction()
    {
        return new RuleAction
        {
            Type = ActionType.StoreDecision,
            Parameters = new Dictionary<string, object> { ["key"] = "value" },
            Order = 1
        };
    }
}
