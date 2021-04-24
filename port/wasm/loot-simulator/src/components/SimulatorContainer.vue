<template>
  <div class="xl:flex xl:max-w-full xl:justify-center w-full mx-auto">
    <div
      class="xl:flex xl:justify-end xl:flex-1 border-t border-b border-gray-100 xl:border-b-0 xl:border-r xl:border-gray-100"
    >
      <div class="w-full px-4 py-4 max-w-4xl mx-auto xl:mx-0">
        <div class="space-x-1">
          <button
            type="button"
            class="inline-flex items-center px-2.5 py-1 border border-transparent text-xs font-medium rounded-full shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none disabled:opacity-50 disabled:cursor-not-allowed"
            :disabled="!validForSubmission || inProgress"
            @click="start"
          >
            Start
          </button>
          <button
            type="button"
            class="inline-flex items-center px-2.5 py-1 border border-transparent text-xs font-medium rounded-full shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none disabled:opacity-50 disabled:cursor-not-allowed"
            :disabled="!inProgress"
            @click="stop"
          >
            Stop
          </button>
        </div>

        <div class="mt-4">
          <label
            for="total_trials"
            class="flex items-center gap-1 text-sm font-medium text-gray-700"
          >
            Number of trials
          </label>
          <div class="mt-1.5">
            <base-integer-input
              id="total_trials"
              name="total_trials"
              class="max-w-sm"
              v-model="totalTrials"
              :min="1"
            />
          </div>
          <p class="mt-1 text-xs text-gray-500">
            Increasing the number of trials improves accuracy of the probability estimate at the
            expense of longer running time.
          </p>
        </div>

        <div class="mt-4">
          <label class="flex items-center gap-1 text-sm font-medium text-gray-700">
            Missions to simulate
          </label>
          <simulator-missions-select class="mt-1.5" v-model="missions" />
        </div>

        <div class="mt-4">
          <label class="flex items-center gap-1 text-sm font-medium text-gray-700">
            Target items
          </label>
          <simulator-items-select class="mt-1.5" v-model="targets" />
        </div>

        <div class="mt-4">
          <label for="seed" class="flex items-center gap-1 text-sm font-medium text-gray-700">
            <span
              class="px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800"
            >
              Advanced
            </span>
            PRNG seed
          </label>
          <div class="mt-1.5">
            <input
              type="text"
              name="seed"
              id="seed"
              class="max-w-sm shadow-sm focus:ring-blue-500 focus:border-blue-500 block w-full sm:text-sm border-gray-300 rounded-md"
              spellcheck="false"
              placeholder="Leave empty to auto-seed"
              v-model.trim="seed"
            />
          </div>
          <p class="mt-1 text-xs text-gray-500">
            Using a fixed seed guarantees reproducible results given the same input parameters (note
            that a simulator update may break the guarantee). Leave empty to generate a new seed for
            each simulation.
          </p>
        </div>
      </div>
    </div>

    <div class="xl:flex xl:flex-1 bg-gray-50">
      <div class="w-full px-4 py-4 max-w-4xl mx-auto xl:mx-0">
        <div v-if="progress !== null">
          <div class="space-y-2">
            <div>
              <p class="text-sm font-medium mb-1">Simulated missions:</p>
              <ul>
                <li v-for="{ id, count } of submittedMissions" :key="id">
                  <mission-item :id="id" :count="count" />
                </li>
              </ul>
              <div class="text-sm mt-1">
                Total mission time:
                {{ formatMissionsDuration(submittedMissionsDurationSeconds) }}
                ({{ formatMissionsDuration(submittedMissionsDurationSeconds / 3) }}
                with Pro Permit)
              </div>
            </div>
            <div>
              <p class="text-sm font-medium mb-1">Target items:</p>
              <ul>
                <li v-for="{ id, count } of submittedTargets" :key="id">
                  <artifact-item :id="id" :count="count" />
                </li>
              </ul>
            </div>
            <div>
              <p class="text-sm">
                <span class="font-medium mr-1">Total trials:</span>
                <span>{{ formatWithThousandSeparators(progress.totalTrials) }}</span>
              </p>
              <p class="text-sm">
                <span class="font-medium mr-1">Random seed:</span>
                <span>{{ submittedSeed }}</span>
              </p>
            </div>
          </div>

          <simulator-progress-bar :progress="progress" class="mt-4" />
          <div class="text-sm tabular-nums mt-2">
            <span class="whitespace-nowrap"
              >{{ formatWithThousandSeparators(progress.successfulTrials) }} successful</span
            >
            /
            <span class="whitespace-nowrap"
              >{{ formatWithThousandSeparators(progress.finishedTrials) }} trials</span
            >,
            <span class="whitespace-nowrap"
              ><span class="text-green-500">{{ successRate }}</span> success rate.</span
            >
          </div>
          <div v-if="missing" class="mt-2">
            <p class="text-sm">Average missing ingredients per failed trial:</p>
            <template v-if="missing.length > 0">
              <ul class="mt-1">
                <li v-for="{ id, count } of missing" :key="id">
                  <artifact-item :id="id" :count="count" :decimals="1" />
                </li>
              </ul>
              <p class="text-xs text-gray-500 mt-1">
                Note that missing ingredients shown are always fully decomposed into primitive,
                non-craftable items.
              </p>
            </template>
            <p v-else class="text-sm mt-1">None.</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, Ref, ref } from "vue";
