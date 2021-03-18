<template>
  <div class="mx-auto max-w-xs">
    <select
      :id="`confidence-${id}`"
      :name="`confidence-${id}`"
      class="mt-1 block w-full pl-3 pr-10 py-1 text-base bg-gray-50 border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md"
      v-model="confidenceLevel"
    >
      <option value="68%">68% confidence level</option>
      <option value="95%">95% confidence level</option>
      <option value="99.7%">99.7% confidence level</option>
    </select>
  </div>

  <v-chart class="chart" :option="option" autoresize />
</template>

<script>
import { use } from "echarts/core";
import { CanvasRenderer } from "echarts/renderers";
import { PieChart, ScatterChart, CustomChart } from "echarts/charts";
import {
  DataZoomComponent,
  DataZoomSliderComponent,
  LegendComponent,
  MarkLineComponent,
  ToolboxComponent,
  TooltipComponent,
} from "echarts/components";
import { computed, ref } from "vue";
import VChart, { THEME_KEY } from "vue-echarts";

use([
  CanvasRenderer,
  PieChart,
  ScatterChart,
  CustomChart,
  DataZoomComponent,
  DataZoomSliderComponent,
  LegendComponent,
  MarkLineComponent,
  ToolboxComponent,
  TooltipComponent,
]);

// SVG shapes from http://thenewcode.com/1037/SVG-Shape-Elements-Polygons
// <polygon> converted to <path> with https://stackoverflow.com/a/10730247.
const tierSymbols = {
  1: {
    // Regular triangle
    path: "M 50 15,100 100,0 100z",
    // (5\sqrt{3}, 15/2)
    size: [8.66, 7.5],
  },
  2: {
    // Square
    path: "M 0 0,0 100,100 100,100 0z",
    // (5\sqrt{2}, 5\sqrt{2})
    size: [7.07, 7.07],
  },
  3: {
    // Regular pentagon
    path: "M 26,86 11.2,40.4 50,12.2 88.8,40.4 74,86 z",
    // (10 \cos(\pi/10), 10 \cos^2(\pi/10))
    size: [9.51, 9.05],
  },
  4: {
    // Regular hexagon
    path: "M 30.1,84.5 10.2,50 30.1,15.5 69.9,15.5 89.8,50 69.9,84.5z",
    // (10, 5\sqrt{3})
    size: [10, 8.66],
  },
};

// z_{1-\alpha/2}, where \alpha = 1 - confidence level.
// http://www.z-table.com/
const zscores = {
  "68%": 1.0,
  "95%": 1.96,
  "99.7%": 2.97,
};

// Render error bar. Adapted from
// https://echarts.apache.org/examples/en/editor.html?c=custom-error-scatter.
function renderItem(params, api) {
  const encode = params.encode;
  const x = api.value(encode.x);
  const yMax = api.value(encode.y[1]);
  const yMin = api.value(encode.y[2]);
  const highPoint = api.coord([x, yMax]);
  const lowPoint = api.coord([x, yMin]);

  const halfWidth = 2;
  const style = api.style({
    stroke: api.visual("color"),
    fill: null,
  });

  return {
    type: "group",
    children: [
      {
        type: "line",
        transition: ["shape"],
        shape: {
          x1: highPoint[0] - halfWidth,
          y1: highPoint[1],
          x2: highPoint[0] + halfWidth,
          y2: highPoint[1],
        },
        style,
      },
      {
        type: "line",
        transition: ["shape"],
        shape: {
          x1: highPoint[0],
          y1: highPoint[1],
          x2: lowPoint[0],
          y2: lowPoint[1],
        },
        style,
      },
      {
        type: "line",
        transition: ["shape"],
        shape: {
          x1: lowPoint[0] - halfWidth,
          y1: lowPoint[1],
          x2: lowPoint[0] + halfWidth,
          y2: lowPoint[1],
        },
        style,
      },
    ],
  };
}

