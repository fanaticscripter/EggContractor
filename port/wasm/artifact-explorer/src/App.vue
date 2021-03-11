<template>
  <div class="flex-1 max-w-7xl w-full mx-auto px-4 xl:px-0 my-4">
    <router-view :missions="missions" :artifacts="artifacts" :lootTable="lootTable"></router-view>
  </div>

  <div class="my-4">
    <h2 class="mb-1 text-center text-base leading-6 font-medium text-gray-900">
      <code>/ei_afx/config</code> data
    </h2>

    <div class="flex items-center justify-center space-x-1">
      <svg viewBox="0 0 20 20" fill="currentColor" class="h-4 w-4 text-gray-500">
        <path
          fill-rule="evenodd"
          d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-11a1 1 0 10-2 0v3.586L7.707 9.293a1 1 0 00-1.414 1.414l3 3a1 1 0 001.414 0l3-3a1 1 0 00-1.414-1.414L11 10.586V7z"
          clip-rule="evenodd"
        />
      </svg>
      <a class="text-sm text-gray-500" :href="missionsCSVBlobURL" download="mission-parameters.csv">
        Export mission parameters as CSV
      </a>
    </div>
    <div class="flex items-center justify-center space-x-1">
      <svg viewBox="0 0 20 20" fill="currentColor" class="h-4 w-4 text-gray-500">
        <path
          fill-rule="evenodd"
          d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-11a1 1 0 10-2 0v3.586L7.707 9.293a1 1 0 00-1.414 1.414l3 3a1 1 0 001.414 0l3-3a1 1 0 00-1.414-1.414L11 10.586V7z"
          clip-rule="evenodd"
        />
      </svg>
      <a
        class="text-sm text-gray-500"
        :href="artifactsCSVBlobURL"
        download="artifact-parameters.csv"
      >
        Export artifact parameters as CSV
      </a>
    </div>
  </div>

  <eiafx-config-table :ships="ships" :missions="missions" :artifacts="artifacts" />
</template>

<script>
import EiafxConfigTable from "@/components/EiafxConfigTable.vue";

import data from "@/app-data.json";

export default {
  components: {
    EiafxConfigTable,
  },

  data() {
    const { ships, missions, artifacts, lootTable, missionsCSV, artifactsCSV } = data;
    return Object.freeze({
      ships,
      missions,
      artifacts,
      lootTable,
      missionsCSVBlobURL: window.URL.createObjectURL(new Blob([missionsCSV], { type: "text/csv" })),
      artifactsCSVBlobURL: window.URL.createObjectURL(
        new Blob([artifactsCSV], { type: "text/csv" })
      ),
    });
  },
};
</script>
