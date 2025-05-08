<template>
  <div>
    <Navbar />
    <section class="bg-gradient-to-r from-green-50 to-primary-50 rounded-2xl shadow-sm">
      <div class="max-w-7xl mx-auto px-4 py-12 sm:px-6 lg:px-8">
        <!-- Error State -->
        <div v-if="error" class="text-center py-12 text-red-600">
          {{ error }}
        </div>

        <!-- Loading State -->
        <div v-else-if="loading" class="text-center py-12">
          <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary-600 mx-auto mb-4"></div>
          <p class="text-gray-600">Loading your delicious dashboard...</p>
        </div>

        <!-- Content -->
        <div v-else>
          <div class="flex flex-col md:flex-row justify-between items-center gap-8">
            <!-- Welcome Message -->
            <div class="flex-1">
              <h1 class="text-3xl md:text-4xl font-bold text-gray-900 mb-2">
                Welcome back, <span class="text-primary-600">{{ user.name }}</span> 👋
              </h1>
              <p class="text-lg text-gray-600 mb-6">
                {{ welcomeMessage }}
              </p>
              
              <!-- Quick Actions -->
              <div class="flex flex-wrap gap-4">
                <NuxtLink 
                  to="/dashboard/recipes/create" 
                  class="inline-flex items-center px-6 py-3 bg-gradient-to-r from-green-600 to-primary-600 text-white font-medium rounded-lg shadow-md hover:from-green-700 hover:to-primary-700 transition-all transform hover:-translate-y-0.5"
                >
                  <PlusIcon class="w-5 h-5 mr-2" />
                  Create New Recipe
                </NuxtLink>
                
                <NuxtLink 
                  to="/dashboard/my-recipes" 
                  class="inline-flex items-center px-6 py-3 bg-white text-gray-700 font-medium rounded-lg border border-gray-300 shadow-sm hover:bg-gray-50 transition-colors"
                >
                  <BookOpenIcon class="w-5 h-5 mr-2 text-primary-500" />
                  My Recipes
                </NuxtLink>
              </div>
            </div>
            
            <!-- Stats Cards -->
            <div class="grid grid-cols-2 gap-4 w-full md:w-auto">
              <div 
                v-for="stat in stats" 
                :key="stat.label"
                class="bg-white p-4 rounded-lg shadow-sm border border-gray-100 text-center hover:shadow-md transition-shadow cursor-pointer"
                @click="handleStatClick(stat.action)"
              >
                <p class="text-2xl font-bold text-primary-600">{{ stat.value }}</p>
                <p class="text-sm text-gray-500">{{ stat.label }}</p>
              </div>
            </div>
          </div>
          
          <!-- Recent Recipes Section -->
          <div v-if="recentRecipes.length > 0" class="mt-12">
            <div class="flex justify-between items-center mb-6">
              <h3 class="text-xl font-bold text-gray-900">Your Recent Recipes</h3>
              <NuxtLink 
                to="/dashboard/recipes" 
                class="text-sm text-primary-600 hover:text-primary-700 font-medium"
              >
                View All Recipes →
              </NuxtLink>
            </div>
            <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
              <RecipeCard 
                v-for="recipe in recentRecipes" 
                :key="recipe.id"
                :recipe="recipe"
                :show-actions="true"
                @edit="handleEditRecipe(recipe.id)"
                @delete="handleDeleteRecipe(recipe.id)"
              />
            </div>
          </div>

          <!-- Empty State -->
          <div v-if="recentRecipes.length === 0" class="mt-12 text-center py-12">
            <div class="mx-auto w-24 h-24 bg-gray-100 rounded-full flex items-center justify-center mb-4">
              <UtensilsIcon class="w-10 h-10 text-gray-400" />
            </div>
            <h3 class="text-lg font-medium text-gray-900 mb-2">No Recipes Yet</h3>
            <p class="text-gray-600 mb-6">Start by creating your first delicious recipe!</p>
            <NuxtLink 
              to="/dashboard/recipes/create" 
              class="inline-flex items-center px-6 py-3 bg-gradient-to-r from-green-600 to-primary-600 text-white font-medium rounded-lg shadow-md hover:from-green-700 hover:to-primary-700 transition-all"
            >
              <PlusIcon class="w-5 h-5 mr-2" />
              Create Your First Recipe
            </NuxtLink>
          </div>
        </div>
      </div>
    </section>
    <Footer />
  </div>
