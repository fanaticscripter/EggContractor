import { habSpaceMultiplier } from "./effects";
import { ei } from "./proto";
import { Artifact, Research, ResearchInstance } from "./types";

type HabId = 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 | 12 | 13 | 14 | 15 | 16 | 17 | 18;

export function isHabId(x: number): x is HabId {
  return Number.isInteger(x) && x >= 0 && x <= 18;
}

export interface Hab {
  id: HabId;
  name: string;
  iconPath: string;
  baseHabSpace: number;
}

// https://egg-inc.fandom.com/wiki/Habitats
export const habVarieties: Hab[] = [
  {
    id: 0,
    name: "Coop",
    iconPath: "egginc/ei_hab_icon_coop.png",
    baseHabSpace: 250,
  },
  {
    id: 1,
    name: "Shack",
    iconPath: "egginc/ei_hab_icon_shack.png",
    baseHabSpace: 500,
  },
  {
    id: 2,
    name: "Super Shack",
    iconPath: "egginc/ei_hab_icon_super_shack.png",
    baseHabSpace: 1e3,
  },
  {
    id: 3,
    name: "Short House",
    iconPath: "egginc/ei_hab_icon_short_house.png",
    baseHabSpace: 2e3,
  },
  {
    id: 4,
    name: "The Standard",
    iconPath: "egginc/ei_hab_icon_the_standard.png",
    baseHabSpace: 5e3,
  },
  {
    id: 5,
    name: "Long House",
    iconPath: "egginc/ei_hab_icon_long_house.png",
    baseHabSpace: 1e4,
  },
  {
    id: 6,
    name: "Double Decker",
    iconPath: "egginc/ei_hab_icon_double_decker.png",
    baseHabSpace: 2e4,
  },
  {
    id: 7,
    name: "Warehouse",
    iconPath: "egginc/ei_hab_icon_warehouse.png",
    baseHabSpace: 5e4,
  },
  {
    id: 8,
    name: "Center",
    iconPath: "egginc/ei_hab_icon_center.png",
    baseHabSpace: 1e5,
  },
  {
    id: 9,
    name: "Bunker",
    iconPath: "egginc/ei_hab_icon_bunker.png",
    baseHabSpace: 2e5,
  },
  {
    id: 10,
    name: "Eggkea",
    iconPath: "egginc/ei_hab_icon_eggkea.png",
    baseHabSpace: 5e5,
  },
  {
    id: 11,
    name: "HAB 1000",
    iconPath: "egginc/ei_hab_icon_hab1k.png",
    baseHabSpace: 1e6,
  },
  {
    id: 12,
    name: "Hangar",
    iconPath: "egginc/ei_hab_icon_hanger.png",
    baseHabSpace: 2e6,
  },
  {
    id: 13,
    name: "Tower",
    iconPath: "egginc/ei_hab_icon_tower.png",
    baseHabSpace: 5e6,
  },
  {
    id: 14,
    name: "HAB 10,000",
    iconPath: "egginc/ei_hab_icon_hab10k.png",
    baseHabSpace: 1e7,
  },
  {
    id: 15,
    name: "Eggtopia",
    iconPath: "egginc/ei_hab_icon_eggtopia.png",
    baseHabSpace: 2.5e7,
  },
  {
    id: 16,
    name: "Monolith",
    iconPath: "egginc/ei_hab_icon_monolith.png",
    baseHabSpace: 5e7,
  },
  {
    id: 17,
    name: "Planet Portal",
    iconPath: "egginc/ei_hab_icon_portal.png",
    baseHabSpace: 1e8,
  },
  {
    id: 18,
    name: "Chicken Universe",
    iconPath: "egginc/ei_hab_icon_chicken_universe.png",
    baseHabSpace: 6e8,
  },
];

function isPortalHab(hab: Hab): boolean {
  return hab.id >= 17;
}

export interface HabSpaceResearch extends Research {
  portalHabsOnly?: boolean;
}

export interface HabSpaceResearchInstance extends ResearchInstance {
  portalHabsOnly?: boolean;
}

const availableHabSpaceResearches: HabSpaceResearch[] = [
  {
    id: "hab_capacity1",
    name: "Hen House Remodel",
    maxLevel: 8,
    perLevel: 0.05,
  },
  {
    id: "microlux",
    name: "Microlux™ Chicken Suites",
    maxLevel: 10,
    perLevel: 0.05,
  },
  {
    id: "grav_plating",
    name: "Grav Plating",
    maxLevel: 25,
    perLevel: 0.02,
  },
  {
    id: "wormhole_dampening",
    name: "Wormhole Dampening",
    maxLevel: 25,
    perLevel: 0.02,
    portalHabsOnly: true,
  },
];

export function farmHabs(farm: ei.Backup.ISimulation): Hab[] {
  const habs = [];
  for (const habId of farm.habs!) {
    if (habId === 19) {
      // 19 is the placeholder for unpurchased habs.
      continue;
    }
    if (!isHabId(habId)) {
      throw new Error(`${habId} is not a recognized hab ID`);
    }
    habs.push(habVarieties[habId]);
  }
  return habs;
}

export function farmHabSpaceResearches(farm: ei.Backup.ISimulation): HabSpaceResearchInstance[] {
  const researches: HabSpaceResearchInstance[] = [];
  for (const r of farm.commonResearch!) {
    for (const rr of availableHabSpaceResearches) {
      if (r.id === rr.id) {
        researches.push({
          ...rr,
          level: r.level!,
        });
      }
    }
  }
  return researches;
}

export function farmHabSpaces(
  habs: Hab[],
  researches: HabSpaceResearchInstance[],
  artifacts: Artifact[]
): number[] {
  let universalMultiplier = 1;
  let portalOnlyMultiplier = 1;
  for (const research of researches) {
    const multiplier = 1 + research.level * research.perLevel;
    if (research.portalHabsOnly) {
      portalOnlyMultiplier *= multiplier;
    } else {
      universalMultiplier *= multiplier;
    }
  }
  const artifactsMultiplier = habSpaceMultiplier(artifacts);
  return habs.map(hab =>
    // Each hab's capacity rounds up individually.
    Math.ceil(
      hab.baseHabSpace *
        universalMultiplier *
        (isPortalHab(hab) ? portalOnlyMultiplier : 1) *
        artifactsMultiplier
    )
  );
}

// Wormhole Dampening levels required to reach 10B hab space, assuming max
// everything else.
export function requiredWDLevelForEnlightenmentDiamond(artifacts: Artifact[]): number {
  const target = 1e10;
  const finalHab = habVarieties[habVarieties.length - 1];
  const finalHabs = [finalHab, finalHab, finalHab, finalHab];
  const maxResearchesWithoutWD: HabSpaceResearchInstance[] = availableHabSpaceResearches
    .slice(0, availableHabSpaceResearches.length - 1)
    .map(research => ({
      ...research,
      level: research.maxLevel,
    }));
  const maxHabSpaceWithoutWD = farmHabSpaces(finalHabs, maxResearchesWithoutWD, artifacts).reduce(
    (total, s) => total + s
  );
  if (maxHabSpaceWithoutWD >= target) {
    return 0;
  }
  return Math.ceil((target / maxHabSpaceWithoutWD - 1) / 0.02);
}
