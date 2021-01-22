import { createApp } from "vue";

import VueTippy from "vue-tippy";
import "tippy.js/dist/tippy.css";
import "tippy.js/themes/translucent.css";

import App from "./App.vue";

(async function () {
  const wasm = __WASM_FILE__;
  let mod;
  if ("compileStreaming" in WebAssembly) {
    mod = await WebAssembly.compileStreaming(fetch(wasm));
  } else {
    // Safari doesn't support WebAssembly.compileStreaming
    const resp = await fetch(wasm);
    mod = await WebAssembly.compile(await resp.arrayBuffer());
  }

  const go = new Go();

  // Ideally the wasm call should go into a web worker, but I'm lazy.
  async function retrieveMissions(playerId) {
    const instance = await WebAssembly.instantiate(mod, go.importObject);
    let payload;
    await new Promise((resolve, reject) => {
      self.wasmArgs = [playerId];
      self.wasmCallback = res => {
        res = JSON.parse(res);
        if (res.successful) {
          payload = res.data;
          resolve();
        } else {
          reject(res.error);
        }
      };
      go.run(instance);
    });
    return payload;
  }

  const app = createApp(App, { retrieveMissions });
  app.use(VueTippy, {
    defaultProps: { theme: "translucent" },
  });
  app.mount("#app");
})();
