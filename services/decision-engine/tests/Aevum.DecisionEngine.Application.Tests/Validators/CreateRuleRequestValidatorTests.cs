using Aevum.DecisionEngine.Application.DTOs;
using Aevum.DecisionEngine.Application.Validators;
using Aevum.DecisionEngine.Domain.Enums;
using Xunit;

namespace Aevum.DecisionEngine.Application.Tests.Validators;

public class CreateRuleRequestValidatorTests
{
    private readonly CreateRuleRequestValidator _validator;

    public CreateRuleRequestValidatorTests()
    {
        _validator = new CreateRuleRequestValidator();
    }

    [Fact]
    public async Task Validate_ValidRequest_Succeeds()
    {
        // Arrange
        var request = new CreateRuleRequest
        {
            Name = "Test Rule",
            Description = "A test rule",
            Priority = 1,
            Conditions = new[]
            {
                new RuleConditionDto
                {
                    Field = "status",
                    Operator = ComparisonOperator.Equals,
                    Value = "active"
                }
            },
            Actions = new[]
            {
                new RuleActionDto
                {
                    Type = ActionType.LogEvent,
                    Parameters = new Dictionary<string, object> { { "message", "Rule triggered" } }
                }
            }
        };

        // Act
        var result = await _validator.ValidateAsync(request);

        // Assert
        Assert.True(result.IsValid);
    }

    [Fact]
    public async Task Validate_MissingName_Fails()
    {
        // Arrange
        var request = new CreateRuleRequest
        {
            Name = "",
            Description = "A test rule",
            Priority = 1,
            Conditions = new[] { new RuleConditionDto { Field = "status", Operator = ComparisonOperator.Equals, Value = "active" } },
            Actions = new[] { new RuleActionDto { Type = ActionType.LogEvent, Parameters = new Dictionary<string, object>() } }
        };

        // Act
        var result = await _validator.ValidateAsync(request);

        // Assert
        Assert.False(result.IsValid);
    }
}
