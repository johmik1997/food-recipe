<template>
    <div class="min-h-screen bg-gray-50">
    <Navbar />
   
  <div class="max-w-3xl mx-auto p-4 md:p-6">
    <h1 class="text-3xl font-bold mb-6 text-gray-800">
      <span class="bg-gradient-to-r from-green-600 to-green-700 bg-clip-text text-transparent">
        Create New Recipe
      </span>
    </h1>

    <Form :validation-schema="schema" @submit="submit" class="bg-white rounded-xl shadow-sm p-6 border border-gray-100">
      <div class="space-y-6">
        <!-- Basic Information Section -->
        <div class="border-b border-gray-200 pb-6">
          <h2 class="text-xl font-semibold mb-4 text-gray-700">Basic Information</h2>
          
          <Field name="title" v-slot="{ field, errorMessage }">
            <label class="block font-medium text-gray-700 mb-1">Title*</label>
            <input 
              v-bind="field" 
              type="text" 
              class="w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-green-500" 
              placeholder="e.g. Grandma's Chocolate Chip Cookies"
            />
            <span v-if="errorMessage" class="text-red-500 text-sm mt-1 block">{{ errorMessage }}</span>
          </Field>

          <Field name="description" v-slot="{ field, errorMessage }" class="mt-4">
            <label class="block font-medium text-gray-700 mb-1">Description</label>
            <textarea 
              v-bind="field" 
              class="w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-green-500" 
              rows="3"
              placeholder="Tell us about your recipe..."
            ></textarea>
            <span v-if="errorMessage" class="text-red-500 text-sm mt-1 block">{{ errorMessage }}</span>
          </Field>
        </div>

        <!-- Time & Servings Section -->
        <div class="border-b border-gray-200 pb-8">
          <h2 class="text-xl font-semibold mb-6 text-gray-800">Time & Servings</h2>
          
          <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
            <!-- Prep Time -->
            <div class="space-y-1">
              <Field name="prep_time" v-slot="{ field, errors, meta }">
                <label class="block text-sm font-medium text-gray-700">
                  Prep Time (minutes)
                  <span class="text-gray-400 ml-1">optional</span>
                </label>
                <div class="relative mt-1">
                  <input 
                    v-bind="field"
                    type="number"
                    class="block w-full rounded-lg border-gray-300 shadow-sm py-2.5 px-3 focus:border-green-500 focus:ring-green-500 sm:text-sm"
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
              <Field name="cook_time" v-slot="{ field, errors, meta }">
                <label class="block text-sm font-medium text-gray-700">
                  Cook Time (minutes)
                  <span class="text-gray-400 ml-1">optional</span>
                </label>
                <div class="relative mt-1">
                  <input 
                    v-bind="field"
                    type="number"
                    class="block w-full rounded-lg border-gray-300 shadow-sm py-2.5 px-3 focus:border-green-500 focus:ring-green-500 sm:text-sm"
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
              <Field name="servings" v-slot="{ field, errors, meta }">
                <label class="block text-sm font-medium text-gray-700">
                  Servings
                  <span class="text-red-500">*</span>
                </label>
                <div class="relative mt-1">
                  <input 
                    v-bind="field"
                    type="number"
                    class="block w-full rounded-lg border-gray-300 shadow-sm py-2.5 px-3 focus:border-green-500 focus:ring-green-500 sm:text-sm"
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

          <!-- Price Field -->
          <div class="mt-6 space-y-1">
            <Field name="price" v-slot="{ field, errors, meta }">
              <label class="block text-sm font-medium text-gray-700">
                Price (Birr)
                <span class="text-red-500">*</span>
              </label>
              <div class="relative mt-1">
                <input 
                  v-bind="field"
                  type="number"
                  class="block w-full rounded-lg border-gray-300 shadow-sm py-2.5 px-3 focus:border-green-500 focus:ring-green-500 sm:text-sm"
                  placeholder="50"
                  min="0"
                  required
                  :class="{
                    'border-red-300 text-red-900 placeholder-red-300 focus:border-red-500 focus:ring-red-500': errors.length,
                    'border-gray-300': !errors.length
                  }"
                >
                <div class="absolute inset-y-0 right-0 flex items-center pr-3 pointer-events-none">
                  <span class="text-gray-500 sm:text-sm">birr</span>
                </div>
              </div>
              <p v-if="errors.length" class="mt-1 text-sm text-red-600">
                {{ errors[0] }}
              </p>
            </Field>
          </div>
        </div>

        <!-- Category Section -->
        <div class="border-b border-gray-200 pb-6">
          <Field name="category_id" v-slot="{ field, errorMessage }">
            <label class="block font-medium text-gray-700 mb-1">Category*</label>
            <select 
              v-bind="field" 
              class="w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-green-500" 
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
        <div class="border-b border-gray-200 pb-6">
          <div class="flex justify-between items-center mb-4">
            <h2 class="text-xl font-semibold text-gray-700">Steps</h2>
            <button 
              type="button" 
              @click="addStep" 
              class="px-3 py-1 border border-gray-300 text-gray-700 font-medium rounded-lg hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 text-sm flex items-center"
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
              class="w-20 px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-green-500 mb-3" 
              min="1"
              @change="reorderSteps"
            />
            
            <textarea 
              v-model="step.instruction" 
              class="w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-green-500 mb-3" 
              rows="3"
              placeholder="Detailed instructions..."
              required
            />
            
            <input 
              v-model="step.image_url" 
              class="w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-green-500" 
              placeholder="Optional image URL (e.g., from Imgur)"
            />
          </div>
        </div>

        <!-- Ingredients Section -->
        <div class="border-b border-gray-200 pb-6">
          <div class="flex justify-between items-center mb-4">
            <h2 class="text-xl font-semibold text-gray-700">Ingredients</h2>
            <button 
              type="button" 
              @click="addIngredient" 
              class="px-3 py-1 border border-gray-300 text-gray-700 font-medium rounded-lg hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 text-sm flex items-center"
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
              class="w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-green-500 mb-3" 
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
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-green-500" 
                  min="0" 
                  placeholder="1.5"
                />
              </div>
              <div>
                <label class="block text-sm text-gray-600 mb-1">Unit</label>
                <input 
                  v-model="ing.unit" 
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-green-500" 
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
              class="px-3 py-1 border border-gray-300 text-gray-700 font-medium rounded-lg hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 text-sm flex items-center"
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
              class="block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-lg file:border-0 file:text-sm file:font-semibold file:bg-green-50 file:text-green-700 hover:file:bg-green-100 mb-3"
            />
            
            <div v-if="img.previewUrl" class="mt-3 flex flex-col items-center">
              <img :src="img.previewUrl" class="max-h-48 rounded-lg shadow" />
              <label class="inline-flex items-center mt-3 cursor-pointer">
                <input 
                  v-model="img.is_featured" 
                  type="radio" 
                  name="featuredImage"
                  class="focus:ring-green-500 h-4 w-4 text-green-600 border-gray-300 mr-2" 
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
            class="w-full px-6 py-3 bg-gradient-to-r from-green-600 to-green-700 text-white font-medium rounded-lg hover:from-green-700 hover:to-green-800 focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2 mt-6 disabled:opacity-50 disabled:cursor-not-allowed"
            :disabled="loading"
          >
            <span v-if="loading" class="flex items-center justify-center">
              <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Creating Recipe...
            </span>
            <span v-else>Create Recipe</span>
          </button>
        </div>
      </div>
    </Form>
  </div>
  <Footer/>
