namespace Aevum.DecisionEngine.Domain.Models;

public sealed record EvaluationContext
{
    public required IReadOnlyDictionary<string, object> Data { get; init; }
    public required string RequestId { get; init; }
    public required DateTimeOffset Timestamp { get; init; }
    public IReadOnlyDictionary<string, string> Metadata { get; init; } = new Dictionary<string, string>();
}
