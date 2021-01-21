<!-- Based on https://css-tricks.com/building-progress-ring-quickly/ -->

<template>
  <svg :height="radius * 2" :width="radius * 2">
    <circle
      class="text-gray-200"
      stroke="currentColor"
      :stroke-width="stroke"
      fill="transparent"
      :r="normalizedRadius"
      :cx="radius"
      :cy="radius"
    />
    <circle
      stroke="currentColor"
      :stroke-dasharray="circumference + ' ' + circumference"
      :style="{ strokeDashoffset: strokeDashoffset }"
      :stroke-width="stroke"
      fill="transparent"
      :r="normalizedRadius"
      :cx="radius"
      :cy="radius"
    />
  </svg>
</template>

<script>
export default {
  props: {
    radius: Number,
    stroke: Number,
    duration: Number,
    deadline: Number,
  },
  data() {
    const normalizedRadius = this.radius - this.stroke * 2;
    const circumference = normalizedRadius * 2 * Math.PI;
    let progress = 0;
    if (this.deadline !== 0) {
      let remainingSeconds = this.deadline - Date.now() / 1000;
      if (remainingSeconds < 0) {
        remainingSeconds = 0;
      }
      progress = (this.duration - remainingSeconds) / this.duration;
    }
    return {
      normalizedRadius,
      circumference,
      progress,
    };
  },
  mounted() {
    if (this.deadline === 0) {
      return;
    }
    const intervalId = setInterval(() => {
      const remainingSeconds = this.deadline - Date.now() / 1000;
      if (remainingSeconds < 0) {
        this.progress = 1;
        clearInterval(intervalId);
      } else {
        this.progress = (this.duration - remainingSeconds) / this.duration;
      }
    }, 1000);
  },
  computed: {
    strokeDashoffset() {
      return this.circumference - this.progress * this.circumference;
    },
  },
};
</script>

<style scoped>
circle {
  transition: stroke-dashoffset 0.2s;
  transform: rotate(90deg);
  transform-origin: 50% 50%;
}
</style>
