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
                      <img class="h-8 w-8" :src="iconURL(tier.iconPath, 128)" alt="" />
                      <div class="text-xs text-gray-500">
                        {{ tier.tier_name }}
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

  methods: {
    iconURL,

    effect(family) {
      const tier = family.tiers[0];
      return tier.has_effects ? tier.effects[0].family_effect : "";
    },
  },
};
</script>
