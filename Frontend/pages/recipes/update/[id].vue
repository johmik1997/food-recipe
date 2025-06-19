<script setup>
import { ref, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { Form, Field } from "vee-validate";
import * as yup from "yup";
import { useToast } from "vue-toast-notification";

const toast = useToast();
const route = useRoute();
const router = useRouter();

const recipeStore = useRecipeStore();
const useAuthStore = authStore();

const recipeId = parseInt(route.params.id);
const user_id = parseInt(useAuthStore.$state.userId);
const MAX_FILE_SIZE = 2 * 1024 * 1024;

const recipe = ref(null);
const categories = ref([]);

const image = ref({
  name: "",
  type: "",
  base64String: "",
});

const getExtension = (filename) => filename.split(".").pop();

const handleFileChange = (event) => {
  const file = event.target.files[0];
  if (file) {
    image.value.name = file.name;
    image.value.type = file.type;
    const reader = new FileReader();
    reader.onload = () => {
      const base64 = reader.result.split("base64,")[1];
      image.value.base64String = base64;
    };
    reader.readAsDataURL(file);
    reader.onerror = (error) => {
      console.error("Error reading image file:", error);
    };
  }
};

const schema = yup.object({
  images: yup
    .mixed()
    .test("fileType", "Please provide a supported file type", (file) => {
      if (!file) return true;
      return ["jpg", "jpeg", "png", "webp", "svg", "gif"].includes(getExtension(file.name));
    })
    .test("fileSize", "File is too big. Max size is 2MB.", (file) => {
      if (!file) return true;
      return file.size <= MAX_FILE_SIZE;
    }),
});

const handleUpdateRecipe = async (values) => {
  try {
    const payload = {
      title: values.title,
      description: values.description,
      prep_time: values.prep_time,
      cook_time: values.cook_time,
      servings: values.servings,
      price: values.price,
      category_id: values.category_id,
      recipe_id: recipeId,
      user_id: user_id,
      imageName: image.value.name,
      imageType: image.value.type,
      base64String: image.value.base64String,
    };

    const res = await recipeStore.updateRecipe(payload);
    if (res) {
      await recipeStore.singleRecipe(recipeId);
      const message = recipeStore.$state.successmessage || "Recipe updated successfully!";
      recipeStore.setSuccessMessage("");
      toast.success(message);
      router.push("/recipes/" + recipeId);
    }
  } catch (error) {
    console.error("Error updating recipe:", error);
    toast.error("Failed to update recipe.");
  }
};

const handleFetchCategories = async () => {
  await recipeStore.getCategories();
  categories.value = recipeStore.categories;
};

const fetchRecipe = async () => {
  await recipeStore.singleRecipe(recipeId);
  recipe.value = recipeStore.$state.recipe;
};

onMounted(async () => {
  await handleFetchCategories();
  await fetchRecipe();
});
</script>
<template>
  <div class="container mx-auto p-4">
    <div class="max-w-4xl mx-auto bg-white dark:bg-[#20161F] shadow-lg rounded-lg p-6 shadow-green-400">
      <h1 class="text-2xl font-bold mb-6 text-center">Update Recipe</h1>

      <Form
        v-if="recipe"
        @submit="handleUpdateRecipe"
        :validation-schema="schema"
        :initial-values="{
          title: recipe.title || '',
          description: recipe.description || '',
          prep_time: recipe.prep_time || 0,
          cook_time: recipe.cook_time || 0,
          servings: recipe.servings || 0,
          category_id: recipe.category_id || '',
          price: recipe.price || 0,
          images: null
        }"
        v-slot="{ errors }"
        class="space-y-6"
      >
        <div class="pl-4 pr-4">
          <div class="flex flex-row gap-2">
            <UIInput name="title" placeholder="Recipe Title" label="Title" class="w-full" type="text" />

            <UIInput name="prep_time" placeholder="Preparation time" label="Prep Time (min)" type="number" class="w-full" />
          </div>

          <div>
            <UIInput name="description" placeholder="Enter description" label="Description" type="text" />
            <UIInput name="cook_time" label="Cook Time (min)" placeholder="Enter cooking time" type="number" />
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <UIInput name="servings" label="Servings" placeholder="e.g. 4" type="number" />
<UIInput name="price" label="Price" placeholder="e.g. 25" type="number" />
            <div>
              <label for="category_id" class="block text-sm font-medium text-gray-700 dark:text-gray-50 mb-2">
                Category
              </label>
              <Field
                as="select"
                name="category_id"
                class="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="">Select a category</option>
                <option v-for="category in categories" :key="category.id" :value="category.id">
                  {{ category.name }}
                </option>
              </Field>
            </div>
          </div>

          <!-- Image Upload -->
          <div class="w-full mt-4">
            <h2 class="block pl-3 ml-px text-sm font-medium text-gray-700 dark:text-gray-200">Images</h2>
            <UIMultipleImages
              name="images"
              @image-changed="handleFileChange"
              :isMultiple="false"
              :error-message="errors.images"
            />
          </div>

          <!-- Submit Button -->
          <div class="flex justify-center mt-6">
            <button
              type="submit"
              class="bg-cyan-500 text-white px-6 py-2 hover:bg-cyan-700 transition-colors rounded-full ring-2 ring-cyan-300"
            >
              <span v-if="!recipeStore.$state.isLoading">Update Recipe</span>
              <span v-else>
                <UILoading />
              </span>
            </button>
          </div>
        </div>
      </Form>
    </div>
  </div>
</template>
