<template>
  <div class="flex-1 max-w-7xl w-full mx-auto px-4 xl:px-0 my-4 text-sm">
    <section id="intro">
      <p class="mx-10%">
        This page is dedicated to the statistical analysis of Egg, Inc. spaceship mission rewards.
      </p>

      <p class="mx-10%">
        Empirical data used for the analysis comes from user contributions to
        <a
          href="https://ei.mikit.app/conribute_data"
          target="_blank"
          class="text-blue-500 hover:text-blue-700"
          >ei.mikit.app</a
        >.
        <span class="text-green-500"
          >It would be much appreciated if you visit and contribute yours;</span
        >
        it is as easy as entering your ID and pressing a few buttons.
      </p>

      <p class="mx-10% mt-2">
        Plotted on this page are per-item expectations (total number received divided by number of
        missions) normalized (i.e. divided) by the odds multiplier, a not-very-well-understood
        internal parameter for each item. You can find the odds multipliers in the config table at
        the bottom of the
        <a
          href="https://wasmegg.netlify.app/artifact-explorer/"
          target="_blank"
          class="text-blue-500 hover:text-blue-700"
          >artifact explorer</a
        >.
      </p>

      <p class="mx-10%">
        Confidence interval for each plotted data point is estimated using the
        <a
          href="https://en.wikipedia.org/wiki/Binomial_proportion_confidence_interval#Wilson_score_interval"
          target="_blank"
          class="text-blue-500 hover:text-blue-700"
          >Wilson score interval</a
        >, at 95% confidence level by default. The underlying assumption is that given a particular
        mission type, each item has a fixed probability to drop, so for each item we're essentially
        estimating that fixed probability from a series of Bernoulli trials (with two outcomes: item
        X or not), where the Wilson score interval is appropriate.
      </p>

      <p class="mx-10% mt-2">
        Plots on this page are powered by
        <a
          href="https://echarts.apache.org/en/index.html"
          target="_blank"
          class="text-blue-500 hover:text-blue-700"
          >Apache ECharts</a
        >.
      </p>

      <p class="mx-10% text-green-500">
        You can click on legends to toggle series on and off. Hover or click on data points to
        reveal details.
      </p>
    </section>

    <section id="toc" class="mx-10% my-2 grid grid-cols-1 gap-1 sm:grid-cols-3">
      <template v-for="mission in missions.slice().reverse()" :key="mission.info.id">
        <a :href="`#${mission.info.id}`" class="text-xs leading-tight">
          {{ mission.info.display }}
        </a>
      </template>
    </section>

    <template v-for="mission in missions.slice().reverse()" :key="mission.info.id">
      <section :id="mission.info.id">
        <header class="mx-4 mt-4 mb-1">
          <a href="#intro" class="block">
            <h2
              :id="mission.info.id"
              class="text-center text-base leading-6 font-medium text-gray-900"
            >
              {{ mission.info.display }}
            </h2>
          </a>
          <mission-summary :mission="mission" />
        </header>
        <loot-chart :items="items" :missionStats="mission" />
      </section>
    </template>
  </div>
</template>

<script>
import MissionSummary from "./components/MissionSummary.vue";
import LootChart from "./components/LootChart.vue";

import data from "@/app-data.json";

export default {
  components: {
    MissionSummary,
    LootChart,
  },

  setup() {
    const items = Object.fromEntries(data.items.map(item => [item.id, item]));
    const missions = data.stats;
    return {
      items,
      missions,
    };
  },
};
</script>
