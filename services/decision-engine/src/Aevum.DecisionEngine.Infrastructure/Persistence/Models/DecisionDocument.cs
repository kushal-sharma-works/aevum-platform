using MongoDB.Bson;
using MongoDB.Bson.Serialization.Attributes;
using Aevum.DecisionEngine.Domain.Enums;

namespace Aevum.DecisionEngine.Infrastructure.Persistence.Models;

public sealed class DecisionDocument
{
    [BsonId]
    [BsonRepresentation(BsonType.String)]
    public string Id { get; set; } = string.Empty;

    [BsonElement("ruleId")]
    public string RuleId { get; set; } = string.Empty;

    [BsonElement("ruleVersion")]
    public int RuleVersion { get; set; }

    [BsonElement("requestId")]
    public string RequestId { get; set; } = string.Empty;

    [BsonElement("status")]
    [BsonRepresentation(BsonType.String)]
    public DecisionStatus Status { get; set; }

    [BsonElement("inputContext")]
    public Dictionary<string, BsonValue> InputContext { get; set; } = [];

    [BsonElement("matchedConditions")]
    public List<string> MatchedConditions { get; set; } = [];

    [BsonElement("executedActions")]
    public List<RuleActionDocument> ExecutedActions { get; set; } = [];

    [BsonElement("evaluatedAt")]
    public DateTimeOffset EvaluatedAt { get; set; }

    [BsonElement("deterministicHash")]
    public string DeterministicHash { get; set; } = string.Empty;

    [BsonElement("errorMessage")]
    public string? ErrorMessage { get; set; }

    [BsonElement("outputData")]
    public Dictionary<string, BsonValue> OutputData { get; set; } = [];

    [BsonElement("evaluationDurationMs")]
    public long EvaluationDurationMs { get; set; }
}
