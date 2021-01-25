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
import { getLocalStorage, setLocalStorage } from "./utils";

export default {
  components: { ArtifactGrid },

  props: {
    progress: Object,
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
};
</script>
