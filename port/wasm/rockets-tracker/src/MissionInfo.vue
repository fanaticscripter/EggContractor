<template>
  <ul v-if="activeMissions && activeMissions.length > 0" class="grid grid-cols-1 gap-6 sm:grid-cols-2 md:grid-cols-3 mx-4 xl:mx-0 my-4" :class="[activeMissions.length >= 4 ? 'lg:grid-cols-4' : 'lg:grid-cols-3']">
    <li v-for="(mission, index) in activeMissions" :key="index" class="col-span-1 flex flex-col text-center bg-gray-50 rounded-2xl shadow-lg divide-y divide-gray-200">
      <div class="flex-1 flex flex-col p-6">
        <div class="w-36 h-36 flex-shrink-0 mx-auto relative" :class="[durationTypeFgClass(mission.durationTypeDisplay)]">
          <img class="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-32 h-32 rounded-full" :src="iconURL(mission.shipIconPath, 256)" :alt="mission.shipName">
          <progress-ring :radius="72" :stroke="2" :duration="mission.durationSeconds" :deadline="mission.returnTimestamp"></progress-ring>
        </div>
        <h3 class="mt-4 text-gray-900 text-sm font-medium">{{ mission.shipName }}</h3>
        <div class="mt-1 flex-grow flex flex-col">
          <div>
            <span class="px-2 py-1 text-white text-xs font-medium rounded-full" :class="[durationTypeBgClass(mission.durationTypeDisplay)]">{{ mission.durationTypeDisplay }}</span>
          </div>
          <div class="mt-2 text-gray-500 text-xs">Capacity: {{ mission.capacity }}</div>
          <div class="mt-1 text-gray-500 text-xs">
            Duration:
            <template v-if="mission.durationSeconds > 0">
              {{ mission.durationDisplay }}
            </template>
            <template v-else>
              &ndash;
            </template>
          </div>
          <div class="mt-1 text-gray-700 text-sm font-medium">{{ mission.statusDisplay }}</div>
          <div v-if="mission.returnTimestamp > 0" class="mt-1 text-gray-700 text-sm font-medium tabular-nums">
            <countdown-timer :deadline="mission.returnTimestamp"></countdown-timer>
          </div>
          <div v-else-if="mission.statusDisplay === 'Fueling' && mission.fuels && mission.fuels.length > 0" class="mt-1">
            <img v-for="fuel in mission.fuels" :key="fuel.egg" class="inline h-4 w-4 align-text-top" :src="iconURL(fuel.eggIconPath, 64)" v-tippy="{ content: fuel.amountDisplay + ' (this value from your save data may lag behind your actual progress)' }">
          </div>
        </div>
      </div>
    </li>
  </ul>

  <div v-else class="text-center mt-6">
    No active mission. You should start one!
  </div>

  <!-- Notifications toggle -->
  <div class="mt-2 space-y-1">
    <div
      v-if="notificationSupportedByBrowser"
      class="flex items-center justify-center space-x-2"
      v-tippy="{
        content:
          `Your browser's notifications feature is used to display system-level notifications when your rockets return. ` +
          `You have to allow notifications from this site when prompted to enable this feature.`,
        allowHTML: true,
      }"
    >
      <button id="notifications" name="notifications" type="button" class="flex-shrink-0 group relative rounded-full inline-flex items-center justify-center h-5 w-10 cursor-pointer focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500" @click="toggleNotifications()">
        <span class="sr-only">Toggle notifications</span>
        <span aria-hidden="true" class="absolute h-4 w-9 mx-auto rounded-full transition-colors ease-in-out duration-200" :class="[notificationsOn ? 'bg-green-400' : 'bg-gray-200']"></span>
        <span aria-hidden="true" class="absolute left-0 inline-block h-5 w-5 border border-gray-200 rounded-full bg-white shadow transform ring-0 transition-transform ease-in-out duration-200" :class="[notificationsOn ? 'translate-x-5' : 'translate-x-0']"></span>
      </button>
      <label for="notifications" class="text-sm text-gray-600">Mission return notifications</label>
    </div>

    <div
      v-if="notificationPermissionDenied"
      class="mt-2 flex items-center justify-center space-x-1"
      v-tippy="{
        content:
          `The site has been denied / did not receive permissions to show you notifications. ` +
          `Reload the page and try to toggle notifications again; or Google ` +
          `“how to allow website notifications in <your browser>” for instructions.`,
      }"
    >
      <svg class="h-4 w-4 text-red-500" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
        <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
      </svg>
      <a href="https://www.google.com/search?q=how+to+allow+website+notifications+in+my+browser" target="_blank" class="text-xs text-red-500 hover:text-red-400">Notifications permission denied</a>
      <svg class="h-4 w-4 text-gray-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
        <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-8-3a1 1 0 00-.867.5 1 1 0 11-1.731-1A3 3 0 0113 8a3.001 3.001 0 01-2 2.83V11a1 1 0 11-2 0v-1a1 1 0 011-1 1 1 0 100-2zm0 8a1 1 0 100-2 1 1 0 000 2z" clip-rule="evenodd" />
      </svg>
    </div>
  </div>

  <div class="mt-6">
    <h2 class="mx-4 my-4 text-center text-md leading-6 font-medium text-gray-900">Mission statistics</h2>

    <div v-if="unlockProgress" class="-mt-2 mb-2 -space-y-1">
      <div v-if="unlockProgress.nextShipToLaunch.name !== unlockProgress.nextShipToUnlock.name" class="text-sm text-center space-x-1">
        <img class="inline w-8 h-8" :src="iconURL(unlockProgress.nextShipToLaunch.iconPath, 128)" alt="">
        <span class="whitespace-nowrap">
          {{ unlockProgress.nextShipToLaunch.name }} unlocked
        </span>
      </div>

      <div v-if="unlockProgress.nextShipToUnlock" class="text-sm text-center space-x-1">
        <img class="inline w-8 h-8" :src="iconURL(unlockProgress.nextShipToUnlock.iconPath, 128)" alt="">
        <span class="whitespace-nowrap">
          {{ unlockProgress.nextShipToUnlock.name }} unlock:
          <span class="font-medium">
            {{ unlockProgress.nextShipToUnlock.launchesDone }} / {{ unlockProgress.nextShipToUnlock.launchesRequired }}
          </span>
          <span
            class="inline-flex items-center space-x-1 ml-1"
            v-tippy="{ content: 'This is the accumulative mission time for all remaining required launches, including mission time for previous locked ships (if any), divided by max concurrency (3 for Pro Permit holders, 1 otherwise). Remaining timer on your active mission(s) is not included.' }"
          >
            ({{ unlockProgress.nextShipToUnlock.accumulativeMissionTimeRequired }} <svg class="h-4 w-4 text-gray-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
              <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-8-3a1 1 0 00-.867.5 1 1 0 11-1.731-1A3 3 0 0113 8a3.001 3.001 0 01-2 2.83V11a1 1 0 11-2 0v-1a1 1 0 011-1 1 1 0 100-2zm0 8a1 1 0 100-2 1 1 0 000 2z" clip-rule="evenodd" />
            </svg>)
          </span>
        </span>
      </div>

      <div class="py-1 text-center text-xs text-gray-500 underline cursor-pointer select-none" @click="furtherShipsToUnlockExpanded = !furtherShipsToUnlockExpanded">
        Click here to {{ furtherShipsToUnlockExpanded ? "collapse": "expand" }} the full list of ships to unlock
      </div>
      <template v-if="furtherShipsToUnlockExpanded && unlockProgress.furtherShipsToUnlock">
        <div v-for="ship in unlockProgress.furtherShipsToUnlock" :key="ship.name" class="text-sm text-center space-x-1">
          <img class="inline w-8 h-8" :src="iconURL(ship.iconPath, 128)" alt="">
          <span class="whitespace-nowrap">
            {{ ship.name }} unlock:
            <span class="font-medium">
              {{ ship.launchesDone }} / {{ ship.launchesRequired }}
            </span>
            <span class="inline-flex items-center space-x-1 ml-1">
              ({{ ship.accumulativeMissionTimeRequired }})
            </span>
          </span>
        </div>
      </template>
    </div>

    <div class="flex flex-col">
      <div class="-my-2 overflow-x-auto xl:-mx-4">
        <div class="py-2 align-middle inline-block min-w-full lg:px-4">
          <div class="shadow overflow-hidden border-b border-gray-200">
            <table class="min-w-full divide-y divide-gray-200">
              <thead class="bg-gray-50">
                <th scope="col" class="px-6 py-2 text-center text-xs font-medium text-gray-500 uppercase">Ship</th>
                <th scope="col" class="px-6 py-2 text-center text-xs font-medium text-gray-500 uppercase">Type</th>
                <th scope="col" class="px-6 py-2 text-center text-xs font-medium text-gray-500 uppercase">Duration</th>
                <th scope="col" class="px-6 py-2 text-center text-xs font-medium text-gray-500 uppercase">Capacity</th>
                <th scope="col" class="px-6 py-2 text-center text-xs font-medium text-gray-500 uppercase">Fuels</th>
                <th scope="col" class="px-6 py-2 text-center text-xs font-medium text-gray-500 uppercase">Launched</th>
              </thead>
              <tbody class="divide-y divide-gray-200">
                <template v-for="ship in missionStats.ships" :key="ship.shipName">
                  <tr class="text-gray-500">
                    <td class="relative px-6 py-1.5 border-r whitespace-nowrap text-sm" :rowspan="ship.types.length + 1">
                      <img class="absolute top-1/2 left-6 transform -translate-y-1/2 w-12 h-12" :src="iconURL(ship.shipIconPath, 128)" :alt="ship.shipName">
                      <span class="pl-14">{{ ship.shipName }}</span>
                    </td>
                    <td class="px-6 py-1.5 bg-gray-50 whitespace-nowrap text-center text-sm">Aggregate</td>
                    <td class="px-6 py-1.5 bg-gray-50 whitespace-nowrap text-center text-sm"></td>
                    <td class="px-6 py-1.5 bg-gray-50 whitespace-nowrap text-center text-sm"></td>
                    <td class="px-6 py-1.5 bg-gray-50 whitespace-nowrap text-center text-sm"></td>
                    <td class="px-6 py-1.5 bg-gray-50 whitespace-nowrap text-center text-sm">{{ ship.count }}</td>
                  </tr>
                  <template v-for="type in ship.types" :key="type.durationTypeDisplay">
                    <tr class="text-gray-500">
                      <td class="px-6 py-1.5 whitespace-nowrap text-center text-sm">{{ type.durationTypeDisplay }}</td>
                      <td class="px-6 py-1.5 whitespace-nowrap text-center text-sm">{{ type.durationDisplay }}</td>
                      <td class="px-6 py-1.5 whitespace-nowrap text-center text-sm">{{ type.capacity }}</td>
                      <td class="px-6 py-1.5 whitespace-nowrap text-center text-sm">
                        <img v-for="fuel in type.fuels" :key="fuel.egg" class="inline h-4 w-4" :src="iconURL(fuel.eggIconPath, 64)" v-tippy="{ content: fuel.amountDisplay }">
                      </td>
                      <td class="px-6 py-1.5 whitespace-nowrap text-center text-sm">{{ type.count }}</td>
                    </tr>
                  </template>
                </template>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>

    <div class="mx-4 my-4 xl:mx-0 text-xs">
      Notes:
      <ul class="list-disc">
        <li>Launched active missions are counted in mission statistics; fueling active missions are not.</li>
        <li>There might be a delay between your local mission state updates and server state pulled by this site.</li>
      </ul>
    </div>
  </div>

  <div class="mx-4 xl:mx-0 my-4">
    <h2 class="mx-4 mt-4 mb-2 text-center text-md leading-6 font-medium text-gray-900">Launch log</h2>
    <div class="w-full max-w-xs mx-auto my-2">
      <label for="launch-log-filter" class="sr-only">Filter launch log</label>
      <select id="launch-log-filter" name="launch-log-filter" class="mt-1 block w-full pl-3 pr-10 py-1 text-sm bg-gray-50 border-gray-300 focus:outline-none focus:ring-green-500 focus:border-green-500 sm:text-sm rounded-md" v-model="launchLogFilter">
        <option value="3d">Limit to 3 days</option>
        <option value="7d">Limit to 7 days</option>
        <option value="30d">Limit to 30 days</option>
        <option value="none">Show all</option>
      </select>
    </div>
    <div>
      <template v-for="date in filteredLaunchLogDates" :key="date.date">
        <div class="my-2 text-sm font-medium text-gray-900">{{ date.date }}</div>
        <div class="grid grid-cols-1 gap-x-6 gap-y-2 sm:grid-cols-2 lg:grid-cols-3">
          <div v-for="(mission, index) in date.missions" :key="index" class="text-xs tabular-nums">
            <span class="mr-2">{{ formatTime(mission.startTimestamp) }}</span>
            <span class="mr-1">{{ mission.shipName }}</span>
            <span :class="[durationTypeFgClass(mission.durationTypeDisplay)]">{{
              mission.durationTypeDisplay
            }}</span>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script>
