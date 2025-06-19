<script setup>
import { ref, onMounted, watch } from "vue";
import { gsap } from "gsap";

const showSignup = ref(false);

definePageMeta({
  layout: false,
});

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
      class="relative flex-1 hidden w-0 lg:flex items-center justify-center bg-gradient-to-r from-teal-500 via-blue-400 to-indigo-500"
    >
      <div class="text-center space-y-6 px-8">
        <h1
          class="animated-text text-7xl font-bold bg-clip-text bg-gradient-to-r from-white to-amber-200"
        >
          Welcome to Kushna!
        </h1>
        <p
          class="animated-text text-4xl font-semibold bg-clip-text bg-gradient-to-r from-amber-100 to-white"
        >
          {{
            showSignup
              ? "Join our culinary community"
              : "Sign in to continue cooking"
          }}
        </p>
      </div>
    </div>

    <!-- Right Section: Form -->
    <div
      class="flex flex-col justify-center px-4 sm:px-6 lg:flex-none lg:px-20 xl:px-24 flex-1 bg-white"
    >
      <!-- Form Toggle Button and Conditional Heading -->
      <div
        class="flex flex-row gap-4 pt-20 sm:items-center sm:justify-center lg:items-start lg:justify-start"
      >
        <h2 class="text-center border-b-2 border-teal-500 text-gray-700">
          {{
            showSignup
              ? "Already have an account?"
              : "New to Kushna?"
          }}
        </h2>
        <button
          @click="toggleForm"
          class="text-teal-600 font-medium hover:text-teal-800 transition-colors"
        >
          {{ showSignup ? "Login" : "Sign Up" }}
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
  text-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

/* Smooth background transition */
.bg-gradient-to-r {
  transition: background 0.5s ease;
}
</style>