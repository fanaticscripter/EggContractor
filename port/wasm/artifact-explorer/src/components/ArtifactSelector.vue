<template>
  <div class="my-4 space-y-2">
    <div class="w-full max-w-xs">
      <label for="artifact" class="block text-sm font-medium text-gray-700">
        What does this artifact do and how do I get it?
      </label>
      <select
        id="artifact"
        name="artifact"
        class="mt-1 block w-full pl-3 pr-10 py-1 text-sm bg-gray-50 border-gray-300 focus:outline-none focus:ring-green-500 focus:border-green-500 rounded-md"
        v-model="artifactId"
      >
        <option value="">-- Select an item --</option>
        <optgroup label="Artifacts">
          <option v-for="item in artifactsSection" :key="item.id" :value="item.id">
            {{ item.selectName }}
          </option>
        </optgroup>
        <optgroup label="Stones">
          <option v-for="item in stonesSection" :key="item.id" :value="item.id">
            {{ item.selectName }}
          </option>
        </optgroup>
        <optgroup label="Ingredients">
          <option v-for="item in ingredientsSection" :key="item.id" :value="item.id">
            {{ item.selectName }}
          </option>
        </optgroup>
      </select>
    </div>
    <p class="text-xs text-gray-500">
      You may use this dropdown or the visual artifacts picker below, or click on any artifact name
      on the page, including, for instance, results in a mission query.
    </p>
  </div>
</template>

<script>
import ArtifactName from "@/components/ArtifactName.vue";

import { stringCmp } from "@/utils";

export default {
  components: {
    ArtifactName,
  },

  props: {
    initialArtifactId: String,
    artifacts: Array,
  },

  data() {
    return {
      commonArtifacts: this.artifacts.filter(artifact => artifact.afxRarity === 0),
      artifactId: this.initialArtifactId || "",
    };
  },

  computed: {
    artifactsSection() {
      return this.sectionByType("Artifact");
    },

    stonesSection() {
      return this.sectionByType("Stone");
    },

    ingredientsSection() {
      return this.sectionByType("Ingredient");
    },
  },

  watch: {
    artifactId() {
      if (this.artifactId === "") {
        this.$router.push({
          name: "home",
        });
      } else {
        this.$router.push({
          name: "artifact",
          params: {
            artifactId: this.artifactId,
          },
        });
      }
    },
  },

  methods: {
    sectionByType(typeName) {
      return this.commonArtifacts
        .filter(item => item.family.type === typeName)
        .map(item => ({
          id: item.id,
          selectName: `${item.family.name}, ${item.tier_name} (T${item.tier_number})`,
        }))
        .sort((item1, item2) => stringCmp(item1.selectName, item2.selectName));
    },
  },
};
</script>
