<script setup>
import { Form } from "vee-validate";
import * as yup from "yup";
import { useToast } from "vue-toast-notification";
import { gsap } from "gsap";

useSeoMeta({
  title: "recipe-app | Reset Password",
  description: "The project app meta.",
});

definePageMeta({
  layout: false,
});

const toast = useToast();
const router = useRouter();
const token = ref("");
const userId = ref("");
const route = useRoute();
const auth = authStore();

const formError = ref({
  password: "",
  passwordConfirmation: "",
});

onMounted(() => {
  const queryParams = new URLSearchParams(window.location.search);
  token.value = queryParams.get("token");
  userId.value = queryParams.get("id");
  if (!token.value || !userId.value) {
    router.push("/");
  }
});
const passwordNumber = /^((?=\S*?[0-9]).{8,})\S$/;
const passwordCapitalLetter = /^((?=\S*?[A-Z]).{8,})\S$/;

const errorMessages = {
  password: {
    required: "Password field is required!",
    min_password_len: "The minimum length of password is 6",
    not_valid_password_capitalizer:
      "Atleast there is 1 capital letter in the password",
    not_vaild_password_digit: "Password also must have a digit in it",
  },

  passwordConfirmation: {
    must_be_same: "Password confirmation must be same with password",
  },
};

const schema = yup.object({
  password: yup
    .string()
    .required(errorMessages.password.required)
    .min(8, errorMessages.password.min_password_len)
    .matches(
      passwordCapitalLetter,
      errorMessages.password.not_valid_password_capitalizer
    )
    .matches(passwordNumber, errorMessages.password.not_vaild_password_digit),
  passwordConfirmation: yup
    .string()
    .required()
    .oneOf(
      [yup.ref("password"), null],
      errorMessages.passwordConfirmation.must_be_same
    ),
});

const handlePasswordReset = async (value) => {
  value.token = token.value;
  value.userId = Number(userId.value);
  if (!isNaN(value.userId)) {
    console.log("the value:", value);
    const result = await auth.resetPassword(value);
    if (result == true) {
      if (auth.$state.successMessage) {
        const message = auth.$state.successMessage;
        auth.setSuccessMessage("");
        toast.success(message);
      } else {
        toast.success(
          "Your password is updated successfully!. Preceed to login"
        );
      }
      router.push("/auth");
    } else {
      toast.error("Something went wrong! please try again.");
    }
  } else {
    toast.error("Invalid user reference!");
  }
};

function animateText() {
  const textElements = document.querySelectorAll(".animated-text");

  gsap.from(textElements, {
    duration: 1,
    opacity: 0,
    y: 50,
    stagger: 0.3,
    ease: "power3.out",
  });
}
// GSAP Animation for the text
onMounted(() => {
  animateText();
});
</script>

<template>
  <div class="flex h-screen w-full dark:bg-[#20161F]">
    <!-- Background Image Section -->
    <div
      class="relative flex-1 hidden w-0 lg:flex items-center justify-center bg-gradient-to-r from-purple-600 via-green-200 to-indigo-400 dark:bg-gradient-to-r dark:from-[#20161F] dark:via-[#2A1C23] dark:to-[#1A0E14]"
    >
      <div class="text-center space-y-6">
        <h1
          class="animated-text text-7xl font-bold bg-clip-text bg-gradient-to-r from-yellow-400 to-pink-600"
        >
          Welcome to Kushna!
        </h1>
        <p
          class="animated-text text-4xl font-semibold bg-clip-text bg-gradient-to-r from-green-400 to-blue-500"
        >
          Reset your Password in no time.
        </p>
      </div>
    </div>

    <!-- Form Section -->
    <div class="flex flex-1 items-center justify-center">
      <div class="w-full max-w-sm lg:w-96">
        <!-- Border Line Section -->
        <div
          class="flex flex-row items-center justify-center gap-2 text-center mb-6"
        >
          <h1 class="text-lg font-semibold">Wanna</h1>
          <NuxtLink
            class="text-cyan-400 font-bold p-2 rounded-full hover:text-cyan-900"
            to="/"
          >
            Signup | Login
          </NuxtLink>
        </div>
        <div class="border-b-2 border-black mb-6 dark:border-gray-100"></div>

        <!-- Form Card -->
        <div
          class="w-full border-2 p-4 border-cyan-100 shadow-lg rounded-lg px-2 sm:px-6 lg:px-4 xl:px-6 dark:bg-[#20161F]"
        >
          <Form
            @submit="handlePasswordReset"
            :validation-schema="schema"
            v-slot="{ errors }"
            class="flex flex-col gap-4 items-center"
          >
            <UIInput
              name="password"
              placeholder="*********"
              label="password"
              :is-required="true"
              type="password"
              :error-message="
                errors.password
                  ? errors.password
                  : formError.password
                  ? formError.password
                  : ''
              "
              :is-password="true"
            />
            <UIInput
              name="passwordConfirmation"
              placeholder="Confirm Your New Password"
              label="Password Confirmation"
              :is-required="true"
              type="password"
              :error-message="
                errors.passwordConfirmation
                  ? errors.passwordConfirmation
                  : formError.passwordConfirmation
                  ? formError.passwordConfirmation
                  : ''
              "
              :is-password="true"
            />
            <div class="flex justify-center mt-4">
              <button
                class="bg-cyan-200 border-2 rounded-full px-6 py-2 hover:bg-cyan-600"
              >
                <span v-if="auth.$state.onLoad == false">
                  Update Password
                </span>
                <div
                  class="w-full flex items-center justify-center gap-2"
                  v-else
                >
                  <UILoading />
                </div>
              </button>
            </div>
          </Form>
        </div>
      </div>
    </div>
  </div>
</template>
<style scoped>
/* Add gradient text styling */
.animated-text {
  background-clip: text;
  -webkit-background-clip: text;
  color: transparent;
}
</style>
