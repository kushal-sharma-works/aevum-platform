using Aevum.DecisionEngine.Domain.Enums;

namespace Aevum.DecisionEngine.Domain.Models;

public sealed record EvaluationResult
{
    public required bool IsMatch { get; init; }
    public required IReadOnlyList<string> MatchedConditions { get; init; }
    public required IReadOnlyList<RuleAction> ActionsToExecute { get; init; }
    public required DecisionStatus Status { get; init; }
    public required string DeterministicHash { get; init; }
    public string? ErrorMessage { get; init; }
    public IReadOnlyDictionary<string, object> OutputData { get; init; } = new Dictionary<string, object>();
}
