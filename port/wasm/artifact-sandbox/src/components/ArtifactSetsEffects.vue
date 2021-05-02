<template>
  <div class="max-w-xs mx-auto my-2 space-y-2">
    <div class="text-base text-center font-medium uppercase">Selected aggregate effects</div>

    <div class="flex flex-wrap items-center justify-center">
      <span class="flex whitespace-nowrap mr-2">
        <img :src="iconURL('egginc/egg_of_prophecy.png', 64)" class="inline h-5 w-5" />
        <span class="text-sm">{{ builds.config.prophecyEggs }}</span>
      </span>
      <span class="flex whitespace-nowrap">
        <img :src="iconURL('egginc/egg_soul.png', 64)" class="inline h-5 w-5" />
        <span class="text-sm">{{ formatEIValue(builds.config.soulEggs) }}</span>
      </span>
    </div>

    <div v-if="builds.config.anyBoostActive()" class="flex flex-col items-center justify-center">
      <div v-if="builds.config.birdFeedActive" class="flex whitespace-nowrap">
        <img
          :src="iconURL('egginc/b_icon_jimbos_orange_big.png', 64)"
          class="inline h-5 w-5 mr-1"
        />
        <span class="text-sm">Bird feed active</span>
      </div>
      <div v-if="builds.config.tachyonPrismActive" class="flex whitespace-nowrap">
        <img
          :src="iconURL('egginc/b_icon_tachyon_prism_orange_big.png', 64)"
          class="inline h-5 w-5 mr-1"
        />
        <span class="text-sm">Tachyon prism active</span>
      </div>
      <div v-if="builds.config.soulBeaconActive" class="flex whitespace-nowrap">
        <img
          :src="iconURL('egginc/b_icon_soul_beacon_orange.png', 64)"
          class="inline h-5 w-5 mr-1"
        />
        <span class="text-sm">Soul beacon active</span>
      </div>
      <div v-if="builds.config.boostBeaconActive" class="flex whitespace-nowrap">
        <img
          :src="iconURL('egginc/b_icon_boost_beacon_orange.png', 64)"
          class="inline h-5 w-5 mr-1"
        />
        <span class="text-sm">Boost beacon active</span>
      </div>
    </div>

    <div
      v-if="builds.config.tachyonDeflectorBonus > 0"
      class="flex flex-wrap items-center justify-center"
    >
      <span class="flex whitespace-nowrap">
        <img :src="iconURL('egginc/afx_tachyon_deflector_4.png', 64)" class="inline h-5 w-5 mr-1" />
        <span class="text-sm">
          Tachyon deflector +{{ formatEIPercentage(builds.config.tachyonDeflectorBonus) }}
        </span>
      </span>
    </div>

    <!-- :set is not a Vue feature, it's a trick to get a scoped temporary
    variable so that we don't need to keep track of footnote numbers manually.
    -->
    <table class="min-w-full rounded-md overflow-hidden" :set="(footnoteNumber = 1)">
      <tbody>
        <tr>
          <td class="px-4 py-1.5 text-sm text-left">
            <effect-with-note
              :noteList="notes"
              :number="footnoteNumber++"
              :showFootnote="showFootnotes"
              >EB</effect-with-note
            >
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
            <effect-with-note
              :noteList="notes"
              :number="footnoteNumber++"
              :showFootnote="showFootnotes"
              >Role</effect-with-note
            >
          </td>
          <td
            v-if="buildValidities[0]"
            class="px-4 py-1.5 text-base text-right whitespace-nowrap"
            :set="(role = earningBonusToFarmerRole(earningBonus(...buildConfig(0))))"
          >
            <span :style="{ color: role.color }">
              {{ role.name }}
            </span>
          </td>
          <td v-else class="px-4 py-1.5 text-base text-right text-red-500 whitespace-nowrap">
            &mdash;
          </td>
        </tr>

        <tr>
          <td class="px-4 py-1.5 text-sm text-left">
            <effect-with-note
              :noteList="notes"
              :number="footnoteNumber++"
              :showFootnote="showFootnotes"
              >Earnings</effect-with-note
            >
          </td>
          <td
            v-if="buildValidities[0]"
            class="px-4 py-1.5 text-base text-right whitespace-nowrap Bonus"
          >
            &times;{{ formatFloat(earningsMultiplier(...buildConfig(0))) }}
          </td>
          <td v-else class="px-4 py-1.5 text-base text-right text-red-500 whitespace-nowrap">
            &mdash;
          </td>
        </tr>

        <tr>
          <td class="px-4 py-1.5 text-sm text-left">
            <effect-with-note
              :noteList="notes"
              :number="footnoteNumber++"
              :showFootnote="showFootnotes"
              >Max RCB</effect-with-note
            >
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
            <effect-with-note
              :noteList="notes"
              :number="footnoteNumber++"
              :showFootnote="showFootnotes"
              >Earnings w/ Max RCB</effect-with-note
            >
          </td>
          <td
            v-if="buildValidities[0]"
            class="px-4 py-1.5 text-base text-right whitespace-nowrap Bonus"
          >
            &times;{{
              formatFloat(earningsWithMaxRunningChickenBonusMultiplier(...buildConfig(0)))
            }}
          </td>
          <td v-else class="px-4 py-1.5 text-base text-right text-red-500 whitespace-nowrap">
            &mdash;
          </td>
        </tr>

        <tr>
          <td class="px-4 py-1.5 text-sm text-left">
            <effect-with-note
              :noteList="notes"
              :number="footnoteNumber++"
              :showFootnote="showFootnotes"
              :dagger="true"
              >SE gain</effect-with-note
            >
          </td>
          <td
            v-if="buildValidities[0]"
            class="px-4 py-1.5 text-base text-right whitespace-nowrap Bonus"
          >
            &times;{{ formatFloat(soulEggsGainMultiplier(...buildConfig(0))) }}
          </td>
          <td v-else class="px-4 py-1.5 text-base text-right text-red-500 whitespace-nowrap">
            &mdash;
          </td>
        </tr>

        <tr>
          <td class="px-4 py-1.5 text-sm text-left">
            <effect-with-note
              :noteList="notes"
              :number="footnoteNumber++"
              :showFootnote="showFootnotes"
              :dagger="true"
              >SE gain w/<br />empty habs start</effect-with-note
            >
          </td>
          <td
            v-if="buildValidities[0]"
            class="px-4 py-1.5 text-base text-right whitespace-nowrap Bonus"
          >
            &times;{{ formatFloat(soulEggsGainWithEmptyHabsStartMultiplier(...buildConfig(0))) }}
          </td>
          <td v-else class="px-4 py-1.5 text-base text-right text-red-500 whitespace-nowrap">
            &mdash;
          </td>
        </tr>

        <tr>
          <td class="px-4 py-1.5 text-sm text-left">
            <effect-with-note
              :noteList="notes"
              :number="footnoteNumber++"
              :showFootnote="showFootnotes"
              >Boost duration</effect-with-note
            >
          </td>
          <td
            v-if="buildValidities[0]"
            class="px-4 py-1.5 text-base text-right whitespace-nowrap Bonus"
          >
            &times;{{ formatFloat(boostDurationMultiplier(...buildConfig(0))) }}
          </td>
          <td v-else class="px-4 py-1.5 text-base text-right text-red-500 whitespace-nowrap">
            &mdash;
          </td>
        </tr>

        <tr>
          <td class="px-4 py-1.5 text-sm text-left">
            <effect-with-note
              :noteList="notes"
              :number="footnoteNumber++"
              :showFootnote="showFootnotes"
              >Research discount</effect-with-note
            >
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
            <effect-with-note
              :noteList="notes"
              :number="footnoteNumber++"
              :showFootnote="showFootnotes"
              >Max hab space</effect-with-note
            >
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
            <effect-with-note
              :noteList="notes"
              :number="footnoteNumber++"
              :showFootnote="showFootnotes"
              >Max IHR</effect-with-note
            >
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
            <effect-with-note
              :noteList="notes"
              :number="footnoteNumber++"
              :showFootnote="showFootnotes"
              >Egg laying rate</effect-with-note
            >
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
            <effect-with-note
              :noteList="notes"
              :number="footnoteNumber++"
              :showFootnote="showFootnotes"
              >Max egg laying rate</effect-with-note
            >
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
            <effect-with-note
              :noteList="notes"
              :number="footnoteNumber++"
              :showFootnote="showFootnotes"
              >Max shipping capacity</effect-with-note
            >
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

  <div class="mt-2 DaggerNote text-sm leading-tight">
    <sup>&dagger;</sup> In order to maximize SE gain from boosted prestiges, you should optimize the
    <span class="DaggerNote--highlight">“SE gain”</span> or
    <span class="DaggerNote--highlight">“SE gain w/ empty habs start”</span> stat. If you start your
    boosts with preloaded, almost full habs, artifact effects over your SE gain is reflected by the
    <span class="DaggerNote--highlight">“SE gain”</span> stat, which you should attempt to maximize;
    if you start your boosts (including at least one tachyon prism) with empty habs, commonly seen
    in multi-prestige or all-in-one single-prestige strategies, you should instead maximize
    <span class="DaggerNote--highlight">“SE gain w/ empty habs start”</span>. Don't forget to
    configure <span class="uppercase">active boost effects</span> when optimizing aforementioned
    stats.
  </div>

  <div v-if="showFootnotes" class="mt-2">
    <ol class="list-decimal list-inside">
      <li v-for="(note, index) in notes" :key="index" class="Note text-xs">
        <span class="text-gray-50">{{ note[0] }}:</span>
        {{ note[1] }}
      </li>
    </ol>
  </div>
