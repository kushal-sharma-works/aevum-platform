namespace Aevum.DecisionEngine.Domain.Interfaces;

public interface IEventTimelineClient
{
    Task<bool> IngestEventAsync(string streamId, string eventType, object data, CancellationToken cancellationToken = default);
}
