<template>
  <div class="flex flex-col flex-1 pt-1 pb-4">
    <div class="flex flex-col flex-1 max-w-4xl w-full mx-auto px-4 xl:px-0 space-y-2 mt-2">
      <div class="max-w-sm">
        <label for="message" class="block text-sm font-medium text-gray-700">
          Protobuf message type
        </label>
        <select
          id="message"
          name="message"
          class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none sm:text-sm rounded-md"
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
        class="px-3 py-2 w-full resize-y border border-gray-300 rounded-md text-base sm:text-xs font-mono"
        placeholder="Paste base64-encoded payload here..."
        spellcheck="false"
        v-model.trim="encodedPayload"
      ></textarea>

      <decode-result
        :message="message"
        :authenticated="authenticated"
        :encodedPayload="encodedPayload"
      />
    </div>
  </div>
</template>

<script>
import DecodeResult from "@/components/DecodeResult.vue";

import { ref, watch } from "vue";
import { messageGroups } from "@/lib/lib";
import { getLocalStorage, setLocalStorage } from "@/utils";

const MESSAGE_LOCALSTORAGE_KEY = "message";
const AUTHENTICATED_LOCALSTORAGE_KEY = "authenticated";
const ENCODED_PAYLOAD_LOCALSTORAGE_KEY = "encoded_payload";
const DEFAULT_MESSAGE = "EggIncFirstContactResponse";

export default {
  components: {
    DecodeResult,
  },

  setup() {
    const message = ref(getLocalStorage(MESSAGE_LOCALSTORAGE_KEY) || DEFAULT_MESSAGE);
    const authenticated = ref(getLocalStorage(AUTHENTICATED_LOCALSTORAGE_KEY) === "true");
    const encodedPayload = ref(getLocalStorage(ENCODED_PAYLOAD_LOCALSTORAGE_KEY) || "");

    watch(message, () => setLocalStorage(MESSAGE_LOCALSTORAGE_KEY, message.value));
    watch(authenticated, () =>
      setLocalStorage(AUTHENTICATED_LOCALSTORAGE_KEY, authenticated.value)
    );
    watch(encodedPayload, () =>
      setLocalStorage(ENCODED_PAYLOAD_LOCALSTORAGE_KEY, encodedPayload.value)
    );

    return {
      messageGroups,
      message,
      authenticated,
      encodedPayload,
    };
  },
};
</script>

<style scoped>
textarea#input {
  min-height: 6rem;
}
</style>
