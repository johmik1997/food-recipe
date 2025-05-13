<template>
  <div class="recipe-card bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow">
    <div class="relative">
      <NuxtImg 
        :src="recipe.featured_image_url || '/placeholder-food.jpg'" 
        :alt="recipe.title"
        class="w-full h-48 object-cover"
        width="400"
        height="192"
      />
      <!-- Action buttons overlay -->
      <div class="absolute top-2 right-2 flex gap-2">
        <button 
          @click.stop="toggleBookmark"
          class="p-2 bg-white/80 rounded-full hover:bg-white transition-colors"
          :aria-label="recipe.is_bookmarked ? 'Remove bookmark' : 'Bookmark recipe'"
        >
          <Icon 
            :name="recipe.is_bookmarked ? 'material-symbols:bookmark-added' : 'material-symbols:bookmark-outline'" 
            class="text-lg"
            :class="recipe.is_bookmarked ? 'text-green-600' : 'text-gray-600'"
          />
        </button>
        <button 
          @click.stop="toggleLike"
          class="p-2 bg-white/80 rounded-full hover:bg-white transition-colors relative"
          :aria-label="recipe.is_liked ? 'Unlike recipe' : 'Like recipe'"
        >
          <Icon 
            :name="recipe.is_liked ? 'material-symbols:favorite' : 'material-symbols:favorite-outline'" 
            class="text-lg"
            :class="recipe.is_liked ? 'text-red-500' : 'text-gray-600'"
          />
          <span 
            v-if="recipe.likes_count > 0"
            class="absolute -bottom-1 -right-1 text-xs bg-white rounded-full px-1 min-w-[16px] text-center"
          >
            {{ recipe.likes_count }}
          </span>
        </button>
      </div>
    </div>

    <div class="p-4">
      <h3 class="text-lg font-semibold mb-2">{{ recipe.title }}</h3>
      
      <!-- Rating display -->
      <div v-if="recipe.ratings_aggregate?.aggregate?.avg?.value" class="flex items-center mb-2">
        <div class="flex items-center mr-2">
          <Icon 
            v-for="i in 5" 
            :key="i"
            :name="i <= Math.round(recipe.ratings_aggregate.aggregate.avg.value) ? 'material-symbols:star' : 'material-symbols:star-outline'"
            class="h-4 w-4"
            :class="i <= Math.round(recipe.ratings_aggregate.aggregate.avg.value) ? 'text-yellow-400' : 'text-gray-300'"
          />
        </div>
        <span class="text-sm text-gray-600">
          {{ recipe.ratings_aggregate.aggregate.avg.value.toFixed(1) }} ({{ recipe.ratings_aggregate.aggregate.count }})
        </span>
      </div>
      
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
            :to="`/user/${recipe.user?.id}`" 
            class="flex items-center"
          >
            <NuxtImg 
              :src="recipe.user?.avatar_image_url || '/placeholder-user.jpg'"
              :alt="recipe.user?.name || 'Anonymous'"
              class="h-8 w-8 rounded-full object-cover"
              width="32"
              height="32"
            />
          </NuxtLink>
        </div>
        <div class="ml-3">
          <NuxtLink 
            :to="`/user/${recipe.user?.id}`" 
            class="text-sm font-medium text-gray-700 hover:text-primary-600"
          >
            {{ recipe.user?.name || 'Anonymous' }}
          </NuxtLink>
          <p class="text-xs text-gray-500">
            {{ formatDate(recipe.created_at) }}
          </p>
        </div>
      </div>
      
      <!-- Action buttons -->
       <div class="flex justify-between mt-4">
      <div class="flex justify-between mt-4">
        <NuxtLink 
          :to="`/recipes/${recipe.id}`" 
          class="inline-flex items-center px-3 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-green-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 transition-colors"
        >
          <Icon name="material-symbols:visibility" class="mr-1" />
          View Recipe
        </NuxtLink>
      </div>
      <div class="flex justify-between mt-4">
        <NuxtLink 
       :to="`/payment/${recipe.id}`"
          class="inline-flex items-center px-3 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-red-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 transition-colors"
        >
          <Icon name="material-symbols:purchase" class="mr-1" />
          Buy Recipe
        </NuxtLink>
      </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useMutation } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { format } from 'date-fns'

