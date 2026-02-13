import * as pulumi from "@pulumi/pulumi"

export function defaultTags(): Record<string, string> {
  return {
    Project: "aevum",
    Environment: pulumi.getStack(),
    ManagedBy: "pulumi",
    Team: "platform"
  }
}
