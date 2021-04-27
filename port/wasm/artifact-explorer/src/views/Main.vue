<template>
  <spoiler-alert />
  <mission-selector
    :initialMissionId="$route.params.missionId"
    :missions="missions"
    :key="$route.path"
  />
  <router-view name="mission" :missions="missions" :artifacts="artifacts" :lootTable="lootTable" />
  <artifact-selector
    :initialArtifactId="$route.params.artifactId"
    :artifacts="artifacts"
    :key="$route.path"
  />
  <div class="my-4 text-xs text-red-900">
    <p class="font-medium">Artifact notes:</p>
    <p>
      * Certain effect values shown may be 1% higher than the corresponding in-game values; those
      are caused by erroneous floating point handling in the game, i.e. values here are correct.
    </p>
    <p>
      &dagger; Artifacts marked with &dagger; are not available from missions, backed by both sound
      theory and observation.
    </p>
    <p>
      &Dagger; Artifacts marked with &Dagger; are not available from missions, backed by observation
      only.
    </p>
  </div>
  <router-view name="artifact" :missions="missions" :artifacts="artifacts" :lootTable="lootTable" />
  <artifact-grid :artifacts="artifacts" />
</template>

<script>
import SpoilerAlert from "@/components/SpoilerAlert.vue";
import ArtifactGrid from "@/components/ArtifactGrid.vue";
import ArtifactSelector from "@/components/ArtifactSelector.vue";
import MissionSelector from "@/components/MissionSelector.vue";

export default {
  components: {
    SpoilerAlert,
    ArtifactGrid,
    ArtifactSelector,
    MissionSelector,
  },

  props: {
    missions: Array,
    artifacts: Array,
    lootTable: Object,
  },
};
</script>
