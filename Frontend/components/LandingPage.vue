<script setup>
import { gsap } from "gsap";
import { onMounted, ref, computed } from "vue";
import {
  RecpiesAllRecipes,
  FilterBreakfast,
  FilterDesserts,
  FilterDinner,
  FilterFasting,
  FilterLunch,
  FilterNonFasting,
} from "#components";
import { useRecipeStore, authStore } from "#imports";

const recipesStore = useRecipeStore();
const auth = authStore();

// Refs for DOM elements
const heroSection = ref(null);
const title = ref(null);
const subtitle = ref(null);
const cta = ref(null);
const feature1 = ref(null);
const feature2 = ref(null);
const feature3 = ref(null);

// Tab and search functionality
const selectedTab = ref("all");
const searchQuery = ref("");
const isAuthenticated = computed(() => auth.isAuthed);
const userId = computed(() => auth.userId);

// Formatting helper function
const formatNumber = (num) => {
  if (!num) return '0'
  return new Intl.NumberFormat().format(num)
}

// Component selection based on tab
const selectedComponent = computed(() => {
  switch (selectedTab.value) {
    case "all":
      return RecpiesAllRecipes;
    case "breakfast":
      return FilterBreakfast;
    case "lunch":
      return FilterLunch;
    case "desserts":
      return FilterDesserts;
    case "fasting":
      return FilterFasting;
    case "non-fasting":
      return FilterNonFasting;
    case "dinner":
      return FilterDinner;
    default:
      return RecpiesAllRecipes;
  }
});

const handleSearch = async () => {
  recipesStore.setSearchRecipe(searchQuery.value);
  await recipesStore.getAllRecipes(searchQuery.value);
};

const fetchCategories = async () => {
  await recipesStore.getCategories();
};

// Mounted lifecycle hook
onMounted(() => {
  // Hero Section Animations
  gsap.from(heroSection.value, {
    duration: 1,
    scale: 1.2,
    opacity: 0,
    ease: "power2.out",
  });

  gsap.from(title.value, {
    duration: 1,
    y: -50,
    opacity: 0,
    ease: "power2.out",
  });

  gsap.from(subtitle.value, {
    duration: 1,
    y: 50,
    opacity: 0,
    ease: "power2.out",
    delay: 0.5,
  });

  gsap.from(cta.value, {
    duration: 1,
    y: 50,
    opacity: 0,
    ease: "power2.out",
    delay: 1,
  });

  // Features Section Animations
  gsap.from(feature1.value, {
    duration: 1,
    x: -100,
    opacity: 0,
    ease: "power2.out",
    delay: 1.5,
  });

  gsap.from(feature2.value, {
    duration: 1,
    y: 100,
    opacity: 0,
    ease: "power2.out",
    delay: 2,
  });

  gsap.from(feature3.value, {
    duration: 1,
    x: 100,
    opacity: 0,
    ease: "power2.out",
    delay: 2.5,
  });
  
  fetchCategories();
});
</script>

