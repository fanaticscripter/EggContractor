<template>
  <form class="sm:mx-auto sm:max-w-xs sm:w-full m-4 space-y-1" @submit.prevent="loadMissions">
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
      <div class="text-center">
        <span
          class="mt-2 inline-flex items-center space-x-1"
          v-tippy="{
            content:
              'The ID asked for here is the unique ID used by Egg, Inc.\'s server to identify your account. You can find it in game screen -> nine dots menu -> Settings -> Privacy & Data, at the very bottom. It should look like EI1234567890123456. Your old game services ID prior to the Artifact Update does not work here. Also note that the ID is case-sensitive.',
          }"
        >
          <svg
            class="h-4 w-4 text-gray-400"
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 20 20"
            fill="currentColor"
            aria-hidden="true"
          >
            <path
              fill-rule="evenodd"
              d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-8-3a1 1 0 00-.867.5 1 1 0 11-1.731-1A3 3 0 0113 8a3.001 3.001 0 01-2 2.83V11a1 1 0 11-2 0v-1a1 1 0 011-1 1 1 0 100-2zm0 8a1 1 0 100-2 1 1 0 000 2z"
              clip-rule="evenodd"
            />
          </svg>
          <span class="text-xs text-gray-500">Where do I find my ID?</span>
        </span>
      </div>
    </div>
    <div>
      <button
        type="submit"
        class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
        :class="{ 'cursor-not-allowed': submitDisabled }"
        :disabled="submitDisabled"
      >
        Load Player Data
      </button>
    </div>
  </form>

  <template v-if="!playerIdSubmitted">
    <!-- Waiting for first submission of playerId -->
    <!-- What's new -->
    <div v-if="activeWhatsNew.length > 0" class="rounded-md bg-green-50 mx-4 my-4 xl:mx-0 p-4">
      <h3 class="flex items-center space-x-2">
        <svg
          class="h-3 w-3 text-green-700"
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 -32 576 576"
          fill="currentColor"
          aria-hidden="true"
        >
          <path
            fill="currentColor"
            d="M32 176c-17.67 0-32 14.33-32 32v96c0 17.67 14.33 32 32 32 11.38 0 20.9-6.28 26.57-15.22l106.99 32.3c-3.35 9.76-5.56 20.04-5.56 30.92 0 52.94 43.06 96 96 96 44.49 0 81.66-30.57 92.5-71.7L480 448V64L58.57 191.22C52.9 182.28 43.38 176 32 176zm179.29 190.88l91.47 27.61C297.95 415.92 278.85 432 256 432c-26.47 0-48-21.53-48-48 0-6.05 1.24-11.79 3.29-17.12zM560 32h-32c-8.84 0-16 7.16-16 16v416c0 8.84 7.16 16 16 16h32c8.84 0 16-7.16 16-16V48c0-8.84-7.16-16-16-16z"
            class=""
          ></path>
        </svg>
        <span class="text-sm font-medium text-green-800">What's new</span>
      </h3>
      <div class="mt-2 text-sm text-green-700">
        <ul class="list-disc pl-5 space-y-1">
          <li v-for="wn in activeWhatsNew" :key="wn.id">
            <span v-if="wn.rawHTML" v-html="wn.message"></span>
            <template v-else>{{ wn.message }}</template>
          </li>
        </ul>
      </div>
    </div>
  </template>
  <template v-else-if="loading">
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
    <mission-info
      :activeMissions="activeMissions"
      :missionStats="missionStats"
      :unlockProgress="unlockProgress"
      :launchLog="launchLog"
    ></mission-info>

    <artifact-info :progress="artifactsProgress" :save="save"></artifact-info>
  </template>
</template>

<script>
import ArtifactInfo from "./ArtifactInfo.vue";
import MissionInfo from "./MissionInfo.vue";
import { getLocalStorage, setLocalStorage } from "./utils";

const whatsNew = [
  {
    id: "notifications",
    message: `<strong class="font-medium underline">OS-level mission return notifications</strong>
    (displayed in macOS Notification Center, Windows 10 Action Center, etc.) are now available!
    This feature is supported on all modern, mainstream desktop browsers, and possibly a select few
    Android browsers. You should see a "Mission return notifications" toggle right under your
    active missions if your browser is supported.
    <strong class="font-medium underline">Toggle it on and allow notifications for this site
    when prompted by your browser</strong> to enable.`,
    rawHTML: true,
    expires: 1612169010000, // Wed Jan 27 08:43:30 UTC 2021 +5days
  },
  {
    id: "artifacting-progress",
    message: `Now you can also <strong class="font-medium underline">track your artifact collection progress</strong>
    with an organized and intuitive interface. Scroll down to the "Artifacting progress" section once data is loaded.`,
    rawHTML: true,
    expires: 1612062059000, // Tue Jan 26 03:00:59 UTC 2021 + 5days
  },
];

export default {
  components: {
    MissionInfo,
    ArtifactInfo,
  },
  props: {
    retrieveMissions: Function,
  },
  data() {
    const playerId =
      new URLSearchParams(window.location.search).get("playerId") ||
      getLocalStorage("playerId") ||
      "";
    return {
      playerId,
      playerIdSubmitted: false,
      activeMissions: [],
      missionStats: null,
      unlockProgress: null,
      launchLog: null,
      artifactsProgress: null,
      save: null,
      loading: false,
      error: "",
      whatsNew,
    };
  },
  computed: {
    submitDisabled() {
      return this.playerId.trim() === "" || this.loading;
    },
    truncatedError() {
      return this.error.length <= 500 ? this.error : `${this.error.substr(0, 497)}...`;
    },
    activeWhatsNew() {
      const now = Date.now();
      return this.whatsNew.filter(wn => now < wn.expires);
    },
  },
  methods: {
    async loadMissions() {
      const playerId = this.playerId.trim();
      if (!playerId) {
        return;
      }
      setLocalStorage("playerId", playerId);
      this.playerIdSubmitted = true;
      this.loading = true;
      this.error = "";
      try {
        const {
          activeMissions,
          missionStats,
          unlockProgress,
          launchLog,
          artifactsProgress,
          save,
        } = await this.retrieveMissions(playerId);
        this.activeMissions = activeMissions;
        this.missionStats = missionStats;
        this.unlockProgress = unlockProgress;
        this.launchLog = launchLog;
        this.artifactsProgress = artifactsProgress;
        this.save = save;
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
