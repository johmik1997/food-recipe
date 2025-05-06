<template>
  <section class="bg-gradient-to-r from-green-50 to-primary-50 rounded-2xl shadow-sm">
    <div class="max-w-7xl mx-auto px-4 py-12 sm:px-6 lg:px-8">
      <div class="flex flex-col md:flex-row justify-between items-center gap-8">
        <!-- Welcome Message -->
        <div class="flex-1">
          <h1 class="text-3xl md:text-4xl font-bold text-gray-900 mb-2">
            Welcome back, <span class="text-primary-600">{{ user.name }}</span> 👋
          </h1>
          <p class="text-lg text-gray-600 mb-6">
            Let's create something delicious today!
          </p>
          
          <!-- Quick Actions -->
          <div class="flex flex-wrap gap-4">
            <NuxtLink 
              to="/dashboard/recipes/create" 
              class="inline-flex items-center px-6 py-3 bg-gradient-to-r from-green-600 to-primary-600 text-white font-medium rounded-lg shadow-md hover:from-green-700 hover:to-primary-700 transition-all transform hover:-translate-y-0.5"
            >
              <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
              </svg>
              Create New Recipe
            </NuxtLink>
            
            <NuxtLink 
              to="/dashboard/bookmarks" 
              class="inline-flex items-center px-6 py-3 bg-white text-gray-700 font-medium rounded-lg border border-gray-300 shadow-sm hover:bg-gray-50 transition-colors"
            >
              <svg class="w-5 h-5 mr-2 text-yellow-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z"></path>
              </svg>
              View Bookmarks
            </NuxtLink>
          </div>
        </div>
        
        <!-- Stats Cards -->
        <div class="grid grid-cols-2 gap-4 w-full md:w-auto">
          <div class="bg-white p-4 rounded-lg shadow-sm border border-gray-100 text-center">
            <p class="text-2xl font-bold text-primary-600">{{ user.recipeCount }}</p>
            <p class="text-sm text-gray-500">Recipes</p>
          </div>
          <div class="bg-white p-4 rounded-lg shadow-sm border border-gray-100 text-center">
            <p class="text-2xl font-bold text-primary-600">{{ user.likeCount }}</p>
            <p class="text-sm text-gray-500">Likes</p>
          </div>
          <div class="bg-white p-4 rounded-lg shadow-sm border border-gray-100 text-center">
            <p class="text-2xl font-bold text-primary-600">{{ user.bookmarkCount }}</p>
            <p class="text-sm text-gray-500">Bookmarks</p>
          </div>
          <div class="bg-white p-4 rounded-lg shadow-sm border border-gray-100 text-center">
            <p class="text-2xl font-bold text-primary-600">{{ user.followerCount }}</p>
            <p class="text-sm text-gray-500">Followers</p>
          </div>
        </div>
      </div>
      
      <!-- Recent Activity (Optional) -->
      <div v-if="recentActivity.length > 0" class="mt-12">
        <h3 class="text-lg font-medium text-gray-900 mb-4">Recent Activity</h3>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          <div 
            v-for="activity in recentActivity" 
            :key="activity.id" 
            class="bg-white p-4 rounded-lg shadow-sm border border-gray-100"
          >
            <div class="flex items-center">
              <img 
                :src="activity.user.avatar" 
                :alt="activity.user.name" 
                class="w-10 h-10 rounded-full mr-3"
              >
              <div>
                <p class="text-sm text-gray-600">
                  <span class="font-medium">{{ activity.user.name }}</span> {{ activity.action }}
                </p>
                <p class="text-xs text-gray-400">{{ formatDate(activity.date) }}</p>
              </div>
            </div>
            <NuxtLink 
              :to="`/recipes/${activity.recipe.id}`" 
              class="mt-3 block text-sm font-medium text-primary-600 hover:text-primary-700"
            >
              {{ activity.recipe.title }}
            </NuxtLink>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
// Example data - replace with actual data from your API
const user = {
  name: 'Sarah',
  recipeCount: 12,
  likeCount: 56,
  bookmarkCount: 8,
  followerCount: 24
}

const recentActivity = [
  {
    id: 1,
    user: {
      name: 'Alex',
      avatar: 'https://i.pravatar.cc/150?img=3'
    },
    action: 'liked your recipe',
    date: new Date(Date.now() - 3600000),
    recipe: {
      id: 'abc123',
      title: 'Spicy Thai Noodles'
    }
  },
  {
    id: 2,
    user: {
      name: 'Jamie',
      avatar: 'https://i.pravatar.cc/150?img=5'
    },
    action: 'commented on your recipe',
    date: new Date(Date.now() - 86400000),
    recipe: {
      id: 'def456',
      title: 'Chocolate Chip Cookies'
    }
  }
]

const formatDate = (date) => {
  return new Intl.DateTimeFormat('en-US', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  }).format(date)
}
</script>