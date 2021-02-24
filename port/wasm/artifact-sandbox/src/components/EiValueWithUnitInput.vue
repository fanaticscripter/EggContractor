<template>
  <input
    :id="id"
    :name="id"
    type="text"
    :pattern="valueWithUnitRegExpPattern"
    :value="raw"
    @input="emitValue"
    class="bg-dark-20 block w-full pl-10 sm:text-sm rounded-md"
    :class="
      invalid
        ? 'border-red-300 text-red-500 placeholder-red-300 focus:outline-none focus:ring-red-500 focus:border-red-500'
        : 'focus:outline-none focus:ring-blue-500 focus:border-blue-500'
    "
  />
</template>

<script>
import { valueWithUnitRegExpPattern, parseValueWithUnit } from "@/lib/utils/utils";

export default {
  props: {
    id: String,
    raw: String,
    value: Number,
  },

  data() {
    return {
      invalid: false,
    };
  },

  emits: ["update:raw", "update:value"],

  valueWithUnitRegExpPattern,

  methods: {
    emitValue(e) {
      const value = e.target.value.trim();
      const parsed = parseValueWithUnit(value);
      if (parsed !== null) {
        this.invalid = false;
        this.$emit("update:raw", value);
        this.$emit("update:value", parsed);
      } else {
        this.invalid = true;
      }
    },
  },
};
</script>
