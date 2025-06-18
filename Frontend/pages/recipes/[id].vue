<script setup>
import { ref, onMounted, computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useToast } from "vue-toast-notification";
import { authStore } from "../stores/auth";

useSeoMeta({
  title: "recipe-app | Recipe Details Page",
  description: "The project app meta.",
});

const toast = useToast();
const route = useRoute();
const recipeId = parseInt(route.params.id);
const router = useRouter();
const recipeStore = useRecipeStore();
const likeStore = useLikeStore();
const auth = authStore();
const commentStore = useCommnentStore();
const ratingStore = useRatingStore();
const bookmarkStore = useBookmarkStore();

// Ensure user data is accessed properly
const userId = computed(() => auth.userId || 0);
const userName = computed(() => auth.userName || '');

const recipeDetails = ref(null);
const isLoading = ref(false);
const comment = ref("");
const isBookmarked = ref(false);
const isLiked = ref(false);
const rating = ref(0);
const modalRef = ref(null);
const commentSubmitted = ref(false);
const rateSubmitted = ref(false);
const bookMarkId = computed(() => bookmarkStore.bookmarkedId);
const likedId = computed(() => likeStore.likeId);

const fetchRecipeDetails = async () => {
  isLoading.value = true;
  try {
    console.log("Fetching recipe details for ID:", recipeId);
    await recipeStore.singleRecipe(recipeId);
    recipeDetails.value = recipeStore.recipe;
  } catch (error) {
    toast.error("Failed to load recipe detail");
  } finally {
    isLoading.value = false;
  }
};

const handleCheckBookmark = async () => {
  try {
    const payload = {
      recipe_id: recipeId,
      user_id: userId.value,
    };
    await bookmarkStore.checkIfBookmarked(payload);
    isBookmarked.value = bookmarkStore.$state.isBookmarked;
  } catch (error) {
    console.error("Error checking bookmark status:", error);
  }
};

const handlesaveBookmark = async () => {
  try {
    const payload = {
      recipe_id: recipeId,
      user_id: userId.value,
    };
    await bookmarkStore.createBookmark(payload);
    isBookmarked.value = true;
    toast.success("recipe saved successfully!");
  } catch (error) {
    console.log("error saving a recipe", error);
    toast.error("error saving recipe!");
  }
};

const handleRemoveBookmark = async () => {
  try {
    if (!bookMarkId.value) {
      toast.error("No bookmark ID found for deletion");
      return;
    }

    await bookmarkStore.removeBookmark(bookMarkId.value);
    toast.success("Bookmark removed successfully!");
    isBookmarked.value = false;
  } catch (error) {
    console.error("Error removing the bookmark:", error);
    toast.error("Error removing the bookmark");
  }
};

const handleComments = async () => {
  try {
    const payload = {
      comment: comment.value,
      recipe_id: recipeId,
      user_id: userId.value,
    };

    const res = await commentStore.createComment(payload);
    if (res) {
      commentSubmitted.value = true;
      comment.value = "";
      toast.success("Comment submitted successfully!");
    }
    fetchRecipeDetails();
  } catch (error) {
    console.error("Error submitting comment:", error);
  }
};

const handleCheckLike = async () => {
  try {
    const payload = {
      recipe_id: recipeId,
      user_id: userId.value,
    };
    await likeStore.checkIfLiked(payload);
    isLiked.value = likeStore.$state.isLiked;
  } catch (error) {
    console.log("error checking if recipe is liked");
  }
};

const handleLikeRecipe = async () => {
  try {
    const payload = {
      recipe_id: recipeId,
      user_id: userId.value,
    };
    const res = await likeStore.likeRecipe(payload);
    if (res) {
      isLiked.value = true;
      toast.success("recipe liked successfully!");
    }
    fetchRecipeDetails();
    handleCheckLike();
  } catch (error) {
    toast.error("error liking the recipe");
  }
};

const handleRemoveLike = async () => {
  try {
    if (!likedId.value) {
      console.log("no like id found to delete");
      return;
    }

    await likeStore.removeLike(likedId.value);
    isLiked.value = false;
    fetchRecipeDetails();
  } catch (error) {
    toast.error("Error removing the like");
  }
};

