using Aevum.DecisionEngine.Domain.Enums;

namespace Aevum.DecisionEngine.Domain.Models;

public sealed record Decision
{
    public required string Id { get; init; }
    public required string RuleId { get; init; }
    public required int RuleVersion { get; init; }
    public required string RequestId { get; init; }
    public required DecisionStatus Status { get; init; }
    public required IReadOnlyDictionary<string, object> InputContext { get; init; }
    public required IReadOnlyList<string> MatchedConditions { get; init; }
    public required IReadOnlyList<RuleAction> ExecutedActions { get; init; }
    public required DateTimeOffset EvaluatedAt { get; init; }
    public required string DeterministicHash { get; init; }
    public string? ErrorMessage { get; init; }
    public IReadOnlyDictionary<string, object> OutputData { get; init; } = new Dictionary<string, object>();
    public long EvaluationDurationMs { get; init; }
}