</template>

<script setup>
import { PlusIcon, BookOpenIcon} from '@heroicons/vue/outline'

// Define components (if using Nuxt auto-imports, you can remove these)
const RecipeCard = defineAsyncComponent(() => import('@/components/RecipeCard.vue'))
const Navbar = defineAsyncComponent(() => import('@/components/Navbar.vue'))
const Footer = defineAsyncComponent(() => import('@/components/Footer.vue'))

const userId = ref(null)
const userName = ref('Chef')
const recipeCount = ref(0)
const loading = ref(true)
const error = ref(null)
const recentRecipes = ref([])

// Set userId from localStorage
onMounted(() => {
  const userStr = localStorage.getItem("user")
  if (userStr) {
    try {
      const user = JSON.parse(userStr)
      userId.value = user.id
      // Simulate loading recent recipes
      setTimeout(() => {
        recentRecipes.value = [
          {
            id: '1',
            title: 'Sample Recipe 1',
            image: '/placeholder-recipe.jpg',
            rating: 4.5,
            prepTime: '30 mins',
            difficulty: 'Medium'
          },
          {
            id: '2',
            title: 'Sample Recipe 2',
            image: '/placeholder-recipe.jpg',
            rating: 4.0,
            prepTime: '45 mins',
            difficulty: 'Easy'
          }
        ]
        loading.value = false
      }, 1500)
    } catch (e) {
      error.value = 'Failed to parse user data'
      loading.value = false
    }
  } else {
    error.value = 'No user data found'
    loading.value = false
  }
})

// User object
const user = computed(() => ({
  name: userName.value,
  recipeCount: recipeCount.value,
  likeCount: 42, // Example data
  bookmarkCount: 18, // Example data
  followerCount: 156 // Example data
}))

// Welcome message
const welcomeMessage = computed(() => {
  const hour = new Date().getHours()
  if (hour < 12) return 'Good morning! What will you cook today?'
  if (hour < 18) return 'Good afternoon! Ready to create something tasty?'
  return 'Good evening! Time to plan tomorrow\'s meals!'
})

// Stats
const stats = computed(() => [
  { label: 'Recipes', value: user.value.recipeCount, action: 'viewRecipes' },
  { label: 'Likes', value: user.value.likeCount, action: 'viewLikes' },
  { label: 'Bookmarks', value: user.value.bookmarkCount, action: 'viewBookmarks' },
  { label: 'Followers', value: user.value.followerCount, action: 'viewFollowers' }
])

// Handle stat card clicks
const handleStatClick = (action) => {
  switch (action) {
    case 'viewRecipes':
      navigateTo('/dashboard/my-recipes')
      break
    case 'viewLikes':
      navigateTo('/dashboard/likes')
      break
    case 'viewBookmarks':
      navigateTo('/dashboard/bookmarks')
      break
    case 'viewFollowers':
      navigateTo('/dashboard/followers')
      break
  }
}

// Recipe actions
const handleEditRecipe = (id) => {
  navigateTo(`/dashboard/recipes/edit/${id}`)
}

const handleDeleteRecipe = (id) => {
  // Implement delete logic
  recentRecipes.value = recentRecipes.value.filter(recipe => recipe.id !== id)
}
</script>

<style scoped>
/* Custom transitions for smoother hover effects */
.recipe-card-enter-active,
.recipe-card-leave-active {
  transition: all 0.3s ease;
}
.recipe-card-enter-from,
.recipe-card-leave-to {
  opacity: 0;
  transform: translateY(10px);
}
</style>