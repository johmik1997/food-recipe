<template>
    <div>
      <div class="mb-6">
        <h1 class="text-2xl font-bold text-gray-900">Add New Recipe</h1>
        <p class="mt-1 text-gray-600">Share your culinary masterpiece with the community</p>
      </div>
      
      <Form @submit="submitRecipe" class="space-y-6">
        <!-- Recipe Basics -->
        <div class="p-6 bg-white rounded-lg shadow">
          <h2 class="text-lg font-medium text-gray-900">Recipe Basics</h2>
          <div class="grid grid-cols-1 gap-6 mt-4 sm:grid-cols-2">
            <div>
              <label for="title" class="block text-sm font-medium text-gray-700">Recipe Title</label>
              <Field
                name="title"
                v-slot="{ field, errors }"
                rules="required|min:3|max:100"
              >
                <input
                  id="title"
                  type="text"
                  v-bind="field"
                  class="block w-full px-3 py-2 mt-1 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm"
                  :class="{ 'border-red-300': errors.length }"
                />
                <ErrorMessage name="title" class="mt-1 text-sm text-red-600" />
              </Field>
            </div>
            
            <div>
              <label for="category" class="block text-sm font-medium text-gray-700">Category</label>
              <Field
                name="category"
                v-slot="{ field, errors }"
                rules="required"
                as="select"
                class="block w-full px-3 py-2 mt-1 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm"
                :class="{ 'border-red-300': errors.length }"
              >
                <option value="">Select a category</option>
                <option v-for="category in categories" :key="category.id" :value="category.id">
                  {{ category.name }}
                </option>
              </Field>
              <ErrorMessage name="category" class="mt-1 text-sm text-red-600" />
            </div>
            
            <div>
              <label for="prepTime" class="block text-sm font-medium text-gray-700">Preparation Time (minutes)</label>
              <Field
                name="prepTime"
                v-slot="{ field, errors }"
                rules="required|numeric|min_value:1"
              >
                <input
                  id="prepTime"
                  type="number"
                  v-bind="field"
                  class="block w-full px-3 py-2 mt-1 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm"
                  :class="{ 'border-red-300': errors.length }"
                />
                <ErrorMessage name="prepTime" class="mt-1 text-sm text-red-600" />
              </Field>
            </div>
            
            <div>
              <label for="servings" class="block text-sm font-medium text-gray-700">Servings</label>
              <Field
                name="servings"
                v-slot="{ field, errors }"
                rules="required|numeric|min_value:1"
              >
                <input
                  id="servings"
                  type="number"
                  v-bind="field"
                  class="block w-full px-3 py-2 mt-1 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm"
                  :class="{ 'border-red-300': errors.length }"
                />
                <ErrorMessage name="servings" class="mt-1 text-sm text-red-600" />
              </Field>
            </div>
            
            <div class="sm:col-span-2">
              <label for="description" class="block text-sm font-medium text-gray-700">Description</label>
              <Field
                name="description"
                v-slot="{ field, errors }"
                rules="required|min:10|max:500"
              >
                <textarea
                  id="description"
                  rows="3"
                  v-bind="field"
                  class="block w-full px-3 py-2 mt-1 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm"
                  :class="{ 'border-red-300': errors.length }"
                />
                <ErrorMessage name="description" class="mt-1 text-sm text-red-600" />
              </Field>
            </div>
          </div>
        </div>
        
        <!-- Recipe Images -->
        <div class="p-6 bg-white rounded-lg shadow">
          <h2 class="text-lg font-medium text-gray-900">Recipe Images</h2>
          <div class="mt-4">
            <label class="block text-sm font-medium text-gray-700">Upload Images</label>
            <div class="flex flex-col items-center justify-center px-6 pt-5 pb-6 mt-1 border-2 border-gray-300 border-dashed rounded-md">
              <div class="space-y-1 text-center">
                <svg
                  class="w-12 h-12 mx-auto text-gray-400"
                  stroke="currentColor"
                  fill="none"
                  viewBox="0 0 48 48"
                  aria-hidden="true"
                >
                  <path
                    d="M28 8H12a4 4 0 00-4 4v20m32-12v8m0 0v8a4 4 0 01-4 4H12a4 4 0 01-4-4v-4m32-4l-3.172-3.172a4 4 0 00-5.656 0L28 28M8 32l9.172-9.172a4 4 0 015.656 0L28 28m0 0l4 4m4-24h8m-4-4v8m-12 4h.02"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  />
                </svg>
                <div class="flex text-sm text-gray-600">
                  <label
                    for="file-upload"
                    class="relative font-medium text-primary-600 bg-white rounded-md cursor-pointer hover:text-primary-500 focus-within:outline-none"
                  >
                    <span>Upload files</span>
                    <input
                      id="file-upload"
                      name="file-upload"
                      type="file"
                      class="sr-only"
                      multiple
                      @change="handleImageUpload"
                    />
                  </label>
                  <p class="pl-1">or drag and drop</p>
                </div>
                <p class="text-xs text-gray-500">PNG, JPG, GIF up to 10MB</p>
              </div>
            </div>
            
            <div v-if="uploadedImages.length > 0" class="mt-4">
              <label class="block text-sm font-medium text-gray-700">Select Featured Image</label>
              <div class="grid grid-cols-2 gap-4 mt-2 sm:grid-cols-3 md:grid-cols-4">
                <div
                  v-for="(image, index) in uploadedImages"
                  :key="index"
                  @click="setFeaturedImage(index)"
                  class="relative overflow-hidden border-2 rounded-md cursor-pointer"
                  :class="{ 'border-primary-500': featuredImageIndex === index, 'border-transparent': featuredImageIndex !== index }"
                >
                  <img :src="image.preview" class="object-cover w-full h-32" />
                  <div
                    v-if="featuredImageIndex === index"
                    class="absolute inset-0 flex items-center justify-center bg-black bg-opacity-50"
                  >
                    <span class="text-white">Featured</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- Ingredients -->
        <div class="p-6 bg-white rounded-lg shadow">
          <h2 class="text-lg font-medium text-gray-900">Ingredients</h2>
          <div class="mt-4 space-y-4">
            <div
              v-for="(ingredient, index) in ingredients"
              :key="index"
              class="flex space-x-4"
            >
              <div class="flex-1">
                <label :for="`ingredient-name-${index}`" class="block text-sm font-medium text-gray-700">Name</label>
                <Field
                  :name="`ingredients[${index}].name`"
                  v-slot="{ field }"
                  rules="required"
                >
                  <input
                    :id="`ingredient-name-${index}`"
                    type="text"
                    v-bind="field"
                    v-model="ingredient.name"
                    class="block w-full px-3 py-2 mt-1 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm"
                  />
                </Field>
              </div>
              <div class="flex-1">
                <label :for="`ingredient-amount-${index}`" class="block text-sm font-medium text-gray-700">Amount</label>
                <Field
                  :name="`ingredients[${index}].amount`"
                  v-slot="{ field }"
                  rules="required"
                >
                  <input
                    :id="`ingredient-amount-${index}`"
                    type="text"
                    v-bind="field"
                    v-model="ingredient.amount"
                    class="block w-full px-3 py-2 mt-1 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm"
                  />
                </Field>
              </div>
              <div class="flex items-end">
                <button
                  type="button"
                  @click="removeIngredient(index)"
                  class="p-2 text-gray-400 hover:text-red-500"
                >
                  <Icon name="mdi:trash-can" class="w-5 h-5" />
                </button>
              </div>
            </div>
            
            <button
              type="button"
              @click="addIngredient"
              class="inline-flex items-center px-3 py-2 text-sm font-medium leading-4 text-white bg-primary-600 border border-transparent rounded-md shadow-sm hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
            >
              <Icon name="mdi:plus" class="w-4 h-4 mr-2" />
              Add Ingredient
            </button>
          </div>
        </div>
        
        <!-- Instructions -->
        <div class="p-6 bg-white rounded-lg shadow">
          <h2 class="text-lg font-medium text-gray-900">Instructions</h2>
          <div class="mt-4 space-y-4">
            <div
              v-for="(step, index) in instructions"
              :key="index"
              class="flex space-x-4"
            >
              <div class="flex items-start pt-2">
                <span class="inline-flex items-center justify-center w-8 h-8 rounded-full bg-primary-100 text-primary-800">
                  {{ index + 1 }}
                </span>
              </div>
              <div class="flex-1">
                <label :for="`step-${index}`" class="block text-sm font-medium text-gray-700">Step {{ index + 1 }}</label>
                <Field
                  :name="`instructions[${index}].description`"
                  v-slot="{ field }"
                  rules="required"
                >
                  <textarea
                    :id="`step-${index}`"
                    rows="3"
                    v-bind="field"
                    v-model="step.description"
                    class="block w-full px-3 py-2 mt-1 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm"
                  />
                </Field>
              </div>
              <div class="flex items-start pt-2">
                <button
                  type="button"
                  @click="removeInstruction(index)"
                  class="p-2 text-gray-400 hover:text-red-500"
                >
                  <Icon name="mdi:trash-can" class="w-5 h-5" />
                </button>
              </div>
            </div>
            
            <button
              type="button"
              @click="addInstruction"
              class="inline-flex items-center px-3 py-2 text-sm font-medium leading-4 text-white bg-primary-600 border border-transparent rounded-md shadow-sm hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
            >
              <Icon name="mdi:plus" class="w-4 h-4 mr-2" />
              Add Step
            </button>
          </div>
        </div>
        
        <!-- Submit Button -->
        <div class="flex justify-end">
          <button
            type="submit"
            class="inline-flex justify-center px-4 py-2 text-sm font-medium text-white bg-primary-600 border border-transparent rounded-md shadow-sm hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
          >
            Save Recipe
          </button>
        </div>
      </Form>
    </div>
  </template>
  
  <script setup>
  import { Form, Field, ErrorMessage } from 'vee-validate';
  
  const categories = ref([
    { id: 1, name: 'Breakfast' },
    { id: 2, name: 'Lunch' },
    { id: 3, name: 'Dinner' },
    { id: 4, name: 'Dessert' },
    { id: 5, name: 'Appetizer' },
    { id: 6, name: 'Salad' },
    { id: 7, name: 'Soup' },
    { id: 8, name: 'Vegetarian' },
    { id: 9, name: 'Vegan' },
    { id: 10, name: 'Gluten-Free' }
  ]);
  
  const uploadedImages = ref([]);
  const featuredImageIndex = ref(null);
  
  const ingredients = ref([
    { name: '', amount: '' }
  ]);
  
  const instructions = ref([
    { description: '' }
  ]);
  
  const handleImageUpload = (event) => {
    const files = event.target.files;
    if (!files) return;
    
    Array.from(files).forEach(file => {
      if (!file.type.match('image.*')) return;
      
      const reader = new FileReader();
      reader.onload = (e) => {
        uploadedImages.value.push({
          file,
          preview: e.target.result
        });
      };
      reader.readAsDataURL(file);
    });
    
    if (uploadedImages.value.length > 0 && featuredImageIndex.value === null) {
      featuredImageIndex.value = 0;
    }
  };
  
  const setFeaturedImage = (index) => {
    featuredImageIndex.value = index;
  };
  
  const addIngredient = () => {
    ingredients.value.push({ name: '', amount: '' });
  };
  
  const removeIngredient = (index) => {
    if (ingredients.value.length > 1) {
      ingredients.value.splice(index, 1);
    } else {
      ingredients.value[0] = { name: '', amount: '' };
    }
  };
  
  const addInstruction = () => {
    instructions.value.push({ description: '' });
  };
  
  const removeInstruction = (index) => {
    if (instructions.value.length > 1) {
      instructions.value.splice(index, 1);
    } else {
      instructions.value[0] = { description: '' };
    }
  };
  
  const submitRecipe = async (values) => {
    if (!featuredImageIndex.value !== null) {
      alert('Please select a featured image');
      return;
    }
    
    const formData = new FormData();
    
    // Add basic recipe info
    formData.append('title', values.title);
    formData.append('category_id', values.category);
    formData.append('prep_time', values.prepTime);
    formData.append('servings', values.servings);
    formData.append('description', values.description);
    
    // Add featured image
    if (featuredImageIndex.value !== null) {
      formData.append('featured_image', uploadedImages.value[featuredImageIndex.value].file);
    }
    
    // Add other images
    uploadedImages.value.forEach((image, index) => {
      if (index !== featuredImageIndex.value) {
        formData.append('images[]', image.file);
      }
    });
    
    // Add ingredients
    ingredients.value.forEach((ingredient, index) => {
      formData.append(`ingredients[${index}][name]`, ingredient.name);
      formData.append(`ingredients[${index}][amount]`, ingredient.amount);
    });
    
    // Add instructions
    instructions.value.forEach((instruction, index) => {
      formData.append(`instructions[${index}][description]`, instruction.description);
      formData.append(`instructions[${index}][step_number]`, index + 1);
    });
    
    try {
      // Submit to Hasura action
      // const { data, error } = await useMutation(ADD_RECIPE_MUTATION, {
      //   variables: { input: formData }
      // });
      
      // if (error) throw error;
      
      alert('Recipe submitted successfully!');
      // Reset form or redirect
    } catch (error) {
      console.error('Error submitting recipe:', error);
      alert('Failed to submit recipe. Please try again.');
    }
  };
  </script>