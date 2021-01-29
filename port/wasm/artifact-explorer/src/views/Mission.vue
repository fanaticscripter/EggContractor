<template>
  <div class="-mx-4 sm:mx-0 bg-gray-50 overflow-hidden sm:rounded-lg sm:shadow-md">
    <div class="bg-gray-100 px-4 py-4 border-b border-gray-200 sm:px-6">
      <mission-name :mission="mission" />
    </div>

    <div class="px-4 py-4 sm:px-6 space-y-2">
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
        <li v-for="artifact in possibleArtifacts" :key="artifact.sortKey">
          <artifact-name :artifact="artifact" :showTier="true" />
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
import ArtifactName from "@/components/ArtifactName.vue";
import MissionName from "@/components/MissionName.vue";

import { getLocalStorage, setLocalStorage, stringCmp } from "@/utils";

const SORT_BY_LOCALSTORAGE_KEY = "artifactSortBy";
const SORT_BY = {
  QUALITY: "quality",
  FAMILY: "family",
  NAME: "name",
};

export default {
  components: {
    ArtifactName,
    MissionName,
  },

  props: {
    missionId: String,
    missions: Array,
    artifacts: Array,
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
          return mission;
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
        });
    },
  },

  watch: {
    sortBy() {
      setLocalStorage(SORT_BY_LOCALSTORAGE_KEY, this.sortBy);
    },
  },
};
</script>
