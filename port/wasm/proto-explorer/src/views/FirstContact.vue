<template>
  <api-requester
    apiEndpoint="/ei/first_contact"
    requestMessage="EggIncFirstContactRequest"
    responseMessage="EggIncFirstContactResponse"
    :persistFormData="persistFormData"
    :getRequestPayloadObject="getRequestPayloadObject"
  >
    <template #form-body>
      <parameter-input
        name="user_id"
        label="User ID"
        placeholder="Ex: EI1234567890123456"
        :required="true"
        v-model.trim="userId"
      />
      <parameter-input
        name="device_id"
        label="Device ID"
        placeholder="Optional; default: a random UUID v4"
        v-model.trim="deviceId"
      />
      <parameter-input
        name="game_services_id"
        label="Game services ID"
        placeholder="Optional"
        v-model.trim="gameServicesId"
      />
      <request-button :formValid="formValid" />
    </template>
  </api-requester>
</template>

<script>
import APIRequester from "@/components/APIRequester.vue";
import ParameterInput from "@/components/ParameterInput.vue";
import RequestButton from "@/components/RequestButton.vue";

import { computed, ref } from "vue";
import { CLIENT_VERSION, PLATFORM, basicRequestInfo } from "@/lib/lib";
import { getLocalStorage, setLocalStorage, uuid4 } from "@/utils";

const USER_ID_LOCALSTORAGE_KEY = "user_id";
const DEVICE_ID_LOCALSTORAGE_KEY = "device_id";
const GAME_SERVICES_ID_LOCALSTORAGE_KEY = "game_services_id";

export default {
  components: {
    "api-requester": APIRequester,
    ParameterInput,
    RequestButton,
  },

  setup() {
    const userId = ref(getLocalStorage(USER_ID_LOCALSTORAGE_KEY) || "");
    const deviceId = ref(getLocalStorage(DEVICE_ID_LOCALSTORAGE_KEY) || "");
    const gameServicesId = ref(getLocalStorage(GAME_SERVICES_ID_LOCALSTORAGE_KEY) || "");
    const formValid = computed(() => userId.value !== "");

    const persistFormData = () => {
      setLocalStorage(USER_ID_LOCALSTORAGE_KEY, userId.value);
      setLocalStorage(DEVICE_ID_LOCALSTORAGE_KEY, deviceId.value);
      setLocalStorage(GAME_SERVICES_ID_LOCALSTORAGE_KEY, gameServicesId.value);
    };

    const getRequestPayloadObject = () => ({
      clientVersion: CLIENT_VERSION,
      platform: PLATFORM,
      eiUserId: userId.value,
      deviceId: deviceId.value || uuid4(),
      gameServicesId: gameServicesId.value || null,
      rinfo: basicRequestInfo(userId.value),
    });

    return {
      userId,
      deviceId,
      gameServicesId,
      formValid,
      persistFormData,
      getRequestPayloadObject,
    };
  },
};
</script>

<style scoped>
::v-deep(textarea#request) {
  min-height: 3rem !important;
}
</style>
