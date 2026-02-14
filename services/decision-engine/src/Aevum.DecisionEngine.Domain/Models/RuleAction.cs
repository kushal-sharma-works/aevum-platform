using Aevum.DecisionEngine.Domain.Enums;

namespace Aevum.DecisionEngine.Domain.Models;

public sealed record RuleAction
{
    public required ActionType Type { get; init; }
    public required IReadOnlyDictionary<string, object> Parameters { get; init; }
    public int Order { get; init; }
    public string? Description { get; init; }
}
