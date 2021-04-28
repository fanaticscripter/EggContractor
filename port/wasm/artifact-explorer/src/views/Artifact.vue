<template>
  <div
    id="artifact-card"
    class="-mx-4 sm:mx-0 bg-gray-50 overflow-hidden sm:rounded-lg sm:shadow-md"
  >
    <div class="bg-gray-100 px-4 py-4 border-b border-gray-200 sm:px-6">
      <div class="-ml-4 -mt-2 flex items-center justify-between flex-wrap sm:flex-nowrap">
        <div class="ml-4 mt-2">
          <artifact-name
            :artifact="artifact"
            :showTier="true"
            :noLink="true"
            :noAvailabilityMarker="true"
          />
        </div>
        <div class="ml-4 mt-2 flex-shrink-0">
          <share :id="artifact.id" :domElementId="'artifact-card'" />
        </div>
      </div>
    </div>

    <template v-if="artifact.has_effects">
      <div class="px-2 py-4 sm:px-4 space-y-2">
        <div class="px-2 text-sm font-medium text-gray-500">
          {{ artifact.effects.length > 1 ? "Effects" : "Effect" }}:
        </div>
        <div class="flex flex-col space-y-2">
          <div class="flex-grow overflow-auto">
            <table class="tabular-nums">
              <tbody>
                <tr
                  v-for="rarity in artifact.effects"
                  :key="rarity.afx_rarity"
                  :class="rarityFgClass(rarity.afx_rarity)"
                >
                  <td
                    v-if="artifact.has_rarities"
                    class="h-5 px-2 whitespace-nowrap text-left text-sm leading-4"
                  >
                    {{ rarity.rarity }}
                  </td>
                  <td
                    v-if="rarity.slots !== null"
                    class="h-5 flex items-center px-2 whitespace-nowrap space-x-0.5"
                    v-tippy="{ content: `${rarity.slots} slots` }"
                  >
                    <div
                      v-for="index in rarity.slots"
                      :key="index"
                      class="flex h-4 w-4 items-center justify-center rounded bg-gray-400"
                    >
                      <img
                        class="h-5 w-5 max-w-none"
                        :src="iconURL('egginc-extras/icon_afx_stone_slot.png', 64)"
                        :style="{ filter: 'brightness(10)' }"
                      />
                    </div>
                  </td>
                  <td class="h-5 px-2 truncate text-left text-sm leading-4">{{ rarity.effect }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
      <hr />
    </template>

    <template v-if="artifact.craftable">
      <div class="px-4 py-4 sm:px-6 space-y-2">
        <div class="text-sm font-medium text-gray-500">Crafting recipe:</div>
        <div>
          <table class="tabular-nums">
            <tbody>
              <tr v-for="ingredient in artifact.recipe.ingredients" :key="ingredient.id">
                <td class="text-left text-sm">{{ ingredient.count }}&times;</td>
                <td class="pl-1">
                  <artifact-name :artifact="id2artifact[ingredient.id]" :showTier="true" />
                </td>
              </tr>
            </tbody>
          </table>
          <div class="my-0.5 -mx-0.5 flex items-center space-x-1">
            <img
              class="h-4 w-4"
              src="https://eggincassets.tcl.sh/64/egginc-extras/icon_golden_egg.png"
            />
            <span class="text-sm">
              {{ artifact.recipe.crafting_price.initial.toLocaleString("en-US") }} &ndash;
              {{ artifact.recipe.crafting_price.minimum.toLocaleString("en-US") }}
            </span>
            <info
              v-tippy="{
                content:
                  `The crafting price is determined by the following formula: ` +
                  `<img class='p-2 bg-white' src='${require('@/assets/crafting-price-formula.svg')}'>`,
                allowHTML: true,
              }"
            />
          </div>
        </div>
      </div>
      <hr />
    </template>

    <template v-if="!artifact.ingredients_available_from_missions">
      <div class="px-4 py-4 sm:px-6 space-y-2">
        <div class="flex items-center space-x-1">
          <span class="text-sm font-medium text-gray-500">Hard dependencies</span>
          <info
            v-tippy="{
              content:
                'For an item unobtainable from missions, the hard dependencies are the highest level mission-obtainable items in the crafting ingredient tree; i.e., you absolutely have to gather these ingredients to craft the item in question, no way to skip them.',
            }"
          />
        </div>
        <div>
          <table class="tabular-nums">
            <tbody>
              <tr v-for="ingredient in artifact.hard_dependencies" :key="ingredient.id">
                <td class="text-left text-sm">{{ ingredient.count }}&times;</td>
                <td class="pl-1">
                  <artifact-name :artifact="id2artifact[ingredient.id]" :showTier="true" />
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
      <hr />
    </template>

    <template v-if="recursiveIngredients.length > 0">
      <div class="px-4 py-4 sm:px-6 space-y-2">
        <div class="flex items-center space-x-1">
          <span class="text-sm font-medium text-gray-500">Recursive ingredients</span>
          <info
            v-tippy="{
              content:
                'Ingredients of ingredients, ingredients of ingredients of ingredients, etc.',
            }"
          />
        </div>
        <ul class="grid grid-cols-1 gap-x-4 gap-y-1 sm:grid-cols-2 xl:grid-cols-3">
          <li v-for="ingredient in recursiveIngredients" :key="ingredient.id">
            <artifact-name :artifact="ingredient" :showTier="true" />
          </li>
        </ul>
      </div>
      <hr />
    </template>

    <div class="px-4 py-4 sm:px-6 space-y-2">
      <template v-if="!artifact.notDroppableInPractice && obtainableMissions.length > 0">
        <div class="text-sm font-medium text-gray-500">Available from the following missions:</div>
        <ul class="grid grid-cols-1 gap-x-4 gap-y-1 sm:grid-cols-2 xl:grid-cols-3">
          <li
            v-for="mission in obtainableMissions"
            :key="mission.id"
            class="flex flex-row items-start flex-wrap"
          >
            <mission-name class="mr-1" :mission="mission" />
            <tippy>
              <span
                class="flex flex-row items-center text-sm"
                :class="mission.notEnoughData ? 'text-gray-500' : null"
                >({{ formatToPrecision(mission.itemsPerMission, mission.precision)
                }}<sup v-if="mission.notEnoughData" class="-top-0.5">?</sup>/<img
                  class="h-4 w-4"
                  :src="iconURL(mission.shipIconPath, 32)"
                />, {{ formatToPrecision(mission.itemsPerDay, mission.precision)
                }}<sup v-if="mission.notEnoughData" class="-top-0.5">?</sup>/d)</span
              >

              <template #content>
                <p>
                  Received {{ mission.itemCount.total }} from {{ mission.missionCount }} missions.
                </p>
                <p v-if="mission.itemsPerMission > 0">
                  1 item per
                  {{ formatToPrecision(1 / mission.itemsPerMission, mission.precision) }} missions,
                  or 1 item per
                  {{ formatToPrecision(1 / mission.itemsPerDay, mission.precision) }} days.
                </p>
                <p v-if="mission.notEnoughData">Not enough data, rate likely far from accurate.</p>
                <ul v-if="mission.itemCount.total > 0" :set="(total = mission.itemCount.total)">
                  <template v-for="rarity in artifact.effects" :key="rarity.afxRarity">
                    <li
                      v-if="rarity.afx_rarity !== 0"
                      :set="(count = mission.itemCount.rarities[rarity.afx_rarity])"
                    >
                      <span :class="rarityFgClassOnDarkBg(rarity.afx_rarity)"
                        >{{ rarity.rarity }}:</span
                      >
                      {{ count }}/{{ total }}, {{ formatPercentage(count, total)
                      }}<template v-if="count > 0"
                        >, 1 per
                        <span :class="rarityFgClassOnDarkBg(rarity.afx_rarity)">
                          {{ formattedNumMissionsPerItem(count, mission.missionCount) }}
                        </span>
                        missions</template
                      >
                    </li>
                  </template>
                </ul>
              </template>
            </tippy>
          </li>
        </ul>
        <loot-data-credit />
      </template>
      <template v-else>
        <div class="text-sm font-medium text-gray-500">Not available from missions :(</div>
      </template>
    </div>

    <template v-if="dependents.length > 0">
      <hr />
      <div class="px-4 py-4 sm:px-6 space-y-2">
        <div class="text-sm font-medium text-gray-500">Used in:</div>
        <ul class="grid grid-cols-1 gap-x-4 gap-y-1 sm:grid-cols-2 xl:grid-cols-3 tabular-nums">
          <li v-for="dependent in dependents" :key="dependent.item.id" class="flex">
            <span class="text-sm pr-1">{{ dependent.count }}x in</span>
            <artifact-name :artifact="dependent.item" :showTier="true" />
          </li>
        </ul>
      </div>
    </template>

    <template v-if="hardDependents.length > 0">
      <hr />
      <div class="px-4 py-4 sm:px-6 space-y-2">
        <div class="flex items-center space-x-1">
          <span class="text-sm font-medium text-gray-500">Hard dependents</span>
          <info
            v-tippy="{
              content:
                'For a mission-obtainable item, a hard dependent is a mission-unobtainable item that must use the indicated number of this item in its crafting tree.',
            }"
          />
        </div>
        <ul class="grid grid-cols-1 gap-x-4 gap-y-1 sm:grid-cols-2 xl:grid-cols-3 tabular-nums">
          <li v-for="dependent in hardDependents" :key="dependent.item.id" class="flex">
            <span class="text-sm pr-1">{{ dependent.count }}x in</span>
            <artifact-name :artifact="dependent.item" :showTier="true" />
          </li>
        </ul>
      </div>
    </template>

    <template v-if="recursiveDependents.length > 0">
      <hr />
      <div class="px-4 py-4 sm:px-6 space-y-2">
        <div class="flex items-center space-x-1">
          <span class="text-sm font-medium text-gray-500">Recursive dependents</span>
          <info
            v-tippy="{
              content:
                'Dependents of dependents (i.e. &quot;used in&quot;), dependents of dependents of dependents, etc.',
            }"
          />
        </div>
        <ul class="grid grid-cols-1 gap-x-4 gap-y-1 sm:grid-cols-2 xl:grid-cols-3">
          <li v-for="dependent in recursiveDependents" :key="dependent.id">
            <artifact-name :artifact="dependent" :showTier="true" />
          </li>
        </ul>
      </div>
    </template>

    <hr />
    <div class="px-4 py-4 sm:px-6 space-y-1">
      <div>
        <span class="text-sm font-medium text-gray-500 mr-1">
          What do I get from consuming this item?
        </span>
        <a
          :href="`https://wasmegg.netlify.app/consumption-sheet/#${artifact.id}`"
          target="_blank"
          class="inline-flex items-center border-dashed border-b border-gray-700 text-sm whitespace-nowrap leading-tight space-x-0.5"
        >
          <span>Consumption sheet</span>
          <svg viewBox="0 0 20 20" fill="currentColor" class="h-4 w-4" data-html2canvas-ignore>
            <path
              d="M11 3a1 1 0 100 2h2.586l-6.293 6.293a1 1 0 101.414 1.414L15 6.414V9a1 1 0 102 0V4a1 1 0 00-1-1h-5z"
            />
            <path
              d="M5 5a2 2 0 00-2 2v8a2 2 0 002 2h8a2 2 0 002-2v-3a1 1 0 10-2 0v3H5V7h3a1 1 0 000-2H5z"
            />
          </svg>
        </a>
      </div>
      <div v-if="artifact.type === 'Stone' || artifact.type === 'Stone ingredient'">
        <span class="text-sm font-medium text-gray-500 mr-1">
          Consuming which items yields this
          {{ artifact.type === "Stone" ? "stone" : "stone fragment" }}?
        </span>
        <a
          :href="`https://wasmegg.netlify.app/consumption-sheet/#${artifact.id}-sources`"
          target="_blank"
          class="inline-flex items-center border-dashed border-b border-gray-700 text-sm whitespace-nowrap leading-tight space-x-0.5"
        >
          <span>Consumption sheet</span>
          <svg viewBox="0 0 20 20" fill="currentColor" class="h-4 w-4" data-html2canvas-ignore>
            <path
              d="M11 3a1 1 0 100 2h2.586l-6.293 6.293a1 1 0 101.414 1.414L15 6.414V9a1 1 0 102 0V4a1 1 0 00-1-1h-5z"
            />
            <path
              d="M5 5a2 2 0 00-2 2v8a2 2 0 002 2h8a2 2 0 002-2v-3a1 1 0 10-2 0v3H5V7h3a1 1 0 000-2H5z"
            />
          </svg>
        </a>
      </div>
    </div>
  </div>
