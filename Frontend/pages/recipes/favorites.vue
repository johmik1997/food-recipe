<template>
  <div class="min-h-screen bg-gray-50">
    <Navbar />
    
    <div class="py-12">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <!-- Header -->
        <div class="text-center mb-16">
          <h1 class="text-4xl md:text-5xl font-bold text-gray-900 mb-4">
            <span class="bg-gradient-to-r from-green-600 to-primary-600 bg-clip-text text-transparent">
              Most Loved Recipes
            </span>
          </h1>
          <p class="max-w-2xl mx-auto text-lg text-gray-600">
            Discover recipes with the most likes from our community
          </p>
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
            <Icon name="material-symbols:error-outline-rounded" class="mx-auto h-16 w-16 text-red-500" />
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
        <div v-else-if="recipes.length === 0" class="text-center py-16 bg-white rounded-xl shadow-sm border border-gray-100">
          <div class="max-w-md mx-auto">
            <Icon name="material-symbols:favorite-outline" class="mx-auto h-16 w-16 text-gray-400" />
            <h3 class="mt-4 text-lg font-medium text-gray-900">No favorite recipes yet</h3>
            <p class="mt-2 text-sm text-gray-600">
              Recipes will appear here as they get liked by users
            </p>
            <NuxtLink 
              to="/recipes" 
              class="mt-6 inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
            >
              Browse All Recipes
            </NuxtLink>
          </div>
        </div>

        <!-- Recipe Grid -->
        <div v-else class="grid gap-8 md:grid-cols-2 lg:grid-cols-3">
          <RecipeCard 
            v-for="recipe in recipes" 
            :key="recipe.id"
            :recipe="{
              ...recipe,
              is_liked: false,
              is_bookmarked: false,
              ratings_aggregate: recipe.ratings_aggregate || { aggregate: null }
            }"
            :show-actions="true"
            @update="handleRecipeUpdate"
          />
        </div>

        <!-- Pagination -->
        <div v-if="recipes.length > 0" class="mt-12 flex justify-center">
          <nav class="inline-flex rounded-md shadow-sm -space-x-px">
            <button 
              @click="prevPage" 
              :disabled="currentPage === 1"
              class="px-4 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Previous
            </button>
            <button 
              v-for="page in visiblePages" 
              :key="page"
              @click="goToPage(page)"
              :class="[
                'px-4 py-2 border-t border-b border-gray-300 text-sm font-medium',
                currentPage === page 
                  ? 'bg-primary-600 text-white' 
                  : 'bg-white text-gray-500 hover:bg-gray-50'
              ]"
            >
              {{ page }}
            </button>
            <button 
              @click="nextPage" 
              :disabled="currentPage === totalPages"
              class="px-4 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
            >
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
const itemsPerPage = 12
const currentPage = ref(1)
import { computed } from 'vue'
import { useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'
// GraphQL query for favorite recipes (sorted by likes)
const FAVORITE_RECIPES_QUERY = gql`
  query GetFavoriteRecipes($limit: Int!, $offset: Int!) {
    recipes(
      order_by: { 
        user_likes_aggregate: { count: desc },
        user_bookmarks_aggregate: { count: desc }
      }
      limit: $limit
      offset: $offset
    ) {
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
        id
        name
        avatar_image_url
      }
      user_bookmarks_aggregate {
        aggregate {
          count
        }
      }
      user_likes_aggregate {
        aggregate {
          count
        }
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

// Apollo query
const { result, loading: pending, error, refetch } = useQuery(
  FAVORITE_RECIPES_QUERY,
  () => ({
    limit: itemsPerPage,
    offset: (currentPage.value - 1) * itemsPerPage
  })
)

// Process recipes data
const recipes = computed(() => {
  if (!result.value?.recipes) return []
  
  return result.value.recipes.map(recipe => {
    return {
      ...recipe,
      bookmarks_count: recipe.user_bookmarks_aggregate?.aggregate?.count || 0,
      likes_count: recipe.user_likes_aggregate?.aggregate?.count || 0,
      rating: recipe.ratings_aggregate?.aggregate?.avg?.value || 0,
      rating_count: recipe.ratings_aggregate?.aggregate?.count || 0
    }
  })
})

// Total recipes count
const totalRecipes = computed(() => result.value?.recipes_aggregate?.aggregate?.count || 0)

// Total pages
const totalPages = computed(() => Math.ceil(totalRecipes.value / itemsPerPage))

// Visible pages for pagination (max 5 pages shown)
const visiblePages = computed(() => {
  const range = 2 // Number of pages to show before/after current
  const start = Math.max(1, currentPage.value - range)
  const end = Math.min(totalPages.value, currentPage.value + range)
  
  const pages = []
  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  return pages
})

// Pagination methods
const prevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }
}

const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }
}

const goToPage = (page) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }
}

// Handle recipe updates from child components
const handleRecipeUpdate = () => {
  refetch()
}
</script>