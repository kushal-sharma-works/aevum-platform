using Aevum.DecisionEngine.Domain.Enums;
using Aevum.DecisionEngine.Domain.Models;

namespace Aevum.DecisionEngine.Domain.Interfaces;

public interface IRuleRepository
{
    Task<Rule?> GetByIdAsync(string id, int? version = null, CancellationToken cancellationToken = default);
    Task<IReadOnlyList<Rule>> GetActiveRulesAsync(CancellationToken cancellationToken = default);
    Task<IReadOnlyList<Rule>> GetByStatusAsync(RuleStatus status, CancellationToken cancellationToken = default);
    Task<Rule> CreateAsync(Rule rule, CancellationToken cancellationToken = default);
    Task<Rule> UpdateAsync(Rule rule, CancellationToken cancellationToken = default);
    Task DeleteAsync(string id, CancellationToken cancellationToken = default);
    Task<int> GetLatestVersionAsync(string id, CancellationToken cancellationToken = default);
}
