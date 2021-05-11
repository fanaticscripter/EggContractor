import { ei } from "../proto";
import { farmResearch } from "./common";

export function researchPriceMultiplierFromResearches(
  farm: ei.Backup.ISimulation,
  progress: ei.Backup.IGame
): number {
  const research = farmResearch(farm, progress, {
    id: "cheaper_research",
    name: "Lab Upgrade",
    maxLevel: 10,
    perLevel: -0.05,
  });
  return 1 + (research ? research.level * research.perLevel : 0);
}
