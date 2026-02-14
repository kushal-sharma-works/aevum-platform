using System.Net.Http.Json;
using System.Text.Json;
using Aevum.DecisionEngine.Domain.Interfaces;

namespace Aevum.DecisionEngine.Infrastructure.Http;

public sealed class EventTimelineClient(HttpClient httpClient) : IEventTimelineClient
{
    private readonly HttpClient _httpClient = httpClient;

    public async Task<bool> IngestEventAsync(
        string streamId,
        string eventType,
        object data,
        CancellationToken cancellationToken = default)
    {
        try
        {
            var request = new IngestEventRequest
            {
                StreamId = streamId,
                EventType = eventType,
                Data = JsonSerializer.SerializeToElement(data)
            };

            var response = await _httpClient.PostAsJsonAsync("/api/v1/events", request, cancellationToken);
            return response.IsSuccessStatusCode;
        }
        catch
        {
            return false;
        }
    }

    private sealed record IngestEventRequest
    {
        public required string StreamId { get; init; }
        public required string EventType { get; init; }
        public required JsonElement Data { get; init; }
    }
}
