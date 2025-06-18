<script setup>
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

const MAX_FILE_SIZE = 2097152;

const getExtension = (filename) => {
  const extension = filename.split(".");
  return extension[1];
};

const image = ref({
  name: "",
  type: "",
  base64String: "",
});

const handleFileChange = (event) => {
  const file = event.target.files[0];
  if (file) {
    image.value.name = file.name;
    image.value.type = file.type;
    const reader = new FileReader();
    reader.onload = () => {
      const base64String = reader.match(/base64,(.*)$/)[1];
      image.value.base64String = base64String;
    };
    reader.onerror = (error) => {
      console.log(" error uploading the image", error);
    };
  }
};

const schema = yup.object({
  images: yup
    .mixed()
    .test({
      message: "please provide a supported file type",
      test: (file, context) => {
        let isValid;
        if (file) {
          isValid = ["jpg", "gif", "png", "jpeg", "svg", "webp"].includes(
            getExtension(file?.name)
          );
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
      message: " File is too big, can't exit 2mb",
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

const handleUpdateRecipe = async (value) => {
  try {
    if (!value) {
      console.log("no input is provided");
      return;
    }
    const res = await recipeStore.updateRecipe({
      title: value.title,
      description: value.description,
      preparation_time: value.preparation_time,
      category_id: value.category_id,
      price: value.price,
      recipe_id: recipeId,
      user_id: user_id,
      imageName: image.value.name,
      imageType: image.value.type,
      base64String: image.value.base64String,
    });
    console.log("updated recipe detail", res);
    if (res) {
      await recipeStore.singleRecipe(recipeId);
      const message =
        recipeStore.$state.successmessage || "recope updated recipe";
      recipeStore.setSuccessMessage("");
      toast.success(message);
      router.push("/recipes/" + recipeId);
    } else {
      console.log("error updating the recipe");
    }
  } catch (error) {
    toast.error("error updating recipe");
    console.log("error updating the recipe", error);
  }
};

const categories = ref([]);
const handleFetchCatagories = async () => {
  await recipeStore.getCategories();
  categories.value = recipeStore.categories;
  console.log("categories", JSON.stringify(categories.value));
};
onMounted(handleFetchCatagories);
</script>

<template>
  <div class="conatiner mx-auto p-4">
    <div
      class="max-w-4xl mx-auto bg-white dark:bg-[#20161F] shadow-lg rounded-lg p-6 shadow-green-400"
    >
      <h1 class="text-2xl font-bold mb-6 text-center">update recipe</h1>
      <Form
        @submit="handleUpdateRecipe"
        :validation-schema="schema"
        class="space-y-6"
      >
        <div class="pl-4 pr-4">
          <div class="flex flex-row gap-2">
            <UIInput
              name="title"
              placeholder="recipe title"
              label="Title"
              class="w-full"
              type="text"
            />
            <UIInput
              name="prep_time"
              placeholder="enter preparation time"
              label="Preparation Time (min)"
              type="number"
              class="w-full"
            />
          </div>
          <div>
            <UIInput
              name="description"
              placeholder="enter description of the recipe"
              label="Description"
              type="text"
            />
          </div>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <UIInput
              name="price"
              label="Price"
              placeholder="enter price here.."
              type="number"
            />

            <div>
              <label
                for="category_id"
                class="block text-sm font-medium text-gray-700 dark:text-gray-50 mb-2"
                >Category
              </label>
              <Field
                as="select"
                name="category_id"
                class="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="">Select a category</option>
                <option
                  v-for="category in categories"
                  :key="category.id"
                  :value="category.id"
                >
                  {{ category.name }}
                </option>
              </Field>
            </div>
          </div>
          <div class="flex justify-center mt-4">
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
