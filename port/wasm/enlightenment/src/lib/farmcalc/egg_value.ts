import { eggValueMultiplier } from "../effects";
import { ei } from "../proto";
import { Artifact, Research, ResearchInstance } from "../types";
import { farmResearches } from "./common";

interface EggValueResearch extends Research {
  compoundMultiplicatively?: boolean;
}
interface EggValueResearchInstance extends ResearchInstance {
  compoundMultiplicatively?: boolean;
}

// https://wasmegg.netlify.app/researches/
// SELECT id, name, levels AS maxLevel, per_level AS perLevel, iif(levels_compound = 'multiplicative', true, NULL) AS compoundMultiplicatively FROM research WHERE categories LIKE '%egg_value%' ORDER BY serial_id;
const availableEggValueResearches: EggValueResearch[] = [
  {
    id: "nutritional_sup",
    name: "Nutritional Supplements",
    maxLevel: 40,
    perLevel: 0.25,
  },
  {
    id: "padded_packaging",
    name: "Padded Packaging",
    maxLevel: 30,
    perLevel: 0.25,
  },
  {
    id: "bigger_eggs",
    name: "Bigger Eggs",
    maxLevel: 1,
    perLevel: 2,
    compoundMultiplicatively: true,
  },
  {
    id: "usde_prime",
    name: "USDE Prime Certification",
    maxLevel: 1,
    perLevel: 3,
    compoundMultiplicatively: true,
  },
  {
    id: "superfeed",
    name: "Super-Feed™ Diet",
    maxLevel: 35,
    perLevel: 0.25,
  },
  {
    id: "improved_genetics",
    name: "Improved Genetics",
    maxLevel: 30,
    perLevel: 0.15,
  },
  {
    id: "shell_fortification",
    name: "Shell Fortification",
    maxLevel: 60,
    perLevel: 0.15,
  },
  {
    id: "even_bigger_eggs",
    name: "Even Bigger Eggs",
    maxLevel: 5,
    perLevel: 2,
    compoundMultiplicatively: true,
  },
  {
    id: "genetic_purification",
    name: "Genetic Purification",
    maxLevel: 100,
    perLevel: 0.1,
  },
  {
    id: "graviton_coating",
    name: "Graviton Coating",
    maxLevel: 7,
    perLevel: 2,
    compoundMultiplicatively: true,
  },
  {
    id: "chrystal_shells",
    name: "Crystalline Shelling",
    maxLevel: 100,
    perLevel: 0.25,
  },
  {
    id: "telepathic_will",
    name: "Telepathic Will",
    maxLevel: 50,
    perLevel: 0.25,
  },
  {
    id: "atomic_purification",
    name: "Atomic Purification",
    maxLevel: 50,
    perLevel: 0.1,
  },
  {
    id: "multi_layering",
    name: "Multiversal Layering",
    maxLevel: 3,
    perLevel: 10,
    compoundMultiplicatively: true,
  },
  {
    id: "eggsistor",
    name: "Eggsistor Miniturization",
    maxLevel: 100,
    perLevel: 0.05,
  },
  {
    id: "matter_reconfig",
    name: "Matter Reconfiguration",
    maxLevel: 500,
    perLevel: 0.01,
  },
  {
    id: "timeline_splicing",
    name: "Timeline Splicing",
    maxLevel: 1,
    perLevel: 10,
    compoundMultiplicatively: true,
  },
];

export function farmEggValueResearches(farm: ei.Backup.ISimulation): EggValueResearchInstance[] {
  return farmResearches(farm, null, availableEggValueResearches);
}

const baseEggValue = 1e-7;

export function farmEggValue(
  researches: EggValueResearchInstance[],
  artifacts: Artifact[]
): number {
  let eggValue = baseEggValue;
  for (const research of researches) {
    if (research.compoundMultiplicatively) {
      eggValue *= research.perLevel ** research.level;
    } else {
      eggValue *= 1 + research.level * research.perLevel;
    }
  }
  eggValue *= eggValueMultiplier(artifacts);
  return eggValue;
}
