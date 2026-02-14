using FluentAssertions;
using Aevum.DecisionEngine.Domain.Enums;
using Aevum.DecisionEngine.Domain.Models;

namespace Aevum.DecisionEngine.Domain.Tests.Models;

public sealed class DecisionTests
{
    [Fact]
    public void Decision_ShouldBeImmutable()
    {
        // Arrange
        var decision = CreateTestDecision();

        // Act & Assert
        decision.Should().NotBeNull();
        decision.Id.Should().Be("test-decision-1");
    }

    [Fact]
    public void Decision_WithMethod_ShouldCreateNewInstance()
    {
        // Arrange
        var original = CreateTestDecision();

        // Act
        var modified = original with { Status = DecisionStatus.Rejected };

        // Assert
        original.Status.Should().Be(DecisionStatus.Approved);
        modified.Status.Should().Be(DecisionStatus.Rejected);
        original.Should().NotBeSameAs(modified);
    }

    private static Decision CreateTestDecision()
    {
        return new Decision
        {
            Id = "test-decision-1",
            RuleId = "test-rule-1",
            RuleVersion = 1,
            RequestId = "req-123",
            Status = DecisionStatus.Approved,
            InputContext = new Dictionary<string, object> { ["amount"] = 150 },
            MatchedConditions = ["amount > 100"],
            ExecutedActions = [],
            EvaluatedAt = DateTimeOffset.UtcNow,
            DeterministicHash = "abc123",
            EvaluationDurationMs = 10
        };
    }
}
