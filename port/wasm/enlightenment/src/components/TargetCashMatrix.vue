<!-- This example requires Tailwind CSS v2.0+ -->
<template>
  <div class="flex flex-col">
    <div class="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
      <div class="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8">
        <div class="shadow overflow-hidden border-b border-gray-200 sm:rounded-lg">
          <table class="min-w-full divide-y divide-gray-200 text-center tabular-nums">
            <thead class="bg-gray-50">
              <tr>
                <th
                  scope="col"
                  class="px-2 py-2 whitespace-pre text-sm font-medium text-gray-500"
                ></th>
                <th
                  scope="col"
                  class="px-2 py-2 whitespace-pre text-sm font-medium text-gray-500"
                  v-for="target in targets"
                  :key="target.description"
                >
                  {{ target.description }}
                </th>
              </tr>
            </thead>
            <tbody>
              <tr class="bg-white">
                <td class="px-2 py-1 whitespace-nowrap text-sm font-medium text-gray-900">
                  Discount
                </td>
                <td
                  class="px-2 py-1 whitespace-nowrap text-sm text-gray-500"
                  v-for="target in targets"
                  :key="target.description"
                >
                  {{ formatPercentage(1 - target.multiplier) }}
                </td>
              </tr>

              <tr class="bg-gray-50">
                <td class="px-2 py-1 whitespace-nowrap text-sm font-medium text-gray-900">
                  Cash target
                </td>
                <td
                  class="px-2 py-1 whitespace-nowrap text-sm text-gray-500"
                  v-for="target in targets"
                  :key="target.description"
                >
                  <base-e-i-value :value="baseTarget * target.multiplier" />
                </td>
              </tr>

              <tr class="bg-white border-b border-gray-200">
                <td class="px-2 py-1 whitespace-nowrap text-sm font-medium text-gray-900">
                  Need to earn
                </td>
                <td
                  class="px-2 py-1 whitespace-nowrap text-sm text-gray-500"
                  v-for="target in targets"
                  :key="target.description"
                >
                  <base-e-i-value :value="needToEarn(target)" />
                </td>
              </tr>
              <tr
                v-for="(m, index) in means"
                :key="m.description"
                :class="index % 2 === 1 ? 'bg-white' : 'bg-gray-50'"
              >
                <td class="px-2 py-1 whitespace-nowrap text-sm font-medium text-gray-900">
                  {{ m.description }}
                </td>
                <td
                  class="px-2 py-1 whitespace-nowrap text-sm text-gray-500"
                  v-for="target in targets"
                  :key="target.description"
                >
                  {{ m.calc(needToEarn(target), m.rate) }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType, toRefs } from "vue";

import { formatPercentage } from "@/utils";
import BaseEIValue from "@/components/BaseEIValue.vue";

type Target = {
  multiplier: number;
  description: string;
};

type Means = {
  rate: number;
  description: string;
  calc: (target: number, rate: number) => string;
};

export default defineComponent({
  components: {
    BaseEIValue,
  },
  props: {
    baseTarget: {
      type: Number,
      required: true,
    },
    current: {
      type: Number,
      required: true,
    },
    targets: {
      type: Array as PropType<Target[]>,
      required: true,
    },
    means: {
      type: Array as PropType<Means[]>,
      required: true,
    },
  },
  setup(props) {
    const { baseTarget, current } = toRefs(props);
    const needToEarn = (target: Target) =>
      Math.max(baseTarget.value * target.multiplier - current.value, 0);
    return {
      needToEarn,
      formatPercentage,
    };
  },
});
</script>
