using Aevum.DecisionEngine.Domain.Enums;

namespace Aevum.DecisionEngine.Application.DTOs;

public sealed record EvaluateDecisionRequest
{
    public required string RuleId { get; init; }
    public int? RuleVersion { get; init; }
    public required Dictionary<string, object> Context { get; init; }
    public required string RequestId { get; init; }
    public Dictionary<string, string>? Metadata { get; init; }
}

public sealed record DecisionResponse
{
    public required string Id { get; init; }
    public required string RuleId { get; init; }
    public required int RuleVersion { get; init; }
    public required string RequestId { get; init; }
    public required DecisionStatus Status { get; init; }
    public required Dictionary<string, object> InputContext { get; init; }
    public required IReadOnlyList<string> MatchedConditions { get; init; }
    public required IReadOnlyList<RuleActionDto> ExecutedActions { get; init; }
    public required DateTimeOffset EvaluatedAt { get; init; }
    public required string DeterministicHash { get; init; }
    public string? ErrorMessage { get; init; }
    public Dictionary<string, object>? OutputData { get; init; }
    public required long EvaluationDurationMs { get; init; }
}