</div>
</template>
<script setup>
import { ref, watch, onMounted } from 'vue'
import { Form, Field } from 'vee-validate'
import * as yup from 'yup'
import { useRouter } from 'vue-router'
import gql from 'graphql-tag'
import { useMutation, useQuery } from '@vue/apollo-composable'
import { PlusIcon, TrashIcon } from '@heroicons/vue/outline'

const router = useRouter()

// Schema validation
const schema = yup.object({
  title: yup.string().required('Title is required').max(100, 'Title must be less than 100 characters'),
  description: yup.string().max(500, 'Description must be less than 500 characters'),
  prep_time: yup.number().min(0).max(999, 'Prep time seems too long').nullable(),
  cook_time: yup.number().min(0).max(999, 'Cook time seems too long').nullable(),
  servings: yup.number().min(1).max(100, 'Servings must be less than 100').required('Servings is required'),
  price: yup.number().min(0).max(9999, 'Price must be less than 9999').required('Price is required'),
  category_id: yup.string().required('Category is required')
})

// Form data
const steps = ref([{ step_number: 1, instruction: '', image_url: '' }])
const ingredients = ref([{ name: '', quantity: null, unit: '' }])
const images = ref([{ file: null, previewUrl: '', is_featured: true }])
const loading = ref(false)

// GraphQL Mutations
const INSERT_RECIPE = gql`
  mutation InsertRecipe(
    $title: String!
    $description: String
    $prep_time: Int
    $cook_time: Int
    $servings: Int!
    $price: Int!
    $user_id: uuid!
  ) {
    insert_recipes_one(object: {
      title: $title
      description: $description
      prep_time: $prep_time
      cook_time: $cook_time
      servings: $servings
      price: $price
      user_id: $user_id
    }) {
      id
      title
    }
  }
`

const INSERT_RECIPE_CATEGORY = gql`
  mutation InsertRecipeCategory($recipe_id: uuid!, $category_id: uuid!) {
    insert_recipe_categories_one(object: {
      recipe_id: $recipe_id
      category_id: $category_id
    }) {
      category_id
    }
  }
`

const INSERT_RECIPE_STEPS = gql`
  mutation InsertRecipeSteps($objects: [recipe_steps_insert_input!]!) {
    insert_recipe_steps(objects: $objects) {
      affected_rows
    }
  }
`

