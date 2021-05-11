import { droneRewardsMultiplier } from "../effects";
import { ei } from "../proto";
import { Artifact } from "../types";
import { farmResearch } from "./common";

// Based on @mikit's research:
// https://discord.com/channels/455380663013736479/455385659079917569/799780925793763379
export function calculateDroneValues(
  farm: ei.Backup.ISimulation,
  progress: ei.Backup.IGame,
  artifacts: Artifact[],
  params: {
    population: number;
    farmValue: number;
    rcb: number;
  }
): {
  tier1: number;
  tier1Prob: number;
  tier2: number;
  tier2Prob: number;
  tier3: number;
  tier3Prob: number;
  elite: number;
} {
  const { population, farmValue, rcb } = params;
  const farmValuePerChicken = farmValue / population;
  const scalingFactor =
    farmValuePerChicken < 3
      ? 30
      : farmValuePerChicken < 15
      ? 15
      : farmValuePerChicken < 30
      ? 8
      : farmValuePerChicken < 100
      ? 5
      : farmValuePerChicken < 200
      ? 2
      : farmValuePerChicken < 1_000
      ? 1
      : farmValuePerChicken < 50_000
      ? 0.75
      : farmValuePerChicken < 1_000_000
      ? 0.5
      : 0.25;
  const base = ((farmValue * rcb ** 0.5) / 50) * droneRewardsMultiplier(artifacts);
  const probMultiplier = biggerDronesProbabilityMultiplier(farm, progress);
  return {
    tier1: 0.02 * base * scalingFactor,
    tier1Prob: 1 - 0.1 * probMultiplier,
    tier2: 0.1 * base * scalingFactor,
    tier2Prob: 0.07 * probMultiplier,
    tier3: 0.5 * base * scalingFactor,
    tier3Prob: 0.03 * probMultiplier,
    elite: 10.0 * base * Math.min(10, scalingFactor),
  };
}

function biggerDronesProbabilityMultiplier(
  farm: ei.Backup.ISimulation,
  progress: ei.Backup.IGame
): number {
  const research = farmResearch(farm, progress, {
    id: "drone_rewards",
    name: "Drone Rewards",
    maxLevel: 20,
    perLevel: 0.1,
  });
  return 1 + (research ? research.level * research.perLevel : 0);
}
