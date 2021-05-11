import { artifactSpecToItem } from "../catalog";
import { clarityEffect, researchPriceMultiplierFromArtifacts } from "../effects";
import { ei } from "../proto";
import { Artifact, Item, Stone } from "../types";
import { requiredWDLevelForEnlightenmentDiamond } from "./hab_space";

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
  const stonedGussets = artifactsFromInventoryWithClarityStones(
    backup,
    ei.ArtifactSpec.Name.ORNATE_GUSSET
  );
  let bestGusset = null;
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

export function bestPossibleCubeForEnlightenment(backup: ei.IBackup): Artifact | null {
  const stonedCubes = artifactsFromInventoryWithClarityStones(
    backup,
    ei.ArtifactSpec.Name.PUZZLE_CUBE
  );
  if (stonedCubes.length === 0) {
    return null;
  }
  let bestCube = null;
  let minPriceMultiplier = 1;
  for (const cube of stonedCubes) {
    const priceMultiplier = researchPriceMultiplierFromArtifacts([cube]);
    if (priceMultiplier < minPriceMultiplier) {
      bestCube = cube;
      minPriceMultiplier = priceMultiplier;
    }
  }
  return bestCube;
}

function newArtifact(hostItem: Item, stones: Stone[]): Artifact {
  return {
    ...hostItem,
    stones,
    clarityEffect: clarityEffect(hostItem, stones),
  };
}

// Given an artifact family, returns a list of owned artifacts of that family
// slotted with the best possible clarity stones owned. Commons are skipped as
// they can't be stoned.
function artifactsFromInventoryWithClarityStones(
  backup: ei.IBackup,
  family: ei.ArtifactSpec.Name
): Artifact[] {
  if (!backup.artifactsDb) {
    return [];
  }
  const inventory = backup.artifactsDb.inventoryItems;
  if (!inventory) {
    return [];
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

  const bareArtifacts = <Artifact[]>[];
  const seenGussetKeys = new Set<string>();
  const recordArtifact = (spec: ei.IArtifactSpec) => {
    // Skip commons as they are useless for enlightenment.
    if (spec.name === family && spec.rarity! > 0) {
      const gusset = newArtifact(artifactSpecToItem(spec), []);
      if (!seenGussetKeys.has(gusset.key)) {
        bareArtifacts.push(gusset);
        seenGussetKeys.add(gusset.key);
      }
    }
  };

  for (const item of inventory) {
    const count = item.quantity!;
    const spec = item.artifact!.spec!;
    recordArtifact(spec);
    recordClarityStone(spec, count);
    for (const stone of item.artifact!.stones || []) {
      recordClarityStone(stone, count);
    }
  }
  if (bareArtifacts.length === 0) {
    return [];
  }
  // Sort artifacts and clarity stones from higher to lower level, then higher to
  // lower rarity.
  bareArtifacts.sort((g1, g2) => g2.afxLevel - g1.afxLevel || g2.afxRarity - g1.afxRarity);
  clarityStones.sort((s1, s2) => s2.afxLevel - s1.afxLevel);
  return bareArtifacts.map(gusset => newArtifact(gusset, clarityStones.slice(0, gusset.slots)));
}
