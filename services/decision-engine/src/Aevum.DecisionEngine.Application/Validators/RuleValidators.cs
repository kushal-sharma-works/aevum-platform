using FluentValidation;
using Aevum.DecisionEngine.Application.DTOs;

namespace Aevum.DecisionEngine.Application.Validators;

public sealed class CreateRuleRequestValidator : AbstractValidator<CreateRuleRequest>
{
    public CreateRuleRequestValidator()
    {
        RuleFor(x => x.Name)
            .NotEmpty().WithMessage("Rule name is required")
            .MaximumLength(200).WithMessage("Rule name must not exceed 200 characters");

        RuleFor(x => x.Description)
            .MaximumLength(1000).WithMessage("Description must not exceed 1000 characters")
            .When(x => !string.IsNullOrEmpty(x.Description));

        RuleFor(x => x.Conditions)
            .NotEmpty().WithMessage("At least one condition is required")
            .Must(c => c.Count <= 50).WithMessage("Cannot exceed 50 conditions");

        RuleForEach(x => x.Conditions)
            .SetValidator(new RuleConditionDtoValidator());

        RuleFor(x => x.Actions)
            .NotEmpty().WithMessage("At least one action is required")
            .Must(a => a.Count <= 20).WithMessage("Cannot exceed 20 actions");

        RuleForEach(x => x.Actions)
            .SetValidator(new RuleActionDtoValidator());

        RuleFor(x => x.Priority)
            .GreaterThanOrEqualTo(0).WithMessage("Priority must be non-negative")
            .LessThanOrEqualTo(1000).WithMessage("Priority must not exceed 1000");

        RuleFor(x => x.EffectiveUntil)
            .Must((req, until) => !until.HasValue || !req.EffectiveFrom.HasValue || until > req.EffectiveFrom)
            .WithMessage("EffectiveUntil must be after EffectiveFrom")
            .When(x => x.EffectiveFrom.HasValue && x.EffectiveUntil.HasValue);
    }
}

public sealed class UpdateRuleRequestValidator : AbstractValidator<UpdateRuleRequest>
{
    public UpdateRuleRequestValidator()
    {
        RuleFor(x => x.Name)
            .NotEmpty().WithMessage("Rule name is required")
            .MaximumLength(200).WithMessage("Rule name must not exceed 200 characters");

        RuleFor(x => x.Description)
            .MaximumLength(1000).WithMessage("Description must not exceed 1000 characters")
            .When(x => !string.IsNullOrEmpty(x.Description));

        RuleFor(x => x.Conditions)
            .NotEmpty().WithMessage("At least one condition is required")
            .Must(c => c.Count <= 50).WithMessage("Cannot exceed 50 conditions");

        RuleForEach(x => x.Conditions)
            .SetValidator(new RuleConditionDtoValidator());

        RuleFor(x => x.Actions)
            .NotEmpty().WithMessage("At least one action is required")
            .Must(a => a.Count <= 20).WithMessage("Cannot exceed 20 actions");

        RuleForEach(x => x.Actions)
            .SetValidator(new RuleActionDtoValidator());

        RuleFor(x => x.Priority)
            .GreaterThanOrEqualTo(0).WithMessage("Priority must be non-negative")
            .LessThanOrEqualTo(1000).WithMessage("Priority must not exceed 1000");

        RuleFor(x => x.Status)
            .IsInEnum().WithMessage("Invalid rule status")
            .When(x => x.Status.HasValue);

        RuleFor(x => x.EffectiveUntil)
            .Must((req, until) => !until.HasValue || !req.EffectiveFrom.HasValue || until > req.EffectiveFrom)
            .WithMessage("EffectiveUntil must be after EffectiveFrom")
            .When(x => x.EffectiveFrom.HasValue && x.EffectiveUntil.HasValue);
    }
}

public sealed class RuleConditionDtoValidator : AbstractValidator<RuleConditionDto>
{
    private const int DefaultMaxNestingDepth = 5;

    public RuleConditionDtoValidator() : this(DefaultMaxNestingDepth)
    {
    }

    private RuleConditionDtoValidator(int remainingDepth)
    {
        RuleFor(x => x.Field)
            .NotEmpty().WithMessage("Field name is required")
            .MaximumLength(100).WithMessage("Field name must not exceed 100 characters");

        RuleFor(x => x.Operator)
            .IsInEnum().WithMessage("Invalid comparison operator");

        RuleFor(x => x.Value)
            .NotNull().WithMessage("Condition value is required");

        RuleFor(x => x.LogicalOperator)
            .IsInEnum().WithMessage("Invalid logical operator")
            .When(x => x.LogicalOperator.HasValue);

        if (remainingDepth > 0)
        {
            RuleForEach(x => x.NestedConditions)
                .SetValidator(new RuleConditionDtoValidator(remainingDepth - 1))
                .When(x => x.NestedConditions is not null && x.NestedConditions.Count > 0);
        }
        else
        {
            RuleFor(x => x.NestedConditions)
                .Must(nested => nested is null || nested.Count == 0)
                .WithMessage($"Nested conditions exceed maximum depth of {DefaultMaxNestingDepth}");
        }
    }
}

public sealed class RuleActionDtoValidator : AbstractValidator<RuleActionDto>
{
    public RuleActionDtoValidator()
    {
        RuleFor(x => x.Type)
            .IsInEnum().WithMessage("Invalid action type");

        RuleFor(x => x.Parameters)
            .NotNull().WithMessage("Action parameters are required")
            .Must(p => p.Count <= 20).WithMessage("Cannot exceed 20 parameters");

        RuleFor(x => x.Order)
            .GreaterThanOrEqualTo(0).WithMessage("Order must be non-negative");

        RuleFor(x => x.Description)
            .MaximumLength(500).WithMessage("Description must not exceed 500 characters")
            .When(x => !string.IsNullOrEmpty(x.Description));
    }
}
