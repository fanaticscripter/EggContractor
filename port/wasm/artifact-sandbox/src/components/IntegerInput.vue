<template>
  <input
    :id="id"
    :name="id"
    type="number"
    :min="min"
    :max="max"
    v-model.number="value"
    class="bg-dark-20 block w-full pl-10 pt-2.5 pb-2 sm:text-sm rounded-md"
    :class="
      invalid
        ? 'border-red-300 text-red-500 placeholder-red-300 focus:outline-none focus:ring-red-500 focus:border-red-500'
        : 'focus:outline-none focus:ring-blue-500 focus:border-blue-500'
    "
  />
</template>

<script>
export default {
  props: {
    id: String,
    modelValue: Number,
    min: Number,
    max: Number,
  },

  data() {
    return {
      value: this.modelValue,
    };
  },

  emits: ["update:modelValue"],

  computed: {
    invalid() {
      if (!Number.isInteger(this.value)) {
        return true;
      }
      if (this.max !== undefined && this.value > this.max) {
        return true;
      }
      if (this.min !== undefined && this.value < this.min) {
        return true;
      }
      return false;
    },
  },

  watch: {
    modelValue() {
      this.value = this.modelValue;
    },

    value() {
      if (this.invalid) {
        return;
      }
      this.$emit("update:modelValue", this.value);
    },
  },
};
</script>
