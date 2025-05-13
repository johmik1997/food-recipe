<template>
  <div class="bg-white shadow-sm">
    <Navbar />  
    <section class="relative bg-gray-900 text-white overflow-hidden">
      <!-- Background Image with Overlay -->
      <div class="absolute inset-0 z-0">
        <img 
          src="https://images.unsplash.com/photo-1546069901-ba9599a7e63c?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1480&q=80" 
          alt="Delicious food assortment"
          class="w-full h-full object-cover opacity-50"
        />
      </div>

      <div class="relative z-10 max-w-7xl mx-auto px-4 py-24 sm:px-6 lg:px-8">
        <div class="text-center">
          <!-- Main Headline -->
          <h1 class="text-4xl md:text-6xl font-bold mb-4 tracking-tight">
            <span class="block">Discover & Share</span>
            <span class="block text-green-400">Delicious Recipes</span>
          </h1>

          <!-- Subtitle -->
          <p class="mt-6 max-w-lg mx-auto text-xl text-gray-300">
            Join our community of food lovers. Find inspiration or contribute your own culinary creations.
          </p>

          <!-- Conditional CTAs -->
          <div class="mt-10 flex flex-col sm:flex-row justify-center gap-4">
            <NuxtLink 
              to="/recipes" 
              class="px-8 py-3 border border-transparent text-base font-medium rounded-md text-white bg-green-600 hover:bg-green-700 md:py-4 md:text-lg md:px-10 transition-colors"
            >
              Browse Recipes
            </NuxtLink>

            <NuxtLink 
              v-if="!isAuthenticated"
              to="/register" 
              class="px-8 py-3 border border-transparent text-base font-medium rounded-md text-green-700 bg-white hover:bg-gray-100 md:py-4 md:text-lg md:px-10 transition-colors"
            >
              Share Your Recipe
            </NuxtLink>

            <NuxtLink 
              v-else
              to="/dashboard/recipes/create" 
              class="px-8 py-3 border border-transparent text-base font-medium rounded-md text-green-700 bg-white hover:bg-gray-100 md:py-4 md:text-lg md:px-10 transition-colors"
            >
              Create New Recipe
            </NuxtLink>
          </div>

          <!-- Stats Bar (Optional) -->
          <div class="mt-16 grid grid-cols-3 gap-8 text-center">
            <div>
              <p class="text-3xl font-bold text-green-400">1,000+</p>
              <p class="mt-2 text-gray-300">Recipes</p>
            </div>
            <div>
              <p class="text-3xl font-bold text-green-400">500+</p>
              <p class="mt-2 text-gray-300">Food Creators</p>
            </div>
            <div>
              <p class="text-3xl font-bold text-green-400">10K+</p>
              <p class="mt-2 text-gray-300">Meals Shared</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Scrolling Indicator (Optional) -->
      <div class="absolute bottom-10 left-1/2 transform -translate-x-1/2 z-10">
        <svg class="w-8 h-8 text-white animate-bounce" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 14l-7 7m0 0l-7-7m7 7V3"></path>
        </svg>
      </div>
    </section>

    <!-- Trending Recipes Section -->
    <section class="max-w-7xl mx-auto px-4 py-16 sm:px-6 lg:px-8">
      <div class="text-center mb-12">
        <h2 class="text-3xl font-bold text-gray-900">
          <span class="bg-gradient-to-r from-green-600 to-primary-600 bg-clip-text text-transparent">
            Trending Now
          </span>
        </h2>
        <p class="mt-4 max-w-2xl mx-auto text-lg text-gray-600">
          Discover what recipes are currently trending in our community
        </p>
      </div>

      <!-- Loading State -->
      <div v-if="trendingPending" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
        <div v-for="n in 4" :key="n" class="bg-white rounded-lg shadow-sm overflow-hidden">
          <div class="animate-pulse">
            <div class="h-48 bg-gray-200"></div>
            <div class="p-4">
              <div class="h-6 bg-gray-200 rounded w-3/4 mb-2"></div>
              <div class="h-4 bg-gray-200 rounded w-1/2"></div>
            </div>
          </div>
        </div>
      </div>

      <!-- Error State -->
      <div v-else-if="trendingError" class="text-center py-12">
        <div class="inline-flex items-center px-4 py-2 rounded-md bg-red-100 text-red-700">
          <Icon name="heroicons:exclamation-circle" class="w-5 h-5 mr-2" />
          Failed to load trending recipes
        </div>
      </div>

      <!-- Trending Recipes Grid -->
      <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
        <div 
          v-for="recipe in trendingRecipes" 
          :key="recipe.recipe_id"
          class="bg-white rounded-lg shadow-sm overflow-hidden hover:shadow-md transition-shadow"
        >
          <NuxtLink :to="`/recipes/${recipe.recipe_id}`">
            <div class="relative">
              <NuxtImg 
                :src="recipe.featured_image_url || '/placeholder-food.jpg'"
                :alt="recipe.title"
                class="w-full h-48 object-cover"
                width="400"
                height="192"
              />
              <div class="absolute bottom-2 left-2 bg-green-600 text-white text-xs font-semibold px-2 py-1 rounded">
                Trending
              </div>
              <div class="absolute top-2 right-2 bg-white/90 text-green-600 text-xs font-bold px-2 py-1 rounded-full">
                {{ Math.round(recipe.trending_score * 100) }}%
              </div>
            </div>
            <div class="p-4">
              <h3 class="text-lg font-semibold text-gray-900 mb-1 line-clamp-1">{{ recipe.title }}</h3>
              <div class="flex items-center text-sm text-gray-500">
                <Icon name="heroicons:fire" class="w-4 h-4 mr-1 text-orange-500" />
                <span>Trending score: {{ Math.round(recipe.trending_score * 100) }}</span>
              </div>
            </div>
          </NuxtLink>
        </div>
      </div>

      <!-- View All Button -->
      <div class="mt-10 text-center">
        <NuxtLink 
          to="/recipes"
          class="inline-flex items-center px-6 py-3 border border-transparent text-base font-medium rounded-md shadow-sm text-white bg-green-600 hover:bg-green-700"
        >
          View All Trending Recipes
          <Icon name="heroicons:arrow-right" class="ml-2 -mr-1 w-5 h-5" />
        </NuxtLink>
      </div>
    </section>
    <section class="max-w-7xl mx-auto px-4 py-16 sm:px-6 lg:px-8 bg-gray-50 rounded-lg">
      <div class="text-center mb-12">
        <h2 class="text-3xl font-bold text-gray-900">
          <span class="bg-gradient-to-r from-green-600 to-primary-600 bg-clip-text text-transparent">
            Community Favorites
          </span>
        </h2>
        <p class="mt-4 max-w-2xl mx-auto text-lg text-gray-600">
          Most bookmarked recipes by our community
        </p>
      </div>

      <div v-if="popularLoading" class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div v-for="n in 3" :key="n" class="bg-white rounded-lg shadow-sm overflow-hidden">
          <div class="animate-pulse">
            <div class="h-48 bg-gray-200"></div>
            <div class="p-4">
              <div class="h-6 bg-gray-200 rounded w-3/4 mb-2"></div>
              <div class="h-4 bg-gray-200 rounded w-1/2"></div>
            </div>
          </div>
        </div>
      </div>

      <div v-else-if="popularResult" class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <RecipeCard 
          v-for="recipe in popularResult.recipes" 
          :key="recipe.id"
          :recipe="recipe"
          :show-author="true"
        />
      </div>

      <div class="mt-10 text-center">
        <NuxtLink 
          to="/recipes?sort=popular"
          class="inline-flex items-center px-6 py-3 border border-transparent text-base font-medium rounded-md shadow-sm text-white bg-green-600 hover:bg-green-700"
        >
          View All Popular Recipes
          <Icon name="heroicons:arrow-right" class="ml-2 -mr-1 w-5 h-5" />
        </NuxtLink>
      </div>
    </section>

    <!-- Latest Recipes Section -->
    <section class="max-w-7xl mx-auto px-4 py-16 sm:px-6 lg:px-8">
      <div class="text-center mb-12">
        <h2 class="text-3xl font-bold text-gray-900">
          <span class="bg-gradient-to-r from-green-600 to-primary-600 bg-clip-text text-transparent">
            Fresh From the Community
          </span>
        </h2>
        <p class="mt-4 max-w-2xl mx-auto text-lg text-gray-600">
          The newest recipes shared by our community members
        </p>
      </div>

      <div v-if="latestLoading" class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div v-for="n in 3" :key="n" class="bg-white rounded-lg shadow-sm overflow-hidden">
          <div class="animate-pulse">
            <div class="h-48 bg-gray-200"></div>
            <div class="p-4">
              <div class="h-6 bg-gray-200 rounded w-3/4 mb-2"></div>
              <div class="h-4 bg-gray-200 rounded w-1/2"></div>
            </div>
          </div>
        </div>
      </div>

      <div v-else-if="latestResult" class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <RecipeCard 
          v-for="recipe in latestResult.recipes" 
          :key="recipe.id"
          :recipe="recipe"
          :show-author="true"
        />
      </div>

      <div class="mt-10 text-center">
        <NuxtLink 
          to="/recipes?sort=newest"
          class="inline-flex items-center px-6 py-3 border border-transparent text-base font-medium rounded-md shadow-sm text-white bg-green-600 hover:bg-green-700"
        >
          View All Recent Recipes
          <Icon name="heroicons:arrow-right" class="ml-2 -mr-1 w-5 h-5" />
        </NuxtLink>
      </div>
    </section>
  <section class="max-w-7xl mx-auto px-4 py-16 sm:px-6 lg:px-8 text-center">
    <div class="bg-gradient-to-r from-green-500 to-green-600 rounded-xl p-8 text-white">
      <h2 class="text-3xl font-bold mb-6">
        {{ isAuthenticated ? 'Ready to Share Your Next Creation?' : 'Join Our Vibrant Food Community' }}
      </h2>
      
      <!-- Loading State -->
      <div v-if="communityStatsLoading" class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <div v-for="n in 3" :key="n" class="bg-white/10 p-4 rounded-lg animate-pulse">
          <div class="h-10 bg-white/20 rounded mb-2 mx-auto w-1/2"></div>
          <div class="h-6 bg-white/20 rounded w-3/4 mx-auto"></div>
        </div>
      </div>

      <!-- Error State -->
      <div v-else-if="communityStatsError" class="mb-8">
        <div class="inline-flex items-center px-4 py-2 rounded-md bg-red-100 text-red-700">
          <Icon name="heroicons:exclamation-circle" class="w-5 h-5 mr-2" />
          Failed to load community stats
        </div>
      </div>

      <!-- Loaded State -->
      <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <div class="bg-white/10 p-4 rounded-lg">
          <p class="text-4xl font-bold">{{ formatNumber(communityStats?.recipes_aggregate.aggregate.count) || '0' }}</p>
          <p class="mt-2">Recipes Shared</p>
        </div>
        <div class="bg-white/10 p-4 rounded-lg">
          <p class="text-4xl font-bold">{{ formatNumber(communityStats?.users_aggregate.aggregate.count) || '0' }}</p>
          <p class="mt-2">Food Creators</p>
        </div>
        <div class="bg-white/10 p-4 rounded-lg">
          <p class="text-4xl font-bold">{{ formatNumber(communityStats?.user_bookmarks_aggregate.aggregate.count) || '0' }}</p>
          <p class="mt-2">Meals Saved</p>
        </div>
      </div>

      <!-- Rest of the content remains the same -->
      <div class="max-w-2xl mx-auto">
        <p class="text-xl mb-6">
          {{ isAuthenticated 
            ? 'Continue your culinary journey with us - share, explore, and get inspired!'
            : 'Discover, share, and connect with fellow food enthusiasts.' 
          }}
        </p>
        
        <div class="flex flex-col sm:flex-row justify-center gap-4">
          <NuxtLink 
            to="/recipes"
            class="px-6 py-3 border border-transparent text-base font-medium rounded-md shadow-sm text-green-600 bg-white hover:bg-gray-100"
          >
            Browse Recipes
          </NuxtLink>
          
          <NuxtLink 
            v-if="!isAuthenticated"
            to="/register" 
            class="px-6 py-3 border border-transparent text-base font-medium rounded-md shadow-sm text-white bg-green-700 hover:bg-green-800"
          >
            Join Now - It's Free!
          </NuxtLink>
          
          <NuxtLink 
            v-else
            to="/dashboard/recipes/create" 
            class="px-6 py-3 border border-transparent text-base font-medium rounded-md shadow-sm text-white bg-green-700 hover:bg-green-800"
          >
            Create New Recipe
          </NuxtLink>
        </div>
      </div>
    </div>
  </section>
  <Footer />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'

