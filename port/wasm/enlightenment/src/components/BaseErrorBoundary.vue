<template>
  <div v-if="error">
    <div class="text-center break-words text-red-500">
      <svg class="relative inline h-4 w-4 -top-px mr-1" viewBox="0 0 20 20" fill="currentColor">
        <path
          fill-rule="evenodd"
          d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z"
          clip-rule="evenodd"
        />
      </svg>
      {{ error.toString() }}
    </div>
  </div>
  <slot v-else></slot>
</template>

<script lang="ts">
import { defineComponent, onErrorCaptured, ref } from "vue";

export default defineComponent({
  setup() {
    const error = ref(null as Error | null);
    onErrorCaptured(err => {
      error.value = err as Error;
      return false;
    });
    return {
      error,
    };
  },
});
</script>