</template>

<script>
import ArtifactName from "@/components/ArtifactName.vue";
import Info from "@/components/Info.vue";
import LootDataCredit from "@/components/LootDataCredit.vue";
import MissionName from "@/components/MissionName.vue";
import Share from "@/components/Share.vue";

import { iconURL, stringCmp } from "@/utils";

export default {
  components: {
    ArtifactName,
    Info,
    LootDataCredit,
    MissionName,
    Share,
  },

  props: {
    artifactId: String,
    missions: Array,
    artifacts: Array,
    lootTable: Object,
  },

  data() {
    const id2artifact = {};
    for (const artifact of this.artifacts) {
      if (artifact.afxRarity === 0) {
        id2artifact[artifact.id] = artifact;
      }
    }
    return {
      id2artifact,
    };
  },

  computed: {
    artifact() {
      for (const artifact of this.artifacts) {
        if (artifact.id === this.artifactId && artifact.afxRarity === 0) {
          return artifact;
        }
      }
      return undefined;
    },

    recursiveIngredients() {
      if (!this.artifact.craftable) {
        return [];
      }
      const queue = this.artifact.recipe.ingredients.map(it => this.id2artifact[it.id]);
      const directIngredients = new Set(queue.map(it => it.id));
      const seen = new Set();
      const ingredients = [];
      while (queue.length > 0) {
        const item = queue.shift();
        if (!directIngredients.has(item.id)) {
          if (seen.has(item.id)) {
            continue;
          }
          ingredients.push(item);
        }
        seen.add(item.id);
        if (!item.craftable) {
          continue;
        }
        queue.push(...item.recipe.ingredients.map(it => this.id2artifact[it.id]));
      }
      return ingredients.sort((it1, it2) => stringCmp(it1.sortKey, it2.sortKey));
    },

    obtainableMissions() {
      if (!this.artifact) {
        return [];
      }
      const missions = this.missions
        .filter(
          mission =>
            mission.minQuality <= this.artifact.quality &&
            this.artifact.quality <= mission.maxQuality
        )
        .map(mission => {
          const missionLoot = this.lootTable[mission.id];
          const missionCount = missionLoot.missionCount;
          const itemCount = missionLoot.items[this.artifactId];
          const itemsPerMission = missionCount === 0 ? 0 : itemCount.total / missionCount;
          const itemsPerDay = (itemsPerMission * 3600 * 24) / mission.durationSeconds;
          const precision = Math.min(itemCount.total.toString().length, 3);
          const notEnoughData = itemCount.total < 20;
          return {
            ...mission,
            missionCount,
            itemCount,
            itemsPerMission,
            itemsPerDay,
            precision,
            notEnoughData,
          };
        });
      return missions.sort((m1, m2) => m2.itemsPerDay - m1.itemsPerDay);
    },

    dependents() {
      return this.calculateDependents(this.artifact);
    },

    hardDependents() {
      const allDependents = [].concat(
        this.dependents.map(it => it.item),
        this.recursiveDependents
      );
      const hard = [];
      for (const dependent of allDependents) {
        if (dependent.hard_dependencies !== null) {
          for (const ingr of dependent.hard_dependencies) {
            if (ingr.id === this.artifact.id) {
              hard.push({
                item: dependent,
                count: ingr.count,
              });
            }
          }
        }
      }
      return hard.sort((it1, it2) => stringCmp(it1.item.sortKey, it2.item.sortKey));
    },

    recursiveDependents() {
      const queue = this.dependents.map(it => it.item);
      const directDependents = new Set(queue.map(it => it.id));
      const seen = new Set();
      const dependents = [];
      while (queue.length > 0) {
        const item = queue.shift();
        if (!directDependents.has(item.id)) {
          if (seen.has(item.id)) {
            continue;
          }
          dependents.push(item);
        }
        seen.add(item.id);
        queue.push(...this.calculateDependents(item).map(it => it.item));
      }
      return dependents.sort((it1, it2) => stringCmp(it1.sortKey, it2.sortKey));
    },
  },

  methods: {
    iconURL,

    rarityFgClass(rarity) {
      switch (rarity) {
        case 1:
          return "text-blue-500";
        case 2:
          return "text-purple-500";
        case 3:
          return "text-yellow-500";
        default:
          return "";
      }
    },

    rarityFgClassOnDarkBg(rarity) {
      switch (rarity) {
        case 1:
          return "text-blue-400";
        case 2:
          return "text-purple-400";
        case 3:
          return "text-yellow-400";
        default:
          return "";
      }
    },

    calculateDependents(item) {
      const dependents = [];
      for (const it of this.artifacts) {
        if (it.afxRarity !== 0 || !it.craftable) {
          continue;
        }
        for (const ingr of it.recipe.ingredients) {
          if (ingr.id === item.id) {
            dependents.push({
              item: it,
              count: ingr.count,
            });
            break;
          }
        }
      }
      return dependents;
    },

    // precision is optional.
    formatPercentage(x, total, precision) {
      const percentage = (x / total) * 100;
      if (!precision) {
        precision = Math.max(
          Math.min(x.toString().length, 3),
          Math.round(percentage).toString().length
        );
      }
      return `${percentage.toPrecision(precision)}%`;
    },

    formatToPrecision(x, precision) {
      const s = x.toPrecision(precision);
      // If the formatted string is a decimal without exponent, or one with a
      // negative exponent, return as is.
      if (s.match(/^\d+\.\d+$/) || s.includes("e-")) {
        return s;
      }
      // If the formatted string is an integer, or has a positive exponent,
      // convert it to non-scientific notation, and add a ~ in front to mark it
      // as an approximate.
      return "~" + parseFloat(s).toFixed();
    },

    formattedNumMissionsPerItem(count, missionCount) {
      if (count === 0) {
        return "\u221e";
      }
      const precision = Math.min(count.toString().length, 3);
      return this.formatToPrecision(missionCount / count, precision);
    },
  },
};
</script>
