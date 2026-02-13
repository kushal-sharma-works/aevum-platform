import * as pulumi from "@pulumi/pulumi"

export function resourceName(resource: string): string {
  return `aevum-${pulumi.getStack()}-${resource}`
}
