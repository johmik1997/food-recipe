<script setup>
import { ref } from "vue";
import { Form, Field } from "vee-validate";
import { useToast } from "vue-toast-notification";
import * as yup from "yup";
import { useRouter } from "vue-router";
useSeoMeta({
  title: "recipe-app | Profile Page",
  description: "The project app meta.",
});

const router = useRouter();
const toast = useToast();
const bookmarkStore = useBookmarkStore();
const bookMarkId = computed(() => bookmarkStore.bookmarkedId);
const recipeStore = useRecipeStore();
const likeStore = useLikeStore();
const likedId = computed(() => likeStore.likeId);
console.log("likedId from profile page", likedId.value);
const likedRecipes = ref([]);
const getExtension = (filename) => {
  const extension = filename.split(".");
  return extension[1];
};
const MAX_FILE_SIZE = 2097152;
const auth = authStore();

const profileError = ref({
  userName: "",
  phone: "",
  profile: "",
});
const nameregex = /^([^\x00-\x7F]|[\a-zA-Z_\ \.\+\-]){2,20}$/;
const phoneregex = /^(^\+251|^251|^0)?(9|7)\d{8}$/;
const errorMessages = {
  userName: {
    invalid_username: "Invalid user name",
  },
  phone: {
    invalid_phone: "Invalid phone address",
  },
};

const schema = yup.object({
  userName: yup
    .string()
    .required()
    .matches(nameregex, errorMessages.userName.invalid_username),
  phone: yup
    .string()
    .required()
    .matches(phoneregex, errorMessages.phone.invalid_phone),
  profile: yup
    .mixed()
    .test({
      message: "Please provide supported file type",
      test: (file, context) => {
        let isValid;
        if (file) {
          isValid = [
            "jpg",
            "gif",
            "png",
            "jpeg",
            "svg",
            "webp",
            "jpeg",
          ].includes(getExtension(file?.name));
        } else {
          isValid = true;
        }
        if (!isValid) {
          context?.createError();
        }
        return isValid;
      },
    })
    .test({
      message: "File is too big, cant exit 2Mb",
      test: (file) => {
        let isValid;
        if (file) {
          isValid = file?.size < MAX_FILE_SIZE;
        } else {
          isValid = true;
        }
        return isValid;
      },
    }),
});

const handleFileChange = (event) => {
  const file = event.target.files[0];
  if (file) {
    image.value.name = file.name;
    image.value.type = file.type;
    const reader = new FileReader();
    reader.onload = () => {
      const base64String = reader.result.match(/base64,(.*)$/)[1];
      image.value.base64String = base64String;
    };
    reader.readAsDataURL(file);
    reader.onerror = (error) => {
      console.log("there is some error in image uploading: ", error);
      reject(error);
    };
  }
};
const handleFectchProfile = async () => {
  try {
    await auth.getProfile();
    likedRecipes.value = recipeStore.$state.recipes;
    console.log("liked recipes of a user", JSON.stringify);
  } catch (error) {}
};
onMounted(async () => {
  handleFectchProfile();
  await recipeStore.getSellerRecipes();
});

const image = ref({
  name: "",
  type: "",
  base64String: "",
});

const handleUpdateProfle = async (value) => {
  auth.changeOnLoad(true);
  let payload;
  if (
    image.value.name !== "" ||
    image.value.type !== "" ||
    image.value.base64String !== ""
  ) {
    payload = {
      userName: value.userName,
      phone: value.phone,
      imageName: image.value.name,
      imageType: image.value.type,
      imageString: image.value.base64String,
    };
  } else {
    payload = {
      userName: value.userName,
      phone: value.phone,
    };
  }
  const result = await auth.updateProfile(payload);
  auth.changeOnLoad(false);
  if (result) {
    await auth.getProfile();
    if (
      auth.$state.successMessage &&
      auth.$state.successMessage.length > 0
    ) {
      const message = auth.$state.successMessage;
      auth.setSuccessMessage("");
      toast.success(message);
    } else {
      toast.success("Your profile updated successfully!");
    }
    showProfileUpdate.value = false;
  } else {
    if (
      auth.$state.errorMessage &&
      auth.$state.errorMessage.length > 0
    ) {
      const message = auth.$state.errorMessage;
      auth.setErrorMessage("");
      toast.error(message);
    } else {
      toast.error("Something wrong! please try again.");
    }
  }
};

