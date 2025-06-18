<script setup>

const recipeStore = useRecipeStore()
const auth = authStore();

const isAuthenticated = computed(() => auth.isAuthed)
const userId = computed(() => auth.userId)


// Formatting helper function
const formatNumber = (num) => {
  if (!num) return '0'
  return new Intl.NumberFormat().format(num)
}
</script>

<template>
  <div class="bg-white shadow-sm">   
    <!-- Hero Section -->
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
          <h1 class="text-4xl md:text-6xl font-bold mb-4 tracking-tight">
            <span class="block">Discover & Share</span>
            <span class="block text-green-400">Delicious Recipes</span>
          </h1>

          <p class="mt-6 max-w-lg mx-auto text-xl text-gray-300">
            Join our community of food lovers. Find inspiration or contribute your own culinary creations.
          </p>

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
              Join Our Community
            </NuxtLink>

            <NuxtLink 
              v-else
              to="/recipes/create" 
              class="px-8 py-3 border border-transparent text-base font-medium rounded-md text-green-700 bg-white hover:bg-gray-100 md:py-4 md:text-lg md:px-10 transition-colors"
            >
              Share Your Recipe
            </NuxtLink>
          </div>
        </div>
      </div>
    </section>
<!-- Kushan Cuisine Section -->
    <section class="max-w-7xl mx-auto px-4 py-16 sm:px-6 lg:px-8">
      <div class="text-center mb-12">
        <h2 class="text-3xl font-bold text-gray-900 mb-4">
          Explore <span class="text-green-600">Kushan Recipe</span>
        </h2>
        <p class="text-xl text-gray-600 max-w-3xl mx-auto">
          Discover the rich flavors of ancient Kushan Empire - where Central Asian, Indian, and Persian culinary traditions blend harmoniously.
        </p>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
        <!-- Recipe Card 1 -->
        <div class="bg-white rounded-xl shadow-md overflow-hidden transition-transform duration-300 hover:scale-105">
          <img 
            src="https://images.unsplash.com/photo-1601050690597-df0568f70950?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1470&q=80" 
            alt="Kushan Dumplings"
            class="w-full h-48 object-cover"
          >
          <div class="p-6">
            <h3 class="text-xl font-bold text-gray-900 mb-2">Kushan Mantu</h3>
            <p class="text-gray-600 mb-4">Steamed dumplings filled with spiced lamb and onions, served with yogurt and mint oil.</p>
            <div class="flex justify-between items-center">
              <span class="text-green-600 font-medium">Difficulty: Medium</span>
              <NuxtLink 
                to="/recipes/" 
                class="text-green-600 hover:text-green-700 font-medium"
              >
                View Recipe →
              </NuxtLink>
            </div>
          </div>
        </div>

        <!-- Recipe Card 2 -->
        <div class="bg-white rounded-xl shadow-md overflow-hidden transition-transform duration-300 hover:scale-105">
          <img 
            src="https://images.unsplash.com/photo-1603105037880-880cd4edfb0d?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=687&q=80" 
            alt="Kushan Pilaf"
            class="w-full h-48 object-cover"
          >
          <div class="p-6">
            <h3 class="text-xl font-bold text-gray-900 mb-2">Kushan Osh</h3>
            <p class="text-gray-600 mb-4">Fragrant rice pilaf with lamb, carrots, and raisins - a royal dish from the Silk Road.</p>
            <div class="flex justify-between items-center">
              <span class="text-green-600 font-medium">Difficulty: Easy</span>
              <NuxtLink 
                to="/recipes/" 
                class="text-green-600 hover:text-green-700 font-medium"
              >
                View Recipe →
              </NuxtLink>
            </div>
          </div>
        </div>

        <!-- Recipe Card 3 -->
        <div class="bg-white rounded-xl shadow-md overflow-hidden transition-transform duration-300 hover:scale-105">
          <img 
            src="https://images.unsplash.com/photo-1565557623262-b51c2513a641?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1371&q=80" 
            alt="Kushan Bread"
            class="w-full h-48 object-cover"
          >
          <div class="p-6">
            <h3 class="text-xl font-bold text-gray-900 mb-2">Kushan Nan</h3>
            <p class="text-gray-600 mb-4">Traditional tandoor-baked flatbread with sesame and nigella seeds.</p>
            <div class="flex justify-between items-center">
              <span class="text-green-600 font-medium">Difficulty: Simple</span>
              <NuxtLink 
                to="/recipes/" 
                class="text-green-600 hover:text-green-700 font-medium"
              >
                View Recipe →
              </NuxtLink>
            </div>
          </div>
        </div>
      </div>

      <div class="text-center mt-12">
        <NuxtLink 
          to="/recipes" 
          class="inline-flex items-center px-6 py-3 border border-transparent text-base font-medium rounded-md shadow-sm text-white bg-green-600 hover:bg-green-700"
        >
          Explore More Kushan Recipes
          <svg class="ml-3 -mr-1 h-5 w-5" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M12.293 5.293a1 1 0 011.414 0l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414-1.414L14.586 11H3a1 1 0 110-2h11.586l-2.293-2.293a1 1 0 010-1.414z" clip-rule="evenodd" />
          </svg>
        </NuxtLink>
      </div>
    </section>
    <!-- CTA Section -->
    <section class="max-w-7xl mx-auto px-4 py-16 sm:px-6 lg:px-8 text-center">
      <div class="bg-gradient-to-r from-green-500 to-green-600 rounded-xl p-8 text-white">
        <h2 class="text-3xl font-bold mb-6">
          {{ isAuthenticated ? 'Ready to Share Your Next Creation?' : 'Join Our Vibrant Food Community' }}
        </h2>
        
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
              to="/auth" 
              class="px-6 py-3 border border-transparent text-base font-medium rounded-md shadow-sm text-white bg-green-700 hover:bg-green-800"
            >
              Join Now - It's Free!
            </NuxtLink>
            
            <NuxtLink 
              v-else
              to="/recipes/create" 
              class="px-6 py-3 border border-transparent text-base font-medium rounded-md shadow-sm text-white bg-green-700 hover:bg-green-800"
            >
              Create New Recipe
            </NuxtLink>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>