const props = defineProps({
  recipe: {
    type: Object,
    required: true
  },
  showActions: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update', 'edit', 'delete'])

const userId = ref(null)
const isProcessing = ref(false)
const errorMessage = ref(null)

onMounted(() => {
  const userStr = localStorage.getItem("user")
  if (userStr) {
    try {
      const user = JSON.parse(userStr)
      userId.value = user.id
    } catch (error) {
      console.error("Failed to parse user data", error)
    }
  }
})

// GraphQL mutations
const { mutate: like } = useMutation(gql`
  mutation LikeRecipe($recipeId: uuid!, $userId: uuid!) {
    insert_user_likes_one(
      object: { recipe_id: $recipeId, user_id: $userId }
      on_conflict: { 
        constraint: user_likes_user_id_recipe_id_key
        update_columns: [] 
      }
    ) {
      id
    }
  }
`)

const { mutate: unlike } = useMutation(gql`
  mutation UnlikeRecipe($recipeId: uuid!, $userId: uuid!) {
    delete_user_likes(
      where: { 
        recipe_id: { _eq: $recipeId }, 
        user_id: { _eq: $userId } 
      }
    ) {
      affected_rows
    }
  }
`)

const { mutate: bookmark } = useMutation(gql`
  mutation BookmarkRecipe($recipeId: uuid!, $userId: uuid!) {
    insert_user_bookmarks_one(
      object: { recipe_id: $recipeId, user_id: $userId }
      on_conflict: { 
        constraint: user_bookmarks_user_id_recipe_id_key
        update_columns: [] 
      }
    ) {
      id
    }
  }
`)

const { mutate: unbookmark } = useMutation(gql`
  mutation UnbookmarkRecipe($recipeId: uuid!, $userId: uuid!) {
    delete_user_bookmarks(
      where: { 
        recipe_id: { _eq: $recipeId }, 
        user_id: { _eq: $userId } 
      }
    ) {
      affected_rows
    }
  }
`)

const toggleLike = async () => {
  if (isProcessing.value) return
  if (!userId.value) return navigateTo('/login')
  
  isProcessing.value = true
  errorMessage.value = null

  try {
    const variables = { recipeId: props.recipe.id, userId: userId.value }
    const newLikeStatus = !props.recipe.is_liked
    
    // Optimistic UI update
    emit('update', {
      id: props.recipe.id,
      is_liked: newLikeStatus,
      likes_count: newLikeStatus 
        ? props.recipe.likes_count + 1 
        : Math.max(0, props.recipe.likes_count - 1)
    })

    // Perform mutation
    if (newLikeStatus) {
      await like(variables)
    } else {
      await unlike(variables)
    }
  } catch (error) {
    console.error('Error toggling like:', error)
    errorMessage.value = "Failed to update like. Please try again."
    // Revert optimistic update
    emit('update', {
      id: props.recipe.id,
      is_liked: props.recipe.is_liked,
      likes_count: props.recipe.likes_count
    })
  } finally {
    isProcessing.value = false
  }
}

const toggleBookmark = async () => {
  if (isProcessing.value) return
  if (!userId.value) return navigateTo('/login')
  
  isProcessing.value = true
  errorMessage.value = null

  try {
    const variables = { recipeId: props.recipe.id, userId: userId.value }
    const newBookmarkStatus = !props.recipe.is_bookmarked
    
    // Optimistic UI update
    emit('update', {
      id: props.recipe.id,
      is_bookmarked: newBookmarkStatus
    })

    // Perform mutation
    if (newBookmarkStatus) {
      await bookmark(variables)
    } else {
      await unbookmark(variables)
    }
  } catch (error) {
    console.error('Error toggling bookmark:', error)
    errorMessage.value = "Failed to update bookmark. Please try again."
    // Revert optimistic update
    emit('update', {
      id: props.recipe.id,
      is_bookmarked: props.recipe.is_bookmarked
    })
  } finally {
    isProcessing.value = false
  }
}

const formatDate = (dateString) => {
  try {
    return format(new Date(dateString), 'MMM d, yyyy')
  } catch {
    return ''
  }
}
</script>

<style scoped>
.recipe-card {
  position: relative;
  transition: all 0.3s ease;
}

.recipe-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.1);
}

/* Smooth transitions for interactive elements */
button {
  transition: all 0.2s ease;
}

/* Like counter badge */
.min-w-\[16px\] {
  min-width: 16px;
}

/* Rating stars */
.rating-star {
  transition: color 0.2s ease;
}

/* Action buttons container */
.action-buttons {
  display: flex;
  gap: 0.5rem;
  margin-top: 1rem;
}
</style>