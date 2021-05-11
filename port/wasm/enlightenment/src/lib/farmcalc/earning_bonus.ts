import { prophecyEggBonusFromArtifacts, soulEggBonusFromArtifacts } from "../effects";
import { ei } from "../proto";
import { Artifact, Research } from "../types";
import { farmResearch } from "./common";
import { accountProphecyEggsCount } from "./prophecy_eggs";

const baseSoulEggBonus = 0.1;
const soulFoodResearch: Research = {
  id: "soul_eggs",
  name: "Soul Food",
  maxLevel: 140,
  perLevel: 0.01,
};
const baseProphecyEggBonus = 0.05;
const prophecyBonusResearch: Research = {
  id: "prophecy_bonus",
  name: "Prophecy Bonus",
  maxLevel: 5,
  perLevel: 0.01,
};

export function farmEarningBonus(
  backup: ei.IBackup,
  farm: ei.Backup.ISimulation,
  progress: ei.Backup.IGame,
  artifacts: Artifact[]
): number {
  const soulEggsCount = (progress.soulEggsD || progress.soulEggs || 0) as number;
  const prophecyEggsCount = accountProphecyEggsCount(backup);
  const soulEggBonus = farmSoulEggBonus(farm, progress, artifacts);
  const prophecyEggBonus = farmProphecyEggBonus(farm, progress, artifacts);
  return soulEggsCount * soulEggBonus * (1 + prophecyEggBonus) ** prophecyEggsCount;
}

function farmSoulEggBonus(
  farm: ei.Backup.ISimulation,
  progress: ei.Backup.IGame,
  artifacts: Artifact[]
): number {
  const research = farmResearch(farm, progress, soulFoodResearch);
  return (
    baseSoulEggBonus +
    (research ? research.perLevel * research.level : 0) +
    soulEggBonusFromArtifacts(artifacts)
  );
}

function farmProphecyEggBonus(
  farm: ei.Backup.ISimulation,
  progress: ei.Backup.IGame,
  artifacts: Artifact[]
): number {
  const research = farmResearch(farm, progress, prophecyBonusResearch);
  return (
    baseProphecyEggBonus +
    (research ? research.perLevel * research.level : 0) +
    prophecyEggBonusFromArtifacts(artifacts)
  );
}

// Role calculate adapted from role.ts of CoopTracker.

// Implements farmer roles from the Egg, Inc. Discord.
//
// !?gethexcodes all
//
// ...
// Farmer: #ca6500
// Farmer II: #c68100
// Farmer III: #c49c00
// Kilofarmer: #c2b900
// Kilofarmer II: #afc300
// Kilofarmer III: #93c400
// Megafarmer: #75c800
// Megafarmer II: #59cd00
// Megafarmer III: #3cd300
// Gigafarmer: #1dda00
// Gigafarmer II: #00e204
// Gigafarmer III: #00eb27
// Terafarmer: #00f44e
// Terafarmer II: #00fe77
// Terafarmer III: #0cfca0
// Petafarmer: #1af7c4
// Petafarmer II: #27f4e3
// Petafarmer III: #33e0f0
// Exafarmer: #3dc4ee
// Exafarmer II: #46acec
// Exafarmer III: #4c96ea
// Zettafarmer: #5181e9
// Zettafarmer II: #546de8
// Zettafarmer III: #5557e8
// Yottafarmer: #6854e8
// Yottafarmer II: #7c51e9
// Yottafarmer III: #914dea
// Xennafarmer: #a746eb
// Xennafarmer II: #c13dee
// Xennafarmer III: #dd33f0
// Weccafarmer: #f327e9
// Weccafarmer II: #f71bcb
// Weccafarmer III: #fb0da8
// Vendafarmer: #fe007f
// Vendafarmer II: #f71bcb
// Vendafarmer III: #e132f1
// Uadafarmer: #ac44ec
// Uadafarmer II: #8150e9
// Uadafarmer III: #5a55e8
// Treidafarmer: #527ae9
// Treidafarmer II: #48a4eb
// Treidafarmer III: #38d3ef
// Quadafarmer: #21f5d5
// Quadafarmer II: #07fd93
// Quadafarmer III: #00f141
// Pendafarmer: #08df00
// Pendafarmer II: #43d100
// Pendafarmer III: #7bc700
// Exedafarmer: #b2c200
// Exedafarmer II: #c49c00
// Exedafarmer III: #ca6500
// Infinifarmer: #546e7a
// ...
//
// pbpaste | perl -lape 's/\/\/ /{name:"/; s/: /",color:"/; s/$/"},/' | pbcopy

