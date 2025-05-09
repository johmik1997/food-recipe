<template>
  <div class="min-h-screen bg-gray-50">
    <Navbar />
    
    <div class="py-12">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <!-- Hero Header -->
        <div class="text-center mb-16">
          <h1 class="text-4xl md:text-5xl font-bold text-gray-900 mb-4">
            <span class="bg-gradient-to-r from-green-600 to-primary-600 bg-clip-text text-transparent">
              Delicious Recipes
            </span>
          </h1>
          <p class="max-w-2xl mx-auto text-lg text-gray-600">
            Discover mouth-watering recipes from our community. Find inspiration for your next meal!
          </p>
        </div>

        <!-- Search and Filter Bar -->
        <div class="mb-12 bg-white p-6 rounded-xl shadow-sm border border-gray-100">
          <div class="flex flex-col md:flex-row gap-4 items-center">
            <!-- Search Input -->
            <div class="relative flex-1 w-full">
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <svg class="h-5 w-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
              </div>
              <input
                type="text"
                v-model="searchQuery"
                placeholder="Search recipes..."
                class="block w-full pl-10 pr-3 py-3 border border-gray-300 rounded-lg bg-gray-50 focus:ring-2 focus:ring-primary-500 focus:border-transparent"
              >
            </div>
            
            <!-- Filters -->
            <div class="flex gap-3 w-full md:w-auto">
              <!-- Preparation Time Filter -->
              <div class="relative w-full">
                <select 
                  v-model="prepTimeFilter"
                  class="block w-full px-3 py-3 border border-gray-300 rounded-lg bg-gray-50 focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                >
                  <option value="">All Times</option>
                  <option value="10">Quick (≤ 10 mins)</option>
                  <option value="15">Quick (≤ 15 mins)</option>
                  <option value="30">Fast (≤ 30 mins)</option>
                  <option value="60">Medium (≤ 1 hour)</option>
                  <option value="120">Long (≤ 2 hours)</option>
                  <option value="121">Very Long (> 2 hours)</option>
                </select>
              </div>

              <!-- Ingredients Filter -->
              <div class="relative w-full">
                <input
                  v-model="ingredientFilter"
                  type="text"
                  placeholder="Ingredients..."
                  class="block w-full px-3 py-3 border border-gray-300 rounded-lg bg-gray-50 focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                />
              </div>

              <!-- Clear Filters Button -->
              <button
                v-if="hasFilters"
                @click="clearFilters"
                class="px-4 py-3 text-sm text-gray-600 hover:text-gray-800 transition-colors whitespace-nowrap"
              >
                Clear Filters
              </button>
            </div>
          </div>
        </div>
        
        <!-- Action Bar -->
        <div class="flex flex-col sm:flex-row justify-between items-center mb-8 gap-4">
          <h2 class="text-2xl font-bold text-gray-800">
            {{ filteredRecipes.length }} {{ filteredRecipes.length === 1 ? 'Recipe' : 'Recipes' }} Available
          </h2>
          <NuxtLink 
            to="/recipes/create" 
            class="inline-flex items-center px-6 py-3 border border-transparent text-base font-medium rounded-xl shadow-sm text-white bg-gradient-to-r from-green-600 to-primary-600 hover:from-green-700 hover:to-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 transition-all duration-200 transform hover:-translate-y-0.5"
          >
            <svg class="-ml-1 mr-2 h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
            </svg>
            Add Recipe
          </NuxtLink>
        </div>

        <!-- Loading State -->
        <div v-if="pending" class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          <div v-for="n in 6" :key="n" class="bg-white rounded-xl overflow-hidden shadow-sm border border-gray-100">
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
            <h3 class="mt-4 text-lg font-medium text-gray-900">Failed to load recipes</h3>
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
        <div v-else-if="filteredRecipes.length === 0" class="text-center py-16 bg-white rounded-xl shadow-sm border border-gray-100">
          <div class="max-w-md mx-auto">
            <svg class="mx-auto h-16 w-16 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
            </svg>
            <h3 class="mt-4 text-lg font-medium text-gray-900">
              {{ hasFilters ? 'No recipes match your filters' : 'No recipes found' }}
            </h3>
            <p class="mt-2 text-sm text-gray-600">
              {{ hasFilters ? 'Try adjusting your search criteria' : 'Get started by creating your first recipe.' }}
            </p>
            <button
              v-if="hasFilters"
              @click="clearFilters"
              class="mt-6 inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
            >
              Clear Filters
            </button>
            <NuxtLink 
              v-else
              to="/recipes/create" 
              class="mt-6 inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
            >
              <svg class="-ml-1 mr-2 h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
              </svg>
              Add Recipe
            </NuxtLink>
          </div>
        </div>

        <!-- Recipe Grid -->
        <div v-else class="grid gap-8 md:grid-cols-2 lg:grid-cols-3">
          <RecipeCard 
            v-for="recipe in filteredRecipes" 
            :key="recipe.id"
            :recipe="recipe"
            :show-actions="true"
          />
        </div>

        <!-- Pagination -->
        <div v-if="filteredRecipes.length > 0" class="mt-12 flex justify-center">
          <nav class="inline-flex rounded-md shadow-sm -space-x-px">
            <button class="px-4 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50">
              Previous
            </button>
            <button class="px-4 py-2 border-t border-b border-gray-300 bg-white text-sm font-medium text-primary-600 hover:bg-gray-50">
              1
            </button>
            <button class="px-4 py-2 border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50">
              2
            </button>
            <button class="px-4 py-2 border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50">
              3
            </button>
            <span class="px-4 py-2 border border-gray-300 bg-white text-sm font-medium text-gray-700">
              ...
            </span>
            <button class="px-4 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50">
              Next
            </button>
          </nav>
        </div>
      </div>
    </div>
    <Footer />
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'

