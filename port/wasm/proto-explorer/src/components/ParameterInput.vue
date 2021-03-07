<template>
  <div>
    <label for="user_id" class="block text-sm font-medium text-gray-700">{{ label }}</label>
    <div class="mt-1">
      <input
        type="text"
        :name="name"
        :id="name"
        class="max-w-lg sm:max-w-xs shadow-sm focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md"
        :placeholder="placeholder"
        :value="modelValue"
        :required="required"
        @input="$emit('update:modelValue', $event.target.value)"
      />
    </div>
  </div>
</template>

<script>
export default {
  props: {
    name: {
      type: String,
      required: true,
    },
    label: {
      type: String,
      required: true,
    },
    required: {
      type: Boolean,
      default: false,
    },
    modelValue: String,
    placeholder: String,
  },

  emits: ["update:modelValue"],
};
</script>
