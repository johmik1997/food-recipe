<template>
    <div class="min-h-screen bg-gray-50">
    <Navbar />
    
  <div class="max-w-4xl mx-auto px-4 py-8 md:py-12">
 
    <!-- Loading state -->
    <div v-if="loading" class="text-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500 mx-auto"></div>
      <p class="mt-4 text-gray-600">Loading recipe...</p>
    </div>

    <!-- Error state -->
    <div v-if="error" class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
      <p class="font-medium">Error loading recipe: {{ error.message }}</p>
      <button 
        @click="fetchRecipe" 
        class="mt-2 px-4 py-2 bg-red-100 hover:bg-red-200 rounded-lg text-red-700 font-medium transition-colors"
      >
        Retry
      </button>
    </div>

    <!-- Recipe content -->
    <div v-if="recipe" class="space-y-10">
      <!-- Recipe header -->
      <div class="space-y-4">
        <div>
          <h1 class="text-3xl md:text-4xl font-bold text-gray-900 leading-tight">{{ recipe.title }}</h1>
          
          <div class="mt-3 flex flex-wrap items-center gap-x-6 gap-y-2 text-gray-600">
            <span v-if="recipe.prep_time" class="flex items-center">
              <ClockIcon class="h-5 w-5 mr-2 text-gray-500" />
              Prep: {{ recipe.prep_time }} min
            </span>
            <span v-if="recipe.cook_time" class="flex items-center">
              <FireIcon class="h-5 w-5 mr-2 text-gray-500" />
              Cook: {{ recipe.cook_time }} min
            </span>
            <span v-if="recipe.servings" class="flex items-center">
              <UsersIcon class="h-5 w-5 mr-2 text-gray-500" />
              Serves: {{ recipe.servings }}
            </span>
          </div>
        </div>

        <!-- Featured image -->
        <div v-if="recipe.featured_image_url" class="rounded-xl overflow-hidden shadow-lg">
          <img 
            :src="recipe.featured_image_url" 
            :alt="recipe.title" 
            class="w-full h-96 object-cover"
            loading="lazy"
          />
        </div>
      </div>

      <!-- Description -->
      <div v-if="recipe.description" class="prose max-w-none text-lg text-gray-700 leading-relaxed">
        <p>{{ recipe.description }}</p>
      </div>

      <!-- Ingredients -->
      <div class="bg-gray-50 p-6 rounded-xl border border-gray-100 shadow-sm">
        <h2 class="text-2xl font-bold mb-5 text-gray-900 flex items-center">
          <ClipboardListIcon class="h-6 w-6 mr-2 text-blue-600" />
          Ingredients
        </h2>
        <ul class="space-y-3">
          <li 
            v-for="(ingredient, index) in recipe.recipe_ingredients" 
            :key="index" 
            class="flex items-start group"
          >
            <span class="inline-flex items-center justify-center h-5 w-5 rounded-full bg-blue-100 mt-0.5 mr-3 flex-shrink-0 group-hover:bg-blue-200 transition-colors">
              <span class="w-1.5 h-1.5 bg-blue-600 rounded-full"></span>
            </span>
            <span class="text-gray-700">
              <span v-if="ingredient.quantity" class="font-medium">{{ ingredient.quantity }}</span>
              <span v-if="ingredient.unit" class="ml-1 text-gray-600">{{ ingredient.unit }}</span>
              <span class="ml-2">{{ ingredient.name }}</span>
            </span>
          </li>
        </ul>
      </div>

      <!-- Steps -->
      <div class="space-y-8">
        <h2 class="text-2xl font-bold mb-6 text-gray-900 flex items-center">
          <BookOpenIcon class="h-6 w-6 mr-2 text-blue-600" />
          Instructions
        </h2>
        <div 
          v-for="step in recipe.recipe_steps" 
          :key="step.step_number" 
          class="flex group hover:bg-gray-50 rounded-lg p-3 transition-colors"
        >
          <div class="flex-shrink-0 mr-4">
            <span class="flex items-center justify-center w-9 h-9 rounded-full bg-blue-600 text-white font-bold group-hover:bg-blue-700 transition-colors">
              {{ step.step_number }}
            </span>
          </div>
          <div class="prose max-w-none flex-1">
            <p class="text-gray-800">{{ step.instruction }}</p>
            <img 
              v-if="step.image_url" 
              :src="step.image_url" 
              :alt="`Step ${step.step_number}`"
              class="mt-4 rounded-lg shadow-md max-w-full border border-gray-200"
              loading="lazy"
            />
          </div>
        </div>
      </div>
      
      <!-- Additional images -->
      <div v-if="recipe.images && recipe.images.length > 0" class="grid grid-cols-2 md:grid-cols-3 gap-4">
        <div 
          v-for="(image, index) in recipe.images" 
          :key="index"
          class="overflow-hidden rounded-xl shadow-md hover:shadow-lg transition-shadow"
        >
          <img 
            :src="image.image_url" 
            :alt="`Recipe image ${index + 1}`"
            class="w-full h-48 md:h-56 object-cover hover:scale-105 transition-transform duration-300"
            loading="lazy"
          />
        </div>
      </div>

      <!-- Reviews Section -->
      <div class="mt-12 pt-8 border-t border-gray-200">
        <div class="flex flex-col md:flex-row md:items-center md:justify-between mb-8 gap-4">
          <h2 class="text-2xl font-bold text-gray-900 flex items-center">
            <StarIcon class="h-6 w-6 mr-2 text-blue-600" />
            Reviews
          </h2>
          <div class="flex flex-col sm:flex-row sm:items-center gap-4">
            <div class="flex items-center">
              <div class="flex items-center mr-2">
                <StarIcon v-for="i in 5" :key="i" 
                  class="h-5 w-5"
                  :class="i <= Math.round(recipe.average_rating) ? 'text-yellow-400' : 'text-gray-300'"
                />
              </div>
              <span class="text-gray-700 ml-1">
                <span class="font-medium">{{ recipe.average_rating?.toFixed(1) || '0.0' }}</span>
                <span class="text-gray-500"> ({{ recipe.review_count || 0 }} reviews)</span>
              </span>
            </div>
            <button 
              @click="handleReviewButtonClick"
              class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors shadow-sm flex items-center justify-center"
            >
              <PencilIcon class="h-5 w-5 mr-2" />
              Write a Review
            </button>
          </div>
        </div>

        <!-- Review Form -->
        <div 
          v-if="showReviewForm" 
          class="bg-white p-6 rounded-xl shadow-md mb-8 border border-gray-100"
        >
          <h3 class="text-xl font-medium mb-5 text-gray-900">Add Your Review</h3>
          <form @submit.prevent="submitReview">
            <div class="mb-5">
              <label class="block text-gray-700 mb-2 font-medium">Rating <span class="text-red-500">*</span></label>
              <div class="flex">
                <button 
                  v-for="i in 5" 
                  :key="i"
                  type="button"
                  @click="newReview.rating = i"
                  class="focus:outline-none transform hover:scale-110 transition-transform"
                >
                  <StarIcon 
                    class="h-8 w-8"
                    :class="i <= newReview.rating ? 'text-yellow-400' : 'text-gray-300'"
                  />
                </button>
              </div>
              <p v-if="!newReview.rating && formSubmitted" class="text-red-500 text-sm mt-1">
                Please select a rating
              </p>
            </div>
            <div class="mb-5">
              <label for="comment" class="block text-gray-700 mb-2 font-medium">Comment (optional)</label>
              <textarea
                id="comment"
                v-model="newReview.comment"
                rows="4"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition"
                placeholder="Share your thoughts about this recipe..."
              ></textarea>
            </div>
            <div class="flex justify-end gap-3">
              <button
                type="button"
                @click="showReviewForm = false"
                class="px-4 py-2 text-gray-700 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
              >
                Cancel
              </button>
              <button
                type="submit"
                class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 transition-colors flex items-center"
                :disabled="submittingReview"
              >
                <span v-if="submittingReview" class="flex items-center">
                  <SpinnerIcon class="animate-spin h-4 w-4 mr-2" />
                  Submitting...
                </span>
                <span v-else>Submit Review</span>
              </button>
            </div>
          </form>
        </div>

        <!-- Reviews List -->
        <div v-if="recipe.combined_reviews?.length > 0" class="space-y-6">
          <div 
            v-for="review in recipe.combined_reviews" 
            :key="review.id" 
            class="bg-white p-6 rounded-xl shadow-sm border border-gray-100 hover:border-gray-200 transition-colors"
          >
            <div class="flex items-start">
              <img 
                :src="review.user.avatar_image_url || '/placeholder-user.jpg'" 
                :alt="review.user.name"
                class="h-10 w-10 rounded-full object-cover mr-4 flex-shrink-0"
              />
              <div class="flex-1">
                <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between">
                  <h4 class="font-medium text-gray-900">{{ review.user.name }}</h4>
                  <span class="text-sm text-gray-500 mt-1 sm:mt-0">{{ formatDate(review.created_at) }}</span>
                </div>
                <div v-if="review.rating" class="flex items-center mt-2 mb-3">
                  <StarIcon v-for="i in 5" :key="i" 
                    class="h-5 w-5"
                    :class="i <= review.rating ? 'text-yellow-400' : 'text-gray-300'"
                  />
                </div>
                <p class="text-gray-700 whitespace-pre-line">{{ review.comment }}</p>
              </div>
            </div>
          </div>
        </div>

        <div v-else class="bg-white p-8 text-center rounded-xl shadow-sm border border-gray-100">
          <ChatBubbleBottomCenterTextIcon class="h-10 w-10 mx-auto text-gray-400" />
          <p class="text-gray-600 mt-3">No reviews yet. Be the first to review!</p>
        </div>
      </div>
    </div>
  </div>
