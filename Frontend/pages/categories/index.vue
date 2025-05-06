<template>
    <div class="min-h-screen bg-gray-50">
      <Navbar />
      
      <div class="py-12">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <!-- Hero Header -->
          <div class="text-center mb-16">
            <h1 class="text-4xl md:text-5xl font-bold text-gray-900 mb-4">
              <span class="bg-gradient-to-r from-green-600 to-primary-600 bg-clip-text text-transparent">
                Recipe Categories
              </span>
            </h1>
            <p class="max-w-2xl mx-auto text-lg text-gray-600">
              Browse recipes by category. Find inspiration for any meal or dietary preference.
            </p>
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
              <h3 class="mt-4 text-lg font-medium text-gray-900">Failed to load categories</h3>
              <p class="mt-2 text-sm text-gray-600">{{ error.message }}</p>
              <button 
                @click="refetch()"
                class="mt-6 inline-flex items-center px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors"
              >
                Try Again
              </button>
            </div>
          </div>
  
          <!-- Category Grid -->
          <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
            <!-- Dynamic Category Cards -->
            <NuxtLink 
                v-for="category in categories"
                :key="category.name"
                :to="`/categories/${category.name.toLowerCase().replace(/ & /g, '-').replace(/ /g, '-')}`"
                class="group relative h-64 rounded-2xl overflow-hidden shadow-md hover:shadow-xl transition-all duration-300"
                >
              <!-- Category Image (you might want to add images to your categories in the database) -->
              <img 
                :src="getCategoryImage(category.name)"
                :alt="category.name"
                class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
              >
              <div class="absolute inset-0 bg-gradient-to-t from-black/60 to-transparent flex items-end p-6">
                <div>
                  <h2 class="text-2xl font-bold text-white">{{ category.name }}</h2>
                  <p class="text-gray-200 mt-1">{{ getRecipeCount(category) }} recipes</p>
                </div>
                <span class="ml-auto text-3xl">{{ getCategoryEmoji(category.name) }}</span>
              </div>
            </NuxtLink>
  
            <!-- View All -->
            <NuxtLink 
              to="/categories/all" 
              class="group relative h-64 rounded-2xl overflow-hidden shadow-md hover:shadow-xl transition-all duration-300 border-2 border-dashed border-gray-300 flex items-center justify-center"
            >
              <div class="text-center p-6">
                <div class="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-green-100 text-green-600 mb-3">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                  </svg>
                </div>
                <h2 class="text-2xl font-bold text-gray-800">View All</h2>
                <p class="text-gray-600 mt-1">Browse all recipes</p>
              </div>
            </NuxtLink>
          </div>
        </div>
      </div>
  
      <Footer />
    </div>
  </template>
  
  <script setup>
  import gql from 'graphql-tag'
  import { useQuery } from '@vue/apollo-composable'
  import { computed } from 'vue'
  
  const GET_CATEGORIES = gql`
    query GetCategories {
      categories {
        name
        recipe_categories {
          recipe {
            id
          }
        }
      }
    }
  `
  
  const { result, loading: pending, error, refetch } = useQuery(GET_CATEGORIES)
  
  const categories = computed(() => {
    return result.value?.categories || []
  })
  
  // Helper function to get emoji for each category
  const getCategoryEmoji = (categoryName) => {
    const emojiMap = {
      'Breakfast': '🥞',
      'Lunch': '🥗',
      'Dinner': '🍝',
      'Desserts': '🍰',
      'Vegetarian': '🥕',
      'Vegan': '🌱',
      'Gluten-Free': '🌾',
      'Salads': '🥗',
      'Soups & Stews': '🍲',
      'Appetizers': '🍤',
      'Main Dishes': '🍛',
      'Side Dishes': '🥔',
      'Healthy': '🥑',
      'Kid-Friendly': '🧒',
      'Snacks': '🍪',
      'vegetable': '🥦'
    }
    return emojiMap[categoryName] || '🍽️'
  }
  
  // Helper function to get recipe count
  const getRecipeCount = (category) => {
    return category.recipe_categories?.length || 0
  }
  
  // Helper function to get category image (you can replace with actual images from your DB)
  const getCategoryImage = (categoryName) => {
    const imageMap = {
      'Breakfast': 'https://images.unsplash.com/photo-1550583724-b2692b85b150?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=687&q=80',
      'Lunch': 'https://images.unsplash.com/photo-1546069901-ba9599a7e63c?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=880&q=80',
      'Dinner': 'https://images.unsplash.com/photo-1547592180-85f173990554?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1470&q=80',
      'Desserts': 'https://images.unsplash.com/photo-1571115177098-24ec42ed204d?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=687&q=80',
      'Vegan': 'https://images.unsplash.com/photo-1543351611-58f69d7c1781?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1374&q=80'
    }
    return imageMap[categoryName] || 'https://images.unsplash.com/photo-1504674900247-0877df9cc836?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1470&q=80'
  }
  </script>