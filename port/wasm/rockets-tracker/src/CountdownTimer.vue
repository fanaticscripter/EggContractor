<template>
  {{ hours }}:{{ minutes.toString().padStart(2, "0") }}:{{ seconds.toString().padStart(2, "0") }}
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
