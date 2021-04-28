<template>
  <component
    :is="noLink ? 'div' : 'router-link'"
    :to="noLink ? undefined : { name: 'mission', params: { missionId: mission.id } }"
  >
    <div class="flex items-center">
      <div
        class="flex-shrink-0 h-5 w-5"
        v-tippy="{
          content: `<img data-src='${iconURL(
            mission.shipIconPath,
            256
          )}' class='h-36 w-36 p-1.5 rounded-full border-2 ${borderClass}'>`,
          allowHTML: true,
        }"
      >
        <img
          class="h-5 w-5 p-px rounded-full border"
          :class="borderClass"
          :src="iconURL(mission.shipIconPath, 32)"
        />
      </div>
      <span class="ml-1 text-sm">{{ mission.display }}</span>
    </div>
  </component>
</template>

<script>
import { iconURL } from "@/utils";

export default {
  props: {
    mission: Object,
    noLink: Boolean,
  },

  computed: {
    borderClass() {
      switch (this.mission.typeId) {
        case 0:
          return "border-blue-500";
        case 1:
          return "border-purple-500";
        case 2:
          return "border-yellow-500";
        default:
          return "";
      }
    },
  },

  methods: {
    iconURL,
  },
};
</script>
