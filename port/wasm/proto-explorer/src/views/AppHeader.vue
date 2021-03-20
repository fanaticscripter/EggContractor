<template>
  <header class="shadow bg-white border-b border-gray-200">
    <h1 class="mx-4 my-4 text-center text-lg leading-6 font-medium text-gray-900">
      Proto explorer
    </h1>
    <div class="max-w-4xl w-full mx-auto px-4 xl:px-0 space-y-2">
      <div>
        <p class="text-sm font-medium text-gray-700">
          App version:
          <code class="text-xs font-mono">{{ appVersion }}</code>
          <span class="font-normal">
            (<a
              href="https://github.com/fanaticscripter/EggContractor/tree/master/misc/protobuf"
              target="_blank"
              class="hover:text-gray-500 border-b border-gray-500 border-dashed"
              >protobuf definitions</a
            >,
            <router-link
              :to="{ name: 'doc' }"
              target="_blank"
              class="hover:text-gray-500 border-b border-gray-500 border-dashed"
              >generated documentation</router-link
            >)
          </span>
        </p>
      </div>

      <div class="mt-2">
        <div class="sm:hidden mb-2">
          <label for="tabs" class="sr-only">Select a tab</label>
          <select
            id="tabs"
            name="tabs"
            class="block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md"
            v-model="selectedRoute"
          >
            <option
              v-for="tab in tabs"
              :key="tab.name"
              :value="tab.name"
              v-html="tab.display"
            ></option>
          </select>
        </div>
        <div class="hidden sm:block">
          <div class="">
            <nav class="-mb-px flex space-x-8" aria-label="Tabs">
              <router-link
                v-for="tab in tabs"
                :key="tab.name"
                :to="{ name: tab.name }"
                href="#"
                class="whitespace-nowrap py-2 border-b-2 font-medium text-sm"
                :class="
                  tab.name === route
                    ? 'border-blue-500 text-blue-600'
                    : 'border-transparent text-gray-700 hover:text-gray-900 hover:border-gray-300'
                "
                aria-current="page"
                v-html="tab.display"
              />
            </nav>
          </div>
        </div>
      </div>
    </div>
  </header>
</template>

<script>
import { ref, watch } from "vue";
import { useRouter } from "vue-router";
import { APP_VERSION } from "@/lib/lib";

export default {
  props: {
    route: {
      type: String,
      required: true,
    },
  },

  setup(props) {
    const appVersion = APP_VERSION;
    const tabs = [
      {
        name: "arbitrary_payload",
        display: "Arbitrary payload",
      },
      {
        name: "first_contact",
        display: "<code>/first_contact</code>",
      },
      {
        name: "get_periodicals",
        display: "<code>/get_periodicals</code>",
      },
      {
        name: "coop_status",
        display: "<code>/coop_status</code>",
      },
    ];
    const selectedRoute = ref(props.route);

    const router = useRouter();
    watch(selectedRoute, () => router.push({ name: selectedRoute.value }));

    return {
      appVersion,
      tabs,
      selectedRoute,
    };
  },
};
</script>