// Search and filter states
const searchQuery = ref('')
const prepTimeFilter = ref('')
const ingredientFilter = ref('')

// GraphQL query
const SEARCHABLE_RECIPES = gql`
  query SearchableRecipes {
    recipes {
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
        name
        avatar_image_url
      }
      recipe_ingredients {
          name
      }
    }
  }
`


const { result, loading: pending, error, refetch } = useQuery(SEARCHABLE_RECIPES)

// Raw recipes data
const recipes = computed(() => result.value?.recipes || [])

// Filtered recipes
const filteredRecipes = computed(() => {
  return recipes.value.filter(recipe => {
    // Search filter
    const matchesSearch = searchQuery.value === '' || 
                         recipe.title.toLowerCase().includes(searchQuery.value.toLowerCase()) || 
                         recipe.description.toLowerCase().includes(searchQuery.value.toLowerCase())
    
    // Preparation time filter
    const prepTime = parseInt(recipe.prep_time) || 0
    let matchesPrepTime = true
    if (prepTimeFilter.value) {
      const filterTime = parseInt(prepTimeFilter.value)
      matchesPrepTime = filterTime === 121 ? prepTime > 120 : prepTime <= filterTime
    }
    
    // Ingredients filter
    let matchesIngredients = true
    if (ingredientFilter.value) {
      const searchTerms = ingredientFilter.value.toLowerCase().split(',').map(term => term.trim())
      const recipeIngredients = recipe.recipe_ingredients?.map(ri => ri?.name?.toLowerCase()) || []
      
      // Filter out any undefined/null values
      const validIngredients = recipeIngredients.filter(Boolean)
      
      matchesIngredients = searchTerms.every(term => 
        validIngredients.some(ingredient => ingredient.includes(term)))
    }
    
    return matchesSearch && matchesPrepTime && matchesIngredients
  })
})

// Check if any filters are active
const hasFilters = computed(() => {
  return prepTimeFilter.value || ingredientFilter.value || searchQuery.value
})

// Clear all filters
function clearFilters() {
  searchQuery.value = ''
  prepTimeFilter.value = ''
  ingredientFilter.value = ''
}
</script>