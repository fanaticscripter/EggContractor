<template>
  <div class="max-w-7xl mx-auto px-4 mb-3 text-sm space-y-1">
    <div v-if="showIntroduction">
      <p>
        The loot simulator simulates loot drops for a set of missions of your choice (collectively
        called a trial) a large number of times, checking whether the dropped items are enough to
        satisfy a set of items of your choice (crafting and demotions allowed) for each trial, in
        order to estimate the probability of gathering these items after running the specified set
        of missions.
      </p>
      <p>
        Drop rates are based on data submitted to
        <a
          href="https://ei.mikit.app/contribute_data"
          target="_blank"
          class="text-blue-600 hover:text-blue-700"
          >ei.mikit.app</a
        >
        which are independently viewable in
        <a
          href="https://wasmegg.netlify.app/artifact-explorer/"
          target="_blank"
          class="text-blue-600 hover:text-blue-700"
          >artifact explorer</a
        >. Note that for the purpose of this simulator, aforementioned drop rates data are treated
        as exact, whereas in reality the empirical data obviously have errors. The accuracy depends
        on the sample size, which is noted for each mission in the mission selector below (the
        number of missions recorded is in parentheses, in the form of "&lt;number&gt; on file").
      </p>
      <p>
        This simulator does not take rarity into account, since (1) the sample size of uncommon
        drops from missions is rather limited (you can check what we have in the artifact explorer);
        (2) we have very limited public data on rarity distributions from crafting; (3) supporting
        rarities increases the complexity of this simulator rather significantly. You can use
        empirical data and values of the odds multiplier parameter (both available in the artifact
        explorer) for a rough idea of rarity distribution.
      </p>
    </div>
    <div class="text-blue-600 hover:text-blue-700 cursor-pointer" @click="toggleIntroduction">
      <template v-if="showIntroduction">Collpase this introduction</template>
      <template v-else>Show introduction</template>
    </div>
  </div>

  <base-error-boundary>
    <template #default>
      <Suspense>
        <template #default>
          <simulator-container class="flex-1" />
        </template>
        <template #fallback>
          <base-loading>Checking browser capabilities...</base-loading>
        </template>
      </Suspense>
    </template>

    <template #error="{ error }">
      <div class="max-w-7xl mx-auto px-4">
        <template v-if="isModuleWorkerNotSupportedError(error)">
          <div class="text-sm text-red-500">
            This simulator is only supported on Chromium-based browsers at the moment. Please try
            Google Chrome instead, or upgrade to the latest version if you are already using it.
            Note that third-party browsers on iOS, including Google Chrome, are reskinned Safari,
            and as such does not support this simulator.
          </div>
          <div class="text-sm mt-2">
            You might want to know: it takes about 457 Extended Henerprise missions to have a ~50%
            chance to craft a Gilded Book of Basan (T4).
          </div>
        </template>
        <template v-else>
          <div class="text-sm mb-1">Unexpected error occurred, please report to the author:</div>
          <pre class="text-xs text-red-500">{{ error.toString() }}</pre>
          <pre class="text-xs text-red-500">{{ error.stack }}</pre>
        </template>
      </div>
    </template>
  </base-error-boundary>
</template>

<script lang="ts">
import { defineComponent, ref } from "vue";

import BaseErrorBoundary from "./components/BaseErrorBoundary.vue";
import BaseLoading from "./components/BaseLoading.vue";
import SimulatorContainer from "@/components/SimulatorContainer.vue";
import { ModuleWorkerNotSupportedError } from "./errors";
import { getLocalStorage, setLocalStorage } from "./storage";

const SHOW_INTRODUCTION_KEY = "showIntroduction";

export default defineComponent({
  components: {
    BaseErrorBoundary,
    BaseLoading,
    SimulatorContainer,
  },
  setup() {
    const showIntroduction = ref(getLocalStorage(SHOW_INTRODUCTION_KEY) !== "false");
    const toggleIntroduction = () => {
      showIntroduction.value = !showIntroduction.value;
      setLocalStorage(SHOW_INTRODUCTION_KEY, showIntroduction.value);
    };
    const isModuleWorkerNotSupportedError = (err: Error) => {
      return err instanceof ModuleWorkerNotSupportedError;
    };
    return {
      showIntroduction,
      toggleIntroduction,
      isModuleWorkerNotSupportedError,
    };
  },
});
</script>
