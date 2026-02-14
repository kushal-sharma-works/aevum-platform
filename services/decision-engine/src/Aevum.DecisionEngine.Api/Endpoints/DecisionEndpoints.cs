using FluentValidation;
using Aevum.DecisionEngine.Application.DTOs;
using Aevum.DecisionEngine.Application.Mapping;
using Aevum.DecisionEngine.Application.Services;
using Aevum.DecisionEngine.Domain.Interfaces;
using Microsoft.AspNetCore.Mvc;

namespace Aevum.DecisionEngine.Api.Endpoints;

public static class DecisionEndpoints
{
    public static RouteGroupBuilder MapDecisionEndpoints(this RouteGroupBuilder group)
    {
        group.MapPost("/evaluate", EvaluateDecisionAsync)
            .WithName("EvaluateDecision")
            .WithOpenApi()
            .Produces<DecisionResponse>(StatusCodes.Status200OK)
            .ProducesValidationProblem()
            .Produces(StatusCodes.Status404NotFound)
            .Produces(StatusCodes.Status422UnprocessableEntity);

        group.MapGet("/{id}", GetDecisionByIdAsync)
            .WithName("GetDecisionById")
            .WithOpenApi()
            .Produces<DecisionResponse>()
            .Produces(StatusCodes.Status404NotFound);

        group.MapGet("/request/{requestId}", GetDecisionByRequestIdAsync)
            .WithName("GetDecisionByRequestId")
            .WithOpenApi()
            .Produces<DecisionResponse>()
            .Produces(StatusCodes.Status404NotFound);

        group.MapGet("/rule/{ruleId}", GetDecisionsByRuleIdAsync)
            .WithName("GetDecisionsByRuleId")
            .WithOpenApi()
            .Produces<List<DecisionResponse>>();

        return group;
    }

    private static async Task<IResult> EvaluateDecisionAsync(
        EvaluateDecisionRequest request,
        [FromServices] IValidator<EvaluateDecisionRequest> validator,
        [FromServices] EvaluationService evaluationService,
        [FromServices] RuleManagementService ruleManagementService,
        [FromServices] TimeProvider timeProvider,
        CancellationToken cancellationToken)
    {
        var validationResult = await validator.ValidateAsync(request, cancellationToken);
        if (!validationResult.IsValid)
        {
            return Results.ValidationProblem(validationResult.ToDictionary());
        }

        var rule = await ruleManagementService.GetRuleAsync(
            request.RuleId,
            request.RuleVersion,
            cancellationToken);

        var context = request.ToEvaluationContext() with
        {
            Timestamp = timeProvider.GetUtcNow()
        };

        var decision = await evaluationService.EvaluateAsync(rule, context, cancellationToken);
        return Results.Ok(decision.ToResponse());
    }

    private static async Task<IResult> GetDecisionByIdAsync(
        string id,
        [FromServices] IDecisionRepository repository,
        CancellationToken cancellationToken)
    {
        var decision = await repository.GetByIdAsync(id, cancellationToken);
        
        if (decision is null)
        {
            return Results.NotFound(new { message = $"Decision with ID '{id}' not found" });
        }

        return Results.Ok(decision.ToResponse());
    }

    private static async Task<IResult> GetDecisionByRequestIdAsync(
        string requestId,
        [FromServices] IDecisionRepository repository,
        CancellationToken cancellationToken)
    {
        var decision = await repository.GetByRequestIdAsync(requestId, cancellationToken);
        
        if (decision is null)
        {
            return Results.NotFound(new { message = $"Decision with RequestId '{requestId}' not found" });
        }

        return Results.Ok(decision.ToResponse());
    }

    private static async Task<IResult> GetDecisionsByRuleIdAsync(
        string ruleId,
        int? version,
        [FromServices] IDecisionRepository repository,
        CancellationToken cancellationToken)
    {
        var decisions = await repository.GetByRuleIdAsync(ruleId, version, cancellationToken);
        var responses = decisions.Select(d => d.ToResponse()).ToList();
        return Results.Ok(responses);
    }
}
