namespace Aevum.DecisionEngine.Domain.Exceptions;

public class DomainException(string message, Exception? innerException = null) 
    : Exception(message, innerException)
{
    public string ErrorCode { get; init; } = "DOMAIN_ERROR";
}
