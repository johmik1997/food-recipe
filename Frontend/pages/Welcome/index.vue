<script setup>
import { onMounted, ref } from "vue";
import { authStore } from "~/stores/auth";
import { gsap } from "gsap";

useSeoMeta({
  title: "recipe-app | Welcome",
  description: "The project app meta.",
});
definePageMeta({
  layout: false,
});

// Reactive state
const auth = authStore();
const isAuthenticated = ref(false);
const id = ref("");
const name = ref("");
const user = ref({});

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

onMounted(() => {
  auth.initAuthState();

  isAuthenticated.value = auth.isAuthed;
  id.value = auth.userId;
  name.value = auth.user?.userName || "";
  user.value = auth.user;
  animateText();
});
</script>

<template>
  <div
    class="flex flex-col items-center justify-center min-h-screen bg-gray-100 p-6 dark:bg-[#20161F]"
  >
    <div class="bg-white shadow-lg rounded-lg p-8 max-w-xl text-center">
      <h1
        class="animated-text text-5xl font-bold bg-clip-text bg-gradient-to-r from-yellow-400 to-pink-600"
      >
        We have sent a verification email to your provided email address
      </h1>
      <p
        class="animated-text text-3xl font-semibold bg-clip-text bg-gradient-to-r from-green-400 to-blue-500"
      >
        Welcome! Please verify your email to proceed.
      </p>
      <a
        href="https://mail.google.com/mail/u/0/#inbox"
        target="_blank"
        rel="noopener noreferrer"
        class="px-40 py-3 mt-5 font-bold text-lg text-[#20161F] flex grid-flow-col border-black gap-4 cursor-pointer transition-colors duration-300 bg-cyan-200 border-1 rounded-full bg-gradient-to-r from-purple-600 via-green-200 to-indigo-400 dark:bg-gradient-to-r dark:from-[#20161F] dark:via-[#2A1C23] dark:to-[#1A0E14] hover:ring-blue-500 hover:border-blue-500 dark:text-white"
      >
        Go to your email inbox
      </a>
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

/* Ensure the link is clickable */
a {
  display: inline-block;
  z-index: 10;
  pointer-events: auto;
}
</style>
