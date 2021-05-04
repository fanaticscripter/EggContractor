import { artifactSpecToItem } from "./catalog";
import { clarityEffect } from "./effects";
import { requiredWDLevelForEnlightenmentDiamond } from "./hab_space";
import { ei } from "./proto";
import { Artifact, Item, Stone } from "./types";

export function homeFarmArtifacts(backup: ei.IBackup): Artifact[] {
  const db = backup.artifactsDb;
  if (!db) {
    return [];
  }
  const inventory = db.inventoryItems;
  if (!inventory) {
    return [];
  }
  if (!db.activeArtifactSets || db.activeArtifactSets.length === 0) {
    return [];
  }
  const activeSet = db.activeArtifactSets[0];
  if (!activeSet.slots) {
    return [];
  }
  const itemIdToArtifact = new Map(inventory.map(item => [item.itemId!, item.artifact!]));
  let artifacts: Artifact[] = [];
  for (const slot of activeSet.slots) {
    if (!slot.occupied) {
      continue;
    }
    const artifact = itemIdToArtifact.get(slot.itemId!);
    if (artifact) {
      const hostItem = artifactSpecToItem(artifact.spec!);
      const stones = (artifact.stones || []).map(spec => artifactSpecToItem(spec));
      artifacts.push(newArtifact(hostItem, stones));
    }
  }
  return artifacts;
}

export function bestPossibleGussetForEnlightenment(backup: ei.IBackup): Artifact | null {
  if (!backup.artifactsDb) {
    return null;
  }
  const inventory = backup.artifactsDb.inventoryItems;
  if (!inventory) {
    return null;
  }

  const clarityStones = <Artifact[]>[];
  const recordClarityStone = (spec: ei.IArtifactSpec, count: number) => {
    if (spec.name === ei.ArtifactSpec.Name.CLARITY_STONE) {
      const stone = newArtifact(artifactSpecToItem(spec), []);
      for (let i = 0; i < count; i++) {
        clarityStones.push(stone);
      }
    }
  };

  const bareGussets = <Artifact[]>[];
  const seenGussetKeys = new Set<string>();
  const recordGusset = (spec: ei.IArtifactSpec) => {
    // Skip common gussets as they are useless for enlightenment.
    if (spec.name === ei.ArtifactSpec.Name.ORNATE_GUSSET && spec.rarity! > 0) {
      const gusset = newArtifact(artifactSpecToItem(spec), []);
      if (!seenGussetKeys.has(gusset.key)) {
        bareGussets.push(gusset);
        seenGussetKeys.add(gusset.key);
      }
    }
  };

  for (const item of inventory) {
    const count = item.quantity!;
    const spec = item.artifact!.spec!;
    recordGusset(spec);
    recordClarityStone(spec, count);
    for (const stone of item.artifact!.stones || []) {
      recordClarityStone(stone, count);
    }
  }
  if (bareGussets.length === 0) {
    return null;
  }
  // Sort gussets and clarity stones from higher to lower level, then higher to
  // lower rarity.
  bareGussets.sort((g1, g2) => g2.afxLevel - g1.afxLevel || g2.afxRarity - g1.afxRarity);
  clarityStones.sort((s1, s2) => s2.afxLevel - s1.afxLevel);

  const stonedGussets: Artifact[] = bareGussets.map(gusset =>
    newArtifact(gusset, clarityStones.slice(0, gusset.slots))
  );

  let bestGusset = stonedGussets[0];
  let minRequiredWDLevel = 25;
  for (const gusset of stonedGussets) {
    const level = requiredWDLevelForEnlightenmentDiamond([gusset]);
    if (level < minRequiredWDLevel) {
      bestGusset = gusset;
      minRequiredWDLevel = level;
    }
  }
  return bestGusset;
}

function newArtifact(hostItem: Item, stones: Stone[]): Artifact {
  return {
    ...hostItem,
    stones,
    clarityEffect: clarityEffect(hostItem, stones),
  };
}
