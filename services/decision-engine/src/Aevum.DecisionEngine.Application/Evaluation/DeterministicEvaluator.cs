using System.Security.Cryptography;
using System.Text;
using System.Text.Json;
using System.Text.RegularExpressions;
using Aevum.DecisionEngine.Domain.Enums;
using Aevum.DecisionEngine.Domain.Exceptions;
using Aevum.DecisionEngine.Domain.Interfaces;
using Aevum.DecisionEngine.Domain.Models;

namespace Aevum.DecisionEngine.Application.Evaluation;

public sealed class DeterministicEvaluator(TimeProvider timeProvider) : IDeterministicEvaluator
{
    private readonly TimeProvider _timeProvider = timeProvider;

    public EvaluationResult Evaluate(Rule rule, EvaluationContext context)
    {
        try
        {
            var matchedConditions = new List<string>();
            var isMatch = EvaluateConditions(rule.Conditions, context.Data, matchedConditions, rule.Id);
            
            var deterministicHash = ComputeHash(rule, context);
            var status = isMatch ? DecisionStatus.Approved : DecisionStatus.Rejected;
            var actionsToExecute = isMatch ? rule.Actions.OrderBy(a => a.Order).ToList() : [];

            return new EvaluationResult
            {
                IsMatch = isMatch,
                MatchedConditions = matchedConditions,
                ActionsToExecute = actionsToExecute,
                Status = status,
                DeterministicHash = deterministicHash,
                OutputData = new Dictionary<string, object>()
            };
        }
        catch (Exception ex) when (ex is not EvaluationException)
        {
            throw new EvaluationException($"Evaluation failed for rule '{rule.Id}': {ex.Message}", ex);
        }
    }

    public string ComputeHash(Rule rule, EvaluationContext context)
    {
        var hashInput = new
        {
            RuleId = rule.Id,
            RuleVersion = rule.Version,
            Context = SortDictionary(context.Data),
            Timestamp = context.Timestamp.ToUnixTimeSeconds()
        };

        var json = JsonSerializer.Serialize(hashInput, new JsonSerializerOptions
        {
            WriteIndented = false,
            PropertyNamingPolicy = JsonNamingPolicy.CamelCase
        });

        var bytes = SHA256.HashData(Encoding.UTF8.GetBytes(json));
        return Convert.ToHexString(bytes).ToLowerInvariant();
    }

    private bool EvaluateConditions(
        IReadOnlyList<RuleCondition> conditions,
        IReadOnlyDictionary<string, object> data,
        List<string> matchedConditions,
        string ruleId)
    {
        if (conditions.Count == 0)
            return true;

        var results = new List<bool>();

        foreach (var condition in conditions)
        {
            bool result;

            if (condition.NestedConditions.Count > 0)
            {
                result = EvaluateConditions(condition.NestedConditions, data, matchedConditions, ruleId);
            }
            else
            {
                result = EvaluateCondition(condition, data, ruleId);
                if (result)
                {
                    matchedConditions.Add($"{condition.Field} {condition.Operator} {condition.Value}");
                }
            }

            if (condition.LogicalOperator == LogicalOperator.Not)
            {
                result = !result;
            }

            results.Add(result);

            if (results.Count > 1)
            {
                var prevLogicalOp = conditions[results.Count - 2].LogicalOperator ?? LogicalOperator.And;
                
                if (prevLogicalOp == LogicalOperator.And)
                {
                    var combinedResult = results[^2] && results[^1];
                    results.RemoveRange(results.Count - 2, 2);
                    results.Add(combinedResult);
                }
                else if (prevLogicalOp == LogicalOperator.Or)
                {
                    var combinedResult = results[^2] || results[^1];
                    results.RemoveRange(results.Count - 2, 2);
                    results.Add(combinedResult);
                }
            }
        }

        return results.Count > 0 && results[0];
    }

    private bool EvaluateCondition(
        RuleCondition condition,
        IReadOnlyDictionary<string, object> data,
        string ruleId)
    {
        if (!data.TryGetValue(condition.Field, out var fieldValue))
        {
            throw new EvaluationException(ruleId, condition.Field, "Field not found in context");
        }

        try
        {
            return condition.Operator switch
            {
                ComparisonOperator.Equals => AreEqual(fieldValue, condition.Value),
                ComparisonOperator.NotEquals => !AreEqual(fieldValue, condition.Value),
                ComparisonOperator.GreaterThan => CompareValues(fieldValue, condition.Value) > 0,
                ComparisonOperator.GreaterThanOrEqual => CompareValues(fieldValue, condition.Value) >= 0,
                ComparisonOperator.LessThan => CompareValues(fieldValue, condition.Value) < 0,
                ComparisonOperator.LessThanOrEqual => CompareValues(fieldValue, condition.Value) <= 0,
                ComparisonOperator.Contains => ContainsValue(fieldValue, condition.Value),
                ComparisonOperator.NotContains => !ContainsValue(fieldValue, condition.Value),
                ComparisonOperator.StartsWith => StartsWithValue(fieldValue, condition.Value),
                ComparisonOperator.EndsWith => EndsWithValue(fieldValue, condition.Value),
                ComparisonOperator.In => InCollection(fieldValue, condition.Value),
                ComparisonOperator.NotIn => !InCollection(fieldValue, condition.Value),
                ComparisonOperator.Regex => MatchesRegex(fieldValue, condition.Value),
                _ => throw new EvaluationException(ruleId, condition.Field, $"Unsupported operator: {condition.Operator}")
            };
        }
        catch (Exception ex) when (ex is not EvaluationException)
        {
            throw new EvaluationException(ruleId, condition.Field, $"Comparison failed: {ex.Message}", ex);
        }
    }