const handleRating = async () => {
  try {
    const payload = {
      rating: rating.value,
      recipe_id: recipeId,
      user_id: userId.value,
    };

    await ratingStore.createRating(payload);
    rateSubmitted.value = true;
    toast.success("thanks for rating");
    fetchRecipeDetails();
    closeModal();
  } catch (error) {
    console.error("Error rating the recipe:", error);
  }
};

const handleBuyRecipe = async () => {
  if (isNaN(recipeId)) {
    console.error("Invalid recipe ID");
    return;
  }

  try {
    isLoading.value = true;
    const res = await recipeStore.buyRecipe({
      buyer_id: userId.value,
      recipe_id: recipeId,
    });
    const paymentId = recipeStore.$state.paymentId;
    await recipeStore.getCheckOutUrl(paymentId);
    const checkout_url = recipeStore.$state.checkoutUrl;
    window.location.href = checkout_url;
  } catch (error) {
    console.error("Error buying recipe:", error);
    toast.error("Something went wrong. Please try again!");
  } finally {
    isLoading.value = false;
  }
};

const handleGetPaymentDetails = async () => {
  try {
    const res = await recipeStore.getPaymentDetails({
      recipe_id: recipeId,
      buyer_id: userId.value,
    });
    
    if (!res?.data?.sold_recipes?.[0]?.payments?.[0]) {
      console.log("No payment details found");
      recipeStore.$state.paymentStatus = "";
      return;
    }

    const paymentDetails = res.data.sold_recipes[0].payments[0];
    recipeStore.paymentDetails = paymentDetails;
    recipeStore.tx_ref = paymentDetails.tx_ref;
    recipeStore.paymentId = paymentDetails.id;
    recipeStore.paymentStatus = paymentDetails.payment_status;
  } catch (error) {
    console.error("Error getting payment details:", error);
  }
};

const handleVerifyPayment = async () => {
  try {
    const tx_ref = recipeStore?.paymentDetails?.tx_ref;
    const paymentId = recipeStore?.paymentDetails?.id;
    
    if (!paymentId || !tx_ref) {
      console.error("Missing payment details for verification");
      return;
    }
    
    await recipeStore.verifyPayment({ id: paymentId, tx_ref });
  } catch (error) {
    console.error("Error verifying payment:", error);
  }
};

// Delete recipe functionality
const recipeToDeleteId = ref(null);
const deleteModalRef = ref(null);

const openDeleteModal = (recipeId) => {
  recipeToDeleteId.value = recipeId;
  deleteModalRef.value?.showModal();
};

const closeDeleteModal = () => {
  deleteModalRef.value?.close();
};

const confirmDelete = async () => {
  if (recipeToDeleteId.value) {
    try {
      await handleDeleteRecipe(recipeToDeleteId.value);
      toast.success("Recipe deleted successfully!");
      router.push("/recipes");
    } catch (error) {
      toast.error("Error deleting recipe. Please try again");
    }
  }
};

const handleDeleteRecipe = async (recipe_id) => {
  try {
    await recipeStore.deleteRecipes({ id: recipe_id, user_id: userId.value });
    await auth.getProfile();
    await recipeStore.getSellerRecipes();
    router.push("/recipes");
  } catch (error) {
    toast.error("Something went wrong! Please try again.");
  }
};

onMounted(() => {
  fetchRecipeDetails();
  handleCheckBookmark();
  handleCheckLike();
  handleGetPaymentDetails();
});
const openModal = () => {
  if (modalRef.value) modalRef.value.showModal();
};

const closeModal = () => {
  if (modalRef.value) modalRef.value.close();
};

</script>


