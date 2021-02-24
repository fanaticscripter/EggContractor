<template>
  <div class="max-w-sm w-full mx-auto my-2">
    <label :for="domId" class="block text-sm font-medium">Artifact #{{ artifactIndex }}</label>
    <select
      :id="domId"
      :name="domId"
      class="mt-1 block w-full pl-3 pr-10 py-2 text-base bg-dark-20 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md"
      v-model="selectedArtifact.id"
    >
      <optgroup
        v-for="group in artifactOptionsGrouped"
        :key="group.familyName"
        :label="group.familyName"
      >
        <option v-for="option in group.options" :key="option.id" :value="option.id">
          {{ option.name }}
        </option>
      </optgroup>
    </select>

    <template v-for="i in 3" :key="i">
      <stone-picker
        v-model:stoneId="selectedArtifact.stones[i - 1]"
        :domId="`${domId}-stone-${i}`"
        :disabled="i > numSlots"
      />
    </template>
  </div>
</template>

<script>
import StonePicker from "@/components/StonePicker.vue";

import { artifactOptions, artifactOptionsGrouped } from "@/lib/data";

export default {
  components: {
    StonePicker,
  },

  props: {
    artifactIndex: {
      type: Number,
      required: true,
    },
    domId: {
      type: String,
      required: true,
    },
    artifact: {
      type: Object,
      required: true,
    },
  },

  data() {
    return {
      artifactOptions,
      artifactOptionsGrouped,
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
      for (const option of this.artifactOptions) {
        if (this.selectedArtifact.id === option.id) {
          return option.slots;
        }
      }
      return 0;
    },
  },
};
</script>
