<!-- This example requires Tailwind CSS v2.0+ -->
<template>
  <div class="fixed right-6 bottom-6" :style="{ right: 'max(calc(50vw - 36.5rem), 1.5rem)' }">
    <div
      class="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-blue-500 shadow-inner cursor-pointer"
      v-tippy="{ content: 'Egg Inc. OoM-aware calculator' }"
      @click="open = true"
    >
      <!-- Heroicon name: outline/calculator -->
      <svg
        xmlns="http://www.w3.org/2000/svg"
        class="h-8 w-8 text-blue-100"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z"
        />
      </svg>
    </div>
  </div>

  <TransitionRoot as="template" :show="open">
    <Dialog
      as="div"
      static
      class="fixed z-20 inset-0 overflow-y-auto"
      @close="open = false"
      :open="open"
    >
      <div
        class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0"
      >
        <TransitionChild
          as="template"
          enter="ease-out duration-300"
          enter-from="opacity-0"
          enter-to="opacity-100"
          leave="ease-in duration-200"
          leave-from="opacity-100"
          leave-to="opacity-0"
        >
          <DialogOverlay class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" />
        </TransitionChild>

        <!-- This element is to trick the browser into centering the modal contents. -->
        <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true"
          >&#8203;</span
        >
        <TransitionChild
          as="template"
          enter="ease-out duration-300"
          enter-from="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
          enter-to="opacity-100 translate-y-0 sm:scale-100"
          leave="ease-in duration-200"
          leave-from="opacity-100 translate-y-0 sm:scale-100"
          leave-to="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
        >
          <div
            class="inline-block align-bottom bg-white rounded-lg px-4 pt-5 pb-4 text-left overflow-hidden shadow-xl transform transition-all w-full max-w-xs sm:my-8 sm:align-middle sm:max-w-xl sm:px-5 sm:py-5"
          >
            <calculator-instance class="w-full" v-model="calculatorExpr" />
          </div>
        </TransitionChild>
      </div>
    </Dialog>
  </TransitionRoot>
</template>

<script lang="ts">
import { defineComponent, ref } from "vue";
import {
  Dialog,
  DialogOverlay,
  DialogTitle,
  TransitionChild,
  TransitionRoot,
} from "@headlessui/vue";

import CalculatorInstance from "@/components/CalculatorInstance.vue";

export default defineComponent({
  components: {
    Dialog,
    DialogOverlay,
    DialogTitle,
    TransitionChild,
    TransitionRoot,
    CalculatorInstance,
  },
  setup() {
    const open = ref(false);
    // Persist between calculator component rerenders.
    const calculatorExpr = ref("");
    return {
      open,
      calculatorExpr,
    };
  },
});
</script>
