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
                        class="h-3 w-3"
                        :src="iconURL('egginc/icon_afx_stone_add_large.png', 32)"
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

    <div class="px-4 pb-4 sm:px-6 text-xs text-gray-500">
      &dagger; Not available from missions.
    </div>
  </div>
</template>

<script>
import ArtifactName from "@/components/ArtifactName.vue";
import MissionName from "@/components/MissionName.vue";
import Info from "@/components/Info.vue";
import Share from "@/components/Share.vue";

import { iconURL, stringCmp } from "@/utils";

export default {
  components: {
    ArtifactName,
    MissionName,
    Info,
    Share,
  },

  props: {
    artifactId: String,
    missions: Array,
    artifacts: Array,
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
      return this.missions.filter(
        mission =>
          mission.minQuality <= this.artifact.quality && this.artifact.quality <= mission.maxQuality
      );
    },

    dependents() {
      return this.calculateDependents(this.artifact);
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
  },
};
</script>
