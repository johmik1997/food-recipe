<script setup>
import { computed, onMounted } from "vue";
import { useRouter } from "vue-router";
import { useToast } from "vue-toast-notification";
import { authStore } from "@/stores/auth";

const toast = useToast();
const auth = authStore();
const router = useRouter();

onMounted(async () => {
  await auth.getProfile();
});

const user = computed(() => auth.user);

const handleLogout = async () => {
  const result = await auth.logout();
  if (result) {
    toast.success("✅ You have logged out successfully");
    router.push("/");
  } else {
    toast.error("❌ Unable to logout");
  }
};
</script>

<template>
  <div class="sticky top-0 z-10">
    <div
      class="navbar bg-base-100 container mx-auto rounded-xl border-1 border-cyan-200 border-b shadow-md custom-navbar dark:bg-[#545454] mt-2"
    >
      <div class="flex-1">
        <!-- Logo -->
      🍲 Kushna
      </div>

      <!-- Navigation Links -->
      <div class="flex-row flex gap-1 mt-4 text-sm">
        <NuxtLink
          to="/"
          class="hidden lg:block btn btn-ghost text-xl dark:text-[#C9CF43]"
          >Home</NuxtLink
        >
        <NuxtLink
          to="/recipes"
          class="hidden lg:block btn btn-ghost text-xl dark:text-[#C9CF43]"
          >Recipes</NuxtLink
        >
        <NuxtLink
          to="/recipes/create"
          class="hidden lg:block btn btn-ghost text-xl dark:text-[#C9CF43]"
          >Let's Cook</NuxtLink
        >
        <NuxtLink
          to="/about"
          class="hidden lg:block btn btn-ghost text-xl dark:text-[#C9CF43]"
          >About Us</NuxtLink
        >
        <NuxtLink
          to="/admin/dashboard"
          class="hidden lg:block btn btn-ghost text-xl dark:text-[#C9CF43]"
          >Dashboard</NuxtLink
        >
      </div>

      <!-- User Dropdown -->
      <div class="dropdown dropdown-end">
        <div tabindex="0" role="button" class="btn btn-ghost btn-circle avatar">
          <div class="w-10 rounded-full">
            <img
              :src="auth.user?.profile || '/images/Sample_User_Icon.png'"
              :alt="user?.username || 'User'"
            />
          </div>
        </div>
        <ul
          tabindex="0"
          class="menu menu-sm dropdown-content bg-base-100 rounded-box mt-3 w-52 p-2 shadow"
        >
          <li>
            <NuxtLink to="/profile"
              >Profile <span class="badge">New</span></NuxtLink
            >
          </li>
          <li><NuxtLink to="/">Home</NuxtLink></li>
          <li><NuxtLink to="/recipes">Recipes</NuxtLink></li>
          <li><NuxtLink to="/about">About Us</NuxtLink></li>
          <li><NuxtLink to="/admin/dashboard">Admin Dashboard</NuxtLink></li>
          <li><NuxtLink to="/admin/create-user">Create New User</NuxtLink></li>
          <li><NuxtLink to="/recipes/create">Let's Cook</NuxtLink></li>
          <li><button @click="handleLogout">Logout</button></li>
        </ul>
      </div>
    </div>
  </div>
</template>

<style scoped>
.navbar {
  transition: all 0.3s ease;
}
</style>