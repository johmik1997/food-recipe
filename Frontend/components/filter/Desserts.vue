<script setup>
import { onMounted, ref } from "vue";
import { useToast } from "vue-toast-notification";
const toast = useToast();

const categories = ref([]);
const recipesStore = useRecipeStore();
const bookmarkStore = useBookmarkStore();

const carouselRef = ref(null);
const userStore = authStore();
const user_id = userStore.$state.userId;

const bookmarkStates = reactive({});
const useAuthStore = authStore();
const bookMarkId = computed(() => bookmarkStore.bookmarkedId);

const handleCheckBookmark = async (recipeId) => {
  try {
    const payload = {
      recipe_id: recipeId,
      user_id,
    };
    await bookmarkStore.checkIfBookmarked(payload);
    bookmarkStates[recipeId] = bookmarkStore.$state.isBookmarked; 
  } catch (error) {
    console.error("Error checking bookmark status:", error);
  }
};

const handlesaveBookmark = async (recipeId) => {
  try {
    const payload = {
      recipe_id: recipeId,
      user_id,
    };
    await bookmarkStore.createBookmark(payload);
    toast.success("Recipe saved successfully!");
    bookmarkStates[recipeId] = true;
  } catch (error) {
    console.log("error saving a recipe", error);
    toast.error("error saving recipe!");
  }
};

const handleRemoveBookmark = async (recipeId) => {
  try {
    if (!bookMarkId.value) {
      toast.error("No bookmark ID found for deletion");
      return;
    }

    await bookmarkStore.removeBookmark(bookMarkId.value);
    toast.success("Bookmark removed successfully!");
    bookmarkStates[recipeId] = false;
  } catch (error) {
    console.error("Error removing the bookmark:", error);
    toast.error("Error removing the bookmark");
  }
};

onMounted(async () => {
  try {
    await recipesStore.getAllRecipes();
    recipesStore.recipes.forEach((recipe) => {
      bookmarkStates[recipe.id] = false;
      handleCheckBookmark(recipe.id);
    });
    await recipesStore.filterByCategory("Dessert");
  } catch (error) {
    console.error("Failed to load or filter categories", error);
  }
});

const scrollCarousel = (direction) => {
  const container = carouselRef.value;
  const scrollAmount = 300;
  if (container) {
    if (direction === "left") {
      container.scrollBy({ left: -scrollAmount, behavior: "smooth" });
    } else {
      container.scrollBy({ left: scrollAmount, behavior: "smooth" });
    }
  }
};
</script>

<template>
  <div>
    <div v-if="recipesStore.recipes.length === 0">
      <div class="flex flex-col items-center gap-4 p-6 bg-white">
        <h1 class="text-center font-poppins-italic text-2xl md:text-4xl font-bold text-green-800">
          No Desserts & Sweets recipes found, be the first to create one!
        </h1>
        <NuxtLink to="/recipes/create" class="relative group cursor-pointer">
          <img
            src="/Chef-pana.svg"
            alt="Create Recipe"
            class="h-[300px] w-[300px] transition-transform duration-300 group-hover:scale-105"
          />
          <div
            class="absolute inset-0 bg-pink-500 bg-opacity-70 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity duration-300"
          >
            <span class="text-white font-bold text-lg">Create Recipe</span>
          </div>
        </NuxtLink>
      </div>
    </div>

    <div v-else class="relative font-poppins bg-white">
      <!-- Left Chevron Button -->
      <button
        @click="scrollCarousel('left')"
        class="absolute left-0 top-1/2 transform -translate-y-1/2 bg-pink-500 text-white p-3 rounded-full shadow-lg hover:bg-pink-600 transition-colors duration-200 z-10"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="h-6 w-6"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M15 19l-7-7 7-7"
          />
        </svg>
      </button>

      <!-- Carousel Container -->
      <div
        ref="carouselRef"
        class="flex overflow-x-auto scroll-smooth snap-x snap-mandatory gap-6 p-4 transition-transform duration-1000 ease-in-out bg-white"
        style="scrollbar-width: none; -ms-overflow-style: none"
      >
        <div
          v-for="recipe in recipesStore.recipes"
          :key="recipe.id"
          class="flex flex-col bg-white w-96 border-2 border-pink-100 rounded-xl overflow-hidden flex-shrink-0 snap-center transition-all duration-300 hover:shadow-lg hover:-translate-y-1"
        >
          <figure class="relative h-60 overflow-hidden">
            <img
              :src="recipe.featured_image || 'https://img.daisyui.com/images/stock/photo-1606107557195-0e29a4b5b4aa.webp'"
              alt="Recipe Image"
              class="w-full h-full object-cover transition-transform duration-300 hover:scale-110"
            />
            <!-- Save Icon -->
            <button
              v-if="useAuthStore.$state.userId"
              @click.stop="bookmarkStates[recipe.id] ? handleRemoveBookmark(recipe.id) : handlesaveBookmark(recipe.id)"
              class="absolute bottom-4 right-4 p-2 rounded-full bg-white/90 hover:bg-white transition-colors duration-200 shadow-md"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="h-6 w-6"
                :class="{
                  'text-pink-500 fill-pink-500': bookmarkStates[recipe.id],
                  'text-pink-300': !bookmarkStates[recipe.id],
                }"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z"
                />
              </svg>
            </button>
          </figure>

          <div class="p-6 flex flex-col flex-1">
            <NuxtLink :to="`/recipes/${recipe.id}`" class="group">
              <h2 class="text-xl font-bold text-pink-800 mb-2 group-hover:text-pink-600 transition-colors">
                {{ recipe.title || "Untitled Recipe" }}
              </h2>
            </NuxtLink>

            <div class="flex justify-between items-center mt-2">
              <div class="flex items-center">
                <div class="flex mr-2">
                  <span
                    v-for="star in 5"
                    :key="star"
                    class="text-xl"
                    :class="{
                      'text-yellow-500': star <= Math.round(recipe.average_rating),
                      'text-pink-100': star > Math.round(recipe.average_rating),
                    }"
                  >
                    ★
                  </span>
                </div>
                <span class="text-pink-700 font-medium text-sm">
                  ({{ recipe.ratings_aggregate.aggregate.count }})
                </span>
              </div>
              <span class="font-bold text-pink-600 text-xl">
                ${{ recipe.price }}
              </span>
            </div>
            
            <!-- Category and Prep Time -->
            <div class="flex justify-between items-center mt-4">
              <span class="text-pink-700 bg-pink-100 px-3 py-1 rounded-full text-sm font-medium">
                {{ recipe.catagory?.name || "Uncategorized" }}
              </span>
              <span class="text-pink-700 text-sm font-medium">
                {{ recipe.prep_time || "0" }} mins
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- Right Chevron Button -->
      <button
        @click="scrollCarousel('right')"
        class="absolute right-0 top-1/2 transform -translate-y-1/2 bg-pink-500 text-white p-3 rounded-full shadow-lg hover:bg-pink-600 transition-colors duration-200 z-10"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="h-6 w-6"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M9 5l7 7-7 7"
          />
        </svg>
      </button>
    </div>
  </div>
</template>

<style>
@import url("https://fonts.googleapis.com/css2?family=Poppins:wght@400;500;600;700&display=swap");

.font-poppins {
  font-family: "Poppins", sans-serif;
}
.font-poppins-italic {
  font-family: "Poppins", sans-serif;
  font-style: italic;
}

/* Hide scrollbar */
::-webkit-scrollbar {
  display: none;
}

/* Smooth card hover effect */
.card-hover {
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}
.card-hover:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 25px rgba(236, 72, 153, 0.1);
}
</style>