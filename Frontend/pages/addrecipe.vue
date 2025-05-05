<template>
  <div class="max-w-3xl mx-auto p-6 bg-white rounded-lg shadow-md">
    <h1 class="text-3xl font-bold mb-6 text-indigo-700">Create Your Recipe</h1>

    <Form :validation-schema="schema" @submit="submit" class="space-y-6">
      <!-- Basic Information Section -->
      <div class="bg-indigo-50 p-4 rounded-lg">
        <h2 class="text-xl font-semibold mb-4 text-indigo-700">Basic Information</h2>
        
        <Field name="title" v-slot="{ field, errorMessage }">
          <label class="block font-medium text-gray-700">Title*</label>
          <input v-bind="field" type="text" class="form-input mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500" placeholder="e.g. Grandma's Chocolate Chip Cookies" />
          <span class="text-red-500 text-sm mt-1">{{ errorMessage }}</span>
        </Field>

        <Field name="description" v-slot="{ field, errorMessage }" class="mt-4">
          <label class="block font-medium text-gray-700">Description</label>
          <textarea v-bind="field" rows="3" class="form-textarea mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500" placeholder="Tell us about your recipe..."></textarea>
          <span class="text-red-500 text-sm mt-1">{{ errorMessage }}</span>
        </Field>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mt-4">
          <Field name="prep_time" v-slot="{ field, errorMessage }">
            <label class="block font-medium text-gray-700">Prep Time (min)</label>
            <input v-bind="field" type="number" class="form-input mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500" min="0" placeholder="15" />
            <span class="text-red-500 text-sm mt-1">{{ errorMessage }}</span>
          </Field>
          <Field name="cook_time" v-slot="{ field, errorMessage }">
            <label class="block font-medium text-gray-700">Cook Time (min)</label>
            <input v-bind="field" type="number" class="form-input mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500" min="0" placeholder="30" />
            <span class="text-red-500 text-sm mt-1">{{ errorMessage }}</span>
          </Field>
          <Field name="servings" v-slot="{ field, errorMessage }">
            <label class="block font-medium text-gray-700">Servings</label>
            <input v-bind="field" type="number" class="form-input mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500" min="1" placeholder="4" />
            <span class="text-red-500 text-sm mt-1">{{ errorMessage }}</span>
          </Field>
        </div>

        <Field name="category_id" v-slot="{ field, errorMessage }" class="mt-4">
          <label class="block font-medium text-gray-700">Category*</label>
          <select v-bind="field" class="form-select mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500" :disabled="categoriesLoading">
            <option disabled value="">
              {{ categoriesLoading ? 'Loading categories...' : 'Select a category' }}
            </option>
            <option v-for="cat in categories" :key="cat.id" :value="cat.id">
              {{ cat.name }}
            </option>
          </select>
          <span class="text-red-500 text-sm mt-1">{{ errorMessage }}</span>
        </Field>
      </div>

      <!-- Steps Section -->
      <div class="bg-indigo-50 p-4 rounded-lg">
        <h2 class="text-xl font-semibold mb-4 text-indigo-700">Cooking Steps</h2>
        
        <div v-for="(step, index) in steps" :key="index" class="mb-4 p-4 bg-white rounded-md border border-gray-200 shadow-sm">
          <div class="flex items-center mb-2">
            <span class="inline-flex items-center justify-center h-8 w-8 rounded-full bg-indigo-100 text-indigo-800 font-bold mr-3">
              {{ step.step_number }}
            </span>
            <input v-model.number="step.step_number" type="number" class="form-input w-20 rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500" placeholder="Step" min="1" />
          </div>
          <textarea v-model="step.instruction" class="form-textarea mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500" rows="3" placeholder="Detailed instruction..." required />
          <input v-model="step.image_url" class="form-input mt-2 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500" placeholder="Optional image URL (e.g., for step photos)" />
          
          <button 
            v-if="steps.length > 1" 
            type="button" 
            @click="removeStep(index)"
            class="mt-2 inline-flex items-center px-3 py-1 border border-transparent text-sm leading-4 font-medium rounded-md text-red-700 bg-red-100 hover:bg-red-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
          >
            <TrashIcon class="h-4 w-4 mr-1" />
            Remove Step
          </button>
        </div>
        
        <button 
          type="button" 
          @click="addStep" 
          class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
        >
          <PlusIcon class="h-4 w-4 mr-1" />
          Add Step
        </button>
      </div>

      <!-- Ingredients Section -->
      <div class="bg-indigo-50 p-4 rounded-lg">
        <h2 class="text-xl font-semibold mb-4 text-indigo-700">Ingredients</h2>
        
        <div v-for="(ing, index) in ingredients" :key="index" class="mb-4 p-4 bg-white rounded-md border border-gray-200 shadow-sm">
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700">Name*</label>
              <input v-model="ing.name" class="form-input mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500" placeholder="e.g. Flour" required />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700">Quantity</label>
              <input v-model.number="ing.quantity" type="number" step="0.1" class="form-input mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500" placeholder="1.5" min="0" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700">Unit</label>
              <input v-model="ing.unit" class="form-input mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500" placeholder="e.g. cups" />
            </div>
          </div>
          
          <button 
            v-if="ingredients.length > 1" 
            type="button" 
            @click="removeIngredient(index)"
            class="mt-3 inline-flex items-center px-3 py-1 border border-transparent text-sm leading-4 font-medium rounded-md text-red-700 bg-red-100 hover:bg-red-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
          >
            <TrashIcon class="h-4 w-4 mr-1" />
            Remove Ingredient
          </button>
        </div>
        
        <button 
          type="button" 
          @click="addIngredient" 
          class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
        >
          <PlusIcon class="h-4 w-4 mr-1" />
          Add Ingredient
        </button>
      </div>

      <!-- Images Section -->
      <div class="bg-indigo-50 p-4 rounded-lg">
        <h2 class="text-xl font-semibold mb-4 text-indigo-700">Images</h2>
        <p class="text-sm text-gray-600 mb-4">Upload photos of your delicious creation. Mark one as the featured image.</p>
        
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div 
            v-for="(img, index) in images" 
            :key="index" 
            class="relative p-4 bg-white rounded-md border-2 transition-all duration-200"
            :class="{'border-indigo-500': img.is_featured, 'border-gray-200': !img.is_featured}"
          >
            <div class="flex justify-between items-start mb-2">
              <label class="inline-flex items-center">
                <input 
                  v-model="img.is_featured" 
                  type="radio" 
                  name="featuredImage"
                  class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                  :value="true"
                  @change="setFeaturedImage(index)"
                />
                <span class="ml-2 text-sm text-gray-700">Featured</span>
              </label>
              
              <button 
                v-if="images.length > 1" 
                type="button" 
                @click="removeImage(index)"
                class="text-red-500 hover:text-red-700"
                title="Remove image"
              >
                <TrashIcon class="h-5 w-5" />
              </button>
            </div>
            
            <div v-if="img.previewUrl" class="mt-2">
              <img :src="img.previewUrl" class="w-full h-48 object-cover rounded-md" />
            </div>
            
            <div class="mt-4">
              <label class="block text-sm font-medium text-gray-700 mb-1">Image {{ index + 1 }}</label>
              <div class="flex items-center">
                <label class="cursor-pointer">
                  <span class="inline-flex items-center px-3 py-1 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
                    <UploadIcon class="h-4 w-4 mr-1" />
                    {{ img.file ? 'Change' : 'Select' }}
                  </span>
                  <input 
                    type="file" 
                    @change="handleImageUpload($event, index)"
                    accept="image/*"
                    class="sr-only"
                  />
                </label>
                <span v-if="img.file" class="ml-2 text-sm text-gray-500 truncate">{{ img.file.name }}</span>
              </div>
            </div>
          </div>
        </div>
        
        <button 
          type="button" 
          @click="addImage" 
          class="mt-4 inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
        >
          <PlusIcon class="h-4 w-4 mr-1" />
          Add Another Image
        </button>
      </div>

      <!-- Submit Button -->
      <div class="flex justify-end">
        <button 
          type="submit" 
          class="inline-flex items-center px-6 py-3 border border-transparent text-base font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
          :disabled="loading"
        >
          <template v-if="loading">
            <Spinner class="h-5 w-5 mr-2" />
            Creating Recipe...
          </template>
          <template v-else>
            <CheckIcon class="h-5 w-5 mr-2" />
            Submit Recipe
          </template>
        </button>
      </div>
    </Form>
  </div>
