# Build stage
FROM mcr.microsoft.com/dotnet/sdk:9.0 AS build
WORKDIR /src

# Copy solution and project files
COPY DecisionEngine.sln ./
COPY src/Aevum.DecisionEngine.Api/Aevum.DecisionEngine.Api.csproj src/Aevum.DecisionEngine.Api/
COPY src/Aevum.DecisionEngine.Domain/Aevum.DecisionEngine.Domain.csproj src/Aevum.DecisionEngine.Domain/
COPY src/Aevum.DecisionEngine.Application/Aevum.DecisionEngine.Application.csproj src/Aevum.DecisionEngine.Application/
COPY src/Aevum.DecisionEngine.Infrastructure/Aevum.DecisionEngine.Infrastructure.csproj src/Aevum.DecisionEngine.Infrastructure/

# Restore dependencies
RUN dotnet restore

# Copy source code
COPY src/ src/

# Build
RUN dotnet build -c Release --no-restore
RUN dotnet publish src/Aevum.DecisionEngine.Api/Aevum.DecisionEngine.Api.csproj -c Release -o /app/publish --no-restore

# Runtime stage
FROM mcr.microsoft.com/dotnet/aspnet:9.0-alpine AS runtime
WORKDIR /app

# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -u 1000 -G appuser -s /bin/sh -D appuser

# Copy published app
COPY --from=build /app/publish .

# Set ownership
RUN chown -R appuser:appuser /app

USER appuser

EXPOSE 8080
EXPOSE 9090

ENV ASPNETCORE_URLS=http://+:8080 \
    ASPNETCORE_ENVIRONMENT=Production \
    DOTNET_RUNNING_IN_CONTAINER=true

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

ENTRYPOINT ["dotnet", "Aevum.DecisionEngine.Api.dll"]
