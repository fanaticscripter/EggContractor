<template>
  <template v-if="days > 0">
    {{ days }}:{{ hours.toString().padStart(2, "0") }}:{{ minutes.toString().padStart(2, "0") }}:{{ seconds.toString().padStart(2, "0") }}
  </template>
  <template v-else>
    {{ hours }}:{{ minutes.toString().padStart(2, "0") }}:{{ seconds.toString().padStart(2, "0") }}
  </template>
</template>

<script>
export default {
  props: {
    deadline: Number,
  },
  data() {
    let remainingSeconds = this.deadline - Date.now() / 1000;
    if (remainingSeconds < 0) {
      remainingSeconds = 0;
    }
    return {
      remainingSeconds,
    };
  },
  mounted() {
    const intervalId = setInterval(() => {
      const remainingSeconds = this.deadline - Date.now() / 1000;
      if (remainingSeconds < 0) {
        this.remainingSeconds = 0;
        clearInterval(intervalId);
      } else {
        this.remainingSeconds = remainingSeconds;
      }
    }, 200);
  },
  computed: {
    days () {
      return Math.floor(this.remainingSeconds / 86400);
    },
    hours() {
      return Math.floor(this.remainingSeconds / 3600) % 24;
    },
    minutes() {
      return Math.floor(this.remainingSeconds / 60) % 60;
    },
    seconds() {
      return Math.floor(this.remainingSeconds) % 60;
    },
  },
};
</script>
