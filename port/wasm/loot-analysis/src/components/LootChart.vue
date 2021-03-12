<template>
  <div class="flex-1 max-w-7xl w-full mx-auto px-4 xl:px-0 my-4">
    <v-chart class="chart" :option="option" />
  </div>
</template>

<script>
import { use } from "echarts/core";
import { CanvasRenderer } from "echarts/renderers";
import { PieChart, ScatterChart } from "echarts/charts";
import { LegendComponent, TitleComponent, TooltipComponent } from "echarts/components";
import { ref } from "vue";
import VChart, { THEME_KEY } from "vue-echarts";

use([CanvasRenderer, PieChart, ScatterChart, LegendComponent, TitleComponent, TooltipComponent]);

export default {
  components: {
    VChart,
  },

  provide: {
    [THEME_KEY]: "macarons",
  },

  setup() {
    const option = ref({
      // title: {
      //   text: "Traffic Sources",
      //   left: "center",
      // },
      // tooltip: {
      //   trigger: "item",
      //   formatter: "{a} <br/>{b} : {c} ({d}%)",
      // },
      // legend: {
      //   orient: "vertical",
      //   left: "left",
      //   data: ["Direct", "Email", "Ad Networks", "Video Ads", "Search Engines"],
      // },
      tooltip: {},
      legend: {},
      xAxis: {},
      yAxis: {},
      series: [
        {
          type: "scatter",
          name: "Artifacts (common)",
          dimensions: ["Quality", "Normalized expectation"],
          data: [
            [8.9, 2],
            [9.3, 4],
            [11.5, 3],
            [12.9, 1],
          ],
          encode: {
            x: "Quality",
            y: "Normalized expectation",
            tooltip: ["Quality", "Normalized expectation"],
          },
        },
      ],
    });

    return { option };
  },
};
</script>

<style scoped>
.chart {
  height: 400px;
}
</style>
