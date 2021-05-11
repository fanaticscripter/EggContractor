<template>
  <span
    v-tippy="{
      content: prefix + formatEIValue(value, { scientific: true }) + suffix,
      allowHTML: true,
    }"
    >{{ prefix }}{{ formatEIValue(value) }}{{ suffix }}</span
  >
</template>

<script lang="ts">
import { formatEIValue } from "@/lib";
import { defineComponent } from "@vue/runtime-core";

// TODO: click to copy
export default defineComponent({
  props: {
    value: {
      type: Number,
      required: true,
    },
    prefix: {
      type: String,
      default: "",
    },
    suffix: {
      type: String,
      default: "",
    },
  },
  setup() {
    return {
      formatEIValue,
    };
  },
});
</script>
