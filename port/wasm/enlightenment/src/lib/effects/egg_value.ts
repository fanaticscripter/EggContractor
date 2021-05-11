import { ei } from "../proto";
import { Artifact } from "../types";
import { multiplicativeEffect } from "./common";

export function eggValueMultiplier(artifacts: Artifact[]) {
  return multiplicativeEffect(artifacts, [
    ei.ArtifactSpec.Name.LIGHT_OF_EGGENDIL,
    ei.ArtifactSpec.Name.DEMETERS_NECKLACE,
    ei.ArtifactSpec.Name.TUNGSTEN_ANKH,
    ei.ArtifactSpec.Name.SHELL_STONE,
  ]);
}
