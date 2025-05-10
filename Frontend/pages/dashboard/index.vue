<template>
  <div>
    <Navbar />
    <section class="bg-gradient-to-r from-green-50 to-primary-50 rounded-2xl shadow-sm">
      <div class="max-w-7xl mx-auto px-4 py-12 sm:px-6 lg:px-8">
        <!-- Error State -->
        <div v-if="error" class="text-center py-12 text-red-600">
          {{ error }}
        </div>

        <!-- Loading State -->
        <div v-else-if="loading" class="text-center py-12">
          <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary-600 mx-auto mb-4"></div>
          <p class="text-gray-600">Loading your delicious dashboard...</p>
        </div>

        <!-- Content -->
        <div v-else>
          <div class="flex flex-col md:flex-row justify-between items-center gap-8">
            <!-- Welcome Message -->
            <div class="flex-1">
              <h1 class="text-3xl md:text-4xl font-bold text-gray-900 mb-2">
                Welcome back, <span class="text-primary-600">{{ userData.name }}</span> 👋
              </h1>
              <p class="text-lg text-gray-600 mb-6">
                {{ welcomeMessage }}
              </p>
              
              <!-- Quick Actions -->
              <div class="flex flex-wrap gap-4">
                <NuxtLink 
                  to="/dashboard/recipes/create" 
                  class="inline-flex items-center px-6 py-3 bg-gradient-to-r from-green-600 to-primary-600 text-white font-medium rounded-lg shadow-md hover:from-green-700 hover:to-primary-700 transition-all transform hover:-translate-y-0.5"
                >
                  <PlusIcon class="w-5 h-5 mr-2" />
                  Create New Recipe
                </NuxtLink>
                
                <NuxtLink 
                  to="/dashboard/my-recipes" 
                  class="inline-flex items-center px-6 py-3 bg-white text-gray-700 font-medium rounded-lg border border-gray-300 shadow-sm hover:bg-gray-50 transition-colors"
                >
                  <BookOpenIcon class="w-5 h-5 mr-2 text-primary-500" />
                  My Recipes
                </NuxtLink>
              </div>
            </div>
            
            <!-- Stats Cards -->
            <div class="grid grid-cols-3 gap-4 w-full md:w-auto">
              <div 
                v-for="stat in stats" 
                :key="stat.label"
                class="bg-white p-4 rounded-lg shadow-sm border border-gray-100 text-center hover:shadow-md transition-shadow cursor-pointer"
                @click="handleStatClick(stat.action)"
              >
                <p class="text-2xl font-bold text-primary-600">{{ stat.value }}</p>
                <p class="text-sm text-gray-500">{{ stat.label }}</p>
              </div>
            </div>
          </div>
          
          <!-- Recipe Sections -->
          <div class="mt-12 space-y-12">
            <!-- Recent Recipes -->
            <div v-if="recentRecipes.length > 0">
              <div class="flex justify-between items-center mb-6">
                <h3 class="text-xl font-bold text-gray-900">Your Recent Recipes</h3>
                <NuxtLink 
                  to="recipes" 
                  class="text-sm text-primary-600 hover:text-primary-700 font-medium"
                >
                  View All Recipes →
                </NuxtLink>
              </div>
              <div class="grid gap-8 md:grid-cols-2 lg:grid-cols-3">
                <RecipeCard 
                  v-for="recipe in recentRecipes" 
                  :key="recipe.id"
                  :recipe="recipe"
                  :show-actions="true"
                  @edit="handleEditRecipe(recipe.id)"
                  @delete="handleDeleteRecipe(recipe.id)"
                />
              </div>
            </div>

            <!-- Popular Recipes -->
            <div v-if="popularRecipes.length > 0">
              <div class="flex justify-between items-center mb-6">
                <h3 class="text-xl font-bold text-gray-900">Most Popular Recipes</h3>
                <NuxtLink 
                  to="/recipes/popular" 
                  class="text-sm text-primary-600 hover:text-primary-700 font-medium"
                >
                  View All Popular →
                </NuxtLink>
              </div>
              <div class="grid gap-8 md:grid-cols-2 lg:grid-cols-3">
                <RecipeCard 
                  v-for="recipe in popularRecipes" 
                  :key="recipe.id"
                  :recipe="recipe"
                  :show-actions="true"
                />
              </div>
            </div>

            <!-- Favorite Recipes -->
            <div v-if="favoriteRecipes.length > 0">
              <div class="flex justify-between items-center mb-6">
                <h3 class="text-xl font-bold text-gray-900">Most Liked Recipes</h3>
                <NuxtLink 
                  to="/recipes/favorites" 
                  class="text-sm text-primary-600 hover:text-primary-700 font-medium"
                >
                  View All Favorites →
                </NuxtLink>
              </div>
              <div class="grid gap-8 md:grid-cols-2 lg:grid-cols-3">
                <RecipeCard 
                  v-for="recipe in favoriteRecipes" 
                  :key="recipe.id"
                  :recipe="recipe"
                  :show-actions="true"
                />
              </div>
            </div>
          </div>

          <!-- Empty State -->
          <div v-if="recentRecipes.length === 0 && popularRecipes.length === 0 && favoriteRecipes.length === 0" 
               class="mt-12 text-center py-12">
            <div class="mx-auto w-24 h-24 bg-gray-100 rounded-full flex items-center justify-center mb-4">
              <UtensilsIcon class="w-10 h-10 text-gray-400" />
            </div>
            <h3 class="text-lg font-medium text-gray-900 mb-2">No Recipes Yet</h3>
            <p class="text-gray-600 mb-6">Start by creating your first delicious recipe!</p>
            <NuxtLink 
              to="/dashboard/recipes/create" 
              class="inline-flex items-center px-6 py-3 bg-gradient-to-r from-green-600 to-primary-600 text-white font-medium rounded-lg shadow-md hover:from-green-700 hover:to-primary-700 transition-all"
            >
              <PlusIcon class="w-5 h-5 mr-2" />
              Create Your First Recipe
            </NuxtLink>
          </div>
        </div>
      </div>
    </section>
    <Footer />
  </div>
