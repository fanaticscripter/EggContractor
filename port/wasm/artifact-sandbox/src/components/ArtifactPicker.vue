<template>
  <div class="max-w-sm w-full mx-auto my-2">
    <label class="block text-sm font-medium">Artifact #{{ artifactIndex }}</label>
    <artifact-picker-item-select v-model="selectedArtifact.id" type="artifact" class="mt-1" />
    <template v-for="i in 3" :key="i">
      <div class="flex items-center space-x-1 mt-1" :class="i > numSlots ? 'opacity-50' : null">
        <div class="flex-shrink-0 h-8 w-8 rounded-lg bg-dark-20">
          <img
            class="h-8 w-8"
            src="https://eggincassets.tcl.sh/64/egginc-extras/icon_afx_stone_slot.png"
          />
        </div>
        <artifact-picker-item-select
          v-model="selectedArtifact.stones[i - 1]"
          type="stone"
          :disabled="i > numSlots"
          class="flex-grow"
        />
      </div>
    </template>
  </div>
</template>

<script>
import { artifactFromId } from "@/lib/data";
import ArtifactPickerItemSelect from "./ArtifactPickerItemSelect.vue";

export default {
  components: {
    ArtifactPickerItemSelect,
  },

  props: {
    artifactIndex: {
      type: Number,
      required: true,
    },
    artifact: {
      type: Object,
      required: true,
    },
  },

  data() {
    return {
      selectedArtifact: this.artifact,
    };
  },

  emits: ["update:artifact"],

  watch: {
    selectedArtifact: {
      handler() {
        this.$emit("update:artifact", this.selectedArtifact);
      },
      deep: true,
    },
  },

  computed: {
    numSlots() {
      return artifactFromId(this.artifact.id)?.slots || 0;
    },
  },
};
</script>
