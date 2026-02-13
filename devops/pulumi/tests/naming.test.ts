import { describe, expect, it, vi } from "vitest"

describe("resourceName", () => {
  it("returns project-env-resource format", async () => {
    vi.resetModules()
    vi.doMock("@pulumi/pulumi", () => ({ getStack: () => "dev" }))
    const { resourceName } = await import("../utils/naming")
    expect(resourceName("events")).toBe("aevum-dev-events")
  })
})