import CountdownTimer from "./CountdownTimer.vue";
import ProgressRing from "./ProgressRing.vue";
import { getLocalStorage, setLocalStorage, iconURL } from "./utils";

// Note: This component must be recreated for notifications to properly
// unregister/re-register.
export default {
  components: {
    CountdownTimer,
    ProgressRing,
  },

  props: {
    activeMissions: Array,
    missionStats: Object,
    unlockProgress: Object,
    launchLog: Object,
  },

  data() {
    return {
      notificationSupportedByBrowser: "Notification" in window,
      notificationsOn: getLocalStorage("notifications") === "true",
      notificationPermissionDenied: false,
      furtherShipsToUnlockExpanded: false,
      launchLogFilter: getLocalStorage("launchLogFilter") || "7d",
    };
  },

  async mounted() {
    if (this.notificationsOn) {
      await this.registerNotifications();
    }
  },

  beforeUnmount() {
    this.unregisterNotifications();
  },

  // setTimeout IDs for scheduled notifications.
  notificationTimeoutIds: [],

  computed: {
    filteredLaunchLogDates () {
      let days = 0;
      switch (this.launchLogFilter) {
        case "3d":
          days = 3;
          break;
        case "7d":
          days = 7;
          break;
        case "30d":
          days = 30;
          break;
      }
      if (days === 0) {
        // Do not filter.
        return this.launchLog.dates;
      }
      const now = new Date();
      const earliestDay = new Date(now.getFullYear(), now.getMonth(), now.getDate() - (days - 1));
      return this.launchLog.dates.filter(date => {
        const [yy, mm, dd] = date.date.split("-").map(parseFloat);
        return new Date(yy, mm - 1, dd) >= earliestDay;
      })
    }
  },

  watch: {
    launchLogFilter() {
      setLocalStorage("launchLogFilter", this.launchLogFilter);
    },
  },

  methods: {
    async toggleNotifications() {
      this.notificationsOn = !this.notificationsOn;
      setLocalStorage("notifications", this.notificationsOn);
      if (this.notificationsOn) {
        await this.registerNotifications();
      } else {
        this.unregisterNotifications();
      }
    },

    async registerNotifications() {
      if (!("Notification" in window)) {
        return;
      }
      // Safari doesn't support the promise version of Notification.requestPermission. Thanks Apple!
      const permission = await new Promise((resolve, _) => {
        Notification.requestPermission(result => resolve(result));
      });
      if (permission === "default" || permission === "denied") {
        console.error(`no permissions to display notifications: ${permission}`);
        this.notificationPermissionDenied = true;
        this.notificationsOn = false;
        setLocalStorage("notifications", this.notificationsOn);
        return;
      }
      this.notificationPermissionDenied = false;
      for (const mission of this.activeMissions) {
        if (mission.returnTimestamp > 0) {
          const desc = `${mission.durationTypeDisplay.toLowerCase()} ${mission.shipName}`;
          const timeout = mission.returnTimestamp * 1000 - Date.now();
          if (timeout < 10000) {
            continue;
          }
          const timeoutId = setTimeout(() => {
            new Notification(`Your ${desc} mission has returned!`, {
              icon: iconURL(mission.shipIconPath),
            });
          }, timeout);
          console.log(
            `scheduled notifaction for ${desc} in ${Math.round(
              timeout
            )}ms (timeout ID ${timeoutId})`
          );
          this.$options.notificationTimeoutIds.push(timeoutId);
        }
      }
    },

    unregisterNotifications() {
      while (this.$options.notificationTimeoutIds.length > 0) {
        const timeoutId = this.$options.notificationTimeoutIds.pop();
        clearTimeout(timeoutId);
        console.log(`unscheduled notification with timeout ID ${timeoutId}`);
      }
    },

    formatTime(timestamp) {
      return new Intl.DateTimeFormat("en-US", {
        hour: "2-digit",
        minute: "2-digit",
        second: "2-digit",
        hourCycle: "h23",
      }).format(new Date(timestamp * 1000));
    },

    durationTypeFgClass(durationType) {
      switch (durationType) {
        case "Tutorial":
        case "Short":
          return "text-blue-500";
        case "Standard":
          return "text-purple-500";
        case "Extended":
          return "text-yellow-500";
        default:
          return "text-black";
      }
    },

    durationTypeBgClass(durationType) {
      switch (durationType) {
        case "Tutorial":
        case "Short":
          return "bg-blue-500";
        case "Standard":
          return "bg-purple-500";
        case "Extended":
          return "bg-yellow-500";
        default:
          return "bg-black";
      }
    },

    iconURL,
  },
};
</script>
