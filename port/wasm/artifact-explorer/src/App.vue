<template>
  <div class="flex-1 max-w-7xl w-full mx-auto px-4 xl:px-0 my-4">
    <router-view :missions="missions" :artifacts="artifacts" :lootTable="lootTable"></router-view>
  </div>

  <h2 class="mt-4 mb-2 text-center text-base leading-6 font-medium text-gray-900">
    <code>/ei_afx/config</code> data
  </h2>
  <eiafx-config-table :ships="ships" :missions="missions" :artifacts="artifacts" />
</template>

<script>
import EiafxConfigTable from "@/components/EiafxConfigTable.vue";

import data from "@/app-data.json";

export default {
  components: {
    EiafxConfigTable,
  },

  data() {
    const { ships, missions, artifacts, lootTable } = data;
    return Object.freeze({
      ships,
      missions,
      artifacts,
      lootTable,
    });
  },
};
</script>
