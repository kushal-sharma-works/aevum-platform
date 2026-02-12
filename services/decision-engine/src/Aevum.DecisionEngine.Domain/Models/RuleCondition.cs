using Aevum.DecisionEngine.Domain.Enums;

namespace Aevum.DecisionEngine.Domain.Models;

public sealed record RuleCondition
{
    public required string Field { get; init; }
    public required ComparisonOperator Operator { get; init; }
    public required object Value { get; init; }
    public LogicalOperator? LogicalOperator { get; init; }
    public IReadOnlyList<RuleCondition> NestedConditions { get; init; } = [];
}
