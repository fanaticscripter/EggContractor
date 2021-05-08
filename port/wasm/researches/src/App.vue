<template>
  <a
    :href="researchedJsonUrl"
    download="researches.json"
    class="block text-xs text-center text-blue-500 hover:text-blue-600 pointer-cursor mb-2"
  >
    Complete JSON data download
  </a>

  <div class="flex-1 xl:flex xl:max-w-full xl:justify-center w-full mx-auto">
    <div
      class="xl:flex xl:justify-end xl:flex-1 border-t border-b border-gray-100 xl:border-b-0 xl:border-r xl:border-gray-100"
    >
      <div class="flex flex-1 w-full max-w-4xl mx-auto xl:mx-0 p-4">
        <div class="relative w-full" :style="{ minHeight: '40rem' }">
          <ace-editor
            :modelValue="researchesJson"
            lang="json"
            :readonly="true"
            :foldAtIndentation="4"
          />
        </div>
      </div>
    </div>

    <div class="xl:flex xl:flex-1 bg-gray-50">
      <div class="flex flex-col flex-1 w-full max-w-4xl mx-auto xl:mx-0 p-4 space-y-2">
        <p class="text-sm">
          You can query the full list of researches as a SQLite database with the following schema:
        </p>
        <div class="relative w-full" :style="{ height: '40rem' }">
          <ace-editor :modelValue="schema" lang="sql" :readonly="true" />
        </div>

        <p class="text-sm">Enter your query here:</p>

        <div class="relative w-full" :style="{ height: '6rem' }">
          <ace-editor
            v-model="query"
            lang="sql"
            :commands="[submitQueryAceCommand]"
            :event-bus="eventBus"
          />
        </div>

        <button
          type="button"
          class="inline-flex items-center justify-center px-2.5 py-1.5 border border-transparent text-xs font-medium rounded shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
          @click="submitQuery"
        >
          Query ({{ submitQueryHotkeyDisplay }})
        </button>

        <div v-if="queryError" class="text-sm text-red-500">{{ queryError.toString() }}</div>

        <div
          v-show="queryError === null"
          class="flex-1 relative w-full"
          :style="{ minHeight: '20rem' }"
        >
          <ace-editor
            :modelValue="queryResultsJson"
            lang="json"
            :readonly="true"
            :foldAtIndentation="4"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, onBeforeUnmount, onMounted, ref } from "vue";
import hotkeys from "hotkeys-js";
import mitt from "mitt";
import { ParamsObject } from "sql.js";

import { executeQuery, researches, schema } from "@/data";
import AceEditor from "@/components/AceEditor.vue";
import { isApplePlatform } from "./utils";

const sampleQuery = `-- All internal hatchery rate researches.
SELECT id, name, levels AS maxLevel, per_level AS perLevel, effect_type = 'multiplicative' AS multiplicative, prices FROM research WHERE categories LIKE '%internal_hatchery_rate%' ORDER BY serial_id;
-- Another more advanced example where we also calculate the total price of each entry:
--SELECT research.id, research.name, research.levels AS maxLevel, research.per_level AS perLevel, research.effect_type = 'multiplicative' AS multiplicative, sum(CAST(prices.value AS REAL)) AS totalPrice FROM research, json_each(prices) AS prices WHERE categories LIKE '%internal_hatchery_rate%' GROUP BY research.serial_id ORDER BY research.serial_id;`;

export default defineComponent({
  name: "App",
  components: {
    AceEditor,
  },
  setup() {
    const researchesJson = JSON.stringify(researches, null, 2);
    const researchedJsonUrl = URL.createObjectURL(
      new Blob([researchesJson], { type: "application/json" })
    );
    onBeforeUnmount(() => {
      URL.revokeObjectURL(researchedJsonUrl);
    });

    const query = ref(sampleQuery);
    const eventBus = mitt();
    const submitQuery = () => {
      eventBus.emit("getValue");
    };

    const executedQuery = computed(() => {
      let result = <ParamsObject[]>[];
      let error: Error | null = null;
      try {
        result = executeQuery(query.value);
      } catch (err) {
        error = err;
      }
      return {
        result,
        error,
      };
    });
    const queryResults = computed(() => executedQuery.value.result);
    const queryResultsJson = computed(() => JSON.stringify(queryResults.value, null, 2));
    const queryError = computed(() => executedQuery.value.error);

    const submitQueryHotkeyDisplay = isApplePlatform() ? "⌘⏎" : "Ctrl+⏎";
    const submitQueryHotkeys = "command+enter, ctrl+enter";
    onMounted(() => {
      hotkeys(submitQueryHotkeys, () => {
        submitQuery();
      });
    });
    onBeforeUnmount(() => {
      hotkeys.unbind(submitQueryHotkeys);
    });
    // Ace traps keys, so we need to configure the keybinding separately for Ace.
    const submitQueryAceCommand = {
      name: "submitQuery",
      bindKey: { win: "Ctrl-Enter", mac: "Cmd-Enter" },
      exec: () => {
        submitQuery();
      },
    };

    return {
      researchesJson,
      researchedJsonUrl,
      schema,
      query,
      submitQuery,
      submitQueryHotkeyDisplay,
      submitQueryAceCommand,
      eventBus,
      queryResultsJson,
      queryError,
    };
  },
});
</script>