<Footer/>
    </div>
</template>

<script setup>
import { useRoute } from 'vue-router'
import { ClockIcon, FireIcon, UsersIcon, StarIcon } from '@heroicons/vue/outline'
import { computed, ref, onMounted } from 'vue'
import gql from 'graphql-tag'
import { useQuery, useMutation } from '@vue/apollo-composable'
import { useHead } from '#imports'

const route = useRoute()
const recipeId = route.params.id
const userId = ref(null)
const isAuthenticated = ref(false)

// Review form state
const showReviewForm = ref(false)
const newReview = ref({
  rating: 0,
  comment: ''
})
const submittingReview = ref(false)
const formSubmitted = ref(false)

// Get user ID from localStorage when component mounts
onMounted(() => {
  const userStr = localStorage.getItem("user")
  if (userStr) {
    try {
      const user = JSON.parse(userStr)
      userId.value = user.id
      isAuthenticated.value = true
    } catch (e) {
      console.error("Failed to parse user data", e)
    }
  }
})

// GraphQL queries
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
      comments(order_by: {created_at: desc}) {
        id
        content
        created_at
        user {
          id
          name
          avatar_image_url
        }
      }
      ratings {
        value
        user_id
      }
      ratings_aggregate {
        aggregate {
          avg {
            value
          }
          count
        }
      }
    }
  }
