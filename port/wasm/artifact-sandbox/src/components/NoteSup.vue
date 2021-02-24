<template>
  <sup class="px-0.5">
    <span v-if="showFootnote">{{ footnoteNumber }}</span>
    <span v-else data-html2canvas-ignore>?</span>
  </sup>
</template>

<script>
export default {
  props: {
    footnoteNumber: Number,
    showFootnote: Boolean,
  },
};
</script>

<style scoped>
sup {
  color: #a6a6a6;
}
</style>
