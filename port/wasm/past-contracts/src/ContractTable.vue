<template>
  <div class="flex justify-center my-3">
    <div class="px-3 py-2 border rounded-md shadow space-y-0.5">
      <div class="flex justify-center mb-1 text-sm font-medium text-gray-900">Color coding</div>
      <div class="relative flex items-start">
        <span class="flex items-center text-green-500">
          <svg viewBox="-32 -32 576 576" class="h-4"><path fill="currentColor" d="M256 8C119 8 8 119 8 256s111 248 248 248 248-111 248-248S393 8 256 8z"></path></svg>
        </span>
        <span class="ml-2 text-xs text-gray-600">Never attempted</span>
      </div>
      <div class="relative flex items-start">
        <span class="flex items-center h-4 text-red-500">
          <svg viewBox="-32 -32 576 576" class="h-4"><path fill="currentColor" d="M256 8C119 8 8 119 8 256s111 248 248 248 248-111 248-248S393 8 256 8z"></path></svg>
        </span>
        <span class="ml-2 text-xs text-gray-600">Attempted, failed to collect prophecy egg(s)</span>
      </div>
      <div class="relative flex items-start">
        <span class="flex items-center h-4 text-yellow-500">
          <svg viewBox="-32 -32 576 576" class="h-4"><path fill="currentColor" d="M256 8C119 8 8 119 8 256s111 248 248 248 248-111 248-248S393 8 256 8z"></path></svg>
        </span>
        <span class="ml-2 text-xs text-gray-600">Attempted, failed to complete all goals</span>
      </div>
      <div class="relative flex items-start">
        <span class="flex items-center h-4 text-gray-500">
          <svg viewBox="-32 -32 576 576" class="h-4"><path fill="currentColor" d="M256 8C119 8 8 119 8 256s111 248 248 248 248-111 248-248S393 8 256 8z"></path></svg>
        </span>
        <span class="ml-2 text-xs text-gray-600">Completed</span>
      </div>
    </div>
  </div>

  <div class="my-3">
    <div class="relative flex items-start justify-center">
      <div class="flex items-center h-5">
        <input id="hideCompleted" name="hideCompleted" v-model="hideCompleted" type="checkbox" class="focus:ring-green-500 h-4 w-4 text-green-600 border-gray-300 rounded">
      </div>
      <div class="ml-2 text-sm">
        <label for="hideCompleted" class="text-gray-600">Hide completed contracts</label>
      </div>
    </div>
    <div class="relative flex items-start justify-center">
      <div class="flex items-center h-5">
        <input id="hideNoPE" name="hideNoPE" v-model="hideNoPE" type="checkbox" class="focus:ring-green-500 h-4 w-4 text-green-600 border-gray-300 rounded">
      </div>
      <div class="ml-2 text-sm">
        <label for="hideNoPE" class="text-gray-600">Hide contracts without prophecy egg reward</label>
      </div>
    </div>
  </div>

  <div class="flex flex-col">
    <div class="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
      <div class="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8">
        <div class="shadow overflow-hidden border-b border-gray-200">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <th scope="col" class="px-6 py-2 text-center text-xs font-medium text-gray-500 uppercase">ID</th>
              <th scope="col" class="px-6 py-2 text-center text-xs font-medium text-gray-500 uppercase">Name</th>
              <th scope="col" class="px-6 py-2 text-center text-xs font-medium text-gray-500 uppercase">Date</th>
              <th scope="col" class="px-6 py-2 text-center text-xs font-medium text-gray-500 uppercase">code</th>
              <th scope="col" class="px-6 py-2 text-center text-xs font-medium text-gray-500 uppercase">#Goals</th>
              <th scope="col" class="px-6 py-2 text-center text-xs font-medium text-gray-500 uppercase">PE</th>
            </thead>
            <tbody>
              <template v-for="(contract, index) in contracts" :key="contract.id">
                <tr v-if="(!hideCompleted || contract.incomplete) && (!hideNoPE || contract.hasProphecyEgg)" :class="[index % 2 === 1 ? 'bg-gray-50' : 'bg-white', contract.prophecyEggNotCollected ? 'text-red-500' : contract.incomplete ? 'text-yellow-500' : 'text-gray-500']">
                  <td class="px-6 py-1 whitespace-nowrap text-center text-sm">{{ contract.id }}</td>
                  <td class="px-6 py-1 whitespace-nowrap text-center text-sm">{{ contract.name }}</td>
                  <td class="px-6 py-1 whitespace-nowrap text-center text-sm">{{ contract.date }}</td>
                  <td class="px-6 py-1 max-w-column truncate text-center text-sm">{{ contract.code }}</td>
                  <td class="px-6 py-1 whitespace-nowrap text-center text-sm">{{ contract.goals }}</td>
                  <td class="px-6 py-1 whitespace-nowrap text-center text-sm">{{ contract.prophecyEgg }}</td>
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
      <li>The contracts listed are past contracts the player has attempted. The ones never seen nor attempted cannot be retrieved. Consult a complete contract list to find out which ones were missed.</li>
      <li>The "Date" column shows the date on which the player started the respective contract farm.</li>
      <li>The "PE" column indicates which reward of the contract, if any, was a prophecy egg. The column is blank if there's no PE associated with the contract. Otherwise, for older contracts without standard/elite tiers, this column should look like "#2", meaning the second reward being a PE; for newer contracts with tiers, this column should look like "std #3", meaning the third reward of standard tier being a PE, or "elt #2", meaning the second reward of elite tier being a PE. The tier shown is the tier the player last attempted the contract on, with the exception that if the player completed none of the goals then the tier shown defaults to elite (since in that case it's harder to tell which tier the player was on at that time).</li>
    </ul>
  </div>
</template>

<script>
export default {
  props: {
    contracts: Array,
  },
  data() {
    return {
      hideCompleted: false,
      hideNoPE: false,
    };
  },
}
</script>