    private static bool AreEqual(object fieldValue, object conditionValue)
    {
        var fieldStr = ConvertToString(fieldValue);
        var conditionStr = ConvertToString(conditionValue);
        return string.Equals(fieldStr, conditionStr, StringComparison.OrdinalIgnoreCase);
    }

    private static int CompareValues(object fieldValue, object conditionValue)
    {
        if (TryConvertToDecimal(fieldValue, out var fieldDecimal) &&
            TryConvertToDecimal(conditionValue, out var conditionDecimal))
        {
            return fieldDecimal.CompareTo(conditionDecimal);
        }

        if (TryConvertToDateTime(fieldValue, out var fieldDateTime) &&
            TryConvertToDateTime(conditionValue, out var conditionDateTime))
        {
            return fieldDateTime.CompareTo(conditionDateTime);
        }

        var fieldStr = ConvertToString(fieldValue);
        var conditionStr = ConvertToString(conditionValue);
        return string.Compare(fieldStr, conditionStr, StringComparison.OrdinalIgnoreCase);
    }

    private static bool ContainsValue(object fieldValue, object conditionValue)
    {
        var fieldStr = ConvertToString(fieldValue);
        var conditionStr = ConvertToString(conditionValue);
        return fieldStr.Contains(conditionStr, StringComparison.OrdinalIgnoreCase);
    }

    private static bool StartsWithValue(object fieldValue, object conditionValue)
    {
        var fieldStr = ConvertToString(fieldValue);
        var conditionStr = ConvertToString(conditionValue);
        return fieldStr.StartsWith(conditionStr, StringComparison.OrdinalIgnoreCase);
    }

    private static bool EndsWithValue(object fieldValue, object conditionValue)
    {
        var fieldStr = ConvertToString(fieldValue);
        var conditionStr = ConvertToString(conditionValue);
        return fieldStr.EndsWith(conditionStr, StringComparison.OrdinalIgnoreCase);
    }

    private static bool InCollection(object fieldValue, object conditionValue)
    {
        if (conditionValue is not System.Collections.IEnumerable enumerable)
            return false;

        var fieldStr = ConvertToString(fieldValue);
        foreach (var item in enumerable)
        {
            if (string.Equals(fieldStr, ConvertToString(item), StringComparison.OrdinalIgnoreCase))
                return true;
        }

        return false;
    }

    private static bool MatchesRegex(object fieldValue, object conditionValue)
    {
        var fieldStr = ConvertToString(fieldValue);
        var pattern = ConvertToString(conditionValue);
        return Regex.IsMatch(fieldStr, pattern, RegexOptions.IgnoreCase, TimeSpan.FromSeconds(1));
    }

    private static string ConvertToString(object value)
    {
        return value switch
        {
            null => string.Empty,
            string s => s,
            JsonElement element => element.ToString(),
            _ => value.ToString() ?? string.Empty
        };
    }

    private static bool TryConvertToDecimal(object value, out decimal result)
    {
        result = 0;
        
        if (value is decimal d)
        {
            result = d;
            return true;
        }

        if (value is int i)
        {
            result = i;
            return true;
        }

        if (value is long l)
        {
            result = l;
            return true;
        }

        if (value is double dbl)
        {
            result = (decimal)dbl;
            return true;
        }

        if (value is JsonElement element && element.ValueKind == JsonValueKind.Number)
        {
            return element.TryGetDecimal(out result);
        }

        return decimal.TryParse(ConvertToString(value), out result);
    }

    private static bool TryConvertToDateTime(object value, out DateTimeOffset result)
    {
        result = default;

        if (value is DateTimeOffset dto)
        {
            result = dto;
            return true;
        }

        if (value is DateTime dt)
        {
            result = new DateTimeOffset(dt);
            return true;
        }

        if (value is JsonElement element && element.ValueKind == JsonValueKind.String)
        {
            return DateTimeOffset.TryParse(element.GetString(), out result);
        }

        return DateTimeOffset.TryParse(ConvertToString(value), out result);
    }

    private static Dictionary<string, object> SortDictionary(IReadOnlyDictionary<string, object> dict)
    {
        var sorted = new Dictionary<string, object>();
        foreach (var key in dict.Keys.OrderBy(k => k))
        {
            sorted[key] = dict[key];
        }
        return sorted;
    }
}
