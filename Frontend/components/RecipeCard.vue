<!-- ~/components/RecipeCard.vue -->
<template>
    <div class="recipe-card bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow">
      <NuxtImg 
        :src="recipe.featured_image_url || '/placeholder-food.jpg'" 
        :alt="recipe.title"
        class="w-full h-48 object-cover"
      />
      <div class="p-4">
        <h3 class="text-lg font-semibold mb-2">{{ recipe.title }}</h3>
        <div class="flex items-center text-sm text-gray-600 mb-2">
          <Icon name="material-symbols:time" class="mr-1" />
          <span>Total Time: {{ recipe.total_time }} mins</span>
        </div>
        <div class="flex gap-4 text-sm mb-3">
          <span class="flex items-center">
            <Icon name="material-symbols:prep" class="mr-1" />
            Prep: {{ recipe.prep_time }}m
          </span>
          <span class="flex items-center">
            <Icon name="material-symbols:cook" class="mr-1" />
            Cook: {{ recipe.cook_time }}m
          </span>
        </div>
        <p class="text-gray-600 text-sm mb-3 line-clamp-2">{{ recipe.description }}</p>
        <NuxtLink 
          :to="`/recipes/${recipe.id}`" 
          class="mt-2 inline-block text-blue-600 hover:underline"
        >
          View Recipe
        </NuxtLink>
      </div>
    </div>
  </template>
  
  <script setup>
  defineProps({
    recipe: {
      type: Object,
      required: true,
      default: () => ({
        id: '',
        title: '',
        featured_image_url: '',
        total_time: 0,
        prep_time: 0,
        cook_time: 0,
        description: ''
      })
    }
  })
  </script>