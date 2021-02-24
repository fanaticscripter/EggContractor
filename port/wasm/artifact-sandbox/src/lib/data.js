import data from "./data.json";

import { stringCmp } from "@/utils";

const artifacts = data.artifacts;
const stones = data.stones;

const artifactsAlphabetical = [...artifacts].sort((a1, a2) => {
  if (a1.family_id !== a2.family_id) {
    return stringCmp(a1.family_name, a2.family_name);
  }
  if (a1.tier_number !== a2.tier_number) {
    return a1.tier_number - a2.tier_number;
  }
  return a1.afx_rarity - a2.afx_rarity;
});

const stonesAlphabetical = [...stones].sort((s1, s2) => {
  if (s1.family_id !== s2.family_id) {
    return stringCmp(s1.family_name, s2.family_name);
  }
  return s1.tier_number - s2.tier_number;
});

const artifactOptions = [
  {
    id: "",
    name: "-- Select an artifact --",
    familyName: "",
    slots: 0,
  },
  ...artifactsAlphabetical.map(artifact => ({
    id: artifact.id,
    name: `${artifact.family_name}, ${artifact.tier_name} (T${artifact.tier_number}), ${artifact.rarity}`,
    familyName: artifact.family_name,
    slots: artifact.slots,
  })),
];

const artifactOptionsGrouped = artifactOptions.reduce((groups, option) => {
  if (groups.length === 0 || option.familyName !== groups[groups.length - 1].familyName) {
    groups.push({ familyName: option.familyName, options: [option] });
    return groups;
  }
  groups[groups.length - 1].options.push(option);
  return groups;
}, []);

const stoneOptions = [
  {
    id: "",
    name: "-- Select a stone --",
    familyName: "",
  },
  ...stonesAlphabetical.map(stone => ({
    id: stone.id,
    name: `${stone.family_name}, ${stone.tier_name} (T${stone.tier_number})`,
    familyName: stone.family_name,
  })),
];

const stoneOptionsGrouped = stoneOptions.reduce((groups, option) => {
  if (groups.length === 0 || option.familyName !== groups[groups.length - 1].familyName) {
    groups.push({ familyName: option.familyName, options: [option] });
    return groups;
  }
  groups[groups.length - 1].options.push(option);
  return groups;
}, []);

const id2artifactMap = new Map(artifacts.map(a => [a.id, a]));

function artifactFromId(id) {
  return id2artifactMap.get(id) || null;
}

function artifactFromAfxIdLevelRarity(afxId, afxLevel, afxRarity) {
  for (const artifact of artifacts) {
    if (
      artifact.afx_id === afxId &&
      artifact.afx_level === afxLevel &&
      artifact.afx_rarity === afxRarity
    ) {
      return artifact;
    }
  }
  return null;
}

const id2stoneMap = new Map(stones.map(s => [s.id, s]));

function stoneFromId(id) {
  return id2stoneMap.get(id) || null;
}

function stoneFromAfxIdLevel(afxId, afxLevel) {
  for (const stone of stones) {
    if (stone.afx_id === afxId && stone.afx_level === afxLevel) {
      return stone;
    }
  }
  return null;
}

export {
  artifacts,
  stones,
  artifactOptions,
  artifactOptionsGrouped,
  stoneOptions,
  stoneOptionsGrouped,
  artifactFromId,
  artifactFromAfxIdLevelRarity,
  stoneFromId,
  stoneFromAfxIdLevel,
};
