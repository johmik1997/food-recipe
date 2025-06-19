<script setup>
import { Field, useField } from "vee-validate";

const props = defineProps({
  name: { type: String, required: true },
  placeholder: { type: String, default: "" },
  type: { type: String, default: "text" },
  label: { type: String, default: "" },
  isRequired: { type: Boolean, default: false },
  errorMessage: { type: String, default: "" },
  isPassword: { type: Boolean, default: false },
  isDisabled: { type: Boolean, default: false },
});

const inputType = ref(props.type);
const togglePasswordShow = ref(false);

const handleToggleType = () => {
  togglePasswordShow.value = !togglePasswordShow.value;
  inputType.value = togglePasswordShow.value ? "text" : props.type;
};

const { value, errorMessage } = useField(props.name);
</script>

<template>
  <div class="relative mb-3">
    <label
      v-if="label"
      :for="`field-${name}`"
      class="block pl-3 ml-px text-sm font-medium text-gray-700 dark:text-gray-100"
      :class="{ 'text-red-500': errorMessage }"
    >
      {{ label }}
    </label>

    <input
      :id="`field-${name}`"
      :type="inputType"
      :name="name"
      v-model="value"
      :placeholder="placeholder"
      :disabled="isDisabled"
      class="block w-full px-8 py-2 input input-bordered border-gray-300 rounded-full shadow-sm focus:ring-green-500 focus:border-green-500 sm:text-sm bg-white"
      :class="{ 'border-red-500 focus:border-red-500': errorMessage }"
    />

    <!-- Password toggle -->
    <div
      class="absolute right-2 top-6 p-2 flex items-center text-neutral-500 z-10"
      v-if="isPassword"
    >
      <button type="button" @click="handleToggleType">
        <span v-if="togglePasswordShow">
          <!-- Eye icon -->
          👁️
        </span>
        <span v-else>
          <!-- Eye off icon -->
          🚫
        </span>
      </button>
    </div>

    <p class="mt-2 text-sm text-red-600" v-if="errorMessage">
      {{ errorMessage }}
    </p>
  </div>
</template>


<style scoped>
.input {
  transition: all 0.3s ease;
  border: 1px solid #d1d5db;
}

.input:focus {
  border-color: #10b981;
  box-shadow: 0 0 0 3px rgba(16, 185, 129, 0.2);
  outline: none;
}

.input:hover:not(:disabled) {
  border-color: #10b981;
}

.input:disabled {
  background-color: #f3f4f6;
  cursor: not-allowed;
}

/* Password toggle button */
button[type="button"] {
  transition: all 0.2s ease;
}

button[type="button"]:hover {
  color: #10b981;
  transform: scale(1.05);
}

button[type="button"]:focus {
  outline: none;
  color: #10b981;
}
</style>