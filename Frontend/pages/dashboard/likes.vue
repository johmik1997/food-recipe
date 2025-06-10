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
          <p class="text-gray-600">Loading your liked recipes...</p>
        </div>

        <!-- Content -->
        <div v-else>
          <div class="flex justify-between items-center mb-8">
            <h1 class="text-2xl md:text-3xl font-bold text-gray-900">
              Your Liked Recipes
            </h1>
            <NuxtLink 
              to="/dashboard" 
              class="text-sm text-primary-600 hover:text-primary-700 font-medium"
            >
              ← Back to Dashboard
            </NuxtLink>
          </div>

          <!-- Liked Recipes Grid -->
          <div v-if="likedRecipes.length > 0" class="grid gap-8 md:grid-cols-2 lg:grid-cols-3">
            <RecipeCard 
              v-for="recipe in likedRecipes" 
              :key="recipe.id"
              :recipe="recipe"
              :show-actions="true"
              @update="handleRecipeUpdate"
            />
          </div>

          <!-- Empty State -->
          <div v-else class="text-center py-12">
            <div class="mx-auto w-24 h-24 bg-gray-100 rounded-full flex items-center justify-center mb-4">
              <Icon name="material-symbols:favorite-outline" class="w-10 h-10 text-gray-400" />
            </div>
            <h3 class="text-lg font-medium text-gray-900 mb-2">No Liked Recipes Yet</h3>
            <p class="text-gray-600 mb-6">Like some recipes to see them appear here!</p>
            <NuxtLink 
              to="/recipes" 
              class="inline-flex items-center px-6 py-3 bg-gradient-to-r from-green-600 to-primary-600 text-white font-medium rounded-lg shadow-md hover:from-green-700 hover:to-primary-700 transition-all"
            >
              <Icon name="heroicons:magnifying-glass" class="w-5 h-5 mr-2" />
              Browse Recipes
            </NuxtLink>
          </div>
        </div>
      </div>
    </section>
    <Footer />
  </div>
</template>

<script setup>
import { useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'

const userId = ref(null)
const loading = ref(true)
const error = ref(null)

// Set userId from localStorage
onBeforeMount(() => {
  try {
    const userStr = localStorage.getItem("user")
    if (userStr) {
      const user = JSON.parse(userStr)
      userId.value = user.id
    } else {
      throw new Error('No user data found')
    }
  } catch (e) {
    error.value = e.message
    loading.value = false
    navigateTo('/login')
  }
})

// GraphQL query to fetch liked recipes
const { result: likedRecipesResult, loading: queryLoading } = useQuery(gql`
  query GetLikedRecipes($userId: uuid!) {
    user_likes(where: { user_id: { _eq: $userId } }) {
      recipe {
        id
        title
        description
        prep_time
        cook_time
        total_time
        servings
        feature_image_url
        created_at
        user {
          id
          name
          avatar_image_url
        }
        user_likes_aggregate {
          aggregate {
            count
          }
        }
        user_bookmarks_aggregate {
          aggregate {
            count
          }
        }
        ratings_aggregate {
          aggregate {
            avg {
              value
            }
            count
          }
        }
      }
    }
  }
`, () => ({
  userId: userId.value
}), () => ({
  enabled: !!userId.value
}))

// Watch loading state
watch(queryLoading, (isLoading) => {
  loading.value = isLoading
})

// Format liked recipes data
const likedRecipes = computed(() => {
  if (!likedRecipesResult.value) return []
  
  return likedRecipesResult.value.user_likes.map(like => ({
    ...like.recipe,
    user: like.recipe.user,
    is_liked: true, // Since these are liked recipes
    is_bookmarked: false, // You might want to fetch this separately
    likes_count: like.recipe.user_likes_aggregate.aggregate.count,
    bookmarks_count: like.recipe.user_bookmarks_aggregate.aggregate.count
  }))
})

// Handle recipe updates (like/unlike)
const handleRecipeUpdate = (updatedRecipe) => {
  // In a real implementation, you would update the Apollo cache here
  console.log('Recipe updated:', updatedRecipe)
  
  // If the recipe was unliked, you might want to refetch the data
  if (updatedRecipe.is_liked === false) {
    // You could trigger a refetch here
    // likedRecipesResult.refetch()
  }
}
</script>

<style scoped>
/* Add any custom styles here */
</style>