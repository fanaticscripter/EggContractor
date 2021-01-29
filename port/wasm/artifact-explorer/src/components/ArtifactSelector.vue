<template>
  <div class="my-4 space-y-2">
    <div class="w-full max-w-xs">
      <label for="artifact" class="block text-sm font-medium text-gray-700">
        How do I get this artifact?
      </label>
      <select
        id="artifact"
        name="artifact"
        class="mt-1 block w-full pl-3 pr-10 py-1 text-sm bg-gray-50 border-gray-300 focus:outline-none focus:ring-green-500 focus:border-green-500 rounded-md"
        v-model="artifactId"
      >
        <option value="">-- Select an artifact --</option>
        <option v-for="artifact in commonArtifacts" :key="artifact.itemId" :value="artifact.itemId">
          <artifact-name :artifact="artifact" :showTier="true" :plainText="true" />
        </option>
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
      artifactId: this.initialArtifactId || "",
    };
  },

  computed: {
    commonArtifacts() {
      return this.artifacts
        .filter(artifact => artifact.rarity === 0)
        .sort((artifact1, artifact2) => stringCmp(artifact1.sortKey, artifact2.sortKey));
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
};
</script>
