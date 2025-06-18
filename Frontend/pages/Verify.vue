<script setup>
import { useRoute, useRouter } from "vue-router";
import { ref, onMounted } from "vue";

useSeoMeta({
  title: "recipe-app | Verify Email",
  description: "The project app meta.",
});
const useAuthStore = authStore();
const route = useRoute();
const router = useRouter();

const verificationStatus = ref(null);

onMounted(async () => {
  try {
    console.log("Route query parameters:", route.query);

    const user_id = parseInt(route.query.user_id, 10);
    const verification_token = route.query.verification_token;

    if (!user_id || !verification_token) {
      console.error("Missing user_id or verification_token:", {
        user_id,
        verification_token,
      });
      verificationStatus.value = false;
      return;
    }

    console.log("Attempting to verify email with:", {
      user_id,
      verification_token,
    });

    await useAuthStore.verifyEmail({ user_id, verification_token });

    verificationStatus.value = useAuthStore.isEmailVerified;
    console.log("Verification status:", verificationStatus.value);
  } catch (error) {
    console.error("Error during email verification:", error);
    verificationStatus.value = false;
  }
});

function gotoHome() {
  console.log("Navigating to home page");
  router.push("/");
}
</script>

<template>
  <div
    class="flex flex-col items-center justify-center h-screen dark:bg-[#20161F]"
  >
    <div class="text-center">
      <h1 class="text-2xl font-bold mb-4">
        {{
          verificationStatus === null
            ? "Verifying..."
            : verificationStatus
            ? "Verification Successful!"
            : "Verification Failed"
        }}
      </h1>
      <p v-if="verificationStatus" class="mb-4 text-green-600 text-xl">
        Your email has been successfully verified. You can now contine to homepage.
      </p>
      <p v-else-if="verificationStatus === false" class="mb-4 text-red-600">
        Verification failed. Please try again or contact support.
      </p>
      <button
        v-if="verificationStatus"
        @click="gotoHome"
        class="px-4 py-2 text-white bg-blue-500 rounded hover:bg-blue-700"
      >
        Go to home
      </button>
    </div>
  </div>
</template>
