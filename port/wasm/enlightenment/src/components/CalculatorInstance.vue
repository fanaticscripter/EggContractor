<template>
  <div>
    <h2 class="text-sm font-medium mb-1">Calculator</h2>
    <p class="text-xs mb-2">
      In this calculator, you can use values carrying Egg, Inc's OoM units (see reference table
      below). Supported operators include +, -, *, /, and ** or ^ for exponentiation. Common
      mathematical functions like exp, log (ln), log10 (lg), max, min, round, floor and ceil are
      also supported.
    </p>
    <textarea
      class="w-full resize border rounded-md border-gray-300 sm:text-sm"
      :style="{ minHeight: '4rem' }"
      autocapitalize="off"
      spellcheck="false"
      placeholder="Example: (1.254Od / 8.769Td)^2 * 2.282s"
      v-model="expr"
    ></textarea>
    <template v-if="result !== null">
      <div class="text-sm">
        {{ result }}<br />
        <span class="text-indigo-700">EI notation:</span> <base-e-i-value :value="result" />
      </div>
    </template>
    <template v-else>
      <div class="text-sm">
        Waiting for valid expression...<br />
        <span class="invisible">.</span>
      </div>
    </template>
    <hr class="my-2" />
    <p class="text-xs font-medium">Reference table</p>
    <div class="grid grid-cols-3 sm:grid-cols-6 text-xs tabular-nums mt-1">
      <div v-for="unit in units">
        <span class="inline-block" :style="{ width: '1.5rem' }">{{ unit.symbol }}</span> 10<sup>{{
          unit.oom
        }}</sup>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, ref, toRefs, watch } from "vue";

import { calculateWithOoMUnits, units } from "@/lib";
import BaseEIValue from "@/components/BaseEIValue.vue";

export default defineComponent({
  components: {
    BaseEIValue,
  },
  props: {
    // The expression.
    modelValue: {
      type: String,
      default: "",
    },
  },
  emits: {
    "update:modelValue": (newValue: string) => true,
  },
  setup(props, { emit }) {
    const { modelValue } = toRefs(props);
    const expr = ref(modelValue.value);
    const result = computed(() => calculateWithOoMUnits(expr.value));
    watch(expr, () => emit("update:modelValue", expr.value));
    return {
      expr,
      result,
      units,
    };
  },
});
</script>
