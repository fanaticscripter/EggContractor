import data from "./catalog.json";
import { ei } from "./proto";
import { Item } from "./types";

export const catalog: Item[] = data;
export const keyToItemMap = new Map(catalog.map(item => [item.key, item]));

export function artifactSpecToItem(spec: ei.IArtifactSpec): Item {
  return keyToItemMap.get(`${spec.name}:${spec.level}:${spec.rarity}`)!;
}
