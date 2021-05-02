<!-- This component is adapted from SimulatorItemSelect.vue from loot-simulator. -->

<template>
  <div class="flex flex-wrap items-center gap-2">
    <div class="flex-grow relative max-w-sm">
      <div class="relative rounded-md shadow-sm">
        <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
          <img
            v-if="selected && searchFilter === selected.display"
            :src="iconURL(selected.iconPath, 64)"
            class="flex-shrink-0 h-6 w-6"
          />
          <!-- Heroicon name: solid/question-mark-circle -->
          <svg
            v-else
            xmlns="http://www.w3.org/2000/svg"
            class="flex-shrink-0 h-6 w-6 text-gray-400"
            viewBox="-2 -2 24 24"
            fill="currentColor"
          >
            <path
              fill-rule="evenodd"
              d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-8-3a1 1 0 00-.867.5 1 1 0 11-1.731-1A3 3 0 0113 8a3.001 3.001 0 01-2 2.83V11a1 1 0 11-2 0v-1a1 1 0 011-1 1 1 0 100-2zm0 8a1 1 0 100-2 1 1 0 000 2z"
              clip-rule="evenodd"
            />
          </svg>
        </div>
        <input
          ref="selectButton"
          type="text"
          class="focus:ring-blue-500 focus:border-blue-500 block w-full pl-11 pr-12 pt-2.5 pb-2 sm:text-sm bg-dark-20 rounded-md placeholder-white"
          :class="
            selected && searchFilter === selected.display && selected.afx_rarity > 0
              ? selected.rarity
              : null
          "
          spellcheck="false"
          :placeholder="
            type === 'artifact' ? '-- Select artifact (type to filter) --' : '-- Select stone --'
          "
          v-model="searchFilter"
          :disabled="!!disabled"
          @focus="openDropdown"
          @blur="closeDropdown"
          @keydown="handleKeydown"
        />
        <div class="absolute inset-y-0 right-0 pr-2 flex items-center pointer-events-none">
          <!-- Heroicon name: solid/selector -->
          <svg
            class="h-5 w-5 text-gray-400"
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 20 20"
            fill="currentColor"
            aria-hidden="true"
          >
            <path
              fill-rule="evenodd"
              d="M10 3a1 1 0 01.707.293l3 3a1 1 0 01-1.414 1.414L10 5.414 7.707 7.707a1 1 0 01-1.414-1.414l3-3A1 1 0 0110 3zm-3.707 9.293a1 1 0 011.414 0L10 14.586l2.293-2.293a1 1 0 011.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z"
              clip-rule="evenodd"
            />
          </svg>
        </div>
        <div
          v-if="selected"
          class="absolute inset-y-0 right-7 flex items-center cursor-pointer"
          @click="clear"
        >
          <!-- Heroicon name: solid/x -->
          <svg
            class="h-4 w-4 text-gray-400"
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 20 20"
            fill="currentColor"
            aria-hidden="true"
          >
            <path
              fill-rule="evenodd"
              d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
              clip-rule="evenodd"
            />
          </svg>
        </div>
      </div>

      <ul
        ref="dropdownList"
        v-show="open"
        class="absolute mt-1 w-full bg-dark-20 shadow-lg rounded-md py-1 text-base border border-gray-500 ring-1 ring-black ring-opacity-5 overflow-auto focus:outline-none sm:text-sm z-10"
        :style="{ maxHeight: '21.5rem' }"
        tabindex="-1"
      >
        <li
          v-for="item in filteredItems"
          :key="item.id"
          class="cursor-default select-none relative py-0.5 pl-3 pr-9"
          :class="item.id === active?.id ? 'bg-blue-600' : null"
          @mousedown="selectItem(item)"
          @mouseenter="activateItem(item)"
          @mouseleave="deactivateItem()"
        >
          <div class="flex items-center">
            <img :src="iconURL(item.iconPath, 64)" class="flex-shrink-0 h-6 w-6" />
            <span
              class="ml-2 block truncate mt-1"
              :class="[
                item.id === selected?.id ? 'font-semibold' : 'font-normal',
                item.afx_rarity > 0 ? item.rarity : null,
              ]"
            >
              {{ item.display }}
            </span>
          </div>

          <span
            v-if="item.id === selected?.id"
            class="absolute inset-y-0 right-0 flex items-center pr-4"
            :class="item.id === active?.id ? 'text-white' : 'text-blue-600'"
          >
            <!-- Heroicon name: solid/check -->
            <svg
              class="h-5 w-5"
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 20 20"
              fill="currentColor"
              aria-hidden="true"
            >
              <path
                fill-rule="evenodd"
                d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                clip-rule="evenodd"
              />
            </svg>
          </span>
        </li>
        <div v-if="filteredItems.length === 0" class="py-0.5 pl-3 pr-9">
          <div class="flex items-center">
            <!-- Heroicon name: solid/exclamation-circle -->
            <svg
              class="flex-shrink-0 h-6 w-6 text-gray-400"
              viewBox="-2 -2 24 24"
              fill="currentColor"
            >
              <path
                fill-rule="evenodd"
                d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z"
                clip-rule="evenodd"
              />
            </svg>
            <span class="ml-2 block truncate">No match</span>
          </div>
        </div>
      </ul>
    </div>
  </div>
