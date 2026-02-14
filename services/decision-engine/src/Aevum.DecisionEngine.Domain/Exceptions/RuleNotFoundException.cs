namespace Aevum.DecisionEngine.Domain.Exceptions;

public sealed class RuleNotFoundException(string ruleId, int? version = null) 
    : DomainException($"Rule with ID '{ruleId}'{(version.HasValue ? $" version {version}" : "")} was not found")
{
    public string RuleId { get; } = ruleId;
    public int? Version { get; } = version;
}
