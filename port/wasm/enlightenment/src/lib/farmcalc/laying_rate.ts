import { eggLayingRateMultiplier } from "../effects";
import { ei } from "../proto";
import { Artifact, Research, ResearchInstance } from "../types";
import { farmResearches } from "./common";

interface EggLayingRateResearch extends Research {}

interface EggLayingRateResearchInstance extends ResearchInstance {}

// https://wasmegg.netlify.app/researches/
// SELECT id, name, levels AS maxLevel, per_level AS perLevel, iif(levels_compound = 'multiplicative', true, NULL) AS compoundMultiplicatively FROM research WHERE categories LIKE '%egg_value%' ORDER BY serial_id;
const availableEggLayingRateResearches: EggLayingRateResearch[] = [
  {
    id: "comfy_nests",
    name: "Comfortable Nests",
    maxLevel: 50,
    perLevel: 0.1,
  },
  {
    id: "hen_house_ac",
    name: "Hen House A/C",
    maxLevel: 50,
    perLevel: 0.05,
  },
  {
    id: "improved_genetics",
    name: "Improved Genetics",
    maxLevel: 30,
    perLevel: 0.15,
  },
  {
    id: "time_compress",
    name: "Time Compression",
    maxLevel: 20,
    perLevel: 0.1,
  },
  {
    id: "timeline_diversion",
    name: "Timeline Diversion",
    maxLevel: 50,
    perLevel: 0.02,
  },
  {
    id: "relativity_optimization",
    name: "Relativity Optimization",
    maxLevel: 10,
    perLevel: 0.1,
  },
  {
    id: "epic_egg_laying",
    name: "Epic Comfy Nests",
    maxLevel: 20,
    perLevel: 0.05,
  },
];

export function farmEggLayingRateResearches(
  farm: ei.Backup.ISimulation,
  progress: ei.Backup.IGame
): EggLayingRateResearchInstance[] {
  return farmResearches(farm, progress, availableEggLayingRateResearches);
}

const baseEggLayingRate = 1 / 30; // 1 egg per 30 seconds

export function farmEggLayingRatePerChicken(
  researches: EggLayingRateResearchInstance[],
  artifacts: Artifact[]
): number {
  let rate = baseEggLayingRate;
  for (const research of researches) {
    rate *= 1 + research.level * research.perLevel;
  }
  rate *= eggLayingRateMultiplier(artifacts);
  return rate;
}

// Not capped by shipping capacity.
export function farmEggLayingRate(
  farm: ei.Backup.ISimulation,
  progress: ei.Backup.IGame,
  artifacts: Artifact[]
): number {
  return (
    (farm.numChickens! as number) *
    farmEggLayingRatePerChicken(farmEggLayingRateResearches(farm, progress), artifacts)
  );
}
