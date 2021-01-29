<template>
  <div class="-mx-4 sm:mx-0 bg-gray-50 overflow-hidden sm:rounded-lg sm:shadow-md">
    <div class="bg-gray-100 px-4 py-4 border-b border-gray-200 sm:px-6">
      <artifact-name :artifact="artifact" :showTier="true" />
    </div>

    <template v-if="artifact.tierNumber > 1">
      <div class="px-4 py-4 sm:px-6 space-y-2">
        <div class="text-sm font-medium text-gray-500">Crafting recipe:</div>
        <div class="text-sm">Coming soon&trade;</div>
      </div>
      <hr />
    </template>

    <div class="px-4 py-4 sm:px-6 space-y-2">
      <template v-if="obtainableMissions.length > 0">
        <div class="text-sm font-medium text-gray-500">Available from the following missions:</div>
        <ul class="grid grid-cols-1 gap-x-4 gap-y-1 sm:grid-cols-2 xl:grid-cols-3">
          <li v-for="mission in obtainableMissions" :key="mission.id">
            <mission-name :mission="mission" />
          </li>
        </ul>
      </template>
      <template v-else>
        <div class="text-sm font-medium text-gray-500">Not available from missions :(</div>
      </template>
    </div>
  </div>
</template>

<script>
import ArtifactName from "@/components/ArtifactName.vue";
import MissionName from "@/components/MissionName.vue";

export default {
  components: {
    ArtifactName,
    MissionName,
  },

  props: {
    artifactId: String,
    missions: Array,
    artifacts: Array,
  },

  computed: {
    artifact() {
      for (const artifact of this.artifacts) {
        if (artifact.itemId === this.artifactId && artifact.rarity === 0) {
          return artifact;
        }
      }
      return undefined;
    },

    obtainableMissions() {
      if (!this.artifact) {
        return [];
      }
      return this.missions.filter(
        mission =>
          mission.minQuality <= this.artifact.quality && this.artifact.quality <= mission.maxQuality
      );
    },
  },
};
</script>
