<template>
  <div class="min-h-screen bg-gray-50">
    <Navbar />
    
    <div class="py-12">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <!-- Header -->
        <div class="text-center mb-16">
          <h1 class="text-4xl md:text-5xl font-bold text-gray-900 mb-4">
            <span class="bg-gradient-to-r from-green-600 to-primary-600 bg-clip-text text-transparent">
              My Bookmarks
            </span>
          </h1>
          <p class="max-w-2xl mx-auto text-lg text-gray-600">
            All your saved recipes in one place
          </p>
        </div>

        <!-- Loading State -->
        <div v-if="pending" class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          <div v-for="n in 3" :key="n" class="bg-white rounded-xl overflow-hidden shadow-sm border border-gray-100">
            <div class="animate-pulse">
              <div class="h-48 bg-gray-200"></div>
              <div class="p-6">
                <div class="h-6 bg-gray-200 rounded w-3/4 mb-4"></div>
                <div class="h-4 bg-gray-200 rounded w-full mb-2"></div>
                <div class="h-4 bg-gray-200 rounded w-5/6 mb-4"></div>
                <div class="flex justify-between">
                  <div class="h-4 bg-gray-200 rounded w-1/4"></div>
                  <div class="h-4 bg-gray-200 rounded w-1/4"></div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Error State -->
        <div v-else-if="error" class="text-center py-16 bg-white rounded-xl shadow-sm border border-gray-100">
          <div class="max-w-md mx-auto">
            <svg class="mx-auto h-16 w-16 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
            <h3 class="mt-4 text-lg font-medium text-gray-900">Failed to load bookmarks</h3>
            <p class="mt-2 text-sm text-gray-600">{{ error.message }}</p>
            <button 
              @click="refetch()"
              class="mt-6 inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
            >
              Retry
            </button>
          </div>
        </div>

        <!-- Empty State -->
        <div v-else-if="bookmarkedRecipes.length === 0" class="text-center py-16 bg-white rounded-xl shadow-sm border border-gray-100">
          <div class="max-w-md mx-auto">
            <svg class="mx-auto h-16 w-16 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
            </svg>
            <h3 class="mt-4 text-lg font-medium text-gray-900">No bookmarks yet</h3>
            <p class="mt-2 text-sm text-gray-600">
              Save your favorite recipes to find them here later
            </p>
            <NuxtLink 
              to="/recipes" 
              class="mt-6 inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
            >
              Browse Recipes
            </NuxtLink>
          </div>
        </div>

        <!-- Bookmarked Recipes Grid -->
        <div v-else class="grid gap-8 md:grid-cols-2 lg:grid-cols-3">
          <RecipeCard 
            v-for="recipe in bookmarkedRecipes" 
            :key="recipe.id"
            :recipe="recipe"
            :show-actions="true"
            @update="handleRecipeUpdate"
          />
        </div>
      </div>
    </div>
    <Footer />
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'

// Get user ID
const userId = ref(null)
onMounted(() => {
  const userStr = localStorage.getItem("user")
  if (userStr) {
    try {
      const user = JSON.parse(userStr)
      userId.value = user.id
    } catch (error) {
      console.error("Failed to parse user data", error)
    }
  }
})

// GraphQL query to get bookmarked recipes
const BOOKMARKED_RECIPES = gql`
  query BookmarkedRecipes($userId: uuid!) {
    user_bookmarks(where: {user_id: {_eq: $userId}}) {
      recipe {
        id
        title
        description
        prep_time
        cook_time
        total_time
        servings
        featured_image_url
        created_at
        user {
          id
          name
          avatar_image_url
        }
        user_bookmarks(where: {user_id: {_eq: $userId}}) {
          id
        }
        user_likes(where: {user_id: {_eq: $userId}}) {
          id
        }
        user_likes_aggregate {
          aggregate {
            count
          }
        }
      }
    }
  }
`

// Fetch bookmarked recipes
const { result, loading: pending, error, refetch } = useQuery(
  BOOKMARKED_RECIPES,
  () => ({
    userId: userId.value
  }),
  () => ({
    enabled: !!userId.value,
    fetchPolicy: 'cache-and-network'
  })
)

// Process bookmarked recipes data
const bookmarkedRecipes = computed(() => {
  if (!result.value?.user_bookmarks) return []
  
  return result.value.user_bookmarks.map(bookmark => {
    const recipe = bookmark.recipe
    return {
      ...recipe,
      is_bookmarked: true,
      is_liked: recipe.user_likes?.length > 0,
      likes_count: recipe.user_likes_aggregate?.aggregate?.count || 0
    }
  })
})

// Handle recipe updates from child components
const handleRecipeUpdate = (updatedRecipe) => {
  // If a recipe is unbookmarked, we should refetch to update the list
  if (updatedRecipe.is_bookmarked === false) {
    refetch()
  }
}
</script>