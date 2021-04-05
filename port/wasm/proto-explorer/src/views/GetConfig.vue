<template>
  <api-requester
    apiEndpoint="/ei/get_config"
    requestMessage="ConfigRequest"
    responseMessage="ConfigResponse"
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
      <request-button :formValid="formValid" />
    </template>
  </api-requester>
</template>

<script>
import APIRequester from "@/components/APIRequester.vue";
import ParameterInput from "@/components/ParameterInput.vue";
import RequestButton from "@/components/RequestButton.vue";

import { computed, ref } from "vue";
import { basicRequestInfo } from "@/lib/lib";
import { getLocalStorage, setLocalStorage } from "@/utils";

const USER_ID_LOCALSTORAGE_KEY = "user_id";

export default {
  components: {
    "api-requester": APIRequester,
    ParameterInput,
    RequestButton,
  },

  setup() {
    const userId = ref(getLocalStorage(USER_ID_LOCALSTORAGE_KEY) || "");
    const formValid = computed(() => userId.value !== "");

    const persistFormData = () => {
      setLocalStorage(USER_ID_LOCALSTORAGE_KEY, userId.value);
    };

    const getRequestPayloadObject = () => ({
      rinfo: basicRequestInfo(userId.value),
    });

    return {
      userId,
      formValid,
      persistFormData,
      getRequestPayloadObject,
    };
  },
};
</script>
