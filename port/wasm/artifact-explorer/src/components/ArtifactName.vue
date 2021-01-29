<template>
  <template v-if="plainText">
    {{ artifact.name }}
    <template v-if="showTier"> (T{{ artifact.tierNumber }}) </template>
    <template v-if="artifact.rarity > 0"> , {{ artifact.rarityDisplay }} </template>
  </template>

  <template v-else>
    <router-link :to="{ name: 'artifact', params: { artifactId: artifact.itemId } }">
      <div class="flex">
        <div
          class="flex items-center"
          v-tippy="{
            content: `<img src='${iconURL(artifact.iconPath, 256)}' class='h-16 w-16'>`,
            allowHTML: true,
          }"
        >
          <div class="flex-shrink-0 h-4 w-4">
            <img class="h-4 w-4" :src="iconURL(artifact.iconPath, 32)" />
          </div>
          <div class="ml-1 text-sm" :class="[rarityFgClass(artifact.rarity)]">
            <span>{{ artifact.name }}</span>
            <template v-if="showTier">
              <span> (T{{ artifact.tierNumber }})</span>
            </template>
            <template v-if="artifact.rarity > 0">
              <span>, {{ artifact.rarityDisplay }}</span>
            </template>
          </div>
        </div>
      </div>
    </router-link>
  </template>
</template>

<script>
import { iconURL } from "@/utils";

export default {
  props: {
    artifact: Object,
    showTier: Boolean,
    plainText: Boolean,
  },

  methods: {
    iconURL,

    rarityFgClass(rarity) {
      switch (rarity) {
        case 1:
          return "text-blue-500";
        case 2:
          return "text-purple-500";
        case 3:
          return "text-yellow-500";
        default:
          return "";
      }
    },
  },
};
</script>
