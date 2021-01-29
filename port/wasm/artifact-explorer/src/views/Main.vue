<template>
  <mission-selector :initialMissionId="$route.params.missionId" :missions="missions" :key="$route.path" />
  <router-view name="mission" :missions="missions" :artifacts="artifacts" />
  <artifact-selector :initialArtifactId="$route.params.artifactId" :artifacts="artifacts" :key="$route.path" />
  <router-view name="artifact" :missions="missions" :artifacts="artifacts" />
  <artifact-grid :artifacts="artifacts" />
</template>

<script>
import ArtifactGrid from "@/components/ArtifactGrid.vue";
import ArtifactSelector from "@/components/ArtifactSelector.vue";
import MissionSelector from "@/components/MissionSelector.vue";

export default {
  components: {
    ArtifactGrid,
    ArtifactSelector,
    MissionSelector,
  },

  props: {
    missions: Array,
    artifacts: Array,
  },
};
</script>
