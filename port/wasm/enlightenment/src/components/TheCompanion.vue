<template>
  <main>
    <p>{{ nickname }}</p>
    <p class="text-sm">
      Last synced to server:
      <span class="whitespace-nowrap" v-tippy="{ content: lastRefreshed.format('LLL') }">
        {{ lastRefreshedRelative }}
      </span>
    </p>
    <p
      class="inline-flex items-center space-x-1"
      v-tippy="{
        content: `
          <p>The game, while active, saves to Egg, Inc.&rsquo;s server every couple of minutes if network condition allows.
          Other than soon after a fresh launch of the game, such server syncs are unpredictable from the user&rsquo;s point of view.
          <span class='text-blue-300'>You can force close then reopen the app to reasonably reliably trigger a sync</span>
          (search for &ldquo;iOS force close app&rdquo; or &ldquo;Android force close app&rdquo; if you need help).</p>

          <p>However, even after an app-initiated sync, it may take an unpredicatible amount of time
          (usually within a minute or so) for the game&rsquo;s server to serve the updated save through its API,
          which is then picked up by this tool. There is no solution other than clicking &ldquo;Load Player Data&rdquo;
          periodically until the fresh save shows up. Please do not refresh too fast, which is not helpful.</p>`,
        allowHTML: true,
      }"
    >
      <base-info />
      <span class="text-xs text-gray-500">Why is my save out of date?</span>
    </p>
    <p v-if="egg !== enlightenmentEgg" class="text-sm text-red-500 inline-flex items-center">
      Current egg is <img :src="eggIconURL" class="inline h-8 w-8" />, not
      <img :src="enlightenmentEggIconURL" class="inline h-8 w-8" />!
    </p>
    <template v-else>
      <p class="text-sm">
        Last save population:
        <span class="text-green-500 tabular-nums">
          {{ formatWithThousandSeparators(lastRefreshedPopulation) }}
        </span>
      </p>
      <p class="text-sm">
        Current population:
        <span class="text-green-500 tabular-nums mr-0.5">
          {{ formatWithThousandSeparators(currentPopulation) }}
        </span>
        <base-info
          class="inline relative -top-px"
          v-tippy="{
            content:
              'The current population is calculated based on the population and offline IHR from the last save. Assuming your IHR did not change since the last save, this number should be slightly ahead of your actual population at the moment, depending on how long you remained active since the last save.',
          }"
        />
      </p>
      <p v-if="!completed" class="text-sm">
        <template v-if="totalHabSpaceSufficient">Enlightenment Diamond Trophy forecast: </template>
        <template v-else>
          Enlightenment Diamond Trophy forecast, assuming sufficient hab space can be unlocked in
          time:
        </template>
        <span class="text-green-500 whitespace-nowrap">
          <template v-if="completionForecast">{{ completionForecast.format("LLL") }}</template>
          <template v-else>Never</template>
        </span>
      </p>

      <hr class="mt-2" />

      <section class="my-2 text-sm">
        <h2 class="font-medium">Habs</h2>
        <div class="flex gap-2 my-2">
          <img
            v-for="(hab, index) in habs"
            :key="index"
            :src="iconURL(hab.iconPath, 128)"
            class="h-16 w-16 bg-gray-50 rounded-lg shadow"
            v-tippy="{
              content: `${hab.name}, space: ${formatWithThousandSeparators(habSpaces[index])}`,
            }"
          />
        </div>
        <p>
          Hab space:
          <span class="text-green-500">{{ formatWithThousandSeparators(totalHabSpace) }}</span>
        </p>
        <unfinished-researches :researches="habSpaceResearches" class="my-1" />
        <p v-if="!totalHabSpaceSufficient">
          Required Wormhole Dampening level:
          <span class="text-blue-500 mr-0.5">{{ requiredWDLevel }}/25</span>
          <base-info
            class="inline relative -top-px"
            v-tippy="{
              content:
                'Minimum Wormhole Dampening level to reach 10B hab space, assuming all habs are final tier, and all other hab space-related researches have been finished.',
            }"
          />
        </p>
      </section>

      <hr />

      <section class="my-2 text-sm">
        <h2 class="font-medium">Internal hatchery</h2>
        <p class="mt-1">
          Active IHR:
          <span class="whitespace-nowrap">
            <span class="text-green-500">{{ formatWithThousandSeparators(onlineIHR, -1) }}</span>
            chickens/min
          </span>
          <!-- Force a space between the two nowrap spans to prevent the two being treated as a whole. -->
          {{ " " }}
          <span class="whitespace-nowrap">
            (<span class="text-green-500">{{
              formatWithThousandSeparators(onlineIHRPerHab, -1)
            }}</span>
            chickens/min/hab)
          </span>
        </p>
        <p>
          Offline IHR:
          <span class="text-green-500">{{ formatWithThousandSeparators(offlineIHR, -1) }}</span>
          chickens/min
        </p>
        <unfinished-researches :researches="internalHatcheryResearches" class="my-1" />
      </section>

      <hr />

      <section class="my-2 text-sm">
        <h2 class="font-medium mb-2">Artifacts</h2>
        <artifacts-gallery :artifacts="artifacts" />
      </section>
    </template>
  </main>
</template>

