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
    bookmarkStates[recipeId] = false; // Update bookmark state for this recipe
  } catch (error) {
    console.error("Error removing the bookmark:", error);
    toast.error("Error removing the bookmark");
  }
};


onMounted(async () => {
  try {
    // Fetch categories
    // await recipesStore.getCategories();
    // categories.value = recipesStore.categories;
    // console.log("all categories", JSON.stringify(categories.value, null, 2));
    await recipesStore.getAllRecipes();
    recipesStore.recipes.forEach((recipe) => {
      bookmarkStates[recipe.id] = false;
      handleCheckBookmark(recipe.id);
    });
    await recipesStore.filterByCategory("Dinner");
    console.log(
      "filtered recipes",
      JSON.stringify(recipesStore.recipes, null, 2)
    );
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
      <div class="flex flex-col items-center gap-4 p-6">
        <h1
          class="text-center dark:text-gray-100 font-poppins-italic text-2xl md:text-4xl font-bold text-gray-800"
        >
          No Dinner recipes found, be the first to create one!
        </h1>
        <NuxtLink to="/recipes/create" class="relative group cursor-pointer">
          <img
            src="/Chef-pana.svg"
            alt="Create Recipe"
            class="h-[300px] w-[300px] transition-transform duration-300 group-hover:scale-105"
          />
          <div
            class="absolute inset-0 bg-black bg-opacity-50 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity duration-300"
          >
            <span class="text-white font-bold text-lg">Create Recipe</span>
          </div>
        </NuxtLink>
      </div>
    </div>
    <div v-else class="relative font-poppins dark:bg-[#20161F]">
      <!-- Left Chevron Button -->
      <button
        @click="scrollCarousel('left')"
        class="absolute left-0 top-1/2 transform -translate-y-1/2 bg-green-500 text-white p-3 rounded-full shadow-lg hover:bg-green-600 transition-colors duration-200 z-10"
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
        class="flex overflow-x-auto scroll-smooth snap-x snap-mandatory gap-6 p-4 transition-transform duration-1000 ease-in-out dark:bg-[#20161F]"
        style="scrollbar-width: none; -ms-overflow-style: none"
      >
        <div
          v-for="recipe in recipesStore.recipes"
          :key="recipe.id"
          class="card bg-base-100 dark:bg-[#20161F] w-96 border-2 transition-shadow duration-300 ease-in-out transform hover:-translate-y-1 rounded-xl overflow-hidden flex-shrink-0 snap-center"
        >
          <figure class="relative">
            <img
              :src="
                recipe.featured_image ||
                'https://img.daisyui.com/images/stock/photo-1606107557195-0e29a4b5b4aa.webp'
              "
              alt="Recipe Image"
              class="w-full h-60 object-cover transition-transform duration-300 hover:scale-105"
            />
            <!-- Save Icon -->
            <button
              v-if="useAuthStore.$state.userId"
              @click.stop="
                bookmarkStates[recipe.id]
                  ? handleRemoveBookmark(recipe.id)
                  : handlesaveBookmark(recipe.id)
              "
              class="absolute bottom-2 right-2 p-2 rounded-full bg-white/80 hover:bg-white/90 transition-colors duration-200 backdrop-blur-sm"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="h-6 w-6"
                :class="{
                  'text-green-500': bookmarkStates[recipe.id],
                  'text-green-300': !bookmarkStates[recipe.id],
                }"
                fill="none"
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

          <div class="card-body p-6">
            <NuxtLink :to="`/recipes/${recipe.id}`">
              <div
                class="flex flex-row justify-between items-center dark:text-[#C9CF43]"
              >
                <h2
                  class="card-title text-lg font-bold text-green-900 dark:text-[#C9CF43]"
                >
                  {{ recipe.title || "Untitled Recipe" }}
                </h2>
                <!-- <h1 class="font-bold text-2xl text-green-700 dark:text-[#C9CF43]">
                    $ {{ recipe.price }}
                  </h1> -->
              </div>
            </NuxtLink>

            <div class="flex flex-row justify-between">
              <div class="flex items-center mt-2">
                <span
                  v-for="star in 5"
                  :key="star"
                  class="text-2xl"
                  :class="{
                    'text-green-500 dark:text-[#C9CF43]':
                      star <= Math.round(recipe.average_rating),
                    'text-green-200': star > Math.round(recipe.average_rating),
                  }"
                >
                  ★
                </span>
                <span class="font-bold"
                  >({{ recipe.ratings_aggregate.aggregate.count }})</span
                >
              </div>
              <!-- <button class="py-1 px-4 r dark:bg-[#C9CF43]">
                  <img
                    src="/icons/cart-large-svgrepo-com.svg"
                    alt=""
                    class="w-12 h-12"
                  />
                </button> -->
              <h1 class="font-bold text-2xl text-green-700 dark:text-[#C9CF43]">
                $ {{ recipe.price }}
              </h1>
            </div>
            <!-- Category Badge and Preparation Time -->
            <div class="flex flex-row justify-between items-center mr-4">
              <div class="flex flex-row gap-8 mt-4 justify-between">
                <div class="text-sm font-semibold dark:text-gray-300">
                  <span></span>
                 {{ recipe.catagory?.name || "Uncategorized" }}
                </div>
                <div class="text-sm font-semibold dark:text-gray-300">
                  <h1>
                   {{ recipe.prep_time || "Uncategorized" }} mins
                  </h1>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Right Chevron Button -->
      <button
        @click="scrollCarousel('right')"
        class="absolute right-0 top-1/2 transform -translate-y-1/2 bg-green-500 text-white p-3 rounded-full shadow-lg hover:bg-green-600 transition-colors duration-200 z-10"
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
</style>
