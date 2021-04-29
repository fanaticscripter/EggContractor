import { ei } from "../proto";
import { Artifact } from "../types";
import { multiplicativeEffect } from "./common";

export function habSpaceMultiplier(artifacts: Artifact[]) {
  return multiplicativeEffect(artifacts, [ei.ArtifactSpec.Name.ORNATE_GUSSET]);
}
