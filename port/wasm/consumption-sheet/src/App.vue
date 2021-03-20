<template>
  <div class="flex-1 max-w-10xl w-full mx-auto px-4 my-4">
    <section id="index" class="mt-4 mb-6 hide-in-screenshot-mode">
      <h2 class="my-2">Index</h2>

      <p>
        <a href="#consumption" class="text-sm text-gray-700 hover:text-gray-500 underline"
          >What do I get from consuming ...?</a
        >
      </p>
      <div
        class="my-2 grid grid-cols-2 gap-1 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 3xl:grid-cols-8"
      >
        <a
          v-for="family in data.families"
          :key="family.id"
          :href="`#${family.id}`"
          class="text-xs text-gray-500 hover:text-gray-400 underline"
          :class="family.type"
        >
          {{ family.name }}
        </a>
      </div>

      <p>
        <a href="#sources" class="text-sm text-gray-700 hover:text-gray-500 underline"
          >Consuming which items yields ... stone?</a
        >
      </p>
      <div
        class="my-2 grid grid-cols-2 gap-1 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 3xl:grid-cols-8"
      >
        <template v-for="family in data.families" :key="family.id">
          <a
            v-if="family.type === 'Stone'"
            :href="`#${family.id}-sources`"
            class="text-xs text-gray-500 hover:text-gray-400 underline"
            :class="family.type"
          >
            {{ family.name }}
          </a>
        </template>
      </div>
    </section>

    <section id="consumption" class="mt-4 mb-6">
      <h2 class="my-2">What do I get from consuming ...?</h2>

      <div class="text-xs text-gray-500 space-y-1">
        <p>
          The following dataset contains the outcomes of consuming each item (different rarities are
          condisered distinct) 100 times. The consumption outcome of an individual item is
          considered deterministic if all trial runs produced the exact same rewards; otherwise, the
          item is marked as nondeterministic
          <svg class="inline h-3" viewBox="0 0 640 512">
            <path
              fill="currentColor"
              d="M592 192H473.26c12.69 29.59 7.12 65.2-17 89.32L320 417.58V464c0 26.51 21.49 48 48 48h224c26.51 0 48-21.49 48-48V240c0-26.51-21.49-48-48-48zM480 376c-13.25 0-24-10.75-24-24 0-13.26 10.75-24 24-24s24 10.74 24 24c0 13.25-10.75 24-24 24zm-46.37-186.7L258.7 14.37c-19.16-19.16-50.23-19.16-69.39 0L14.37 189.3c-19.16 19.16-19.16 50.23 0 69.39L189.3 433.63c19.16 19.16 50.23 19.16 69.39 0L433.63 258.7c19.16-19.17 19.16-50.24 0-69.4zM96 248c-13.25 0-24-10.75-24-24 0-13.26 10.75-24 24-24s24 10.74 24 24c0 13.25-10.75 24-24 24zm128 128c-13.25 0-24-10.75-24-24 0-13.26 10.75-24 24-24s24 10.74 24 24c0 13.25-10.75 24-24 24zm0-128c-13.25 0-24-10.75-24-24 0-13.26 10.75-24 24-24s24 10.74 24 24c0 13.25-10.75 24-24 24zm0-128c-13.25 0-24-10.75-24-24 0-13.26 10.75-24 24-24s24 10.74 24 24c0 13.25-10.75 24-24 24zm128 128c-13.25 0-24-10.75-24-24 0-13.26 10.75-24 24-24s24 10.74 24 24c0 13.25-10.75 24-24 24z"
            /></svg
          >, and the displayed number of each byproduct is the expectation value computed from all
          trial runs. The exact output from each trial run can be viewed by clicking on the "expand"
          button
          <svg viewBox="0 0 192 512" class="inline h-4 text-green-500">
            <path
              fill="currentColor"
              d="M192 127.338v257.324c0 17.818-21.543 26.741-34.142 14.142L29.196 270.142c-7.81-7.81-7.81-20.474 0-28.284l128.662-128.662c12.599-12.6 34.142-3.676 34.142 14.142z"
            /></svg
          >.
        </p>
        <p>
          Note that due to the limited sample size, certain expectation values may be highly
          inaccurate, certain byproducts may have been missed altogether, and certain
          nondeterministic items may have been marked as deterministic.
        </p>
      </div>

      <div
        class="my-4 grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 3xl:grid-cols-5"
      >
        <div
          v-for="family in data.families"
          :key="family.id"
          :id="family.id"
          class="-mx-4 sm:mx-0 bg-gray-50 overflow-hidden sm:rounded-lg sm:shadow-md ScrollTarget"
          :class="family.type"
        >
          <div class="bg-gray-100 px-6 py-3 border-b border-gray-200">
            <div class="-ml-4 -mt-2 flex items-center justify-between flex-wrap sm:flex-nowrap">
              <a href="#index">
                <div class="ml-4 mt-2">{{ family.name }}</div>
              </a>
              <div class="ml-4 mt-2 flex-shrink-0"></div>
            </div>
          </div>

          <ul class="py-2" :class="family.type === 'Artifact' ? 'divide-y' : 'space-y-2'">
            <li
              v-for="tier in family.tiers"
              :key="tier.id"
              class="px-4"
              :class="family.type === 'Artifact' ? 'py-2' : null"
            >
              <div :id="tier.id" class="flex space-x-2 ScrollTarget">
                <div class="flex-shrink-0" v-tippy="{ content: tier.name }">
                  <img class="h-10 w-10" :src="iconURL(`egginc/${tier.icon_filename}`, 128)" />
                </div>

                <ul v-if="family.type === 'Artifact'" class="space-y-1">
                  <li v-for="rarity in tier.rarities" :key="rarity.item.afx_rarity">
                    <consumption-outcome :outcome="rarity" />
                  </li>
                </ul>

                <div v-else class="flex items-center">
                  <consumption-outcome :outcome="tier.rarities[0]" :hideRarity="true" />
                </div>
              </div>
            </li>
          </ul>
        </div>
      </div>
    </section>

    <hr />

    <section id="sources" class="mt-4 mb-6">
      <h2 class="my-4">Consuming which items yields ... stone?</h2>

      <div class="text-xs text-gray-500 space-y-1">
        <p>
          This section is the reverse of the previous section, showing a compilation of all items
          that could yield an individual stone (fragment) upon consumption, along with the expected
          yields, subject to the same limitations in accuracy. Item rarities are color coded:
          <span class="inline-block relative top-0.5 h-3 w-3 rounded bg-common"></span> for common,
          <span class="inline-block relative top-0.5 h-3 w-3 rounded bg-rare"></span> for rare,
          <span class="inline-block relative top-0.5 h-3 w-3 rounded bg-epic"></span> for epic, and
          <span class="inline-block relative top-0.5 h-3 w-3 rounded bg-legendary"></span> for
          legendary. Nondeterministic items are marked with a dot in the corner:
          <span class="inline-block relative top-0.5 h-3 w-3 rounded bg-legendary">
            <span
              class="absolute -top-0.5 -right-0.5 block h-1.5 w-1.5 rounded-full ring-1 ring-white bg-green-400"
            ></span> </span
          >.
        </p>
      </div>

      <div
        class="my-4 grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 3xl:grid-cols-5"
      >
        <template v-for="family in data.families" :key="family.id">
          <div
            v-if="family.type === 'Stone'"
            :id="`${family.id}-sources`"
            class="-mx-4 sm:mx-0 bg-gray-50 overflow-hidden sm:rounded-lg sm:shadow-md ScrollTarget"
          >
            <div class="bg-gray-100 px-6 py-3 border-b border-gray-200">
              <div class="-ml-4 -mt-2 flex items-center justify-between flex-wrap sm:flex-nowrap">
                <a href="#index">
                  <div class="ml-4 mt-2">{{ family.name }}</div>
                </a>
                <div class="ml-4 mt-2 flex-shrink-0"></div>
              </div>
            </div>

            <ul class="py-2 divide-y">
              <li v-for="tier in family.tiers" :key="tier.id" class="px-4 py-2">
                <div :id="`${tier.id}-sources`" class="flex space-x-2 ScrollTarget">
                  <div class="flex-shrink-0" v-tippy="{ content: tier.name }">
                    <img class="h-10 w-10" :src="iconURL(`egginc/${tier.icon_filename}`, 128)" />
                  </div>

                  <sources :sources="tier.sources" />
                </div>
              </li>
            </ul>
          </div>
        </template>
      </div>
    </section>
  </div>
</template>

<script>
import ConsumptionOutcome from "./ConsumptionOutcome.vue";
import Sources from "./Sources.vue";

export default {
  components: {
    ConsumptionOutcome,
    Sources,
  },

  props: {
    data: Object,
  },
};
</script>

<style scoped>
.Artifact + .Stone,
.Stone + .Ingredient {
  grid-column-start: 1;
}

.ScrollTarget {
  scroll-margin-top: 5rem;
}
</style>