// Reactive state for toggling sections
const showProfileUpdate = ref(false);

const toggleProfileUpdate = () => {
  showProfileUpdate.value = !showProfileUpdate.value;
};

// Unlike a recipe
const unlikeRecipe = async (recipeId, likedId) => {
  try {
    const res = await recipeStore.unlikeRecipe({
      recipe_id: recipeId,
      id: likedId,
    });
    if (res) {
      await auth.getProfile();
      likedRecipes.value = likedRecipes.value.filter(
        (recipe) => recipe.id !== recipeId
      );

      console.log("recipe id to unlike", recipeId);
      console.log("Recipe unliked successfully!");
      toast.success("Recipe unliked successfully!");
    }
  } catch (error) {
    console.error("Error unliking recipe:", error);
    toast.error("Failed to unlike recipe. Please try again.");
  }
};

// Unsave a recipe
const unsaveRecipe = async (recipeId, bookmarkId) => {
  console.log("recipe id to unsave", recipeId, "bookmark id", bookmarkId);
  try {
    const res = await recipeStore.unsaveRecipe({
      recipe_id: recipeId,
      id: bookmarkId,
    });
    if (res) {
      toast.success("Recipe unsaved successfully!");
      await auth.getProfile();
    }
  } catch (error) {
    console.error("Error unsaving recipe:", error);
    toast.error("Something went wrong! Please try again.");
  }
};

const user_id = Number(auth.$state.userId);

// Delete Confirmation Modal Logic
const showDeleteModal = ref(false);
const recipeToDeleteId = ref(null);
const deleteModalRef = ref(null);

const openDeleteModal = (recipeId) => {
  recipeToDeleteId.value = recipeId;
  deleteModalRef.value.showModal();
};

const closeDeleteModal = () => {
  deleteModalRef.value.close();
};

const confirmDelete = async () => {
  if (recipeToDeleteId.value) {
    try {
      await recipeStore.deleteRecipes({ id: recipeToDeleteId.value, user_id });
      toast.success("Recipe deleted successfully!");
      await auth.getProfile();
      await recipeStore.getSellerRecipes();
      closeDeleteModal();
    } catch (error) {
      console.error("Error deleting recipe:", error);
      toast.error("Something went wrong! Please try again.");
    }
  }
  showDeleteModal.value = false;
};
</script>

