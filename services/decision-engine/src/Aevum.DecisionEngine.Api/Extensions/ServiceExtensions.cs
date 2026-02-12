using FluentValidation;
using Aevum.DecisionEngine.Application.Evaluation;
using Aevum.DecisionEngine.Application.Services;
using Aevum.DecisionEngine.Application.Validators;
using Aevum.DecisionEngine.Domain.Interfaces;
using Aevum.DecisionEngine.Infrastructure.Http;
using Aevum.DecisionEngine.Infrastructure.Persistence;
using Aevum.DecisionEngine.Infrastructure.Persistence.Repositories;
using Polly;
using Polly.Extensions.Http;

namespace Aevum.DecisionEngine.Api.Extensions;

public static class ServiceExtensions
{
    public static IServiceCollection AddApplicationServices(this IServiceCollection services, IConfiguration configuration)
    {
        // TimeProvider for deterministic evaluation
        services.AddSingleton(TimeProvider.System);

        // Validators
        services.AddValidatorsFromAssemblyContaining<CreateRuleRequestValidator>();

        // Application services
        services.AddSingleton<IDeterministicEvaluator, DeterministicEvaluator>();
        services.AddScoped<EvaluationService>();
        services.AddScoped<RuleManagementService>();

        // Infrastructure - MongoDB
        var mongoConnectionString = configuration["MongoDB:ConnectionString"] 
            ?? throw new InvalidOperationException("MongoDB:ConnectionString is required");
        var mongoDatabaseName = configuration["MongoDB:DatabaseName"] 
            ?? throw new InvalidOperationException("MongoDB:DatabaseName is required");

        services.AddSingleton(sp => new MongoDbContext(mongoConnectionString, mongoDatabaseName));
        services.AddSingleton<IRuleRepository, MongoDbRuleRepository>();
        services.AddSingleton<IDecisionRepository, MongoDbDecisionRepository>();

        // Infrastructure - Event Timeline HTTP Client
        var eventTimelineBaseUrl = configuration["EventTimeline:BaseUrl"];
        if (!string.IsNullOrEmpty(eventTimelineBaseUrl))
        {
            services.AddHttpClient<IEventTimelineClient, EventTimelineClient>(client =>
            {
                client.BaseAddress = new Uri(eventTimelineBaseUrl);
                client.Timeout = TimeSpan.FromSeconds(10);
            })
            .AddPolicyHandler(GetRetryPolicy())
            .AddPolicyHandler(GetCircuitBreakerPolicy());
        }
        else
        {
            // Fallback no-op client if Event Timeline is not configured
            services.AddSingleton<IEventTimelineClient>(new NoOpEventTimelineClient());
        }

        return services;
    }

    private static IAsyncPolicy<HttpResponseMessage> GetRetryPolicy()
    {
        return HttpPolicyExtensions
            .HandleTransientHttpError()
            .WaitAndRetryAsync(3, retryAttempt => TimeSpan.FromMilliseconds(100 * Math.Pow(2, retryAttempt)));
    }

    private static IAsyncPolicy<HttpResponseMessage> GetCircuitBreakerPolicy()
    {
        return HttpPolicyExtensions
            .HandleTransientHttpError()
            .CircuitBreakerAsync(5, TimeSpan.FromSeconds(30));
    }

    private sealed class NoOpEventTimelineClient : IEventTimelineClient
    {
        public Task<bool> IngestEventAsync(string streamId, string eventType, object data, CancellationToken cancellationToken = default)
        {
            return Task.FromResult(true);
        }
    }
}
