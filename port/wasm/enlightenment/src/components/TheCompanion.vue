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
      class="flex items-center space-x-1"
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
        <template v-if="completionForecast">
          <span class="text-green-500 whitespace-nowrap">
            {{ completionForecast.format("LLL z") }}
          </span>
          <template v-if="completionForecastDays > 0">
            ({{ completionForecastDays.toFixed(1) }} days)
          </template>
        </template>
        <template v-else>Never</template>
      </p>

      <hr class="mt-2" />

      <collapsible-section
        sectionTitle="Habs"
        :visible="isVisibleSection('habs')"
        @toggle="toggleSectionVisibility('habs')"
        class="my-2 text-sm"
      >
        <div class="flex my-2 space-x-2">
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
        <template v-if="!totalHabSpaceSufficient">
          <p>
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
          <template v-if="minimumRequiredWDLevel < requiredWDLevel">
            <p>
              Note that the level above assumes your current set of artifacts. The minimum WD level
              required is
              <span class="text-blue-500 mr-0.5">{{ minimumRequiredWDLevel }}/25</span>, assuming
              you equip your most effective gusset as pictured below (stone rearrangement possibly
              needed):
            </p>
            <artifacts-gallery :artifacts="bestPossibleGussetSet" class="mt-2 mb-3" />
          </template>
        </template>
      </collapsible-section>

      <hr />

      <collapsible-section
        sectionTitle="Earnings"
        :visible="isVisibleSection('earnings')"
        @toggle="toggleSectionVisibility('earnings')"
        class="my-2 text-sm"
      >
        <p>
          Earning bonus:
          <base-e-i-value class="text-green-500" :value="earningBonus * 100" suffix="%" />,
          <span class="whitespace-nowrap" :style="{ color: farmerRole.color }">{{
            farmerRole.name
          }}</span>
        </p>
        <p>Farm value: <base-e-i-value class="text-green-500" :value="farmValue" /></p>
        <p>Cash on hand: <base-e-i-value class="text-green-500" :value="cashOnHand" /></p>
        <p>Egg value: <base-e-i-value class="text-green-500" :value="eggValue" /></p>
        <p>
          Earning rate (active, no running chicken):
          <base-e-i-value class="text-green-500" :value="earningRateOnlineBaseline" suffix="/s" />
        </p>
        <p>
          Earning rate (active, max RCB <span class="text-green-500">{{ maxRCB }}x</span>):
          <base-e-i-value class="text-green-500" :value="earningRateOnlineMaxRCB" suffix="/s" />
        </p>
        <p>
          Earning rate (offline):
          <base-e-i-value class="text-green-500" :value="earningRateOffline" suffix="/s" />
        </p>
        <p class="mt-1">Drone values at max RCB:</p>
        <ul>
          <li>
            Elite: <base-e-i-value class="text-green-500" :value="droneValuesAtMaxRCB.elite" />
          </li>
          <li>
            Regular tier 1:
            <base-e-i-value class="text-green-500" :value="droneValuesAtMaxRCB.tier1" /> ({{
              formatPercentage(droneValuesAtMaxRCB.tier1 / droneValuesAtMaxRCB.elite)
            }}
            of elite), {{ formatPercentage(droneValuesAtMaxRCB.tier1Prob) }} chance
          </li>
          <li>
            Regular tier 2:
            <base-e-i-value class="text-green-500" :value="droneValuesAtMaxRCB.tier2" />({{
              formatPercentage(droneValuesAtMaxRCB.tier2 / droneValuesAtMaxRCB.elite)
            }}
            of elite),
            {{ formatPercentage(droneValuesAtMaxRCB.tier2Prob) }} chance
          </li>
          <li>
            Regular tier 3:
            <base-e-i-value class="text-green-500" :value="droneValuesAtMaxRCB.tier3" />({{
              formatPercentage(droneValuesAtMaxRCB.tier3 / droneValuesAtMaxRCB.elite)
            }}
            of elite),
            {{ formatPercentage(droneValuesAtMaxRCB.tier3Prob) }} chance
          </li>
        </ul>
        <p class="text-xs text-gray-500 my-1">
          Drone values are based on your current equipped set of artifacts. You may increase their
          values with an Aurelian brooch or a Mercury's lens (drone reward is proportional to farm
          value) or a Vial of Martian dust / terra stones (drone reward is proportional to the
          square root of active running chicken bonus, so increasing max RCB helps to a small
          extend); or increase their frequency with a Neodymium medallion. Farming drones during a
          Generous Drones event is also immensely helpful.
        </p>
        <p class="text-xs text-gray-500 my-1">
          Note: Farm value and drone values are calculated based on mikit#7826's research on game
          version v1.12.13 (pre-artifacts), and drone probabilities were speculative at that time.
          No in-depth research has been carried out since the artifact update, so values may be
          inaccurate in certain edge cases.
        </p>

        <template v-if="cashTargetPreDiscount > 0">
          <p>
            Cash required to reach minimum required Wormhole Dampening level
            <span class="text-blue-500">{{ minimumRequiredWDLevel }}/25</span>
            (before discounts):
            <base-e-i-value class="text-pink-500" :value="cashTargetPreDiscount" />
          </p>
          <target-cash-matrix
            :baseTarget="cashTargetPreDiscount"
            :current="cashOnHand"
            :targets="cashTargets"
            :means="cashMeans"
            class="my-2"
          />
          <template v-if="betterCubePossible">
            <p>Your best cube possible is pictured below (stone rearrangement possibly needed):</p>
            <artifacts-gallery :artifacts="bestPossibleCubeSet" class="mt-2 mb-3" />
          </template>
          <div class="text-sm mt-2">
            <a
              href="https://docs.google.com/spreadsheets/d/157K4r3Z5wfCNKhUWb34mlxM08DEA1AWamsA20xjQIhw/edit?usp=sharing"
              target="_blank"
              class="text-blue-500 hover:text-blue-600"
              >Sami#2336's spreadsheet</a
            >
            may provide more detailed help regarding execution, at the expense of requiring manual
            input for many parameters.
          </div>
        </template>
      </collapsible-section>

      <hr />

      <collapsible-section
        sectionTitle="Internal hatchery"
        :visible="isVisibleSection('internal_hatchery')"
        @toggle="toggleSectionVisibility('internal_hatchery')"
        class="my-2 text-sm"
      >
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
      </collapsible-section>

      <hr />

      <collapsible-section
        sectionTitle="Artifacts"
        :visible="isVisibleSection('artifacts')"
        @toggle="toggleSectionVisibility('artifacts')"
        class="my-2 text-sm"
      >
        <artifacts-gallery :artifacts="artifacts" />
      </collapsible-section>
    </template>
  </main>
