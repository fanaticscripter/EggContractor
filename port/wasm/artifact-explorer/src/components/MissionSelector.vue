<template>
  <div class="my-2 space-y-2">
    <div class="w-full max-w-xs">
      <label for="mission" class="block text-sm font-medium text-gray-700">
        What can I possibly get from this mission?
      </label>
      <select
        id="mission"
        name="mission"
        class="mt-1 block w-full pl-3 pr-10 py-1 text-sm bg-gray-50 border-gray-300 focus:outline-none focus:ring-green-500 focus:border-green-500 rounded-md"
        v-model="missionId"
      >
        <option value="">-- Select a mission --</option>
        <option v-for="mission in missions" :key="mission.id" :value="mission.id">
          {{ mission.display }}
        </option>
      </select>
    </div>
    <p class="text-xs text-gray-500">
      You may use this dropdown, or click on any mission name on the page, including, for instance,
      results in an artifact query.
    </p>
  </div>
</template>

<script>
export default {
  props: {
    initialMissionId: String,
    missions: Array,
  },

  data() {
    return {
      missionId: this.initialMissionId || "",
    };
  },

  watch: {
    missionId() {
      if (this.missionId === "") {
        this.$router.push({
          name: "home",
        });
      } else {
        this.$router.push({
          name: "mission",
          params: {
            missionId: this.missionId,
          },
        });
      }
    },
  },
};
</script>
