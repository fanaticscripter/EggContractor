<template>
  <form
    class="mx-4 sm:mx-auto sm:max-w-xs sm:w-full mt-2 mb-4 space-y-1"
    @submit.prevent="submit(playerId)"
  >
    <div>
      <label for="email" class="sr-only">Player ID</label>
      <input
        type="text"
        name="playerId"
        id="playerId"
        v-model.trim="playerId"
        class="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
        placeholder="Player ID"
      />
      <div class="text-center">
        <span
          class="mt-2 inline-flex items-center space-x-1"
          v-tippy="{
            content:
              'The ID asked for here is the unique ID used by Egg, Inc.\'s server to identify your account. You can find it in game screen -> nine dots menu -> Settings -> Privacy & Data, at the very bottom. It should look like EI1234567890123456. Your old game services ID prior to the Artifact Update does not work here. Also note that the ID is case-sensitive.',
          }"
        >
          <base-info />
          <span class="text-xs text-gray-500">Where do I find my ID?</span>
        </span>
      </div>
    </div>
    <div>
      <button
        type="submit"
        class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
        :class="{ 'cursor-not-allowed': !submittable }"
        :disabled="!submittable"
      >
        Load Player Data
      </button>
    </div>
  </form>
</template>

<script lang="ts">
import { computed, defineComponent, PropType, ref, toRefs } from "vue";
import BaseInfo from "./BaseInfo.vue";

export default defineComponent({
  components: { BaseInfo },
  props: {
    playerIdPreload: {
      type: String,
      default: "",
    },
    submit: {
      type: Function as PropType<(playerId: string) => void>,
      required: true,
    },
  },
  setup(props) {
    const { playerIdPreload } = toRefs(props);
    const playerId = ref(playerIdPreload.value);
    const submittable = computed(() => {
      return playerId.value !== "";
    });
    return {
      playerId,
      submittable,
    };
  },
});
</script>
