<template>
  <div class="max-w-xs mx-auto my-2 space-y-2">
    <div class="text-base text-center font-medium uppercase">Selected aggregate effects</div>

    <div class="flex flex-wrap itesm-center justify-center">
      <span class="flex whitespace-nowrap mr-2">
        <img :src="iconURL('egginc/egg_of_prophecy.png', 64)" class="inline h-5 w-5" />
        <span class="text-sm">{{ builds.config.prophecyEggs }}</span>
      </span>
      <span class="flex whitespace-nowrap">
        <img :src="iconURL('egginc/egg_soul.png', 64)" class="inline h-5 w-5" />
        <span class="text-sm">{{ formatEIValue(builds.config.soulEggs) }}</span>
      </span>
    </div>

    <table class="min-w-full rounded-md overflow-hidden">
      <tbody>
        <tr>
          <td class="px-4 py-1.5 text-sm text-left">
            <span
              v-tippy="{
                content: notes[0][1],
              }"
            >
              EB<note-sup :footnoteNumber="1" :showFootnote="showFootnotes" />
            </span>
          </td>
          <td v-if="buildValidities[0]" class="px-4 py-1.5 text-base text-right whitespace-nowrap">
            <span class="Value">{{ formatEIPercentage(earningBonus(...buildConfig(0))) }}</span>
            <span class="Bonus"
              >&nbsp;(&times;{{ formatFloat(earningBonusMultiplier(...buildConfig(0))) }})</span
            >
          </td>
          <td v-else class="px-4 py-1.5 text-base text-right text-red-500 whitespace-nowrap">
            &mdash;
          </td>
        </tr>

        <tr>
          <td class="px-4 py-1.5 text-sm text-left">
            <span
              v-tippy="{
                content: notes[1][1],
              }"
            >
              Earnings<note-sup :footnoteNumber="2" :showFootnote="showFootnotes" />
            </span>
          </td>
          <td
            v-if="buildValidities[0]"
            class="px-4 py-1.5 text-base text-right whitespace-nowrap Bonus"
          >
            &times;{{ formatFloat(earningsMultipler(...buildConfig(0))) }}
          </td>
          <td v-else class="px-4 py-1.5 text-base text-right text-red-500 whitespace-nowrap">
            &mdash;
          </td>
        </tr>

        <tr>
          <td class="px-4 py-1.5 text-sm text-left">
            <span
              v-tippy="{
                content: notes[2][1],
              }"
            >
              Max RCB<note-sup :footnoteNumber="3" :showFootnote="showFootnotes" />
            </span>
          </td>
          <td v-if="buildValidities[0]" class="px-4 py-1.5 text-base text-right whitespace-nowrap">
            <span class="Value">{{ maxRunningChickenBonus(...buildConfig(0)) }}</span>
            <span class="Bonus"
              >&nbsp;(&times;{{
                formatFloat(maxRunningChickenBonusMultiplier(...buildConfig(0)))
              }})</span
            >
          </td>
          <td v-else class="px-4 py-1.5 text-base text-right text-red-500 whitespace-nowrap">
            &mdash;
          </td>
        </tr>

        <tr>
          <td class="px-4 py-1.5 text-sm text-left">
            <span
              v-tippy="{
                content: notes[3][1],
              }"
            >
              Earnings w/ max RCB<note-sup :footnoteNumber="4" :showFootnote="showFootnotes" />
            </span>
          </td>
          <td
            v-if="buildValidities[0]"
            class="px-4 py-1.5 text-base text-right whitespace-nowrap Bonus"
          >
            &times;{{ formatFloat(earningsWithMaxRunningChickenBonusMultipler(...buildConfig(0))) }}
          </td>
          <td v-else class="px-4 py-1.5 text-base text-right text-red-500 whitespace-nowrap">
            &mdash;
          </td>
        </tr>

        <tr>
          <td class="px-4 py-1.5 text-sm text-left">
            <span
              v-tippy="{
                content: notes[4][1],
              }"
            >
              SE gain<note-sup :footnoteNumber="5" :showFootnote="showFootnotes" />
            </span>
          </td>
          <td
            v-if="buildValidities[0]"
            class="px-4 py-1.5 text-base text-right whitespace-nowrap Bonus"
          >
            &times;{{ formatFloat(soulEggsGainMultipler(...buildConfig(0))) }}
          </td>
          <td v-else class="px-4 py-1.5 text-base text-right text-red-500 whitespace-nowrap">
            &mdash;
          </td>
        </tr>

        <tr>
          <td class="px-4 py-1.5 text-sm text-left">
            <span
              v-tippy="{
                content: notes[5][1],
              }"
            >
              Research discount<note-sup :footnoteNumber="6" :showFootnote="showFootnotes" />
            </span>
          </td>
          <td
            v-if="buildValidities[0]"
            class="px-4 py-1.5 text-base text-right whitespace-nowrap Bonus"
          >
            {{ formatFloat(researchPriceDiscount(...buildConfig(0)) * 100) }}%
          </td>
          <td v-else class="px-4 py-1.5 text-base text-right text-red-500 whitespace-nowrap">
            &mdash;
          </td>
        </tr>

        <tr>
          <td class="px-4 py-1.5 text-sm text-left">
            <span
              v-tippy="{
                content: notes[6][1],
              }"
            >
              Max hab space<note-sup :footnoteNumber="7" :showFootnote="showFootnotes" />
            </span>
          </td>
          <td v-if="buildValidities[0]" class="px-4 py-1.5 text-base text-right whitespace-nowrap">
            <span class="Value">{{ maxHabSpace(...buildConfig(0)).toLocaleString("en-US") }}</span>
            <span class="Bonus"
              >&nbsp;(&times;{{ formatFloat(habSpaceMultiplier(...buildConfig(0))) }})</span
            >
          </td>
          <td v-else class="px-4 py-1.5 text-base text-right text-red-500 whitespace-nowrap">
            &mdash;
          </td>
        </tr>

        <tr>
          <td class="px-4 py-1.5 text-sm text-left">
            <span
              v-tippy="{
                content: notes[7][1],
              }"
            >
              Max IHR<note-sup :footnoteNumber="8" :showFootnote="showFootnotes" />
            </span>
          </td>
          <td v-if="buildValidities[0]" class="px-4 py-1.5 text-base text-right whitespace-nowrap">
            <span class="Value"
              >{{
                maxInternalHatcheryRatePerMinPerHab(...buildConfig(0)).toLocaleString("en-US")
              }}/min/hab</span
            >
            <span class="Bonus"
              >&nbsp;(&times;{{
                formatFloat(internalHatcheryRateMultiplier(...buildConfig(0)))
              }})</span
            >
          </td>
          <td v-else class="px-4 py-1.5 text-base text-right text-red-500 whitespace-nowrap">
            &mdash;
          </td>
        </tr>

        <tr>
          <td class="px-4 py-1.5 text-sm text-left">
            <span
              v-tippy="{
                content: notes[8][1],
              }"
            >
              Egg laying rate<note-sup :footnoteNumber="9" :showFootnote="showFootnotes" />
            </span>
          </td>
          <td
            v-if="buildValidities[0]"
            class="px-4 py-1.5 text-base text-right whitespace-nowrap Bonus"
          >
            &times;{{ formatFloat(layingRateMultiplier(...buildConfig(0))) }}
          </td>
          <td v-else class="px-4 py-1.5 text-base text-right text-red-500 whitespace-nowrap">
            &mdash;
          </td>
        </tr>

        <tr>
          <td class="px-4 py-1.5 text-sm text-left">
            <span
              v-tippy="{
                content: notes[9][1],
              }"
            >
              Max egg laying rate<note-sup :footnoteNumber="10" :showFootnote="showFootnotes" />
            </span>
          </td>
          <td
            v-if="buildValidities[0]"
            class="px-4 py-1.5 text-base text-right whitespace-nowrap Bonus"
          >
            <span class="Value"
              >{{ formatEIValue(maxHourlyLayingRate(...buildConfig(0))) }}/hr</span
            >
            <span class="Bonus"
              >&nbsp;(&times;{{ formatFloat(maxLayingRateMultiplier(...buildConfig(0))) }})</span
            >
          </td>
          <td v-else class="px-4 py-1.5 text-base text-right text-red-500 whitespace-nowrap">
            &mdash;
          </td>
        </tr>

        <tr>
          <td class="px-4 py-1.5 text-sm text-left">
            <span
              v-tippy="{
                content: notes[10][1],
              }"
            >
              Max shipping capacity<note-sup :footnoteNumber="11" :showFootnote="showFootnotes" />
            </span>
          </td>
          <td
            v-if="buildValidities[0]"
            class="px-4 py-1.5 text-base text-right whitespace-nowrap Bonus"
          >
            <span class="Value"
              >{{ formatEIValue(maxHourlyShippingCapacity(...buildConfig(0))) }}/hr</span
            >
            <span class="Bonus"
              >&nbsp;(&times;{{ formatFloat(shippingCapacityMultiplier(...buildConfig(0))) }})</span
            >
          </td>
          <td v-else class="px-4 py-1.5 text-base text-right text-red-500 whitespace-nowrap">
            &mdash;
          </td>
        </tr>
      </tbody>
    </table>
  </div>

  <div v-if="showFootnotes">
    <ol class="list-decimal list-inside">
      <li v-for="(note, index) in notes" :key="index" class="Note text-xs">
        <span class="text-gray-50">{{ note[0] }}:</span>
        {{ note[1] }}
      </li>
    </ol>
  </div>
