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
        </div>
      </div>
    </li>
  </ul>

  <div v-else class="text-center mt-6">
    No active mission. You should start one!
  </div>

  <div class="mt-6">
    <h2 class="mx-4 my-4 text-center text-md leading-6 font-medium text-gray-900">Mission statistics</h2>

    <div v-if="unlockProgress && unlockProgress.nextShipToLaunchName !== unlockProgress.nextShipToUnlockName" class="-mt-4 text-sm text-center">
      <img class="inline w-12 h-12" :src="iconURL(unlockProgress.nextShipToLaunchIconPath, 128)" alt="">
      <span class="whitespace-nowrap">
        {{ unlockProgress.nextShipToLaunchName }} unlocked
      </span>
    </div>
    <div v-if="unlockProgress && unlockProgress.hasShipToUnlock" class="-mt-4 text-sm text-center">
      <img class="inline w-12 h-12" :src="iconURL(unlockProgress.nextShipToUnlockIconPath, 128)" alt="">
      <span class="whitespace-nowrap">
        {{ unlockProgress.nextShipToUnlockName }} unlock:
        <span class="font-medium">{{ unlockProgress.launchesDone }} / {{ unlockProgress.launchesRequiredToUnlock }}</span>
      </span>
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
    <div>
      <template v-for="date in launchLog.dates" :key="date.date">
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
import { iconURL } from "./utils";

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
    return {};
  },

  methods: {
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