<template>
  <div class="bg-white shadow-sm">
    <!-- Hero Section -->
    <section ref="heroSection" class="relative bg-gray-900 text-white overflow-hidden">
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
          <h1 ref="title" class="text-4xl md:text-6xl font-bold mb-4 tracking-tight">
            <span class="block">Discover & Share</span>
            <span class="block text-green-400">Delicious Recipes</span>
          </h1>

          <p ref="subtitle" class="mt-6 max-w-lg mx-auto text-xl text-gray-300">
            Join our community of food lovers. Find inspiration or contribute your own culinary creations.
          </p>

          <div ref="cta" class="mt-10 flex flex-col sm:flex-row justify-center gap-4">
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

    <!-- Recipe Tabs Section -->
    <section class="py-16 dark:bg-[#20161F]">
      <div class="container mx-auto">
        <!-- Tabs -->
        <div class="tabs-container items-center justify-center mx-auto">
          <div
            class="tabs-box flex flex-wrap md:flex-row space-x-6 p-1 bg-gray-400 dark:bg-[#422f40] items-center justify-center w-full rounded-full font-poppins-italic"
          >
            <!-- All Tab -->
            <input
              type="radio"
              name="my_tabs_1"
              id="all"
              class="hidden"
              v-model="selectedTab"
              value="all"
            />
            <label
              for="all"
              class="tab-label md:text-xl px-6 py-2 rounded-full dark:text-black bg-gray-400 dark:bg-[#422f40] text-gray-700 cursor-pointer transition-all duration-300 ease-in-out hover:bg-gray-200 hover:text-gray-700 hover:shadow-lg"
              :class="{
                'bg-white dark:bg-white dark:text-black text-gray-700 shadow-lg':
                  selectedTab === 'all',
                'bg-gray-400 text-gray-700 hover:bg-gray-200 hover:shadow-lg':
                  selectedTab !== 'all',
              }"
            >
              All
            </label>

            <!-- Breakfast Tab -->
            <input
              type="radio"
              name="my_tabs_1"
              id="breakfast"
              class="hidden"
              v-model="selectedTab"
              value="breakfast"
            />
            <label
              for="breakfast"
              class="tab-label md:text-xl px-6 py-2 dark:text-black rounded-full bg-gray-400 dark:bg-[#422f40] text-gray-700 cursor-pointer transition-all duration-300 ease-in-out hover:bg-gray-200 hover:text-gray-700 hover:shadow-lg"
              :class="{
                'bg-white dark:bg-white text-gray-700 shadow-lg':
                  selectedTab === 'breakfast',
                'bg-gray-400 text-gray-700 hover:bg-gray-200 hover:shadow-lg':
                  selectedTab !== 'breakfast',
              }"
            >
              Breakfast
            </label>

            <!-- Lunch Tab -->
            <input
              type="radio"
              name="my_tabs_1"
              id="lunch"
              class="hidden"
              v-model="selectedTab"
              value="lunch"
            />
            <label
              for="lunch"
              class="tab-label md:text-xl px-6 py-2 dark:text-black rounded-full dark:bg-[#422f40] bg-gray-400 text-gray-700 cursor-pointer transition-all duration-300 ease-in-out hover:bg-gray-200 hover:text-gray-700 hover:shadow-lg"
              :class="{
                'bg-white dark:bg-white text-gray-700 shadow-lg':
                  selectedTab === 'lunch',
                'bg-gray-400 text-gray-700 hover:bg-gray-200 hover:shadow-lg':
                  selectedTab !== 'lunch',
              }"
            >
              Lunch
            </label>

            <!-- Desserts & Sweets Tab -->
            <input
              type="radio"
              name="my_tabs_1"
              id="desserts"
              class="hidden"
              v-model="selectedTab"
              value="desserts"
            />
            <label
              for="desserts"
              class="tab-label md:text-xl px-6 py-2 dark:text-black rounded-full dark:bg-[#422f40] bg-gray-400 text-gray-700 cursor-pointer transition-all duration-300 ease-in-out hover:bg-gray-200 hover:text-gray-700 hover:shadow-lg"
              :class="{
                'bg-white dark:bg-white text-gray-700 shadow-lg':
                  selectedTab === 'desserts',
                'bg-gray-400 text-gray-700 hover:bg-gray-200 hover:shadow-lg':
                  selectedTab !== 'desserts',
              }"
            >
              Desserts & Sweets
            </label>

            <!-- Fasting Dishes Tab -->
            <input
              type="radio"
              name="my_tabs_1"
              id="fasting"
              class="hidden"
              v-model="selectedTab"
              value="fasting"
            />
            <label
              for="fasting"
              class="tab-label md:text-xl px-6 py-2 dark:text-black rounded-full dark:bg-[#422f40] bg-gray-400 text-gray-700 cursor-pointer transition-all duration-300 ease-in-out hover:bg-gray-200 hover:text-gray-700 hover:shadow-lg"
              :class="{
                'bg-white dark:bg-white text-gray-700 shadow-lg':
                  selectedTab === 'fasting',
                'bg-gray-400  text-gray-700 hover:bg-gray-200 hover:shadow-lg':
                  selectedTab !== 'fasting',
              }"
            >
              Fasting Dishes
            </label>

            <!-- Non-Fasting Dishes Tab -->
            <input
              type="radio"
              name="my_tabs_1"
              id="non-fasting"
              class="hidden"
              v-model="selectedTab"
              value="non-fasting"
            />
            <label
              for="non-fasting"
              class="tab-label md:text-xl px-6 py-2 dark:text-black rounded-full dark:bg-[#422f40] bg-gray-400 text-gray-700 cursor-pointer transition-all duration-300 ease-in-out hover:bg-gray-200 hover:text-gray-700 hover:shadow-lg"
              :class="{
                'bg-white dark:bg-white text-gray-700 shadow-lg':
                  selectedTab === 'non-fasting',
                'bg-gray-400 text-gray-700 hover:bg-gray-300 hover:shadow-lg':
                  selectedTab !== 'non-fasting',
              }"
            >
              Non-Fasting Dishes
            </label>

            <!-- Dinner Tab -->
            <input
              type="radio"
              name="my_tabs_1"
              id="dinner"
              class="hidden"
              v-model="selectedTab"
              value="dinner"
            />
            <label
              for="dinner"
              class="tab-label md:text-xl px-6 py-2 rounded-full dark:text-black dark:bg-[#422f40] bg-gray-400 text-gray-700 cursor-pointer transition-all duration-300 ease-in-out hover:bg-gray-200 hover:text-gray-700 hover:shadow-lg"
              :class="{
                'bg-white dark:bg-white text-gray-700 shadow-lg':
                  selectedTab === 'dinner',
                'bg-gray-400 text-gray-700 hover:bg-gray-300 hover:shadow-lg':
                  selectedTab !== 'dinner',
              }"
            >
              Dinner
            </label>
          </div>
        </div>
        <div class="mb-6 mt-6 flex justify-center dark:bg-[#20161F]">
          <input
            @input="handleSearch"
            v-model="searchQuery"
            type="text"
            placeholder="what are we cooking today?"
            class="w-full sm:w-1/2 p-3 border-2 border-green-300 rounded-full shadow-sm focus:outline-none focus:ring-2 focus:ring-green-500 dark:bg-[#20161F]"
          />
        </div>

        <!-- Dynamic component based on selected tab -->
        <div class="mt-8">
          <component :is="selectedComponent" />
        </div>
      </div>
    </section>

    <!-- Kushan Cuisine Section -->
    <section ref="feature1" class="max-w-7xl mx-auto px-4 py-16 sm:px-6 lg:px-8">
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
        <div ref="feature2" class="bg-white rounded-xl shadow-md overflow-hidden transition-transform duration-300 hover:scale-105">
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
        <div ref="feature3" class="bg-white rounded-xl shadow-md overflow-hidden transition-transform duration-300 hover:scale-105">
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

<style scoped>
@import url("https://fonts.googleapis.com/css2?family=Poppins:wght@400;500;600;700&display=swap");

.font-poppins {
  font-family: "Poppins", sans-serif;
}
.font-poppins-italic {
  font-family: "Poppins", sans-serif;
  font-style: italic;
}

@media (max-width: 768px) {
  /* Adjust for small screens */
  section[ref="heroSection"] {
    background-size: contain;
    background-repeat: no-repeat;
  }
  
  .tabs-box {
    flex-direction: column;
    align-items: stretch;
    gap: 0.5rem;
  }
  
  .tab-label {
    width: 100%;
    text-align: center;
  }
}

/* Hide scrollbar */
::-webkit-scrollbar {
  display: none;
}
</style>