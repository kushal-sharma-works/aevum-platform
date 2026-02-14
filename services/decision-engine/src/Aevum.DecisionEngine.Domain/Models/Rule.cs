using Aevum.DecisionEngine.Domain.Enums;

namespace Aevum.DecisionEngine.Domain.Models;

public sealed record Rule
{
    public required string Id { get; init; }
    public required string Name { get; init; }
    public string? Description { get; init; }
    public required IReadOnlyList<RuleCondition> Conditions { get; init; }
    public required IReadOnlyList<RuleAction> Actions { get; init; }
    public required RuleStatus Status { get; init; }
    public required int Version { get; init; }
    public required int Priority { get; init; }
    public required DateTimeOffset CreatedAt { get; init; }
    public required DateTimeOffset UpdatedAt { get; init; }
    public string? CreatedBy { get; init; }
    public string? UpdatedBy { get; init; }
    public IReadOnlyDictionary<string, string> Metadata { get; init; } = new Dictionary<string, string>();
    public DateTimeOffset? EffectiveFrom { get; init; }
    public DateTimeOffset? EffectiveUntil { get; init; }
}
