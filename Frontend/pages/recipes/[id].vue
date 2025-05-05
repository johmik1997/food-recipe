<template>
    <div class="max-w-4xl mx-auto px-4 py-8">
      <!-- Loading state -->
      <div v-if="loading" class="text-center py-12">
        <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500 mx-auto"></div>
        <p class="mt-4 text-gray-600">Loading recipe...</p>
      </div>
  
      <!-- Error state -->
      <div v-if="error" class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
        <p>Error loading recipe: {{ error.message }}</p>
        <button 
          @click="fetchRecipe" 
          class="mt-2 px-3 py-1 bg-red-100 hover:bg-red-200 rounded text-red-700"
        >
          Retry
        </button>
      </div>
  
      <!-- Recipe content -->
      <div v-if="recipe" class="space-y-8">
        <!-- Recipe header -->
        <div>
          <h1 class="text-3xl font-bold text-gray-900">{{ recipe.title }}</h1>
          
          <div class="mt-2 flex items-center space-x-4 text-gray-600">
            <span v-if="recipe.prep_time" class="flex items-center">
              <ClockIcon class="h-5 w-5 mr-1" />
              Prep: {{ recipe.prep_time }} min
            </span>
            <span v-if="recipe.cook_time" class="flex items-center">
              <FireIcon class="h-5 w-5 mr-1" />
              Cook: {{ recipe.cook_time }} min
            </span>
            <span v-if="recipe.servings" class="flex items-center">
              <UsersIcon class="h-5 w-5 mr-1" />
              Serves: {{ recipe.servings }}
            </span>
          </div>
        </div>
  
        <!-- Featured image -->
        <div v-if="recipe.featured_image_url" class="rounded-lg overflow-hidden shadow-lg">
          <img 
            :src="recipe.featured_image_url" 
            :alt="recipe.title" 
            class="w-full h-96 object-cover"
          />
        </div>
  
        <!-- Description -->
        <div v-if="recipe.description" class="prose max-w-none">
          <p class="text-lg text-gray-700">{{ recipe.description }}</p>
        </div>
  
        <!-- Ingredients -->
        <div class="bg-gray-50 p-6 rounded-lg">
          <h2 class="text-2xl font-bold mb-4 text-gray-900">Ingredients</h2>
          <ul class="space-y-2">
            <li v-for="(ingredient, index) in recipe.ingredients" :key="index" class="flex items-start">
              <span class="inline-block w-2 h-2 bg-blue-500 rounded-full mt-2 mr-2"></span>
              <span>
                <span v-if="ingredient.quantity">{{ ingredient.quantity }}</span>
                <span v-if="ingredient.unit" class="ml-1">{{ ingredient.unit }}</span>
                <span class="ml-2 font-medium">{{ ingredient.name }}</span>
              </span>
            </li>
          </ul>
        </div>
  
        <!-- Steps -->
        <div class="space-y-6">
          <h2 class="text-2xl font-bold mb-4 text-gray-900">Instructions</h2>
          <div v-for="step in recipe.steps" :key="step.step_number" class="flex">
            <div class="flex-shrink-0 mr-4">
              <span class="flex items-center justify-center w-8 h-8 rounded-full bg-blue-500 text-white font-bold">
                {{ step.step_number }}
              </span>
            </div>
            <div class="prose max-w-none flex-1">
              <p>{{ step.instruction }}</p>
              <img 
                v-if="step.image_url" 
                :src="step.image_url" 
                :alt="`Step ${step.step_number}`"
                class="mt-3 rounded-lg shadow-md max-w-full"
              />
            </div>
          </div>
        </div>
  
        <!-- Additional images -->
        <div v-if="recipe.images && recipe.images.length > 0" class="grid grid-cols-2 md:grid-cols-3 gap-4">
          <div 
            v-for="(image, index) in recipe.images" 
            :key="index"
            class="overflow-hidden rounded-lg shadow-md"
          >
            <img 
              :src="image.image_url" 
              :alt="`Recipe image ${index + 1}`"
              class="w-full h-48 object-cover hover:scale-105 transition-transform duration-300"
            />
          </div>
        </div>
      </div>
    </div>
  </template>
  
  <script setup>
  import { useRoute } from 'vue-router'
  import { ClockIcon, FireIcon, UsersIcon } from '@heroicons/vue/24/outline'
  import gql from 'graphql-tag'
import { useQuery } from '@vue/apollo-composable'
  
  const route = useRoute()
  const recipeId = route.params.id
  
  // GraphQL query
  const GET_RECIPE_BY_ID = gql`
    query GetRecipeById($id: uuid!) {
      recipes_by_pk(id: $id) {
        id
        title
        description
        prep_time
        cook_time
        servings
        featured_image_url
        created_at
        user {
          id
          name
        }
        recipe_ingredients {
          name
          quantity
          unit
        }
       recipe_steps(order_by: {step_number: asc}) {
          step_number
          instruction
          image_url
        }
        recipe_images(where: {is_featured: {_eq: false}}) {
          image_url
        }
        recipe_categories {
      category {
        id
        name
      }
    }
      }
    }
  `
  
  const { result, loading, error, refetch } = useQuery(GET_RECIPE_BY_ID, {
    id: recipeId
  })
  
  const recipe = computed(() => result.value?.recipes_by_pk)
  
  // Function to retry if error occurs
  const fetchRecipe = () => {
    refetch()
  }
  
  // Set page title
  useHead({
    title: computed(() => recipe.value ? `${recipe.value.title} - FoodRecipe` : 'Loading Recipe...'),
    meta: [
      {
        name: 'description',
        content: computed(() => recipe.value?.description || 'View this delicious recipe')
      }
    ]
  })
  </script>
  
  <style scoped>
  .prose {
    color: #374151; /* gray-700 */
    line-height: 1.75;
  }
  
  .prose p {
    margin-bottom: 1rem;
  }
  </style>