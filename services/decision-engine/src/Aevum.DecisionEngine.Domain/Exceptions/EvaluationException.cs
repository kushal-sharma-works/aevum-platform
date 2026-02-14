namespace Aevum.DecisionEngine.Domain.Exceptions;

public sealed class EvaluationException(string message, Exception? innerException = null) 
    : DomainException(message, innerException)
{
    public EvaluationException(string ruleId, string fieldName, string details, Exception? innerException = null)
        : this($"Evaluation failed for rule '{ruleId}' on field '{fieldName}': {details}", innerException)
    {
        RuleId = ruleId;
        FieldName = fieldName;
    }

    public string? RuleId { get; }
    public string? FieldName { get; }
}
