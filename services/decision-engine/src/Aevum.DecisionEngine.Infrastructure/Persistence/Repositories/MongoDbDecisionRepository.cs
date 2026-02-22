using MongoDB.Driver;
using Aevum.DecisionEngine.Domain.Interfaces;
using Aevum.DecisionEngine.Domain.Models;
using Aevum.DecisionEngine.Infrastructure.Persistence.Mapping;
using Aevum.DecisionEngine.Infrastructure.Persistence.Models;

namespace Aevum.DecisionEngine.Infrastructure.Persistence.Repositories;

public sealed class MongoDbDecisionRepository : IDecisionRepository
{
    private readonly IMongoCollection<DecisionDocument> _collection;

    public MongoDbDecisionRepository(MongoDbContext context)
    {
        _collection = context.GetCollection<DecisionDocument>("decisions");
        CreateIndexes();
    }

    private void CreateIndexes()
    {
        var requestIdIndexKeys = Builders<DecisionDocument>.IndexKeys.Ascending(d => d.RequestId);
        var requestIdIndexModel = new CreateIndexModel<DecisionDocument>(
            requestIdIndexKeys,
            new CreateIndexOptions { Unique = true });
        _collection.Indexes.CreateOne(requestIdIndexModel);

        var hashIndexKeys = Builders<DecisionDocument>.IndexKeys.Ascending(d => d.DeterministicHash);
        var hashIndexModel = new CreateIndexModel<DecisionDocument>(hashIndexKeys);
        _collection.Indexes.CreateOne(hashIndexModel);

        var ruleIndexKeys = Builders<DecisionDocument>.IndexKeys
            .Ascending(d => d.RuleId)
            .Ascending(d => d.RuleVersion);
        var ruleIndexModel = new CreateIndexModel<DecisionDocument>(ruleIndexKeys);
        _collection.Indexes.CreateOne(ruleIndexModel);
    }

    public async Task<Decision?> GetByIdAsync(string id, CancellationToken cancellationToken = default)
    {
        var filter = Builders<DecisionDocument>.Filter.Eq(d => d.Id, id);
        var doc = await _collection.Find(filter).FirstOrDefaultAsync(cancellationToken);
        return doc?.ToDomain();
    }

    public async Task<Decision?> GetByRequestIdAsync(string requestId, CancellationToken cancellationToken = default)
    {
        var filter = Builders<DecisionDocument>.Filter.Eq(d => d.RequestId, requestId);
        var doc = await _collection.Find(filter).FirstOrDefaultAsync(cancellationToken);
        return doc?.ToDomain();
    }

    public async Task<IReadOnlyList<Decision>> GetByRuleIdAsync(string ruleId, int? version = null, CancellationToken cancellationToken = default)
    {
        var filter = Builders<DecisionDocument>.Filter.Eq(d => d.RuleId, ruleId);
        
        if (version.HasValue)
        {
            filter &= Builders<DecisionDocument>.Filter.Eq(d => d.RuleVersion, version.Value);
        }

        var sort = Builders<DecisionDocument>.Sort.Descending(d => d.EvaluatedAt);
        var docs = await _collection.Find(filter).Sort(sort).Limit(100).ToListAsync(cancellationToken);
        return docs.Select(d => d.ToDomain()).ToList();
    }

    public async Task<Decision> CreateAsync(Decision decision, CancellationToken cancellationToken = default)
    {
        var doc = decision.ToDocument();
        await _collection.InsertOneAsync(doc, cancellationToken: cancellationToken);
        return doc.ToDomain();
    }
}
