import { farmValueMultiplier } from "../effects/farm_value";
import { ei } from "../proto";
import { Artifact } from "../types";
import { farmResearch } from "./common";
import { farmEarningBonus } from "./earning_bonus";
import { farmEggValue, farmEggValueResearches } from "./egg_value";
import { farmHabs, farmHabSpaceResearches, farmHabSpaces } from "./hab_space";
import { farmInternalHatcheryRates, farmInternalHatcheryResearches } from "./internal_hatchery";
import {
  farmEggLayingRate,
  farmEggLayingRatePerChicken,
  farmEggLayingRateResearches,
} from "./laying_rate";
import { farmMaxRCB, farmMaxRCBResearches } from "./max_rcb";
import { farmShippingCapacity } from "./shipping_capacity";

// https://cdn.discordapp.com/attachments/455385659079917569/798052527156494336/fv.png
// by @mikit
// Note that "MaxRunningChickenBonus" in the image should be replaced with maxRCB-4.
export function calculateFarmValue(
  backup: ei.IBackup,
  farm: ei.Backup.ISimulation,
  progress: ei.Backup.IGame,
  artifacts: Artifact[]
): number {
  // All calculations need to be done with no artifacts, except a flat bonus
  // from Mercury's lens if equipped.
  const population = farm.numChickens! as number;
  const shippingCapacity = farmShippingCapacity(farm, progress, []);
  const eggLayingRate = farmEggLayingRate(farm, progress, []);
  const populationEffective = Math.floor(
    population * Math.min(1, shippingCapacity / eggLayingRate)
  );
  const populationUndeliverable = population - populationEffective;
  const totalHabCapacity = farmHabSpaces(farmHabs(farm), farmHabSpaceResearches(farm), []).reduce(
    (total, s) => total + s
  );
  const populationVacant = Math.max(totalHabCapacity - population, 0);
  const internalHatcheryRate = farmInternalHatcheryRates(
    farmInternalHatcheryResearches(farm, progress),
    []
  ).onlineRate;
  const populationProjected = internalHatcheryRate * maxAwayTime(farm, progress);
  const eggConstMultiplier = 20; // 20 for the enlightenment egg.
  const eggValue = farmEggValue(farmEggValueResearches(farm), []);
  const eggLayingRatePerChicken = farmEggLayingRatePerChicken(
    farmEggLayingRateResearches(farm, progress),
    []
  );
  const earningBonus = farmEarningBonus(backup, farm, progress, []);
  const maxRCB = farmMaxRCB(farmMaxRCBResearches(farm, progress), []);
  return (
    30000 *
    accountingTrickMultiplier(farm, progress) *
    eggValue *
    eggLayingRatePerChicken *
    (earningBonus + 1) *
    (maxRCB - 4) ** 0.25 *
    eggConstMultiplier *
    (populationEffective +
      0.2 * populationUndeliverable +
      populationVacant ** 0.6 +
      0.25 * populationProjected) *
    farmValueMultiplier(artifacts)
  );
}

// Max away time, in minutes.
function maxAwayTime(farm: ei.Backup.ISimulation, progress: ei.Backup.IGame): number {
  let awayTimePerSilo = 60; // Unupgraded silo can last one hour.
  const research = farmResearch(farm, progress, {
    id: "silo_capacity",
    name: "Silo Capacity",
    maxLevel: 20,
    perLevel: 6,
  });
  if (research) {
    awayTimePerSilo += research.level * research.perLevel;
  }
  return farm.silosOwned! * awayTimePerSilo;
}

function accountingTrickMultiplier(farm: ei.Backup.ISimulation, progress: ei.Backup.IGame): number {
  const research = farmResearch(farm, progress, {
    id: "accounting_tricks",
    name: "Accounting Tricks",
    maxLevel: 20,
    perLevel: 0.05,
  });
  return 1 + (research ? research.level * research.perLevel : 0);
}