type FarmerRole = {
  name: string;
  color: string;
};

const roles: FarmerRole[] = [
  { name: "Farmer", color: "#ca6500" },
  { name: "Farmer II", color: "#c68100" },
  { name: "Farmer III", color: "#c49c00" },
  { name: "Kilofarmer", color: "#c2b900" },
  { name: "Kilofarmer II", color: "#afc300" },
  { name: "Kilofarmer III", color: "#93c400" },
  { name: "Megafarmer", color: "#75c800" },
  { name: "Megafarmer II", color: "#59cd00" },
  { name: "Megafarmer III", color: "#3cd300" },
  { name: "Gigafarmer", color: "#1dda00" },
  { name: "Gigafarmer II", color: "#00e204" },
  { name: "Gigafarmer III", color: "#00eb27" },
  { name: "Terafarmer", color: "#00f44e" },
  { name: "Terafarmer II", color: "#00fe77" },
  { name: "Terafarmer III", color: "#0cfca0" },
  { name: "Petafarmer", color: "#1af7c4" },
  { name: "Petafarmer II", color: "#27f4e3" },
  { name: "Petafarmer III", color: "#33e0f0" },
  { name: "Exafarmer", color: "#3dc4ee" },
  { name: "Exafarmer II", color: "#46acec" },
  { name: "Exafarmer III", color: "#4c96ea" },
  { name: "Zettafarmer", color: "#5181e9" },
  { name: "Zettafarmer II", color: "#546de8" },
  { name: "Zettafarmer III", color: "#5557e8" },
  { name: "Yottafarmer", color: "#6854e8" },
  { name: "Yottafarmer II", color: "#7c51e9" },
  { name: "Yottafarmer III", color: "#914dea" },
  { name: "Xennafarmer", color: "#a746eb" },
  { name: "Xennafarmer II", color: "#c13dee" },
  { name: "Xennafarmer III", color: "#dd33f0" },
  { name: "Weccafarmer", color: "#f327e9" },
  { name: "Weccafarmer II", color: "#f71bcb" },
  { name: "Weccafarmer III", color: "#fb0da8" },
  { name: "Vendafarmer", color: "#fe007f" },
  { name: "Vendafarmer II", color: "#f71bcb" },
  { name: "Vendafarmer III", color: "#e132f1" },
  { name: "Uadafarmer", color: "#ac44ec" },
  { name: "Uadafarmer II", color: "#8150e9" },
  { name: "Uadafarmer III", color: "#5a55e8" },
  { name: "Treidafarmer", color: "#527ae9" },
  { name: "Treidafarmer II", color: "#48a4eb" },
  { name: "Treidafarmer III", color: "#38d3ef" },
  { name: "Quadafarmer", color: "#21f5d5" },
  { name: "Quadafarmer II", color: "#07fd93" },
  { name: "Quadafarmer III", color: "#00f141" },
  { name: "Pendafarmer", color: "#08df00" },
  { name: "Pendafarmer II", color: "#43d100" },
  { name: "Pendafarmer III", color: "#7bc700" },
  { name: "Exedafarmer", color: "#b2c200" },
  { name: "Exedafarmer II", color: "#c49c00" },
  { name: "Exedafarmer III", color: "#ca6500" },
  { name: "Infinifarmer", color: "#546e7a" },
];

export function earningBonusToFarmerRole(earningBonus: number): FarmerRole {
  const rank = Math.floor(Math.max(Math.log10(earningBonus), 0));
  return rank < roles.length ? roles[rank] : roles[roles.length - 1];
}
