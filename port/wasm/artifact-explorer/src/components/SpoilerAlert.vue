<template>
  <div v-if="show" class="rounded-md bg-blue-50 p-4 my-4">
    <div class="flex">
      <div class="flex-shrink-0">
        <svg
          class="h-5 w-5 text-blue-400"
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 20 20"
          fill="currentColor"
          aria-hidden="true"
        >
          <path
            fill-rule="evenodd"
            d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z"
            clip-rule="evenodd"
          />
        </svg>
      </div>
      <div class="ml-3">
        <h3 class="text-sm font-medium text-blue-800">SPOILER ALERT</h3>
        <div class="mt-2 text-sm text-blue-700">
          <p>
            This site contains a shocking amount of spoilers. By further interacting with this site,
            you agree to waive your natural and/or legal right to a spoiler-free gaming experience
            and the excitement or despair of discovery.
          </p>
        </div>
        <div class="mt-4">
          <div class="-mx-2 -my-1.5 flex">
            <button
              class="bg-blue-50 px-2 py-1.5 rounded-md text-sm font-medium text-blue-800 hover:bg-blue-100 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-blue-50 focus:ring-blue-600"
              @click="giveConsent()"
            >
              Get out of the way already!
            </button>
            <button
              class="ml-3 bg-blue-50 px-2 py-1.5 rounded-md text-sm font-medium text-blue-800 hover:bg-blue-100 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-blue-50 focus:ring-blue-600"
              @click="unmountApp()"
            >
              Nope
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { render } from "vue";

import { getLocalStorage, setLocalStorage } from "@/utils";

const SPOILER_ALERT_CONSENT_LOCALSTORAGE_KEY = "spoilerAlertConsent";

export default {
  data() {
    return {
      show: getLocalStorage(SPOILER_ALERT_CONSENT_LOCALSTORAGE_KEY) !== "true",
    };
  },

  methods: {
    giveConsent() {
      setLocalStorage(SPOILER_ALERT_CONSENT_LOCALSTORAGE_KEY, true);
      this.show = false;
    },

    unmountApp() {
      const app = document.getElementById("app");
      render(null, app);
      app.innerHTML = `<div class="text-center">Enjoy your spoiler-free experience.</div>`;
    },
  },
};
</script>
