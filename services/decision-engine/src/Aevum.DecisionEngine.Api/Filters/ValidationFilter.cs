using FluentValidation;
using Microsoft.AspNetCore.Http.HttpResults;

namespace Aevum.DecisionEngine.Api.Filters;

public static class ValidationFilter
{
    public static async Task<Results<Ok<TResponse>, ValidationProblem>> ValidateAsync<TRequest, TResponse>(
        TRequest request,
        IValidator<TRequest> validator,
        Func<TRequest, Task<TResponse>> handler)
    {
        var validationResult = await validator.ValidateAsync(request);
        
        if (!validationResult.IsValid)
        {
            var errors = validationResult.Errors
                .GroupBy(e => e.PropertyName)
                .ToDictionary(
                    g => g.Key,
                    g => g.Select(e => e.ErrorMessage).ToArray()
                );

            return TypedResults.ValidationProblem(errors);
        }

        var response = await handler(request);
        return TypedResults.Ok(response);
    }
}
