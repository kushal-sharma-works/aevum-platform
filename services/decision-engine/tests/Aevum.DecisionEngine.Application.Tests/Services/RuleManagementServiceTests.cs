using FluentAssertions;
using NSubstitute;
using Aevum.DecisionEngine.Application.Services;
using Aevum.DecisionEngine.Domain.Enums;
using Aevum.DecisionEngine.Domain.Exceptions;
using Aevum.DecisionEngine.Domain.Interfaces;
using Aevum.DecisionEngine.Domain.Models;

namespace Aevum.DecisionEngine.Application.Tests.Services;

public sealed class RuleManagementServiceTests
{
    private static Rule CreateRule(string id, RuleStatus status = RuleStatus.Draft, int version = 1) => new()
    {
        Id = id,
        Name = "Rule",
        Description = "desc",
        Conditions = [new RuleCondition { Field = "amount", Operator = ComparisonOperator.GreaterThan, Value = 10 }],
        Actions = [new RuleAction { Type = ActionType.StoreDecision, Parameters = new Dictionary<string, object> { ["k"] = "v" }, Order = 1 }],
        Status = status,
        Version = version,
        Priority = 10,
        CreatedAt = DateTimeOffset.UtcNow.AddMinutes(-10),
        UpdatedAt = DateTimeOffset.UtcNow.AddMinutes(-5)
    };

    [Fact]
    public async Task CreateRuleAsync_ShouldSetDefaults()
    {
        var repository = Substitute.For<IRuleRepository>();
        Rule captured = null!;
        repository.CreateAsync(Arg.Any<Rule>(), Arg.Any<CancellationToken>())
            .Returns(ci =>
            {
                captured = ci.Arg<Rule>();
                return captured;
            });

        var service = new RuleManagementService(repository, TimeProvider.System);
        var input = CreateRule(string.Empty, status: RuleStatus.Active, version: 99);

        var created = await service.CreateRuleAsync(input);

        created.Status.Should().Be(RuleStatus.Draft);
        created.Version.Should().Be(1);
        created.Id.Should().NotBeNullOrWhiteSpace();
        captured.Should().NotBeNull();
    }

    [Fact]
    public async Task ActivateRuleAsync_ShouldUpdateStatus()
    {
        var repository = Substitute.For<IRuleRepository>();
        var existing = CreateRule("rule-1", RuleStatus.Draft, 1);

        repository.GetByIdAsync("rule-1", null, Arg.Any<CancellationToken>()).Returns(existing);
        repository.UpdateAsync(Arg.Any<Rule>(), Arg.Any<CancellationToken>())
            .Returns(ci => ci.Arg<Rule>());

        var service = new RuleManagementService(repository, TimeProvider.System);

        var activated = await service.ActivateRuleAsync("rule-1");

        activated.Status.Should().Be(RuleStatus.Active);
        await repository.Received(1).UpdateAsync(Arg.Is<Rule>(r => r.Status == RuleStatus.Active), Arg.Any<CancellationToken>());
    }

    [Fact]
    public async Task UpdateRuleAsync_ShouldCreateNewVersion()
    {
        var repository = Substitute.For<IRuleRepository>();
        var existing = CreateRule("rule-1", RuleStatus.Active, 2);
        var updates = CreateRule("rule-1", RuleStatus.Active, 0) with { Name = "Updated Name" };

        repository.GetByIdAsync("rule-1", null, Arg.Any<CancellationToken>()).Returns(existing);
        repository.GetLatestVersionAsync("rule-1", Arg.Any<CancellationToken>()).Returns(2);
        repository.CreateAsync(Arg.Any<Rule>(), Arg.Any<CancellationToken>())
            .Returns(ci => ci.Arg<Rule>());

        var service = new RuleManagementService(repository, TimeProvider.System);

        var updated = await service.UpdateRuleAsync("rule-1", updates);

        updated.Version.Should().Be(3);
        updated.CreatedAt.Should().Be(existing.CreatedAt);
        updated.CreatedBy.Should().Be(existing.CreatedBy);
        updated.Name.Should().Be("Updated Name");
        await repository.Received(1).CreateAsync(Arg.Is<Rule>(r => r.Version == 3), Arg.Any<CancellationToken>());
        await repository.DidNotReceive().UpdateAsync(Arg.Any<Rule>(), Arg.Any<CancellationToken>());
    }

    [Fact]
    public async Task DeactivateRuleAsync_ShouldUpdateStatus()
    {
        var repository = Substitute.For<IRuleRepository>();
        var existing = CreateRule("rule-1", RuleStatus.Active, 1);

        repository.GetByIdAsync("rule-1", null, Arg.Any<CancellationToken>()).Returns(existing);
        repository.UpdateAsync(Arg.Any<Rule>(), Arg.Any<CancellationToken>())
            .Returns(ci => ci.Arg<Rule>());

        var service = new RuleManagementService(repository, TimeProvider.System);

        var deactivated = await service.DeactivateRuleAsync("rule-1");

        deactivated.Status.Should().Be(RuleStatus.Inactive);
        await repository.Received(1).UpdateAsync(Arg.Is<Rule>(r => r.Status == RuleStatus.Inactive), Arg.Any<CancellationToken>());
    }

    [Fact]
    public async Task DeleteRuleAsync_ShouldDeleteWhenFound()
    {
        var repository = Substitute.For<IRuleRepository>();
        var existing = CreateRule("rule-1");

        repository.GetByIdAsync("rule-1", null, Arg.Any<CancellationToken>()).Returns(existing);

        var service = new RuleManagementService(repository, TimeProvider.System);

        await service.DeleteRuleAsync("rule-1");

        await repository.Received(1).DeleteAsync("rule-1", Arg.Any<CancellationToken>());
    }

    [Fact]
    public async Task GetRuleAsync_ShouldThrowWhenMissing()
    {
        var repository = Substitute.For<IRuleRepository>();
        repository.GetByIdAsync("missing", null, Arg.Any<CancellationToken>()).Returns((Rule?)null);

        var service = new RuleManagementService(repository, TimeProvider.System);

        var act = async () => await service.GetRuleAsync("missing");

        await act.Should().ThrowAsync<RuleNotFoundException>();
    }
}
