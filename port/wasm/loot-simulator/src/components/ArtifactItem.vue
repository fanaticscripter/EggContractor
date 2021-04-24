<template>
  <span class="flex items-center space-x-1">
    <a :href="itemURL" target="_blank" class="inline h-6 w-6">
      <img :src="itemIconURL" class="inline h-6 w-6" />
    </a>
    <a :href="itemURL" target="_blank" class="text-sm truncate">{{ item.name }}</a>
    <span v-if="count !== undefined" class="text-sm whitespace-nowrap">
      &times; {{ count.toFixed(decimals) }}
    </span>
  </span>
</template>

<script lang="ts">
import { computed, defineComponent, PropType, toRefs } from "vue";

import { itemIdToItem } from "@/data";
import { ItemId } from "@/types";
import { iconURL } from "@/utils";

export default defineComponent({
  props: {
    id: {
      type: String as PropType<ItemId>,
      required: true,
    },
    count: {
      type: Number,
      required: false,
    },
    decimals: {
      type: Number,
      default: 0,
    },
  },
  setup(props) {
    const { id } = toRefs(props);
    const item = computed(() => itemIdToItem.get(id.value)!);
    const itemURL = computed(
      () => `https://wasmegg.netlify.app/artifact-explorer/#/artifact/${id.value}/`
    );
    const itemIconURL = computed(() => iconURL(item.value.iconPath, 64));
    return {
      item,
      itemURL,
      itemIconURL,
    };
  },
});
</script>
