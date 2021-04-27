<template>
  <template v-if="plainText">
    {{ artifact.name }}
    <template v-if="showTier"> (T{{ artifact.tier_number }}) </template>
    <template v-if="artifact.afxRarity > 0"> , {{ artifact.rarity }} </template>
  </template>

  <template v-else>
    <component
      :is="noLink ? 'div' : 'router-link'"
      :to="noLink ? undefined : { name: 'artifact', params: { artifactId: artifact.id } }"
    >
      <div class="flex">
        <div
          class="flex items-center"
          v-tippy="{
            content: `<img src='${iconURL(artifact.iconPath, 256)}' class='h-32 w-32'>`,
            allowHTML: true,
          }"
        >
          <div class="flex-shrink-0 h-4 w-4">
            <img class="h-4 w-4" :src="iconURL(artifact.iconPath, 32)" />
          </div>
          <div
            class="ml-1 text-sm"
            :class="[
              rarityFgClass(artifact.afxRarity),
              noAvailabilityMarker ? null : availabilityClass,
            ]"
          >
            <span>{{ artifact.name }}</span>
            <template v-if="showTier">
              <span> (T{{ artifact.tier_number }})</span>
            </template>
            <template v-if="artifact.afxRarity > 0">
              <span>, {{ artifact.rarity }}</span>
            </template>
          </div>
        </div>
      </div>
    </component>
  </template>
</template>

<script>
import { iconURL } from "@/utils";

export default {
  props: {
    artifact: Object,
    showTier: Boolean,
    plainText: Boolean,
    noLink: Boolean,
    noAvailabilityMarker: Boolean,
  },

  computed: {
    availabilityClass() {
      return !this.artifact.available_from_missions
        ? "text-red-900 dagger"
        : this.artifact.notDroppableInPractice
        ? "text-red-900 Dagger"
        : "";
    },
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

<style scoped>
.dagger:after {
  content: "\2020";
}

.Dagger:after {
  content: "\2021";
}
</style>
