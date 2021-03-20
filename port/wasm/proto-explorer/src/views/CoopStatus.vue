<template>
  <api-requester
    apiEndpoint="/ei/coop_status"
    requestMessage="ContractCoopStatusRequest"
    responseMessage="ContractCoopStatusResponse"
    :persistFormData="persistFormData"
    :getRequestPayloadObject="getRequestPayloadObject"
  >
    <template #form-body>
      <parameter-input
        name="contract_id"
        label="Contract ID"
        placeholder="Ex: graviton-project"
        :required="true"
        v-model.trim="contractId"
      />
      <parameter-input
        name="coop_code"
        label="Coop code"
        :required="true"
        v-model.trim="coopCode"
      />
      <parameter-input
        name="user_id"
        label="User ID"
        placeholder="Optional, random by default"
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
import { CLIENT_VERSION, basicRequestInfo, randomEiUserId } from "@/lib/lib";
import { getLocalStorage, setLocalStorage } from "@/utils";

const CONTRACT_ID_LOCALSTORAGE_KEY = "contract_id";
const COOP_CODE_LOCALSTORAGE_KEY = "coop_code";
const USER_ID_LOCALSTORAGE_KEY = "user_id";

export default {
  components: {
    "api-requester": APIRequester,
    ParameterInput,
    RequestButton,
  },

  setup() {
    const contractId = ref(getLocalStorage(CONTRACT_ID_LOCALSTORAGE_KEY) || "");
    const coopCode = ref(getLocalStorage(COOP_CODE_LOCALSTORAGE_KEY) || "");
    const userId = ref(getLocalStorage(USER_ID_LOCALSTORAGE_KEY) || "");
    const formValid = computed(() => contractId.value !== "" && coopCode.value != "");

    const persistFormData = () => {
      setLocalStorage(CONTRACT_ID_LOCALSTORAGE_KEY, contractId.value);
      setLocalStorage(COOP_CODE_LOCALSTORAGE_KEY, coopCode.value);
      setLocalStorage(USER_ID_LOCALSTORAGE_KEY, userId.value);
    };

    const getRequestPayloadObject = () => {
      const uid = userId.value || randomEiUserId();
      return {
        contractIdentifier: contractId.value,
        coopIdentifier: coopCode.value,
        userId: uid,
        currentClientVersion: CLIENT_VERSION,
        rinfo: basicRequestInfo(uid),
      };
    };

    return {
      contractId,
      coopCode,
      userId,
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
