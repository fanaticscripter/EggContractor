import { createApp } from "vue";

import VueTippy from "vue-tippy";
import "tippy.js/dist/tippy.css";
import "tippy.js/themes/translucent.css";

import App from "./App.vue";
import router from "./router";

window.onload = () => {
  const app = createApp(App);
  app.use(router);
  app.use(VueTippy, {
    defaultProps: { theme: "translucent" },
  });
  app.mount("#app");
};
