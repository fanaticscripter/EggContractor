<template>
  <input
    :id="id"
    :name="id"
    type="text"
    :pattern="valueWithUnitRegExpPattern"
    v-model="input"
    class="bg-dark-20 block w-full pl-10 pt-2.5 pb-2 sm:text-sm rounded-md"
    :class="
      !parsed
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
      input: this.raw,
      invalid: false,
    };
  },

  emits: ["update:raw", "update:value"],

  valueWithUnitRegExpPattern,

  computed: {
    parsed() {
      return parseValueWithUnit(this.input);
    },
  },

  watch: {
    input() {
      if (this.parsed === null) {
        return;
      }
      this.$emit("update:raw", this.input);
      this.$emit("update:value", this.parsed);
    },
  },

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
