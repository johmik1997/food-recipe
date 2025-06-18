<template>
  <div class="bg-white rounded-lg shadow-sm overflow-hidden hover:shadow-md transition-shadow duration-300 h-full flex flex-col">
    <!-- Recipe Image -->
    <NuxtLink :to="`/recipes/${recipe.id}`" class="block relative aspect-w-4 aspect-h-3">
      <img 
        :src="recipe.featured_image || '/placeholder-recipe.jpg'" 
        :alt="recipe.title"
        class="w-full h-48 object-cover hover:scale-105 transition-transform duration-300"
      >
      <!-- Overlay with quick actions -->
      <div class="absolute inset-0 bg-black bg-opacity-0 hover:bg-opacity-10 transition-all duration-300 flex items-end p-3">
        <div class="flex space-x-2">
          <span class="bg-green-600 text-white text-xs px-2 py-1 rounded-full">
            {{ recipe.catagory?.name || 'Uncategorized' }}
          </span>
          <span v-if="recipe.price > 0" class="bg-white text-green-600 text-xs px-2 py-1 rounded-full font-bold">
            ${{ recipe.price }}
          </span>
        </div>
      </div>
    </NuxtLink>

    <!-- Recipe Content -->
    <div class="p-4 flex-grow flex flex-col">
      <!-- Title and Author -->
      <div class="mb-3">
        <NuxtLink :to="`/recipes/${recipe.id}`">
          <h3 class="text-lg font-bold text-gray-900 line-clamp-2 hover:text-green-600 transition-colors">
            {{ recipe.title }}
          </h3>
        </NuxtLink>
        <div v-if="showAuthor" class="flex items-center mt-2">
          <img 
            :src="recipe.user?.profile || '/placeholder-avatar.jpg'" 
            :alt="recipe.user?.username"
            class="w-6 h-6 rounded-full object-cover mr-2"
          >
          <span class="text-sm text-gray-600">{{ recipe.user?.username }}</span>
        </div>
      </div>

      <!-- Stats -->
      <div class="flex items-center justify-between mt-auto">
        <div class="flex items-center space-x-4 text-sm text-gray-500">
          <div class="flex items-center">
            <Icon name="heroicons:clock" class="w-4 h-4 mr-1" />
            <span>{{ recipe.prep_time}} min</span>
          </div>
          <div class="flex items-center">
            <Icon name="heroicons:star" class="w-4 h-4 mr-1 text-yellow-400" />
            <span>{{ recipe.average_rating || '0' }}</span>
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="flex items-center space-x-2">
          <button 
            @click.stop="toggleLike"
            class="p-2 rounded-full hover:bg-green-50 transition-colors"
            :class="{ 'text-red-500': isLiked, 'text-gray-400': !isLiked }"
            aria-label="Like recipe"
          >
            <Icon :name="isLiked ? 'heroicons:heart-20-solid' : 'heroicons:heart'" class="w-5 h-5" />
          </button>
          <button 
            @click.stop="toggleBookmark"
            class="p-2 rounded-full hover:bg-green-50 transition-colors"
            :class="{ 'text-green-600': isBookmarked, 'text-gray-400': !isBookmarked }"
            aria-label="Bookmark recipe"
          >
            <Icon :name="isBookmarked ? 'heroicons:bookmark-solid' : 'heroicons:bookmark'" class="w-5 h-5" />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useToast } from 'vue-toast-notification'

const props = defineProps({
  recipe: {
    type: Object,
    required: true
  },
  showAuthor: {
    type: Boolean,
    default: true
  }
})

const toast = useToast()
const userStore = authStore()
const recipeStore = useRecipeStore()
const bookmarkStore = useBookmarkStore()
const likeStore = useLikeStore()

const isAuthenticated = computed(() => userStore.isAuthenticated)
const isBookmarked = ref(false)
const isLiked = ref(false)

// Check initial bookmark status
onMounted(async () => {
  if (isAuthenticated.value) {
    await checkBookmarkStatus()
    await checkLikeStatus()
  }
})

const checkBookmarkStatus = async () => {
  try {
    const payload = {
      recipe_id: props.recipe.id,
      user_id: userStore.userId
    }
    await bookmarkStore.checkIfBookmarked(payload)
    isBookmarked.value = bookmarkStore.$state.isBookmarked
  } catch (error) {
    console.error('Error checking bookmark status:', error)
  }
}

const checkLikeStatus = async () => {
  try {
    const payload = {
      recipe_id: props.recipe.id,
      user_id: userStore.userId
    }
    await likeStore.checkIfLiked(payload)
    isLiked.value = likeStore.$state.isLiked
  } catch (error) {
    console.error('Error checking like status:', error)
  }
}

const toggleBookmark = async () => {
  if (!isAuthenticated.value) {
    toast.info('Please login to save recipes')
    return
  }

  try {
    if (isBookmarked.value) {
      await bookmarkStore.removeBookmark(bookmarkStore.bookmarkedId)
      toast.success('Recipe removed from favorites')
    } else {
      await bookmarkStore.createBookmark({
        recipe_id: props.recipe.id,
        user_id: userStore.userId
      })
      toast.success('Recipe saved to favorites')
    }
    isBookmarked.value = !isBookmarked.value
  } catch (error) {
    toast.error('Failed to update bookmark')
    console.error('Bookmark error:', error)
  }
}

const toggleLike = async () => {
  if (!isAuthenticated.value) {
    toast.info('Please login to like recipes')
    return
  }

  try {
    if (isLiked.value) {
      await likeStore.removeLike(likeStore.likeId)
      toast.success('Recipe unliked')
    } else {
      await likeStore.likeRecipe({
        recipe_id: props.recipe.id,
        user_id: userStore.userId
      })
      toast.success('Recipe liked')
    }
    isLiked.value = !isLiked.value
  } catch (error) {
    toast.error('Failed to update like')
    console.error('Like error:', error)
  }
}
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.aspect-w-4 {
  position: relative;
  padding-bottom: 75%; /* 4:3 aspect ratio */
}

.aspect-w-4 > * {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}
</style>