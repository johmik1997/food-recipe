<template>
  <nav class="w-full bg-white shadow px-4 py-3 flex items-center justify-between sticky top-0 z-50">
    <!-- Left: Logo and Main Links -->
    <div class="flex items-center gap-4">
      <NuxtLink to="/" class="text-xl font-bold text-green-600 flex items-center gap-2">
        <span>🍲</span>
        <span class="hidden sm:inline">Kushna</span>
      </NuxtLink>

      <div class="hidden md:flex items-center gap-3">
        <NuxtLink 
          to="/recipes" 
          class="hover:text-green-600 transition-colors text-sm sm:text-base"
          active-class="text-green-600 font-medium"
        >
          Browse Recipes
        </NuxtLink>

        <NuxtLink 
          to="/categories" 
          class="hover:text-green-600 transition-colors text-sm sm:text-base"
          active-class="text-green-600 font-medium"
        >
          Categories
        </NuxtLink>
      </div>
    </div>

    <!-- Right: Auth or User Actions -->
    <div class="flex items-center gap-3">
      <template v-if="isAuthenticated">
        <!-- My Recipes Link -->
        <NuxtLink 
          to="/dashboard/my-recipe" 
          class="hover:text-green-600 transition-colors hidden md:block text-sm sm:text-base"
          active-class="text-green-600 font-medium"
        >
          My Recipes
        </NuxtLink>
        
        <!-- Bookmarks Link -->
        <NuxtLink 
          to="/dashboard/bookmarks" 
          class="hover:text-green-600 transition-colors hidden md:block text-sm sm:text-base"
          active-class="text-green-600 font-medium"
        >
          Bookmarks
        </NuxtLink>

        <!-- Mobile Menu Button -->
        <button @click="toggleMobileMenu" class="md:hidden text-gray-600 hover:text-green-600 p-1">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
          </svg>
        </button>

        <!-- Avatar Dropdown -->
        <div class="relative hidden md:block" @click.stop="toggleDropdown">
          <div class="flex items-center gap-1 cursor-pointer">
            <img
              :src="userAvatar || '/default-avatar.png'"
              alt="User Avatar"
              class="w-8 h-8 rounded-full"
            />
            <span class="text-sm font-medium hidden lg:inline">Hi, {{ userName }}</span>
            <svg 
              xmlns="http://www.w3.org/2000/svg" 
              class="h-4 w-4 transition-transform duration-200"
              :class="{ 'rotate-180': dropdownOpen }"
              fill="none" 
              viewBox="0 0 24 24" 
              stroke="currentColor"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
            </svg>
          </div>
          
          <!-- Dropdown Menu -->
          <div
            v-if="dropdownOpen"
            class="absolute right-0 mt-2 w-48 bg-white border border-gray-200 rounded-lg shadow-lg z-50 py-1 text-sm"
            @click.stop
          >
            <NuxtLink 
              to="/dashboard" 
              class="block px-4 py-2 hover:bg-gray-50 text-gray-700"
              @click="dropdownOpen = false"
            >
              Dashboard
            </NuxtLink>
            <NuxtLink 
              to="/dashboard/purchases" 
              class="block px-4 py-2 hover:bg-gray-50 text-gray-700"
              @click="dropdownOpen = false"
            >
              Purchased Recipes
            </NuxtLink>
            <NuxtLink 
              to="/dashboard/updateProfile" 
              class="block px-4 py-2 hover:bg-gray-50 text-gray-700"
              @click="dropdownOpen = false"
            >
              Update Profile
            </NuxtLink>
            <div class="border-t border-gray-200 my-1"></div>
            <button 
              @click="logout"
              class="w-full text-left px-4 py-2 hover:bg-gray-50 text-gray-700"
            >
              Logout
            </button>
          </div>
        </div>
      </template>

      <template v-else>
        <!-- Login Link -->
        <NuxtLink 
          to="/login" 
          class="hover:text-green-600 transition-colors hidden md:block text-sm sm:text-base"
          active-class="text-green-600 font-medium"
        >
          Login
        </NuxtLink>
        
        <!-- Sign Up Button -->
        <NuxtLink 
          to="/signup" 
          class="bg-green-600 text-white px-3 py-1.5 sm:px-4 sm:py-2 rounded-lg hover:bg-green-700 transition-colors hidden md:block text-sm sm:text-base"
        >
          Sign Up
        </NuxtLink>

        <!-- Mobile Menu Button for non-auth users -->
        <button @click="toggleMobileMenu" class="md:hidden text-gray-600 hover:text-green-600 p-1">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
          </svg>
        </button>
      </template>
    </div>

    <!-- Mobile Menu Overlay -->
    <div 
      v-if="mobileMenuOpen" 
      class="fixed inset-0 bg-black bg-opacity-50 z-40 md:hidden"
      @click="mobileMenuOpen = false"
    ></div>

    <!-- Mobile Menu Sidebar -->
    <div 
      v-if="mobileMenuOpen" 
      class="fixed top-0 right-0 h-full w-64 bg-white z-50 shadow-lg transform transition-transform duration-300 ease-in-out"
      :class="{ 'translate-x-0': mobileMenuOpen, 'translate-x-full': !mobileMenuOpen }"
    >
      <div class="flex justify-between items-center p-4 border-b">
        <NuxtLink to="/" class="text-xl font-bold text-green-600" @click="mobileMenuOpen = false">
          Kushna
        </NuxtLink>
        <button @click="mobileMenuOpen = false" class="text-gray-500 hover:text-gray-700">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <div class="p-4 space-y-2">
        <!-- Common Links -->
        <NuxtLink 
          to="/recipes" 
          class="block py-2 px-3 rounded hover:bg-gray-100 text-gray-700"
          active-class="text-green-600 font-medium bg-green-50"
          @click="mobileMenuOpen = false"
        >
          Browse Recipes
        </NuxtLink>
        <NuxtLink 
          to="/categories" 
          class="block py-2 px-3 rounded hover:bg-gray-100 text-gray-700"
          active-class="text-green-600 font-medium bg-green-50"
          @click="mobileMenuOpen = false"
        >
          Categories
        </NuxtLink>

        <template v-if="isAuthenticated">
          <!-- User Links -->
          <div class="pt-2 mt-2 border-t">
            <NuxtLink 
              to="/dashboard/my-recipe" 
              class="block py-2 px-3 rounded hover:bg-gray-100 text-gray-700"
              active-class="text-green-600 font-medium bg-green-50"
              @click="mobileMenuOpen = false"
            >
              My Recipes
            </NuxtLink>
            <NuxtLink 
              to="/dashboard/bookmarks" 
              class="block py-2 px-3 rounded hover:bg-gray-100 text-gray-700"
              active-class="text-green-600 font-medium bg-green-50"
              @click="mobileMenuOpen = false"
            >
              Bookmarks
            </NuxtLink>
            <NuxtLink 
              to="/dashboard" 
              class="block py-2 px-3 rounded hover:bg-gray-100 text-gray-700"
              active-class="text-green-600 font-medium bg-green-50"
              @click="mobileMenuOpen = false"
            >
              Dashboard
            </NuxtLink>
            <NuxtLink 
              to="/dashboard/purchases" 
              class="block py-2 px-3 rounded hover:bg-gray-100 text-gray-700"
              active-class="text-green-600 font-medium bg-green-50"
              @click="mobileMenuOpen = false"
            >
              Purchased Recipes
            </NuxtLink>
            <button 
              @click="logout"
              class="w-full text-left py-2 px-3 rounded hover:bg-gray-100 text-red-600"
            >
              Logout
            </button>
          </div>
        </template>
        <template v-else>
          <!-- Auth Links -->
          <div class="pt-2 mt-2 border-t">
            <NuxtLink 
              to="/login" 
              class="block py-2 px-3 rounded hover:bg-gray-100 text-gray-700"
              active-class="text-green-600 font-medium bg-green-50"
              @click="mobileMenuOpen = false"
            >
              Login
            </NuxtLink>
            <NuxtLink 
              to="/signup" 
              class="block py-2 px-3 rounded bg-green-600 text-white hover:bg-green-700"
              @click="mobileMenuOpen = false"
            >
              Sign Up
            </NuxtLink>
          </div>
        </template>
      </div>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'