</template>

<script>
import EffectWithNote from "./EffectWithNote.vue";
import { Builds } from "@/lib/models";
import {
  earningBonus,
  earningBonusMultiplier,
  earningsMultiplier,
  earningsWithMaxRunningChickenBonusMultiplier,
  maxRunningChickenBonus,
  maxRunningChickenBonusMultiplier,
  soulEggsGainMultiplier,
  soulEggsGainWithEmptyHabsStartMultiplier,
  boostDurationMultiplier,
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
import { earningBonusToFarmerRole } from "@/lib/role";
import { formatEIValue, formatEIPercentage, formatFloat } from "@/lib/utils/utils";

export default {
  components: {
    EffectWithNote,
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
          "Role",
          "This is the role (rank) corresponding to the earning bonus, as used in the Egg, Inc. Discord.",
        ],
        [
          "Earnings",
          "Aggregate effect on bock earning rate from earning bonus increase, egg value increase, and egg laying rate increase (not considering shipping-limited scenarios). Running chicken bonus is not taken into account here; see “Earnings w/ max RCB” instead. Indirect bonus from boosted chicken population growth is not included.",
        ],
        ["Max RCB", "Max running chicken bonus."],
        [
          "Earnings w/ max RCB",
          "Increase in bock earning rate assuming the respective max running chicken bonus is attained with and without artifacts.",
        ],
        [
          "SE gain",
          "Soul egg earning rate multiplier from bock earning rate bonus (with max RCB) and soul egg collection rate bonus. The late game dampening exponent (0.21) is used; may not be accurate for early game players. Indirect bonus from boosted chicken population growth is not included; see “SE gain w/ empty habs start” instead.",
        ],
        [
          "SE gain w/ empty habs start",
          "Same as “SE gain” except for taking into account the indirect earnings bonus from faster chicken population growth when there is a monocle-boosted tachyon prism active. Assumes the tachyon prism is activated at zero population, and population never hits the hab space cap; otherwise, the actual effect is between this stat and “SE gain”.",
        ],
        [
          "Boost duration",
          "Affects the duration of any boost activated while this artifact set is equipped. Artifact changes after activation have no effect on the duration.",
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
        [
          "Egg laying rate",
          "Per-chicken multiplier from egg laying rate bonuses. Tachyon deflector bonus included.",
        ],
        [
          "Max egg laying rate",
          "Total hourly egg laying rate with full habs. Tachyon deflector bonus included.",
        ],
        ["Max shipping capacity", "Actual egg shipping rate is capped at this value."],
      ],
    };
  },

  methods: {
    earningBonus,
    earningBonusMultiplier,
    earningBonusToFarmerRole,
    earningsMultiplier,
    earningsWithMaxRunningChickenBonusMultiplier,
    maxRunningChickenBonus,
    maxRunningChickenBonusMultiplier,
    soulEggsGainMultiplier,
    soulEggsGainWithEmptyHabsStartMultiplier,
    boostDurationMultiplier,
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

.DaggerNote {
  color: #b38c00;
}

.DaggerNote--highlight {
  color: #ffc601;
}
</style>
