import { artifactSpecToItem } from "./catalog";
import { clarityEffect } from "./effects";
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

function newArtifact(hostItem: Item, stones: Stone[]): Artifact {
  return {
    ...hostItem,
    stones,
    clarityEffect: clarityEffect(hostItem, stones),
  };
}
