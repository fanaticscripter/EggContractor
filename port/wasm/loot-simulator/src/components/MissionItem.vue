<template>
  <span class="flex items-center space-x-1">
    <a :href="missionURL" target="_blank" class="inline h-6 w-6">
      <img :src="missionIconURL" class="inline h-6 w-6" />
    </a>
    <a :href="missionURL" target="_blank" class="text-sm truncate">{{ mission.display }}</a>
    <span v-if="count !== undefined" class="text-sm whitespace-nowrap"> &times; {{ count }}</span>
  </span>
</template>

<script lang="ts">
import { computed, defineComponent, PropType, toRefs } from "vue";

import { missionIdToMission } from "@/data";
import { MissionId } from "@/types";
import { iconURL } from "@/utils";

export default defineComponent({
  props: {
    id: {
      type: String as PropType<MissionId>,
      required: true,
    },
    count: {
      type: Number,
      required: false,
    },
  },
  setup(props) {
    const { id } = toRefs(props);
    const mission = computed(() => missionIdToMission.get(id.value)!);
    const missionURL = computed(
      () => `https://wasmegg.netlify.app/artifact-explorer/#/mission/${id.value}/`
    );
    const missionIconURL = computed(() => iconURL(mission.value.iconPath, 64));
    return {
      mission,
      missionURL,
      missionIconURL,
    };
  },
});
</script>