<template>
  <div class="min-h-50vh container mx-auto p-6 bg-white">
    <!-- Header Section -->
    <div class="flex items-center space-x-4">
      <div class="avatar online">
        <div class="w-24 rounded-full border-2 border-green-500">
          <img
            :src="auth.$state.user.profile || '/images/default-avatar.png'"
            :alt="`${auth.$state.user.username}'s profile`"
            class="transform hover:scale-150 hover:rotate-2 transition duration-300 ease-in-out rounded-full w-24 h-24 object-cover"
          />
        </div>
      </div>
      <div>
        <h1 class="text-2xl font-bold text-green-800">
          {{ auth.$state.user.username }}
        </h1>
        <p class="text-green-600">
          {{ auth.$state.user.email }}
        </p>
        <p class="text-green-600">
          {{ auth.$state.user.phone }}
        </p>
        <p class="text-green-700 font-medium">{{ auth.$state.user.role }}</p>

        <p class="text-green-800">
          Joined in:
          <span class="text-green-600">
            {{
              new Date(auth.$state.user.created_at).toDateString()
            }}</span
          >
        </p>
      </div>
      <button
        @click="toggleProfileUpdate"
        class="ml-auto px-4 py-2 bg-green-100 rounded-full text-green-800 hover:bg-green-200 border border-green-300 transition-colors duration-200"
      >
        Update Profile
      </button>
    </div>

    <!-- Profile Update Section -->
    <transition name="slide">
      <div
        v-if="showProfileUpdate"
        class="mt-6 bg-white shadow rounded-lg p-6 border border-green-200"
      >
        <h2 class="text-lg font-bold mb-4 text-green-800">Update Profile</h2>
        <Form
          @submit="handleUpdateProfle"
          :validation-schema="schema"
          v-slot="{ errors }"
          class="space-y-4"
        >
          <div>
            <UIInput
              name="userName"
              type="text"
              placeholder="enter username"
              label="User Name"
              :is-required="true"
              :edit-mode="true"
              :value="
                auth.$state.user.username
                  ? auth.$state.username
                  : ''
              "
              :error-message="
                errors.userName
                  ? errors.userName
                  : profileError.userName
                  ? profileError.userName
                  : ''
              "
              class="mt-1 block w-full border-green-300 rounded-md shadow-sm focus:border-green-500 focus:ring-green-500"
            />
          </div>
          <div>
            <UIInput
              name="phone"
              placeholder="Insert your new phonenumber"
              label="Phone Number"
              :is-required="true"
              :edit-mode="true"
              :value="
                auth.$state.user.phone
                  ? auth.$state.user.phone
                  : ''
              "
              type="text"
              :error-message="
                errors.phone
                  ? errors.phone
                  : profileError.phone
                  ? profileError.phone
                  : ''
              "
            />
          </div>
          <div class="mt-6">
            <button
              class="w-full px-6 py-3 text-sm font-medium text-white capitalize transition-colors duration-300 transform bg-green-600 rounded-full hover:bg-green-700 focus:outline-none focus:ring focus:ring-green-300 focus:ring-opacity-50"
            >
              <span v-if="auth.$state.onLoad == false">
                Update Profile
              </span>
              <div class="w-full flex items-center justify-center gap-2" v-else>
                <UILoading />
              </div>
            </button>
          </div>
        </Form>
      </div>
    </transition>

    <!-- Delete Confirmation Modal -->
    <dialog ref="deleteModalRef" class="modal">
      <div class="modal-box bg-white border border-green-200">
        <!-- Close button -->
        <form method="dialog">
          <button
            class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2 text-green-600 hover:bg-green-100"
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

    <!-- My Recipes Section -->
    <section class="container mx-auto p-4">
      <div class="mt-8">
        <h2 class="text-3xl font-bold text-green-800 mb-8 text-center">
          My Recipes
        </h2>
        <div
          v-if="
            auth.$state.user.recipes &&
            auth.$state.user.recipes.length == 0
          "
          class="text-center flex flex-col justify-center items-center relative group"
        >
          <img
            src="public\Empty-cuate.svg" 
            class="h-[300px] w-[300px] transition-transform duration-300 group-hover:scale-105"
            alt="No recipes created yet"
          />

          <!-- Overlay -->
          <div
            class="absolute inset-0 bg-green-800 bg-opacity-70 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity duration-300 rounded-lg"
          >
            <span class="text-white font-bold text-lg"
              >No recipes yet. Time to get cooking!</span
            >
          </div>
        </div>
        <div
          v-else
          class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-6"
        >
          <div
            v-for="recipe in auth.$state.user.recipes"
            :key="recipe.id"
            class="bg-white p-4 border border-green-100 rounded-lg shadow-sm hover:shadow-md transition-shadow duration-300 cursor-pointer relative"
          >
            <!-- Delete Button -->
            <button
              @click.stop="openDeleteModal(recipe.id)"
              class="absolute top-3 right-3 p-2 bg-white rounded-full shadow-md hover:bg-green-100 transition-colors duration-300 text-green-800"
              aria-label="Delete recipe"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                stroke-width="1.5"
                stroke="currentColor"
                class="w-6 h-6 text-green-700"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0"
                />
              </svg>
            </button>

            <!-- Recipe Image -->
            <img
              :src="
                recipe.featured_image ||
                '/public/Chef-cuate.svg'
              "
              :alt="recipe.title"
              class="w-full h-48 object-cover rounded-lg mb-4"
            />

            <!-- Recipe Title -->
            <NuxtLink :to="`/recipes/${recipe.id}`">
              <h3 class="text-xl font-bold text-green-800 mb-2 truncate">
                {{ recipe.title }}
              </h3>
            </NuxtLink>

            <!-- Recipe Details -->
            <div class="flex items-center justify-between">
              <p class="text-sm text-green-600">
                Price: ${{ recipe.price }}
              </p>
              <p class="text-sm text-green-600">
                Rating: {{ recipe.average_rating }}
              </p>
            </div>
          </div>
        </div>
      </div>

      <!-- Liked Recipes Section -->
      <div class="mt-12">
        <h2 class="text-3xl font-bold text-green-800 mb-8 text-center">
          Liked Recipes
        </h2>
        <div
          v-if="
            auth.$state.user.likes &&
            auth.$state.user.likes.length === 0
          "
          class="text-center flex flex-col justify-center items-center relative group"
        >
          <img
            src="public\Money stress-bro (1).svg"
            class="h-[300px] w-[300px] transition-transform duration-300 group-hover:scale-105"
            alt="No favorite recipes yet"
          />

          <!-- Overlay -->
          <div
            class="absolute inset-0 bg-green-800 bg-opacity-70 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity duration-300 rounded-lg"
          >
            <span class="text-white font-bold text-lg"
              >No favorite recipes yet. Start exploring!</span
            >
          </div>
        </div>
        <div
          v-else
          class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-6"
        >
          <div
            v-for="like in auth.$state.user.likes"
            :key="like.recipe_id"
            class="bg-white p-4 border border-green-100 rounded-lg shadow-sm hover:shadow-md transition-shadow duration-300 cursor-pointer relative"
          >
            <!-- Unlike Button -->
            <button
              @click.stop="unlikeRecipe(like.recipe_id, like.id)"
              class="absolute top-3 right-3 p-2 bg-white rounded-full shadow-md hover:bg-green-100 transition-colors duration-300 text-green-800"
              aria-label="Unlike recipe"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                stroke-width="1.5"
                stroke="currentColor"
                class="w-6 h-6 text-green-700"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M21 8.25c0-2.485-2.099-4.5-4.688-4.5-1.935 0-3.597 1.126-4.312 2.733-.715-1.607-2.377-2.733-4.313-2.733C5.1 3.75 3 5.765 3 8.25c0 7.22 9 12 9 12s9-4.78 9-12Z"
                />
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M3 3l18 18"
                />
              </svg>
            </button>

            <!-- Recipe Image -->
            <img
              :src="
                like.recipe.featured_image ||
                '/\public\Chef-cuate.svg'
              "
              :alt="like.recipe.title"
              class="w-full h-48 object-cover rounded-lg mb-4"
            />

            <!-- Recipe Title -->
            <NuxtLink :to="`/recipes/${like.recipe.id}`">
              <h3 class="text-xl font-bold text-green-800 mb-2 truncate">
                {{ like.recipe.title }}
              </h3>
            </NuxtLink>

            <!-- Recipe Details -->
            <div class="flex items-center justify-between">
              <p class="text-sm text-green-600">
                Price: ${{ like.recipe.price }}
              </p>
              <p class="text-sm text-green-600">
                Rating: {{ like.recipe.average_rating }}
              </p>
            </div>
          </div>
        </div>
      </div>

      <!-- Bookmarked Recipes Section -->
      <div class="mt-12">
        <h2 class="text-3xl font-bold text-green-800 mb-8 text-center">
          Saved Recipes
        </h2>
        <div
          v-if="
            auth.$state.user.bookmarks &&
            auth.$state.user.bookmarks.length === 0
          "
          class="text-center flex flex-col justify-center items-center relative group"
        >
          <img
            src="\public\Empty-cuate.svg"
            class="h-[300px] w-[300px] transition-transform duration-300 group-hover:scale-105"
            alt="No saved recipes yet"
          />

          <!-- Overlay -->
          <div
            class="absolute inset-0 bg-green-800 bg-opacity-70 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity duration-300 rounded-lg"
          >
            <span class="text-white font-bold text-lg">Save recipes to view them here</span>
          </div>
        </div>
        <div
          v-else
          class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-6"
        >
          <div
            v-for="bookmark in auth.$state.user.bookmarks"
            :key="bookmark.recipe_id"
            class="bg-white p-4 border border-green-100 rounded-lg shadow-sm hover:shadow-md transition-shadow duration-300 cursor-pointer relative"
          >
            <!-- Unsave Button -->
            <button
              @click.stop="unsaveRecipe(bookmark.recipe_id, bookmark.id)"
              class="absolute top-3 right-3 p-2 bg-white rounded-full shadow-md hover:bg-green-100 transition-colors duration-300 text-green-800"
              aria-label="Unsave recipe"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                stroke-width="1.5"
                stroke="currentColor"
                class="w-6 h-6 text-green-700"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M17.593 3.322c1.1.128 1.907 1.077 1.907 2.185V21L12 17.25 4.5 21V5.507c0-1.108.806-2.057 1.907-2.185a48.507 48.507 0 0 1 11.186 0Z"
                />
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M3 3l18 18"
                />
              </svg>
            </button>

            <!-- Recipe Image -->
            <img
              :src="
                bookmark.recipe.featured_image ||
                '/public\Chef-cuate.svg'
              "
              :alt="bookmark.recipe.title"
              class="w-full h-48 object-cover rounded-lg mb-4"
            />

            <!-- Recipe Title -->
            <NuxtLink :to="`/recipes/${bookmark.recipe.id}`">
              <h3 class="text-xl font-bold text-green-800 mb-2 truncate">
                {{ bookmark.recipe.title }}
              </h3>
            </NuxtLink>

            <!-- Recipe Details -->
            <div class="flex items-center justify-between">
              <p class="text-sm text-green-600">
                Price: ${{ bookmark.recipe.price }}
              </p>
              <p class="text-sm text-green-600">
                Rating: {{ bookmark.recipe.average_rating }}
              </p>
            </div>
          </div>
        </div>
      </div>

      <!-- Sold Recipes Section -->
      <div class="mt-12">
        <h2 class="text-3xl font-bold text-green-800 mb-8 text-center">
          Sold Recipes
        </h2>
        <div
          v-if="
            recipeStore.$state.soldRecipes &&
            recipeStore.$state.soldRecipes.length === 0
          "
          class="text-center flex flex-col justify-center items-center relative group"
        >
          <img
            src="/public\Empty-cuate.svg"
            class="h-[300px] w-[300px] transition-transform duration-300 group-hover:scale-105"
            alt="No recipes sold yet"
          />

          <!-- Overlay -->
          <div
            class="absolute inset-0 bg-green-800 bg-opacity-70 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity duration-300 rounded-lg"
          >
            <span class="text-white font-bold text-lg"
              >Your delicious recipes could be here!</span
            >
          </div>
        </div>

        <div
          v-else
          class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-6"
        >
          <div
            v-for="recipe in recipeStore.$state.soldRecipes"
            :key="recipe.recipe.id"
            class="bg-white p-4 border border-green-100 rounded-lg shadow-sm hover:shadow-md transition-shadow duration-300 cursor-pointer"
          >
            <!-- Recipe Image -->
            <img
              :src="
                recipe.recipe.featured_image ||
                '/images/recipe-placeholder.jpg'
              "
              :alt="recipe.recipe.title"
              class="w-full h-48 object-cover rounded-lg mb-4"
            />

            <!-- Recipe Title -->
            <h3 class="text-xl font-bold text-green-800 mb-2 truncate">
              {{ recipe.recipe.title }}
            </h3>

            <!-- Recipe Details -->
            <div class="flex items-center justify-between">
              <div class="flex items-center space-x-2">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  class="h-6 w-6 text-green-600"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M5 13l4 4L19 7"
                  />
                </svg>
                <span class="text-lg font-semibold text-green-600">
                  {{ recipe.recipe.sold_recipes_aggregate.aggregate.count }}
                  Sold
                </span>
              </div>

              <span
                class="text-lg font-bold text-green-700"
              >
                ${{ recipe.recipe.price }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<style>
.slide-enter-active,
.slide-leave-active {
  transition: all 0.3s ease;
}
.slide-enter-from,
.slide-leave-to {
  transform: translateY(-20px);
  opacity: 0.3;
}

/* Custom styles */
.avatar.online::before {
  background-color: #10b981 !important;
  border-color: white;
}

.card {
  transition: all 0.3s ease;
}

.card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 25px -5px rgba(5, 150, 105, 0.1);
}

/* Hide scrollbar */
::-webkit-scrollbar {
  display: none;
}

/* Button transitions */
button {
  transition: all 0.2s ease;
}

/* Image hover effects */
img {
  transition: transform 0.3s ease;
}

img:hover {
  transform: scale(1.03);
}
</style>