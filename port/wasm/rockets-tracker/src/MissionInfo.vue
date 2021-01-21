<template>
  <ul v-if="activeMissions && activeMissions.length > 0" class="grid grid-cols-1 gap-8 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-3 mx-4 xl:mx-0 my-4">
    <li v-for="(mission, index) in activeMissions" :key="index" class="col-span-1 flex flex-col text-center bg-gray-50 rounded-2xl shadow-lg divide-y divide-gray-200">
      <div class="flex-1 flex flex-col p-6">
        <div class="w-36 h-36 flex-shrink-0 mx-auto relative" :class="[durationTypeFgClass(mission.durationTypeDisplay)]">
          <img class="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-32 h-32 rounded-full" :src="mission.shipIconPath" :alt="mission.shipName">
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

    <div v-if="unlockProgress" class="-mt-4 text-sm text-center">
      <img class="inline w-12 h-12" :src="unlockProgress.nextShipIconPath" alt="">
      <span class="whitespace-nowrap">
        {{ unlockProgress.nextShipName }} unlock:
        <span class="font-medium">{{ unlockProgress.launchesDone }} / {{ unlockProgress.launchesRequired }}</span>
      </span>
    </div>

    <div class="flex flex-col">
      <div class="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
        <div class="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8">
          <div class="shadow overflow-hidden border-b border-gray-200">
            <table class="min-w-full divide-y divide-gray-200">
              <thead class="bg-gray-50">
                <th scope="col" class="px-6 py-2 text-center text-xs font-medium text-gray-500 uppercase">Ship</th>
                <th scope="col" class="px-6 py-2 text-center text-xs font-medium text-gray-500 uppercase">Type</th>
                <th scope="col" class="px-6 py-2 text-center text-xs font-medium text-gray-500 uppercase">Duration</th>
                <th scope="col" class="px-6 py-2 text-center text-xs font-medium text-gray-500 uppercase">Capacity</th>
                <th scope="col" class="px-6 py-2 text-center text-xs font-medium text-gray-500 uppercase">Launched</th>
              </thead>
              <tbody class="divide-y divide-gray-200">
                <template v-for="ship in missionStats.ships" :key="ship.shipName">
                  <tr class="text-gray-500">
                    <td class="relative px-6 py-1.5 border-r whitespace-nowrap text-sm" :rowspan="ship.types.length + 1">
                      <img class="absolute top-1/2 left-6 transform -translate-y-1/2 w-12 h-12" :src="ship.shipIconPath" :alt="ship.shipName">
                      <span class="pl-14">{{ ship.shipName }}</span>
                    </td>
                    <td class="px-6 py-1.5 bg-gray-50 whitespace-nowrap text-center text-sm">Aggregate</td>
                    <td class="px-6 py-1.5 bg-gray-50 whitespace-nowrap text-center text-sm"></td>
                    <td class="px-6 py-1.5 bg-gray-50 whitespace-nowrap text-center text-sm"></td>
                    <td class="px-6 py-1.5 bg-gray-50 whitespace-nowrap text-center text-sm">{{ ship.count }}</td>
                  </tr>
                  <template v-for="type in ship.types" :key="type.durationTypeDisplay">
                    <tr class="text-gray-500">
                      <td class="px-6 py-1.5 whitespace-nowrap text-center text-sm">{{ type.durationTypeDisplay }}</td>
                      <td class="px-6 py-1.5 whitespace-nowrap text-center text-sm">{{ type.durationDisplay }}</td>
                      <td class="px-6 py-1.5 whitespace-nowrap text-center text-sm">{{ type.capacity }}</td>
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
</template>

<script>
import CountdownTimer from "./CountdownTimer.vue";
import ProgressRing from './ProgressRing.vue';

export default {
  components: {
    CountdownTimer,
    ProgressRing,
  },
  props: {
    activeMissions: Array,
    missionStats: Object,
    unlockProgress: Object,
  },

  data() {
    return {};
  },

  methods: {
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
  },
};
</script>
