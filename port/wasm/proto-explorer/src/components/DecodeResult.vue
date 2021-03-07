<template>
  <div class="flex flex-col flex-1 space-y-2">
    <div class="text-sm font-medium text-gray-700 break-words">
      <div
        v-if="decodeError !== null"
        class="text-xs text-red-500 font-mono font-medium whitespace-pre-wrap break-words"
      >
        {{ decodeError }}
      </div>

      <div v-if="decodedMAC !== null">
        Message authentication code:
        <span class="text-xs font-mono">{{ decodedMAC }}</span>
      </div>

      <div v-if="message">
        <router-link
          :to="{ name: 'doc', hash: `#ei.${message}` }"
          target="_blank"
          class="hover:text-gray-500 border-b border-gray-500 border-dashed"
        >
          <code class="text-xs font-mono">{{ message }}</code> documentation
        </router-link>
      </div>
    </div>

    <div
      v-show="decodedPayload !== null"
      class="flex flex-col flex-1 relative"
      :style="{ minHeight: '24rem' }"
    >
      <div id="editor" class="absolute h-full w-full border border-gray-300 rounded-md"></div>
      <copy-button
        class="absolute top-1 left-1 z-50"
        :content="formattedDecodedPayload"
        tooltip="Copy decoded payload as JSON"
      />
    </div>

    <div
      v-show="decodedPayload !== null"
      class="text-sm text-gray-700 flex flex-row flex-wrap items-center"
    >
      <label for="ei-raw-value" class="mr-2">Format EI value:</label>
      <div class="flex flex-row flex-1 min-w-full sm:min-w-max items-center">
        <input
          class="flex-1 max-w-xs shadow-sm px-2 py-1 text-base sm:text-sm border-gray-300 rounded appearance-none mr-1"
          type="number"
          placeholder="Ex. 10000000000000 (=10T)"
          id="ei-raw-value"
          v-model.number="eiValue"
        />
        <template v-if="formattedEIValue !== ''"> = {{ formattedEIValue }} </template>
      </div>
    </div>
  </div>
</template>

<script>
import CopyButton from "./CopyButton.vue";

import { computed, onBeforeUnmount, onMounted, ref, toRefs, watch } from "vue";
import { decodeMessage, formatEIValue } from "@/lib/lib";

export default {
  components: {
    CopyButton,
  },

  props: {
    message: String,
    authenticated: Boolean,
    encodedPayload: String,
  },

  setup(props) {
    const { message, authenticated, encodedPayload } = toRefs(props);
    const eiValue = ref("");

    const decodeResult = computed(() => {
      if (!message.value || !encodedPayload.value) {
        return {};
      }
      const result = decodeMessage(message.value, encodedPayload.value, authenticated.value);
      // If decoding failed, see if we can decode as authenticated instead.
      if (result.error !== undefined && !authenticated.value) {
        const resultAsAuthenticated = decodeMessage(message.value, encodedPayload.value, true);
        if (resultAsAuthenticated.error === undefined) {
          resultAsAuthenticated.error =
            `Failed to decode directlly, but successfully decoded as authenticated message. Forgot to check the box?\n` +
            `(${result.error})`;
          return resultAsAuthenticated;
        }
      }
      return result;
    });

    const decodedPayload = computed(() => {
      const { payload = null } = decodeResult.value;
      return payload;
    });

    const decodedMAC = computed(() => {
      const { code = null } = decodeResult.value;
      return code;
    });

    const decodeError = computed(() => {
      const { error = null } = decodeResult.value;
      return error;
    });

    const formattedDecodedPayload = computed(() => {
      const decoded = decodedPayload.value;
      return decoded === null ? "" : JSON.stringify(decoded, null, 2);
    });

    const formattedEIValue = computed(() => {
      const val = eiValue.value;
      return val === "" || isNaN(val) ? "" : formatEIValue(val);
    });

    let editor;

    onMounted(() => {
      editor = ace.edit("editor");
      editor.setReadOnly(true);
      editor.setOption("tabSize", 2);
      editor.session.setMode("ace/mode/json");
      editor.session.setUseWrapMode(true);
      editor.session.setValue(formattedDecodedPayload.value);
    });

    watch(formattedDecodedPayload, () => editor.session.setValue(formattedDecodedPayload.value));

    onBeforeUnmount(() => {
      editor.destroy();
    });

    return {
      eiValue,
      decodedPayload,
      decodedMAC,
      decodeError,
      formattedDecodedPayload,
      formattedEIValue,
    };
  },
};
</script>

<style scoped>
/* Disable spin button for number input */
input[type="number"]::-webkit-outer-spin-button,
input[type="number"]::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}

input[type="number"] {
  -moz-appearance: textfield;
}
</style>
