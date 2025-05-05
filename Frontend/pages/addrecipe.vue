<template>
  <div class="max-w-3xl mx-auto p-6">
    <h1 class="text-3xl font-bold mb-6">Create Recipe</h1>

    <Form :validation-schema="schema" @submit="submit">
      <div class="space-y-4">
        <Field name="title" v-slot="{ field, errorMessage }">
          <label class="block font-medium">Title*</label>
          <input v-bind="field" type="text" class="input" />
          <span class="text-red-500 text-sm">{{ errorMessage }}</span>
        </Field>

        <Field name="description" v-slot="{ field, errorMessage }">
          <label class="block font-medium">Description</label>
          <textarea v-bind="field" class="input"></textarea>
          <span class="text-red-500 text-sm">{{ errorMessage }}</span>
        </Field>

        <div class="grid grid-cols-3 gap-4">
          <Field name="prep_time" v-slot="{ field, errorMessage }">
            <label class="block font-medium">Prep Time (min)</label>
            <input v-bind="field" type="number" class="input" min="0" />
            <span class="text-red-500 text-sm">{{ errorMessage }}</span>
          </Field>
          <Field name="cook_time" v-slot="{ field, errorMessage }">
            <label class="block font-medium">Cook Time (min)</label>
            <input v-bind="field" type="number" class="input" min="0" />
            <span class="text-red-500 text-sm">{{ errorMessage }}</span>
          </Field>
          <Field name="servings" v-slot="{ field, errorMessage }">
            <label class="block font-medium">Servings</label>
            <input v-bind="field" type="number" class="input" min="1" />
            <span class="text-red-500 text-sm">{{ errorMessage }}</span>
          </Field>
        </div>

        <Field name="category_id" v-slot="{ field, errorMessage }">
          <label class="block font-medium">Category*</label>
          <select v-bind="field" class="input" :disabled="categoriesLoading">
            <option disabled value="">
              {{ categoriesLoading ? 'Loading categories...' : 'Select a category' }}
            </option>
            <option v-for="cat in categories" :key="cat.id" :value="cat.id">
              {{ cat.name }}
            </option>
          </select>
          <span class="text-red-500 text-sm">{{ errorMessage }}</span>
        </Field>

        <!-- Steps -->
        <div>
          <label class="block font-bold mb-2">Steps</label>
          <div v-for="(step, index) in steps" :key="index" class="mb-4 p-4 border rounded">
            <input v-model.number="step.step_number" type="number" class="input mb-2" placeholder="Step Number" min="1" />
            <textarea v-model="step.instruction" class="input mb-2" placeholder="Instruction" required />
            <input v-model="step.image_url" class="input" placeholder="Optional Image URL" />
            <button 
              v-if="steps.length > 1" 
              type="button" 
              @click="removeStep(index)"
              class="text-red-500 text-sm mt-1"
            >
              Remove Step
            </button>
          </div>
          <button type="button" @click="addStep" class="btn mt-2">Add Step</button>
        </div>

        <!-- Ingredients -->
        <div>
          <label class="block font-bold mb-2">Ingredients</label>
          <div v-for="(ing, index) in ingredients" :key="index" class="mb-4 p-4 border rounded">
            <input v-model="ing.name" class="input mb-2" placeholder="Name*" required />
            <input v-model.number="ing.quantity" type="number" step="0.1" class="input mb-2" placeholder="Quantity" min="0" />
            <input v-model="ing.unit" class="input" placeholder="Unit (e.g., grams)" />
            <button 
              v-if="ingredients.length > 1" 
              type="button" 
              @click="removeIngredient(index)"
              class="text-red-500 text-sm mt-1"
            >
              Remove Ingredient
            </button>
          </div>
          <button type="button" @click="addIngredient" class="btn mt-2">Add Ingredient</button>
        </div>

        <!-- Images -->
        <div>
          <label class="block font-bold mb-2">Images</label>
          <p class="text-sm text-gray-600 mb-2">Mark one image as featured</p>
          
          <div v-for="(img, index) in images" :key="index" class="mb-4 p-4 border rounded">
            <input 
              type="file" 
              @change="handleImageUpload($event, index)"
              accept="image/*"
              class="input mb-2"
            />
            <div v-if="img.previewUrl" class="mt-2">
              <img :src="img.previewUrl" class="max-h-32" />
            </div>
            <label class="inline-flex items-center mt-2">
              <input 
                v-model="img.is_featured" 
                type="radio" 
                name="featuredImage"
                class="mr-2" 
                :value="true"
                @change="setFeaturedImage(index)"
              />
              Featured
            </label>
            <button 
              v-if="images.length > 1" 
              type="button" 
              @click="removeImage(index)"
              class="text-red-500 text-sm mt-1 block"
            >
              Remove Image
            </button>
          </div>
          <button type="button" @click="addImage" class="btn mt-2">Add Image</button>
        </div>

        <button type="submit" class="btn btn-primary mt-6" :disabled="loading">
          {{ loading ? 'Submitting...' : 'Submit Recipe' }}
        </button>
      </div>
    </Form>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { Form, Field } from 'vee-validate'
import * as yup from 'yup'
import { useRouter } from 'vue-router'
import gql from 'graphql-tag'
import { useMutation, useQuery } from '@vue/apollo-composable'

const router = useRouter()

// Schema validation
const schema = yup.object({
  title: yup.string().required('Title is required'),
  description: yup.string(),
  prep_time: yup.number().min(0).nullable(),
  cook_time: yup.number().min(0).nullable(),
  servings: yup.number().min(1).nullable(),
  category_id: yup.string().required('Category is required')
})

