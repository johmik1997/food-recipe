<script setup>
import { useToast } from "vue-toast-notification";
import { onMounted, ref, computed, reactive } from "vue";

import { useRoute } from "vue-router";
const toast = useToast();
const route = useRoute();
const recipesStore = useRecipeStore();

const searchQuery = ref("");
const userStore = authStore();
const bookmarkStore = useBookmarkStore();
const likeStore = useLikeStore();
const user_id = userStore.$state.userId;
  console.log("✅ userId after mounted from the index page", user_id);
const bookmarkStates = reactive({});
const likeStates = reactive({});

const bookMarkId = computed(() => bookmarkStore.bookmarkedId);
const likedId = computed(() => likeStore.likeId);

const handleSearch = async () => {
  recipesStore.setSearchRecipe(searchQuery.value);
  console.log("search item from the comp", searchQuery.value);
  await recipesStore.getAllRecipes(searchQuery.value);
};

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
      likeStates[recipe.id] = false;
      handleCheckLIke(recipe.id);
    });
  } catch (error) {
    toast.error("Failed to load recipes");
    console.error("Failed to load recipes", error);
  }
});
</script>

<template>
  <div class="font-poppins p-6 lg:container lg:mx-auto bg-white">
    <!-- Search Input -->
    <div class="mb-6 flex justify-center">
      <input
        @input="handleSearch"
        v-model="searchQuery"
        type="text"
        placeholder="What are we cooking today?"
        class="w-full sm:w-1/2 p-3 border-2 border-green-300 rounded-full shadow-sm focus:outline-none focus:ring-2 focus:ring-green-500 text-green-800 placeholder-green-400 bg-white"
      />
    </div>

    <!-- Recipe Grid -->
    <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <!-- Recipe Card -->
      <div
        v-for="recipe in recipesStore.recipes"
        :key="recipe.id"
        class="card bg-white border-2 border-green-100 rounded-xl overflow-hidden shadow-md hover:shadow-lg transition-all duration-300"
      >
        <!-- Recipe Image with Action Buttons -->
        <figure class="relative">
          <img
            :src="recipe.featured_image || '/images/recipe-placeholder.jpg'"
            alt="Recipe Image"
            class="w-full h-60 object-cover transition-transform duration-300 hover:scale-105"
          />
          
          <!-- Like Button -->
          <button
            @click.stop="likeStates[recipe.id] ? handleRemoveLike(recipe.id) : handleLikeRecipe(recipe.id)"
            class="absolute bottom-2 left-2 p-2 rounded-full bg-white/90 hover:bg-white transition-all duration-200 shadow-md"
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
          
          <!-- Bookmark Button -->
          <button
            @click.stop="bookmarkStates[recipe.id] ? handleRemoveBookmark(recipe.id) : handlesaveBookmark(recipe.id)"
            class="absolute bottom-2 right-2 p-2 rounded-full bg-white/90 hover:bg-white transition-all duration-200 shadow-md"
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

        <!-- Recipe Details -->
        <div class="card-body p-6">
          <!-- Recipe Title -->
          <NuxtLink :to="{ name: 'recipes-id', params: { id: recipe.id } }">
            <h2 class="card-title text-xl font-bold text-green-800 mb-2 hover:text-green-600 transition-colors">
              {{ recipe.title || "Untitled Recipe" }}
            </h2>
          </NuxtLink>

          <!-- Rating and Price -->
          <div class="flex justify-between items-center mb-3">
            <!-- Rating Stars -->
            <div class="flex items-center">
              <div class="flex mr-1">
                <span
                  v-for="star in 5"
                  :key="star"
                  class="text-lg"
                  :class="{
                    'text-yellow-400': star <= Math.round(recipe.average_rating),
                    'text-gray-300': star > Math.round(recipe.average_rating),
                  }"
                >
                  ★
                </span>
              </div>
              <span class="text-sm text-green-700 font-medium">
                ({{ recipe.ratings_aggregate.aggregate.count }})
              </span>
            </div>
            
            <!-- Price -->
            <span class="text-xl font-bold text-green-700">
              ${{ recipe.price }}
            </span>
          </div>

          <!-- Category and Prep Time -->
          <div class="flex justify-between text-sm">
            <div class="flex items-center text-green-700">
              <span class="mr-1">🥘</span>
              {{ recipe.catagory?.name || "Uncategorized" }}
            </div>
            <div class="flex items-center text-green-700">
              <span class="mr-1">⏱️</span>
              {{ recipe.prep_time || "N/A" }} mins
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style>
@import url('https://fonts.googleapis.com/css2?family=Poppins:wght@400;500;600;700&display=swap');

.font-poppins {
  font-family: 'Poppins', sans-serif;
}

.card {
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 20px rgba(5, 150, 105, 0.1);
}

/* Smooth transitions for interactive elements */
button, a {
  transition: all 0.2s ease;
}

/* Rating stars styling */
.text-yellow-400 {
  color: #facc15;
}

/* Input placeholder styling */
::placeholder {
  color: #86efac;
  opacity: 1;
}
</style>