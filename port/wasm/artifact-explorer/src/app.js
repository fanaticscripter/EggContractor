import { createApp } from "vue";

import VueTippy, { Tippy } from "vue-tippy";
import "tippy.js/dist/tippy.css";
import "tippy.js/themes/translucent.css";

import App from "./App.vue";
import router from "./router";

const app = createApp(App);
app.use(router);
app.use(VueTippy, {
  defaultProps: {
    theme: "translucent",
    onShow: instance => {
      // Image lazyloading.
      for (const img of instance.popper.querySelectorAll("img")) {
        if (!img.src && img.dataset.src) {
          img.src = img.dataset.src;
        }
      }
    },
  },
});
app.component("tippy", Tippy);
app.mount("#app");
