<template>
  <div class="max-w-xs mx-auto my-4">
    <h3 class="my-2 text-center text-sm font-medium uppercase">Configurations</h3>

    <div class="mt-1 relative rounded-md shadow-sm">
      <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
        <img :src="iconURL('egginc/egg_of_prophecy.png', 64)" class="h-5 w-5" />
      </div>
      <integer-input id="prophecy_eggs" :min="0" :max="9999" v-model="conf.prophecyEggs" />
    </div>

    <div class="mt-1 relative rounded-md shadow-sm">
      <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
        <img :src="iconURL('egginc/egg_soul.png', 64)" class="h-5 w-5" />
      </div>
      <ei-value-with-unit-input
        id="soul_eggs"
        v-model:raw="conf.soulEggsInput"
        v-model:value="conf.soulEggs"
      />
    </div>

    <div class="mt-2 flex items-center justify-center">
      <input
        id="is_enlightenment"
        name="is_enlightenment"
        type="checkbox"
        class="h-4 w-4 bg-dark-20 text-blue-600 focus:ring-blue-500 focus:ring-offset-dark-30 rounded"
        v-model="conf.isEnlightenment"
      />
      <label for="is_enlightenment" class="ml-2 block text-sm">Enlightenment farm</label>
    </div>

    <div class="mt-4 flex justify-center">
      <div class="space-y-0.5">
        <h4 class="text-center text-sm uppercase">Active boost effects</h4>
        <div class="relative flex items-start">
          <input
            id="bird_feed_active"
            name="bird_feed_active"
            type="checkbox"
            class="h-4 w-4 bg-dark-20 text-blue-600 focus:ring-blue-500 focus:ring-offset-dark-30 rounded"
            v-model="conf.birdFeedActive"
          />
          <label for="bird_feed_active" class="ml-2 block text-sm">Bird feed (earnings)</label>
        </div>
        <div class="relative flex items-start">
          <input
            id="tachyon_prism_active"
            name="tachyon_prism_active"
            type="checkbox"
            class="h-4 w-4 bg-dark-20 text-blue-600 focus:ring-blue-500 focus:ring-offset-dark-30 rounded"
            v-model="conf.tachyonPrismActive"
          />
          <label for="tachyon_prism_active" class="ml-2 block text-sm"
            >Tachyon prism (internal hatchery)</label
          >
        </div>
        <div class="relative flex items-start">
          <input
            id="soul_beacon_active"
            name="soul_beacon_active"
            type="checkbox"
            class="h-4 w-4 bg-dark-20 text-blue-600 focus:ring-blue-500 focus:ring-offset-dark-30 rounded"
            v-model="conf.soulBeaconActive"
          />
          <label for="soul_beacon_active" class="ml-2 block text-sm">Soul beacon</label>
        </div>
        <div class="relative flex items-start">
          <input
            id="boost_beacon_active"
            name="boost_beacon_active"
            type="checkbox"
            class="h-4 w-4 bg-dark-20 text-blue-600 focus:ring-blue-500 focus:ring-offset-dark-30 rounded"
            v-model="conf.boostBeaconActive"
          />
          <label for="boost_beacon_active" class="ml-2 block text-sm">Boost beacon</label>
        </div>
      </div>
    </div>

    <div class="mt-4 flex justify-center">
      <div class="space-y-0.5">
        <h4 class="text-center text-sm uppercase">Other bonuses</h4>
        <div>
          <label
            for="tachyon_deflector_percentage"
            class="block text-center"
            v-tippy="{
              content:
                'This is the total tachyon deflector bonus from other players in your coop. Note that the value displayed in-game may be 1% less than the actual value due to a long standing floating-point representation bug in the game. You can always find the correct value from https://eicoop.netlify.app/.',
            }"
          >
            Tachyon deflector bonus<sup class="px-0.5">?</sup>
          </label>
          <div class="mt-1 relative rounded-md shadow-sm">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <img :src="iconURL('egginc/afx_tachyon_deflector_4.png', 64)" class="h-5 w-5" />
            </div>
            <integer-input
              id="tachyon_deflector_percentage"
              :min="0"
              :modelValue="round(conf.tachyonDeflectorBonus * 100)"
              @update:modelValue="value => (conf.tachyonDeflectorBonus = value / 100)"
              class="pr-4"
            />
            <div class="absolute inset-y-px right-0 pr-2 pt-2.5 pb-2 sm:text-sm text-gray-200">
              %
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { Config } from "@/lib/models";
import IntegerInput from "@/components/IntegerInput.vue";
import EiValueWithUnitInput from "@/components/EiValueWithUnitInput.vue";

export default {
  components: {
    IntegerInput,
    EiValueWithUnitInput,
  },

  props: {
    config: {
      type: Config,
      required: true,
    },
  },

  data() {
    return {
      conf: this.config,
    };
  },

  emits: ["update:config"],

  watch: {
    conf: {
      handler() {
        this.$emit("update:config", this.conf);
      },
      deep: true,
    },
  },

  methods: {
    round(x) {
      return Math.round(x);
    },
  },
};
</script>

<style scoped>
sup {
  color: #a6a6a6;
}
</style>
