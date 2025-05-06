<!-- ~/pages/recipes/index.vue -->
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
            to="/signup" 
            class="px-8 py-3 border border-transparent text-base font-medium rounded-md text-green-700 bg-white hover:bg-gray-100 md:py-4 md:text-lg md:px-10 transition-colors"
          >
            Share Your Recipe
          </NuxtLink>

          <NuxtLink 
            v-else
            to="/dashboard/new-recipe" 
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
  </div>
  <Footer />
</template>

<script setup>
// Replace with your actual auth check
const isAuthenticated = ref(false)
</script>