<template>
  <section>
    <h2 class="font-medium">
      {{ sectionTitle }}
      <span
        class="ml-0.5 text-indigo-500 hover:text-indigo-600 font-normal underline cursor-pointer select-none"
        @click="$emit('toggle')"
      >
        {{ visible ? "collapse section" : "expand section" }}
      </span>
    </h2>
    <div v-show="visible">
      <slot></slot>
    </div>
  </section>
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  props: {
    sectionTitle: {
      type: String,
      required: true,
    },
    visible: {
      type: Boolean,
      required: true,
    },
  },
  emits: ["toggle"],
});
</script>
