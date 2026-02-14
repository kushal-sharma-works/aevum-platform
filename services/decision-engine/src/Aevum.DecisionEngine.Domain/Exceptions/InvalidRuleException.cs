namespace Aevum.DecisionEngine.Domain.Exceptions;

public sealed class InvalidRuleException(string message) 
    : DomainException(message)
{
    public IReadOnlyList<string> ValidationErrors { get; init; } = [];
}
