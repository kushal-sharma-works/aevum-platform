using Aevum.DecisionEngine.Domain.Models;

namespace Aevum.DecisionEngine.Domain.Interfaces;

public interface IDecisionRepository
{
    Task<Decision?> GetByIdAsync(string id, CancellationToken cancellationToken = default);
    Task<Decision?> GetByRequestIdAsync(string requestId, CancellationToken cancellationToken = default);
    Task<IReadOnlyList<Decision>> GetByRuleIdAsync(string ruleId, int? version = null, CancellationToken cancellationToken = default);
    Task<Decision> CreateAsync(Decision decision, CancellationToken cancellationToken = default);
}
