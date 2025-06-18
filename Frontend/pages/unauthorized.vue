<script setup>
import { useRouter } from "vue-router";
import { ref, onMounted } from "vue";

useSeoMeta({
  title: "recipe-app | Unauthorized",
  description: "The project app meta.",
});

definePageMeta({
  layout: false,
});
const router = useRouter();

// Redirect the user back to the previous page
const goBack = () => {
  const previousRoute = window.sessionStorage.getItem("previousRoute");
  if (previousRoute) {
    router.push(previousRoute);
    window.sessionStorage.removeItem("previousRoute");
  } else {
    router.push("/");
  }
};

const isVisible = ref(false);

onMounted(() => {
  setTimeout(() => {
    isVisible.value = true;
  }, 100);
});

// Redirect the user to the home page
const goHome = () => {
  router.push("/");
};
</script>

<template>
  <div
    class="unauthorized-page h-screen flex flex-col items-center justify-center text-center p-[50px] text-red-950 border-1 dark:bg-[#20161F]"
  >

    <img
      src="/Do not enter sign-pana.svg"
      alt="Restricted Area"
      class="w-[80%] max-w-md transition-all duration-1000 ease-in-out hover:rotate-6 hover:scale-105"
      :class="{
        'opacity-0 scale-90': !isVisible,
        'opacity-100 scale-100': isVisible,
      }"
    />
    <h1
      class="text-4xl md:text-6xl font-bold dark:text-white mb-8 transition-all duration-1000 ease-in-out"
      :class="{
        'opacity-0 translate-y-10': !isVisible,
        'opacity-100 translate-y-0': isVisible,
      }"
    >
      Oops! Our kitchen is currently a disaster—only family members can handle
      this chaos.
    </h1>
    <p
      class="text-sm md:text-xl font-bold dark:text-white mb-8 transition-all duration-1000 ease-in-out"
      :class="{
        'opacity-0 translate-y-10': !isVisible,
        'opacity-100 translate-y-0': isVisible,
      }"
    >
      Looks like you're not on the VIP guest list. No peeking at our secret
      recipes!
    </p>

    <div class="flex space-x-4">

      <transition name="bounce">
        <button
          v-if="isVisible"
          @click="goBack"
          class="px-6 py-3 bg-gradient-to-r from-purple-500 to-pink-500 rounded-full text-white hover:from-purple-600 hover:to-pink-600 focus:outline-none transform hover:scale-105 transition-all duration-300"
        >
          Go Back
        </button>
      </transition>

      <!-- Go to Home Button -->
      <transition name="bounce">
        <button
          v-if="isVisible"
          @click="goHome"
          class="px-6 py-3 bg-gradient-to-r from-blue-500 to-teal-500 rounded-full text-white hover:from-blue-600 hover:to-teal-600 focus:outline-none transform hover:scale-105 transition-all duration-300"
        >
          Go to Home
        </button>
      </transition>
    </div>
  </div>
</template>

<style scoped>
/* Optional: Add custom animations */
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.fade-in {
  animation: fadeIn 1s ease-in-out;
}

.bounce-enter-active {
  animation: bounce 0.6s ease-in-out;
}

@keyframes bounce {
  0% {
    transform: scale(0.8);
    opacity: 0;
  }
  50% {
    transform: scale(1.1);
    opacity: 1;
  }
  100% {
    transform: scale(1);
  }
}

/* Add a subtle background animation */
.unauthorized-page {
  /* background: linear-gradient(-45deg, #20161F, #2e141e, #35063b, #10011b); */
  background-size: 400% 400%;
  animation: gradientBG 15s ease infinite;
}

@keyframes gradientBG {
  0% {
    background-position: 0% 50%;
  }
  50% {
    background-position: 100% 50%;
  }
  100% {
    background-position: 0% 50%;
  }
}
</style>