// Form data
const steps = ref([{ step_number: 1, instruction: '', image_url: '' }])
const ingredients = ref([{ name: '', quantity: null, unit: '' }])
const images = ref([{ file: null, previewUrl: '', is_featured: true }])

// GraphQL Queries
const CREATE_RECIPE_MUTATION = gql`
  mutation CreateRecipe($object: CreateRecipeInput!) {
    createRecipe(object: $object) {
      id
      title
      total_time
      featured_image_url
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

const UPDATE_RECIPE = gql`
  mutation UpdateRecipe($id: uuid!, $featured_image_url: String) {
    update_recipes_by_pk(
      pk_columns: {id: $id}
      _set: {featured_image_url: $featured_image_url}
    ) {
      id
      featured_image_url
    }
  }
`

// Mutations
const { mutate: createRecipe, loading: creatingRecipe } = useMutation(CREATE_RECIPE_MUTATION)
const { mutate: uploadImage, loading: uploadingImage } = useMutation(UPLOAD_IMAGE)
const { mutate: updateRecipe, loading: updatingRecipe } = useMutation(UPDATE_RECIPE)

const loading = ref(false)

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
  const userStr = localStorage.getItem("user");
  if (userStr) {
    const user = JSON.parse(userStr);
    userId.value = user.id
  }
})
console.log(userId)

const handleImageUpload = (event, index) => {
  const file = event.target.files[0]
  if (!file) return

  const validTypes = ['image/jpeg', 'image/png', 'image/webp']
  if (!validTypes.includes(file.type)) {
    alert('Only JPG, PNG, or WEBP images are allowed')
    return
  }

  if (file.size > 5 * 1024 * 1024) {
    alert('Image size should not exceed 5MB')
    return
  }

  const reader = new FileReader()
  reader.onload = (e) => {
    images.value[index].file = file
    images.value[index].previewUrl = e.target.result
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
  if (images.value[index].is_featured && images.value.length > 1) {
    const newIndex = index === 0 ? 1 : 0
    images.value[newIndex].is_featured = true
  }
  images.value.splice(index, 1)
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

    // Validate form data
    if (steps.value.some(step => !step.instruction.trim())) {
      throw new Error('All steps must have instructions')
    }
    if (ingredients.value.some(ing => !ing.name.trim())) {
      throw new Error('All ingredients must have names')
    }
    if (images.value.length === 0) {
      throw new Error('Please add at least one image')
    }
    if (images.value.some(img => !img.file)) {
      throw new Error('Please select an image file for all images')
    }
    
    const recipePayload = {
      title: values.title,
      description: values.description || null,
      user_id: userId.value,
      prep_time: values.prep_time || null,
      cook_time: values.cook_time || null,
      servings: values.servings || null,
      category_ids: [values.category_id],
      featured_image_url: null, // Will be set after upload
      steps: steps.value.map(step => ({
        step_number: step.step_number,
        instruction: step.instruction,
        image_url: step.image_url || null
      })),
      ingredients: ingredients.value.map(ing => ({
        name: ing.name,
        quantity: ing.quantity || null,
        unit: ing.unit || null
      })),
      images: [] // Will be populated after upload
    }

    const { data: createData } = await createRecipe({ object: recipePayload })
    const createdRecipe = createData?.createRecipe
    
    if (!createdRecipe?.id) {
      throw new Error('Failed to create recipe')
    }
    const uploadImages = async (recipeId) => {
  const uploadedImages = []
  
  for (const [index, img] of images.value.entries()) {
    if (!img.file) continue
    
    try {
      const base64str = await toBase64(img.file)
      const filename = `${Date.now()}_${img.file.name.replace(/\s+/g, '_')}`
      
      const { data } = await uploadImage({
        recipe_id: recipeId,
        is_featured: img.is_featured || false,
        base64str,
        filename
      })
      
      if (data?.uploadRecipeImage?.image_url) {
        uploadedImages.push({
          image_url: data.uploadRecipeImage.image_url,
          is_featured: img.is_featured || false
        })
      }
    } catch (err) {
      console.error(`Error uploading image ${index + 1}:`, err)
      throw new Error(`Failed to upload image ${index + 1}`)
    }
  }
  
  return uploadedImages
}
    // 2. Upload all images with the recipe's ID
    const uploadedImages = []
    let featuredImageUrl = null

    for (const [index, img] of images.value.entries()) {
      const base64str = await toBase64(img.file)
      const filename = `${Date.now()}_${img.file.name.replace(/\s+/g, '_')}`

      const { data } = await uploadImage({
        recipe_id: createdRecipe.id,
        is_featured: img.is_featured,
        base64str,
        filename
      })

      const imageUrl = data?.uploadRecipeImage?.image_url
      if (imageUrl) {
        uploadedImages.push({
          image_url: imageUrl,
          is_featured: img.is_featured
        })

        if (img.is_featured) {
          featuredImageUrl = imageUrl
        }
      }
    }

    // 3. Update recipe with the featured image URL
    if (featuredImageUrl) {
      await updateRecipe({
        id: createdRecipe.id,
        featured_image_url: featuredImageUrl
      })
    }
 alert(`Recipe created successfully! ID: ${createdRecipe.id}`)
    router.push(`/recipes/${createdRecipe.id}`)

  } catch (err) {
    alert(err.message || 'An error occurred while submitting the recipe.')
    console.error(err)
  } finally {
    loading.value = false
  }
}

</script>