import { onBeforeUnmount, ref, Ref } from "vue";
import { APIClient, encodeMessage } from "@/lib/lib";
import { getLocalStorage, setLocalStorage } from "@/utils";

const ENCODED_REQUEST_PAYLOAD_LOCALSTORAGE_KEY = "encoded_request_payload";
const ENCODED_RESPONSE_PAYLOAD_LOCALSTORAGE_KEY = "encoded_response_payload";

/**
 * @callback payloadObjectGetter
 * @returns {Object}
 */

/**
 * @typedef {Object} persistedValues
 * @property {String} encodedRequestPayload
 * @property {String} encodedResponsePayload
 */

/**
 * @callback persistedValuesGetter
 * @returns {persistedValues}
 */

/**
 * @callback persistedValuesSetter
 * @param {persistedValues} values
 * @returns {void}
 */

/**
 * @typedef {Object} useAPIClientReturnValue
 * @property {!Ref<String>} encodedRequestPayload
 * @property {!Ref<String>} encodedResponsePayload
 * @property {!Ref<String>} requestError
 * @property {!Ref<Boolean>} loading
 * @property {!Function} sendRequest
 */

/**
 * @param {!String} endpoint
 * @param {!String} requestMessageName
 * @param {!payloadObjectGetter} getRequestPayloadObject
 * @param {persistedValuesGetter} [getPersistedValues]
 * @param {persistedValuesSetter} [persistValues]
 * @returns {useAPIClientReturnValue}
 */
function useAPIClient(
  endpoint,
  requestMessageName,
  getRequestPayloadObject,
  getPersistedValues,
  persistValues
) {
  /** @type {persistedValues} */
  const persistedValues = typeof getPersistedValues === "function" ? getPersistedValues() : {};
  const {
    encodedRequestPayload: defaultEncodedRequestPayload = "",
    encodedResponsePayload: defaultEncodedResponsePayload = "",
  } = persistedValues;

  let encodedRequestPayload = ref(defaultEncodedRequestPayload);
  let encodedResponsePayload = ref(defaultEncodedResponsePayload);
  let requestError = ref("");
  let loading = ref(false);

  const client = new APIClient();
  // requestId is used to track whether the request being handled has been
  // superseded by another request.
  let requestId;

  const sendRequest = async () => {
    const thisRequestId = Math.floor(Math.random() * 65536);
    requestId = thisRequestId;
    loading.value = true;
    encodedResponsePayload.value = null;
    requestError.value = null;
    const result = await request(
      client,
      endpoint,
      requestMessageName,
      getRequestPayloadObject,
      encodedRequestPayload
    );
    if (requestId === thisRequestId) {
      encodedResponsePayload.value = result.payload;
      requestError.value = result.error;
      if (typeof persistValues === "function") {
        if (result.error === null) {
          persistValues({
            encodedRequestPayload: encodedRequestPayload.value,
            encodedResponsePayload: encodedResponsePayload.value,
          });
        } else {
          persistValues({
            encodedRequestPayload: "",
            encodedResponsePayload: "",
          });
        }
      }
    }
    loading.value = false;
  };

  onBeforeUnmount(() => {
    client.destroy();
  });

  return {
    encodedRequestPayload,
    encodedResponsePayload,
    requestError,
    loading,
    sendRequest,
  };
}

/**
 * @typedef {Object} requestResult
 * @property {String} payload - Encoded response payload.
 * @property {String} error
 */

/**
 * @async
 * @param {!APIClient} client
 * @param {!String} endpoint
 * @param {!String} messageName
 * @param {!payloadObjectGetter} getRequestPayloadObject
 * @param {!Ref<String>} encodedRequestPayload
 * @returns {Promise<!requestResult>}
 */
async function request(
  client,
  endpoint,
  messageName,
  getRequestPayloadObject,
  encodedRequestPayload
) {
  try {
    encodedRequestPayload.value = encodeMessage(messageName, getRequestPayloadObject());
    const encodedResponsePayload = await client.request(
      endpoint,
      encodedRequestPayload.value,
      true /* flush */
    );
    return {
      payload: encodedResponsePayload,
      error: null,
    };
  } catch (e) {
    return {
      payload: null,
      error: `${e}`,
    };
  }
}

/**
 * @returns {persistedValues}
 */
function getPayloadsFromLocalStorage() {
  return {
    encodedRequestPayload: getLocalStorage(ENCODED_REQUEST_PAYLOAD_LOCALSTORAGE_KEY) || "",
    encodedResponsePayload: getLocalStorage(ENCODED_RESPONSE_PAYLOAD_LOCALSTORAGE_KEY) || "",
  };
}

/**
 * @param {persistedValues} values
 * @returns {void}
 */
function persistPayloadsToLocalStorage(values) {
  setLocalStorage(ENCODED_REQUEST_PAYLOAD_LOCALSTORAGE_KEY, values.encodedRequestPayload);
  setLocalStorage(ENCODED_RESPONSE_PAYLOAD_LOCALSTORAGE_KEY, values.encodedResponsePayload);
}

export { useAPIClient, getPayloadsFromLocalStorage, persistPayloadsToLocalStorage };
