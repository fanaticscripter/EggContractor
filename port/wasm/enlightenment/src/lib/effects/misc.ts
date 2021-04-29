import { ei } from "../proto";
import { Artifact, Item, Stone } from "../types";

export function clarityEffect(hostItem: Item, stones: Stone[]): number {
  if (hostItem.afxId === ei.ArtifactSpec.Name.LIGHT_OF_EGGENDIL) {
    return 1;
  }
  let effect = 0;
  for (const stone of stones) {
    if (stone.afxId === ei.ArtifactSpec.Name.CLARITY_STONE) {
      effect += stone.effectDelta;
    }
  }
  return effect;
}

export function hasNoEffect(artifact: Artifact): boolean {
  return artifact.clarityEffect === 0;
}

export function hasIneffectiveClarityStones(artifact: Artifact): boolean {
  if (artifact.afxId !== ei.ArtifactSpec.Name.LIGHT_OF_EGGENDIL) {
    return false;
  }
  for (const stone of artifact.stones) {
    if (stone.afxId === ei.ArtifactSpec.Name.CLARITY_STONE) {
      return true;
    }
  }
  return false;
}
