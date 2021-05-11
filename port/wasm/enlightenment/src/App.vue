<template>
  <div class="max-w-5xl w-full px-4 pb-4 xl:px-0 mx-auto">
    <the-player-id-form :playerIdPreload="playerIdPreload" :submit="submitPlayerId" />

    <!-- Use a key to recreate on data loading -->
    <base-error-boundary v-if="playerId" :key="`${playerId}:${refreshId}`">
      <Suspense>
        <template #default>
          <the-companion :playerId="playerId" />
        </template>
        <template #fallback>
          <base-loading />
        </template>
      </Suspense>
    </base-error-boundary>

    <template v-else>
      <div class="text-sm mt-4">
        This tool pulls the latest save for the specified player, and generates a report with
        information useful for completing the enlightenment diamond trophy, e.g. hab space, internal
        hatchery rate, equipped artifacts, projected completion date, etc.
      </div>

      <div class="text-sm mt-2">
        If you are not yet sure about feasibility of the trophy, or earning enough cash for the
        required WD level appears to be very difficult,
        <a
          href="https://docs.google.com/spreadsheets/d/157K4r3Z5wfCNKhUWb34mlxM08DEA1AWamsA20xjQIhw/edit?usp=sharing"
          target="_blank"
          class="text-blue-500 hover:text-blue-600"
          >Sami#2336's detailed spreadsheet</a
        >
        may help with your decision and execution.
      </div>
    </template>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from "vue";

import { getLocalStorage, setLocalStorage } from "@/utils";
import BaseErrorBoundary from "@/components/BaseErrorBoundary.vue";
import BaseLoading from "@/components/BaseLoading.vue";
import ThePlayerIdForm from "@/components/ThePlayerIdForm.vue";
import TheCompanion from "./components/TheCompanion.vue";

const PLAYER_ID_LOCALSTORAGE_KEY = "playerId";

export default defineComponent({
  components: {
    BaseErrorBoundary,
    BaseLoading,
    ThePlayerIdForm,
    TheCompanion,
  },
  setup() {
    const playerIdPreload =
      new URLSearchParams(window.location.search).get("playerId") ||
      getLocalStorage(PLAYER_ID_LOCALSTORAGE_KEY) ||
      "";
    const playerId = ref("");
    const refreshId = ref(Date.now());
    const submitPlayerId = (id: string) => {
      playerId.value = id;
      refreshId.value = Date.now();
      setLocalStorage(PLAYER_ID_LOCALSTORAGE_KEY, id);
    };
    return {
      playerIdPreload,
      playerId,
      refreshId,
      submitPlayerId,
    };
  },
});
</script>
