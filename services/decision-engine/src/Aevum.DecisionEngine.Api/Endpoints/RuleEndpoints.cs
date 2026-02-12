using FluentValidation;
using Aevum.DecisionEngine.Application.DTOs;
using Aevum.DecisionEngine.Application.Mapping;
using Aevum.DecisionEngine.Application.Services;
using Aevum.DecisionEngine.Domain.Enums;

namespace Aevum.DecisionEngine.Api.Endpoints;

public static class RuleEndpoints
{
    public static RouteGroupBuilder MapRuleEndpoints(this RouteGroupBuilder group)
    {
        group.MapPost("/", CreateRuleAsync)
            .WithName("CreateRule")
            .WithOpenApi()
            .Produces<RuleResponse>(StatusCodes.Status201Created)
            .ProducesValidationProblem();

        group.MapGet("/{id}", GetRuleByIdAsync)
            .WithName("GetRuleById")
            .WithOpenApi()
            .Produces<RuleResponse>()
            .Produces(StatusCodes.Status404NotFound);

        group.MapPut("/{id}", UpdateRuleAsync)
            .WithName("UpdateRule")
            .WithOpenApi()
            .Produces<RuleResponse>()
            .ProducesValidationProblem()
            .Produces(StatusCodes.Status404NotFound);

        group.MapDelete("/{id}", DeleteRuleAsync)
            .WithName("DeleteRule")
            .WithOpenApi()
            .Produces(StatusCodes.Status204NoContent)
            .Produces(StatusCodes.Status404NotFound);

        group.MapGet("/", GetRulesAsync)
            .WithName("GetRules")
            .WithOpenApi()
            .Produces<List<RuleResponse>>();

        group.MapPost("/{id}/activate", ActivateRuleAsync)
            .WithName("ActivateRule")
            .WithOpenApi()
            .Produces<RuleResponse>()
            .Produces(StatusCodes.Status404NotFound);

        group.MapPost("/{id}/deactivate", DeactivateRuleAsync)
            .WithName("DeactivateRule")
            .WithOpenApi()
            .Produces<RuleResponse>()
            .Produces(StatusCodes.Status404NotFound);

        return group;
    }

    private static async Task<IResult> CreateRuleAsync(
        CreateRuleRequest request,
        IValidator<CreateRuleRequest> validator,
        RuleManagementService service,
        CancellationToken cancellationToken)
    {
        var validationResult = await validator.ValidateAsync(request, cancellationToken);
        if (!validationResult.IsValid)
        {
            return Results.ValidationProblem(validationResult.ToDictionary());
        }

        var rule = request.ToDomain();
        var created = await service.CreateRuleAsync(rule, cancellationToken);
        var response = created.ToResponse();

        return Results.Created($"/api/v1/rules/{response.Id}", response);
    }

    private static async Task<IResult> GetRuleByIdAsync(
        string id,
        int? version,
        RuleManagementService service,
        CancellationToken cancellationToken)
    {
        var rule = await service.GetRuleAsync(id, version, cancellationToken);
        return Results.Ok(rule.ToResponse());
    }

    private static async Task<IResult> UpdateRuleAsync(
        string id,
        UpdateRuleRequest request,
        IValidator<UpdateRuleRequest> validator,
        RuleManagementService service,
        CancellationToken cancellationToken)
    {
        var validationResult = await validator.ValidateAsync(request, cancellationToken);
        if (!validationResult.IsValid)
        {
            return Results.ValidationProblem(validationResult.ToDictionary());
        }

        var updates = request.ToDomain(id);
        var updated = await service.UpdateRuleAsync(id, updates, cancellationToken);
        return Results.Ok(updated.ToResponse());
    }

    private static async Task<IResult> DeleteRuleAsync(
        string id,
        RuleManagementService service,
        CancellationToken cancellationToken)
    {
        await service.DeleteRuleAsync(id, cancellationToken);
        return Results.NoContent();
    }

    private static async Task<IResult> GetRulesAsync(
        RuleStatus? status,
        RuleManagementService service,
        CancellationToken cancellationToken)
    {
        var rules = status.HasValue
            ? await service.GetRulesByStatusAsync(status.Value, cancellationToken)
            : await service.GetActiveRulesAsync(cancellationToken);

        var responses = rules.Select(r => r.ToResponse()).ToList();
        return Results.Ok(responses);
    }

    private static async Task<IResult> ActivateRuleAsync(
        string id,
        RuleManagementService service,
        CancellationToken cancellationToken)
    {
        var activated = await service.ActivateRuleAsync(id, cancellationToken);
        return Results.Ok(activated.ToResponse());
    }

    private static async Task<IResult> DeactivateRuleAsync(
        string id,
        RuleManagementService service,
        CancellationToken cancellationToken)
    {
        var deactivated = await service.DeactivateRuleAsync(id, cancellationToken);
        return Results.Ok(deactivated.ToResponse());
    }
}
