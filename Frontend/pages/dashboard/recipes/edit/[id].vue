<template>
  <div class="max-w-3xl mx-auto p-4 md:p-6">
    <h1 class="text-3xl font-bold mb-6 text-gray-800">Edit Recipe</h1>

    <Form :validation-schema="schema" @submit="submit" class="bg-white rounded-lg shadow-md p-6">
      <div class="space-y-6">
        <!-- Basic Information Section -->
        <div class="border-b pb-6">
          <h2 class="text-xl font-semibold mb-4 text-gray-700">Basic Information</h2>
          
          <Field name="title" v-slot="{ field, errorMessage }">
            <label class="block font-medium text-gray-700 mb-1">Title*</label>
            <input 
              v-bind="field" 
              v-model="form.title"
              type="text" 
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500" 
              placeholder="e.g. Grandma's Chocolate Chip Cookies"
            />
            <span v-if="errorMessage" class="text-red-500 text-sm mt-1 block">{{ errorMessage }}</span>
          </Field>

          <Field name="description" v-slot="{ field, errorMessage }" class="mt-4">
            <label class="block font-medium text-gray-700 mb-1">Description</label>
            <textarea 
              v-bind="field" 
              v-model="form.description"
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500" 
              rows="3"
              placeholder="Tell us about your recipe..."
            ></textarea>
            <span v-if="errorMessage" class="text-red-500 text-sm mt-1 block">{{ errorMessage }}</span>
          </Field>
        </div>

        <!-- Time & Servings Section -->
        <div class="border-b pb-8">
          <h2 class="text-xl font-semibold mb-6 text-gray-800">Time & Servings</h2>
          
          <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
            <!-- Prep Time -->
            <div class="space-y-1">
              <Field name="prep_time" v-slot="{ field, errors, meta }"
              v-model="form.prep_time">
                <label class="block text-sm font-medium text-gray-700">
                  Prep Time (minutes)
                  <span class="text-gray-400 ml-1">optional</span>
                </label>
                <div class="relative mt-1">
                  <input 
                    v-bind="field"
                    type="number"
                    class="block w-full rounded-md border-gray-300 shadow-sm py-2.5 px-3 focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                    placeholder="15"
                    min="0"
                    :class="{
                      'border-red-300 text-red-900 placeholder-red-300 focus:border-red-500 focus:ring-red-500': errors.length,
                      'border-gray-300': !errors.length
                    }"
                  >
                  <div class="absolute inset-y-0 right-0 flex items-center pr-3 pointer-events-none">
                    <span class="text-gray-500 sm:text-sm">min</span>
                  </div>
                </div>
                <p v-if="errors.length" class="mt-1 text-sm text-red-600">
                  {{ errors[0] }}
                </p>
              </Field>
            </div>

            <!-- Cook Time -->
            <div class="space-y-1">
              <Field name="cook_time" v-slot="{ field, errors, meta }"
              v-model="form.cook_time">
                <label class="block text-sm font-medium text-gray-700">
                  Cook Time (minutes)
                  <span class="text-gray-400 ml-1">optional</span>
                </label>
                <div class="relative mt-1">
                  <input 
                    v-bind="field"
                    type="number"
                    class="block w-full rounded-md border-gray-300 shadow-sm py-2.5 px-3 focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                    placeholder="30"
                    min="0"
                    :class="{
                      'border-red-300 text-red-900 placeholder-red-300 focus:border-red-500 focus:ring-red-500': errors.length,
                      'border-gray-300': !errors.length
                    }"
                  >
                  <div class="absolute inset-y-0 right-0 flex items-center pr-3 pointer-events-none">
                    <span class="text-gray-500 sm:text-sm">min</span>
                  </div>
                </div>
                <p v-if="errors.length" class="mt-1 text-sm text-red-600">
                  {{ errors[0] }}
                </p>
              </Field>
            </div>

            <!-- Servings -->
            <div class="space-y-1">
              <Field name="servings" v-slot="{ field, errors, meta }"
              v-model="form.servings">
                <label class="block text-sm font-medium text-gray-700">
                  Servings
                  <span class="text-red-500">*</span>
                </label>
                <div class="relative mt-1">
                  <input 
                    v-bind="field"
                    type="number"
                    class="block w-full rounded-md border-gray-300 shadow-sm py-2.5 px-3 focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                    placeholder="4"
                    min="1"
                    required
                    :class="{
                      'border-red-300 text-red-900 placeholder-red-300 focus:border-red-500 focus:ring-red-500': errors.length,
                      'border-gray-300': !errors.length
                    }"
                  >
                  <div class="absolute inset-y-0 right-0 flex items-center pr-3 pointer-events-none">
                    <span class="text-gray-500 sm:text-sm">people</span>
                  </div>
                </div>
                <p v-if="errors.length" class="mt-1 text-sm text-red-600">
                  {{ errors[0] }}
                </p>
              </Field>
            </div>
          </div>
        </div>

        <!-- Category Section -->
        <div class="border-b pb-6">
          <Field name="category_id" v-slot="{ field, errorMessage }"
          >
            <label class="block font-medium text-gray-700 mb-1">Category*</label>
            <select 
              v-bind="field" 
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500" 
              :disabled="categoriesLoading"
              :class="{ 'opacity-50': categoriesLoading }"
            >
              <option disabled value="">
                {{ categoriesLoading ? 'Loading categories...' : 'Select a category' }}
              </option>
              <option v-for="cat in categories" :key="cat.id" :value="cat.id">
                {{ cat.name }}
              </option>
            </select>
            <span v-if="errorMessage" class="text-red-500 text-sm mt-1 block">{{ errorMessage }}</span>
          </Field>
        </div>

        <!-- Steps Section -->
        <div class="border-b pb-6">
          <div class="flex justify-between items-center mb-4">
            <h2 class="text-xl font-semibold text-gray-700">Steps</h2>
            <button 
              type="button" 
              @click="addStep" 
              class="px-3 py-1 border border-gray-300 text-gray-700 font-medium rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 text-sm flex items-center"
            >
              <PlusIcon class="w-4 h-4 mr-1" />
              Add Step
            </button>
          </div>
          
          <div v-for="(step, index) in steps" :key="index" class="mb-4 p-4 border rounded-lg bg-gray-50">
            <div class="flex justify-between items-center mb-2">
              <label class="font-medium text-gray-700">Step {{ index + 1 }}</label>
              <button 
                v-if="steps.length > 1" 
                type="button" 
                @click="removeStep(index)"
                class="text-red-500 hover:text-red-700 text-sm flex items-center"
              >
                <TrashIcon class="w-4 h-4 mr-1" />
                Remove
              </button>
            </div>
            
            <input 
              v-model.number="step.step_number" 
              type="number" 
              class="w-20 px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 mb-3" 
              min="1"
              @change="reorderSteps"
            />
            
            <textarea 
              v-model="step.instruction" 
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 mb-3" 
              rows="3"
              placeholder="Detailed instructions..."
              required
            />
            
            <input 
              v-model="step.image_url" 
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500" 
              placeholder="Optional image URL (e.g., from Imgur)"
            />
          </div>
        </div>

        <!-- Ingredients Section -->
        <div class="border-b pb-6">
          <div class="flex justify-between items-center mb-4">
            <h2 class="text-xl font-semibold text-gray-700">Ingredients</h2>
            <button 
              type="button" 
              @click="addIngredient" 
              class="px-3 py-1 border border-gray-300 text-gray-700 font-medium rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 text-sm flex items-center"
            >
              <PlusIcon class="w-4 h-4 mr-1" />
              Add Ingredient
            </button>
          </div>
          
          <div v-for="(ing, index) in ingredients" :key="index" class="mb-4 p-4 border rounded-lg bg-gray-50">
            <div class="flex justify-between items-center mb-2">
              <label class="font-medium text-gray-700">Ingredient {{ index + 1 }}</label>
              <button 
                v-if="ingredients.length > 1" 
                type="button" 
                @click="removeIngredient(index)"
                class="text-red-500 hover:text-red-700 text-sm flex items-center"
              >
                <TrashIcon class="w-4 h-4 mr-1" />
                Remove
              </button>
            </div>
            
            <input 
              v-model="ing.name" 
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 mb-3" 
              placeholder="Name* (e.g., Flour)"
              required
            />
            
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm text-gray-600 mb-1">Quantity</label>
                <input 
                  v-model.number="ing.quantity" 
                  type="number" 
                  step="0.1" 
                  class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500" 
                  min="0" 
                  placeholder="1.5"
                />
              </div>
              <div>
                <label class="block text-sm text-gray-600 mb-1">Unit</label>
                <input 
                  v-model="ing.unit" 
                  class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500" 
                  placeholder="e.g., cups, grams"
                />
              </div>
            </div>
          </div>
        </div>

        <!-- Images Section -->
        <div>
          <div class="flex justify-between items-center mb-4">
            <h2 class="text-xl font-semibold text-gray-700">Images</h2>
            <button 
              type="button" 
              @click="addImage" 
              class="px-3 py-1 border border-gray-300 text-gray-700 font-medium rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 text-sm flex items-center"
            >
              <PlusIcon class="w-4 h-4 mr-1" />
              Add Image
            </button>
          </div>
          
          <p class="text-sm text-gray-600 mb-4">Upload at least one image and mark one as featured.</p>
          
          <div v-for="(img, index) in images" :key="index" class="mb-4 p-4 border rounded-lg bg-gray-50">
            <div class="flex justify-between items-center mb-3">
              <label class="font-medium text-gray-700">Image {{ index + 1 }}</label>
              <button 
                v-if="images.length > 1" 
                type="button" 
                @click="removeImage(index)"
                class="text-red-500 hover:text-red-700 text-sm flex items-center"
              >
                <TrashIcon class="w-4 h-4 mr-1" />
                Remove
              </button>
            </div>
            
            <input 
              type="file" 
              @change="handleImageUpload($event, index)"
              accept="image/jpeg, image/png, image/webp"
              class="block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-md file:border-0 file:text-sm file:font-semibold file:bg-primary-50 file:text-primary-700 hover:file:bg-primary-100 mb-3"
            />
            
            <div v-if="img.previewUrl" class="mt-3 flex flex-col items-center">
              <img :src="img.previewUrl" class="max-h-48 rounded-md shadow" />
              <label class="inline-flex items-center mt-3 cursor-pointer">
                <input 
                  v-model="img.is_featured" 
                  type="radio" 
                  name="featuredImage"
                  class="focus:ring-primary-500 h-4 w-4 text-primary-600 border-gray-300 mr-2" 
                  :value="true"
                  @change="setFeaturedImage(index)"
                />
                <span class="ml-2 text-gray-700">Set as featured image</span>
              </label>
            </div>
          </div>
        </div>

        <!-- Submit Button -->
        <div class="pt-4">
          <button 
            type="submit" 
            class="px-4 py-2 bg-blue-600 text-white font-medium rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 mt-6 disabled:opacity-50 disabled:cursor-not-allowed"
            :disabled="loading"
          >
            {{ loading ? 'Saving...' : 'Save Changes' }}
          </button>
        </div>
      </div>
    </Form>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Form, Field } from 'vee-validate'
