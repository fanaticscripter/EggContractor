<template>
  <div class="mx-4 xl:mx-0 my-4">
    <h2 class="mx-4 mt-4 mb-2 text-center text-md leading-6 font-medium text-gray-900">
      Artifacting progress
    </h2>

    <div class="flex justify-center my-2">
      <div class="relative flex items-start">
        <div class="flex items-center h-5">
          <input
            id="spoilers"
            name="spoilers"
            v-model="spoilers"
            type="checkbox"
            class="focus:ring-green-500 h-4 w-4 text-green-600 border-gray-300 rounded"
          />
        </div>
        <div class="ml-2 text-sm">
          <label for="spoilers" class="text-gray-600">Show unseen items (SPOILERS)</label>
        </div>
      </div>
    </div>

    <div class="my-2 space-y-1">
      <div
        class="flex flex-wrap items-center justify-center space-x-1"
        v-tippy="{
          content:
            'This estimate does not take into discounts from crafting sales (as the save file simply does not have that level of granularity). It may also be inaccurate if crafting cost parameters were ever changed server-side.<br>Note that stone-setting expenses are not included.',
          allowHTML: true,
        }"
      >
        <span class="text-sm text-gray-900 truncate"
          >Estimated total crafting expense so far (pre-discount):</span
        >
        <span class="inline-flex items-center space-x-1 text-sm text-gray-900">
          <img class="h-4 w-4" :src="iconURL('egginc-extras/icon_golden_egg.png', 64)" />
          <span>{{ totalCraftingCost.toLocaleString("en-US") }}</span>
          <info />
        </span>
      </div>

      <div
        v-if="!isNaN(accountBalance)"
        class="flex flex-wrap items-center justify-center space-x-1"
      >
        <span class="text-sm text-gray-900 truncate">Account balance:</span>
        <span class="inline-flex items-center space-x-1 text-sm text-gray-900">
          <img class="h-4 w-4" :src="iconURL('egginc-extras/icon_golden_egg.png', 64)" />
          {{ accountBalance.toLocaleString("en-US") }}
        </span>
      </div>

      <div
        v-if="inventoryScore !== undefined"
        class="flex flex-wrap items-center justify-center space-x-1"
        v-tippy="{
          content:
            'The inventory score is an internal value used to determine the outer appearance of the hall of artifacts. The hall expands from two to four segments once the score reaches 1000, and the segments become colored once the score reaches 5000. The effect of this score is entirely cosmetic.',
        }"
      >
        <span class="text-sm text-gray-900 truncate">Inventory score:</span>
        <span class="inline-flex items-center space-x-1 text-sm text-gray-900">
          <span>{{ Math.floor(inventoryScore) }}</span>
          <info />
        </span>
      </div>
    </div>

    <h3 class="my-2 text-sm font-medium text-gray-900">Artifacts</h3>
    <artifact-grid :items="progress.artifacts" :spoilers="spoilers"></artifact-grid>

    <h3 class="my-2 text-sm font-medium text-gray-900">Stones &amp; stone fragments</h3>
    <artifact-grid :items="progress.stones" :spoilers="spoilers"></artifact-grid>

    <h3 class="my-2 text-sm font-medium text-gray-900">Ingredients</h3>
    <artifact-grid :items="progress.ingredients" :spoilers="spoilers"></artifact-grid>
  </div>
</template>

<script>
import ArtifactGrid from "./ArtifactGrid.vue";
import Info from "../../artifact-explorer/src/components/Info.vue";

import { getLocalStorage, setLocalStorage, iconURL } from "./utils";

export default {
  components: {
    ArtifactGrid,
    Info,
  },

  props: {
    progress: Object,
    save: Object,
  },

  computed: {
    accountBalance() {
      try {
        return (
          this.save.progress.lifetime_golden_eggs - this.save.progress.lifetime_golden_eggs_spent
        );
      } catch (e) {
        console.error(e);
        return NaN;
      }
    },

    totalCraftingCost() {
      const classes = [].concat(
        this.progress.artifacts,
        this.progress.stones,
        this.progress.ingredients
      );
      let sum = 0;
      for (const cls of classes) {
        for (const tier of cls.tiers) {
          sum += tier.craftingCost;
        }
      }
      return sum;
    },

    inventoryScore() {
      try {
        return this.save.artifacts.inventory_score;
      } catch (e) {
        console.error(e);
        return undefined;
      }
    },
  },

  data() {
    return {
      spoilers: getLocalStorage("spoilers") === "true",
    };
  },

  watch: {
    spoilers() {
      setLocalStorage("spoilers", this.spoilers);
    },
  },

  methods: {
    iconURL,
  },
};
</script>
