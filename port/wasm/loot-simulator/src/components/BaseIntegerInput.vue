<template>
  <input
    type="number"
    v-model.number="value"
    class="shadow-sm focus:ring-blue-500 focus:border-blue-500 block w-full sm:text-sm border-gray-300 rounded-md focus:outline-none"
    :class="
      invalid
        ? 'border-red-300 text-red-900 placeholder-red-300 focus:ring-red-500 focus:border-red-500'
        : 'focus:ring-blue-500 focus:border-blue-500'
    "
  />
</template>

<script lang="ts">
import { computed, defineComponent, PropType, ref, toRefs, watch } from "vue";

export default defineComponent({
  props: {
    modelValue: {
      type: Number,
      required: true,
    },
    max: {
      type: Number as PropType<Number | undefined>,
      default: undefined,
    },
    min: {
      type: Number as PropType<Number | undefined>,
      default: undefined,
    },
  },
  emits: {
    "update:modelValue": (payload: number) => true,
  },
  setup(props, { emit }) {
    const { modelValue, max, min } = toRefs(props);
    const value = ref(modelValue.value);
    const invalid = computed(() => {
      const v = value.value;
      if (!Number.isInteger(v)) {
        return true;
      }
      if (max.value !== undefined && v > max.value) {
        return true;
      }
      if (min.value !== undefined && v < min.value) {
        return true;
      }
      return false;
    });
    watch(modelValue, () => {
      value.value = modelValue.value;
    });
    watch(value, () => {
      if (invalid.value) {
        return;
      }
      emit("update:modelValue", value.value);
    });
    return {
      value,
      invalid,
    };
  },
});
</script>
