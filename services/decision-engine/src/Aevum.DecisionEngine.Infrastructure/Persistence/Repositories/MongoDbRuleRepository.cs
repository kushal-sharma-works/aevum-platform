using MongoDB.Driver;
using Aevum.DecisionEngine.Domain.Enums;
using Aevum.DecisionEngine.Domain.Interfaces;
using Aevum.DecisionEngine.Domain.Models;
using Aevum.DecisionEngine.Infrastructure.Persistence.Mapping;
using Aevum.DecisionEngine.Infrastructure.Persistence.Models;

namespace Aevum.DecisionEngine.Infrastructure.Persistence.Repositories;

public sealed class MongoDbRuleRepository : IRuleRepository
{
    private readonly IMongoCollection<RuleDocument> _collection;

    public MongoDbRuleRepository(MongoDbContext context)
    {
        _collection = context.GetCollection<RuleDocument>("rules");
        CreateIndexes();
    }

    private void CreateIndexes()
    {
        var indexKeys = Builders<RuleDocument>.IndexKeys
            .Ascending(r => r.Status)
            .Ascending(r => r.Priority);
        var indexModel = new CreateIndexModel<RuleDocument>(indexKeys);
        _collection.Indexes.CreateOne(indexModel);

        var versionIndexKeys = Builders<RuleDocument>.IndexKeys
            .Ascending(r => r.Id)
            .Descending(r => r.Version);
        var versionIndexModel = new CreateIndexModel<RuleDocument>(versionIndexKeys);
        _collection.Indexes.CreateOne(versionIndexModel);
    }

    public async Task<Rule?> GetByIdAsync(string id, int? version = null, CancellationToken cancellationToken = default)
    {
        var filter = Builders<RuleDocument>.Filter.Eq(r => r.Id, id);
        
        if (version.HasValue)
        {
            filter &= Builders<RuleDocument>.Filter.Eq(r => r.Version, version.Value);
        }

        var sort = Builders<RuleDocument>.Sort.Descending(r => r.Version);
        var doc = await _collection.Find(filter).Sort(sort).FirstOrDefaultAsync(cancellationToken);
        return doc?.ToDomain();
    }

    public async Task<IReadOnlyList<Rule>> GetActiveRulesAsync(CancellationToken cancellationToken = default)
    {
        var filter = Builders<RuleDocument>.Filter.Eq(r => r.Status, RuleStatus.Active);
        var sort = Builders<RuleDocument>.Sort.Descending(r => r.Priority);
        var docs = await _collection.Find(filter).Sort(sort).ToListAsync(cancellationToken);
        return docs.Select(d => d.ToDomain()).ToList();
    }

    public async Task<IReadOnlyList<Rule>> GetByStatusAsync(RuleStatus status, CancellationToken cancellationToken = default)
    {
        var filter = Builders<RuleDocument>.Filter.Eq(r => r.Status, status);
        var sort = Builders<RuleDocument>.Sort.Descending(r => r.UpdatedAt);
        var docs = await _collection.Find(filter).Sort(sort).ToListAsync(cancellationToken);
        return docs.Select(d => d.ToDomain()).ToList();
    }

    public async Task<Rule> CreateAsync(Rule rule, CancellationToken cancellationToken = default)
    {
        var doc = rule.ToDocument();
        await _collection.InsertOneAsync(doc, cancellationToken: cancellationToken);
        return doc.ToDomain();
    }

    public async Task<Rule> UpdateAsync(Rule rule, CancellationToken cancellationToken = default)
    {
        var doc = rule.ToDocument();
        await _collection.InsertOneAsync(doc, cancellationToken: cancellationToken);
        return doc.ToDomain();
    }

    public async Task DeleteAsync(string id, CancellationToken cancellationToken = default)
    {
        var filter = Builders<RuleDocument>.Filter.Eq(r => r.Id, id);
        await _collection.DeleteManyAsync(filter, cancellationToken);
    }

    public async Task<int> GetLatestVersionAsync(string id, CancellationToken cancellationToken = default)
    {
        var filter = Builders<RuleDocument>.Filter.Eq(r => r.Id, id);
        var sort = Builders<RuleDocument>.Sort.Descending(r => r.Version);
        var doc = await _collection.Find(filter).Sort(sort).FirstOrDefaultAsync(cancellationToken);
        return doc?.Version ?? 0;
    }
}
