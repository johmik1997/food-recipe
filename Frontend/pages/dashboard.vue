<script setup>
import { ref } from 'vue'
import { useMutation } from '@vue/apollo-composable'
import gql from 'graphql-tag'

// States
const file = ref(null)
const isFeatured = ref(false)
const recipeId = "519d7df8-1071-4097-8e82-f43de76aaa96"
const response = ref(null)
const errorMessage = ref(null)
const isLoading = ref(false)

// Updated mutation — removed `image_url` from input
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

const { mutate } = useMutation(UPLOAD_IMAGE)

const handleFile = (e) => {
  const selected = e.target.files[0]
  if (!selected) return

  const validTypes = ['image/jpeg', 'image/png', 'image/webp']
  if (!validTypes.includes(selected.type)) {
    errorMessage.value = 'Only JPG, PNG, or WEBP images are allowed.'
    file.value = null
    return
  }

  if (selected.size > 5 * 1024 * 1024) {
    errorMessage.value = 'Image size should not exceed 5MB.'
    file.value = null
    return
  }

  file.value = selected
}


const uploadImage = async () => {
  if (!file.value) {
    errorMessage.value = 'Please select a file'
    return
  }

  isLoading.value = true
  errorMessage.value = null
  response.value = null

  try {
    const reader = new FileReader()
    reader.onload = async () => {
      const base64str = reader.result.split(',')[1]
      const filename = `${Date.now()}_${file.value.name}`

      const { data, errors } = await mutate({
        recipe_id: recipeId,
        is_featured: isFeatured.value,
        base64str,
        filename,
      })

      if (errors) {
        errorMessage.value = errors[0].message
      } else if (data?.uploadRecipeImage) {
        response.value = data.uploadRecipeImage
      }
      isLoading.value = false
    }

    reader.onerror = () => {
      errorMessage.value = 'Failed to read file'
      isLoading.value = false
    }

    reader.readAsDataURL(file.value)
  } catch (err) {
    errorMessage.value = err.message
    console.error('Upload error:', err)
    isLoading.value = false
  }
  if (data?.uploadRecipeImage) {
  response.value = data.uploadRecipeImage
  file.value = null
}

}
</script>

<template>
  <div>
    <h2>Upload Recipe Image</h2>
    <input type="file" @change="handleFile($event)" accept="image/*" />

    <label>
      <input type="checkbox" v-model="isFeatured" />
      Featured?
    </label>
    <button @click="uploadImage" :disabled="isLoading">
      {{ isLoading ? 'Uploading...' : 'Upload' }}
    </button>

    <div v-if="errorMessage" class="error-message">
      {{ errorMessage }}
    </div>

    <div v-if="response">
      <p>{{ response.message }}</p>
      <img :src="response.image_url" v-if="response.image_url" style="max-width: 300px; margin-top: 1rem;" />
    </div>
    <div v-if="file">
  <p><strong>Selected:</strong> {{ file.name }}</p>
  <p><strong>Featured:</strong> {{ isFeatured ? 'Yes' : 'No' }}</p>
</div>
    <div v-else>
      <p>No file selected</p>
    </div>  
    
  </div>
</template>

<style>
.error-message {
  color: red;
  margin-top: 1rem;
}
</style>