<script lang="ts">
import { computed, defineComponent, onBeforeUnmount, ref } from "vue";
import dayjs, { Dayjs } from "dayjs";
import localizedFormat from "dayjs/plugin/localizedFormat";
import relativeTime from "dayjs/plugin/relativeTime";

import {
  eggIconPath,
  ei,
  farmHabs,
  farmHabSpaces,
  farmHabSpaceResearches,
  homeFarmArtifacts,
  requestFirstContact,
  requiredWDLevelForEnlightenmentDiamond,
  farmInternalHatcheryResearches,
  farmInternalHatcheryRates,
} from "@/lib";
import { iconURL } from "@/utils";
import ArtifactsGallery from "./ArtifactsGallery.vue";
import UnfinishedResearches from "./UnfinishedResearches.vue";
import BaseInfo from "./BaseInfo.vue";

dayjs.extend(localizedFormat);
dayjs.extend(relativeTime);

export default defineComponent({
  components: {
    ArtifactsGallery,
    UnfinishedResearches,
    BaseInfo,
  },
  props: {
    playerId: {
      type: String,
      required: true,
    },
  },
  // This async component does not respond to playerId changes.
  async setup({ playerId }) {
    // Interval id used for refreshing lastRefreshedRelative.
    let refreshIntervalId: number | undefined;
    onBeforeUnmount(() => {
      clearInterval(refreshIntervalId);
    });

    const data = await requestFirstContact(playerId);
    if (!data.backup || !data.backup.game) {
      throw new Error(`${playerId}: backup is empty`);
    }
    const backup = data.backup;
    const nickname = backup.userName;
    const progress = backup.game;
    if (!progress) {
      throw new Error(`${playerId}: no game progress in backup`);
    }
    if (!backup.farms || backup.farms.length === 0) {
      throw new Error(`${playerId}: no farm info in backup`);
    }

    const farm = backup.farms[0]; // Home farm
    const egg = farm.eggType!;
    const eggIconURL = iconURL(eggIconPath(egg), 128);
    const enlightenmentEgg = ei.Egg.ENLIGHTENMENT;
    const enlightenmentEggIconURL = iconURL(eggIconPath(enlightenmentEgg), 128);
    const lastRefreshedTimestamp = farm.lastStepTime! * 1000;
    const lastRefreshed = dayjs(Math.min(lastRefreshedTimestamp, Date.now()));
    const currentTimestamp = ref(Date.now());
    const lastRefreshedRelative = ref(lastRefreshed.fromNow());
    const artifacts = homeFarmArtifacts(backup);

    refreshIntervalId = setInterval(() => {
      currentTimestamp.value = Date.now();
      lastRefreshedRelative.value = lastRefreshed.fromNow();
    }, 200);

    const habs = farmHabs(farm);
    const habSpaceResearches = farmHabSpaceResearches(farm);
    const habSpaces = farmHabSpaces(habs, habSpaceResearches, artifacts);
    const totalHabSpace = Math.round(habSpaces.reduce((total, s) => total + s));
    const totalHabSpaceSufficient = totalHabSpace >= 1e10;
    const requiredWDLevel = requiredWDLevelForEnlightenmentDiamond(artifacts);

    const internalHatcheryResearches = farmInternalHatcheryResearches(farm, progress);
    const {
      onlineRatePerHab: onlineIHRPerHab,
      onlineRate: onlineIHR,
      offlineRate: offlineIHR,
    } = farmInternalHatcheryRates(internalHatcheryResearches, artifacts);

    const lastRefreshedPopulation = farm.numChickens! as number;
    const targetPopulation = 1e10;
    const completed = lastRefreshedPopulation >= targetPopulation;
    let completionForecast: Dayjs | undefined;
    if (!completed && offlineIHR > 0) {
      const timeToCompleteSeconds =
        ((targetPopulation - lastRefreshedPopulation) / offlineIHR) * 60;
      completionForecast = dayjs(lastRefreshedTimestamp + timeToCompleteSeconds * 1000);
    }
    const currentPopulation = computed(
      () =>
        lastRefreshedPopulation +
        (offlineIHR / 60_000) * (currentTimestamp.value - lastRefreshedTimestamp)
    );

    return {
      nickname,
      lastRefreshed,
      lastRefreshedRelative,
      egg,
      eggIconURL,
      enlightenmentEgg,
      enlightenmentEggIconURL,
      artifacts,
      habs,
      habSpaceResearches,
      totalHabSpace,
      totalHabSpaceSufficient,
      habSpaces,
      requiredWDLevel,
      internalHatcheryResearches,
      onlineIHR,
      onlineIHRPerHab,
      offlineIHR,
      lastRefreshedPopulation,
      currentPopulation,
      completed,
      completionForecast,
      formatWithThousandSeparators,
      iconURL,
    };
  },
});

enum RoundingMode {
  Down = -1,
  Nearest = 0,
  Up = 1,
}

function formatWithThousandSeparators(x: number, roundingMode = RoundingMode.Nearest): string {
  let rounded: number;
  switch (roundingMode) {
    case RoundingMode.Down:
      rounded = Math.floor(x);
      break;
    case RoundingMode.Nearest:
      rounded = Math.round(x);
      break;
    case RoundingMode.Up:
      rounded = Math.ceil(x);
      break;
  }
  return rounded.toLocaleString("en-US");
}
</script>