<template>
  <div>
    <!-- Delete Confirmation Modal -->
    <dialog ref="deleteModalRef" class="modal">
      <div class="modal-box bg-white">
        <!-- Close button -->
        <form method="dialog">
          <button
            class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2 text-green-600"
            @click="closeDeleteModal"
          >
            ✕
          </button>
        </form>

        <!-- Modal title -->
        <h3 class="text-lg font-bold text-green-800">Confirm Deletion</h3>

        <!-- Modal content -->
        <p class="py-4 text-green-700">Are you sure you want to delete this recipe?</p>

        <!-- Modal action buttons -->
        <div class="flex justify-end gap-4">
          <button
            @click="closeDeleteModal"
            class="btn bg-green-100 text-green-800 hover:bg-green-200"
          >
            Cancel
          </button>
          <button
            @click="confirmDelete"
            class="btn bg-red-100 text-red-800 hover:bg-red-200"
          >
            Delete
          </button>
        </div>
      </div>
    </dialog>

    <div v-if="isLoading" class="flex justify-center items-center h-screen">
      <div
        class="animate-spin rounded-full h-[80px] w-[80px] border-t-2 border-b-2 border-green-500"
      ></div>
    </div>

    <div v-else-if="recipeDetails" class="font-poppins container mx-auto bg-white min-h-screen">
      <!-- Recipe Header -->
      <div class="bg-white overflow-hidden pt-6">
        <div class="flex flex-col md:flex-row gap-2 md:gap-12">
          <!-- left side -->
          <div class="flex-[3.5]">
            <div class="mb-6">
              <img
                :src="recipeDetails.featured_image || '/images/recipe-placeholder.jpg'"
                :alt="recipeDetails.title"
                class="w-full h-full object-cover rounded-lg"
              />
            </div>
          </div>
          
          <!-- Middle Section: Recipe Title and Details -->
          <div class="flex-[5] flex flex-col items-center justify-center">
            <div class="flex flex-col items-center">
              <h1 class="text-3xl md:text-4xl lg:text-5xl font-bold text-green-800 mb-4 pl-10 font-poppins-italic">
                {{ recipeDetails.title }}
              </h1>
              <div>
                <p class="text-green-700 font-bold text-xl mt-2">
                  by {{ recipeDetails.user.username }}
                </p>
              </div>
            </div>

            <div class="flex flex-col items-center space-x-6 mb-6 w-full">
              <div class="flex flex-row justify-between items-center space-x-4 p-4 mt-2">
                <!-- Category Section -->
                <div class="flex flex-col items-center flex-2">
                  <span class="text-green-800 font-semibold">
                    <img src="/icons/fork-and-spoon-meal-2-svgrepo-com.svg" class="w-8 h-8" alt=""/>
                  </span>
                  <span class="text-green-700 w-full">
                    {{ recipeDetails.catagory.name }}
                  </span>
                </div>

                <!-- Divider Line -->
                <div class="h-10 w-1 bg-green-300 mx-2"></div>

                <!-- Price Section -->
                <div class="flex flex-[1] flex-col items-center">
                  <span class="text-green-800 font-semibold">
                    <img src="/icons/money-bag-svgrepo-com.svg" class="w-8 h-8" alt=""/>
                  </span>
                  <span class="text-green-700">${{ recipeDetails.price }}</span>
                </div>

                <!-- Divider Line -->
                <div class="h-10 w-1 bg-green-300 mx-2"></div>

                <!-- Cook Time Section -->
                <div class="flex flex-col items-center flex-2">
                  <span class="text-green-800 font-semibold">
                    <img src="/icons/clock-lines-svgrepo-com.svg" class="w-8 h-8" alt=""/>
                  </span>
                  <span class="text-green-700">
                    {{ recipeDetails.prep_time }} min
                  </span>
                </div>
              </div>

              <div class="flex items-center space-x-2 mt-2">
                <div class="flex flex-row gap-1">
                  <span
                    v-for="star in 5"
                    :key="star"
                    class="text-2xl"
                    :class="{
                      'text-yellow-500': star <= Math.round(recipeDetails.rating.rating),
                      'text-yellow-100': star > Math.round(recipeDetails.rating.rating),
                    }"
                  >
                    ★
                  </span>
                  <span class="font-semibold mt-1 text-green-700">
                    ({{ recipeDetails.ratings_aggregate.aggregate.count }})
                  </span>
                </div>
              </div>

              <div class="flex flex-row pt-4 md:pt-10 gap-2 md:gap-4 items-center">
                <!-- Owner Actions -->
                <div
                  class="flex flex-row gap-4"
                  v-if="Number(recipeDetails.user.id) === Number(auth.$state.userId)"
                >
                  <NuxtLink
                    :to="`/recipes/update/${recipeDetails.id}`"
                    class="btn bg-green-100 ring-1 ring-green-500 py-2 px-4 rounded-full hover:bg-green-300"
                  >
                    <img src="/icons/edit-user-2-svgrepo-com.svg" alt="update icon" class="w-8 h-8"/>
                    <span class="text-xs text-green-700 hidden md:block">update recipe</span>
                  </NuxtLink>
                  <button
                    @click="openDeleteModal(recipeDetails.id)"
                    class="btn bg-green-100 ring-1 ring-green-500 py-2 px-4 rounded-full hover:bg-green-300"
                  >
                    <img src="/icons/delete-svgrepo-com.svg" class="h-8 w-8" />
                    <span class="text-xs text-green-700 hidden md:block">delete recipe</span>
                  </button>
                </div>

                <!-- Visitor Actions -->
                <div v-else class="flex flex-row md:pt-4 gap-2 md:gap-4 lg:gap-5 items-center">
                  <!-- Rating Button -->
                  <div class="">
                    <button class="flex flex-col items-center justify-center" @click="openModal">
                      <img src="/icons/rating-svgrepo-com.svg" class="h-8 w-8" alt=""/>
                      <span class="text-xs text-yellow-500 hidden md:block">rate</span>
                      <div class="rating">
  <input type="radio" name="rating-2" value="1" class="mask mask-star-2 bg-yellow-400" v-model="rating" />
  <input type="radio" name="rating-2" value="2" class="mask mask-star-2 bg-yellow-400" v-model="rating" />
  <input type="radio" name="rating-2" value="3" class="mask mask-star-2 bg-yellow-400" v-model="rating" />
  <input type="radio" name="rating-2" value="4" class="mask mask-star-2 bg-yellow-400" v-model="rating" />
  <input type="radio" name="rating-2" value="5" class="mask mask-star-2 bg-yellow-400" v-model="rating" />
