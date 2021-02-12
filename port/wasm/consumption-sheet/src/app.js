import { createApp } from "vue";

import VueTippy from "vue-tippy";
import "tippy.js/dist/tippy.css";
import "tippy.js/themes/translucent.css";

import App from "./App.vue";
import data from "./app-data.json";

import { iconURL } from "./utils";

const app = createApp(App, { data });
app.use(VueTippy, {
  defaultProps: { theme: "translucent" },
});
app.mixin({
  methods: {
    iconURL,
  },
});
app.mount("#app");

window.activateScreenshotMode = () => {
  for (const el of document.querySelectorAll(".hide-in-screenshot-mode")) {
    el.dataset.display = el.style.display;
    el.style.display = "none";
  }
};

window.deactivateScreenshotMode = () => {
  for (const el of document.querySelectorAll(".hide-in-screenshot-mode")) {
    el.style.display = el.dataset.display;
  }
};
