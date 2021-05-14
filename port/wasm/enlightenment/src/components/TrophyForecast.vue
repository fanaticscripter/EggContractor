<template>
  <p v-if="lastRefreshedPopulation < targetPopulation" class="text-sm">
    <template v-if="habSpace >= targetPopulation"
      >Enlightenment {{ trophyLevel }} Trophy forecast:
    </template>
    <template v-else>
      Enlightenment Diamond Trophy forecast, assuming sufficient hab space can be unlocked in time:
    </template>
    <template v-if="completionForecast">
      <span class="text-green-500 whitespace-nowrap">
        {{ completionForecast.format("LLL z") }}
      </span>
      <template v-if="completionForecastDays !== null && completionForecastDays > 0">
        ({{ completionForecastDays.toFixed(1) }} days)
      </template>
    </template>
    <template v-else>Never</template>
  </p>
  <p
    v-if="
      trophyLevel === 'Diamond' &&
      completionForecastDays !== null &&
      completionForecastDays > 0 &&
      completionForecastDays < 1
    "
    class="text-base"
  >
    &#x1f90f;&#x1f90f;&#x1f3fb;&#x1f90f;&#x1f3fc;&#x1f90f;&#x1f3fd;&#x1f90f;&#x1f3fe;&#x1f90f;&#x1f3ff;
  </p>
</template>

<script lang="ts">
import { computed, defineComponent, onBeforeUnmount, ref, toRefs } from "vue";
import dayjs from "dayjs";
import advancedFormat from "dayjs/plugin/advancedFormat";
import localizedFormat from "dayjs/plugin/localizedFormat";
import relativeTime from "dayjs/plugin/relativeTime";
import timezone from "dayjs/plugin/timezone";
import utc from "dayjs/plugin/utc";

// Note that timezone abbreviation may not work due to
// https://github.com/iamkun/dayjs/issues/1154, in which case the GMT offset is
// shown.
dayjs.extend(advancedFormat);
dayjs.extend(localizedFormat);
dayjs.extend(relativeTime);
dayjs.extend(timezone);
dayjs.extend(utc);

export default defineComponent({
  props: {
    trophyLevel: {
      type: String,
      required: true,
    },
    lastRefreshedPopulation: {
      type: Number,
      required: true,
    },
    lastRefreshedTimestamp: {
      type: Number,
      required: true,
    },
    targetPopulation: {
      type: Number,
      required: true,
    },
    habSpace: {
      type: Number,
      required: true,
    },
    offlineIHR: {
      type: Number,
      required: true,
    },
  },
  setup(props) {
    const {
      lastRefreshedPopulation,
      lastRefreshedTimestamp,
      targetPopulation,
      offlineIHR,
    } = toRefs(props);
    const completionForecast = computed(() => {
      if (lastRefreshedPopulation.value < targetPopulation.value && offlineIHR.value > 0) {
        const timeToCompleteSeconds =
          ((targetPopulation.value - lastRefreshedPopulation.value) / offlineIHR.value) * 60;
        return dayjs(lastRefreshedTimestamp.value + timeToCompleteSeconds * 1000);
      } else {
        return null;
      }
    });
    const now = ref(dayjs());
    const completionForecastDays = computed(() =>
      completionForecast.value !== null
        ? completionForecast.value.diff(now.value, "day", true)
        : null
    );
    const refreshIntervalId = setInterval(() => {
      now.value = dayjs();
    }, 60000);
    onBeforeUnmount(() => {
      clearInterval(refreshIntervalId);
    });
    return {
      completionForecast,
      completionForecastDays,
    };
  },
});
</script>
