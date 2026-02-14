using System.Net;
using System.Net.Http.Json;
using FluentAssertions;
using Microsoft.AspNetCore.Mvc.Testing;
using Microsoft.Extensions.DependencyInjection;
using Testcontainers.MongoDb;
using Aevum.DecisionEngine.Application.DTOs;
using Aevum.DecisionEngine.Domain.Enums;
using Aevum.DecisionEngine.Infrastructure.Persistence;

namespace Aevum.DecisionEngine.Integration.Tests;

public sealed class DecisionEngineIntegrationTests : IAsyncLifetime
{
    private WebApplicationFactory<Program> _factory = null!;
    private HttpClient _client = null!;
    private MongoDbContainer _mongoContainer = null!;

    public async Task InitializeAsync()
    {
        _mongoContainer = new MongoDbBuilder()
            .WithImage("mongo:8.0")
            .Build();

        await _mongoContainer.StartAsync();

        _factory = new WebApplicationFactory<Program>()
            .WithWebHostBuilder(builder =>
            {
                builder.ConfigureServices(services =>
                {
                    // Replace MongoDB context with test container
                    var descriptor = services.SingleOrDefault(d => d.ServiceType == typeof(MongoDbContext));
                    if (descriptor != null)
                    {
                        services.Remove(descriptor);
                    }

                    services.AddSingleton(_ => new MongoDbContext(
                        _mongoContainer.GetConnectionString(),
                        "decision_engine_test"));
                });
            });

        _client = _factory.CreateClient();
    }

    public async Task DisposeAsync()
    {
        _client?.Dispose();
        await _factory.DisposeAsync();
        await _mongoContainer.DisposeAsync();
    }

    [Fact]
    public async Task HealthCheck_ShouldReturnHealthy()
    {
        // Act
        var response = await _client.GetAsync("/health");

        // Assert
        response.StatusCode.Should().Be(HttpStatusCode.OK);
    }

    [Fact]
    public async Task CreateRule_ShouldSucceed()
    {
        // Arrange
        var request = new CreateRuleRequest
        {
            Name = "Test Rule",
            Description = "Integration test rule",
            Conditions =
            [
                new RuleConditionDto
                {
                    Field = "amount",
                    Operator = ComparisonOperator.GreaterThan,
                    Value = 100
                }
            ],
            Actions =
            [
                new RuleActionDto
                {
                    Type = ActionType.StoreDecision,
                    Parameters = new Dictionary<string, object> { ["key"] = "value" },
                    Order = 1
                }
            ],
            Priority = 10
        };

        // Act
        var response = await _client.PostAsJsonAsync("/api/v1/rules", request);

        // Assert
        response.StatusCode.Should().Be(HttpStatusCode.Created);
        var rule = await response.Content.ReadFromJsonAsync<RuleResponse>();
        rule.Should().NotBeNull();
        rule!.Name.Should().Be("Test Rule");
        rule.Status.Should().Be(RuleStatus.Draft);
    }

    [Fact]
    public async Task EvaluateDecision_ShouldSucceed()
    {
        // Arrange - Create a rule first
        var createRuleRequest = new CreateRuleRequest
        {
            Name = "Evaluation Test Rule",
            Conditions =
            [
                new RuleConditionDto
                {
                    Field = "score",
                    Operator = ComparisonOperator.GreaterThanOrEqual,
                    Value = 80
                }
            ],
            Actions =
            [
                new RuleActionDto
                {
                    Type = ActionType.StoreDecision,
                    Parameters = new Dictionary<string, object> { ["result"] = "pass" },
                    Order = 1
                }
            ],
            Priority = 10
        };

        var createResponse = await _client.PostAsJsonAsync("/api/v1/rules", createRuleRequest);
        var rule = await createResponse.Content.ReadFromJsonAsync<RuleResponse>();

        // Activate the rule
        await _client.PostAsync($"/api/v1/rules/{rule!.Id}/activate", null);

        var evaluateRequest = new EvaluateDecisionRequest
        {
            RuleId = rule.Id,
            Context = new Dictionary<string, object> { ["score"] = 85 },
            RequestId = Guid.NewGuid().ToString()
        };

        // Act
        var response = await _client.PostAsJsonAsync("/api/v1/decisions/evaluate", evaluateRequest);

        // Assert
        response.StatusCode.Should().Be(HttpStatusCode.OK);
        var decision = await response.Content.ReadFromJsonAsync<DecisionResponse>();
        decision.Should().NotBeNull();
        decision!.Status.Should().Be(DecisionStatus.Approved);
        decision.MatchedConditions.Should().NotBeEmpty();
    }
}
