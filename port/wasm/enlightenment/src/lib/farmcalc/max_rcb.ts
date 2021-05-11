import { maxRCBBonusFromArtifacts } from "../effects/max_rcb";
import { ei } from "../proto";
import { Artifact, Research, ResearchInstance } from "../types";
import { farmResearches } from "./common";

export interface MaxRCBResearch extends Research {}
export interface MaxRCBResearchInstance extends ResearchInstance {}

const availableMaxRCBResearches: MaxRCBResearch[] = [
  {
    id: "coordinated_clucking",
    name: "Coordinated Clucking",
    maxLevel: 50,
    perLevel: 0.2,
  },
  {
    id: "motivational_clucking",
    name: "Motivational Clucking",
    maxLevel: 50,
    perLevel: 0.5,
  },
  {
    id: "enlightened_chickens",
    name: "Enlightened Chickens",
    maxLevel: 150,
    perLevel: 2,
  },
  {
    id: "epic_multiplier",
    name: "Epic Multiplier",
    maxLevel: 100,
    perLevel: 2,
  },
];

export function farmMaxRCBResearches(
  farm: ei.Backup.ISimulation,
  progress: ei.Backup.IGame
): MaxRCBResearchInstance[] {
  return farmResearches(farm, progress, availableMaxRCBResearches);
}

const baseMaxRCB = 5;

export function farmMaxRCB(researches: MaxRCBResearchInstance[], artifacts: Artifact[]): number {
  let rcb = baseMaxRCB;
  for (const research of researches) {
    rcb += research.level * research.perLevel;
  }
  rcb += maxRCBBonusFromArtifacts(artifacts);
  return rcb;
}