</template>

<script>
import { computed, defineComponent, nextTick, ref, toRefs } from "vue";
import scrollIntoView from "scroll-into-view-if-needed";

import {
  artifacts,
  artifactIdToArtifact,
  stones,
  stoneIdToStone,
  searchArtifacts,
  searchStones,
} from "@/lib/data";
import { iconURL } from "@/utils";

export default defineComponent({
  props: {
    // modelValue is the item ID.
    modelValue: {
      type: String,
      required: true,
    },
    type: {
      type: String,
      required: true,
      validator: value => ["artifact", "stone"].includes(value),
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  emits: {
    "update:modelValue": itemId => true,
  },
  setup(props, { emit }) {
    // Note that we do NOT react to type change.
    const { type } = props;
    const { modelValue } = toRefs(props);

    const items = type === "artifact" ? artifacts : stones;
    const itemIdToItem = type === "artifact" ? artifactIdToArtifact : stoneIdToStone;
    const searchItems = type === "artifact" ? searchArtifacts : searchStones;

    const selectButton = ref(null);
    const dropdownList = ref(null);

    const selected = computed(() => (modelValue.value ? itemIdToItem.get(modelValue.value) : null));
    const searchFilter = ref(selected.value ? selected.value.display : "");
    const filteredItems = computed(() =>
      searchFilter.value !== "" ? searchItems(searchFilter.value) : items
    );
    const active = ref(null);

    const dropdownListEntryIndex = id => {
      const entries = filteredItems.value;
      for (let i = 0; i < entries.length; i++) {
        if (entries[i].id === id) {
          return i;
        }
      }
      return -1;
    };
    const scrollDropdownListEntryIntoViewIfNeeded = (index, block /* ScrollLogicalPosition? */) => {
      block ||= "nearest";
      const node = dropdownList.value?.querySelector(`:scope > li:nth-child(${index + 1})`);
      if (node) {
        scrollIntoView(node, {
          scrollMode: "if-needed",
          block,
          inline: "nearest",
          boundary: dropdownList.value,
        });
      }
    };

    const open = ref(false);
    const openDropdown = () => {
      searchFilter.value = "";
      active.value = selected.value;
      open.value = true;
      nextTick(() => {
        if (active.value) {
          const index = dropdownListEntryIndex(active.value.id);
          if (index !== -1) {
            scrollDropdownListEntryIntoViewIfNeeded(index, "center");
          }
        } else {
          scrollDropdownListEntryIntoViewIfNeeded(0);
        }
      });
    };
    const closeDropdown = () => {
      searchFilter.value = selected.value !== null ? selected.value.display : "";
      open.value = false;
    };
    const selectItem = item => {
      emit("update:modelValue", item.id);
    };
    const activateItem = item => {
      active.value = item;
    };
    const deactivateItem = () => {
      active.value = null;
    };
    const handleKeydown = event => {
      switch (event.key) {
        case "Enter":
          event.preventDefault();
          // Do nothing if there are no matching entries.
          if (filteredItems.value.length === 0) {
            return;
          }
          // Select the active entry, or the first entry.
          const item = active.value ? active.value : filteredItems.value[0];
          if (item !== null) {
            selectItem(item);
          }
          // Wait for selected entry to be updated before closing the dropdown.
          nextTick(() => {
            closeDropdown();
            selectButton.value?.blur();
          });
          return;
        case "ArrowDown":
        case "ArrowUp":
          event.preventDefault();
          // Do nothing if there are no matching entries.
          if (filteredItems.value.length === 0) {
            return;
          }
          const entries = filteredItems.value;
          let currentIndex = active.value ? dropdownListEntryIndex(active.value.id) : -1;
          if (currentIndex === -1) {
            // No entry currently active.
            currentIndex = event.key === "ArrowDown" ? -1 : entries.length;
          }
          const newIndex =
            (((event.key === "ArrowDown" ? currentIndex + 1 : currentIndex - 1) % entries.length) +
              entries.length) %
            entries.length;
          active.value = entries[newIndex];
          scrollDropdownListEntryIntoViewIfNeeded(newIndex);
          return;
      }
    };

    const clear = () => {
      emit("update:modelValue", "");
      searchFilter.value = "";
    };

    return {
      selectButton,
      dropdownList,
      selected,
      searchFilter,
      filteredItems,
      active,
      open,
      openDropdown,
      closeDropdown,
      selectItem,
      activateItem,
      deactivateItem,
      handleKeydown,
      clear,
      iconURL,
    };
  },
});
</script>

<style scoped>
.Rare {
  color: hsl(209, 100%, 70%);
}

.Epic {
  color: hsl(300, 100%, 70%);
}

.Legendary {
  color: hsl(37, 100%, 70%);
}
</style>
