import { ei } from "../proto";
import { Artifact } from "../types";
import { multiplicativeEffect } from "./common";

export function farmValueMultiplier(artifacts: Artifact[]) {
  return multiplicativeEffect(artifacts, [ei.ArtifactSpec.Name.MERCURYS_LENS]);
}
