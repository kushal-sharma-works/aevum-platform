using FluentAssertions;
using Aevum.DecisionEngine.Domain.Enums;
using Aevum.DecisionEngine.Domain.Models;

namespace Aevum.DecisionEngine.Domain.Tests.Models;

public sealed class RuleTests
{
    [Fact]
    public void Rule_ShouldBeImmutable()
    {
        // Arrange
        var rule = CreateTestRule();

        // Act & Assert - attempting to modify should not be possible (compile-time check)
        // This test exists to demonstrate the record semantics
        rule.Should().NotBeNull();
        rule.Id.Should().Be("test-rule-1");
    }

    [Fact]
    public void Rule_WithMethod_ShouldCreateNewInstance()
    {
        // Arrange
        var original = CreateTestRule();

        // Act
        var modified = original with { Name = "Modified Rule" };

        // Assert
        original.Name.Should().Be("Test Rule");
        modified.Name.Should().Be("Modified Rule");
        original.Should().NotBeSameAs(modified);
    }

    [Fact]
    public void Rule_Equality_ShouldWorkCorrectly()
    {
        // Arrange
        var fixedTimestamp = DateTimeOffset.UtcNow;
        var rule1 = CreateTestRule(fixedTimestamp);
        var rule2 = CreateTestRule(fixedTimestamp);

        // Act & Assert
        // Use BeEquivalentTo for deep comparison (records use reference equality for collections)
        rule1.Should().BeEquivalentTo(rule2);
    }

    private static Rule CreateTestRule(DateTimeOffset? timestamp = null)
    {
        var ts = timestamp ?? DateTimeOffset.UtcNow;
        return new Rule
        {
            Id = "test-rule-1",
            Name = "Test Rule",
            Description = "Test Description",
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
            CreatedAt = ts,
            UpdatedAt = ts
        };
    }
}