`

const ADD_RATING = gql`
  mutation AddRating($recipeId: uuid!, $value: Int!, $userId: uuid!) {
    insert_ratings_one(object: {
      recipe_id: $recipeId,
      value: $value,
      user_id: $userId
    }) {
      id
      value
      created_at
      user {
        id
        name
      }
    }
  }
`

const ADD_COMMENT = gql`
  mutation AddComment($recipeId: uuid!, $content: String!, $userId: uuid!) {
    insert_comments_one(object: {
      recipe_id: $recipeId,
      content: $content,
      user_id: $userId
    }) {
      id
      content
      created_at
      user {
        id
        name
        avatar_image_url
      }
    }
  }
`

const { result, loading, error, refetch } = useQuery(GET_RECIPE_BY_ID, {
  id: recipeId
})

const { mutate: addRating } = useMutation(ADD_RATING)
const { mutate: addComment } = useMutation(ADD_COMMENT)

const recipe = computed(() => {
  if (!result.value?.recipes_by_pk) return null
  
  return {
    ...result.value.recipes_by_pk,
    average_rating: result.value.recipes_by_pk.ratings_aggregate.aggregate?.avg?.value || 0,
    review_count: result.value.recipes_by_pk.ratings_aggregate.aggregate?.count || 0,
    combined_reviews: result.value.recipes_by_pk.comments.map(comment => {
      if (!comment || !comment.user) {
        return {
          id: comment?.id || 'unknown',
          user: {
            id: 'unknown',
            name: 'Anonymous',
            avatar_image_url: '/placeholder-user.jpg'
          },
          rating: null,
          comment: comment?.content || '',
          created_at: comment?.created_at || new Date().toISOString()
        }
      }
      const userRating = result.value.recipes_by_pk.ratings.find(
        r => r.user_id === comment.user.id
      )
      return {
        id: comment.id,
        user: comment.user,
        rating: userRating?.value || null,
        comment: comment.content,
        created_at: comment.created_at
      }
    })
  }
})

// Handle review button click with authentication check
const handleReviewButtonClick = () => {
  if (!isAuthenticated.value) {
    alert('Please log in to write a review')
    return
  }
  showReviewForm.value = !showReviewForm.value
}

const submitReview = async () => {
  // Check if user is authenticated
  if (!isAuthenticated.value) {
    alert('Please log in to submit a review')
    return
  }

  formSubmitted.value = true
  
  if (!newReview.value.rating) {
    return
  }

  submittingReview.value = true
  try {
    // Submit rating
    await addRating({
      recipeId: recipeId,
      value: newReview.value.rating,
      userId: userId.value
    })
    
    // Submit comment if provided
    if (newReview.value.comment.trim()) {
      await addComment({
        recipeId: recipeId,
        content: newReview.value.comment,
        userId: userId.value
      })
    }
    
    // Reset form and refresh data
    showReviewForm.value = false
    newReview.value = { rating: 0, comment: '' }
    formSubmitted.value = false
    await refetch()
  } catch (error) {
    console.error('Error submitting review:', error)
    alert('Failed to submit review. Please try again.')
  } finally {
    submittingReview.value = false
  }
}

const fetchRecipe = () => {
  refetch()
}

const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', { 
    year: 'numeric', 
    month: 'short', 
    day: 'numeric' 
  })
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