import * as yup from 'yup'
import { PlusIcon, TrashIcon } from '@heroicons/vue/outline'
import gql from 'graphql-tag'
import { useQuery, useMutation } from '@vue/apollo-composable'

const route = useRoute()
const router = useRouter()
const recipeId = route.params.id

// Schema validation
const schema = yup.object({
  title: yup.string().required('Title is required').max(100, 'Title must be less than 100 characters'),
  description: yup.string().max(500, 'Description must be less than 500 characters'),
  prep_time: yup.number().min(0).max(999, 'Prep time seems too long').nullable(),
  cook_time: yup.number().min(0).max(999, 'Cook time seems too long').nullable(),
  servings: yup.number().min(1).max(100, 'Servings must be less than 100').required('Servings is required'),
  category_id: yup.string().required('Category is required')
})

// Form data
const form = ref({
  title: '',
  description: '',
  prep_time: null,
  cook_time: null,
  servings: null,
  category_id: ''
})

const steps = ref([])
const ingredients = ref([])
const images = ref([])
const loading = ref(false)
const saving = ref(false)

// GraphQL Queries
const GET_RECIPE = gql`
  query GetRecipe($id: uuid!) {
    recipes_by_pk(id: $id) {
      id
      title
      description
      prep_time
      cook_time
      servings
      featured_image_url
      recipe_categories {
        category {
          id
          name
        }
      }
      recipe_ingredients {
        id
        name
        quantity
        unit
      }
      recipe_steps(order_by: {step_number: asc}) {
        id
        step_number
        instruction
        image_url
      }
      recipe_images {
        id
        image_url
        is_featured
      }
    }
  }
`

