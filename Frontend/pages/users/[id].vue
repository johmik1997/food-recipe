<template>
  <div class="min-h-screen bg-gray-50">
    <Navbar />
    
    <div class="py-12">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <!-- User Header -->
        <div class="text-center mb-12">
          <div class="mx-auto h-24 w-24 rounded-full overflow-hidden bg-gray-200 mb-4">
            <img 
              :src="userData?.avatar_image_url || '/placeholder-user.jpg'" 
              :alt="userData?.name"
              class="h-full w-full object-cover"
            />
          </div>
          <h1 class="text-3xl font-bold text-gray-900">
            Recipes by {{ userData?.name || 'this user' }}
          </h1>
          <p class="mt-2 text-gray-600">
            {{ filteredRecipes.length }} recipes shared
          </p>
        </div>

        <!-- Recipe Grid -->
        <div v-if="pending" class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          <div v-for="n in 6" :key="n" class="bg-white rounded-xl shadow-sm animate-pulse h-64"></div>
        </div>

        <div v-else-if="error" class="text-center py-16">
          <p class="text-red-500">Error loading recipes: {{ error.message }}</p>
          <button @click="refetch()" class="mt-4 px-4 py-2 bg-primary-600 text-white rounded">
            Try Again
          </button>
        </div>

        <div v-else-if="filteredRecipes.length === 0" class="text-center py-16">
          <p class="text-gray-500">No recipes found for this user</p>
        </div>

        <div v-else class="grid gap-8 md:grid-cols-2 lg:grid-cols-3">
          <RecipeCard 
            v-for="recipe in filteredRecipes" 
            :key="recipe.id"
            :recipe="recipe"
          />
        </div>
      </div>
    </div>
    <Footer />
  </div>
</template>

<script setup>
import { useRoute } from 'vue-router'

const route = useRoute()
const userId = route.params.id // Changed to use 'id' instead of 'userId'

// Query to get user details
const GET_USER = gql`
  query GetUser($id: uuid!) {
    users_by_pk(id: $id) {
      id
      name
      avatar_image_url
    }
  }
`

// Query to get recipes by user
const GET_RECIPES_BY_USER = gql`
  query GetRecipesByUser($userId: uuid!) {
    recipes(where: {user_id: {_eq: $userId}}, order_by: {created_at: desc}) {
      id
      title
      description
      prep_time
      cook_time
      featured_image_url
      user {
        name
        avatar_image_url
      }
      ratings_aggregate {
        aggregate {
          avg {
            value
          }
        }
      }
    }
  }
`

// Fetch user data
const { result: userResult } = useQuery(GET_USER, {
  id: userId
})

// Fetch recipes
const { result, pending, error, refetch } = useQuery(GET_RECIPES_BY_USER, {
  userId: userId
})

const userData = computed(() => userResult.value?.users_by_pk)
const filteredRecipes = computed(() => result.value?.recipes || [])
</script>