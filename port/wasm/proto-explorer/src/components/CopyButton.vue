<template>
  <button @click="copy()" v-tippy="{ content: tooltip }">
    <svg class="h-4 w-4 text-gray-500 hover:text-gray-700" viewBox="0 0 20 20" fill="currentColor">
      <path d="M8 2a1 1 0 000 2h2a1 1 0 100-2H8z" />
      <path
        d="M3 5a2 2 0 012-2 3 3 0 003 3h2a3 3 0 003-3 2 2 0 012 2v6h-4.586l1.293-1.293a1 1 0 00-1.414-1.414l-3 3a1 1 0 000 1.414l3 3a1 1 0 001.414-1.414L10.414 13H15v3a2 2 0 01-2 2H5a2 2 0 01-2-2V5zM15 11h2a1 1 0 110 2h-2v-2z"
      />
    </svg>
  </button>
</template>

<script>
import copyTextToClipboard from "copy-text-to-clipboard";

export default {
  props: {
    content: {
      type: String,
      required: true,
    },
    tooltip: {
      type: String,
      default: "Copy",
    },
  },

  methods: {
    copy() {
      copyTextToClipboard(this.content);
    },
  },
};
</script>
