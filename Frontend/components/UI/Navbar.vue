<script setup>
import { ref, computed } from "vue";
import { useToast } from "vue-toast-notification";

const searchQuery = ref("");

const toast = useToast();
const auth = authStore();
const router = useRouter();

onMounted(async () => {
  await auth.getProfile();
});

const id = auth.userId;
const isAuthenticated = computed(() => auth.isAuthed);
const isAdmin = computed(() => auth.isAdmin);
console.log("the user is authenticated?", isAuthenticated);
const user = computed(() => auth.user);

const handleLogout = async () => {
  const result = await auth.logout();
  if (result) {
    toast.success("You have logged out successfully");
    router.push("/");
  } else {
    toast.error("Unable to logout");
  }
};
</script>

<template>
  <div class="sticky top-0 z-10 bg-white shadow-sm">
    <div class="container mx-auto px-4">
      <nav class="flex items-center justify-between h-16">
        <!-- Logo/Brand -->
        <NuxtLink 
          to="/" 
          class="flex items-center space-x-2 text-green-700 hover:text-green-800"
        >
          <span class="text-2xl">🍲</span>
          <span class="text-xl font-bold hidden sm:block">Kushna</span>
        </NuxtLink>

        <!-- Desktop Navigation -->
        <div class="hidden md:flex items-center space-x-6">
          <NuxtLink 
            v-if="isAuthenticated"
            to="/recipes" 
            class="text-green-700 hover:text-green-800 font-medium"
          >
            Recipes
          </NuxtLink>
          
          <NuxtLink 
            v-if="isAuthenticated"
            to="/recipes/create" 
            class="text-green-700 hover:text-green-800 font-medium"
          >
            Create Recipe
        </NuxtLink>
          
          <NuxtLink 
            v-if="isAdmin"
            to="/admin/dashboard" 
            class="text-green-700 hover:text-green-800 font-medium"
          >
            Dashboard
          </NuxtLink>
          
          <NuxtLink 
            v-if="!isAuthenticated"
            to="/auth" 
            class="text-green-700 hover:text-green-800 font-medium"
          >
            Signup/Login
          </NuxtLink>
        </div>

        <!-- User Profile Dropdown -->
        <div class="flex items-center space-x-4">
          <div v-if="isAuthenticated" class="dropdown dropdown-end">
            <div tabindex="0" class="flex items-center space-x-2 cursor-pointer">
              <div class="avatar">
                <div class="w-10 h-10 rounded-full border-2 border-green-600">
                  <img 
                    :src="auth.$state.user?.profile || '/images/Sample_User_Icon.png'"
                    :alt="auth.$state.user ? `${auth.$state.user.username}'s profile` : 'Default profile'"
                  />
                </div>
              </div>
            </div>
            
            <ul 
              tabindex="0" 
              class="dropdown-content menu p-2 shadow bg-white rounded-box w-52 border border-green-100 mt-2"
            >
              <li v-if="isAuthenticated">
                <NuxtLink :to="`/profile`" class="text-green-700 hover:bg-green-50">
                  Profile
                </NuxtLink>
              </li>
              <li v-if="isAuthenticated">
                <NuxtLink :to="`/recipes`" class="text-green-700 hover:bg-green-50">
                  Recipes
                </NuxtLink>
              </li>
              <li v-if="isAdmin">
                <NuxtLink 
                  :to="`/admin/dashboard`" 
                  class="text-green-700 hover:bg-green-50"
                >
                  Admin Dashboard
                </NuxtLink>
              </li>
              <li v-if="isAdmin">
                <NuxtLink 
                  :to="`/admin/create-user`" 
                  class="text-green-700 hover:bg-green-50"
                >
                  Create New User
                </NuxtLink>
              </li>
              <li v-if="isAuthenticated">
                <NuxtLink 
                  to="/recipes/create" 
                  class="text-green-700 hover:bg-green-50"
                >
                  Let's Cook
                </NuxtLink>
              </li>
              <li v-if="isAuthenticated">
                <button 
                  @click="handleLogout" 
                  class="text-green-700 hover:bg-green-50 text-left"
                >
                  Logout
                </button>
              </li>
              <li v-if="!isAuthenticated">
                <NuxtLink 
                  :to="`/auth`" 
                  class="text-green-700 hover:bg-green-50"
                >
                  Login/Signup
                </NuxtLink>
              </li>
            </ul>
          </div>
        </div>
      </nav>
    </div>
  </div>
</template>

<style scoped>
/* Custom styles */
nav {
  transition: all 0.3s ease;
}

.dropdown-content {
  display: none;
}

.dropdown:hover .dropdown-content,
.dropdown:focus-within .dropdown-content {
  display: block;
}

.avatar img {
  transition: transform 0.3s ease;
}

.avatar:hover img {
  transform: scale(1.05);
}

/* Active link styling */
.router-link-active {
  @apply text-green-800 font-semibold;
}

/* Mobile menu styles */
@media (max-width: 767px) {
  .mobile-menu {
    @apply absolute top-16 left-0 right-0 bg-white shadow-lg py-2 px-4;
  }
}
</style>