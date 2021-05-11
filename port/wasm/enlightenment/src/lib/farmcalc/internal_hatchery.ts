import { internalHatcheryRateMultiplier } from "../effects";
import { ei } from "../proto";
import { Artifact, Research, ResearchInstance } from "../types";
import { farmResearches } from "./common";

export interface InternalHatcheryResearch extends Research {
  multiplicative?: boolean;
  offlineOnly?: boolean;
}

export interface InternalHatcheryResearchInstance extends ResearchInstance {
  multiplicative?: boolean;
  offlineOnly?: boolean;
}

const availableInternalHatcheryResearches: InternalHatcheryResearch[] = [
  {
    id: "internal_hatchery1",
    name: "Internal Hatcheries",
    maxLevel: 10,
    perLevel: 2,
  },
  {
    id: "internal_hatchery2",
    name: "Internal Hatchery Upgrades",
    maxLevel: 10,
    perLevel: 5,
  },
  {
    id: "internal_hatchery3",
    name: "Internal Hatchery Expansion",
    maxLevel: 15,
    perLevel: 10,
  },
  {
    id: "internal_hatchery4",
    name: "Internal Hatchery Expansion",
    maxLevel: 30,
    perLevel: 25,
  },
  {
    id: "internal_hatchery5",
    name: "Machine Learning Incubators",
    maxLevel: 250,
    perLevel: 5,
  },
  {
    id: "neural_linking",
    name: "Neural Linking",
    maxLevel: 30,
    perLevel: 50,
  },
  {
    id: "epic_internal_incubators",
    name: "Epic Int. Hatcheries",
    maxLevel: 20,
    perLevel: 0.05,
    epic: true,
    multiplicative: true,
  },
  {
    id: "int_hatch_calm",
    name: "Internal Hatchery Calm",
    maxLevel: 20,
    perLevel: 0.1,
    epic: true,
    multiplicative: true,
    offlineOnly: true,
  },
];

export function farmInternalHatcheryResearches(
  farm: ei.Backup.ISimulation,
  progress: ei.Backup.IGame
): InternalHatcheryResearchInstance[] {
  return farmResearches(farm, progress, availableInternalHatcheryResearches);
}

// Rates are measured in chickens/min.
export function farmInternalHatcheryRates(
  researches: InternalHatcheryResearchInstance[],
  artifacts: Artifact[]
): {
  onlineRatePerHab: number;
  onlineRate: number;
  offlineRatePerHab: number;
  offlineRate: number;
} {
  let baseRate = 0;
  let onlineMultiplier = 1;
  let offlineMultiplier = 1;
  for (const research of researches) {
    if (research.multiplicative) {
      const multiplier = 1 + research.level * research.perLevel;
      if (research.offlineOnly) {
        offlineMultiplier *= multiplier;
      } else {
        onlineMultiplier *= multiplier;
      }
    } else {
      baseRate += research.level * research.perLevel;
    }
  }
  const artifactsMultiplier = internalHatcheryRateMultiplier(artifacts);
  // With max internal hatchery sharing, four internal hatcheries are constantly
  // at work even if not all habs are bought;
  const onlineRatePerHab = baseRate * onlineMultiplier * artifactsMultiplier;
  const onlineRate = 4 * onlineRatePerHab;
  const offlineRatePerHab = onlineRatePerHab * offlineMultiplier;
  const offlineRate = onlineRate * offlineMultiplier;
  return {
    onlineRatePerHab,
    onlineRate,
    offlineRatePerHab,
    offlineRate,
  };
}
