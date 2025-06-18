<script setup>
import { ref, onMounted, computed, reactive } from "vue";
import { useToast } from "vue-toast-notification";

const toast = useToast();
const recipesStore = useRecipeStore();
const useAuthStore = authStore();
const bookmarkStore = useBookmarkStore();

const carouselRef = ref(null);

const userStore = authStore();
const likeStore = useLikeStore();
const user_id = userStore.$state.userId;

// Reactive object to track bookmark state for each recipe
const bookmarkStates = reactive({});
const likeStates = reactive({});

const bookMarkId = computed(() => bookmarkStore.bookmarkedId);
const likedId = computed(() => likeStore.likeId);

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

const handleCheckLIke = async (recipeId) => {
  try {
    const payload = {
      recipe_id: recipeId,
      user_id,
    };
    await likeStore.checkIfLiked(payload);
    likeStates[recipeId] = likeStore.$state.isLiked;
  } catch (error) {
    console.error("Error checking bookmark status:", error);
  }
};

const handleLikeRecipe = async (recipeId) => {
  try {
    const payload = {
      recipe_id: recipeId,
      user_id,
    };
    await likeStore.likeRecipe(payload);
    toast.success("recipe liked successfully!");
    likeStates[recipeId] = true;
  } catch (error) {
    console.log("error liking the recipe", error);
  }
};

const handleRemoveLike = async (recipeId) => {
  try {
    await likeStore.removeLike(likedId.value);
    toast.success("Recipe unliked successfully!");
    likeStates[recipeId] = false;
  } catch (error) {
    console.log("error liking the recipe", error);
  }
};

onMounted(async () => {
  console.log("Fetching all recipes...");
  try {
    await recipesStore.getAllRecipes();
    console.log(
      "Recipes loaded successfully",
      JSON.stringify(recipesStore.recipes, null, 2)
    );

    recipesStore.recipes.forEach((recipe) => {
      bookmarkStates[recipe.id] = false;
      handleCheckBookmark(recipe.id);

      // Initialize like state
      likeStates[recipe.id] = false;
      handleCheckLIke(recipe.id);
    });
  } catch (error) {
    toast.error("Failed to load recipes");
    console.error("Failed to load recipes", error);
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
  <div class="relative font-poppins bg-white">
    <!-- Left Chevron Button -->
    <button
      @click="scrollCarousel('left')"
      class="absolute left-0 top-1/2 transform -translate-y-1/2 bg-green-600 text-white p-3 rounded-full shadow-lg hover:bg-green-700 transition-colors duration-200 z-10"
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
        class="card border border-green-100 bg-white w-[360px] transition-shadow duration-300 ease-in-out transform hover:-translate-y-1 rounded-xl overflow-hidden flex-shrink-0 snap-center shadow-sm hover:shadow-md"
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
          <button
            v-if="useAuthStore.$state.userId"
            @click.stop="
              likeStates[recipe.id]
                ? handleRemoveLike(recipe.id)
                : handleLikeRecipe(recipe.id)
            "
            class="absolute bottom-2 left-2 p-2 rounded-full bg-white/80 hover:bg-white transition-colors duration-200 backdrop-blur-sm shadow-sm"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="h-6 w-6"
              :class="{
                'text-red-500 fill-red-500': likeStates[recipe.id],
                'text-gray-400': !likeStates[recipe.id],
              }"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"
              />
            </svg>
          </button>
          <!-- Save Icon -->
          <button
            v-if="useAuthStore.$state.userId"
            @click.stop="
              bookmarkStates[recipe.id]
                ? handleRemoveBookmark(recipe.id)
                : handlesaveBookmark(recipe.id)
            "
            class="absolute bottom-2 right-2 p-2 rounded-full bg-white/80 hover:bg-white transition-colors duration-200 backdrop-blur-sm shadow-sm"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="h-6 w-6"
              :class="{
                'text-green-600 fill-green-600': bookmarkStates[recipe.id],
                'text-gray-400': !bookmarkStates[recipe.id],
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

        <div class="card-body p-6">
          <NuxtLink :to="`/recipes/${recipe.id}`">
            <div class="flex flex-row justify-between items-center">
              <h2 class="card-title text-xl font-bold text-green-800">
                {{ recipe.title || "Untitled Recipe" }}
              </h2>
            </div>
          </NuxtLink>

          <div class="flex flex-row justify-between">
            <div class="flex items-center gap-2 mt-2">
              <!-- Star Icon & Rating -->
              <span class="flex items-center text-lg font-semibold">
                <span class="mr-1 text-2xl text-yellow-400"> ⭐ </span>
                <span class="text-gray-700">
                  {{ recipe.average_rating || 0 }} / 5
                </span>
              </span>

              <!-- Ratings Count -->
              <span
                class="bg-green-100 text-green-800 text-xs font-semibold px-3 py-1 rounded-full"
              >
                {{ recipe.ratings_aggregate.aggregate.count }} Reviews
              </span>
            </div>

            <h1 class="font-semibold text-xl text-green-600">
              $ {{ recipe.price }}
            </h1>
          </div>
          
          <!-- Category Badge and Preparation Time -->
          <div class="flex flex-row justify-between items-center mr-4">
            <div class="flex flex-row gap-8 mt-4 justify-between">
              <div class="text-sm font-semibold text-green-700">
                <span></span>
                {{ recipe.catagory?.name || "Uncategorized" }}
              </div>
              <div class="text-sm font-semibold text-green-700">
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
      class="absolute right-0 top-1/2 transform -translate-y-1/2 bg-green-600 text-white p-3 rounded-full shadow-lg hover:bg-green-700 transition-colors duration-200 z-10"
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
</template>

<style>
@import url("https://fonts.googleapis.com/css2?family=Poppins:wght@400;500;600;700&display=swap");

.font-poppins {
  font-family: "Poppins", sans-serif;
}

/* Hide scrollbar */
::-webkit-scrollbar {
  display: none;
}

/* Card hover effect */
.card:hover {
  box-shadow: 0 10px 15px -3px rgba(5, 150, 105, 0.1), 0 4px 6px -2px rgba(5, 150, 105, 0.05);
}

/* Transition for buttons */
button {
  transition: all 0.2s ease-in-out;
}

/* Rating star color */
.text-yellow-400 {
  color: #facc15;
}
</style>