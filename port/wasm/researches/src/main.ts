import { createApp } from "vue";

import App from "./App.vue";
import "./index.css";

import { initDatabase } from "./data";

(async () => {
  await initDatabase();
  createApp(App).mount("#app");
})();