// Auth state
const isAuthenticated = ref(false)
const userId = ref<string | null>(null)
const userName = ref<string>('User')
const userAvatar = ref<string>('')

// UI states
const dropdownOpen = ref(false)
const mobileMenuOpen = ref(false)

// GraphQL query
const { result } = useQuery(gql`
  query getUserName($userId: uuid!) {
    users_by_pk(id: $userId) {
      name
      avatar_image_url
    }
  }
`, 
{ userId },
() => ({
  enabled: !!userId.value
}))

// Watch for user data changes
watch(result, (newResult) => {
  if (newResult?.users_by_pk?.name) {
    userName.value = newResult.users_by_pk.name
    userAvatar.value = newResult.users_by_pk.avatar_image_url || ''
  }
})

// Initialize auth state
onMounted(() => {
  const userStr = localStorage.getItem("user")
  if (userStr) {
    try {
      const user = JSON.parse(userStr)
      userId.value = user.id
      isAuthenticated.value = true
    } catch (e) {
      console.error("Failed to parse user data", e)
    }
  }
})

// Methods
const toggleDropdown = () => {
  dropdownOpen.value = !dropdownOpen.value
}

const toggleMobileMenu = () => {
  mobileMenuOpen.value = !mobileMenuOpen.value
  dropdownOpen.value = false
  // Prevent body scroll when menu is open
  if (mobileMenuOpen.value) {
    document.body.style.overflow = 'hidden'
  } else {
    document.body.style.overflow = ''
  }
}

const logout = () => {
  localStorage.removeItem("user")
  isAuthenticated.value = false
  dropdownOpen.value = false
  mobileMenuOpen.value = false
  document.body.style.overflow = ''
  navigateTo('/login')
}

// Close dropdown when clicking outside
const handleClickOutside = (event: MouseEvent) => {
  if (dropdownOpen.value) {
    dropdownOpen.value = false
  }
}

// Set up click outside listener
onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside)
  document.body.style.overflow = ''
})
</script>

<style scoped>
/* Transition for dropdown */
.dropdown-enter-active,
.dropdown-leave-active {
  transition: all 0.2s ease;
}
.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

/* Mobile menu transition */
.slide-enter-active,
.slide-leave-active {
  transition: transform 0.3s ease;
}
.slide-enter-from,
.slide-leave-to {
  transform: translateX(100%);
}
</style>