using Aevum.DecisionEngine.Api.Endpoints;
using Aevum.DecisionEngine.Api.Extensions;
using Aevum.DecisionEngine.Api.Middleware;
using Prometheus;
using Serilog;

var builder = WebApplication.CreateBuilder(args);

// Add Serilog
builder.AddSerilog();

// Add services
builder.Services.AddApplicationServices(builder.Configuration);
builder.Services.AddObservability(builder.Configuration);

// Add minimal API support
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddOpenApi();

// Add CORS
builder.Services.AddCors(options =>
{
    options.AddDefaultPolicy(policy =>
    {
        policy.AllowAnyOrigin()
              .AllowAnyMethod()
              .AllowAnyHeader();
    });
});

// Health checks
builder.Services.AddHealthChecks();

var app = builder.Build();

// Configure middleware pipeline
app.UseSerilogRequestLogging();
app.UseGlobalExceptionHandler();

if (app.Environment.IsDevelopment())
{
    app.MapOpenApi();
}

app.UseCors();

// Prometheus metrics endpoint
app.UseHttpMetrics();
app.MapMetrics();

// Health check endpoints
app.MapHealthChecks("/health");
app.MapHealthChecks("/health/ready");
app.MapHealthChecks("/health/live");

// API endpoints
var apiV1 = app.MapGroup("/api/v1");

apiV1.MapGroup("/rules")
    .WithTags("Rules")
    .WithOpenApi()
    .MapRuleEndpoints();

apiV1.MapGroup("/decisions")
    .WithTags("Decisions")
    .WithOpenApi()
    .MapDecisionEndpoints();

// Root endpoint
app.MapGet("/", () => new
{
    service = "Decision Engine",
    version = "1.0.0",
    status = "healthy"
})
.ExcludeFromDescription();

Log.Information("Starting Decision Engine service on {Port}", builder.Configuration["ASPNETCORE_URLS"] ?? "http://localhost:5000");

app.Run();

// Make Program accessible for integration tests
public partial class Program { }