</template>

<script setup>
import { PlusIcon, BookOpenIcon} from '@heroicons/vue/outline'
import { useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'

const userId = ref(null)
const loading = ref(true)
const error = ref(null)

// GraphQL queries
const USER_DATA_QUERY = gql`
  query GetUserDashboardData($userId: uuid!) {
    user: users_by_pk(id: $userId) {
      id
      name
      avatar_image_url
    }
    recipes: recipes_aggregate(where: { user_id: { _eq: $userId } }) {
      aggregate {
        count
      }
    }
    likes: user_likes_aggregate(where: { user_id: { _eq: $userId } }) {
      aggregate {
        count
      }
    }
    bookmarks: user_bookmarks_aggregate(where: { user_id: { _eq: $userId } }) {
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
`

const RECENT_RECIPES_QUERY = gql`
  query GetRecentRecipes($userId: uuid!) {
    recipes(
      where: { user_id: { _eq: $userId } }
      order_by: { created_at: desc }
      limit: 3
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
      user_likes_aggregate {
        aggregate {
          count
        }
      }
      user_bookmarks_aggregate {
        aggregate {
          count
        }
      }
    }
  }
`

const POPULAR_RECIPES_QUERY = gql`
  query GetPopularRecipes($userId: uuid!) {
    recipes(
      where: { user_id: { _eq: $userId } }
      order_by: { user_bookmarks_aggregate: { count: desc } }
      limit: 3
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
      user_likes_aggregate {
        aggregate {
          count
        }
      }
      user_bookmarks_aggregate {
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

const FAVORITE_RECIPES_QUERY = gql`
  query GetFavoriteRecipes($userId: uuid!) {
    recipes(
      where: { user_id: { _eq: $userId } }
      order_by: { user_likes_aggregate: { count: desc } }
      limit: 3
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
      user_likes_aggregate {
        aggregate {
          count
        }
      }
      user_bookmarks_aggregate {
        aggregate {
          count
        }
      }
    }
  }
`

// Set userId from localStorage
onMounted(() => {
  const userStr = localStorage.getItem("user")
  if (userStr) {
    try {
      const user = JSON.parse(userStr)
      userId.value = user.id
      fetchData()
    } catch (e) {
      error.value = 'Failed to parse user data'
      loading.value = false
    }
  } else {
    error.value = 'No user data found'
    loading.value = false
  }
})

// Fetch all data
let userResult, recentResult, popularResult, favoriteResult

const fetchData = () => {
  if (!userId.value) return
  
  const { result: userRes } = useQuery(USER_DATA_QUERY, () => ({
    userId: userId.value
  }))
  userResult = userRes
  
  const { result: recentRes } = useQuery(RECENT_RECIPES_QUERY, () => ({
    userId: userId.value
  }))
  recentResult = recentRes
  
  const { result: popularRes } = useQuery(POPULAR_RECIPES_QUERY, () => ({
    userId: userId.value
  }))
  popularResult = popularRes
  
  const { result: favoriteRes } = useQuery(FAVORITE_RECIPES_QUERY, () => ({
    userId: userId.value
  }))
  favoriteResult = favoriteRes
  
  // Watch for all results to complete
  watch([userResult, recentResult, popularResult, favoriteResult], () => {
    if (userResult.value && recentResult.value && popularResult.value && favoriteResult.value) {
      loading.value = false
    }
  }, { immediate: true })
}

// Computed properties
const userData = computed(() => ({
  name: userResult?.value?.user?.name || 'Chef',
  recipeCount: userResult?.value?.recipes?.aggregate?.count || 0,
  likeCount: userResult?.value?.likes?.aggregate?.count || 0,
  bookmarkCount: userResult?.value?.bookmarks?.aggregate?.count || 0
}))

const recentRecipes = computed(() => {
  return recentResult?.value?.recipes?.map(recipe => ({
    ...recipe,
    likes_count: recipe.user_likes_aggregate?.aggregate?.count || 0,
    bookmarks_count: recipe.user_bookmarks_aggregate?.aggregate?.count || 0
  })) || []
})

const popularRecipes = computed(() => {
  return popularResult?.value?.recipes?.map(recipe => ({
    ...recipe,
    likes_count: recipe.user_likes_aggregate?.aggregate?.count || 0,
    bookmarks_count: recipe.user_bookmarks_aggregate?.aggregate?.count || 0
  })) || []
})

const favoriteRecipes = computed(() => {
  return favoriteResult?.value?.recipes?.map(recipe => ({
    ...recipe,
    likes_count: recipe.user_likes_aggregate?.aggregate?.count || 0,
    bookmarks_count: recipe.user_bookmarks_aggregate?.aggregate?.count || 0
  })) || []
})

// Welcome message
const welcomeMessage = computed(() => {
  const hour = new Date().getHours()
  if (hour < 12) return 'Good morning! What will you cook today?'
  if (hour < 18) return 'Good afternoon! Ready to create something tasty?'
  return 'Good evening! Time to plan tomorrow\'s meals!'
})

// Stats - now only 3 stats (removed followers)
const stats = computed(() => [
  { label: 'Recipes', value: userData.value.recipeCount, action: 'viewRecipes' },
  { label: 'Likes', value: userData.value.likeCount, action: 'viewLikes' },
  { label: 'Bookmarks', value: userData.value.bookmarkCount, action: 'viewBookmarks' }
])

// Handle stat card clicks
const handleStatClick = (action) => {
  switch (action) {
    case 'viewRecipes':
      navigateTo('/dashboard/my-recipes')
      break
    case 'viewLikes':
      navigateTo('/dashboard/likes')
      break
    case 'viewBookmarks':
      navigateTo('/dashboard/bookmarks')
      break
  }
}

// Recipe actions
const handleEditRecipe = (id) => {
  navigateTo(`/dashboard/recipes/edit/${id}`)
}

const handleDeleteRecipe = (id) => {
  // Implement delete logic
  // After deletion, you might want to refetch the data
  fetchData()
}
</script>

<style scoped>
/* Custom transitions for smoother hover effects */
.recipe-card-enter-active,
.recipe-card-leave-active {
  transition: all 0.3s ease;
}
.recipe-card-enter-from,
.recipe-card-leave-to {
  opacity: 0;
  transform: translateY(10px);
}
</style>