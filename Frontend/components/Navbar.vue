<template>
    <nav class="w-full bg-white shadow px-6 py-3 flex items-center justify-between sticky top-0 z-50">
      <!-- Left: Logo and Main Links -->
      <div class="flex items-center gap-6">
        <NuxtLink to="/" class="text-xl font-bold text-green-600 flex items-center gap-2">
          <span>🍲</span>
          <span>FoodShare</span>
        </NuxtLink>
        
        <NuxtLink 
          to="/recipes" 
          class="hover:text-green-600 transition-colors"
          active-class="text-green-600 font-medium"
        >
          Browse Recipes
        </NuxtLink>
        
        <NuxtLink 
          to="/categories" 
          class="hover:text-green-600 transition-colors"
          active-class="text-green-600 font-medium"
        >
          Categories
        </NuxtLink>
      </div>
  
      <!-- Center: Search Bar -->
      <div class="flex-1 max-w-xl mx-6">
        <div class="relative">
          <input
            type="text"
            placeholder="Search recipes..."
            class="w-full py-2 px-4 border border-gray-300 rounded-full focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-transparent"
          >
          <button class="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-green-600">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </button>
        </div>
      </div>
  
      <!-- Right: Auth or User Actions -->
      <div class="flex items-center gap-4">
        <template v-if="isAuthenticated">
          <!-- My Recipes Link -->
          <NuxtLink 
            to="/dashboard/my-recipe" 
            class="hover:text-green-600 transition-colors hidden md:block"
            active-class="text-green-600 font-medium"
          >
            My Recipes
          </NuxtLink>
          
          <!-- Bookmarks Link -->
          <NuxtLink 
            to="/dashboard/bookmarks" 
            class="hover:text-green-600 transition-colors hidden md:block"
            active-class="text-green-600 font-medium"
          >
            Bookmarks
          </NuxtLink>
  
          <!-- Mobile Menu Button (hidden on desktop) -->
          <button @click="toggleMobileMenu" class="md:hidden text-gray-600 hover:text-green-600">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
            </svg>
          </button>
  
          <!-- Avatar Dropdown -->
          <div class="relative hidden md:block" @click.stop="toggleDropdown">
            <div class="flex items-center gap-2 cursor-pointer">
              <img
                src="https://i.pravatar.cc/32"
                alt="User Avatar"
                class="w-8 h-8 rounded-full"
              />
              <span class="text-sm font-medium">Hi,{{ userName }}</span>
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
              class="absolute right-0 mt-2 w-56 bg-white border border-gray-200 rounded-lg shadow-lg z-50 py-1"
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
               update profile
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
            class="hover:text-green-600 transition-colors hidden md:block"
            active-class="text-green-600 font-medium"
          >
            Login
          </NuxtLink>
          
          <!-- Sign Up Button -->
          <NuxtLink 
            to="/signup" 
            class="bg-green-600 text-white px-4 py-2 rounded-lg hover:bg-green-700 transition-colors hidden md:block"
          >
            Sign Up
          </NuxtLink>
        </template>
      </div>
  
      <!-- Mobile Menu (shown only on mobile when toggled) -->
      <div 
        v-if="mobileMenuOpen" 
        class="fixed inset-0 bg-white z-40 flex flex-col pt-20 px-6 md:hidden"
        @click="mobileMenuOpen = false"
      >
        <div class="space-y-4">
          <template v-if="isAuthenticated">
            <NuxtLink 
              to="/dashboard/my-recipe" 
              class="block py-3 text-lg border-b border-gray-100"
              @click="mobileMenuOpen = false"
            >
              My Recipes
            </NuxtLink>
            <NuxtLink 
              to="/dashboard/bookmarks" 
              class="block py-3 text-lg border-b border-gray-100"
              @click="mobileMenuOpen = false"
            >
              Bookmarks
            </NuxtLink>
            <NuxtLink 
              to="/dashboard" 
              class="block py-3 text-lg border-b border-gray-100"
              @click="mobileMenuOpen = false"
            >
              Dashboard
            </NuxtLink>
            <NuxtLink 
              to="/dashboard/purchases" 
              class="block py-3 text-lg border-b border-gray-100"
              @click="mobileMenuOpen = false"
            >
              Purchased Recipes
            </NuxtLink>
            <button 
              @click="logout"
              class="w-full text-left py-3 text-lg border-b border-gray-100 text-red-600"
            >
              Logout
            </button>
          </template>
          <template v-else>
            <NuxtLink 
              to="/login" 
              class="block py-3 text-lg border-b border-gray-100"
              @click="mobileMenuOpen = false"
            >
              Login
            </NuxtLink>
            <NuxtLink 
              to="/signup" 
              class="block py-3 text-lg border-b border-gray-100 text-green-600 font-medium"
              @click="mobileMenuOpen = false"
            >
              Sign Up
            </NuxtLink>
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

// UI states
const dropdownOpen = ref(false)
const mobileMenuOpen = ref(false)

// GraphQL query
const { result } = useQuery(gql`
  query getUserName($userId: uuid!) {
    users_by_pk(id: $userId) {
      name
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
}

const logout = () => {
  localStorage.removeItem("user")
  isAuthenticated.value = false
  dropdownOpen.value = false
  mobileMenuOpen.value = false
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
  </style>