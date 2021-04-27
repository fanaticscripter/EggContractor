<template>
  <div class="flex flex-col h-screen">
    <div class="flex-grow overflow-auto">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="divide-y">
          <tr class="divide-x">
            <th
              scope="col"
              rowspan="2"
              class="sticky top-0 z-10 sm:left-0 px-4 py-1 text-left text-xs font-medium text-gray-500 bg-gray-50"
            >
              <div
                class="flex items-center justify-center cursor-pointer"
                @click="toggleSort($options.SORT_BY.ITEM)"
              >
                Item<sort-arrow :direction="sortButtonDirection($options.SORT_BY.ITEM)" />
              </div>
            </th>

            <th
              scope="col"
              rowspan="2"
              class="sticky top-0 px-4 py-1 text-center text-xs font-medium text-gray-500 bg-gray-50"
            >
              <div
                class="flex items-center justify-center cursor-pointer"
                @click="toggleSort($options.SORT_BY.TIER)"
              >
                Tier<sort-arrow :direction="sortButtonDirection($options.SORT_BY.TIER)" />
              </div>
            </th>

            <th
              scope="col"
              rowspan="2"
              class="sticky top-0 px-4 py-1 text-center text-xs font-medium text-gray-500 bg-gray-50"
            >
              <div
                class="flex items-center justify-center cursor-pointer"
                @click="toggleSort($options.SORT_BY.QUALITY)"
              >
                Base quality<sort-arrow
                  :direction="sortButtonDirection($options.SORT_BY.QUALITY)"
                />
              </div>
            </th>

            <th
              scope="col"
              rowspan="2"
              class="sticky top-0 px-4 py-1 text-center text-xs font-medium text-gray-500 bg-gray-50"
            >
              <div
                class="flex items-center justify-center cursor-pointer"
                @click="toggleSort($options.SORT_BY.ODDS_MULTIPLIER)"
              >
                Odds multiplier<sort-arrow
                  :direction="sortButtonDirection($options.SORT_BY.ODDS_MULTIPLIER)"
                />
              </div>
            </th>

            <th
              v-for="ship in ships"
              :key="ship.name"
              scope="col"
              colspan="3"
              class="sticky top-0 px-4 py-1 text-center text-xs font-medium text-gray-500 bg-gray-50"
              v-tippy="{
                content: `
                  <div>
                    <img src='${iconURL(ship.iconPath, 256)}' class='h-16 w-16 mx-auto'>
                    <div class='text-center'>${ship.name}<div>
                  </div>`,
                allowHTML: true,
              }"
            >
              {{ ship.abbrevName }}
            </th>

            <th
              scope="col"
              rowspan="2"
              class="sticky top-0 px-4 py-1 text-center text-xs font-medium text-gray-500 bg-gray-50"
              @click="toggleSort($options.SORT_BY.VALUE)"
            >
              <div class="flex items-center justify-center cursor-pointer">
                Value<sort-arrow :direction="sortButtonDirection($options.SORT_BY.VALUE)" />
              </div>
            </th>

            <th
              scope="col"
              colspan="4"
              class="sticky top-0 px-4 py-1 text-center text-xs font-medium text-gray-500 bg-gray-50 cursor-pointer"
              @click="toggleSort($options.SORT_BY.CRAFTING_PRICE)"
            >
              <div class="flex items-center justify-center">
                Crafting price<sort-arrow
                  :direction="sortButtonDirection($options.SORT_BY.CRAFTING_PRICE)"
                />
              </div>
            </th>
          </tr>

          <tr class="divide-x">
            <th
              v-for="mission in missions"
              :key="mission.id"
              class="sticky top-6 px-2 py-1 text-center text-xs font-medium text-gray-500 bg-gray-50 border border-gray-200 cursor-pointer"
              v-tippy="{
                content:
                  `${mission.display}, ` +
                  `${mission.minQuality.toFixed(1)} - ${mission.maxQuality.toFixed(1)}`,
              }"
            >
              <router-link :to="{ name: 'mission', params: { missionId: mission.id } }">
                {{ mission.abbrevType }}
              </router-link>
            </th>

            <th
              class="sticky top-6 px-2 py-1 text-center text-xs font-medium text-gray-500 bg-gray-50 border"
            >
              Base
            </th>

            <th
              class="sticky top-6 px-2 py-1 text-center text-xs font-medium text-gray-500 bg-gray-50 border border-gray-200"
            >
              Low
            </th>
            <th
              class="sticky top-6 px-2 py-1 text-center text-xs font-medium text-gray-500 bg-gray-50 border"
            >
              Domain
            </th>

            <th
              class="sticky top-6 px-2 py-1 text-center text-xs font-medium text-gray-500 bg-gray-50 border border-gray-200"
            >
              Curve
            </th>
          </tr>
        </thead>
        <tbody class="divide-y">
          <tr v-for="artifact in sortedArtifacts" :key="artifact.sortKey" class="divide-x">
            <td class="sm:sticky left-0 px-4 py-1 whitespace-nowrap text-sm text-gray-500 bg-white">
              <artifact-name :artifact="artifact" />
            </td>

            <td class="px-4 py-1 whitespace-nowrap text-center text-sm text-gray-500">
              {{ artifact.tier_number }}
            </td>

            <td class="px-4 py-1 whitespace-nowrap text-center text-sm text-gray-500">
              {{ artifact.quality }}
            </td>

            <td class="px-4 py-1 whitespace-nowrap text-center text-sm text-gray-500">
              {{ formatSmallFloat(artifact.params.odds_multiplier) }}
            </td>

            <template v-for="mission in missions" :key="mission.id">
              <td
                v-if="
                  artifact.quality >= mission.minQuality && artifact.quality <= mission.maxQuality
                "
                class="px-2 py-1 whitespace-nowrap text-center text-sm text-gray-500 bg-green-300"
                v-tippy="{
                  content: `${mission.display}<br>${artifact.name}, ${artifact.rarity}`,
                  allowHTML: true,
                }"
              >
                &check;
              </td>
              <td
                v-else
                class="px-2 py-1 whitespace-nowrap text-center text-sm bg-gray-300"
                v-tippy="{
                  content: `${mission.display}<br>${artifact.name}, ${artifact.rarity}`,
                  allowHTML: true,
                }"
              ></td>
            </template>

            <td class="px-4 py-1 whitespace-nowrap text-center text-sm text-gray-500">
              {{ artifact.params.value.toFixed(1) }}
            </td>

            <td
              class="px-4 py-1 whitespace-nowrap text-center text-sm"
              :class="[craftable(artifact) ? 'text-gray-500' : 'text-gray-300']"
            >
              {{ artifact.params.crafting_price.toFixed(2) }}
            </td>

            <td
              class="px-4 py-1 whitespace-nowrap text-center text-sm"
              :class="[craftable(artifact) ? 'text-gray-500' : 'text-gray-300']"
            >
              {{ artifact.params.crafting_price_low.toFixed(2) }}
            </td>

            <td
              class="px-4 py-1 whitespace-nowrap text-center text-sm"
              :class="[craftable(artifact) ? 'text-gray-500' : 'text-gray-300']"
            >
              {{ artifact.params.crafting_price_domain }}
            </td>

            <td
              class="px-4 py-1 whitespace-nowrap text-center text-sm"
              :class="[craftable(artifact) ? 'text-gray-500' : 'text-gray-300']"
            >
              {{ artifact.params.crafting_price_curve }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script>
import ArtifactName from "@/components/ArtifactName.vue";
import SortArrow from "@/components/SortArrow.vue";

import { iconURL, stringCmp } from "@/utils";

const SORT_BY = {
  ITEM: "item",
  TIER: "tier",
  QUALITY: "quality",
  ODDS_MULTIPLIER: "odds_multiplier",
  VALUE: "value",
  CRAFTING_PRICE: "crafting_price",
};

export default {
  components: {
    ArtifactName,
    SortArrow,
  },

  props: {
    ships: Array,
    missions: Array,
    artifacts: Array,
  },

  data() {
    return {
      sortBy: SORT_BY.QUALITY,
      sortAscending: true,
    };
  },

  SORT_BY,

  computed: {
    sortedArtifacts() {
      if (!this.artifacts) {
        return this.artifacts;
      }
      return [...this.artifacts].sort((a1, a2) => {
        let cmp = 0;
        switch (this.sortBy) {
          case SORT_BY.ITEM:
            cmp = stringCmp(a1.sortKey, a2.sortKey);
            break;
          case SORT_BY.TIER:
            cmp = a1.tier_number - a2.tier_number;
            break;
          case SORT_BY.ODDS_MULTIPLIER:
            cmp = a1.params.odds_multiplier - a2.params.odds_multiplier;
            break;
          case SORT_BY.VALUE:
            cmp = a1.params.value - a2.params.value;
            break;
          case SORT_BY.CRAFTING_PRICE:
            cmp = a1.params.crafting_price - a2.params.crafting_price;
            break;
        }
        // QUALITY as secondary condition.
        if (cmp === 0) {
          cmp = a1.params.base_quality - a2.params.base_quality;
        }
        return this.sortAscending ? cmp : -cmp;
      });
    },
  },

  methods: {
    iconURL,

    sortButtonDirection(criterion) {
      if (this.sortBy === criterion) {
        return this.sortAscending ? 1 : -1;
      }
      return 0;
    },

    toggleSort(criterion) {
      if (this.sortBy === criterion) {
        this.sortAscending = !this.sortAscending;
      } else {
        this.sortBy = criterion;
        this.sortAscending = true;
      }
    },

    craftable(artifact) {
      return artifact.tier_number > 1 && artifact.afxRarity === 0;
    },

    formatSmallFloat(x) {
      if (x === 0) {
        return "0";
      }
      if (x >= 1) {
        let s = x.toFixed(2);
        return s.replace(/\.?0+$/, "");
      }
      let s = x.toPrecision(2);
      return s.replace(/0+$/, "");
    },
  },
};
</script>
