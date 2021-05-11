import { ei } from "../proto";
import { Artifact } from "../types";
import { aggregateEffect } from "./common";

export function researchPriceMultiplierFromArtifacts(artifacts: Artifact[]) {
  return aggregateEffect(
    artifacts,
    [ei.ArtifactSpec.Name.PUZZLE_CUBE],
    (aggregate, effect) =>
      effect.multiplier <= 1
        ? (1 + effect.delta * effect.multiplier) * aggregate
        : ((1 + effect.delta) / effect.multiplier) * aggregate,
    1
  );
}