</template>

<script lang="ts">
import { computed, defineComponent, onBeforeUnmount, ref } from "vue";
import dayjs, { Dayjs } from "dayjs";
import advancedFormat from "dayjs/plugin/advancedFormat";
import localizedFormat from "dayjs/plugin/localizedFormat";
import relativeTime from "dayjs/plugin/relativeTime";
import timezone from "dayjs/plugin/timezone";
import utc from "dayjs/plugin/utc";

import {
  bestPossibleCubeForEnlightenment,
  bestPossibleGussetForEnlightenment,
  calculateDroneValues,
  calculateFarmValue,
  calculateWDLevelsCost,
  earningBonusToFarmerRole,
  eggIconPath,
  ei,
  farmCurrentWDLevel,
  farmEarningBonus,
  farmEarningRate,
  farmEggValue,
  farmEggValueResearches,
  farmHabs,
  farmHabSpaceResearches,
  farmHabSpaces,
  farmInternalHatcheryRates,
  farmInternalHatcheryResearches,
  farmMaxRCB,
  farmMaxRCBResearches,
  homeFarmArtifacts,
  requestFirstContact,
  requiredWDLevelForEnlightenmentDiamond,
  researchPriceMultiplierFromArtifacts,
  researchPriceMultiplierFromResearches,
} from "@/lib";
import {
  iconURL,
  formatPercentage,
  formatWithThousandSeparators,
  formatDurationAuto,
  getLocalStorage,
  setLocalStorage,
} from "@/utils";
import CollapsibleSection from "@/components/CollapsibleSection.vue";
import ArtifactsGallery from "@/components/ArtifactsGallery.vue";
import UnfinishedResearches from "@/components/UnfinishedResearches.vue";
import TargetCashMatrix from "@/components/TargetCashMatrix.vue";
import BaseInfo from "@/components/BaseInfo.vue";
import BaseEIValue from "@/components/BaseEIValue.vue";