</template>

<script>
import { Form, Field } from 'vee-validate'
import * as yup from 'yup'
import gql from 'graphql-tag'
import { TrashIcon, PlusIcon, UploadIcon, CheckIcon } from '@heroicons/vue/solid'
import Spinner from '@/components/Spinner.vue'

export default {
  components: {
    Form,
    Field,
    TrashIcon,
    PlusIcon,
    UploadIcon,
    CheckIcon,
    Spinner
  },
  
  data() {
    return {
      schema: yup.object({
        title: yup.string().required('Title is required'),
        description: yup.string(),
        prep_time: yup.number().min(0).nullable(),
        cook_time: yup.number().min(0).nullable(),
        servings: yup.number().min(1).nullable(),
        category_id: yup.string().required('Category is required')
      }),
      steps: [{ step_number: 1, instruction: '', image_url: '' }],
      ingredients: [{ name: '', quantity: null, unit: '' }],
      images: [{ file: null, previewUrl: '', is_featured: true }],
      categories: [],
      categoriesLoading: false,
      userId: null,
      creatingRecipe: false,
      uploadingImage: false,
      updatingRecipe: false
    }
  },
  
  computed: {
    loading() {
      return this.creatingRecipe || this.uploadingImage || this.updatingRecipe
    }
  },
  
  mounted() {
    const userStr = localStorage.getItem("user")
    if (userStr) {
      const user = JSON.parse(userStr)
      this.userId = user.id
    }
    
    // Load categories
    this.loadCategories()
  },
  
  methods: {
    async loadCategories() {
      this.categoriesLoading = true
      try {
        const result = await this.$apollo.query({
          query: gql`
            query GetCategories {
              categories {
                id
                name
              }
            }
          `
        })
        this.categories = result.data.categories
      } catch (error) {
        console.error('Error loading categories:', error)
      } finally {
        this.categoriesLoading = false
      }
    },
    
    handleImageUpload(event, index) {
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
        this.$set(this.images[index], 'file', file)
        this.$set(this.images[index], 'previewUrl', e.target.result)
      }
      reader.readAsDataURL(file)
    },
    
    addStep() {
      this.steps.push({ 
        step_number: this.steps.length + 1, 
        instruction: '', 
        image_url: '' 
      })
    },
    
    removeStep(index) {
      this.steps.splice(index, 1)
      // Re-number remaining steps
      this.steps.forEach((step, i) => {
        step.step_number = i + 1
      })
    },
    
    addIngredient() {
      this.ingredients.push({ name: '', quantity: null, unit: '' })
    },
    
    removeIngredient(index) {
      this.ingredients.splice(index, 1)
    },
    
    addImage() {
      this.images.push({ file: null, previewUrl: '', is_featured: false })
    },
    
    removeImage(index) {
      if (this.images[index].is_featured && this.images.length > 1) {
        const newIndex = index === 0 ? 1 : 0
        this.$set(this.images[newIndex], 'is_featured', true)
      }
      this.images.splice(index, 1)
    },
    
    setFeaturedImage(index) {
      this.images.forEach((img, i) => {
        img.is_featured = i === index
      })
    },
    
    // Base64 helper
    toBase64(file) {
      return new Promise((resolve, reject) => {
        const reader = new FileReader()
        reader.readAsDataURL(file)
        reader.onload = () => resolve(reader.result.split(',')[1])
        reader.onerror = error => reject(error)
      })
    },
    
    async submit(values) {
      try {
        // Validate form data
        if (this.steps.some(step => !step.instruction.trim())) {
          throw new Error('All steps must have instructions')
        }
        if (this.ingredients.some(ing => !ing.name.trim())) {
          throw new Error('All ingredients must have names')
        }
        if (this.images.length === 0) {
          throw new Error('Please add at least one image')
        }
        if (this.images.some(img => !img.file)) {
          throw new Error('Please select an image file for all images')
        }
        
        const recipePayload = {
          title: values.title,
          description: values.description || null,
          user_id: this.userId,
          prep_time: values.prep_time || null,
          cook_time: values.cook_time || null,
          servings: values.servings || null,
          category_ids: [values.category_id],
          featured_image_url: null, // Will be set after upload
          steps: this.steps.map(step => ({
            step_number: step.step_number,
            instruction: step.instruction,
            image_url: step.image_url || null
          })),
          ingredients: this.ingredients.map(ing => ({
            name: ing.name,
            quantity: ing.quantity || null,
            unit: ing.unit || null
          })),
          images: [] // Will be populated after upload
        }

        // 1. Create the recipe first
        this.creatingRecipe = true
        const { data: createData } = await this.$apollo.mutate({
          mutation: gql`
            mutation CreateRecipe($object: CreateRecipeInput!) {
              createRecipe(object: $object) {
                id
                title
                total_time
                featured_image_url
              }
            }
          `,
          variables: {
            object: recipePayload
          }
        })
        const createdRecipe = createData?.createRecipe
        
        if (!createdRecipe?.id) {
          throw new Error('Failed to create recipe')
        }

        // 2. Upload all images with the recipe's ID
        const uploadedImages = []
        let featuredImageUrl = null

        for (const [index, img] of this.images.entries()) {
          if (!img.file) continue
          
          const base64str = await this.toBase64(img.file)
          const filename = `${Date.now()}_${img.file.name.replace(/\s+/g, '_')}`

          this.uploadingImage = true
          const { data } = await this.$apollo.mutate({
            mutation: gql`
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
            `,
            variables: {
              recipe_id: createdRecipe.id,
              is_featured: img.is_featured,
              base64str,
              filename
            }
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
          this.updatingRecipe = true
          await this.$apollo.mutate({
            mutation: gql`
              mutation UpdateRecipe($id: uuid!, $featured_image_url: String) {
                update_recipes_by_pk(
                  pk_columns: {id: $id}
                  _set: {featured_image_url: $featured_image_url}
                ) {
                  id
                  featured_image_url
                }
              }
            `,
            variables: {
              id: createdRecipe.id,
              featured_image_url: featuredImageUrl
            }
          })
        }

        // Show success and redirect
        this.$router.push({
          path: `/recipes/${createdRecipe.id}`,
          query: { created: 'true' }
        })

      } catch (err) {
        alert(err.message || 'An error occurred while submitting the recipe.')
        console.error(err)
      } finally {
        this.creatingRecipe = false
        this.uploadingImage = false
        this.updatingRecipe = false
      }
    }
  }
}
</script>
<style scoped>
.form-input,
.form-textarea,
.form-select {
  display: block;
  width: 100%;
  border-radius: 0.375rem; /* rounded-md */
  border: 1px solid #d1d5db; /* border-gray-300 */
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05); /* shadow-sm */
  outline: none;
}

.form-input:focus,
.form-textarea:focus,
.form-select:focus {
  border-color: #6366f1; /* indigo-500 */
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.5); /* focus:ring-indigo-500 */
}

.form-input,
.form-textarea,
.form-select {
  margin-top: 0.25rem; /* mt-1 */
}

.btn {
  display: inline-flex;
  align-items: center;
  padding: 0.5rem 1rem; /* px-4 py-2 */
  border: 1px solid transparent;
  font-size: 0.875rem; /* text-sm */
  font-weight: 500; /* font-medium */
  border-radius: 0.375rem; /* rounded-md */
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05); /* shadow-sm */
  color: #fff;
  background-color: #4f46e5; /* bg-indigo-600 */
  cursor: pointer;
  outline: none;
}

.btn:hover {
  background-color: #4338ca; /* hover:bg-indigo-700 */
}

.btn:focus {
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.5); /* focus:ring-indigo-500 */
}

.btn:disabled {
  opacity: 0.75;
  cursor: not-allowed;
}
</style>
