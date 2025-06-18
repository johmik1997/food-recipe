<script setup>
import { Field } from "vee-validate";

const props = defineProps({
  name: {
    type: String,
    required: true,
  },
  value: {
    type: String,
    required: true,
  },
  placeholder: {
    type: String,
    default: "",
  },
  type: {
    type: String,
    default: "text",
  },
  label: {
    type: String,
    default: "",
  },
  isRequired: {
    type: Boolean,
    default: false,
  },
  errorMessage: {
    type: String,
    default: "",
  },
  isPassword: {
    type: Boolean,
    default: false,
  },
  editMode: {
    type: Boolean,
    default: false,
  },
  isDisabled: {
    type: Boolean,
    default: false,
  },
});

const inputType = ref(props.type);
const togglePasswordShow = ref(false);

const handleToggleType = () => {
  togglePasswordShow.value = !togglePasswordShow.value;
  inputType.value = togglePasswordShow.value ? "text" : props.type;
};
</script>
<template>
  <div class="relative mb-3">
    <label
      v-if="label"
      :for="`field-${name}`"
      class="block pl-3 ml-px text-sm font-medium text-gray-700 dark:text-gray-100"
      :class="{ 'text-red-500 focus:text-red-500': errorMessage }"
    >
      {{ label }}
    </label>
    <Field
      :type="inputType"
      :value="editMode ? value : ''"
      :name="name"
      :id="`field-${name}`"
      class="block w-full px-8 py-2 input input-bordered border-gray-300 rounded-full shadow-sm focus:ring-green-500 focus:border-green-500 sm:text-sm bg-white"
      :placeholder="placeholder"
      :disabled="isDisabled"
      :class="{ 'border-red-500 focus:border-red-500': errorMessage }"
    />
    <div
      class="absolute right-2 top-6 p-2 flex items-center text-neutral-500 z-10"
      v-if="isPassword"
    >
      <button type="button" @click="handleToggleType">
        <span v-if="togglePasswordShow">
          <!-- Eye Icon -->
          <svg
            width="20"
            height="20"
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 14 14"
          >
            <g
              fill="none"
              stroke="currentColor"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path
                d="M13.23 6.246c.166.207.258.476.258.754c0 .279-.092.547-.258.754C12.18 9.025 9.79 11.5 7 11.5c-2.79 0-5.18-2.475-6.23-3.746A1.208 1.208 0 0 1 .512 7c0-.278.092-.547.258-.754C1.82 4.975 4.21 2.5 7 2.5c2.79 0 5.18 2.475 6.23 3.746"
              />
              <path d="M7 9a2 2 0 1 0 0-4a2 2 0 0 0 0 4" />
            </g>
          </svg>
        </span>
        <span v-else>
          <!-- Eye Off Icon -->
          <svg
            width="20"
            height="20"
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 14 14"
          >
            <g
              fill="none"
              stroke="currentColor"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path
                d="M3.63 3.624C4.621 2.98 5.771 2.5 7 2.5c2.79 0 5.18 2.475 6.23 3.746c.166.207.258.476.258.754c0 .279-.092.547-.258.754c-.579.7-1.565 1.767-2.8 2.583m-1.93.933c-.482.146-.984.23-1.5.23c-2.79 0-5.18-2.475-6.23-3.746A1.208 1.208 0 0 1 .512 7c0-.278.092-.547.258-.754c.333-.402.8-.926 1.372-1.454"
              />
              <path d="M8.414 8.414a2 2 0 0 0-2.828-2.828M13.5 13.5L.5.5" />
            </g>
          </svg>
        </span>
      </button>
    </div>
    <p class="mt-2 text-sm text-red-600" v-if="errorMessage">
      <span class="font-medium"> </span> {{ errorMessage }}
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