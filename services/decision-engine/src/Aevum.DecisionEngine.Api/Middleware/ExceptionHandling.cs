using System.Net;
using System.Text.Json;
using Microsoft.AspNetCore.Diagnostics;
using Microsoft.AspNetCore.Mvc;
using Aevum.DecisionEngine.Domain.Exceptions;
using FluentValidation;

namespace Aevum.DecisionEngine.Api.Middleware;

public static class ExceptionHandlingExtensions
{
    public static IApplicationBuilder UseGlobalExceptionHandler(this IApplicationBuilder app)
    {
        app.UseExceptionHandler(errorApp =>
        {
            errorApp.Run(async context =>
            {
                var exceptionHandlerFeature = context.Features.Get<IExceptionHandlerFeature>();
                var exception = exceptionHandlerFeature?.Error;

                if (exception is null)
                    return;

                var (statusCode, problemDetails) = exception switch
                {
                    ValidationException validationEx => (
                        HttpStatusCode.BadRequest,
                        CreateValidationProblemDetails(context, validationEx)
                    ),
                    RuleNotFoundException notFoundEx => (
                        HttpStatusCode.NotFound,
                        CreateProblemDetails(context, HttpStatusCode.NotFound, "Rule not found", notFoundEx.Message)
                    ),
                    EvaluationException evalEx => (
                        HttpStatusCode.UnprocessableEntity,
                        CreateProblemDetails(context, HttpStatusCode.UnprocessableEntity, "Evaluation failed", evalEx.Message)
                    ),
                    InvalidRuleException invalidEx => (
                        HttpStatusCode.BadRequest,
                        CreateProblemDetails(context, HttpStatusCode.BadRequest, "Invalid rule", invalidEx.Message)
                    ),
                    DomainException domainEx => (
                        HttpStatusCode.BadRequest,
                        CreateProblemDetails(context, HttpStatusCode.BadRequest, "Domain error", domainEx.Message)
                    ),
                    _ => (
                        HttpStatusCode.InternalServerError,
                        CreateProblemDetails(context, HttpStatusCode.InternalServerError, "Internal server error", "An unexpected error occurred")
                    )
                };

                context.Response.StatusCode = (int)statusCode;
                context.Response.ContentType = "application/problem+json";

                await context.Response.WriteAsync(JsonSerializer.Serialize(problemDetails, new JsonSerializerOptions
                {
                    PropertyNamingPolicy = JsonNamingPolicy.CamelCase
                }));
            });
        });

        return app;
    }

    private static ProblemDetails CreateProblemDetails(
        HttpContext context,
        HttpStatusCode statusCode,
        string title,
        string detail)
    {
        return new ProblemDetails
        {
            Status = (int)statusCode,
            Type = $"https://httpstatuses.com/{(int)statusCode}",
            Title = title,
            Detail = detail,
            Instance = context.Request.Path
        };
    }

    private static ValidationProblemDetails CreateValidationProblemDetails(
        HttpContext context,
        ValidationException validationException)
    {
        var errors = validationException.Errors
            .GroupBy(e => e.PropertyName)
            .ToDictionary(
                g => g.Key,
                g => g.Select(e => e.ErrorMessage).ToArray()
            );

        return new ValidationProblemDetails(errors)
        {
            Status = (int)HttpStatusCode.BadRequest,
            Type = "https://httpstatuses.com/400",
            Title = "Validation failed",
            Detail = "One or more validation errors occurred",
            Instance = context.Request.Path
        };
    }
}
