import { ei } from "../proto";
import { Research } from "../types";

export function farmResearches<R extends Research>(
  farm: ei.Backup.ISimulation,
  progress: ei.Backup.IGame | null,
  availableResearches: R[]
): (R & { level: number })[] {
  const researches: (R & { level: number })[] = [];
  for (const r of farm.commonResearch!.concat(progress?.epicResearch || [])) {
    for (const rr of availableResearches) {
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

export function farmResearch<R extends Research>(
  farm: ei.Backup.ISimulation,
  progress: ei.Backup.IGame | null,
  availableResearch: R
): (R & { level: number }) | null {
  const researches = farmResearches(farm, progress, [availableResearch]);
  return researches.length > 0 ? researches[0] : null;
}
