import { ei } from "../proto";
import { Artifact } from "../types";
import { multiplicativeEffect } from "./common";

export function eggLayingRateMultiplier(artifacts: Artifact[]) {
  return multiplicativeEffect(artifacts, [
    ei.ArtifactSpec.Name.QUANTUM_METRONOME,
    ei.ArtifactSpec.Name.TACHYON_STONE,
  ]);
}