import { proxy, wrap } from "comlink";

import SimulationWorker from "@/worker?worker";
import { missionIdToMission } from "@/data";
import { ModuleWorkerNotSupportedError } from "@/errors";
import {
  getMissions,
  getSeed,
  getTargets,
  getTotalTrials,
  setMissions,
  setSeed,
  setTargets,
  setTotalTrials,
} from "@/storage";
import {
  ItemId,
  ItemSpec,
  MissionId,
  MissionSpec,
  ItemSelectSpec,
  MissionSelectSpec,
  SimulationProgress,
  SimulationWorkerInterface,
} from "@/types";
import SimulatorMissionsSelect from "@/components/SimulatorMissionsSelect.vue";
import SimulatorItemsSelect from "@/components/SimulatorItemsSelect.vue";
import SimulatorProgressBar from "@/components/SimulatorProgressBar.vue";
import ArtifactItem from "@/components/ArtifactItem.vue";
import MissionItem from "@/components/MissionItem.vue";
import BaseIntegerInput from "@/components/BaseIntegerInput.vue";

function dedupeTargets(targets: ItemSelectSpec[]): ItemSpec[] {
  const deduped = new Map<ItemId, number>();
  for (const { id, count } of targets) {
    if (id !== null) {
      deduped.set(id, (deduped.get(id) || 0) + count);
    }
  }
  return [...deduped.entries()].map(([id, count]) => ({ id, count }));
}

function dedupeMissions(missions: MissionSelectSpec[]): MissionSpec[] {
  const deduped = new Map<MissionId, number>();
  for (const { id, count } of missions) {
    if (id !== null) {
      deduped.set(id, (deduped.get(id) || 0) + count);
    }
  }
  return [...deduped.entries()].map(([id, count]) => ({ id, count }));
}

function formatMissionsDuration(seconds: number): string {
  if (seconds < 3600) {
    return formatMaxDecimalDigits(seconds / 60, 0) + "m";
  }
  if (seconds < 86400) {
    return formatMaxDecimalDigits(seconds / 3600, 1) + "h";
  }
  if (seconds < 86400 * 365) {
    return formatMaxDecimalDigits(seconds / 86400, 1) + "d";
  }
  const years = Math.floor(seconds / (86400 * 365));
  const days = (seconds % (86400 * 365)) / 86400;
  return `${years.toFixed(0)}yr ${days.toFixed(0)}d`;
}

function formatMaxDecimalDigits(x: number, digits: number) {
  if (digits === 0) {
    return x.toFixed(0);
  }
  let s = x.toFixed(digits);
  s = s.replace(/0+$/, "");
  s = s.replace(/\.$/, "");
  return s;
}

function formatWithThousandSeparators(x: number): string {
  return x.toLocaleString("en-US");
}

function randomSeed(): string {
  return String((Date.now() * 2903131747) % 4291103573);
}