export default {
  components: {
    VChart,
  },

  provide: {
    [THEME_KEY]: "infographic",
  },

  props: {
    items: Object,
    missionStats: Object,
  },

  setup({ items, missionStats }) {
    const confidenceLevel = ref("95%");
    const zscore = computed(() => zscores[confidenceLevel.value]);

    const missionCount = missionStats.missionCount;
    const info = missionStats.info;
    const id = info.id;
    const capacity = info.capacity;
    const total = missionCount * capacity;

    const dimensions = [
      "Name",
      "Tier",
      "Quality",
      "Odds multipler",
      "Count",
      "Expectation",
      "Normalized expectation",
      "Normalized expectation (upper)",
      "Normalized expectation (lower)",
    ];
    const tooltipDimensions = [
      "Quality",
      "Odds multipler",
      "Count",
      "Expectation",
      "Normalized expectation",
      "Normalized expectation (upper)",
      "Normalized expectation (lower)",
    ];
    const data = computed(() =>
      missionStats.categories.map(category => ({
        name: category.categoryName,
        data: category.stats.map(entry => {
          const itemId = entry.itemId;
          const item = items[itemId];
          // Confidence interval is estimated using Wilson score interval.
          // https://en.wikipedia.org/wiki/Binomial_proportion_confidence_interval#Wilson_score_interval
          const p = entry.count / total;
          const n = total;
          const z = zscore.value;
          const w = (p + (z * z) / (2 * n)) / (1 + (z * z) / n);
          const pm = (z / (1 + (z * z) / n)) * Math.sqrt((p * (1 - p)) / n + (z * z) / (4 * n * n));
          const wilsonPlus = w + pm;
          const wilsonMinus = w - pm;
          const expectation = p * capacity;
          const precision = Math.min(entry.count.toString().length, 3);
          const normalizedExpectation = expectation / item.oddsMultiplier;
          const normalizedExpectationUpper = (wilsonPlus * capacity) / item.oddsMultiplier;
          // In theory this is nonnegative so the Math.max(..., 0) is unnecessary;
          // in practice we may get a tiny negative number due to floating point
          // imprecision, stretching plots into the negative, which is undesirable.
          const normalizedExpectationLower = Math.max(
            (wilsonMinus * capacity) / item.oddsMultiplier,
            0
          );
          // calc and display positive and negative error
          return [
            item.name,
            item.tier.tier_number,
            item.quality,
            item.oddsMultiplier.toPrecision(3),
            `${entry.count} / ${missionCount}`,
            expectation.toPrecision(precision),
            normalizedExpectation.toPrecision(3),
            normalizedExpectationUpper.toPrecision(3),
            normalizedExpectationLower.toPrecision(3),
          ];
        }),
      }))
    );

    const option = computed(() => ({
      color: [
        "#27727b",
        "#60a5fa",
        "#a78bfa",
        "#fbbf24",
        "#c1232b",
        "#9bca63",
        "#e87c25",

        "#aaaaaa",
        "#aaaaaa",
        "#aaaaaa",
        "#aaaaaa",
      ],
      dataZoom: [
        {
          xAxisIndex: [0],
          type: "slider",
          filterMode: "filter",
        },
        {
          yAxisIndex: [0],
          type: "slider",
        },
      ],
      grid: {
        top: 40,
      },
      legend: {
        left: "left",
        orient: "vertical",
        selected: {
          "Artifacts (Rare)": false,
          "Artifacts (Epic)": false,
          "Artifacts (Legendary)": false,
        },
      },
      toolbox: {
        feature: {
          dataZoom: {},
          restore: {},
          saveAsImage: {
            title: "Save",
            type: "png",
            name: `Loot analysis - ${info.display}`,
            pixelRatio: 2,
          },
        },
      },
      tooltip: {},
      xAxis: {
        min: parseFloat(info.minQuality.toFixed(1)),
        max: parseFloat(info.maxQuality.toFixed(1)),
      },
      yAxis: {},
      series: data.value
        .map(category => [
          {
            type: "scatter",
            name: category.name,
            dimensions,
            encode: {
              x: "Quality",
              y: "Normalized expectation",
              tooltip: tooltipDimensions,
              itemName: "Name",
            },
            data: category.data,
            symbol: (value, params) => {
              const tier = value[1];
              return tier in tierSymbols ? `path://${tierSymbols[tier].path}` : "circle";
            },
            symbolSize: (value, params) => {
              const tier = value[1];
              return tier in tierSymbols ? tierSymbols[tier].size : 10;
            },
            markLine: {
              silent: true,
              lineStyle: {
                color: "#000",
                type: "dashed",
                width: 1.5,
              },
              symbol: "none",
              data: [
                [
                  {
                    name: `Mission target quality (${info.quality.toFixed(1)})`,
                    label: {
                      position: "start",
                      offset: [0, -20],
                    },
                    xAxis: info.quality,
                    yAxis: 0,
                  },
                  {
                    xAxis: info.quality,
                    yAxis: "max",
                  },
                ],
              ],
            },
          },
          {
            type: "custom",
            name: category.name,
            dimensions,
            encode: {
              x: "Quality",
              y: [
                "Normalized expectation",
                "Normalized expectation (upper)",
                "Normalized expectation (lower)",
              ],
              tooltip: tooltipDimensions,
              itemName: 0,
            },
            data: category.data,
            renderItem,
            z: 100,
          },
        ])
        .flat()
        .concat(
          // Phantom series used to create legends for tier shapes.
          [1, 2, 3, 4].map(tier => ({
            type: "scatter",
            name: `Tier ${tier}`,
            symbol: `path://${tierSymbols[tier].path}`,
          }))
        ),
    }));

    return { id, confidenceLevel, option };
  },
};
</script>

<style scoped>
.chart {
  height: 800px;
}

@media (max-height: 900px) {
  .chart {
    height: calc(100vh - 100px);
  }
}
</style>
