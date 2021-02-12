<template>
  <div class="flex flex-wrap items-center text-xs text-gray-500 leading-9 tabular-nums">
    <span
      v-for="source in sources"
      :key="source.id"
      class="inline-flex items-center mr-1.5 whitespace-nowrap"
    >
      <span class="inline-block h-8 w-8 relative mr-0.5">
        <a :href="`#${source.id}`" class="h-8 w-8">
          <img
            class="h-8 w-8 rounded-md"
            :class="rarityBgClass(source.afx_rarity)"
            :src="iconURL(`egginc/${source.icon_filename}`, 64)"
            v-tippy="{
              content: `${source.name} (T${source.tier_number})`,
            }"
          />
        </a>
        <span
          v-if="!source.deterministic"
          class="absolute -top-0.5 -right-0.5 block h-2 w-2 rounded-full ring-2 ring-white bg-green-400"
        ></span>
      </span>
      <span class="sr-only">{{ source.name }}</span>
      <span class="w-9">{{ source.expected_yield }}&times;</span>
    </span>
  </div>
</template>

<script>
export default {
  props: {
    sources: Array,
  },

  methods: {
    rarityBgClass(rarity) {
      switch (rarity) {
        case 1:
          return "bg-rare";
        case 2:
          return "bg-epic";
        case 3:
          return "bg-legendary";
        default:
          return "bg-common";
      }
    },
  },
};
</script>
