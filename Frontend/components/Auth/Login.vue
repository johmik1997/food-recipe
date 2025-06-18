<script setup>
import { ref } from "vue";
import { Form } from "vee-validate";
import * as yup from "yup";
import { useToast } from "vue-toast-notification";
import { authStore } from "../../stores/auth";

useSeoMeta({
  title: "recipe-app | Login",
  description: "my cool recipe app.",
});

definePageMeta({
  layout: false,
});

const toast = useToast();
const useAuthStore = authStore();
const router = useRouter();

const loginError = ref({
  email: "",
  password: "",
});

const emailRegex = /^((?!\.)[\w_.]*[^.])(@\w+)(\.\w+(\.\w+)?[^.\W])$/gm;
const passwordCapitalLetter = /^((?=\S*?[A-Z]).{8,})\S$/;
const passwordNumber = /^((?=\S*?[0-9]).{8,})\S$/;

const errorMessages = {
  email: {
    required: "Email field is required!",
    invalid_email: "Invalid email address!",
  },
  password: {
    required_password: "Password field is required!",
    min_password_len: "Password must be at least 8 characters long!",
    not_valid_password_capitalzer:
      "Password must contain at least one capital letter!",
    not_vaild_password_digit: "Password must contain at least one digit!",
  },
};

const Schema = yup.object({
  email: yup
    .string()
    .required(errorMessages.email.required)
    .matches(emailRegex, errorMessages.email.invalid_email),
  password: yup
    .string()
    .required(errorMessages.password.required_password)
    .min(8, errorMessages.password.min_password_len)
    .matches(
      passwordCapitalLetter,
      errorMessages.password.not_valid_password_capitalzer
    )
    .matches(passwordNumber, errorMessages.password.not_vaild_password_digit),
});

const ResetPassword = ref(false);

async function handleLogin(values) {
  useAuthStore.changeOnLoad(true);
  const result = await useAuthStore.loginUser(values);
  if (result) {
    toast.success("Welcome, you have logged in successfully!");
    if (useAuthStore.$state.isAdmin) {
      router.push("/admin/dashboard");
    } else {
      router.push("/");
    }
  } else {
    const message = useAuthStore.$state.errorMessage || "Invalid credentials";
    useAuthStore.setErrorMessage("");
    toast.error(message);
  }
}

function toggleForm() {
  ResetPassword.value = !ResetPassword.value;
}
</script>

<template>
  <div
    class="w-full border-2 border-green-400 shadow-lg rounded-lg items-center justify-center px-2 sm:px-6 lg:flex-none lg:px-4 xl:px-6 dark:bg-[#20161F]]"
  ></div>
    <Form
      class="pt-5 space-y-6"
      @submit="handleLogin"
      :validation-schema="Schema"
      v-slot="{ errors }"
    >
      <UIInput
        name="email"
        label="Email"
        placeholder="Enter your email"
        :error-message="errors.email || loginError.email"
        :is-required="true"
      />
      <UIInput
        name="password"
        label="Password"
        placeholder="********"
        type="password"
        :error-message="errors.password || loginError.password"
        :is-password="true"
      />

      <button
        @click.prevent="toggleForm"
        class="text-green-400 text-xs hover:text-black cursor-pointer"
      >
        Forgot password?
      </button>

      <div class="flex justify-center mt-4">
        <button
          type="submit"
          class="bg-green-400 text-black border-2 rounded-full px-6 mb-2 py-2 hover:bg-green-400 hover:ring-blue-500 hover:border-blue-500"
        >
          <span v-if="!useAuthStore.$state.onLoad">Login</span>
          <div
            class="w-full flex items-center justify-center gap-2"
            v-else
            aria-live="polite"
          >
            <svg
              aria-hidden="true"
              role="status"
              class="inline w-4 h-4 me-3 text-black animate-spin dark:text-green-400"
              viewBox="0 0 100 101"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
            >
              <!-- SVG paths -->
            </svg>
            Loading...
          </div>
        </button>
      </div>
    </Form>
    <div v-if="ResetPassword">
      <AuthForgotPassword />
    </div>
</template>