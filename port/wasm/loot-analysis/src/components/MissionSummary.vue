<template>
  <div class="text-center text-xs text-gray-700">
    <p>{{ mission.missionCount }} missions</p>
    <template
      v-for="category in [
        { id: 'Artifacts (Rare)', display: 'rares' },
        { id: 'Artifacts (Epic)', display: 'epics' },
        { id: 'Artifacts (Legendary)', display: 'legendaries' },
      ]"
      :key="category.id"
    >
      <p>
        {{ categoryTotal(category.id) }} {{ category.display }},
        {{ categoryPerMission(category.id).toPrecision(2) }} per mission
      </p>
    </template>
  </div>
</template>

<script>
import { toRefs } from "vue";

export default {
  props: {
    mission: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { mission } = toRefs(props);
    const categoryTotal = categoryName => {
      for (const category of mission.value.categories) {
        if (category.categoryName === categoryName) {
          return category.stats.reduce((sum, item) => sum + item.count, 0);
        }
      }
      return 0;
    };
    const categoryPerMission = categoryName =>
      categoryTotal(categoryName) / mission.value.missionCount;
    return {
      categoryTotal,
      categoryPerMission,
    };
  },
};
</script>
