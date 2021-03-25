<template>
  <div
    class="CalendarMonth mx-auto px-4 py-4 space-y-4 bg-gray-50 border md:rounded-lg md:shadow md:px-6"
    :class="forceFullWidth ? 'CalendarMonth--full-width' : null"
  >
    <h2 class="text-center text-sm font-medium">{{ month.format("MMMM YYYY") }}</h2>
    <div class="grid grid-cols-7 gap-2">
      <div
        v-for="(weekday, index) in ['S', 'M', 'T', 'W', 'T', 'F', 'S']"
        :key="index"
        class="text-center text-sm text-gray-500"
      >
        {{ weekday }}
      </div>

      <div v-for="index in month.day()" :key="index"></div>

      <div v-for="d in dates" :key="d.date.date()">
        <div class="text-center text-sm">{{ d.date.format("D") }}</div>
        <div class="space-y-1">
          <template v-for="(event, index) in d.events" :key="index">
            <tippy
              v-if="eventTypesOn[event.type] !== false"
              class="flex items-center md:justify-start space-x-1"
              :class="forceFullWidth ? 'justify-start' : 'justify-center'"
            >
              <event-badge :event="event" />
              <span
                class="text-xs truncate md:inline"
                :class="[eventFgClass(event.type), forceFullWidth ? 'inline' : 'hidden']"
              >
                <template v-if="event.type != 'app-update'">
                  {{ eventCaption(event.type, event.multiplier) }}
                </template>
                <template v-else> v{{ event.version }} </template>
              </span>

              <template #content>
                <template v-if="event.type != 'app-update'">
                  <div :class="eventBrightFgClass(event.type)">
                    <p class="flex items-center space-x-1">
                      <event-badge :event="event" />
                      <span class="text-xs truncate">
                        {{ event.message }}
                      </span>
                    </p>
                    <p>
                      <span class="text-white">Starting time:</span>
                      {{ event.startTime.local().format("MM-DD HH:mm Z") }}
                    </p>
                    <p>
                      <span class="text-white">Duration:</span>
                      {{ formatDuration(event.durationSeconds) }}
                    </p>
                  </div>
                </template>

                <template v-else>
                  <div :class="eventBrightFgClass(event.type)">
                    <p class="flex items-center space-x-1">
                      <img :src="iconURL('egginc/ei_app_icon.png', 64)" class="h-4 w-4 rounded-sm" />
                      <span>App version {{ event.version }}</span>
                    </p>
                    <p>
                      <span class="text-white">Release date:</span>
                      {{ event.startTime.format("MM-DD") }}
                    </p>
                    <p>
                      <span class="text-white">Release notes:</span><br />
                      <span class="whitespace-pre-line">{{ event.releaseNotes }}</span>
                    </p>
                  </div>
                </template>
              </template>
            </tippy>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { Tippy } from "vue-tippy";
import EventBadge from "@/components/EventBadge.vue";

import { computed, toRefs } from "vue";
import dayjs from "dayjs";

import { eventCaption, eventFgClass, eventBrightFgClass } from "@/lib";
import { iconURL } from "@/utils";

export default {
  components: {
    EventBadge,
    Tippy,
  },

  props: {
    // YYYY-MM.
    monthStr: {
      type: String,
      required: true,
    },
    date2events: {
      type: Object,
      required: true,
    },
    eventTypesOn: {
      type: Object,
      required: true,
    },
    forceFullWidth: Boolean,
  },

  setup(props) {
    const { monthStr, date2events } = toRefs(props);
    const month = computed(() => dayjs(monthStr.value));
    const dates = computed(() => {
      const monthNumber = month.value.month();
      const dates = [];
      for (let date = month.value.clone(); date.month() === monthNumber; date = date.add(1, "d")) {
        dates.push({
          date,
          events: date2events.value[date.date()] || [],
        });
      }
      return dates;
    });

    const formatDuration = seconds => {
      if (seconds < 0) {
        return "-" + formatDuration(-seconds);
      }
      if (seconds < 1) {
        return "0m";
      }
      const dd = Math.floor(seconds / 86400);
      seconds -= dd * 86400;
      const hh = Math.floor(seconds / 3600);
      seconds -= hh * 3600;
      const mm = Math.floor(seconds / 60);
      let s = "";
      if (dd > 0) {
        s += `${dd}d`;
      }
      if (hh > 0) {
        s += `${hh}h`;
      }
      if (mm > 0) {
        s += `${mm}m`;
      }
      return s;
    };

    return {
      month,
      dates,
      eventCaption,
      eventFgClass,
      eventBrightFgClass,
      formatDuration,
      iconURL,
    };
  },
};
</script>

<style scoped>
.CalendarMonth {
  width: 100%;
}

.CalendarMonth.CalendarMonth--full-width {
  width: 744px;
}

@media (min-width: 768px) {
  .CalendarMonth {
    width: 744px;
  }
}
</style>
