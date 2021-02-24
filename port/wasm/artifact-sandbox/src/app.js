import { createApp } from "vue";

import VueTippy from "vue-tippy";
import "tippy.js/dist/tippy.css";
import "tippy.js/themes/translucent.css";

import App from "./App.vue";
import router from "./router";

import { iconURL } from "./utils";

// Expose proto to global scope.
require("./lib/schema_pb");

const app = createApp(App);
app.use(router);
app.use(VueTippy, {
  defaultProps: { theme: "translucent" },
});
app.mixin({
  methods: {
    iconURL,
  },
});
app.mount("#app");