const GET_CATEGORIES = gql`
  query GetCategories {
    categories {
      id
      name
    }
  }
`

// Mutations
const UPDATE_RECIPE = gql`
  mutation UpdateRecipe(
    $id: uuid!
    $title: String!
    $description: String
    $prep_time: Int
    $cook_time: Int
    $servings: Int!
    $category_id: uuid!
  ) {
    update_recipes_by_pk(
      pk_columns: {id: $id}
      _set: {
        title: $title
        description: $description
        prep_time: $prep_time
        cook_time: $cook_time
        servings: $servings
      }
    ) {
      id
      title
      description
      prep_time
      cook_time
      servings
    }
    
    # Update category relationship
    delete_recipe_categories(where: {recipe_id: {_eq: $id}}) {
      affected_rows
    }
    
    insert_recipe_categories_one(object: {recipe_id: $id, category_id: $category_id}) {
      id
    }
  }
`

const UPDATE_INGREDIENTS = gql`
  mutation UpdateIngredients($recipeId: uuid!, $ingredients: [recipe_ingredients_insert_input!]!) {
    delete_recipe_ingredients(where: {recipe_id: {_eq: $recipeId}}) {
      affected_rows
    }
    insert_recipe_ingredients(objects: $ingredients) {
      affected_rows
    }
  }
`

