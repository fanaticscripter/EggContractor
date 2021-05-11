import { ei } from "../proto";
import { Artifact } from "../types";
import { additiveEffect } from "./common";

export function maxRCBBonusFromArtifacts(artifacts: Artifact[]) {
  return additiveEffect(artifacts, [
    ei.ArtifactSpec.Name.VIAL_MARTIAN_DUST,
    ei.ArtifactSpec.Name.TERRA_STONE,
  ]);
}
