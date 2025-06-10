<template>
  <div class="max-w-md mx-auto p-6 bg-white rounded-lg shadow-md">
    <h2 class="text-2xl font-bold mb-6 text-center">Update Profile Picture</h2>

    <form @submit.prevent="uploadProfileImage" class="space-y-6">
      <!-- Profile Picture Upload -->
      <div class="flex flex-col items-center">
        <div class="relative mb-4">
          <img 
            :src="previewImage || currentUser.avatar_image_url || '/default-avatar.png'" 
            class="w-32 h-32 rounded-full object-cover border-2 border-gray-200"
            alt="Profile picture"
          >
          <label 
            for="avatar-upload"
            class="absolute bottom-0 right-0 bg-white p-2 rounded-full shadow-md cursor-pointer hover:bg-gray-100 transition-colors"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
            <input 
              id="avatar-upload"
              type="file"
              accept="image/jpeg, image/png, image/webp"
              class="hidden"
              @change="handleImageUpload"
              aria-label="Upload profile picture"
            >
          </label>
        </div>
        <p v-if="uploadError" class="text-red-500 text-sm">{{ uploadError }}</p>
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
          Update Profile Picture
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
import { ref, onMounted } from 'vue'
import { useMutation } from '@vue/apollo-composable'
import gql from 'graphql-tag'

// Reactive state
const currentUser = ref({
  avatar_image_url: ''
})

const selectedFile = ref(null)
const previewImage = ref(null)
const uploadError = ref(null)
const errorMessage = ref(null)
const successMessage = ref(null)
const loading = ref(false)
const userId = ref(null)

// GraphQL Mutation
const UPLOAD_PROFILE_IMAGE = gql`
  mutation UploadProfileImage($input: UploadProfileInput!) {
    uploadProfileImage(input: $input) {
      success
      message
      image_url
    }
  }
`

const { mutate: uploadProfileImageMutation, onError, onDone } = useMutation(UPLOAD_PROFILE_IMAGE)

// Get user ID from localStorage
onMounted(() => {
  const userStr = localStorage.getItem("user")
  if (userStr) {
    try {
      const user = JSON.parse(userStr)
      userId.value = user.id
      if (!userId.value) {
        errorMessage.value = "User session error. Please log in again."
      }
      currentUser.value.avatar_image_url = user.avatar_image_url || ''
    } catch (e) {
      errorMessage.value = "Session data error. Please log in again."
    }
  } else {
    errorMessage.value = "No user session found. Please log in."
  }
})

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

// Set up mutation handlers
onError((error) => {
  errorMessage.value = error.message || 'Failed to upload profile image'
  if (error.graphQLErrors && error.graphQLErrors.length > 0) {
    errorMessage.value = error.graphQLErrors[0].message
  }
  loading.value = false
})

onDone((result) => {
  if (result.data?.uploadProfileImage?.success) {
    const response = result.data.uploadProfileImage
    console.log("Upload successful:", response)
    if(response.success==true) {
      console.log("Profile picture updated successfully:", response.avatar_image_url)
    
    successMessage.value = response.message || 'Profile picture updated successfully'
    currentUser.value.avatar_image_url = response.avatar_image_url
    }
    // Update local storage
    const userData = JSON.parse(localStorage.getItem('user') || '{}')
    userData.avatar_image_url = response.avatar_image_url
    localStorage.setItem('user', JSON.stringify(userData))

    // Clear selected file after successful upload
    selectedFile.value = null
  } else {
    errorMessage.value = result.data?.uploadProfileImage?.message || 'Upload failed'
  }
  loading.value = false
})
const uploadProfileImage = async () => {
  if (!selectedFile.value) {
    errorMessage.value = 'Please select an image to upload'
    return
  }

  if (!userId.value) {
    errorMessage.value = 'User session error. Please log in again.'
    return
  }

  loading.value = true
  errorMessage.value = ''
  successMessage.value = ''

  try {
    // Convert image to base64
    const base64str = await toBase64(selectedFile.value)
    const filename = `${Date.now()}_${selectedFile.value.name.replace(/\s+/g, '_')}`

    // Properly structured input
    const variables = {
      input: {  // Note the "input" wrapper
        userId: userId.value,
        base64str: base64str,
        filename: filename
      }
    }

    console.log("Sending payload:", variables)

    const { data, errors } = await uploadProfileImageMutation(variables)

    if (errors) throw new Error(errors[0].message)
    if (!data?.uploadProfileImage?.success) {
      throw new Error(data?.uploadProfileImage?.message || 'Upload failed')
    }

    // Update UI and local storage
    successMessage.value = data.uploadProfileImage.message || 'Profile picture updated successfully'
    currentUser.value.avatar_image_url = data.uploadProfileImage.avatar_image_url
    
    // Update local storage
    const userData = JSON.parse(localStorage.getItem('user') || '{}')
    userData.avatar_image_url = data.uploadProfileImage.avatar_image_url
    localStorage.setItem('user', JSON.stringify(userData))

    // Clear selected file after successful upload
    selectedFile.value = null

  } catch (error) {
    errorMessage.value = error.message || 'Failed to upload profile image'
    console.error('Upload error:', error)
  } finally {
    loading.value = false
  }
}
</script>