</div>

                    </button>
                    <!-- Rating Modal -->
                    <dialog ref="modalRef" id="my_modal_1" class="modal">
                      <div class="modal-box bg-white">
                        <!-- Close button -->
                        <form method="dialog">
                          <button
                            class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2 text-green-600"
                            @click="closeModal"
                          >
                            ✕
                          </button>
                        </form>

                        <!-- Modal title -->
                        <h3 class="text-lg font-bold text-green-800">Rate this Recipe</h3>

                        <!-- Rating input -->
                        <div class="py-4 items-center">
                          <div class="rating">
                            <input
                              type="radio"
                              name="rating-2"
                              class="mask mask-star-2 bg-green-400"
                              value="1"
                              v-model="rating"
                              aria-label="1 star"
                            />
                            <!-- ... other rating inputs ... -->
                          </div>
                        </div>

                        <!-- Submit button -->
                        <div class="flex justify-end mt-2">
                          <button
                            @click.prevent="handleRating"
                            class="bg-green-100 ring-1 ring-green-500 py-2 px-4 rounded-full hover:bg-green-300 text-green-700"
                          >
                            Submit
                          </button>
                        </div>
                      </div>
                    </dialog>
                  </div>

                  <!-- Like/Unlike Button -->
                  <div>
                    <button
                      v-if="isLiked"
                      @click="handleRemoveLike"
                      class="flex flex-col items-center justify-center"
                    >
                      <img src="\public\icons\like-svgrepo-com.svg" class="w-8 h-8"/>
                      <span class="text-xs text-green-700 hidden md:block">
                        unlike({{ recipeDetails.likes_aggregate.aggregate.count }})
                      </span>
                    </button>
                    <button
                      v-else
                      @click="handleLikeRecipe"
                      class="flex flex-col items-center justify-center"
                    >
                      <img src="\public\icons\like-svgrepo-com.svg" class="w-8 h-8"/>
                      <span class="text-xs text-green-700 hidden md:block">
                        like({{ recipeDetails.likes_aggregate.aggregate.count }})
                      </span>
                    </button>
                  </div>

                  <!-- Comment Button -->
                  <div>
                    <button class="flex flex-col items-center justify-center" onclick="my_modal_3.showModal()">
                      <img src="/icons/comment-svgrepo-com.svg" class="w-8 h-8" alt=""/>
                      <span class="text-xs text-green-700 hidden md:block">Comment</span>
                    </button>

                    <dialog id="my_modal_3" class="modal">
                      <div class="modal-box bg-white">
                        <form method="dialog">
                          <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2 text-green-600">
                            ✕
                          </button>
                        </form>

                        <h3 class="text-lg font-bold text-green-800">
                          Hello {{ auth.$state.userName }} !
                        </h3>

                        <div v-if="commentSubmitted">
                          <p class="py-4 text-green-700">Thanks for your comment!</p>
                        </div>
                        <div v-else>
                          <input
                            v-model="comment"
                            type="text"
                            class="p-2 w-full border border-green-300 rounded-lg focus:ring-green-500 focus:border-green-500"
                            placeholder="Write a comment..."
                          />
                        </div>

                        <div class="flex justify-end mt-2">
                          <button
                            @click.prevent="handleComments"
                            class="bg-green-100 ring-1 ring-green-500 py-2 px-4 rounded-full hover:bg-green-300 text-green-700"
                          >
                            Submit
                          </button>
                        </div>
                      </div>
                    </dialog>
                  </div>

                  <!-- Bookmark Button -->
                  <button
                    v-if="isBookmarked"
                    @click="handleRemoveBookmark"
                    class="flex flex-col items-center justify-center"
                  >
                    <img src="/icons/bookkkmark.svg" class="h-8 w-8" alt="" />
                    <span class="text-xs text-green-700 hidden md:block">unsave</span>
                  </button>
                  <button
                    v-else
                    @click="handlesaveBookmark"
                    class="flex flex-col items-center justify-center"
                  >
                    <img src="/icons/save.svg" class="h-8 w-8" alt="" />
                    <span class="text-xs text-green-700 hidden md:block">save</span>
                  </button>

                  <!-- Buy Button -->
                  <button
                    v-if="!recipeStore.$state.paymentStatus === 'paid'"
                    @click="handleBuyRecipe"
                    class="btn bg-green-100 ring-1 text-center ring-green-500 py-1 md:py-2 px-2 md:px-4 rounded-full hover:bg-green-300 text-green-700"
                  >
                    {{ isLoading ? "Loading..." : "Buy Now" }}
                  </button>
                </div>
              </div>
            </div>
          </div>
          
          <!-- Right Side: Additional Recipe Images -->
          <div class="flex-[2.5] flex justify-center items-center">
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-1 lg:grid-rows-4 gap-2 w-full container mx-auto">
              <div
                v-for="(image, index) in recipeDetails.recipe_images.slice(0, 4)"
                :key="image.id"
                class="relative overflow-hidden rounded-lg shadow-md hover:shadow-lg transition-shadow duration-300 w-full"
              >
                <img
                  :src="image.image_url"
                  :alt="`Recipe Image ${index + 1}`"
                  class="w-full h-[110px] object-cover transform transition-transform duration-300 hover:scale-105"
                />
                <div class="absolute inset-0 bg-black bg-opacity-0 hover:bg-opacity-20 transition-all duration-300"></div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- Description Section -->
        <div class="flex flex-col">
          <p class="text-green-800 text-sm mb-6">
            {{ recipeDetails.description }}
          </p>
          <div class="flex items-center space-x-2">
            <div class="flex flex-row">
              <span class="text-green-800 font-semibold">Created at:</span>
              <span class="text-green-700 ml-2">
                {{ new Date(recipeDetails.created_at).toDateString() }}
              </span>
            </div>
          </div>
        </div>

        <!-- Recipe Content (Ingredients & Steps) -->
        <div>
          <!-- Payment Status Messages -->
          <button
            @click="handleVerifyPayment"
            v-if="recipeStore?.paymentDetails?.payment_status === 'pending'"
            class="p-4 bg-yellow-100 border-l-4 border-yellow-500 text-yellow-700 w-full text-center"
          >
            Click here to verify your payment and unlock the full recipe.
          </button>

          <!-- Unlocked Content -->
          <div
            v-else-if="recipeStore?.paymentDetails?.payment_status === 'paid' || 
                     Number(recipeDetails.user.id) === Number(auth.$state.userId)"
            class="p-8 flex flex-col md:flex-row"
          >
            <!-- Ingredients -->
            <div class="mb-8 md:mb-0 md:flex-1 md:pr-8 md:border-r-4 md:border-green-300">
              <h2 class="text-2xl font-bold text-green-800 mb-4 relative">
                Ingredients
                <span class="absolute bottom-0 left-0 w-16 h-1 bg-green-500 mt-2"></span>
              </h2>
              <ul class="flex flex-col gap-4">
                <li v-for="(ingredient, index) in recipeDetails.ingredients" :key="index">
                  <span class="text-green-800 font-semibold">{{ ingredient.name }}</span>
                  <span class="text-green-600"> -> </span>
                  <span class="text-green-700 ml-2">{{ ingredient.quantity }}</span>
                </li>
              </ul>
            </div>

            <!-- Steps -->
            <div class="mb-8 md:flex-1 md:pl-8">
              <h2 class="text-2xl font-bold text-green-800 mb-4 relative">
                Steps
                <span class="absolute bottom-0 left-0 w-16 h-1 bg-green-500"></span>
              </h2>
              <ol class="space-y-4">
                <li v-for="(step, index) in recipeDetails.steps" :key="index">
                  <span class="text-green-800 font-bold text-lg">Step {{ step.step_number }}:</span>
                  <span class="text-green-700 ml-2">{{ step.instruction }}</span>
                </li>
              </ol>
            </div>
          </div>
          
          <!-- Locked Content Message -->
          <div v-else class="flex justify-center items-center">
            <div class="p-4 bg-blue-100 border-l-4 w-1/2 border-blue-500 text-blue-700 flex flex-col items-center mt-3 md:mt-2">
              <p>Purchase this recipe to unlock the full instructions and ingredients.</p>
              <button
                @click="handleBuyRecipe"
                class="btn bg-green-100 mt-2 ring-1 text-center ring-green-500 py-1 md:py-2 px-2 md:px-4 rounded-full hover:bg-green-300 text-green-700"
              >
                {{ isLoading ? "Loading..." : "Buy Now" }}
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Comments Section -->
      <h1 class="text-center text-xl md:text-3xl text-green-800 font-bold pt-4 mb-2 font-poppins-italic">
        What others say about this recipe
      </h1>
      
      <div v-if="recipeDetails.comments && recipeDetails.comments.length > 0" class="flex flex-col md:flex-row gap-4 items-center justify-center">
        <div
          v-for="(comment, index) in recipeDetails.comments.slice(0, 3)"
          :key="index"
          class="card bg-green-50 text-green-800 w-96 mt-5"
        >
          <div class="card-body items-center text-center">
            <div class="flex items-center gap-3">
              <img
                :src="comment.user.profile || '/images/default-avatar.png'"
                alt="User Profile"
                class="h-10 w-10 rounded-full"
              />
              <div>
                <h4 class="font-semibold text-green-800">
                  {{ comment.user.username }}
                </h4>
                <p class="text-green-600">
                  {{ comment.comment }}
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <div v-else class="flex flex-col items-center justify-center text-green-400 text-center relative group">
        <img
          src="\public\Empty-cuate.svg"
          alt="No comments"
          class="h-[300px] w-[300px] transition-transform duration-300 group-hover:scale-105"
        />
        <div class="absolute inset-0 bg-green-800 bg-opacity-70 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity duration-300">
          <span class="text-white font-bold text-lg">
            No comments yet. Be the first to leave one!
          </span>
        </div>
      </div>

      <!-- More Recipes by Author -->
      <div class="p-6">
        <h1 class="text-xl md:text-3xl font-bold font-poppins-italic mb-8 text-green-800 text-center">
          More recipes by 
          <span class="text-green-800 text-xl md:text-3xl underline decoration-green-500 decoration-4">
            {{ recipeDetails.user.username }}
          </span>
        </h1>

        <div class="flex flex-col md:flex-row gap-4 items-center justify-center">
          <NuxtLink
            v-for="(recipe, index) in recipeDetails.user.recipes.slice(0, 5)"
            :key="index"
            :to="`/recipes/${recipe.id}`"
            class="block overflow-hidden hover:-translate-y-2"
          >
            <div class="relative h-48 overflow-hidden">
              <img
                :src="recipe.featured_image || '\public\Chef-cuate.svg'"
                :alt="recipe.title"
                class="w-full h-full object-cover"
              />
              <div class="absolute inset-0 bg-black bg-opacity-0 hover:bg-opacity-20 transition-all duration-300"></div>
            </div>
            <div class="p-4">
              <h2 class="text-xl font-semibold text-green-800 mb-2">
                {{ recipe.title }}
              </h2>
            </div>
          </NuxtLink>
        </div>
      </div>
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

/* Custom green-themed styles */
.modal-box {
  background-color: white;
  color: #1a5632;
}

.card {
  background-color: #f0fff4;
  border: 1px solid #c6f6d5;
}

.btn {
  transition: all 0.3s ease;
}

.rating input:checked ~ input {
  color: #c6f6d5;
}
</style>