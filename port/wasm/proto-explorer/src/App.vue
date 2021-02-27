<template>
  <!-- TODO: show app version -->
  <div class="flex flex-col flex-1 max-w-4xl w-full mx-auto px-4 xl:px-0 space-y-2">
    <div>
      <p class="text-sm font-medium text-gray-700">
        App version:
        <code class="text-xs font-mono">{{ appVersion }}</code>
        <span class="font-normal">
          (<a
            href="https://github.com/fanaticscripter/EggContractor/tree/master/misc/protobuf"
            target="_blank"
            class="hover:text-gray-500 border-b border-gray-500 border-dashed"
            >protobuf definitions</a
          >)
        </span>
      </p>
    </div>

    <div class="max-w-sm">
      <label for="message" class="block text-sm font-medium text-gray-700">
        Protobuf message type
      </label>
      <select
        id="message"
        name="message"
        class="mt-1 block w-full pl-3 pr-10 py-2 text-base bg-gray-50 border-gray-300 focus:outline-none sm:text-sm rounded-md"
        v-model="message"
      >
        <optgroup v-for="group in messageGroups" :key="group.label" :label="group.label">
          <option v-for="message in group.messages" :key="message">{{ message }}</option>
        </optgroup>
      </select>
    </div>

    <div class="flex items-center">
      <input
        id="authenticated"
        name="authenticated"
        type="checkbox"
        class="h-4 w-4 text-blue-600 focus:outline-none border-gray-300 rounded"
        v-model="authenticated"
      />
      <label for="authenticated" class="ml-2 block text-sm text-gray-900">
        Decode as authenticated message
      </label>
    </div>

    <textarea
      id="input"
      class="px-3 py-2 w-full resize-y bg-gray-50 border border-gray-300 rounded-md text-base sm:text-xs font-mono"
      placeholder="Paste base64-encoded payload here..."
      spellcheck="false"
      v-model="encodedPayload"
    ></textarea>

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
        <a
          :href="`doc.html#ei.${message}`"
          target="_blank"
          class="hover:text-gray-500 border-b border-gray-500 border-dashed"
        >
          <code class="text-xs font-mono">{{ message }}</code> documentation
        </a>
      </div>
    </div>

    <div v-show="decodedPayload !== null" class="flex flex-col flex-1 relative">
      <div id="editor" class="flex-1 border border-gray-300 rounded-md"></div>
      <button
        class="absolute top-0 right-0 px-2 py-2 z-50"
        @click="copyDecodedPayload()"
        v-tippy="{ content: 'Copy decoded payload as JSON' }"
      >
        <svg
          class="h-5 w-5 text-gray-500 hover:text-gray-700"
          viewBox="0 0 20 20"
          fill="currentColor"
        >
          <path d="M8 2a1 1 0 000 2h2a1 1 0 100-2H8z" />
          <path
            d="M3 5a2 2 0 012-2 3 3 0 003 3h2a3 3 0 003-3 2 2 0 012 2v6h-4.586l1.293-1.293a1 1 0 00-1.414-1.414l-3 3a1 1 0 000 1.414l3 3a1 1 0 001.414-1.414L10.414 13H15v3a2 2 0 01-2 2H5a2 2 0 01-2-2V5zM15 11h2a1 1 0 110 2h-2v-2z"
          />
        </svg>
      </button>
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
          v-model="eiValue"
        />
        <template v-if="formattedEIValue !== ''"> = {{ formattedEIValue }} </template>
      </div>
    </div>
  </div>
</template>

<script>
import copyTextToClipboard from "copy-text-to-clipboard";
import { APP_VERSION, decodeMessage, messageGroups, formatEIValue } from "@/lib/lib";
import { getLocalStorage, setLocalStorage } from "./utils";

const MESSAGE_LOCALSTORAGE_KEY = "message";
const AUTHENTICATED_LOCALSTORAGE_KEY = "authenticated";
const ENCODED_PAYLOAD_LOCALSTORAGE_KEY = "encoded_payload";
const DEFAULT_MESSAGE = "EggIncFirstContactResponse";

export default {
  data() {
    return {
      appVersion: APP_VERSION,

      messageGroups,
      message: getLocalStorage(MESSAGE_LOCALSTORAGE_KEY) || DEFAULT_MESSAGE,
      authenticated: getLocalStorage(AUTHENTICATED_LOCALSTORAGE_KEY) === "true",
      encodedPayload: getLocalStorage(ENCODED_PAYLOAD_LOCALSTORAGE_KEY) || "",

      eiValue: "",
    };
  },

  editor: null,

  mounted() {
    const editor = ace.edit("editor");
    this.$options.editor = editor;
    editor.setReadOnly(true);
    editor.setOption("tabSize", 2);
    editor.session.setMode("ace/mode/json");
    editor.session.setUseWrapMode(true);
    editor.session.setValue(this.formattedDecodedPayload);
  },

  watch: {
    message() {
      setLocalStorage(MESSAGE_LOCALSTORAGE_KEY, this.message);
    },

    authenticated() {
      setLocalStorage(AUTHENTICATED_LOCALSTORAGE_KEY, this.authenticated);
    },

    encodedPayload() {
      setLocalStorage(ENCODED_PAYLOAD_LOCALSTORAGE_KEY, this.encodedPayload);
    },

    formattedDecodedPayload() {
      this.$options.editor.session.setValue(this.formattedDecodedPayload);
    },
  },

  computed: {
    decodeResult() {
      if (!this.message || !this.encodedPayload) {
        return {};
      }
      const result = decodeMessage(this.message, this.encodedPayload, this.authenticated);
      // If decoding failed, see if we can decode as authenticated instead.
      if (result.error !== undefined && !this.authenticated) {
        const resultAsAuthenticated = decodeMessage(this.message, this.encodedPayload, true);
        if (resultAsAuthenticated.error === undefined) {
          resultAsAuthenticated.error =
            `Failed to decode directlly, but successfully decoded as authenticated message. Forgot to check the box?\n` +
            `(${result.error})`;
          return resultAsAuthenticated;
        }
      }
      return result;
    },

    decodedPayload() {
      const { payload = null } = this.decodeResult;
      return payload;
    },

    decodedMAC() {
      const { code = null } = this.decodeResult;
      return code;
    },

    decodeError() {
      const { error = null } = this.decodeResult;
      return error;
    },

    formattedDecodedPayload() {
      const decoded = this.decodedPayload;
      return decoded === null ? "" : JSON.stringify(decoded, null, 2);
    },

    formattedEIValue() {
      const val = parseFloat(this.eiValue.trim());
      if (isNaN(val)) {
        return "";
      }
      return formatEIValue(val);
    },
  },

  methods: {
    copyDecodedPayload() {
      copyTextToClipboard(this.formattedDecodedPayload);
    },
  },
};
</script>

<style scoped>
textarea#input {
  min-height: 6rem;
}

#editor {
  min-height: 24rem;
}

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
