import { expose } from "comlink";
import seedrandom, { prng } from "seedrandom";

import { ItemSpec, SimulationWorkerInterface } from "@/types";
import { items, missionIdToMission } from "@/data";

const ABORT_CHECK_INTERVAL = 200;
const REPORT_INTERVAL = 200;

const numItems = items.length;

const indices = [...new Array(numItems).keys()];
const idToNumeric = new Map(indices.map(i => [items[i].id, i]));
const numericToId = new Map(indices.map(i => [i, items[i].id]));
const recipes = indices.map(i => {
  const recipe = items[i].recipe;
  if (recipe === null || recipe.length === 0) {
    return null;
  }
  return recipe.map(ingredient => ({
    item: idToNumeric.get(ingredient.id)!,
    count: ingredient.count,
  }));
});

function shuffle<T>(arr: T[], rng: prng.Prng) {
  let index = arr.length;
  while (index > 1) {
    let randomIndex = Math.floor(rng() * index);
    index--;
    [arr[index], arr[randomIndex]] = [arr[randomIndex], arr[index]];
  }
}

const exposed: SimulationWorkerInterface = {
  ping() {},

  async runSimulations(missions, targets, totalTrials, report, seed, signal?) {
    // Note: for performance reasons, all code in the hot path is inlined
    // instead of being split into functional units.
    const rng = seedrandom.xor128(seed);

    // Making recipes local has a huge effect on performance, even when recipes
    // is in a never accessed branch.
    const localRecipes = [...recipes];

    const missionParams = [];
    for (const m of missions) {
      const missionCount = m.count;
      const mission = missionIdToMission.get(m.id)!;
      const missionCapacity = mission.capcity;
      const missionLootItems = mission.loot.map(item => ({
        item: idToNumeric.get(item.id)!,
        count: item.count,
      }));
      // Shuffle loot items to cancel out any systematic bias in the PRNG or our
      // usage of it.
      shuffle(missionLootItems, rng);

      const missionLootItemsCount = missionLootItems.length;
      const missionLootItemIndexToNumericId = missionLootItems.map(item => item.item);
      const thresholds = [];
      let partialSum = 0;
      for (const item of missionLootItems) {
        thresholds.push(partialSum);
        partialSum += item.count;
      }
      const modulus = partialSum;
      const rejectionThreshold = 2147483648 - (2147483648 % modulus);
      missionParams.push({
        missionCount,
        missionCapacity,
        missionLootItemsCount,
        missionLootItemIndexToNumericId,
        modulus,
        rejectionThreshold,
        thresholds,
      });
    }

    const numericTargets = new Map(targets.map(({ id, count }) => [idToNumeric.get(id)!, count]));

    let finishedTrials = 0;
    let successfulTrials = 0;
    const counter: number[] = new Array(numItems);
    const missingTally: number[] = new Array(numItems).fill(0);
    const timeStart = performance.now();
    let timeLastReport = timeStart;
    let timeLastAbortCheck = timeStart;
    while (finishedTrials < totalTrials) {
      counter.fill(0);
      for (const params of missionParams) {
        const {
          missionCount,
          missionCapacity,
          missionLootItemsCount,
          missionLootItemIndexToNumericId,
          modulus,
          rejectionThreshold,
          thresholds,
        } = params;
        for (let i = 0; i < missionCount; i++) {
          let j = 0;
          while (j < missionCapacity) {
            // Generate random unsigned int32.
            let rand = rng.int32();
            if (rand < 0) {
              rand += 2147483648;
            }

            if (rand >= rejectionThreshold) {
              continue;
            }

            // Binary search for item.
            const x = rand % modulus;
            let index = undefined;
            let l = 0;
            let r = missionLootItemsCount;
            while (l < r - 1) {
              const m = (l + r) >> 1;
              const threshold = thresholds[m];
              if (threshold === x) {
                index = m;
                break;
              }
              if (threshold > x) {
                r = m;
              } else {
                l = m;
              }
            }
            if (index === undefined) {
              index = l;
            }

            counter[missionLootItemIndexToNumericId[index]]++;
            j++;
          }
        }
      }

      const requirements = new Map(numericTargets);
      let hasMissing = false;
      while (true) {
        const next = requirements.entries().next();
        if (next.done) {
          break;
        }
        const [item, requiredCount] = next.value;
        requirements.delete(item);
        const hasCount = counter[item];
        if (hasCount >= requiredCount) {
          counter[item] -= requiredCount;
        } else {
          counter[item] = 0;
          const deficitCount = requiredCount - hasCount;
          const recipe = localRecipes[item];
          if (recipe === null) {
            hasMissing = true;
            missingTally[item] += deficitCount;
            continue;
          }
          for (const ingredient of recipe) {
            requirements.set(
              ingredient.item,
              (requirements.get(ingredient.item) || 0) + deficitCount * ingredient.count
            );
          }
        }
      }
      finishedTrials++;
      if (!hasMissing) {
        successfulTrials++;
      }

      const now = performance.now();
      if (now - timeLastAbortCheck >= ABORT_CHECK_INTERVAL) {
        const aborted = await signal?.aborted;
        if (aborted) {
          break;
        }
        timeLastAbortCheck = now;
      }
      if (now - timeLastReport >= REPORT_INTERVAL) {
        report({
          totalTrials,
          finishedTrials,
          successfulTrials,
          secondsElapsed: (now - timeStart) / 1000,
          stopped: false,
        });
        timeLastReport = now;
      }
    }
    const failedTrials = finishedTrials - successfulTrials;
    report(
      {
        totalTrials,
        finishedTrials,
        successfulTrials,
        secondsElapsed: (performance.now() - timeStart) / 1000,
        stopped: true,
      },
      missingTally.reduce((entries: ItemSpec[], count, index) => {
        if (count > 0) {
          entries.push({ id: numericToId.get(index)!, count: count / failedTrials });
        }
        return entries;
      }, [])
    );
  },
};

expose(exposed);
