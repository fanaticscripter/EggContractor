<template>
  <div
    id="mission-card"
    class="-mx-4 sm:mx-0 bg-gray-50 overflow-hidden sm:rounded-lg sm:shadow-md"
  >
    <div class="bg-gray-100 px-4 py-4 border-b border-gray-200 sm:px-6">
      <div class="-ml-4 -mt-2 flex items-center justify-between flex-wrap sm:flex-nowrap">
        <div class="ml-4 mt-2 flex items-center space-x-1">
          <mission-name :mission="mission" noLink="true" />
          <div class="text-sm">({{ mission.durationDisplay }})</div>
        </div>
        <div class="ml-4 mt-2 flex-shrink-0">
          <share :id="mission.id" :domElementId="'mission-card'" />
        </div>
      </div>
    </div>

    <div class="px-4 py-4 sm:px-6 space-y-2">
      <loot-data-credit />

      <div class="flex items-center">
        <label for="artifact-sort-by" class="text-xs font-medium text-gray-500 uppercase"
          >Sort by</label
        >
        <select
          id="artifact-sort-by"
          name="artifact-sort-by"
          class="ml-2 pl-2 pr-8 py-1 text-xs text-gray-500 uppercase bg-gray-50 border-gray-300 focus:outline-none focus:ring-green-500 focus:border-green-500 rounded-md"
          v-model="sortBy"
        >
          <option :value="$options.SORT_BY.QUALITY">quality</option>
          <option :value="$options.SORT_BY.FAMILY">family</option>
          <option :value="$options.SORT_BY.NAME">name</option>
        </select>
      </div>

      <ul class="grid grid-cols-1 gap-x-4 gap-y-1 sm:grid-cols-2 xl:grid-cols-3">
        <li
          v-for="artifact in possibleArtifacts"
          :key="artifact.sortKey"
          class="flex flex-row items-start flex-wrap"
        >
          <artifact-name class="mr-1" :artifact="artifact" :showTier="true" />
          <tippy>
            <span
              class="flex flex-row items-center text-sm"
              :class="artifact.notEnoughData ? 'text-gray-500' : null"
              >({{ artifact.itemsPerMission.toPrecision(artifact.precision)
              }}<sup v-if="artifact.notEnoughData" class="-top-0.5">?</sup>/<img
                class="h-4 w-4"
                :src="iconURL(mission.shipIconPath, 32)"
              />, {{ artifact.itemsPerDay.toPrecision(artifact.precision)
              }}<sup v-if="artifact.notEnoughData" class="-top-0.5">?</sup>/d)</span
            >

            <template #content>
              <p>Received {{ artifact.itemCount }} from {{ artifact.missionCount }} missions.</p>
              <p v-if="artifact.notEnoughData">Not enough data, rate likely far from accurate.</p>
            </template>
          </tippy>
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
import ArtifactName from "@/components/ArtifactName.vue";
import LootDataCredit from "@/components/LootDataCredit.vue";
import MissionName from "@/components/MissionName.vue";
import Share from "@/components/Share.vue";

import { getLocalStorage, setLocalStorage, iconURL, stringCmp } from "@/utils";

const SORT_BY_LOCALSTORAGE_KEY = "artifactSortBy";
const SORT_BY = {
  QUALITY: "quality",
  FAMILY: "family",
  NAME: "name",
};

export default {
  components: {
    ArtifactName,
    LootDataCredit,
    MissionName,
    Share,
  },

  props: {
    missionId: String,
    missions: Array,
    artifacts: Array,
    lootTable: Object,
  },

  data() {
    return {
      sortBy: getLocalStorage(SORT_BY_LOCALSTORAGE_KEY) || SORT_BY.QUALITY,
    };
  },

  SORT_BY,

  computed: {
    mission() {
      for (const mission of this.missions) {
        if (mission.id === this.missionId) {
          return {
            ...mission,
            loot: this.lootTable[mission.id],
          };
        }
      }
      return undefined;
    },

    possibleArtifacts() {
      if (!this.mission) {
        return [];
      }
      return this.artifacts
        .filter(
          artifact =>
            this.mission.minQuality <= artifact.quality &&
            artifact.quality <= this.mission.maxQuality
        )
        .sort((artifact1, artifact2) => {
          switch (this.sortBy) {
            case SORT_BY.FAMILY:
              return stringCmp(artifact1.sortKey, artifact2.sortKey);
            case SORT_BY.NAME:
              return stringCmp(artifact1.name, artifact2.name);
            default:
              return artifact1.quality - artifact2.quality;
          }
        })
        .map(artifact => {
          const missionCount = this.mission.loot.missionCount;
          const itemCount = this.mission.loot.items[artifact.id].rarities[artifact.afxRarity];
          const itemsPerMission = missionCount === 0 ? 0 : itemCount / missionCount;
          const itemsPerDay = (itemsPerMission * 3600 * 24) / this.mission.durationSeconds;
          const precision = Math.min(itemCount.toString().length, 3);
          const notEnoughData = itemCount < 20;
          return {
            ...artifact,
            missionCount,
            itemCount,
            itemsPerMission,
            itemsPerDay,
            precision,
            notEnoughData,
          };
        });
    },
  },

  watch: {
    sortBy() {
      setLocalStorage(SORT_BY_LOCALSTORAGE_KEY, this.sortBy);
    },
  },

  methods: {
    iconURL,
  },
};
</script>