const UPDATE_STEPS = gql`
  mutation UpdateSteps($recipeId: uuid!, $steps: [recipe_steps_insert_input!]!) {
    delete_recipe_steps(where: {recipe_id: {_eq: $recipeId}}) {
      affected_rows
    }
    insert_recipe_steps(objects: $steps) {
      affected_rows
    }
  }
`

const UPLOAD_IMAGE = gql`
  mutation UploadRecipeImage(
    $recipe_id: uuid!
    $is_featured: Boolean!
    $base64str: String!
    $filename: String!
  ) {
    uploadRecipeImage(
      recipe_id: $recipe_id
      is_featured: $is_featured
      base64str: $base64str
      filename: $filename
    ) {
      success
      message
      image_url
    }
  }
`

const DELETE_IMAGE = gql`
  mutation DeleteRecipeImage($id: uuid!) {
    delete_recipe_images_by_pk(id: $id) {
      id
    }
  }
`

const UPDATE_FEATURED_IMAGE = gql`
  mutation UpdateFeaturedImage($id: uuid!, $featured_image_url: String!) {
    update_recipes_by_pk(
      pk_columns: {id: $id}
      _set: {featured_image_url: $featured_image_url}
    ) {
      id
      featured_image_url
    }
  }
`

// Fetch recipe data
const { result: recipeResult, loading: recipeLoading } = useQuery(
  GET_RECIPE,
  () => ({ id: recipeId })
)

// Fetch categories
const { result: categoriesResult, loading: categoriesLoading } = useQuery(GET_CATEGORIES)

const categories = computed(() => categoriesResult.value?.categories || [])

