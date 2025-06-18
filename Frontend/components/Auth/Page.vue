<script setup>
import { ref, onMounted, watch } from "vue";
import { gsap } from "gsap";

const showSignup = ref(false);

function toggleForm() {
  showSignup.value = !showSignup.value;
}

// GSAP Animation for the text
onMounted(() => {
  animateText();
});

// Watch for changes in `showSignup` and re-animate the text
watch(showSignup, () => {
  animateText();
});

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
</script>

<template>
  <div class="flex h-screen w-full">
    <!-- Left Section: Animated Text -->
    <div
      class="relative flex-1 hidden w-0 lg:flex items-center justify-center bg-gradient-to-r from-green-600 via-green-200 to-indigo-400 dark:bg-gradient-to-r dark:from-[#20161F] dark:via-[#2A1C23] dark:to-[#1A0E14]"
    >
      <div class="text-center space-y-6">
        <h1
          class="animated-text text-7xl font-bold bg-clip-text bg-gradient-to-r from-green-400 to-black-600"
        >
          Welcome to Kushna!
        </h1>
        <p
          class="animated-text text-4xl font-semibold bg-clip-text bg-gradient-to-r from-green-400 to-blue-700"
        >
          {{
            showSignup
              ? "Join our community today."
              : "Sign in to continue your journey."
          }}
        </p>
      </div>
    </div>

    <!-- Right Section: Form -->
    <div
      class="flex flex-col justify-center px-4 sm:px-6 lg:flex-none lg:px-20 xl:px-24 flex-1 dark:bg-[#20161F]"
    >
      <!-- Form Toggle Button and Conditional Heading -->
      <div
        class="flex flex-row gap-4 pt-20 sm:items-center sm:justify-center lg:items-start lg:justify-start"
      >
        <h2
          class="text-center border-b-2 border-black dark:text-gray-200 dark:border-gray-200"
        >
          {{
            showSignup
              ? "Already have an account?"
              : "You haven't signed up yet?"
          }}
        </h2>
        <button
          @click="toggleForm"
          class="text-cyan-400 p-2 rounded-full hover:text-cyan-900"
        >
          {{ showSignup ? "Login" : "Signup" }}
        </button>
      </div>

      <!-- Form Container -->
      <div class="flex items-center w-full max-w-sm mx-auto lg:w-96 mt-6">
        <div v-if="showSignup">
          <AuthSignup />
        </div>
        <div v-else>
          <AuthLogin />
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