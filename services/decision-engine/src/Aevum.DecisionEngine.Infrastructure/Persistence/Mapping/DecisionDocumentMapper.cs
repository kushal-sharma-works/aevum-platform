using MongoDB.Bson;
using System.Text.Json;
using Aevum.DecisionEngine.Domain.Models;
using Aevum.DecisionEngine.Infrastructure.Persistence.Models;

namespace Aevum.DecisionEngine.Infrastructure.Persistence.Mapping;

public static class DecisionDocumentMapper
{
    public static DecisionDocument ToDocument(this Decision decision)
    {
        return new DecisionDocument
        {
            Id = decision.Id,
            RuleId = decision.RuleId,
            RuleVersion = decision.RuleVersion,
            RequestId = decision.RequestId,
            Status = decision.Status,
            InputContext = decision.InputContext.ToDictionary(kvp => kvp.Key, kvp => BsonValue.Create(NormalizeValue(kvp.Value))),
            MatchedConditions = decision.MatchedConditions.ToList(),
            ExecutedActions = decision.ExecutedActions.Select(a => a.ToDocument()).ToList(),
            EvaluatedAt = decision.EvaluatedAt,
            DeterministicHash = decision.DeterministicHash,
            ErrorMessage = decision.ErrorMessage,
            OutputData = decision.OutputData.ToDictionary(kvp => kvp.Key, kvp => BsonValue.Create(NormalizeValue(kvp.Value))),
            EvaluationDurationMs = decision.EvaluationDurationMs
        };
    }

    private static object? NormalizeValue(object? value)
    {
        if (value is null)
        {
            return null;
        }

        if (value is JsonElement jsonElement)
        {
            return JsonElementToObject(jsonElement);
        }

        return value;
    }

    private static object? JsonElementToObject(JsonElement element)
    {
        return element.ValueKind switch
        {
            JsonValueKind.String => element.GetString(),
            JsonValueKind.Number when element.TryGetInt32(out var intValue) => intValue,
            JsonValueKind.Number when element.TryGetInt64(out var longValue) => longValue,
            JsonValueKind.Number when element.TryGetDecimal(out var decimalValue) => decimalValue,
            JsonValueKind.Number => element.GetDouble(),
            JsonValueKind.True => true,
            JsonValueKind.False => false,
            JsonValueKind.Array => element.EnumerateArray().Select(JsonElementToObject).ToList(),
            JsonValueKind.Object => element.EnumerateObject().ToDictionary(p => p.Name, p => JsonElementToObject(p.Value)),
            JsonValueKind.Null => null,
            _ => element.GetRawText()
        };
    }

    public static Decision ToDomain(this DecisionDocument doc)
    {
        return new Decision
        {
            Id = doc.Id,
            RuleId = doc.RuleId,
            RuleVersion = doc.RuleVersion,
            RequestId = doc.RequestId,
            Status = doc.Status,
            InputContext = doc.InputContext.ToDictionary(kvp => kvp.Key, kvp => BsonValueToObject(kvp.Value)),
            MatchedConditions = doc.MatchedConditions,
            ExecutedActions = doc.ExecutedActions.Select(a => a.ToDomain()).ToList(),
            EvaluatedAt = doc.EvaluatedAt,
            DeterministicHash = doc.DeterministicHash,
            ErrorMessage = doc.ErrorMessage,
            OutputData = doc.OutputData.ToDictionary(kvp => kvp.Key, kvp => BsonValueToObject(kvp.Value)),
            EvaluationDurationMs = doc.EvaluationDurationMs
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
