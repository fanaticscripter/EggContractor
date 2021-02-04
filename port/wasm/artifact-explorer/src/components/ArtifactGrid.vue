<template>
  <div>
    <div class="mt-4 mb-2 text-sm font-medium text-gray-700">Visual artifacts picker</div>
    <div class="my-2 text-xs text-gray-500">
      &dagger; Not available from missions.
    </div>
    <div class="my-2 text-sm text-gray-700">Artifacts</div>
    <artifact-grid-section :families="sectionArtifacts" />
    <div class="my-2 text-sm text-gray-700">Stones &amp; fragments</div>
    <artifact-grid-section :families="sectionStones" />
    <div class="my-2 text-sm text-gray-700">Ingredients</div>
    <artifact-grid-section :families="sectionIngredients" />
  </div>
</template>

<script>
import ArtifactGridSection from "@/components/ArtifactGridSection.vue";

import { stringCmp } from "@/utils";

export default {
  components: {
    ArtifactGridSection,
  },

  props: {
    artifacts: Array,
  },

  computed: {
    sections() {
      const items = this.artifacts
        .filter(artifact => artifact.afxRarity === 0)
        .sort((artifact1, artifact2) => stringCmp(artifact1.sortKey, artifact2.sortKey));
      const families = [];
      let family = {};
      let tiers = [];
      for (const item of items) {
        if (item.family.id !== family.id) {
          if (tiers.length > 0) {
            family.tiers = tiers;
            families.push(family);
            tiers = [];
          }
          family = item.family;
        }
        tiers.push(item);
      }
      if (tiers.length > 0) {
        family.tiers = tiers;
        families.push(family);
        tiers = [];
      }

      const artifacts = [];
      const stones = [];
      const ingredients = [];
      for (const family of families) {
        switch (family.type) {
          case "Artifact":
            artifacts.push(family);
            break;
          case "Stone":
            stones.push(family);
            break;
          case "Ingredient":
            ingredients.push(family);
            break;
        }
      }
      return {
        artifacts,
        stones,
        ingredients,
      };
    },

    sectionArtifacts() {
      const { artifacts } = this.sections;
      return artifacts;
    },

    sectionStones() {
      const { stones } = this.sections;
      return stones;
    },

    sectionIngredients() {
      const { ingredients } = this.sections;
      return ingredients;
    },
  },
};
</script>
