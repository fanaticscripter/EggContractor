<template>
  <input
    :id="id"
    :name="id"
    type="text"
    pattern="\d+"
    :value="modelValue"
    @input="emitValue"
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
      invalid: false,
    };
  },

  emits: ["update:modelValue"],

  methods: {
    emitValue(e) {
      const value = e.target.value.trim();
      if (value.match(/^\d+$/) !== null) {
        const intValue = parseInt(value);
        if (this.max !== undefined && intValue > this.max) {
          this.invalid = true;
        } else if (this.min !== undefined && intValue < this.min) {
          this.invalid = true;
        } else {
          this.invalid = false;
          this.$emit("update:modelValue", intValue);
        }
      } else {
        this.invalid = true;
      }
    },
  },
};
</script>
