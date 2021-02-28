<template>
  <span
    v-tippy="{
      content: noteList[number - 1][1],
    }"
    class="relative"
    :class="dagger ? 'Dagger' : null"
  >
    <slot>{{ noteList[number - 1][0] }}</slot
    ><sup class="px-0.5">
      <span v-if="showFootnote">{{ number }}</span>
      <span v-else data-html2canvas-ignore>?</span>
    </sup>
  </span>
</template>

<script>
export default {
  props: {
    noteList: {
      type: Array,
      required: true,
    },
    number: {
      type: Number,
      required: true,
    },
    dagger: Boolean,
    showFootnote: Boolean,
  },
};
</script>

<style scoped>
sup {
  color: #a6a6a6;
}

.Dagger {
  color: #ffc601;
}

.Dagger::before {
  content: "\2020";
  position: absolute;
  top: -0.25em;
  left: -0.75em;
  font-size: 75%;
}
</style>
