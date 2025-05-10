<template>
    <div class="min-h-screen bg-gray-50">
      <Navbar />
      
      <div class="py-12">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <!-- Category Header -->
          <div class="text-center mb-12">
            <h1 class="text-4xl font-bold text-gray-900 mb-2">
              <span class="bg-gradient-to-r from-green-600 to-primary-600 bg-clip-text text-transparent">
                {{ categoryName }} Recipes
              </span>
            </h1>
            <p class="text-lg text-gray-600">
              Browse all our delicious {{ categoryName.toLowerCase() }} recipes
            </p>
          </div>
  
          <!-- Recipe Grid -->
          <div v-if="recipes.length > 0" class="grid gap-8 md:grid-cols-2 lg:grid-cols-3">
            <RecipeCard 
              v-for="recipe in recipes" 
              :key="recipe.id" 
              :recipe="recipe"
            />
          </div>
  
          <!-- Empty State -->
          <div v-else class="text-center py-16 bg-white rounded-xl shadow-sm">
            <div class="max-w-md mx-auto">
              <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
              </svg>
              <h3 class="mt-4 text-lg font-medium text-gray-900">No recipes found</h3>
              <p class="mt-2 text-sm text-gray-500">We don't have any {{ categoryName.toLowerCase() }} recipes yet.</p>
            </div>
          </div>
        </div>
      </div>
  
      <Footer />
    </div>
  </template>
  
  <script setup>
  import { useRoute } from 'vue-router'
  import gql from 'graphql-tag'
  import { useQuery } from '@vue/apollo-composable'
  
  const route = useRoute()
  const categoryName = computed(() => {
    // Convert URL slug back to display name
    return route.params.name
      .replace(/-/g, ' ')
      .replace(/\b\w/g, l => l.toUpperCase())
  })
  
  const GET_CATEGORY_RECIPES = gql`
    query GetCategoryRecipes($categoryName: String!) {
      recipes(where: {recipe_categories: {category: {name: {_eq: $categoryName}}}}) {
        id
        title
        description
        prep_time
        cook_time
        servings
        featured_image_url
        user {
          name
          avatar_image_url
          id
        }
      }
    }
  `
  
  const { result, loading, error } = useQuery(
    GET_CATEGORY_RECIPES,
    () => ({
      categoryName: categoryName.value
    })
  )
  
  const recipes = computed(() => result.value?.recipes || [])
  </script>