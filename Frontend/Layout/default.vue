<template>
    <div class="min-h-screen bg-gray-50">
      <!-- Navigation Sidebar -->
      <div class="fixed inset-y-0 left-0 w-64 bg-white shadow-lg">
        <div class="flex items-center justify-center h-16 px-4 bg-primary-600">
          <h1 class="text-white font-bold text-xl">FoodieRecipes</h1>
        </div>
        <nav class="mt-6">
          <div class="px-4 space-y-1">
            <NuxtLink 
              to="/dashboard" 
              class="flex items-center px-4 py-3 text-sm font-medium rounded-lg text-gray-700 hover:bg-gray-100 group"
              active-class="bg-primary-50 text-primary-600"
            >
              <Icon name="mdi:view-dashboard" class="w-5 h-5 mr-3 text-gray-400 group-hover:text-gray-500" />
              Dashboard
            </NuxtLink>
            
            <NuxtLink 
              to="/dashboard/recipes/add" 
              class="flex items-center px-4 py-3 text-sm font-medium rounded-lg text-gray-700 hover:bg-gray-100 group"
              active-class="bg-primary-50 text-primary-600"
            >
              <Icon name="mdi:plus-circle" class="w-5 h-5 mr-3 text-gray-400 group-hover:text-gray-500" />
              Add Recipe
            </NuxtLink>
            
            <NuxtLink 
              to="/dashboard/recipes/my" 
              class="flex items-center px-4 py-3 text-sm font-medium rounded-lg text-gray-700 hover:bg-gray-100 group"
              active-class="bg-primary-50 text-primary-600"
            >
              <Icon name="mdi:book" class="w-5 h-5 mr-3 text-gray-400 group-hover:text-gray-500" />
              My Recipes
            </NuxtLink>
            
            <NuxtLink 
              to="/dashboard/bookmarks" 
              class="flex items-center px-4 py-3 text-sm font-medium rounded-lg text-gray-700 hover:bg-gray-100 group"
              active-class="bg-primary-50 text-primary-600"
            >
              <Icon name="mdi:bookmark" class="w-5 h-5 mr-3 text-gray-400 group-hover:text-gray-500" />
              Bookmarks
            </NuxtLink>
            
            <NuxtLink 
              to="/dashboard/profile" 
              class="flex items-center px-4 py-3 text-sm font-medium rounded-lg text-gray-700 hover:bg-gray-100 group"
              active-class="bg-primary-50 text-primary-600"
            >
              <Icon name="mdi:account" class="w-5 h-5 mr-3 text-gray-400 group-hover:text-gray-500" />
              Profile
            </NuxtLink>
          </div>
        </nav>
        
        <div class="absolute bottom-0 left-0 right-0 p-4 border-t border-gray-200">
          <button 
            @click="logout"
            class="flex items-center w-full px-4 py-2 text-sm font-medium text-red-600 rounded-lg hover:bg-red-50"
          >
            <Icon name="mdi:logout" class="w-5 h-5 mr-3" />
            Sign Out
          </button>
        </div>
      </div>
      
      <!-- Main Content -->
      <div class="pl-64">
        <!-- Top Header -->
        <header class="flex items-center justify-between h-16 px-6 bg-white border-b border-gray-200">
          <div class="flex items-center">
            <h2 class="text-lg font-semibold text-gray-900">{{ pageTitle }}</h2>
          </div>
          
          <div class="flex items-center space-x-4">
            <button class="p-1 text-gray-400 rounded-full hover:text-gray-500 focus:outline-none">
              <Icon name="mdi:bell" class="w-6 h-6" />
            </button>
            
            <div class="relative">
              <button 
                @click="toggleUserMenu" 
                class="flex items-center focus:outline-none"
              >
                <img 
                  class="w-8 h-8 rounded-full" 
                  :src="user.avatar || '/default-avatar.png'" 
                  alt="User profile"
                >
                <span class="ml-2 text-sm font-medium text-gray-700">{{ user.name }}</span>
              </button>
              
              <div 
                v-if="isUserMenuOpen" 
                class="absolute right-0 z-10 w-48 mt-2 origin-top-right bg-white rounded-md shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none"
              >
                <div class="py-1">
                  <NuxtLink 
                    to="/dashboard/profile" 
                    class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                  >
                    Your Profile
                  </NuxtLink>
                  <button 
                    @click="logout"
                    class="block w-full px-4 py-2 text-sm text-left text-red-600 hover:bg-red-50"
                  >
                    Sign out
                  </button>
                </div>
              </div>
            </div>
          </div>
        </header>
        
        <!-- Page Content -->
        <main class="p-6">
          <slot />
        </main>
      </div>
    </div>
  </template>
  
  <script setup>
  import { ref } from 'vue';
  
  const pageTitle = ref('Dashboard');
  const isUserMenuOpen = ref(false);
  const user = ref({
    name: 'John Doe',
    avatar: null
  });
  
  const toggleUserMenu = () => {
    isUserMenuOpen.value = !isUserMenuOpen.value;
  };
  
  const logout = async () => {
    // Implement logout logic
    try {
      // Clear JWT token, etc.
      await navigateTo('/login');
    } catch (error) {
      console.error('Logout error:', error);
    }
  };
  </script>
  
  <style>
  /* Custom styles if needed */
  </style>