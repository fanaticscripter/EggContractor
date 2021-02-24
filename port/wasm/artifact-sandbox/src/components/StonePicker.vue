<template>
  <div class="flex items-center space-x-1 mt-1" :class="disabled ? 'opacity-50' : null">
    <div class="flex-shrink-0 h-8 w-8 rounded-lg bg-dark-20">
      <img
        class="h-8 w-8"
        src="https://eggincassets.tcl.sh/64/egginc-extras/icon_afx_stone_slot.png"
      />
    </div>
    <select
      :id="domId"
      :name="domId"
      class="block w-full flex-shrink pl-3 pr-10 py-2 text-base bg-dark-20 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md"
      :value="stoneId"
      @input="$emit('update:stoneId', $event.target.value)"
      :disabled="!!disabled"
    >
      <optgroup
        v-for="group in stoneOptionsGrouped"
        :key="group.familyName"
        :label="group.familyName"
      >
        <option v-for="option in group.options" :key="option.id" :value="option.id">
          {{ option.name }}
        </option>
      </optgroup>
    </select>
  </div>
</template>

<script>
import { stoneOptionsGrouped } from "@/lib/data";

export default {
  props: {
    stoneId: String,
    domId: {
      type: String,
      required: true,
    },
    disabled: Boolean,
  },

  data() {
    return {
      stoneOptionsGrouped,
    };
  },

  emits: ["update:stoneId"],
};
</script>