const INSERT_RECIPE_INGREDIENTS = gql`
  mutation InsertRecipeIngredients($objects: [recipe_ingredients_insert_input!]!) {
    insert_recipe_ingredients(objects: $objects) {
      affected_rows
    }
  }
`

const UPLOAD_RECIPE_IMAGE = gql`
  mutation UploadRecipeImage($input: UploadRecipeImageInput!) {
    uploadRecipeImage(input: $input) {
      success
      message
      image_url
    }
  }
`

// Mutations
const { mutate: insertRecipe } = useMutation(INSERT_RECIPE)
const { mutate: insertRecipeCategory } = useMutation(INSERT_RECIPE_CATEGORY)
const { mutate: insertRecipeSteps } = useMutation(INSERT_RECIPE_STEPS)
const { mutate: insertRecipeIngredients } = useMutation(INSERT_RECIPE_INGREDIENTS)
const { mutate: uploadRecipeImage } = useMutation(UPLOAD_RECIPE_IMAGE)

// Categories
const categories = ref([])
const { result, loading: categoriesLoading } = useQuery(gql`
  query GetCategories {
    categories {
      id
      name
    }
  }
`)

watch(result, (res) => {
  if (res?.categories) {
    categories.value = res.categories
  }
})

const userId = ref(null)

onMounted(() => {
  const userStr = localStorage.getItem("user")
  if (userStr) {
    const user = JSON.parse(userStr)
    userId.value = user.id
  }
})

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

const addStep = () => {
  steps.value.push({ 
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

const addIngredient = () => {
  ingredients.value.push({ name: '', quantity: null, unit: '' })
}

const removeIngredient = (index) => {
  ingredients.value.splice(index, 1)
}

const addImage = () => {
  images.value.push({ file: null, previewUrl: '', is_featured: false })
}

const removeImage = (index) => {
  const wasFeatured = images.value[index].is_featured
  images.value.splice(index, 1)
  
  if (wasFeatured && images.value.length > 0) {
    images.value[0].is_featured = true
  }
}

const setFeaturedImage = (index) => {
  images.value.forEach((img, i) => {
    img.is_featured = i === index
  })
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

    // Check featured image
    const hasFeaturedImage = images.value.some(img => img.is_featured)
    if (!hasFeaturedImage) {
      throw new Error('Please mark one image as featured')
    }

    // 1. Create the base recipe
    const { data: recipeData, errors: recipeErrors } = await insertRecipe({
      title: values.title,
      description: values.description || null,
      prep_time: values.prep_time || null,
      cook_time: values.cook_time || null,
      servings: values.servings,
      price: values.price,
      user_id: userId.value
    })
    
    if (recipeErrors) {
      throw new Error('Failed to create recipe: ' + recipeErrors[0].message)
    }

    const recipeId = recipeData?.insert_recipes_one?.id
    if (!recipeId) {
      throw new Error('Failed to create recipe - no ID returned')
    }

    // 2. Add category relationship
    const { errors: categoryErrors } = await insertRecipeCategory({
      recipe_id: recipeId,
      category_id: values.category_id
    })
    
    if (categoryErrors) {
      throw new Error('Failed to add category: ' + categoryErrors[0].message)
    }

    // 3. Add steps
    const stepsInput = steps.value.map(step => ({
      recipe_id: recipeId,
      step_number: step.step_number,
      instruction: step.instruction,
      image_url: step.image_url || null
    }))
    
    const { errors: stepsErrors } = await insertRecipeSteps({
      objects: stepsInput
    })
    
    if (stepsErrors) {
      throw new Error('Failed to add steps: ' + stepsErrors[0].message)
    }

    // 4. Add ingredients
    const ingredientsInput = ingredients.value.map(ing => ({
      recipe_id: recipeId,
      name: ing.name,
      quantity: ing.quantity || null,
      unit: ing.unit || null
    }))
    
    const { errors: ingredientsErrors } = await insertRecipeIngredients({
      objects: ingredientsInput
    })
    
    if (ingredientsErrors) {
      throw new Error('Failed to add ingredients: ' + ingredientsErrors[0].message)
    }

    // 5. Upload images (backend handles featured image)
    for (const [index, img] of images.value.entries()) {
      if (!img.file) continue
      
      try {
        const base64str = await toBase64(img.file)
        const filename = `${Date.now()}_${img.file.name.replace(/\s+/g, '_')}`

        await uploadRecipeImage({
          input: {
            recipe_id: recipeId,
            is_featured: img.is_featured,
            base64str,
            filename
          }
        })

        // Update step image URL if this image is for a step
        if (steps.value[index]?.image_url === '') {
          steps.value[index].image_url = filename // Or use the returned URL if your backend provides it
        }
      } catch (err) {
        console.error(`Error uploading image ${index + 1}:`, err)
      }
    }

    // Success!
    router.push({
      path: `/recipes/${recipeId}`,
      query: { created: 'true' }
    })

  } catch (err) {
    alert(err.message || 'An error occurred while creating the recipe.')
    console.error('Recipe creation error:', err)
  } finally {
    loading.value = false
  }
}
</script>