<template>
  <form
    class="sm:mx-auto sm:max-w-xs sm:w-full m-4 space-y-3"
    @submit.prevent="loadContractTable"
  >
    <div>
      <label for="email" class="sr-only">Player ID</label>
      <input
        type="text"
        name="playerId"
        id="playerId"
        v-model="playerId"
        class="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
        placeholder="Player ID"
      />
    </div>
    <div>
      <button
        type="submit"
        class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
        :class="{ 'cursor-not-allowed': submitDisabled }"
        :disabled="submitDisabled"
      >
        Load Contracts
      </button>
    </div>
  </form>

  <template v-if="!playerIdSubmitted">
    <!-- Waiting for first submission of playerId -->
  </template>
  <template v-else-if="contractsLoading">
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
    <contract-table :contracts="contracts" :csvdata="csvdata"></contract-table>
  </template>
</template>

<script>
import ContractTable from "./ContractTable.vue";

function getLocalStorage(key) {
  try {
    return localStorage[`${window.location.pathname}_${key}`];
  } catch (err) {
    console.error(err);
    return undefined;
  }
}

function setLocalStorage(key, val) {
  try {
    localStorage[`${window.location.pathname}_${key}`] = val;
  } catch (err) {
    console.error(err);
  }
}

export default {
  components: {
    ContractTable,
  },
  props: {
    retrievePastContracts: Function,
  },
  data() {
    const playerId =
      new URLSearchParams(window.location.search).get("playerId") ||
      getLocalStorage("playerId") ||
      "";
    return {
      playerId,
      playerIdSubmitted: false,
      contracts: [],
      contractsLoading: false,
      error: "",
    };
  },
  computed: {
    submitDisabled() {
      return this.playerId.trim() === "" || this.contractsLoading;
    },
    truncatedError() {
      return this.error.length <= 500 ? this.error : `${this.error.substr(0, 497)}...`;
    },
  },
  methods: {
    async loadContractTable() {
      const playerId = this.playerId.trim();
      if (!playerId) {
        return;
      }
      setLocalStorage("playerId", playerId);
      this.playerIdSubmitted = true;
      this.contractsLoading = true;
      this.error = "";
      try {
        const { contracts, csv: csvdata } = await this.retrievePastContracts(playerId);
        this.contracts = contracts;
        this.csvdata = csvdata;
        this.contractsLoading = false;
        this.error = null;
      } catch (err) {
        console.error(err);
        this.contractsLoading = false;
        this.error = err;
      }
    },
  },
};
</script>
