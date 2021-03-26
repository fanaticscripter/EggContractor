<template>
  <div class="max-w-4xl px-4 lg:px-0 mx-auto my-4 space-y-2">
    <div class="flex justify-start sm:justify-center">
      <div>
        <div class="relative flex items-start">
          <div class="flex items-center h-5">
            <input
              id="useUtcDates"
              name="useUtcDates"
              type="checkbox"
              class="focus:ring-green-500 h-4 w-4 text-green-600 border-gray-300 rounded"
              v-model="useUtcDates"
            />
          </div>
          <div class="ml-2 flex items-center space-x-1">
            <label for="useUtcDates" class="text-sm text-gray-600">Use UTC dates</label>
            <info
              class="cursor-help"
              v-tippy="{
                content: `Events are marked on the calendar by their respective starting dates.
                  If this option is checked, the date is calculated under UTC, which basically
                  always avoids putting events of consecutive days under the same date, or
                  leaving some dates empty (these happens when, for instance, two days' worth
                  of events start at 1:00am and 23:00pm of the same day in local timezone).
                  Uncheck to use local timezone instead. Note that you can always hover/click
                  on an event to reveal the exact starting time in local timezone.`,
              }"
            />
          </div>
        </div>

        <div class="relative flex md:hidden items-start">
          <div class="flex items-center h-5">
            <input
              id="forceFullWidth"
              name="forceFullWidth"
              type="checkbox"
              class="focus:ring-green-500 h-4 w-4 text-green-600 border-gray-300 rounded"
              v-model="forceFullWidth"
            />
          </div>
          <div class="ml-2 flex items-center space-x-1">
            <label for="forceFullWidth" class="text-sm text-gray-600">
              Display boost multipliers
            </label>
            <info
              class="cursor-help"
              v-tippy="{
                content: `Display boost multipliers in addition to event labels directly in the
                  calendar. Due to screen width restrictions, checking this option comes with
                  the downside of requiring horizontal scrolling. Regardless of whether you
                  turn this on, you can always click on an event label to reveal more details,
                  including the multiplier.`,
              }"
            />
          </div>
        </div>

        <div class="relative hidden 2col:flex items-start">
          <div class="flex items-center h-5">
            <input
              id="forceSingleColumn"
              name="forceSingleColumn"
              type="checkbox"
              class="focus:ring-green-500 h-4 w-4 text-green-600 border-gray-300 rounded"
              v-model="forceSingleColumn"
            />
          </div>
          <div class="ml-2 flex items-center space-x-1">
            <label for="forceSingleColumn" class="text-sm text-gray-600">
              Single column view
            </label>
            <info
              class="cursor-help"
              v-tippy="{
                content: `When checked, show only one month per row.`,
              }"
            />
          </div>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4">
      <div v-for="[type, name] in eventTypes" :key="type" class="relative flex items-start">
        <div class="flex items-center h-5">
          <input
            :id="`show-${type}`"
            :name="`show-${type}`"
            type="checkbox"
            class="focus:ring-green-500 h-4 w-4 text-green-600 border-gray-300 rounded"
            v-model="eventTypesOn[type]"
            @change="persistEventTypeOn(type, $event.target.checked)"
          />
        </div>
        <div class="ml-2">
          <label :for="`show-${type}`" class="flex items-center space-x-1">
            <event-badge :event="{ type }" />
            <span class="text-sm text-gray-600">{{ capitalize(name.toLowerCase()) }}</span>
          </label>
        </div>
      </div>
    </div>

    <div class="flex items-center justify-center space-x-2">
      <button
        type="button"
        class="inline-flex items-center px-2.5 py-1 border border-transparent text-xs font-medium rounded shadow-sm text-white bg-gray-500 hover:bg-gray-700 focus:outline-none"
        @click="turnOnAllEventTypes()"
      >
        Select all
      </button>
      <button
        type="button"
        class="inline-flex items-center px-2.5 py-1 border border-transparent text-xs font-medium rounded shadow-sm text-white bg-gray-500 hover:bg-gray-700 focus:outline-none"
        @click="turnOffAllEventTypes()"
      >
        Unselect all
      </button>
    </div>

    <div>
      <p class="text-center text-xs text-gray-500">
        Data for 2020 events were extracted from spreadsheet provided by Discord user @lamCube.
        Start timestamps and durations are provided on a best-effort basis and not accurate. 2019 or
        earlier events are omitted due to event scarcity and data incompleteness.
      </p>
      <p class="text-center text-xs text-gray-700">
        Tip: Hover over / click on event labels to reveal details.
      </p>
    </div>
  </div>

  <div class="overflow-x-auto pb-6 mt-4">
    <div
      class="Calendar grid content-evenly gap-6 mx-auto"
      :class="[
        forceFullWidth ? 'Calendar--full-width' : null,
        forceSingleColumn ? 'Calendar--single-column' : null,
      ]"
    >
      <template v-for="[month, date2events] in months" :key="month">
        <calendar-month
          :monthStr="month"
          :date2events="date2events"
          :eventTypesOn="eventTypesOn"
          :forceFullWidth="forceFullWidth"
        />
      </template>
    </div>
  </div>
