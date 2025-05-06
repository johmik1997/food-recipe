<template>
  <div class="min-h-screen bg-gray-50">
    <Navbar />
    
    <div class="py-8">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <!-- Header -->
        <div class="flex justify-between items-center mb-8">
          <h1 class="text-3xl font-bold text-gray-900">My Recipes</h1>
          <NuxtLink 
            to="/dashboard/recipes/create" 
            class="inline-flex items-center px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors"
          >
            <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
            </svg>
            Add New
          </NuxtLink>
        </div>

        <!-- Loading State -->
        <div v-if="pending" class="flex justify-center py-12">
          <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-green-500"></div>
        </div>

        <!-- Error State -->
        <div v-else-if="error" class="text-center py-12 bg-white rounded-xl shadow-sm">
          <div class="max-w-md mx-auto">
            <svg class="mx-auto h-12 w-12 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
            <h3 class="mt-4 text-lg font-medium text-gray-900">Failed to load recipes</h3>
            <p class="mt-2 text-sm text-gray-600">{{ error.message }}</p>
            <button 
              @click="refetch()"
              class="mt-6 inline-flex items-center px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors"
            >
              Try Again
            </button>
          </div>
        </div>

        <!-- Empty State -->
        <div v-else-if="recipes.length === 0" class="text-center py-16 bg-white rounded-xl shadow-sm">
          <div class="max-w-md mx-auto">
            <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
            </svg>
            <h3 class="mt-2 text-lg font-medium text-gray-900">No recipes yet</h3>
            <p class="mt-1 text-sm text-gray-500">Get started by creating your first recipe.</p>
            <NuxtLink 
              to="/dashboard/recipes/create" 
              class="mt-6 inline-flex items-center px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors"
            >
              Add Recipe
            </NuxtLink>
          </div>
        </div>

        <!-- Recipe List -->
        <div v-else class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          <div 
            v-for="recipe in recipes" 
            :key="recipe.id" 
            class="bg-white rounded-xl shadow-sm overflow-hidden hover:shadow-md transition-shadow"
          >
            <NuxtLink :to="`/recipes/${recipe.id}`">
              <img 
                :src="recipe.featured_image_url || '/placeholder-recipe.jpg'" 
                :alt="recipe.title"
                class="w-full h-48 object-cover"
              >
            </NuxtLink>
            <div class="p-6">
              <div class="flex justify-between items-start">
                <div>
                  <NuxtLink 
                    :to="`/recipes/${recipe.id}`" 
                    class="text-xl font-medium text-gray-900 hover:text-green-600"
                  >
                    {{ recipe.title }}
                  </NuxtLink>
                  <p class="mt-2 text-gray-600 line-clamp-2">{{ recipe.description }}</p>
                </div>
                <div class="flex space-x-2">
                  <NuxtLink 
                    :to="`/dashboard/recipes/edit/${recipe.id}`" 
                    class="p-2 text-gray-500 hover:text-green-600 rounded-full hover:bg-gray-100"
                    title="Edit"
                  >
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
                    </svg>
                  </NuxtLink>
                  <button 
                    @click="confirmDelete(recipe.id)"
                    class="p-2 text-gray-500 hover:text-red-600 rounded-full hover:bg-gray-100"
                    title="Delete"
                  >
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
                    </svg>
                  </button>
                </div>
              </div>
              <div class="mt-4 flex items-center text-sm text-gray-500">
                <span class="flex items-center mr-4">
                  <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                  </svg>
                  {{ recipe.prep_time + recipe.cook_time }} mins
                </span>
                <span class="flex items-center">
                  <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"></path>
                  </svg>
                  {{ recipe.servings }} servings
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <Footer />
  </div>
</template>

<script setup>
import gql from 'graphql-tag'
import { useQuery, useMutation } from '@vue/apollo-composable'
import { ref, computed, onMounted } from 'vue'

const userId = ref(null)

onMounted(() => {
  try {
    const userStr = localStorage.getItem("user")
    if (userStr) {
      const user = JSON.parse(userStr)
      userId.value = user.id
    }
  } catch (e) {
    console.error("Failed to read user data", e)
  }
})

const GET_MY_RECIPES = gql`
  query GetMyRecipes($userId: uuid!) {
    recipes(where: {user: {id: {_eq: $userId}}}) {
      id
      title
      featured_image_url
      description
      prep_time
      cook_time
      total_time
      servings
      user {
        name
      }
    }
  }
`

const { result, loading: pending, error, refetch } = useQuery(
  GET_MY_RECIPES,
  () => ({
    userId: userId.value
  }),
  {
    enabled: computed(() => !!userId.value)
  }
)

const DELETE_RECIPE = gql`
  mutation DeleteRecipe($id: uuid!) {
    delete_recipes_by_pk(id: $id) {
      id
    }
  }
`

const { mutate: deleteRecipe } = useMutation(DELETE_RECIPE)

const recipes = computed(() => {
  return result.value?.recipes ? [...result.value.recipes] : []
})

const confirmDelete = async (recipeId) => {
  if (confirm('Are you sure you want to delete this recipe?')) {
    try {
      await deleteRecipe({ id: recipeId })
      await refetch()
    } catch (err) {
      console.error('Error deleting recipe:', err)
      alert('Failed to delete recipe. Please try again.')
    }
  }
}
</script>