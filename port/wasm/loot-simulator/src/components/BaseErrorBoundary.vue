<template>
  <slot v-if="error" name="error" :error="error">
    <div v-if="error" class="text-sm text-red-500">{{ error.toString() }}</div>
  </slot>
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
