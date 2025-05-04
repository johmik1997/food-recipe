<template>
    <div>
      <div class="mb-8">
        <h1 class="text-2xl font-bold text-gray-900">Welcome back, {{ user.name }}!</h1>
        <p class="mt-1 text-gray-600">Discover, create, and share delicious recipes</p>
      </div>
      
      <div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
        <!-- Stats Cards -->
        <div class="p-6 bg-white rounded-lg shadow">
          <div class="flex items-center">
            <div class="p-3 rounded-full bg-primary-50 text-primary-600">
              <Icon name="mdi:book" class="w-6 h-6" />
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-500">My Recipes</p>
              <p class="text-2xl font-semibold text-gray-900">24</p>
            </div>
          </div>
        </div>
        
        <div class="p-6 bg-white rounded-lg shadow">
          <div class="flex items-center">
            <div class="p-3 rounded-full bg-green-50 text-green-600">
              <Icon name="mdi:bookmark" class="w-6 h-6" />
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-500">Bookmarks</p>
              <p class="text-2xl font-semibold text-gray-900">42</p>
            </div>
          </div>
        </div>
        
        <div class="p-6 bg-white rounded-lg shadow">
          <div class="flex items-center">
            <div class="p-3 rounded-full bg-blue-50 text-blue-600">
              <Icon name="mdi:heart" class="w-6 h-6" />
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-500">Likes Received</p>
              <p class="text-2xl font-semibold text-gray-900">156</p>
            </div>
          </div>
        </div>
        
        <div class="p-6 bg-white rounded-lg shadow">
          <div class="flex items-center">
            <div class="p-3 rounded-full bg-purple-50 text-purple-600">
              <Icon name="mdi:comment" class="w-6 h-6" />
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-500">Comments</p>
              <p class="text-2xl font-semibold text-gray-900">32</p>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Recent Activity -->
      <div class="mt-8">
        <h2 class="text-lg font-medium text-gray-900">Recent Activity</h2>
        <div class="mt-4 overflow-hidden bg-white shadow sm:rounded-md">
          <ul class="divide-y divide-gray-200">
            <li v-for="activity in recentActivities" :key="activity.id">
              <div class="px-4 py-4 sm:px-6">
                <div class="flex items-center justify-between">
                  <p class="text-sm font-medium text-primary-600 truncate">{{ activity.action }}</p>
                  <div class="flex flex-shrink-0 ml-2">
                    <p class="inline-flex px-2 text-xs font-semibold leading-5 text-gray-500">
                      {{ activity.time }}
                    </p>
                  </div>
                </div>
                <div class="mt-2 sm:flex sm:justify-between">
                  <div class="sm:flex">
                    <p class="flex items-center text-sm text-gray-500">
                      <Icon :name="activity.icon" class="flex-shrink-0 mr-1.5 h-5 w-5 text-gray-400" />
                      {{ activity.details }}
                    </p>
                  </div>
                </div>
              </div>
            </li>
          </ul>
        </div>
      </div>
      
      <!-- Quick Actions -->
      <div class="mt-8">
        <h2 class="text-lg font-medium text-gray-900">Quick Actions</h2>
        <div class="grid grid-cols-1 gap-4 mt-4 sm:grid-cols-2">
          <NuxtLink 
            to="/addrecipe" 
            class="flex items-center p-4 bg-white rounded-lg shadow hover:bg-gray-50"
          >
            <div class="p-3 rounded-full bg-primary-100 text-primary-600">
              <Icon name="mdi:plus-circle" class="w-6 h-6" />
            </div>
            <div class="ml-4">
              <h3 class="text-sm font-medium text-gray-900">Add New Recipe</h3>
              <p class="text-sm text-gray-500">Share your culinary creation</p>
            </div>
          </NuxtLink>
          
          <NuxtLink 
            to="/dashboard/bookmarks" 
            class="flex items-center p-4 bg-white rounded-lg shadow hover:bg-gray-50"
          >
            <div class="p-3 rounded-full bg-green-100 text-green-600">
              <Icon name="mdi:bookmark" class="w-6 h-6" />
            </div>
            <div class="ml-4">
              <h3 class="text-sm font-medium text-gray-900">View Bookmarks</h3>
              <p class="text-sm text-gray-500">Your saved recipes</p>
            </div>
          </NuxtLink>
        </div>
      </div>
    </div>
  </template>
  
  <script setup>
  const user = ref({
    name: 'John Doe'
  });
  
  const recentActivities = ref([
    {
      id: 1,
      action: 'New recipe added',
      details: 'Homemade Pizza',
      time: '2 hours ago',
      icon: 'mdi:book-plus'
    },
    {
      id: 2,
      action: 'Recipe liked',
      details: 'Chocolate Cake by Sarah',
      time: '5 hours ago',
      icon: 'mdi:heart'
    },
    {
      id: 3,
      action: 'New comment',
      details: 'On your Pasta recipe',
      time: '1 day ago',
      icon: 'mdi:comment'
    },
    {
      id: 4,
      action: 'Bookmark added',
      details: 'Vegetable Curry',
      time: '2 days ago',
      icon: 'mdi:bookmark'
    }
  ]);
  </script>