</template>

<script>
import NoteSup from "@/components/NoteSup.vue";
import { Builds } from "@/lib/models";
import {
  earningBonus,
  earningBonusMultiplier,
  earningsMultipler,
  earningsWithMaxRunningChickenBonusMultipler,
  maxRunningChickenBonus,
  maxRunningChickenBonusMultiplier,
  soulEggsGainMultipler,
  researchPriceDiscount,
  maxHabSpace,
  habSpaceMultiplier,
  maxInternalHatcheryRatePerMinPerHab,
  internalHatcheryRateMultiplier,
  layingRateMultiplier,
  maxLayingRateMultiplier,
  maxHourlyLayingRate,
  shippingCapacityMultiplier,
  maxHourlyShippingCapacity,
} from "@/lib/effects/effects";
import { formatEIValue, formatEIPercentage, formatFloat } from "@/lib/utils/utils";

export default {
  components: {
    NoteSup,
  },

  props: {
    builds: {
      type: Builds,
      required: true,
    },
    showFootnotes: Boolean,
  },

  data() {
    return {
      notes: [
        ["EB", "Earning bonus, as shown on the prestige screen."],
        [
          "Earnings",
          "Aggregate effect on bock earning rate from earning bonus increase, egg value increase, and egg laying rate increase (not considering shipping-limited scenarios). Running chicken bonus is not taken into account here; see “Earnings w/ max RCB” instead.",
        ],
        ["Max RCB", "Max running chicken bonus."],
        [
          "Earnings w/ max RCB",
          "Increase in bock earning rate assuming the respective max running chicken bonus is attained with and without artifacts.",
        ],
        [
          "SE gain",
          "Soul egg earning rate multiplier from bock earning rate bonus (with max RCB) and soul egg collection rate bonus. The late game dampening exponent (0.21) is used; may not be accurate for early game players.",
        ],
        [
          "Research discount",
          "Only applies to common research. Compounds with research sale events: e.g. 60% discount compounded with a 65% research sale results in an aggregate 1−(1−60%)×(1−65%) = 86% discount.",
        ],
        ["Max hab space", "Hab space from four chicken universe habs."],
        [
          "Max IHR",
          "Internal hatchery rate as shown on the stats page, without taking internal hatchery calm into account.",
        ],
        ["Egg laying rate", "Per-chicken multiplier from egg laying rate bonuses."],
        ["Max egg laying rate", "Total hourly egg laying rate with full habs."],
        ["Max shipping capacity", "Actual egg shipping rate is capped at this value."],
      ],
    };
  },

  methods: {
    earningBonus,
    earningBonusMultiplier,
    earningsMultipler,
    earningsWithMaxRunningChickenBonusMultipler,
    maxRunningChickenBonus,
    maxRunningChickenBonusMultiplier,
    soulEggsGainMultipler,
    researchPriceDiscount,
    maxHabSpace,
    habSpaceMultiplier,
    maxInternalHatcheryRatePerMinPerHab,
    internalHatcheryRateMultiplier,
    layingRateMultiplier,
    maxLayingRateMultiplier,
    maxHourlyLayingRate,
    shippingCapacityMultiplier,
    maxHourlyShippingCapacity,

    formatEIValue,
    formatEIPercentage,
    formatFloat,

    buildConfig(i) {
      return [this.builds.builds[i], this.builds.config];
    },
  },

  computed: {
    buildValidities() {
      return this.builds.builds.map(build => !build.hasDuplicates());
    },
  },
};
</script>

<style scoped>
tr:nth-child(odd) {
  background-color: hsl(0, 0%, 20%);
}

tr:nth-child(even) {
  background-color: hsl(0, 0%, 22%);
}

.Value {
  color: #2d87ee;
}

.Bonus {
  color: #1e9c11;
}

.Note {
  color: #a6a6a6;
}
</style>
