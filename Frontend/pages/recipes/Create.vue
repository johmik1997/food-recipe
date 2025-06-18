<script setup>
import { Form, Field } from "vee-validate";
import * as yup from "yup";
import { useToast } from "vue-toast-notification";

const toast = useToast();

useSeoMeta({
  title: "recipe-app | Create Recipe",
  description: "Create a new recipe for my cool recipe app",
});

const recipeStore = useRecipeStore();
const auth = authStore();

const MAX_FILE_SIZE = 2097152;
const router = useRouter();

const getExtension = (filename) => {
  return filename.split(".").pop();
};

const thumbnail = ref({
  name: "",
  type: "",
  base64String: "",
});
const images = ref([]);

const handleFileChange = ({
  thumbnail: selectedThumbnail,
  images: uploadedImages,
}) => {
  thumbnail.value = selectedThumbnail || {
    name: "",
    type: "",
    base64String: "",
  };

  images.value = uploadedImages || [];
};

// error messages
const errorMessages = {
  title: {
    required: "Title is required!",
  },
  description: {
    required: "Description is required!",
  },
  prep_time: {
    required: "Preparation time is required!",
  },
  cook_time: {
    required: "Cooking time is required!",
  },
  servings: {
    required: "Servings is required!",
  },
  category_id: {
    required: "Category is required!",
  },
  price: {
    required: "Price is required!",
  },
  ingredients: {
    required: "At least one ingredient is required!",
  },
  steps: {
    required: "At least one step is required!",
  },
};

// validation schema
const schema = yup.object({
  title: yup.string().required(errorMessages.title.required),
  description: yup.string().required(errorMessages.description.required),
  prep_time: yup.number().required(errorMessages.prep_time.required),
  cook_time: yup.number().required(errorMessages.cook_time.required),
  servings: yup.number().required(errorMessages.servings.required),
  category_id: yup.number().required(errorMessages.category_id.required),
  price: yup.number().required(errorMessages.price.required),
  ingredients: yup
    .array()
    .of(
      yup.object({
        name: yup.string().required("Ingredient name is required!"),
        quantity: yup.string().required("Ingredient quantity is required!"),
      })
    )
    .min(1, errorMessages.ingredients.required),
  steps: yup
    .array()
    .of(
      yup.object({
        step_number: yup.number().required("Step number is required!"),
        instruction: yup.string().required("Step instruction is required!"),
      })
    )
    .min(1, errorMessages.steps.required),
  thumbnail: yup
    .mixed()
    .test({
      message: "Please provide a supported file type",
      test: (file) => {
        if (!file) return true;
        return ["jpg", "gif", "png", "jpeg", "svg"].includes(
          getExtension(file.name)
        );
      },
    })
    .test({
      message: "File is too big, must be less than 2MB",
      test: (file) => {
        if (!file) return true;
        return file.size <= MAX_FILE_SIZE;
      },
    }),
  images: yup.array().of(
    yup
      .mixed()
      .test({
        message: "Please provide a supported file type",
        test: (file) => {
          if (!file) return true;
          return ["jpg", "gif", "png", "jpeg", "svg", "webp"].includes(
            getExtension(file.name)
          );
        },
      })
      .test({
        message: "File is too big, must be less than 2MB",
        test: (file) => {
          if (!file) return true;
          return file.size <= MAX_FILE_SIZE;
        },
      })
  ),
});

