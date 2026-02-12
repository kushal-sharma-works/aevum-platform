using FluentValidation;
using Aevum.DecisionEngine.Application.DTOs;

namespace Aevum.DecisionEngine.Application.Validators;

public sealed class EvaluateDecisionRequestValidator : AbstractValidator<EvaluateDecisionRequest>
{
    public EvaluateDecisionRequestValidator()
    {
        RuleFor(x => x.RuleId)
            .NotEmpty().WithMessage("RuleId is required");

        RuleFor(x => x.RuleVersion)
            .GreaterThan(0).WithMessage("RuleVersion must be positive")
            .When(x => x.RuleVersion.HasValue);

        RuleFor(x => x.Context)
            .NotEmpty().WithMessage("Context is required")
            .Must(c => c.Count <= 100).WithMessage("Context cannot exceed 100 fields");

        RuleFor(x => x.RequestId)
            .NotEmpty().WithMessage("RequestId is required")
            .MaximumLength(100).WithMessage("RequestId must not exceed 100 characters");
    }
}