// Initialize form with recipe data
watch(recipeResult, (result) => {
  if (result?.recipes_by_pk) {
    const recipe = result.recipes_by_pk
    
    // Basic info
    form.value = {
      title: recipe.title,
      description: recipe.description || '',
      prep_time: recipe.prep_time,
      cook_time: recipe.cook_time,
      servings: recipe.servings,
      category_id: recipe.recipe_categories[0]?.category.id || ''
    }
    
    // Steps
    steps.value = recipe.recipe_steps.map(step => ({
      id: step.id,
      step_number: step.step_number,
      instruction: step.instruction,
      image_url: step.image_url || ''
    }))
    
    // Ingredients
    ingredients.value = recipe.recipe_ingredients.map(ing => ({
      id: ing.id,
      name: ing.name,
      quantity: ing.quantity,
      unit: ing.unit || ''
    }))
    
    // Images
    images.value = recipe.recipe_images.map(img => ({
      id: img.id,
      previewUrl: img.image_url,
      is_featured: img.is_featured,
      file: null // No file since it's already uploaded
    }))
  }
})

// Mutations
const { mutate: updateRecipe } = useMutation(UPDATE_RECIPE)
const { mutate: updateIngredients } = useMutation(UPDATE_INGREDIENTS)
const { mutate: updateSteps } = useMutation(UPDATE_STEPS)
const { mutate: uploadImage } = useMutation(UPLOAD_IMAGE)
const { mutate: deleteImage } = useMutation(DELETE_IMAGE)
const { mutate: updateFeaturedImage } = useMutation(UPDATE_FEATURED_IMAGE)

// Helper methods
const addStep = () => {
  steps.value.push({ 
    id: null,
    step_number: steps.value.length + 1, 
    instruction: '', 
    image_url: '' 
  })
}

const removeStep = (index) => {
  steps.value.splice(index, 1)
  // Re-number remaining steps
  steps.value.forEach((step, i) => {
    step.step_number = i + 1
  })
}

const reorderSteps = () => {
  // Sort steps by step_number
  steps.value.sort((a, b) => a.step_number - b.step_number)
}

const addIngredient = () => {
  ingredients.value.push({ 
    id: null,
    name: '', 
    quantity: null, 
    unit: '' 
  })
}

const removeIngredient = (index) => {
  ingredients.value.splice(index, 1)
}

const addImage = () => {
  images.value.push({ 
    id: null,
    file: null, 
    previewUrl: '', 
    is_featured: false 
  })
}

const removeImage = async (index) => {
  const image = images.value[index]
  
  // If this is an existing image (has ID), delete it from server
  if (image.id) {
    try {
      await deleteImage({ id: image.id })
    } catch (err) {
      console.error('Error deleting image:', err)
      alert('Failed to delete image. Please try again.')
      return
    }
  }
  
  const wasFeatured = image.is_featured
  images.value.splice(index, 1)
  
  // If we removed the featured image and there are images left
  if (wasFeatured && images.value.length > 0) {
    // Set the first remaining image as featured
    images.value[0].is_featured = true
    
    // Update featured image in database if it's an existing image
    if (images.value[0].id) {
      try {
        await updateFeaturedImage({
          id: recipeId,
          featured_image_url: images.value[0].previewUrl
        })
      } catch (err) {
        console.error('Error updating featured image:', err)
      }
    }
  }
}

const setFeaturedImage = async (index) => {
  images.value.forEach((img, i) => {
    img.is_featured = i === index
  })
  
  // Update featured image in database
  if (images.value[index].id || images.value[index].previewUrl) {
    try {
      await updateFeaturedImage({
        id: recipeId,
        featured_image_url: images.value[index].previewUrl
      })
    } catch (err) {
      console.error('Error updating featured image:', err)
    }
  }
}

const handleImageUpload = (event, index) => {
  const file = event.target.files[0]
  if (!file) return

  const validTypes = ['image/jpeg', 'image/png', 'image/webp']
  if (!validTypes.includes(file.type)) {
    alert('Only JPG, PNG, or WEBP images are allowed')
    event.target.value = '' // Reset file input
    return
  }

  if (file.size > 5 * 1024 * 1024) {
    alert('Image size should not exceed 5MB')
    event.target.value = '' // Reset file input
    return
  }

  const reader = new FileReader()
  reader.onload = (e) => {
    images.value[index].file = file
    images.value[index].previewUrl = e.target.result
  }
  reader.onerror = () => {
    alert('Error reading image file')
    event.target.value = '' // Reset file input
  }
  reader.readAsDataURL(file)
}

