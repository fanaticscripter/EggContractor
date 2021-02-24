<template>
  <div class="max-w-4xl w-full mx-auto px-4 xl:px-0 my-4">
    <artifact-set-builder :key="key" v-model:build="builds.builds[0]" />
    <configurator :key="key" v-model:config="builds.config" />
  </div>

  <hr class="border-dark-30" />

  <div class="mt-4 mb-2 text-center">
    <button
      type="button"
      class="inline-flex items-center px-3 py-2 border border-transparent shadow-sm text-sm leading-4 font-medium rounded-md bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 focus:ring-offset-dark-30"
      @click="showShareSheet = true"
    >
      Share this build
      <svg
        class="ml-2 -mr-0.5 h-4 w-4"
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 20 20"
        fill="currentColor"
        aria-hidden="true"
      >
        <path
          d="M15 8a3 3 0 10-2.977-2.63l-4.94 2.47a3 3 0 100 4.319l4.94 2.47a3 3 0 10.895-1.789l-4.94-2.47a3.027 3.027 0 000-.74l4.94-2.47C13.456 7.68 14.19 8 15 8z"
        />
      </svg>
    </button>
  </div>

  <div
    id="builds"
    class="max-w-3xl w-full mx-auto px-4 py-4 overflow-hidden bg-dark-25 rounded-xl"
  >
    <artifact-set-display :key="key" :build="builds.builds[0]" :config="builds.config" />
    <hr class="border-dark-30" />
    <artifact-sets-effects
      :key="key"
      :builds="builds"
      :showFootnotes="showShareSheet && showFootnotesWhenSharing"
    />
    <div class="mt-2 text-center text-xs text-dark-60">Built on https://ei.tcl.sh/sandbox</div>
  </div>

  <share-sheet
    :key="key"
    v-model:show="showShareSheet"
    v-model:showFootnotes="showFootnotesWhenSharing"
    :builds="builds"
  />
</template>

<script>
import ArtifactSetBuilder from "@/components/ArtifactSetBuilder.vue";
import ArtifactSetDisplay from "@/components/ArtifactSetDisplay.vue";
import ArtifactSetsEffects from "@/components/ArtifactSetsEffects.vue";
import Configurator from "@/components/Configurator.vue";
import ShareSheet from "@/components/ShareSheet.vue";

import { Builds } from "@/lib/models";

export default {
  components: {
    ArtifactSetBuilder,
    ArtifactSetDisplay,
    ArtifactSetsEffects,
    Configurator,
    ShareSheet,
  },

  props: {
    serializedBuilds: String,
  },

  data() {
    return {
      // Use an initial-path dependent key to work around the problem of
      // artifact-set-builder and configurator not updating upon manual
      // hashchange.
      key: this.serializedBuilds || "",
      builds: this.deserializeBuilds(this.serializedBuilds),
      showShareSheet: false,
      showFootnotesWhenSharing: true,
    };
  },

  methods: {
    deserializeBuilds(s) {
      let builds = Builds.newDefaultBuilds();
      if (s !== undefined) {
        try {
          builds = Builds.deserialize(s);
        } catch (e) {
          console.error(`error deserializing ${s}: ${e}`);
        }
      }
      return builds;
    },
  },

  watch: {
    builds: {
      handler() {
        window.history.replaceState(
          {},
          null,
          this.$router.resolve({
            name: "builds",
            params: { serializedBuilds: this.builds.serialize() },
          }).href
        );
      },
      deep: true,
    },
  },

  beforeRouteUpdate(to, from) {
    // Rerender on manual hashchange.
    this.builds = this.deserializeBuilds(to.params.serializedBuilds);
    this.key = to.params.serializedBuilds || "";
  },
};
</script>
