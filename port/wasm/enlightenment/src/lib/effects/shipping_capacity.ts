import { ei } from "../proto";
import { Artifact } from "../types";
import { multiplicativeEffect } from "./common";

export function shippingCapacityMultiplier(artifacts: Artifact[]) {
  return multiplicativeEffect(artifacts, [
    ei.ArtifactSpec.Name.INTERSTELLAR_COMPASS,
    ei.ArtifactSpec.Name.QUANTUM_STONE,
  ]);
}
