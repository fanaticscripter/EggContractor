<template>
  <div class="space-y-2">
    <simulator-item-select
      v-for="(entry, index) in modelValue"
      :key="entry.rowid"
      :modelValue="entry"
      @update:modelValue="updateItem(index, $event)"
      @delete="() => deleteItem(index)"
    />
    <button
      type="button"
      class="-ml-1 bg-white p-1 rounded-md flex items-center focus:outline-none focus:ring-2 focus:ring-blue-500"
      @click="addItem"
    >
      <span
        class="w-7 h-7 rounded-full border-2 border-dashed border-gray-300 flex items-center justify-center text-gray-400"
      >
        <!-- Heroicon name: solid/plus -->
        <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
          <path
            fill-rule="evenodd"
            d="M10 5a1 1 0 011 1v3h3a1 1 0 110 2h-3v3a1 1 0 11-2 0v-3H6a1 1 0 110-2h3V6a1 1 0 011-1z"
            clip-rule="evenodd"
          />
        </svg>
      </span>
      <span class="ml-2 text-sm text-gray-900">Add another artifact</span>
    </button>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType, toRefs } from "vue";
import { v4 as uuidv4 } from "uuid";

import { ItemSelectSpec } from "@/types";
import SimulatorItemSelect from "@/components/SimulatorItemSelect.vue";

export default defineComponent({
  components: {
    SimulatorItemSelect,
  },
  props: {
    modelValue: {
      type: Array as PropType<ItemSelectSpec[]>,
      required: true,
    },
  },
  emits: {
    "update:modelValue": (payload: ItemSelectSpec[]) => true,
  },
  setup(props, { emit }) {
    const { modelValue } = toRefs(props);
    const addItem = () => {
      const updated = [...modelValue.value];
      updated.push({
        id: null,
        count: 1,
        rowid: uuidv4(),
      });
      emit("update:modelValue", updated);
    };
    const updateItem = (index: number, spec: ItemSelectSpec) => {
      const updated = [...modelValue.value];
      updated[index] = spec;
      emit("update:modelValue", updated);
    };
    const deleteItem = (index: number) => {
      const updated = [...modelValue.value];
      updated.splice(index, 1);
      emit("update:modelValue", updated);
    };
    return {
      addItem,
      updateItem,
      deleteItem,
    };
  },
});
</script>