const handleCreateRecipe = async (value) => {
  if (!value) {
    console.log("no input is provided!");
    return;
  }

  try {
    // step 1: Create the recipe
    const createRecipeResult = await recipeStore.createRecipe({
      title: value.title,
      description: value.description,
      prep_time: value.prep_time,
      cook_time: value.cook_time,
      servings: value.servings,
      category_id: value.category_id,
      price: value.price,
      user_id: Number(auth.$state.userId),
      ingredients: value.ingredients
        .map((ingredient) => ({
          name: ingredient.name || "",
          quantity: ingredient.quantity || "",
        }))
        .filter((ingredient) => ingredient.name && ingredient.quantity),
      steps: value.steps
        .map((step) => ({
          step_number: step.step_number || 0,
          instruction: step.instruction || "",
        }))
        .filter((step) => step.step_number && step.instruction),
      imageName: thumbnail.value.name,
      imageType: thumbnail.value.type,
      base64string: thumbnail.value.base64String,
    });

    console.log("createRecipeResult", createRecipeResult);
    if (createRecipeResult) {
      const message = "Recipe added successfully!  uploading the images...";
      recipeStore.setSuccessMessage("");
      toast.success(message);
      const recipe_id = recipeStore.$state.createdRecipe.id;
      console.log("created_recipe id", recipe_id);

      const featuredImageIndex = images.value.findIndex(
        (img) => img.name === thumbnail.value.name
      );
      // Step 2: upload images
      const uploadPayload = {
        recipe_id,
        featuredImageIndex,
        images: images.value
          .filter((img) => img.name && img.type && img.base64String)
          .map((img) => ({
            name: img.name,
            type: img.type,
            base64String: img.base64String,
          })),
      };
      console.log("recipe id", recipe_id);
      console.log("uploadPayload", uploadPayload);

      const uploadResult = await recipeStore.uploadRecipeImages(uploadPayload);
      console.log(uploadResult);

      await recipeStore.getAllRecipes();
    }

    toast.success("Recipe created successfully!");
    router.push("/recipes");
  } catch (error) {
    console.error("Error creating recipe or uploading images:", error);
    toast.error(error.message || "Something went wrong.");
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
  <div class="container mx-auto p-4">
    <div
      class="max-w-4xl mx-auto bg-white dark:bg-gray-800 shadow-lg rounded-lg p-6 shadow-green-500"
    >
      <h1 class="text-2xl font-bold mb-6 text-center text-green-600">
        Create a New Recipe
      </h1>
      <Form
        @submit="handleCreateRecipe"
        :validation-schema="schema"
        v-slot="{ errors }"
        class="space-y-6"
      >
        <!-- Title and Preparation Time -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <UIInput
            name="title"
            placeholder="Recipe Title"
            label="Title"
            :error-message="errors.title"
            :is-required="true"
          />
          <UIInput
            name="prep_time"
            placeholder="Preparation Time"
            label="Preparation Time (min)"
            :error-message="errors.prep_time"
            :is-required="true"
            type="number"
          />
          <UIInput
            name="cook_time"
            placeholder="cooking Time"
            label="Cooking Time (min)"
            :error-message="errors.cook_time"
            :is-required="true"
            type="number"
          />
        </div>

        <!-- Servings -->
        <UIInput
          name="servings"
          placeholder="Servings"
          label="Servings (eg: 4)"
          :error-message="errors.servings"
          :is-required="true"
          type="number"
        />

        <!-- Description -->
        <UIInput
          name="description"
          placeholder="Recipe Description"
          label="Description"
          :error-message="errors.description"
          :is-required="true"
          type="textarea"
        />

        <!-- Category and Price -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label
              for="category_id"
              class="block text-sm font-medium text-gray-700 dark:text-gray-50 mb-2"
            >
              Category
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
            <span v-if="errors.category_id" class="text-red-500 text-sm">
              {{ errors.category_id }}
            </span>
          </div>
          <div>
            <UIInput
              name="price"
              placeholder="Price"
              label="Price"
              :error-message="errors.price"
              :is-required="true"
              type="number"
            />
          </div>
        </div>

        <!-- Ingredients -->
        <div>
          <h2
            class="block pl-3 ml-px text-sm font-medium text-gray-700 dark:text-gray-200"
          >
            Ingredients
          </h2>
          <FieldArray name="ingredients" v-slot="{ fields, push, remove }">
            <div class="space-y-4">
              <div
                v-for="(field, idx) in fields"
                :key="field.key"
                class="flex flex-col md:flex-row gap-4"
              >
                <UIInput
                  :name="`ingredients[${idx}].name`"
                  placeholder="Ingredient Name"
                  label="Ingredient Name"
                  :error-message="errors.ingredients?.[idx]?.name"
                  :is-required="true"
                  class="flex-1"
                />
                <UIInput
                  :name="`ingredients[${idx}].quantity`"
                  placeholder="Ingredient Quantity"
                  label="Ingredient Quantity"
                  :error-message="errors.ingredients?.[idx]?.quantity"
                  :is-required="true"
                  class="flex-1"
                />
                <button
                  type="button"
                  @click="remove(idx)"
                  class="text-red-600 hover:text-red-800 rounded-lg transition-colors"
                >
                  Remove
                </button>
              </div>
              <button
                type="button"
                @click="push({ name: '', quantity: '' })"
                class="bg-green-100 w-full text-green-700 px-4 py-2 rounded-lg hover:bg-green-200 transition-colors"
              >
                + Add Ingredient
              </button>
            </div>
          </FieldArray>
        </div>

        <!-- Steps -->
        <div>
          <h2
            class="block pl-3 ml-px text-sm font-medium text-gray-700 dark:text-gray-200"
          >
            Steps
          </h2>
          <FieldArray name="steps" v-slot="{ fields, push, remove }">
            <div class="space-y-4">
              <div
                v-for="(field, idx) in fields"
                :key="field.key"
                class="flex flex-col md:flex-row gap-4"
              >
                <UIInput
                  :name="`steps[${idx}].step_number`"
                  placeholder="Step Number"
                  label="Step Number"
                  :error-message="errors.steps?.[idx]?.step_number"
                  :is-required="true"
                  type="number"
                  class="flex-1"
                />
                <UIInput
                  :name="`steps[${idx}].instruction`"
                  placeholder="Step Instruction"
                  label="Step Instruction"
                  :error-message="errors.steps?.[idx]?.instruction"
                  :is-required="true"
                  class="flex-1"
                />
                <button
                  type="button"
                  @click="remove(idx)"
                  class="text-red-600 hover:text-red-800 transition-colors"
                >
                  Remove
                </button>
              </div>
              <button
                type="button"
                @click="push({ step_number: '', instruction: '' })"
                class="bg-green-100 w-full text-green-700 px-4 py-2 rounded-lg hover:bg-green-200 transition-colors"
              >
                + Add Step
              </button>
            </div>
          </FieldArray>
        </div>

        <div class="w-full">
          <h2
            class="block pl-3 ml-px text-sm font-medium text-gray-700 dark:text-gray-200"
          >
            Images
          </h2>
          <UIMultipleImages
            name="images"
            @image-changed="handleFileChange"
            :isMultiple="true"
            :error-message="errors.images"
          />
        </div>

        <!-- Submit Button -->
        <div class="flex justify-center">
          <button
            type="submit"
            class="bg-green-600 text-white px-8 py-3 hover:bg-green-700 transition-colors rounded-lg shadow-md"
          >
            <span v-if="!recipeStore.$state.isLoading"> Create Recipe </span>
            <div v-else class="flex items-center gap-2">
              <UILoading />
            </div>
          </button>
        </div>
      </Form>
    </div>
  </div>
</template>
