using Aevum.DecisionEngine.Domain.Enums;
using Aevum.DecisionEngine.Domain.Exceptions;
using Aevum.DecisionEngine.Domain.Interfaces;
using Aevum.DecisionEngine.Domain.Models;

namespace Aevum.DecisionEngine.Application.Services;

public sealed class RuleManagementService(
    IRuleRepository ruleRepository,
    TimeProvider timeProvider)
{
    private readonly IRuleRepository _ruleRepository = ruleRepository;
    private readonly TimeProvider _timeProvider = timeProvider;

    public async Task<Rule> CreateRuleAsync(Rule rule, CancellationToken cancellationToken = default)
    {
        var now = _timeProvider.GetUtcNow();
        var ruleToCreate = rule with
        {
            Id = string.IsNullOrEmpty(rule.Id) ? Guid.NewGuid().ToString() : rule.Id,
            Version = 1,
            CreatedAt = now,
            UpdatedAt = now,
            Status = RuleStatus.Draft
        };

        return await _ruleRepository.CreateAsync(ruleToCreate, cancellationToken);
    }

    public async Task<Rule> UpdateRuleAsync(string id, Rule updates, CancellationToken cancellationToken = default)
    {
        var existingRule = await _ruleRepository.GetByIdAsync(id, version: null, cancellationToken)
            ?? throw new RuleNotFoundException(id);

        var latestVersion = await _ruleRepository.GetLatestVersionAsync(id, cancellationToken);
        var now = _timeProvider.GetUtcNow();

        var updatedRule = updates with
        {
            Id = id,
            Version = latestVersion + 1,
            CreatedAt = existingRule.CreatedAt,
            CreatedBy = existingRule.CreatedBy,
            UpdatedAt = now
        };

        return await _ruleRepository.CreateAsync(updatedRule, cancellationToken);
    }

    public async Task<Rule> GetRuleAsync(string id, int? version = null, CancellationToken cancellationToken = default)
    {
        var rule = await _ruleRepository.GetByIdAsync(id, version, cancellationToken);
        return rule ?? throw new RuleNotFoundException(id, version);
    }

    public async Task<IReadOnlyList<Rule>> GetActiveRulesAsync(CancellationToken cancellationToken = default)
    {
        return await _ruleRepository.GetActiveRulesAsync(cancellationToken);
    }

    public async Task<IReadOnlyList<Rule>> GetRulesByStatusAsync(RuleStatus status, CancellationToken cancellationToken = default)
    {
        return await _ruleRepository.GetByStatusAsync(status, cancellationToken);
    }

    public async Task DeleteRuleAsync(string id, CancellationToken cancellationToken = default)
    {
        var exists = await _ruleRepository.GetByIdAsync(id, version: null, cancellationToken);
        if (exists is null)
        {
            throw new RuleNotFoundException(id);
        }

        await _ruleRepository.DeleteAsync(id, cancellationToken);
    }

    public async Task<Rule> ActivateRuleAsync(string id, CancellationToken cancellationToken = default)
    {
        var rule = await GetRuleAsync(id, version: null, cancellationToken);
        var activated = rule with
        {
            Status = RuleStatus.Active,
            UpdatedAt = _timeProvider.GetUtcNow()
        };

        return await _ruleRepository.UpdateAsync(activated, cancellationToken);
    }

    public async Task<Rule> DeactivateRuleAsync(string id, CancellationToken cancellationToken = default)
    {
        var rule = await GetRuleAsync(id, version: null, cancellationToken);
        var deactivated = rule with
        {
            Status = RuleStatus.Inactive,
            UpdatedAt = _timeProvider.GetUtcNow()
        };

        return await _ruleRepository.UpdateAsync(deactivated, cancellationToken);
    }
}
