<template>
  <div
    v-if="!hideRarity"
    class="flex items-center space-x-1 text-xs mb-0.5"
    :class="rarityFgClass(outcome.item.afx_rarity)"
  >
    <span>{{ outcome.item.rarity }}</span>
    <span v-if="!outcome.deterministic">
      <svg class="inline h-3" viewBox="0 0 640 512">
        <path
          fill="currentColor"
          d="M592 192H473.26c12.69 29.59 7.12 65.2-17 89.32L320 417.58V464c0 26.51 21.49 48 48 48h224c26.51 0 48-21.49 48-48V240c0-26.51-21.49-48-48-48zM480 376c-13.25 0-24-10.75-24-24 0-13.26 10.75-24 24-24s24 10.74 24 24c0 13.25-10.75 24-24 24zm-46.37-186.7L258.7 14.37c-19.16-19.16-50.23-19.16-69.39 0L14.37 189.3c-19.16 19.16-19.16 50.23 0 69.39L189.3 433.63c19.16 19.16 50.23 19.16 69.39 0L433.63 258.7c19.16-19.17 19.16-50.24 0-69.4zM96 248c-13.25 0-24-10.75-24-24 0-13.26 10.75-24 24-24s24 10.74 24 24c0 13.25-10.75 24-24 24zm128 128c-13.25 0-24-10.75-24-24 0-13.26 10.75-24 24-24s24 10.74 24 24c0 13.25-10.75 24-24 24zm0-128c-13.25 0-24-10.75-24-24 0-13.26 10.75-24 24-24s24 10.74 24 24c0 13.25-10.75 24-24 24zm0-128c-13.25 0-24-10.75-24-24 0-13.26 10.75-24 24-24s24 10.74 24 24c0 13.25-10.75 24-24 24zm128 128c-13.25 0-24-10.75-24-24 0-13.26 10.75-24 24-24s24 10.74 24 24c0 13.25-10.75 24-24 24z"
        />
      </svg>
    </span>
  </div>

  <div v-if="outcome.gold > 0" class="flex">
    <span class="inline-flex items-center space-x-1 text-xs text-gray-500">
      <img class="h-4 w-4 -ml-0.5" :src="iconURL('egginc-extras/icon_golden_egg.png', 64)" />
      {{ Math.floor(outcome.gold).toLocaleString("en-US") }}
    </span>
  </div>

  <div v-else>
    <div class="flex flex-wrap text-xs text-gray-500 leading-7 tabular-nums">
      <span
        v-for="(byproduct, index) in outcome.expected_byproducts"
        :key="byproduct.id"
        class="inline-flex items-center mr-1.5 whitespace-nowrap"
      >
        <a :href="`#${byproduct.id}-sources`" class="h-6 w-6 -m-1 mr-0">
          <img
            class="h-6 w-6"
            :src="iconURL(`egginc/${byproduct.icon_filename}`, 64)"
            v-tippy="{ content: `${byproduct.name} (T${byproduct.tier_number})` }"
          />
        </a>
        <!-- Facilitate copying as text -->
        <span class="sr-only">{{ byproduct.name }}</span>
        <span class="w-9">&times;{{ byproduct.expected_count }}</span>

        <!-- Expand/collapse button -->
        <span
          v-if="index === outcome.expected_byproducts.length - 1 && !outcome.deterministic"
          @click="showSamples = !showSamples"
          class="h-4 select-none cursor-pointer text-green-500 hide-in-screenshot-mode"
          v-tippy="{ content: 'Expand/collapse sample runs' }"
        >
          <svg v-if="showSamples" viewBox="0 0 320 512" class="h-4 ml-1.5">
            <path
              fill="currentColor"
              d="M31.3 192h257.3c17.8 0 26.7 21.5 14.1 34.1L174.1 354.8c-7.8 7.8-20.5 7.8-28.3 0L17.2 226.1C4.6 213.5 13.5 192 31.3 192z"
            />
          </svg>
          <svg v-else viewBox="0 0 192 512" class="h-4 ml-1.5">
            <path
              fill="currentColor"
              d="M192 127.338v257.324c0 17.818-21.543 26.741-34.142 14.142L29.196 270.142c-7.81-7.81-7.81-20.474 0-28.284l128.662-128.662c12.599-12.6 34.142-3.676 34.142 14.142z"
            />
          </svg>
        </span>
      </span>
    </div>

    <ul
      v-if="outcome.sample_byproducts !== null && showSamples"
      class="p-2 bg-gray-200 rounded-lg divide-y divide-gray-300"
    >
      <li
        v-for="(sample, index) in outcome.sample_byproducts"
        :key="index"
        class="flex flex-wrap text-xs text-gray-500 leading-7 tabular-nums"
      >
        <span
          v-for="byproduct in sample"
          :key="byproduct.id"
          class="inline-flex items-center mr-1.5"
        >
          <a :href="`#${byproduct.id}-sources`" class="h-6 w-6 -m-1 mr-0">
            <img
              class="h-6 w-6"
              :src="iconURL(`egginc/${byproduct.icon_filename}`, 64)"
              v-tippy="{
                content: `${byproduct.name} (T${byproduct.tier_number})`,
              }"
            />
          </a>
          <!-- Facilitate copying as text -->
          <span class="sr-only">{{ byproduct.name }}</span>
          <span class="w-9">&times;{{ byproduct.count }}</span>
        </span>
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  props: {
    outcome: Object,
    hideRarity: Boolean,
  },

  data() {
    return {
      showSamples: false,
    };
  },

  methods: {
    rarityFgClass(rarity) {
      switch (rarity) {
        case 1:
          return "text-blue-500";
        case 2:
          return "text-purple-500";
        case 3:
          return "text-yellow-500";
        default:
          return "text-gray-500";
      }
    },
  },
};
</script>
