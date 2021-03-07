<template>
  <div class="flex flex-col flex-1 pt-1 pb-4">
    <div class="flex flex-col flex-1 max-w-4xl w-full mx-auto px-4 xl:px-0 space-y-2 mt-2">
      <form
        class="space-y-2"
        @submit.prevent="
          persistFormData();
          sendRequest();
        "
      >
        <slot name="form-body"></slot>
      </form>

      <div v-if="encodedRequestPayload">
        <label for="message" class="flex items-center text-sm font-medium text-gray-700">
          Request payload
          <copy-button class="ml-1" :content="encodedRequestPayload" />
        </label>
        <textarea
          id="request"
          class="mt-1 px-3 py-2 w-full resize-y border border-gray-300 rounded-md text-base sm:text-xs font-mono"
          rows="1"
          spellcheck="false"
          v-model="encodedRequestPayload"
          readonly
        ></textarea>
      </div>

      <loading-spinner v-if="loading">Waiting for server response...</loading-spinner>

      <div v-if="!loading && encodedResponsePayload">
        <label for="message" class="flex items-center text-sm font-medium text-gray-700">
          Response payload
          <copy-button class="ml-1" :content="encodedResponsePayload" />
        </label>
        <textarea
          id="response"
          class="mt-1 px-3 py-2 w-full resize-y border border-gray-300 rounded-md text-base sm:text-xs font-mono"
          spellcheck="false"
          v-model="encodedResponsePayload"
          readonly
        ></textarea>
      </div>

      <div
        v-if="requestError !== null"
        class="text-xs text-red-500 font-mono font-medium whitespace-pre-wrap break-words"
      >
        {{ requestError }}
      </div>

      <decode-result
        v-show="!loading && encodedResponsePayload"
        :message="responseMessage"
        :authenticated="true"
        :encodedPayload="encodedResponsePayload"
      />
    </div>
  </div>
</template>

<script>
import CopyButton from "@/components/CopyButton.vue";
import DecodeResult from "@/components/DecodeResult.vue";
import LoadingSpinner from "./LoadingSpinner.vue";
import RequestButton from "@/components/RequestButton.vue";

import {
  useAPIClient,
  getPayloadsFromLocalStorage,
  persistPayloadsToLocalStorage,
} from "@/composables/api";

export default {
  components: {
    CopyButton,
    DecodeResult,
    LoadingSpinner,
    RequestButton,
  },

  props: {
    apiEndpoint: {
      type: String,
      required: true,
    },
    requestMessage: {
      type: String,
      required: true,
    },
    responseMessage: {
      type: String,
      required: true,
    },
    persistFormData: {
      type: Function,
      required: true,
    },
    getRequestPayloadObject: {
      type: Function,
      required: true,
    },
  },

  setup({
    apiEndpoint,
    requestMessage,
    responseMessage,
    persistFormData,
    getRequestPayloadObject,
  }) {
    const {
      encodedRequestPayload,
      encodedResponsePayload,
      requestError,
      loading,
      sendRequest,
    } = useAPIClient(
      apiEndpoint,
      requestMessage,
      getRequestPayloadObject,
      getPayloadsFromLocalStorage,
      persistPayloadsToLocalStorage
    );

    return {
      responseMessage,
      encodedRequestPayload,
      encodedResponsePayload,
      requestError,
      loading,
      persistFormData,
      sendRequest,
    };
  },
};
</script>

<style scoped>
textarea#request {
  min-height: 2rem;
}

textarea#response {
  min-height: 6rem;
}
</style>
