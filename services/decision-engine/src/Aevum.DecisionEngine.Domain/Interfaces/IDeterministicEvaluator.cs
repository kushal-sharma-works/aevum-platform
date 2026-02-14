using Aevum.DecisionEngine.Domain.Models;

namespace Aevum.DecisionEngine.Domain.Interfaces;

public interface IDeterministicEvaluator
{
    EvaluationResult Evaluate(Rule rule, EvaluationContext context);
    string ComputeHash(Rule rule, EvaluationContext context);
}
