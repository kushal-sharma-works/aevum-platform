using MongoDB.Bson;
using Aevum.DecisionEngine.Domain.Models;
using Aevum.DecisionEngine.Infrastructure.Persistence.Models;

namespace Aevum.DecisionEngine.Infrastructure.Persistence.Mapping;

public static class RuleDocumentMapper
{
    public static RuleDocument ToDocument(this Rule rule)
    {
        return new RuleDocument
        {
            Id = rule.Id,
            Name = rule.Name,
            Description = rule.Description,
            Conditions = rule.Conditions.Select(c => c.ToDocument()).ToList(),
            Actions = rule.Actions.Select(a => a.ToDocument()).ToList(),
            Status = rule.Status,
            Version = rule.Version,
            Priority = rule.Priority,
            CreatedAt = rule.CreatedAt,
            UpdatedAt = rule.UpdatedAt,
            CreatedBy = rule.CreatedBy,
            UpdatedBy = rule.UpdatedBy,
            Metadata = rule.Metadata.ToDictionary(kvp => kvp.Key, kvp => kvp.Value),
            EffectiveFrom = rule.EffectiveFrom,
            EffectiveUntil = rule.EffectiveUntil
        };
    }

    public static Rule ToDomain(this RuleDocument doc)
    {
        return new Rule
        {
            Id = doc.Id,
            Name = doc.Name,
            Description = doc.Description,
            Conditions = doc.Conditions.Select(c => c.ToDomain()).ToList(),
            Actions = doc.Actions.Select(a => a.ToDomain()).ToList(),
            Status = doc.Status,
            Version = doc.Version,
            Priority = doc.Priority,
            CreatedAt = doc.CreatedAt,
            UpdatedAt = doc.UpdatedAt,
            CreatedBy = doc.CreatedBy,
            UpdatedBy = doc.UpdatedBy,
            Metadata = doc.Metadata,
            EffectiveFrom = doc.EffectiveFrom,
            EffectiveUntil = doc.EffectiveUntil
        };
    }

    public static RuleConditionDocument ToDocument(this RuleCondition condition)
    {
        return new RuleConditionDocument
        {
            Field = condition.Field,
            Operator = condition.Operator,
            Value = BsonValue.Create(condition.Value),
            LogicalOperator = condition.LogicalOperator,
            NestedConditions = condition.NestedConditions.Select(nc => nc.ToDocument()).ToList()
        };
    }

    public static RuleCondition ToDomain(this RuleConditionDocument doc)
    {
        return new RuleCondition
        {
            Field = doc.Field,
            Operator = doc.Operator,
            Value = BsonValueToObject(doc.Value),
            LogicalOperator = doc.LogicalOperator,
            NestedConditions = doc.NestedConditions.Select(nc => nc.ToDomain()).ToList()
        };
    }

    public static RuleActionDocument ToDocument(this RuleAction action)
    {
        return new RuleActionDocument
        {
            Type = action.Type,
            Parameters = action.Parameters.ToDictionary(kvp => kvp.Key, kvp => BsonValue.Create(kvp.Value)),
            Order = action.Order,
            Description = action.Description
        };
    }

    public static RuleAction ToDomain(this RuleActionDocument doc)
    {
        return new RuleAction
        {
            Type = doc.Type,
            Parameters = doc.Parameters.ToDictionary(kvp => kvp.Key, kvp => BsonValueToObject(kvp.Value)),
            Order = doc.Order,
            Description = doc.Description
        };
    }

    private static object BsonValueToObject(BsonValue value)
    {
        return value.BsonType switch
        {
            BsonType.String => value.AsString,
            BsonType.Int32 => value.AsInt32,
            BsonType.Int64 => value.AsInt64,
            BsonType.Double => value.AsDouble,
            BsonType.Boolean => value.AsBoolean,
            BsonType.DateTime => value.ToUniversalTime(),
            BsonType.Array => value.AsBsonArray.Select(BsonValueToObject).ToList(),
            BsonType.Document => value.AsBsonDocument.ToDictionary(e => e.Name, e => BsonValueToObject(e.Value)),
            BsonType.Null => null!,
            _ => value.ToString()!
        };
    }
}
