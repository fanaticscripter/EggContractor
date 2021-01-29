<template>
  <router-link :to="{ name: 'mission', params: { missionId: mission.id } }">
    <div class="flex items-center">
      <div
        class="flex-shrink-0 h-4 w-4"
        v-tippy="{
          content: `<img src='${iconURL(mission.shipIconPath, 256)}' class='h-16 w-16'>`,
          allowHTML: true,
        }"
      >
        <img class="h-4 w-4" :src="iconURL(mission.shipIconPath, 32)" />
      </div>
      <span class="ml-1 text-sm">{{ mission.display }}</span>
    </div>
  </router-link>
</template>

<script>
import { iconURL } from "@/utils";

export default {
  props: {
    mission: Object,
  },

  methods: {
    iconURL,
  },
};
</script>
