using Aevum.DecisionEngine.Application.DTOs;
using Aevum.DecisionEngine.Domain.Models;

namespace Aevum.DecisionEngine.Application.Mapping;

public static class RuleMapper
{
    public static Rule ToDomain(this CreateRuleRequest request)
    {
        return new Rule
        {
            Id = string.Empty, // Will be set by service
            Name = request.Name,
            Description = request.Description,
            Conditions = request.Conditions.Select(c => c.ToDomain()).ToList(),
            Actions = request.Actions.Select(a => a.ToDomain()).ToList(),
            Status = Domain.Enums.RuleStatus.Draft,
            Version = 1,
            Priority = request.Priority,
            CreatedAt = default, // Will be set by service
            UpdatedAt = default, // Will be set by service
            CreatedBy = request.CreatedBy,
            Metadata = request.Metadata ?? new Dictionary<string, string>(),
            EffectiveFrom = request.EffectiveFrom,
            EffectiveUntil = request.EffectiveUntil
        };
    }

    public static Rule ToDomain(this UpdateRuleRequest request, string id)
    {
        return new Rule
        {
            Id = id,
            Name = request.Name,
            Description = request.Description,
            Conditions = request.Conditions.Select(c => c.ToDomain()).ToList(),
            Actions = request.Actions.Select(a => a.ToDomain()).ToList(),
            Status = request.Status ?? Domain.Enums.RuleStatus.Draft,
            Version = 0, // Will be set by service
            Priority = request.Priority,
            CreatedAt = default, // Will be set by service
            UpdatedAt = default, // Will be set by service
            UpdatedBy = request.UpdatedBy,
            Metadata = request.Metadata ?? new Dictionary<string, string>(),
            EffectiveFrom = request.EffectiveFrom,
            EffectiveUntil = request.EffectiveUntil
        };
    }

    public static RuleResponse ToResponse(this Rule rule)
    {
        return new RuleResponse
        {
            Id = rule.Id,
            Name = rule.Name,
            Description = rule.Description,
            Conditions = rule.Conditions.Select(c => c.ToDto()).ToList(),
            Actions = rule.Actions.Select(a => a.ToDto()).ToList(),
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

    public static RuleCondition ToDomain(this RuleConditionDto dto)
    {
        return new RuleCondition
        {
            Field = dto.Field,
            Operator = dto.Operator,
            Value = dto.Value,
            LogicalOperator = dto.LogicalOperator,
            NestedConditions = dto.NestedConditions?.Select(nc => nc.ToDomain()).ToList() ?? []
        };
    }

    public static RuleConditionDto ToDto(this RuleCondition condition)
    {
        return new RuleConditionDto
        {
            Field = condition.Field,
            Operator = condition.Operator,
            Value = condition.Value,
            LogicalOperator = condition.LogicalOperator,
            NestedConditions = condition.NestedConditions.Select(nc => nc.ToDto()).ToList()
        };
    }

    public static RuleAction ToDomain(this RuleActionDto dto)
    {
        return new RuleAction
        {
            Type = dto.Type,
            Parameters = dto.Parameters,
            Order = dto.Order,
            Description = dto.Description
        };
    }

    public static RuleActionDto ToDto(this RuleAction action)
    {
        return new RuleActionDto
        {
            Type = action.Type,
            Parameters = action.Parameters.ToDictionary(kvp => kvp.Key, kvp => kvp.Value),
            Order = action.Order,
            Description = action.Description
        };
    }
}
