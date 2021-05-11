import { awayEarningsMultiplier } from "../effects";
import { ei } from "../proto";
import { Artifact } from "../types";
import { farmEarningBonus } from "./earning_bonus";
import { farmEggValue, farmEggValueResearches } from "./egg_value";
import { farmEggLayingRate } from "./laying_rate";
import { farmMaxRCB, farmMaxRCBResearches } from "./max_rcb";
import { farmShippingCapacity } from "./shipping_capacity";

export function farmEarningRate(
  backup: ei.IBackup,
  farm: ei.Backup.ISimulation,
  progress: ei.Backup.IGame,
  artifacts: Artifact[]
): {
  onlineBaseline: number;
  onlineMaxRCB: number;
  offline: number;
} {
  const eggValue = farmEggValue(farmEggValueResearches(farm), artifacts);
  const eggLayingRate = farmEggLayingRate(farm, progress, artifacts);
  const shippingCapacity = farmShippingCapacity(farm, progress, artifacts);
  const earningBonus = farmEarningBonus(backup, farm, progress, artifacts);
  const onlineBaseline = eggValue * Math.min(eggLayingRate, shippingCapacity) * earningBonus;
  const maxRCB = farmMaxRCB(farmMaxRCBResearches(farm, progress), artifacts);
  const onlineMaxRCB = onlineBaseline * maxRCB;
  // Standard permit earnings halved while offline.
  const offline =
    onlineBaseline * (progress.permitLevel === 1 ? 1 : 0.5) * awayEarningsMultiplier(artifacts);
  return {
    onlineBaseline,
    onlineMaxRCB,
    offline,
  };
}