// Note that timezone abbreviation may not work due to
// https://github.com/iamkun/dayjs/issues/1154, in which case the GMT offset is
// shown.
dayjs.extend(advancedFormat);
dayjs.extend(localizedFormat);
dayjs.extend(relativeTime);
dayjs.extend(timezone);
dayjs.extend(utc);

export default defineComponent({
  components: {
    CollapsibleSection,
    ArtifactsGallery,
    UnfinishedResearches,
    TargetCashMatrix,
    BaseInfo,
    BaseEIValue,
  },
  props: {
    playerId: {
      type: String,
      required: true,
    },
  },
  // This async component does not respond to playerId changes.
  async setup({ playerId }) {
    // Validate and sanitize player ID.
    if (!playerId.match(/^EI\d+$/i)) {
      throw new Error(
        `ID ${playerId} is not in the form EI1234567890123456; please consult "Where do I find my ID?"`
      );
    }
    playerId = playerId.toUpperCase();

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
    const currentWDLevel = farmCurrentWDLevel(farm);
    const requiredWDLevel = requiredWDLevelForEnlightenmentDiamond(artifacts);
    const bestPossibleGusset = bestPossibleGussetForEnlightenment(backup);
    const bestPossibleGussetSet = bestPossibleGusset ? [bestPossibleGusset] : [];
    const minimumRequiredWDLevel = bestPossibleGusset
      ? requiredWDLevelForEnlightenmentDiamond([bestPossibleGusset])
      : requiredWDLevel;

    const earningBonus = farmEarningBonus(backup, farm, progress, artifacts);
    const farmerRole = earningBonusToFarmerRole(earningBonus);
    const farmValue = calculateFarmValue(backup, farm, progress, artifacts);
    const cashOnHand = farm.cashEarned! - farm.cashSpent!;
    const eggValue = farmEggValue(farmEggValueResearches(farm), artifacts);
    const maxRCB = farmMaxRCB(farmMaxRCBResearches(farm, progress), artifacts);
    const {
      onlineBaseline: earningRateOnlineBaseline,
      onlineMaxRCB: earningRateOnlineMaxRCB,
      offline: earningRateOffline,
    } = farmEarningRate(backup, farm, progress, artifacts);
    const droneValuesAtMaxRCB = calculateDroneValues(farm, progress, artifacts, {
      population: farm.numChickens! as number,
      farmValue,
      rcb: maxRCB,
    });
    const cashTargetPreDiscount =
      calculateWDLevelsCost(currentWDLevel, minimumRequiredWDLevel) *
      researchPriceMultiplierFromResearches(farm, progress);
    const currentPriceMultiplier = researchPriceMultiplierFromArtifacts(artifacts);
    const bestPossibleCube = bestPossibleCubeForEnlightenment(backup);
    const bestPossibleCubeSet = bestPossibleCube ? [bestPossibleCube] : [];
    const bestPriceMultiplier = researchPriceMultiplierFromArtifacts(bestPossibleCubeSet);
    const cashTargets = [
      { multiplier: 1, description: "No research sale\nno artifacts" },
      { multiplier: 0.35, description: "65% research sale\n no artifacts" },
    ];
    if (currentPriceMultiplier < 1) {
      cashTargets.push(
        { multiplier: currentPriceMultiplier, description: "No research sale\ncurrent artifacts" },
        {
          multiplier: currentPriceMultiplier * 0.35,
          description: "65% research sale\ncurrent artifacts",
        }
      );
    }
    const betterCubePossible = bestPriceMultiplier < currentPriceMultiplier;
    if (betterCubePossible) {
      cashTargets.push(
        { multiplier: bestPriceMultiplier, description: "No research sale\nbest cube possible" },
        {
          multiplier: bestPriceMultiplier * 0.35,
          description: "65% research sale\nbest cube possible",
        }
      );
    }
    const calculateAndFormatDuration = (target: number, rate: number): string => {
      if (target <= 0) {
        return "-";
      }
      return formatDurationAuto(target / rate);
    };
    const calculateAndFormatNumDrones = (target: number, rate: number): string => {
      let count: number | string;
      if (target <= 0) {
        count = 0;
      } else if (rate === 0) {
        count = "\u221E";
      } else {
        count = Math.ceil(target / rate);
      }
      return `\u00D7${count}`;
    };
    const cashMeans = [
      {
        rate: earningRateOnlineMaxRCB,
        description: "Active earnings at max RCB",
        calc: calculateAndFormatDuration,
      },
      {
        rate: earningRateOffline,
        description: "Offline earnings",
        calc: calculateAndFormatDuration,
      },
      {
        rate: droneValuesAtMaxRCB.elite,
        description: "Elite drone at max RCB",
        calc: calculateAndFormatNumDrones,
      },
    ];

    const internalHatcheryResearches = farmInternalHatcheryResearches(farm, progress);
    const {
      onlineRatePerHab: onlineIHRPerHab,
      onlineRate: onlineIHR,
      offlineRate: offlineIHR,
    } = farmInternalHatcheryRates(internalHatcheryResearches, artifacts);

    const lastRefreshedPopulation = farm.numChickens! as number;
    const targetPopulation = 1e10;
    const completed = lastRefreshedPopulation >= targetPopulation;
    const now = dayjs();
    let completionForecast: Dayjs | undefined;
    let completionForecastDays = 0;
    if (!completed && offlineIHR > 0) {
      const timeToCompleteSeconds =
        ((targetPopulation - lastRefreshedPopulation) / offlineIHR) * 60;
      completionForecast = dayjs(lastRefreshedTimestamp + timeToCompleteSeconds * 1000);
      completionForecastDays = completionForecast.diff(now, "day", true);
    }
    const currentPopulation = computed(
      () =>
        lastRefreshedPopulation +
        (offlineIHR / 60_000) * (currentTimestamp.value - lastRefreshedTimestamp)
    );

    const sectionVisibility = ref(loadSectionVisibilityFromLocalStorage());
    const isVisibleSection = (section: string) => {
      return sectionVisibility.value[section] !== false;
    };
    const toggleSectionVisibility = (section: string) => {
      const current = sectionVisibility.value[section] !== false;
      sectionVisibility.value[section] = !current;
      persistSectionVisibilityToLocalStorage(sectionVisibility.value);
    };

    return {
      nickname,
      lastRefreshed,
      lastRefreshedRelative,
      egg,
      eggIconURL,
      enlightenmentEgg,
      enlightenmentEggIconURL,
      artifacts,
      completed,
      completionForecast,
      completionForecastDays,
      habs,
      habSpaceResearches,
      totalHabSpace,
      totalHabSpaceSufficient,
      habSpaces,
      currentWDLevel,
      requiredWDLevel,
      bestPossibleGussetSet,
      minimumRequiredWDLevel,
      earningBonus,
      farmerRole,
      farmValue,
      cashOnHand,
      eggValue,
      maxRCB,
      earningRateOnlineBaseline,
      earningRateOnlineMaxRCB,
      earningRateOffline,
      droneValuesAtMaxRCB,
      cashTargetPreDiscount,
      cashTargets,
      cashMeans,
      betterCubePossible,
      bestPossibleCubeSet,
      internalHatcheryResearches,
      onlineIHR,
      onlineIHRPerHab,
      offlineIHR,
      lastRefreshedPopulation,
      currentPopulation,
      sectionVisibility,
      isVisibleSection,
      toggleSectionVisibility,
      formatWithThousandSeparators,
      formatPercentage,
      iconURL,
    };
  },
});

type SectionVisibility = { [section: string]: boolean };

const SECTION_VISIBILITY_LOCALSTORAGE_KEY = "sectionVisibility";

function loadSectionVisibilityFromLocalStorage(): SectionVisibility {
  const encoded = getLocalStorage(SECTION_VISIBILITY_LOCALSTORAGE_KEY) || "{}";
  try {
    const visibility: SectionVisibility = {};
    for (const [key, val] of Object.entries(JSON.parse(encoded))) {
      if (val === false) {
        visibility[key] = val;
      }
    }
    return visibility;
  } catch (e) {
    console.error(`error loading sectionVisibility from localStorage: ${e}`);
    return {};
  }
}

function persistSectionVisibilityToLocalStorage(value: SectionVisibility) {
  setLocalStorage(SECTION_VISIBILITY_LOCALSTORAGE_KEY, JSON.stringify(value));
}
</script>
