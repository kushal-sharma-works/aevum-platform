using MongoDB.Bson;
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
            InputContext = decision.InputContext.ToDictionary(kvp => kvp.Key, kvp => BsonValue.Create(kvp.Value)),
            MatchedConditions = decision.MatchedConditions.ToList(),
            ExecutedActions = decision.ExecutedActions.Select(a => a.ToDocument()).ToList(),
            EvaluatedAt = decision.EvaluatedAt,
            DeterministicHash = decision.DeterministicHash,
            ErrorMessage = decision.ErrorMessage,
            OutputData = decision.OutputData.ToDictionary(kvp => kvp.Key, kvp => BsonValue.Create(kvp.Value)),
            EvaluationDurationMs = decision.EvaluationDurationMs
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
