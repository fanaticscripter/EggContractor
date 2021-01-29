<template>
  <template v-if="loading">
    <div class="flex items-center justify-center">
      <svg
        class="animate-spin -ml-1 mr-3 h-5 w-5"
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
      >
        <circle
          class="opacity-25"
          cx="12"
          cy="12"
          r="10"
          stroke="currentColor"
          stroke-width="4"
        ></circle>
        <path
          class="opacity-75"
          fill="currentColor"
          d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
        ></path>
      </svg>
      Loading data...
    </div>
  </template>
  <template v-else-if="error">
    <div class="m-4 flex items-start justify-center">
      <svg
        class="h-5 w-5 flex-shrink-0 mr-1 text-red-500"
        x-description="Heroicon name: exclamation-circle"
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 20 20"
        fill="currentColor"
        aria-hidden="true"
      >
        <path
          fill-rule="evenodd"
          d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z"
          clip-rule="evenodd"
        ></path>
      </svg>
      <div class="text-sm text-red-500 leading-4 break-all">
        {{ truncatedError }}
      </div>
    </div>
  </template>
  <template v-else>
    <div class="flex-1 max-w-7xl w-full mx-auto px-4 xl:px-0 my-4">
      <router-view :missions="missions" :artifacts="artifacts"></router-view>
    </div>


    <div class="my-4">
      <h2 class="mb-1 text-center text-base leading-6 font-medium text-gray-900"><code>/ei_afx/config</code> data</h2>

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
          :href="missionsCSVBlobURL"
          download="mission-parameters.csv"
        >
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
</template>

<script>
import EiafxConfigTable from "@/components/EiafxConfigTable.vue";

export default {
  components: {
    EiafxConfigTable,
  },

  props: {
    retrieveData: Function,
  },

  data() {
    return {
      missions: null,
      artifacts: null,
      missionsCSVBlobURL: "",
      artifactsCSVBlobURL: "",
      loading: false,
      error: "",
    };
  },

  computed: {
    truncatedError() {
      return this.error.length <= 500 ? this.error : `${this.error.substr(0, 497)}...`;
    },
  },

  async mounted() {
    await this.loadData();
  },

  methods: {
    async loadData() {
      this.loading = true;
      this.error = "";
      try {
        const { ships, missions, artifacts, missionsCSV, artifactsCSV } = await this.retrieveData();
        this.ships = ships;
        this.missions = missions;
        this.artifacts = artifacts;
        this.missionsCSVBlobURL = window.URL.createObjectURL(
          new Blob([missionsCSV], { type: "text/csv" })
        );
        this.artifactsCSVBlobURL = window.URL.createObjectURL(
          new Blob([artifactsCSV], { type: "text/csv" })
        );
        this.loading = false;
        this.error = null;
      } catch (err) {
        console.error(err);
        this.loading = false;
        this.error = err;
      }
    },
  },
};
</script>
