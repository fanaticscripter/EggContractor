<template>
  <ul class="my-4 grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
    <template v-for="cls in items" :key="cls.name">
      <li
        v-if="spoilers || cls.unlocked"
        class="col-span-1 bg-gray-50 rounded-lg shadow divide-y divide-gray-200"
      >
        <div class="w-full flex items-center justify-between px-6 py-4 space-x-4">
          <div class="flex-1 truncate">
            <div class="flex items-center space-x-3" :class="{ 'opacity-30': !cls.unlocked }">
              <h3 class="text-gray-900 text-sm font-medium truncate" :title="cls.effect">
                {{ cls.name }}
              </h3>
            </div>
            <p
              class="mt-1 text-gray-500 text-xs truncate"
              :class="{ 'opacity-30': !cls.unlocked }"
              :title="cls.effect"
            >
              {{ cls.effect }}
            </p>
            <div class="mt-2 space-y-1.5">
              <template v-for="tier in cls.tiers" :key="tier.tierNumber">
                <div
                  v-if="spoilers || tier.previousTierUnlocked"
                  class="flex items-center space-x-2"
                  :class="{ 'opacity-30': !tier.unlocked }"
                >
                  <a
                    v-if="spoilers || tier.unlocked"
                    :href="iconURL(tier.iconPath)"
                    target="_blank"
                  >
                    <img class="h-12 w-12" :src="iconURL(tier.iconPath, 128)" />
                  </a>
                  <img
                    v-else
                    class="h-12 w-12 silhouette"
                    :src="iconURL(tier.iconPath, 128)"
                    v-tippy="{ content: 'turn on &quot;show unseen items&quot; to unlock' }"
                  />

                  <div v-if="spoilers || tier.unlocked">
                    <div class="text-xs text-gray-500">
                      {{ tier.name }} &times;
                      <span
                        v-tippy="{
                          content:
                            'total number of this item currently owned, including &quot;shiny&quot; ones',
                        }"
                        >{{ tier.count }}</span
                      >
                    </div>
                    <div class="text-xs text-gray-400">{{ rarityReport(tier) }}</div>
                    <div
                      v-if="tier.craftedCount > 0"
                      class="text-xs text-gray-400"
                      v-tippy="{
                        content:
                          `You spent an estimated ${tier.craftingCost.toLocaleString('en-US')} ` +
                          `golden eggs on crafting this item!`,
                      }"
                    >
                      Crafted {{ tier.craftedCount }}
                    </div>
                  </div>

                  <div
                    v-else
                    class="text-xs text-gray-500"
                    v-tippy="{ content: 'turn on &quot;show unseen items&quot; to unlock' }"
                  >
                    ?
                  </div>
                </div>
              </template>
            </div>
          </div>
        </div>
      </li>
    </template>
  </ul>
</template>

<script>
import { iconURL } from "./utils";

export default {
  props: {
    items: Array,
    spoilers: Boolean,
  },

  methods: {
    rarityReport(tier) {
      let clauses = [];
      if (tier.rarityCounts[1] > 0) {
        clauses.push(`${tier.rarityCounts[1]} Rare`);
      }
      if (tier.rarityCounts[2] > 0) {
        clauses.push(`${tier.rarityCounts[2]} Epic`);
      }
      if (tier.rarityCounts[3] > 0) {
        clauses.push(`${tier.rarityCounts[3]} Legendary`);
      }
      return clauses.join(", ");
    },

    iconURL,
  },
};
</script>

<style scoped>
img.silhouette {
  filter: contrast(0%) brightness(50%);
}
</style>