const isAuthenticated = ref(false)
const userId = ref(null)

// Trending recipes query
const TRENDING_RECIPES_QUERY = gql`
  query GetTrendingRecipes {
    trending_recipes(limit: 4) {
      recipe_id
      title
      featured_image_url
      trending_score
    }
  }
`

// Fetch trending recipes
const { 
  result: trendingResult, 
  loading: trendingPending, 
  error: trendingError 
} = useQuery(TRENDING_RECIPES_QUERY)

// Process trending recipes data
const trendingRecipes = computed(() => {
  return trendingResult.value?.trending_recipes || []
})

// Query for community popular recipes
const { result: popularResult, loading: popularLoading } = useQuery(gql`
  query GetCommunityPopularRecipes {
    recipes(
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
      user {
        id
        name
        avatar_image_url
      }
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
`)

// Query for community latest recipes
const { result: latestResult, loading: latestLoading } = useQuery(gql`
  query GetCommunityLatestRecipes {
    recipes(
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
      user {
        id
        name
        avatar_image_url
      }
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
`)
const COMMUNITY_STATS_QUERY = gql`
   query GetCommunityStats {
    recipes_aggregate {
      aggregate {
        count
      }
    }
    users_aggregate {
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
`

const { 
  result: communityStatsResult,
  loading: communityStatsLoading,
  error: communityStatsError 
} = useQuery(COMMUNITY_STATS_QUERY)

const communityStats = computed(() => communityStatsResult.value)

// Formatting helper function
const formatNumber = (num) => {
  if (!num) return '0'
  return new Intl.NumberFormat().format(num)
}
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
</script>