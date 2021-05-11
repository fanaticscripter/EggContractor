import { ei } from "../proto";
import { Artifact } from "../types";
import { multiplicativeEffect } from "./common";

export function droneRewardsMultiplier(artifacts: Artifact[]) {
  return multiplicativeEffect(artifacts, [ei.ArtifactSpec.Name.AURELIAN_BROOCH]);
}
