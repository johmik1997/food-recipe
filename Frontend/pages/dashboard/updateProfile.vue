<template>
    <div class="max-w-md mx-auto p-6 bg-white rounded-lg shadow-md">
      <h2 class="text-2xl font-bold mb-6 text-center">Update Profile Picture</h2>
      
      <form @submit.prevent="updateProfile" class="space-y-6">
        <!-- Profile Picture Upload -->
        <div class="flex flex-col items-center">
          <div class="relative mb-4">
            <img 
              :src="previewImage || currentUser.avatar || '/default-avatar.png'" 
              class="w-32 h-32 rounded-full object-cover border-2 border-gray-200"
            >
            <label 
              for="avatar-upload"
              class="absolute bottom-0 right-0 bg-white p-2 rounded-full shadow-md cursor-pointer hover:bg-gray-100"
            >
              <CameraIcon class="w-5 h-5 text-gray-600" />
              <input 
                id="avatar-upload"
                type="file"
                accept="image/jpeg, image/png, image/webp"
                class="hidden"
                @change="handleImageUpload"
              >
            </label>
          </div>
          <p v-if="uploadError" class="text-red-500 text-sm">{{ uploadError }}</p>
        </div>
  
        <!-- Submit Button -->
        <button
          type="submit"
          :disabled="loading"
          class="w-full px-6 py-3 bg-primary-600 text-white font-medium rounded-md hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <span v-if="!loading">Update Profile Picture</span>
          <span v-else class="flex items-center justify-center">
            <Spinner class="w-5 h-5 mr-2" />
            Uploading...
          </span>
        </button>
  
        <!-- Success Message -->
        <div v-if="successMessage" class="p-4 bg-green-100 text-green-700 rounded-md">
          {{ successMessage }}
        </div>
  
        <!-- Error Message -->
        <div v-if="errorMessage" class="p-4 bg-red-100 text-red-700 rounded-md">
          {{ errorMessage }}
        </div>
      </form>
    </div>
  </template>
  <script setup>
  import { ref, onMounted } from 'vue'
  import { CameraIcon } from '@heroicons/vue/outline'
  import Spinner from '@/components/Spinner.vue'
  import gql from 'graphql-tag'
  import { useMutation } from '@vue/apollo-composable'
  
  // GraphQL mutations
  const UPLOAD_PROFILE_IMAGE = gql`
    mutation UploadProfileImage($base64str: String!, $filename: String!) {
      uploadProfileImage(base64str: $base64str, filename: $filename) {
        success
        message
        image_url
      }
    }
  `
  
  const UPDATE_USER_PROFILE = gql`
    mutation UpdateUserProfile($avatar: String!) {
      updateUserProfile(avatar: $avatar) {
        id
        avatar_image_url
      }
    }
  `
  
  const { mutate: uploadProfileImage } = useMutation(UPLOAD_PROFILE_IMAGE)
  const { mutate: updateUserProfile } = useMutation(UPDATE_USER_PROFILE)
  
  // Component state
  const currentUser = ref({
    id: '',
    name: '',
    avatar: ''
  })
  
  const previewImage = ref(null)
  const selectedFile = ref(null)
  const uploadError = ref(null)
  const loading = ref(false)
  const successMessage = ref('')
  const errorMessage = ref('')
  
  // Fetch current user data
  onMounted(() => {
    const userStr = localStorage.getItem("user")
    if (userStr) {
      try {
        const user = JSON.parse(userStr)
        // Set all user properties including avatar if available
        currentUser.value = {
          id: user.id,
          name: user.name || '',
          avatar_image_url: user.avatar || ''
        }
      } catch (error) {
        errorMessage.value = 'Failed to load profile data'
        console.error('Error parsing user data:', error)
      }
    }
  })
  
  // Handle image upload
  const handleImageUpload = (event) => {
    const file = event.target.files[0]
    uploadError.value = null
  
    // Validate file
    if (!file) return
    
    const validTypes = ['image/jpeg', 'image/png', 'image/webp']
    if (!validTypes.includes(file.type)) {
      uploadError.value = 'Only JPG, PNG, or WEBP images are allowed'
      return
    }
    
    if (file.size > 5 * 1024 * 1024) { // 5MB limit
      uploadError.value = 'Image size must be less than 5MB'
      return
    }
  
    // Store the file for upload
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
  
  // Convert file to base64
  const toBase64 = (file) => new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.readAsDataURL(file)
    reader.onload = () => resolve(reader.result.split(',')[1])
    reader.onerror = error => reject(error)
  })
  
  // Submit form
  const updateProfile = async () => {
    if (!selectedFile.value) {
      errorMessage.value = 'Please select an image to upload'
      return
    }
  
    loading.value = true
    successMessage.value = ''
    errorMessage.value = ''
  
    try {
      // 1. Convert image to base64
      const base64str = await toBase64(selectedFile.value)
      const filename = `${Date.now()}_${selectedFile.value.name.replace(/\s+/g, '_')}`
  
      // 2. Upload image
      const { data: uploadData, errors: uploadErrors } = await uploadProfileImage({
        base64str,
        filename
      })
  
      if (uploadErrors) {
        throw new Error(uploadErrors[0].message || 'Image upload failed')
      }
  
      if (!uploadData?.uploadProfileImage?.image_url) {
        throw new Error('Failed to get image URL after upload')
      }
  
      const imageUrl = uploadData.uploadProfileImage.image_url
  
      // 3. Update user profile with new avatar URL
      const { data: profileData, errors: profileErrors } = await updateUserProfile({
        avatar: imageUrl
      })
  
      if (profileErrors) {
        throw new Error(profileErrors[0].message || 'Profile update failed')
      }
  
      // Update local state and localStorage
      currentUser.value.avatar_image_url = imageUrl
      const updatedUser = {
        ...JSON.parse(localStorage.getItem("user")),
        avatar_image_url: imageUrl
      }
      localStorage.setItem("user", JSON.stringify(updatedUser))
      
      successMessage.value = 'Profile picture updated successfully!'
      selectedFile.value = null
    } catch (error) {
      errorMessage.value = error.message || 'Failed to update profile picture'
      console.error('Profile update error:', error)
    } finally {
      loading.value = false
    }
  }
  </script>