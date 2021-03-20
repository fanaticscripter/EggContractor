<template>
  <table class="min-w-full divide-y divide-gray-200">
    <thead class="bg-gray-50">
      <tr>
        <template v-for="label in labels" :key="label.sortBy">
          <th
            v-if="label.name !== 'Offline' || displayOfflineColumn"
            @click="setSortBy(label.sortBy)"
            scope="col"
            class="px-6 py-2 text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer"
            :class="label.name === 'Player' ? 'text-left' : 'text-center'"
          >
            {{ label.name }}
            <!-- Use visibility for the arrow so that column widths don't change when sorting a different column. -->
            <span
              class="inline-block w-0 text-gray-400"
              :class="{ invisible: label.sortBy != sortBy }"
              >&nbsp;&#x25BC;</span
            >
          </th>
        </template>
      </tr>
    </thead>
    <tbody>
      <tr
        v-for="(member, index) in sortedMembers"
        :key="member.id"
        :class="index % 2 === 1 ? 'bg-gray-50' : 'bg-white'"
      >
        <td
          class="px-6 py-1 whitespace-nowrap text-sm text-gray-500 hover:text-gray-400 cursor-pointer"
          :class="{ 'CoopTable__member--snoozing': !member.isActive }"
          v-tippy="{ content: `${member.id} (click to copy)` }"
          @click="copy(member.id, `Copied ID of ${member.name}`)"
        >
          {{ member.name }}
        </td>
        <td class="px-6 py-1 whitespace-nowrap text-center text-sm text-gray-500">
          {{ member.eggsLaidStr }}
        </td>
        <td class="px-6 py-1 whitespace-nowrap text-center text-sm text-gray-500">
          {{ member.eggsPerHourStr }}
        </td>
        <td class="px-6 py-1 whitespace-nowrap text-center text-sm text-gray-500">
          {{ member.earningBonusPercentageStr }}
        </td>
        <td class="px-6 py-1 whitespace-nowrap text-center text-sm text-gray-500">
          {{ member.tokens }}
        </td>
        <td
          v-if="displayOfflineColumn"
          class="px-6 py-1 whitespace-nowrap text-center text-sm text-gray-500"
        >
          {{ member.offlineTimeStr }}
        </td>
      </tr>
    </tbody>
  </table>
  <transition name="fade">
    <div
      v-show="popupShow"
      class="fixed top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 px-2 py-1 rounded text-sm text-green-500 bg-white bg-opacity-80"
    >
      {{ popupMessage }}
    </div>
  </transition>
</template>

<script>
import { directive } from "vue-tippy";
import copyTextToClipboard from "copy-text-to-clipboard";
import { getSessionStorage, setSessionStorage } from "./utils";

export default {
  directives: {
    tippy: directive,
  },

  props: {
    contractId: String,
    code: String,
    members: Array,
  },

  data() {
    const labels = [
      { sortBy: "name", name: "Player" },
      { sortBy: "eggsLaid", name: "Laid" },
      { sortBy: "eggsPerHour", name: "Rate/hr" },
      { sortBy: "earningBonusPercentage", name: "EB%" },
      { sortBy: "tokens", name: "Tokens" },
      { sortBy: "offlineSeconds", name: "Offline" },
    ];
    const validSortBys = labels.map(l => l.sortBy);
    const sortBySessionStorageKey = `${this.contractId}:${this.code}_sortBy`;
    let sortBy = getSessionStorage(sortBySessionStorageKey);
    if (!validSortBys.includes(sortBy)) {
      sortBy = "eggsLaid";
    }
    return {
      labels,
      sortBy,
      sortBySessionStorageKey,
      popupTimeoutId: null,
      popupShow: false,
      popupMessage: "",
    };
  },

  computed: {
    displayOfflineColumn() {
      for (const m of this.members) {
        if (m.offlineTimeStr !== "") {
          return true;
        }
      }
      return false;
    },

    sortedMembers() {
      return [...this.members].sort((m1, m2) => {
        if (this.sortBy === "name") {
          return m1[this.sortBy].localeCompare(m2[this.sortBy]);
        } else {
          return m2[this.sortBy] - m1[this.sortBy];
        }
      });
    },
  },

  methods: {
    setSortBy(sortBy) {
      this.sortBy = sortBy;
      setSessionStorage(this.sortBySessionStorageKey, sortBy);
    },

    copy(s, msg) {
      copyTextToClipboard(s);
      this.popupMessage = msg;
      this.popupShow = true;
      if (this.popupTimeoutId !== null) {
        clearTimeout(this.popupTimeoutId);
      }
      this.popupTimeoutId = setTimeout(() => {
        this.popupShow = false;
      }, 3000);
    },
  },
};
</script>

<style scoped>
.fade-enter-active {
  transition: opacity 0.1s ease;
}

.fade-leave-active {
  transition: opacity 0.5s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
