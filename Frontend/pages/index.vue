<!-- ~/pages/recipes/index.vue -->
<template>
  <div class="py-8">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center mb-8">
        <h1 class="text-3xl font-bold text-gray-900">Recipes</h1>
        <NuxtLink 
          to="/recipes/create" 
          class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
        >
          Add Recipe
        </NuxtLink>
      </div>

      <div v-if="pending" class="flex justify-center py-12">
        <Spinner class="h-8 w-8 text-primary-500" />
      </div>

      <div v-else-if="error" class="text-center py-12 text-red-500">
        Error loading recipes: {{ error.message }}
      </div>

      <div v-else-if="recipes.length === 0" class="text-center py-12">
        <div class="max-w-md mx-auto">
          <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
          </svg>
          <h3 class="mt-2 text-lg font-medium text-gray-900">No recipes yet</h3>
          <p class="mt-1 text-sm text-gray-500">Get started by creating your first recipe.</p>
        </div>
      </div>

      <div v-else class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        <RecipeCard 
          v-for="recipe in recipes" 
          :key="recipe.id" 
          :recipe="recipe" 
          class="hover:shadow-lg transition-shadow duration-300"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'

const GET_RECIPES = gql`
  query GetRecipes {
    recipes(distinct_on: title) {
      id
      title
      description
      prep_time
      cook_time
      total_time
      servings
      featured_image_url
      created_at
    }
  }
`

const { result, loading: pending, error } = useQuery(GET_RECIPES)

const recipes = computed(() => result.value?.recipes || [])
</script>