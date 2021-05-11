import { ei } from "../proto";
import { Artifact } from "../types";
import { multiplicativeEffect } from "./common";

export function awayEarningsMultiplier(artifacts: Artifact[]) {
  return multiplicativeEffect(artifacts, [
    ei.ArtifactSpec.Name.LUNAR_TOTEM,
    ei.ArtifactSpec.Name.LUNAR_STONE,
  ]);
}
