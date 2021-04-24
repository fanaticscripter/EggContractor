<template>
  <div class="">
    <div class="h-3 relative rounded-full overflow-hidden">
      <div class="w-full h-full bg-gray-200 absolute"></div>
      <div
        class="h-full absolute rounded-full bg-green-500"
        :class="{ 'ProgressBar--striped': !progress.stopped }"
        :style="{ width: finishedPercentage }"
      ></div>
    </div>
    <div class="flex text-xs tabular-nums">
      <span class="flex flex-1 justify-start pl-0.5">Elapsed {{ elapsed }}</span>
      <span class="flex flex-1 justify-center">{{ finishedPercentage }} @ {{ rateDisplay }}</span>
      <span class="flex flex-1 justify-end pr-0.5">
        <template v-if="progress.stopped">&ndash;</template>
        <template v-else>ETA {{ eta }}</template>
      </span>
    </div>
  </div>
</template>

<script lang="ts">
import { SimulationProgress } from "@/types";
import { computed, defineComponent, PropType, toRefs } from "vue";

function formatDuration(seconds: number) {
  const mm = Math.floor(seconds / 60);
  const ss = Math.floor(seconds - mm * 60);
  return `${mm}:${String(ss).padStart(2, "0")}`;
}

export default defineComponent({
  props: {
    progress: {
      type: Object as PropType<SimulationProgress>,
      required: true,
    },
  },
  setup(props) {
    const { progress } = toRefs(props);
    const elapsed = computed(() => formatDuration(progress.value.secondsElapsed));
    const finishedPercentage = computed(() =>
      progress.value.finishedTrials >= progress.value.totalTrials
        ? "100%"
        : ((progress.value.finishedTrials / progress.value.totalTrials) * 100).toFixed(1) + "%"
    );
    const rate = computed(() =>
      progress.value.secondsElapsed > 0
        ? progress.value.finishedTrials / progress.value.secondsElapsed
        : 0
    );
    const rateDisplay = computed(() => (rate.value > 0 ? rate.value.toFixed(0) : "0") + "/s");
    const eta = computed(() =>
      rate.value === 0
        ? "\u2013"
        : formatDuration((progress.value.totalTrials - progress.value.finishedTrials) / rate.value)
    );
    return {
      elapsed,
      finishedPercentage,
      rateDisplay,
      eta,
    };
  },
});
</script>

<style lang="postcss" scoped>
/* Animation largely borrowed from Bootstrap. */
.ProgressBar--striped {
  background-image: linear-gradient(
    45deg,
    rgba(255, 255, 255, 0.15) 25%,
    transparent 25%,
    transparent 50%,
    rgba(255, 255, 255, 0.15) 50%,
    rgba(255, 255, 255, 0.15) 75%,
    transparent 75%,
    transparent
  );
  background-size: 0.75rem 0.75rem;
  animation: progress-bar-stripes 1s linear infinite;
  transition: width 200ms;
}

@keyframes progress-bar-stripes {
  0% {
    background-position: 0.75rem 0;
  }
  100% {
    background-position: 0 0;
  }
}
</style>
