<template>
  <table class="min-w-full divide-y divide-gray-200">
    <thead class="bg-gray-50">
      <tr>
        <th v-for="label in labels" :key="label.sortBy" @click="sortBy = label.sortBy" scope="col" class="px-6 py-2 text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer" :class="label.name === 'Player' ? 'text-left' : 'text-center'">
          {{ label.name }}
          <!-- Use visibility for the arrow so that column widths don't change when sorting a different column. -->
          <span class="inline-block w-0 text-gray-400" :class="{ invisible: label.sortBy != sortBy }">&nbsp;&#x25BC;</span>
        </th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="(member, index) in sortedMembers" :key="member.id" :class="index % 2 === 1 ? 'bg-gray-50' : 'bg-white'">
        <td class="px-6 py-1 whitespace-nowrap text-sm text-gray-500" :class="{ 'CoopTable__member--snoozing': !member.isActive }" :title="member.id">{{ member.name }}</td>
        <td class="px-6 py-1 whitespace-nowrap text-center text-sm text-gray-500">{{ member.eggsLaidStr }}</td>
        <td class="px-6 py-1 whitespace-nowrap text-center text-sm text-gray-500">{{ member.eggsPerHourStr }}</td>
        <td class="px-6 py-1 whitespace-nowrap text-center text-sm text-gray-500">{{ member.earningBonusPercentageStr }}</td>
        <td class="px-6 py-1 whitespace-nowrap text-center text-sm text-gray-500">{{ member.tokens }}</td>
      </tr>
    </tbody>
  </table>
</template>

<script>
export default {
  props: {
    members: Array,
  },

  data() {
    return {
      labels: [
        { sortBy: "name", name: "Player" },
        { sortBy: "eggsLaid", name: "Laid" },
        { sortBy: "eggsPerHour", name: "Rate/hr" },
        { sortBy: "earningBonusPercentage", name: "EB%" },
        { sortBy: "tokens", name: "Tokens" },
      ],
      sortBy: "eggsLaid",
    };
  },

  computed: {
    sortedMembers() {
      return [...this.members].sort((m1, m2) => {
        if (this.sortBy === "name") {
          return m1[this.sortBy].localeCompare(m2[this.sortBy])
        } else {
          return m2[this.sortBy] - m1[this.sortBy]
        }
      });
    },
  },
};
</script>
