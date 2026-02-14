using Aevum.DecisionEngine.Domain.Enums;

namespace Aevum.DecisionEngine.Application.DTOs;

public sealed record CreateRuleRequest
{
    public required string Name { get; init; }
    public string? Description { get; init; }
    public required IReadOnlyList<RuleConditionDto> Conditions { get; init; }
    public required IReadOnlyList<RuleActionDto> Actions { get; init; }
    public required int Priority { get; init; }
    public string? CreatedBy { get; init; }
    public Dictionary<string, string>? Metadata { get; init; }
    public DateTimeOffset? EffectiveFrom { get; init; }
    public DateTimeOffset? EffectiveUntil { get; init; }
}

public sealed record UpdateRuleRequest
{
    public required string Name { get; init; }
    public string? Description { get; init; }
    public required IReadOnlyList<RuleConditionDto> Conditions { get; init; }
    public required IReadOnlyList<RuleActionDto> Actions { get; init; }
    public required int Priority { get; init; }
    public RuleStatus? Status { get; init; }
    public string? UpdatedBy { get; init; }
    public Dictionary<string, string>? Metadata { get; init; }
    public DateTimeOffset? EffectiveFrom { get; init; }
    public DateTimeOffset? EffectiveUntil { get; init; }
}

public sealed record RuleConditionDto
{
    public required string Field { get; init; }
    public required ComparisonOperator Operator { get; init; }
    public required object Value { get; init; }
    public LogicalOperator? LogicalOperator { get; init; }
    public IReadOnlyList<RuleConditionDto>? NestedConditions { get; init; }
}

public sealed record RuleActionDto
{
    public required ActionType Type { get; init; }
    public required Dictionary<string, object> Parameters { get; init; }
    public int Order { get; init; }
    public string? Description { get; init; }
}

public sealed record RuleResponse
{
    public required string Id { get; init; }
    public required string Name { get; init; }
    public string? Description { get; init; }
    public required IReadOnlyList<RuleConditionDto> Conditions { get; init; }
    public required IReadOnlyList<RuleActionDto> Actions { get; init; }
    public required RuleStatus Status { get; init; }
    public required int Version { get; init; }
    public required int Priority { get; init; }
    public required DateTimeOffset CreatedAt { get; init; }
    public required DateTimeOffset UpdatedAt { get; init; }
    public string? CreatedBy { get; init; }
    public string? UpdatedBy { get; init; }
    public Dictionary<string, string>? Metadata { get; init; }
    public DateTimeOffset? EffectiveFrom { get; init; }
    public DateTimeOffset? EffectiveUntil { get; init; }
}
