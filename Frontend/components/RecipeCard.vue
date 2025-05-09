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
      
   <!-- User info section -->
<div class="flex items-center mt-4 pt-3 border-t border-gray-100">
  <div class="flex-shrink-0">
    <NuxtLink 
  :to="`/users/${recipe.user?.id}`" 
  class="flex items-center"
>
  <NuxtImg 
    :src="recipe.user?.avatar_image_url || '/placeholder-user.jpg'"
    :alt="recipe.user?.name || 'Anonymous'"
    class="h-8 w-8 rounded-full object-cover"
  />
</NuxtLink>
  </div>
  <div class="ml-3">
    <NuxtLink 
      :to="`/users/${recipe.user?.id}`" 
      class="text-sm font-medium text-gray-700 hover:text-primary-600"
    >
      {{ recipe.user?.name || 'Anonymous' }}
    </NuxtLink>
    <p class="text-xs text-gray-500">
      {{ formatDate(recipe.created_at) }}
    </p>
  </div>
</div>
      
      <NuxtLink 
        :to="`/recipes/${recipe.id}`" 
        class="mt-3 inline-block text-blue-600 hover:underline text-sm"
      >
        View Recipe
      </NuxtLink>
    </div>
  </div>
</template>

<script setup>
import { format } from 'date-fns'

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
      description: '',
      created_at: new Date().toISOString(),
      user: {
        name: '',
        avatar_image_url: ''
      }
    })
  }
})

const formatDate = (dateString) => {
  try {
    return format(new Date(dateString), 'MMM d, yyyy')
  } catch {
    return ''
  }
}
</script>