// Base64 helper
const toBase64 = (file) => new Promise((resolve, reject) => {
  const reader = new FileReader()
  reader.readAsDataURL(file)
  reader.onload = () => resolve(reader.result.split(',')[1])
  reader.onerror = error => reject(error)
})

// Submit handler
const submit = async (values) => {
  try {
    loading.value = true

    // Validate steps
    if (steps.value.some(step => !step.instruction.trim())) {
      throw new Error('All steps must have instructions')
    }

    // Validate ingredients
    if (ingredients.value.some(ing => !ing.name.trim())) {
      throw new Error('All ingredients must have names')
    }

    // Validate images
    const hasImages = images.value.some(img => img.previewUrl)
    if (!hasImages) {
      throw new Error('Please add at least one image')
    }

    // Check if at least one image is marked as featured
    const hasFeaturedImage = images.value.some(img => img.is_featured)
    if (!hasFeaturedImage) {
      throw new Error('Please mark one image as featured')
    }

    // 1. Update basic recipe info and category
    const { errors: updateErrors } = await updateRecipe({
      id: recipeId,
      title: values.title,
      description: values.description || null,
      prep_time: values.prep_time || null,
      cook_time: values.cook_time || null,
      servings: values.servings,
      category_id: values.category_id
    })
    
    if (updateErrors) {
      throw new Error('Failed to update recipe: ' + updateErrors[0].message)
    }

    // 2. Update ingredients
    const ingredientsInput = ingredients.value.map(ing => ({
      recipe_id: recipeId,
      name: ing.name,
      quantity: ing.quantity || null,
      unit: ing.unit || null
    }))
    
    const { errors: ingredientsErrors } = await updateIngredients({
      recipeId,
      ingredients: ingredientsInput
    })
    
    if (ingredientsErrors) {
      throw new Error('Failed to update ingredients: ' + ingredientsErrors[0].message)
    }

    // 3. Update steps
    const stepsInput = steps.value.map(step => ({
      recipe_id: recipeId,
      step_number: step.step_number,
      instruction: step.instruction,
      image_url: step.image_url || null
    }))
    
    const { errors: stepsErrors } = await updateSteps({
      recipeId,
      steps: stepsInput
    })
    
    if (stepsErrors) {
      throw new Error('Failed to update steps: ' + stepsErrors[0].message)
    }

    // 4. Handle images
    let featuredImageUrl = null
    
    for (const [index, img] of images.value.entries()) {
      // Skip if this is an existing image that wasn't changed
      if (img.id && !img.file) {
        if (img.is_featured) {
          featuredImageUrl = img.previewUrl
        }
        continue
      }
      
      // Skip if no file to upload
      if (!img.file) continue
      
      try {
        const base64str = await toBase64(img.file)
        const filename = `${Date.now()}_${img.file.name.replace(/\s+/g, '_')}`

        const { data: uploadData, errors: uploadErrors } = await uploadImage({
          recipe_id: recipeId,
          is_featured: img.is_featured,
          base64str,
          filename
        })

        if (uploadErrors) {
          console.error(`Error uploading image ${index + 1}:`, uploadErrors[0].message)
          continue
        }

        if (uploadData?.uploadRecipeImage?.image_url) {
          if (img.is_featured) {
            featuredImageUrl = uploadData.uploadRecipeImage.image_url
          }
        }
      } catch (err) {
        console.error(`Error uploading image ${index + 1}:`, err)
      }
    }

    // 5. Update featured image if needed
    if (featuredImageUrl) {
      try {
        await updateFeaturedImage({
          id: recipeId,
          featured_image_url: featuredImageUrl
        })
      } catch (err) {
        console.error('Error updating featured image:', err)
      }
    }

    // Success!
    router.push({
      path: `/recipes/${recipeId}`,
      query: { updated: 'true' }
    })

  } catch (err) {
    alert(err.message || 'An error occurred while updating the recipe.')
    console.error('Recipe update error:', err)
  } finally {
    loading.value = false
  }
}
</script>