using MongoDB.Bson;
using MongoDB.Bson.Serialization.Attributes;
using Aevum.DecisionEngine.Domain.Enums;

namespace Aevum.DecisionEngine.Infrastructure.Persistence.Models;

public sealed class RuleDocument
{
    [BsonId]
    [BsonRepresentation(BsonType.String)]
    public string Id { get; set; } = string.Empty;

    [BsonElement("name")]
    public string Name { get; set; } = string.Empty;

    [BsonElement("description")]
    public string? Description { get; set; }

    [BsonElement("conditions")]
    public List<RuleConditionDocument> Conditions { get; set; } = [];

    [BsonElement("actions")]
    public List<RuleActionDocument> Actions { get; set; } = [];

    [BsonElement("status")]
    [BsonRepresentation(BsonType.String)]
    public RuleStatus Status { get; set; }

    [BsonElement("version")]
    public int Version { get; set; }

    [BsonElement("priority")]
    public int Priority { get; set; }

    [BsonElement("createdAt")]
    public DateTimeOffset CreatedAt { get; set; }

    [BsonElement("updatedAt")]
    public DateTimeOffset UpdatedAt { get; set; }

    [BsonElement("createdBy")]
    public string? CreatedBy { get; set; }

    [BsonElement("updatedBy")]
    public string? UpdatedBy { get; set; }

    [BsonElement("metadata")]
    public Dictionary<string, string> Metadata { get; set; } = [];

    [BsonElement("effectiveFrom")]
    public DateTimeOffset? EffectiveFrom { get; set; }

    [BsonElement("effectiveUntil")]
    public DateTimeOffset? EffectiveUntil { get; set; }
}

public sealed class RuleConditionDocument
{
    [BsonElement("field")]
    public string Field { get; set; } = string.Empty;

    [BsonElement("operator")]
    [BsonRepresentation(BsonType.String)]
    public ComparisonOperator Operator { get; set; }

    [BsonElement("value")]
    public BsonValue Value { get; set; } = BsonNull.Value;

    [BsonElement("logicalOperator")]
    [BsonRepresentation(BsonType.String)]
    public LogicalOperator? LogicalOperator { get; set; }

    [BsonElement("nestedConditions")]
    public List<RuleConditionDocument> NestedConditions { get; set; } = [];
}

public sealed class RuleActionDocument
{
    [BsonElement("type")]
    [BsonRepresentation(BsonType.String)]
    public ActionType Type { get; set; }

    [BsonElement("parameters")]
    public Dictionary<string, BsonValue> Parameters { get; set; } = [];

    [BsonElement("order")]
    public int Order { get; set; }

    [BsonElement("description")]
    public string? Description { get; set; }
}
