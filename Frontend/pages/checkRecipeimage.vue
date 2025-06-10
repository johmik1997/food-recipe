<template>
  <div class="max-w-md mx-auto p-6 bg-white rounded-lg shadow-md">
    <h2 class="text-2xl font-bold mb-6 text-center">Upload Recipe Image</h2>

    <form @submit.prevent="uploadRecipeImage" class="space-y-6">
     <!-- Recipe Image Upload -->
<div class="flex flex-col items-center">
  <!-- Image Upload Box -->
  <div class="relative mb-4 w-full">
    <div class="border-2 border-dashed border-gray-300 rounded-lg p-4 text-center">
      <div v-if="!previewImage" class="py-8">
        <svg xmlns="http://www.w3.org/2000/svg" class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
        </svg>
        <p class="mt-2 text-sm text-gray-600">Drag and drop your image here, or click to select</p>
      </div>
      <img 
        v-else
        :src="previewImage" 
        class="mx-auto max-h-64 rounded-lg object-contain"
        alt="Recipe image preview"
      >
      <input 
        type="file"
        accept="image/jpeg, image/png, image/webp"
        class="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
        @change="handleImageUpload"
        aria-label="Upload recipe image"
      >
    </div>
    <div v-if="previewImage" class="mt-2 text-right">
      <button 
        type="button"
        @click="clearImage"
        class="text-sm text-red-600 hover:text-red-800"
      >
        Remove
      </button>
    </div>
    <p v-if="uploadError" class="mt-1 text-sm text-red-500">{{ uploadError }}</p>
  </div>

  <!-- Featured Checkbox Section -->
  <div class="w-full mb-4">
    <label class="flex items-center space-x-2">
      <input 
        type="checkbox" 
        v-model="isFeatured" 
        class="rounded text-blue-600 focus:ring-blue-500"
      >
      <span class="text-sm text-gray-700">Set as featured image</span>
    </label>
    <p class="text-sm text-gray-500 mt-1">Featured: {{ isFeatured }}</p>
  </div>
</div>


      <!-- Submit Button -->
      <button
        type="submit"
        :disabled="loading || !selectedFile"
        class="w-full px-6 py-3 bg-blue-600 text-white font-medium rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
      >
        <span v-if="!loading" class="flex items-center justify-center">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
          </svg>
          Upload Image
        </span>
        <span v-else class="flex items-center justify-center">
          <svg class="animate-spin -ml-1 mr-2 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          Uploading...
        </span>
      </button>

      <!-- Status Messages -->
      <div v-if="successMessage" class="p-4 bg-green-100 text-green-700 rounded-md flex items-start">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2 mt-0.5 flex-shrink-0" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
        </svg>
        <span>{{ successMessage }}</span>
      </div>

      <div v-if="errorMessage" class="p-4 bg-red-100 text-red-700 rounded-md flex items-start">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2 mt-0.5 flex-shrink-0" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
        </svg>
        <span>{{ errorMessage }}</span>
      </div>
    </form>
  </div>
</template>


<script setup>
import { ref } from 'vue'
import { useMutation } from '@vue/apollo-composable'
import gql from 'graphql-tag'

// Reactive state
const selectedFile = ref(null)
const previewImage = ref(null)
const isFeatured = ref(true)
const uploadError = ref(null)
const errorMessage = ref(null)
const successMessage = ref(null)
const loading = ref(false)

// HARDCODED RECIPE ID FOR TESTING
const testRecipeId = 'dbe4967f-1fdb-4956-9f25-16ef9f6b9d53'

// GraphQL Mutation
const UPLOAD_RECIPE_IMAGE = gql`
  mutation UploadRecipeImage($input: UploadRecipeImageInput!) {
    uploadRecipeImage(input: $input) {
      success
      message
      image_url
    }
  }
`

const { mutate: uploadRecipeImageMutation } = useMutation(UPLOAD_RECIPE_IMAGE)

// Handle image upload and validation
const handleImageUpload = (e) => {
  const file = e.target.files[0]
  
  // Reset states
  uploadError.value = null
  errorMessage.value = ''
  successMessage.value = ''
  selectedFile.value = null
  previewImage.value = null

  if (!file) return
  
  // Validate file type
  const validTypes = ['image/jpeg', 'image/png', 'image/webp']
  if (!validTypes.includes(file.type)) {
    uploadError.value = 'Only JPG, PNG, or WEBP images are allowed'
    return
  }
  
  // Validate file size (5MB max)
  if (file.size > 5 * 1024 * 1024) {
    uploadError.value = 'Image size must be less than 5MB'
    return
  }

  selectedFile.value = file

  // Create preview
  const reader = new FileReader()
  reader.onload = (e) => {
    previewImage.value = e.target.result
  }
  reader.onerror = () => {
    uploadError.value = 'Error reading image file'
  }
  reader.readAsDataURL(file)
}

// Clear selected image
const clearImage = () => {
  selectedFile.value = null
  previewImage.value = null
  uploadError.value = null
}

// Convert file to base64 (returns raw base64 without data URL prefix)
const toBase64 = (file) => new Promise((resolve, reject) => {
  const reader = new FileReader()
  reader.readAsDataURL(file)
  reader.onload = () => {
    const result = reader.result
    // Extract just the base64 part after the comma
    resolve(result.split(',')[1])
  }
  reader.onerror = error => reject(error)
})

const uploadRecipeImage = async () => {
  if (!selectedFile.value) {
    errorMessage.value = 'Please select an image to upload'
    return
  }

  loading.value = true
  errorMessage.value = ''
  successMessage.value = ''

  try {
    console.log('[DEBUG] Starting image upload process...')
    
    // Convert image to base64
    const base64str = await toBase64(selectedFile.value)
    const filename = `${Date.now()}_${selectedFile.value.name.replace(/\s+/g, '_')}`

    console.log('[DEBUG] Prepared file data:', {
      filename,
      size: selectedFile.value.size,
      type: selectedFile.value.type
    })
console.log(isFeatured.value)
    // Prepare input with hardcoded recipeId
    const variables = {
      input: {
        recipe_id: testRecipeId, // Using hardcoded ID
        is_featured: isFeatured.value,
        base64str: base64str,
        filename: filename
      }
    }

    console.log('[DEBUG] Sending mutation with variables:', variables)

    const { data, errors } = await uploadRecipeImageMutation(variables)

    console.log('[DEBUG] Mutation response:', { data, errors })

    if (errors) {
      throw new Error(errors[0].message)
    }

    if (data?.uploadRecipeImage?.success) {
      successMessage.value = data.uploadRecipeImage.message || 'Image uploaded successfully!'
      console.log('[SUCCESS] Image uploaded:', data.uploadRecipeImage.image_url)
      clearImage()
    } else {
      errorMessage.value = data?.uploadRecipeImage?.message || 'Upload failed (no success flag)'
    }

  } catch (error) {
    errorMessage.value = error.graphQLErrors?.[0]?.message || 
                        error.networkError?.message || 
                        error.message || 
                        'Failed to upload image'
    console.error('[ERROR] Upload failed:', error)
  } finally {
    loading.value = false
  }
}
</script>