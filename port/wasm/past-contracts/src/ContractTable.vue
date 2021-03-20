<template>
  <div class="my-3 text-center">
    <div class="text-sm">
      Contracts completed:
      <strong class="font-semibold">{{ completedContractCount }}/{{ totalContractCount }}</strong>
    </div>
    <div class="mt-1 text-sm font-semibold">Prophecy eggs</div>
    <div class="text-sm">
      From contracts:
      <strong class="font-semibold">{{ collectedContractPEs }}/{{ totalContractPEs }}</strong>
    </div>
    <div class="text-sm" v-tippy="{ content: trophyDetails, allowHTML: true }">
      From trophies:
      <strong class="font-semibold"
        >{{ otherPEProgress.trophies.collected }}/{{ otherPEProgress.trophies.total }}</strong
      >
    </div>
    <div class="text-sm">
      From daily gifts:
      <strong class="font-semibold"
        >{{ otherPEProgress.gifts.collected }}/{{ otherPEProgress.gifts.total }}</strong
      >
    </div>
    <div class="text-sm">
      Total: <strong class="font-semibold">{{ collectedPEs }}/{{ totalPEs }}</strong>
    </div>
  </div>

  <div class="sm:mx-auto sm:max-w-xs sm:w-full mx-4 my-3">
    <a
      class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-gray-500 hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-400 disabled:opacity-50"
      :href="csvBlobURL"
      download="past-contracts.csv"
    >
      Export as CSV
    </a>
  </div>

  <div class="flex justify-center my-3">
    <div class="px-3 py-2 border rounded-md shadow space-y-0.5">
      <div class="flex justify-center mb-1 text-sm font-medium text-gray-900">Color coding</div>
      <div class="relative flex items-start">
        <span class="flex items-center text-green-500">
          <svg viewBox="-32 -32 576 576" class="h-4">
            <path
              fill="currentColor"
              d="M256 8C119 8 8 119 8 256s111 248 248 248 248-111 248-248S393 8 256 8z"
            ></path>
          </svg>
        </span>
        <span class="ml-2 text-xs text-gray-600">Never attempted</span>
      </div>
      <div class="relative flex items-start">
        <span class="flex items-center h-4 text-red-500">
          <svg viewBox="-32 -32 576 576" class="h-4">
            <path
              fill="currentColor"
              d="M256 8C119 8 8 119 8 256s111 248 248 248 248-111 248-248S393 8 256 8z"
            ></path>
          </svg>
        </span>
        <span class="ml-2 text-xs text-gray-600">Attempted, failed to collect prophecy egg(s)</span>
      </div>
      <div class="relative flex items-start">
        <span class="flex items-center h-4 text-yellow-500">
          <svg viewBox="-32 -32 576 576" class="h-4">
            <path
              fill="currentColor"
              d="M256 8C119 8 8 119 8 256s111 248 248 248 248-111 248-248S393 8 256 8z"
            ></path>
          </svg>
        </span>
        <span class="ml-2 text-xs text-gray-600">Attempted, failed to complete all goals</span>
      </div>
      <div class="relative flex items-start">
        <span class="flex items-center h-4 text-gray-500">
          <svg viewBox="-32 -32 576 576" class="h-4">
            <path
              fill="currentColor"
              d="M256 8C119 8 8 119 8 256s111 248 248 248 248-111 248-248S393 8 256 8z"
            ></path>
          </svg>
        </span>
        <span class="ml-2 text-xs text-gray-600">Completed</span>
      </div>
    </div>
  </div>

  <div class="flex justify-center my-3">
    <div class="space-y-0.5">
      <div class="relative flex items-start">
        <div class="flex items-center h-5">
          <input
            id="hideUnattempted"
            name="hideUnattempted"
            v-model="hideUnattempted"
            type="checkbox"
            class="focus:ring-green-500 h-4 w-4 text-green-600 border-gray-300 rounded"
          />
        </div>
        <div class="ml-2 text-sm">
          <label for="hideUnattempted" class="text-gray-600">Hide unattempted contracts</label>
        </div>
      </div>
      <div class="relative flex items-start">
        <div class="flex items-center h-5">
          <input
            id="hideCompleted"
            name="hideCompleted"
            v-model="hideCompleted"
            type="checkbox"
            class="focus:ring-green-500 h-4 w-4 text-green-600 border-gray-300 rounded"
          />
        </div>
        <div class="ml-2 text-sm">
          <label for="hideCompleted" class="text-gray-600">Hide completed contracts</label>
        </div>
      </div>
      <div class="relative flex items-start">
        <div class="flex items-center h-5">
          <input
            id="hideNoPE"
            name="hideNoPE"
            v-model="hideNoPE"
            type="checkbox"
            class="focus:ring-green-500 h-4 w-4 text-green-600 border-gray-300 rounded"
          />
        </div>
        <div class="ml-2 text-sm">
          <label for="hideNoPE" class="text-gray-600"
            >Hide contracts without prophecy egg reward</label
          >
        </div>
      </div>
    </div>
  </div>

  <div class="flex flex-col">
    <div class="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
      <div class="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8">
        <div class="shadow overflow-hidden border-b border-gray-200">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <th
                scope="col"
                class="px-6 py-2 text-center text-xs font-medium text-gray-500 uppercase"
              >
                ID
              </th>
              <th
                scope="col"
                class="px-6 py-2 text-center text-xs font-medium text-gray-500 uppercase"
              >
                Name
              </th>
              <th
                scope="col"
                class="px-6 py-2 text-center text-xs font-medium text-gray-500 uppercase"
              >
                Date
              </th>
              <th
                scope="col"
                class="px-6 py-2 text-center text-xs font-medium text-gray-500 uppercase"
              >
                code
              </th>
              <th
                scope="col"
                class="px-6 py-2 text-center text-xs font-medium text-gray-500 uppercase"
              >
                Goals
              </th>
              <th
                scope="col"
                class="px-6 py-2 text-center text-xs font-medium text-gray-500 uppercase"
              >
                PE
              </th>
            </thead>
            <tbody>
              <template v-for="contract in taggedContracts" :key="contract.id">
                <tr v-show="!contract.hidden" :class="[contract.fgClass, contract.bgClass]">
                  <td
                    class="px-6 py-1 whitespace-nowrap text-center text-sm cursor-pointer"
                    title="click to copy"
                    @click="copy(contract.id, `Copied ID '${contract.id}'`)"
                  >
                    {{ contract.id }}
                  </td>
                  <td class="px-6 py-1 whitespace-nowrap text-center text-sm">
                    <a
                      :class="[contract.fgHoverClass]"
                      :href="wikiLink(contract)"
                      target="_blank"
                      >{{ contract.name }}</a
                    >
                  </td>
                  <td class="px-6 py-1 whitespace-nowrap text-center text-sm tabular-nums">
                    {{ contract.date }}
                  </td>
                  <td
                    class="px-6 py-1 max-w-column truncate text-center text-sm cursor-pointer"
                    title="click to copy"
                    @click="copy(contract.code, `Copied code '${contract.code}'`)"
                  >
                    {{ contract.code }}
                  </td>
                  <td class="px-6 py-1 whitespace-nowrap text-center text-sm">
                    {{ contract.goalsInfo }}
                  </td>
                  <td class="px-6 py-1 whitespace-nowrap text-center text-sm">
                    {{ contract.prophecyEggInfo }}
                  </td>
                </tr>
              </template>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>

  <div class="mx-4 my-4 xl:mx-0 text-xs">
    Notes:
    <ul class="list-disc">
      <li>Contracts currently being worked on show up as never attempted. This is expected.</li>
      <li>
        For contracts with multiple incarnations, i.e. original run and leggacy run(s), only one
        incarnation is listed. If the player has attempted the contract, the last attempted
        incarnation is shown; otherwise, the latest incarnation is shown.
      </li>
      <li>
        The "Date" column shows the date on which the player last started a contract farm for the
        contract, or the estimated date the contract was offered (which may not be accurate) if it
        was never attempted.
      </li>
      <li>
        The "PE" column indicates which reward of the contract, if any, was one or more prophecy
        eggs (the number of prophecy eggs is noted in parentheses if it's more than 1). The column
        is blank if there's no PE associated with the contract. Otherwise, for older contracts
        without standard/elite tiers, this column should look like "#2", meaning the second reward
        being a PE; for newer contracts with tiers, this column should look like "std #3", meaning
        the third reward of standard tier being a PE, or "elt #2", meaning the second reward of
        elite tier being a PE. The tier shown is the tier the player last attempted the contract on,
        with the exception that if the player completed none of the goals then the tier shown
        defaults to elite (since in that case it's harder to tell which tier the player was on at
        the time, if they did make an attempt).
      </li>
      <li>
        You may
        <strong class="font-semibold">click on a contract ID or a coop code to copy it</strong>. To
        find more details about a contract, you may
        <strong class="font-semibold">click on the name of the contract</strong> which should take
        you to its wiki page, or
        <strong class="font-semibold"
          >refer to the
          <a
            href="https://github.com/fanaticscripter/EggContractor/blob/master/misc/ContractAggregator/data/contracts.csv"
            class="text-blue-500"
            target="_blank"
            >full contract archive</a
          ></strong
        >
        compiled by the author (raw CSV version
        <a href="contracts.csv" class="text-blue-500" target="_blank">here</a>).
      </li>
    </ul>
  </div>

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
import copyTextToClipboard from "copy-text-to-clipboard";

export default {
  props: {
    contracts: Array,
    csvdata: String,
    otherPEProgress: Object,
  },
  data() {
    return {
      hideUnattempted: false,
      hideCompleted: false,
      hideNoPE: false,
      popupTimeoutId: null,
      popupShow: false,
      popupMessage: "",
      csvBlobURL: window.URL.createObjectURL(new Blob([this.csvdata], { type: "text/csv" })),
    };
  },
  computed: {
    totalContractCount() {
      return this.contracts.length;
    },
    completedContractCount() {
      return this.contracts.filter(contract => !contract.incomplete).length;
    },
    totalContractPEs() {
      return this.contracts.reduce((total, contract) => total + contract.prophecyEggCount, 0);
    },
    collectedContractPEs() {
      return this.contracts.reduce(
        (total, contract) =>
          total + contract.prophecyEggCount * (contract.prophecyEggNotCollected ? 0 : 1),
        0
      );
    },
    totalPEs() {
      return (
        this.totalContractPEs +
        this.otherPEProgress.trophies.total +
        this.otherPEProgress.gifts.total
      );
    },
    collectedPEs() {
      return (
        this.collectedContractPEs +
        this.otherPEProgress.trophies.collected +
        this.otherPEProgress.gifts.collected
      );
    },
    taggedContracts() {
      const contracts = JSON.parse(JSON.stringify(this.contracts));
      let visibleIndex = 0;
      for (const contract of contracts) {
        contract.fgClass = !contract.attempted
          ? "text-green-500"
          : contract.prophecyEggNotCollected
          ? "text-red-500"
          : contract.incomplete
          ? "text-yellow-500"
          : "text-gray-500";
        contract.fgHoverClass = !contract.attempted
          ? "hover:text-green-400"
          : contract.prophecyEggNotCollected
          ? "hover:text-red-400"
          : contract.incomplete
          ? "hover:text-yellow-400"
          : "hover:text-gray-400";
        contract.hidden = this.shouldHideContract(contract);
        if (!contract.hidden) {
          contract.bgClass = visibleIndex % 2 === 1 ? "bg-gray-50" : "bg-white";
          visibleIndex++;
        } else {
          contract.bgClass = "bg-white";
        }
      }
      return contracts;
    },
    trophyDetails() {
      const details = this.otherPEProgress.trophies.eggs.map(
        e => `${e.egg}: ${e.trophy}, ${e.collected}/${e.total}`
      );
      return '<div class="">' + details.join("<br>") + "</div>";
    },
  },
  methods: {
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

    shouldHideContract(contract) {
      return (
        (this.hideUnattempted && !contract.attempted) ||
        (this.hideCompleted && !contract.incomplete) ||
        (this.hideNoPE && contract.prophecyEggCount === 0)
      );
    },

    wikiLink(contract) {
      const underscoredName = contract.name.replace(/^LEGG?ACY: /, "").replace(" ", "_");
      return `https://egg-inc.fandom.com/wiki/Contracts/${encodeURIComponent(underscoredName)}`;
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