</template>

<script>
import CalendarMonth from "@/components/CalendarMonth.vue";
import EventBadge from "@/components/EventBadge.vue";
import Info from "@/components/Info.vue";

import { computed, ref, watch } from "vue";
import dayjs from "dayjs";
import utc from "dayjs/plugin/utc";

import { events, eventTypes } from "@/lib";
import { getLocalStorage, setLocalStorage } from "@/utils";

dayjs.extend(utc);

const USE_UTC_DATES_LOCALSTORAGE_KEY = "useUtcDates";
const FORCE_FULL_WIDTH_LOCALSTORAGE_KEY = "forceFullWidth";
const FORCE_SINGLE_COLUMN_LOCALSTORAGE_KEY = "forceSingleColumn";

const eventTypeOnLocalStorageKey = eventType => `show-${eventType}`;

const getEventTypesOn = () => {
  const eventTypesOn = {};
  for (const ev of eventTypes) {
    eventTypesOn[ev[0]] = getLocalStorage(eventTypeOnLocalStorageKey(ev[0])) !== "false";
  }
  return eventTypesOn;
};

export default {
  components: {
    CalendarMonth,
    EventBadge,
    Info,
  },

  setup() {
    const useUtcDates = ref(getLocalStorage(USE_UTC_DATES_LOCALSTORAGE_KEY) !== "false");
    watch(useUtcDates, () => setLocalStorage(USE_UTC_DATES_LOCALSTORAGE_KEY, useUtcDates.value));
    const forceFullWidth = ref(getLocalStorage(FORCE_FULL_WIDTH_LOCALSTORAGE_KEY) === "true");
    watch(forceFullWidth, () => setLocalStorage(FORCE_FULL_WIDTH_LOCALSTORAGE_KEY, forceFullWidth));
    const forceSingleColumn = ref(getLocalStorage(FORCE_SINGLE_COLUMN_LOCALSTORAGE_KEY) === "true");
    watch(forceSingleColumn, () =>
      setLocalStorage(FORCE_SINGLE_COLUMN_LOCALSTORAGE_KEY, forceSingleColumn.value)
    );

    const eventTypesOn = ref(getEventTypesOn());
    const persistEventTypeOn = (type, on) => setLocalStorage(eventTypeOnLocalStorageKey(type), on);
    const turnOnAllEventTypes = () => {
      for (const ev of eventTypes) {
        eventTypesOn.value[ev[0]] = true;
        persistEventTypeOn(ev[0], true);
      }
    };
    const turnOffAllEventTypes = () => {
      for (const ev of eventTypes) {
        eventTypesOn.value[ev[0]] = false;
        persistEventTypeOn(ev[0], false);
      }
    };

    const months = computed(() => {
      const months = [];
      let currentMonth;
      let date2events;
      for (const event of events) {
        let startTime = dayjs(event.startTimestamp * 1000);
        if (useUtcDates.value) {
          startTime = startTime.utc();
        }
        const month = startTime.format("YYYY-MM");
        const date = startTime.date();
        if (month !== currentMonth) {
          if (currentMonth) {
            months.push([currentMonth, date2events]);
          }
          currentMonth = month;
          date2events = {};
        }
        if (!(date in date2events)) {
          date2events[date] = [];
        }
        date2events[date].push({
          ...event,
          startTime,
          durationSeconds: (event.endTimestamp || event.startTimestamp) - event.startTimestamp,
        });
      }
      if (events.length > 0) {
        months.push([currentMonth, date2events]);
      }
      return months.reverse();
    });

    return {
      useUtcDates,
      forceFullWidth,
      forceSingleColumn,
      eventTypes,
      eventTypesOn,
      persistEventTypeOn,
      turnOnAllEventTypes,
      turnOffAllEventTypes,
      months,
      capitalize: s => s.charAt(0).toUpperCase() + s.slice(1),
    };
  },
};
</script>

<style scoped>
.Calendar.Calendar--full-width {
  /* 744px * 3 + 24px * 2 */
  max-width: 2280px;
  grid-template-columns: repeat(auto-fit, minmax(744px, 1fr));
}

.Calendar.Calendar--single-column,
.Calendar.Calendar--single-column.Calendar.Calendar--full-width {
  grid-template-columns: 1fr;
}

@media (min-width: 768px) {
  .Calendar {
    max-width: 2280px;
    grid-template-columns: repeat(auto-fit, minmax(744px, 1fr));
  }
}
</style>
