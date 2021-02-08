<template>
  <ul class="my-2 grid grid-cols-1 gap-4 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4">
    <template v-for="family in families" :key="family.name">
      <li class="col-span-1 bg-gray-50 rounded-lg shadow divide-y divide-gray-200">
        <div class="w-full flex items-center justify-between px-6 py-4 space-x-4">
          <div class="flex-1 truncate">
            <div class="flex items-center space-x-3">
              <h3 class="text-gray-900 text-sm font-medium truncate" :title="effect(family)">
                {{ family.name }}
              </h3>
            </div>
            <p class="mt-1 text-gray-500 text-xs truncate" :title="effect(family)">
              {{ effect(family) }}
            </p>
            <div class="mt-2 space-y-1.5">
              <template v-for="tier in family.tiers" :key="tier.tier_number">
                <router-link :to="{ name: 'artifact', params: { artifactId: tier.id } }">
                  <div class="flex">
                    <div
                      class="flex items-center space-x-2"
                      v-tippy="{
                        content: `<img src='${iconURL(tier.iconPath, 256)}' class='h-16 w-16'>`,
                        allowHTML: true,
                      }"
                    >
                      <img class="h-12 w-12" :src="iconURL(tier.iconPath, 128)" />
                      <div>
                        <div
                          class="text-xs"
                          :class="
                            tier.available_from_missions
                              ? 'text-gray-500'
                              : ['text-red-900', 'dagger']
                          "
                        >
                          {{ tier.tier_name }}
                        </div>
                        <div class="text-xs text-gray-400 tabular-nums truncate">
                          <template
                            v-for="(rarity, index) in tier.effects"
                            :key="rarity.afx_rarity"
                          >
                            <template v-if="index !== 0">, </template>
                            <span :class="rarityFgClass(rarity.afx_rarity)">{{
                              rarity.effect_size
                            }}</span>
                          </template>
                        </div>
                      </div>
                    </div>
                  </div>
                </router-link>
              </template>
            </div>
          </div>
        </div>
      </li>
    </template>
  </ul>
</template>

<script>
import { iconURL } from "@/utils";

export default {
  props: {
    families: Array,
  },

  data() {
    console.log(this.families[0]);
  },

  methods: {
    iconURL,

    effect(family) {
      // Get family effect from the last tier. The first tier's effects may be
      // null (stone fragment).
      const tier = family.tiers[family.tiers.length - 1];
      return tier.has_effects ? tier.effects[0].family_effect : "";
    },

    rarityFgClass(rarity) {
      switch (rarity) {
        case 1:
          return "text-blue-400";
        case 2:
          return "text-purple-400";
        case 3:
          return "text-yellow-400";
        default:
          return "";
      }
    },
  },
};
</script>

<style scoped>
.dagger:after {
  content: "\2020";
}
</style>
