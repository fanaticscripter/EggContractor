<template>
  <button
    type="submit"
    class="block px-2.5 py-1.5 border border-transparent text-xs font-medium rounded shadow-sm text-white bg-blue-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
    :class="formValid ? 'hover:bg-blue-700' : 'cursor-default'"
    :disabled="!formValid"
  >
    Request
  </button>
</template>

<script>
export default {
  props: {
    formValid: {
      type: Boolean,
      default: true,
    },
  },
};
</script>
