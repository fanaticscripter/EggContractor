import { ei } from "../proto";
import { Artifact } from "../types";

type Effect = {
  delta: number;
  multiplier: number;
};

// Only effects on the enlightenment farm are considered here.
function gatherRelevantEffects(artifacts: Artifact[], afxIds: ei.ArtifactSpec.Name[]): Effect[] {
  const deltas = [];
  for (const artifact of artifacts) {
    const effectMultiplier = artifact.clarityEffect;
    if (effectMultiplier === 0) {
      continue;
    }
    if (afxIds.includes(artifact.afxId)) {
      deltas.push({
        delta: artifact.effectDelta,
        multiplier: effectMultiplier,
      });
    }
    for (const stone of artifact.stones) {
      if (afxIds.includes(stone.afxId)) {
        deltas.push({
          delta: stone.effectDelta,
          multiplier: effectMultiplier,
        });
      }
    }
  }
  return deltas;
}

export function aggregateEffect(
  artifacts: Artifact[],
  afxIds: ei.ArtifactSpec.Name[],
  aggregator: (aggregate: number, effect: Effect) => number,
  initialValue: number
): number {
  return gatherRelevantEffects(artifacts, afxIds).reduce(aggregator, initialValue);
}

export function additiveEffect(artifacts: Artifact[], afxIds: ei.ArtifactSpec.Name[]): number {
  return aggregateEffect(
    artifacts,
    afxIds,
    (aggregate, effect) => aggregate + effect.delta * effect.multiplier,
    0
  );
}

export function multiplicativeEffect(
  artifacts: Artifact[],
  afxIds: ei.ArtifactSpec.Name[]
): number {
  return aggregateEffect(
    artifacts,
    afxIds,
    (aggregate, effect) => aggregate * (1 + effect.delta * effect.multiplier),
    1
  );
}
