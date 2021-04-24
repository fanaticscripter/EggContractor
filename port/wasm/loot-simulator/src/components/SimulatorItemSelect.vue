<template>
  <div class="flex flex-wrap items-center gap-2">
    <div
      class="flex-grow relative max-w-sm"
      :style="{ minWidth: 'min(calc(100vw - 2rem), 20rem)' }"
    >
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
          class="focus:ring-blue-500 focus:border-blue-500 block w-full pl-11 pr-10 sm:text-sm border-gray-300 rounded-md"
          spellcheck="false"
          placeholder="Select artifact (type to filter)"
          v-model="searchFilter"
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
      </div>

      <ul
        ref="dropdownList"
        v-show="open"
        class="absolute mt-1 w-full bg-white shadow-lg rounded-md py-1 text-base ring-1 ring-black ring-opacity-5 overflow-auto focus:outline-none sm:text-sm z-10"
        :style="{ maxHeight: '21.5rem' }"
        tabindex="-1"
      >
        <li
          v-for="item in filteredItems"
          :key="item.id"
          class="cursor-default select-none relative py-0.5 pl-3 pr-9"
          :class="item.id === active?.id ? 'text-white bg-blue-600' : 'text-gray-900'"
          @mousedown="selectItem(item)"
          @mouseenter="activateItem(item)"
          @mouseleave="deactivateItem()"
        >
          <div class="flex items-center">
            <img :src="iconURL(item.iconPath, 64)" class="flex-shrink-0 h-6 w-6" />
            <span
              class="ml-2 block truncate"
              :class="item.id === selected?.id ? 'font-semibold' : 'font-normal'"
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

    <div class="relative rounded-md shadow-sm self-stretch">
      <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
        <span class="text-gray-500 sm:text-sm">&times;</span>
      </div>
      <base-integer-input
        class="h-full pl-7 min-w-0"
        :style="{ maxWidth: '6rem' }"
        placeholder="count"
        :modelValue="modelValue.count"
        :min="1"
        @update:modelValue="updateCount"
      />
    </div>

    <div>
      <svg
        class="h-5 w-5 text-gray-400 hover:text-gray-500 cursor-pointer"
        viewBox="0 0 20 20"
        fill="currentColor"
        @click="deleteEntry"
      >
        <path
          fill-rule="evenodd"
          d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
          clip-rule="evenodd"
        />
      </svg>
    </div>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, nextTick, PropType, Ref, ref, toRefs } from "vue";
import scrollIntoView from "scroll-into-view-if-needed";

import { itemIdToItem, items, searchItems } from "@/data";
import { iconURL } from "@/utils";
import { Item, ItemId, ItemSelectSpec } from "@/types";
import BaseIntegerInput from "@/components/BaseIntegerInput.vue";

export default defineComponent({
  components: {
    BaseIntegerInput,
  },
  props: {
    modelValue: {
      type: Object as PropType<ItemSelectSpec>,
      required: true,
    },
  },
  emits: {
    "update:modelValue": (payload: ItemSelectSpec) => true,
    delete: () => true,
  },
  setup(props, { emit }) {
    const { modelValue } = toRefs(props);

    const selectButton: Ref<HTMLElement | null> = ref(null);
    const dropdownList: Ref<HTMLElement | null> = ref(null);

    const selected = computed(() =>
      modelValue.value.id !== null ? itemIdToItem.get(modelValue.value.id)! : null
    );
    const searchFilter = ref(selected.value ? selected.value.display : "");
    const filteredItems = computed(() =>
      searchFilter.value !== "" ? searchItems(searchFilter.value) : items
    );
    const active: Ref<Item | null> = ref(null);

    const updateCount = (count: number) => {
      emit("update:modelValue", {
        ...modelValue.value,
        count,
      });
    };

    const dropdownListEntryIndex = (id: ItemId): number => {
      const entries = filteredItems.value;
      for (let i = 0; i < entries.length; i++) {
        if (entries[i].id === id) {
          return i;
        }
      }
      return -1;
    };
    const scrollDropdownListEntryIntoViewIfNeeded = (
      index: number,
      block: ScrollLogicalPosition = "nearest"
    ) => {
      const node = dropdownList.value?.querySelector(`:scope > li:nth-child(${index + 1})`);
      if (node) {
        scrollIntoView(node, {
          scrollMode: "if-needed",
          block,
          inline: "nearest",
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
    const selectItem = (item: Item) => {
      emit("update:modelValue", {
        ...modelValue.value,
        id: item.id,
      });
    };
    const activateItem = (item: Item) => {
      active.value = item;
    };
    const deactivateItem = () => {
      active.value = null;
    };
    const handleKeydown = (event: KeyboardEvent) => {
      switch (event.key) {
        case "Enter":
          event.preventDefault();
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

    const deleteEntry = () => {
      emit("delete");
    };

    return {
      selectButton,
      dropdownList,
      selected,
      searchFilter,
      filteredItems,
      active,
      updateCount,
      open,
      openDropdown,
      closeDropdown,
      selectItem,
      activateItem,
      deactivateItem,
      handleKeydown,
      deleteEntry,
      iconURL,
    };
  },
});
</script>
