<template>
  <svg
    aria-hidden="true"
    focusable="false"
    data-prefix="fas"
    data-icon="sort"
    role="img"
    xmlns="http://www.w3.org/2000/svg"
    viewBox="0 0 320 512"
    class="ml-1 h-3.5 w-3.5"
  >
    <path
      v-if="direction === 1"
      fill="currentColor"
      d="M279 224H41c-21.4 0-32.1-25.9-17-41L143 64c9.4-9.4 24.6-9.4 33.9 0l119 119c15.2 15.1 4.5 41-16.9 41z"
    />
    <path
      v-else-if="direction === -1"
      fill="currentColor"
      d="M41 288h238c21.4 0 32.1 25.9 17 41L177 448c-9.4 9.4-24.6 9.4-33.9 0L24 329c-15.1-15.1-4.4-41 17-41z"
    />
    <path
      v-else
      fill="currentColor"
      d="M41 288h238c21.4 0 32.1 25.9 17 41L177 448c-9.4 9.4-24.6 9.4-33.9 0L24 329c-15.1-15.1-4.4-41 17-41zm255-105L177 64c-9.4-9.4-24.6-9.4-33.9 0L24 183c-15.1 15.1-4.4 41 17 41h238c21.4 0 32.1-25.9 17-41z"
    />
  </svg>
</template>

<script>
export default {
  props: {
    // 1: Ascending (up arrow)
    // 0: Not sorted (both arrows)
    // -1: Descending (down arrow)
    direction: Number,
  },
};
</script>