export default defineComponent({
  components: {
    SimulatorMissionsSelect,
    SimulatorItemsSelect,
    SimulatorProgressBar,
    ArtifactItem,
    MissionItem,
    BaseIntegerInput,
  },
  async setup() {
    const worker = wrap<SimulationWorkerInterface>(new SimulationWorker());

    // Check worker actually works.
    if (
      (await Promise.race([
        worker.ping(),
        new Promise(resolve => setTimeout(resolve, 1000, "timeout")),
      ])) === "timeout"
    ) {
      throw new ModuleWorkerNotSupportedError();
    }

    const targets: Ref<ItemSelectSpec[]> = ref(getTargets());
    const missions: Ref<MissionSelectSpec[]> = ref(getMissions());
    const submittedTargets: Ref<ItemSpec[]> = ref([]);
    const submittedMissions: Ref<MissionSpec[]> = ref([]);

    const submittedMissionsDurationSeconds = computed(() => {
      let total = 0;
      for (const mission of submittedMissions.value) {
        total += missionIdToMission.get(mission.id)!.durationSeconds * mission.count;
      }
      return total;
    });

    const totalTrials = ref(getTotalTrials());
    const submittedTotalTrials = ref(0);

    const seed = ref(getSeed());
    const submittedSeed = ref("");

    const validForSubmission = computed(() => {
      const targetsToSubmit = dedupeTargets(targets.value);
      if (targetsToSubmit.length === 0) {
        return false;
      }
      for (const target of targetsToSubmit) {
        if (target.count <= 0) {
          return false;
        }
      }
      const missionsToSubmit = dedupeMissions(missions.value);
      if (missionsToSubmit.length === 0) {
        return false;
      }
      for (const mission of missionsToSubmit) {
        if (mission.count <= 0) {
          return false;
        }
      }
      if (totalTrials.value <= 0) {
        return false;
      }
      return true;
    });

    const inProgress = ref(false);
    const progress: Ref<SimulationProgress | null> = ref(null);
    const missing: Ref<ItemSpec[] | null> = ref(null);
    const successRate = computed(() => {
      const p = progress.value;
      if (p === null) {
        return 0;
      }
      if (p.finishedTrials === 0) {
        return "\u2013";
      }
      if (p.successfulTrials === 0) {
        return "0%";
      }
      const percentage = (p.successfulTrials / p.finishedTrials) * 100;
      if (percentage < 1e-6) {
        return "~0%";
      }
      const precision = Math.min(3, String(p.successfulTrials).length);
      return percentage.toPrecision(precision) + "%";
    });

    let controller: AbortController;
    const start = async () => {
      if (!validForSubmission.value) {
        return;
      }
      setTargets(targets.value);
      setMissions(missions.value);
      setTotalTrials(totalTrials.value);
      setSeed(seed.value);
      submittedTargets.value = dedupeTargets(targets.value);
      submittedMissions.value = dedupeMissions(missions.value);
      submittedTotalTrials.value = totalTrials.value;
      submittedSeed.value = seed.value !== "" ? seed.value : randomSeed();
      inProgress.value = true;
      progress.value = {
        totalTrials: submittedTotalTrials.value,
        finishedTrials: 0,
        successfulTrials: 0,
        secondsElapsed: 0,
        stopped: false,
      };
      missing.value = null;
      controller = new AbortController();
      // TODO: investigate why performance is degraded on subsequent runs
      // (maybe create a fresh worker)
      await worker.runSimulations(
        // A JSON roundtrip is required to turn vue's proxied objects, which
        // Comlink cannot clone, into plain objects.
        JSON.parse(JSON.stringify(submittedMissions.value)) as MissionSpec[],
        JSON.parse(JSON.stringify(submittedTargets.value)) as ItemSpec[],
        submittedTotalTrials.value,
        proxy(async (updatedProgress, updatedMissing) => {
          progress.value = updatedProgress;
          if (updatedMissing) {
            missing.value = updatedMissing;
          }
          if (updatedProgress.stopped) {
            inProgress.value = false;
          }
        }),
        submittedSeed.value,
        proxy(controller.signal)
      );
    };

    const stop = async () => {
      controller?.abort();
    };

    return {
      targets,
      missions,
      submittedTargets,
      submittedMissions,
      submittedMissionsDurationSeconds,
      totalTrials,
      submittedTotalTrials,
      seed,
      submittedSeed,
      validForSubmission,
      inProgress,
      progress,
      successRate,
      missing,
      start,
      stop,
      formatMissionsDuration,
      formatWithThousandSeparators,
    };
  },
});
</script>
