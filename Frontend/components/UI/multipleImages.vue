<script setup>
import { ref, onMounted } from 'vue'
import { Field } from "vee-validate"

const emit = defineEmits(["image-changed"])

const imagePreviews = ref([])
const selectedThumbnail = ref(null)
const images = ref([])
const thumbnail = ref(null)

const props = defineProps({
  name: {
    type: String,
    default: "images",
  },
  isVideo: {
    type: Boolean,
    default: false,
  },
  errorMessage: {
    type: String,
  },
  isRequired: {
    type: Boolean,
    default: false,
  },
  isMultiple: {
    type: Boolean,
    default: false,
  },
})

onMounted(() => {
  const fileInput = document.getElementById("dropzone-file")
  if (fileInput) {
    fileInput.addEventListener("change", handleFileChange)
  }
})

const handleFileChange = (event) => {
  const files = event.target.files
  if (!files || files.length === 0) return

  // Reset previous previews
  imagePreviews.value = []
  images.value = []

  for (let i = 0; i < files.length; i++) {
    const file = files[i]
    if (!file.type.startsWith("image/")) continue

    const reader = new FileReader()
    reader.onload = () => {
      const base64String = reader.result.match(/base64,(.*)$/)[1]
      const image = {
        name: file.name,
        type: file.type,
        base64String,
      }

      images.value.push(image)
      imagePreviews.value.push(reader.result)
    }
    reader.readAsDataURL(file)
  }

  emit("image-changed", {
    thumbnail: selectedThumbnail.value,
    images: images.value,
  })
}

const selectThumbnail = (index) => {
  selectedThumbnail.value = images.value[index]
  thumbnail.value = imagePreviews.value[index]
  emit("image-changed", {
    thumbnail: selectedThumbnail.value,
    images: images.value,
  })
}

const removeImage = (index) => {
  if (selectedThumbnail.value === images.value[index]) {
    selectedThumbnail.value = null
    thumbnail.value = null
  }

  imagePreviews.value.splice(index, 1)
  images.value.splice(index, 1)

  emit("image-changed", {
    thumbnail: selectedThumbnail.value,
    images: images.value,
  })
}
</script>

<template>
  <div class="flex flex-col gap-4 items-center justify-center w-full">
    <label
      for="dropzone-file"
      class="flex flex-col items-center justify-center w-full h-32 border-2 border-dashed rounded-xl cursor-pointer bg-white hover:bg-gray-50 border-gray-300 hover:border-green-400 transition-colors"
    >
      <div class="flex flex-col items-center justify-center pt-5 pb-6">
        <svg
          class="w-8 h-8 mb-4 text-green-500"
          aria-hidden="true"
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 20 16"
        >
          <path
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M13 13h3a3 3 0 0 0 0-6h-.025A5.56 5.56 0 0 0 16 6.5 5.5 5.5 0 0 0 5.207 5.021C5.137 5.017 5.071 5 5 5a4 4 0 0 0 0 8h2.167M10 15V6m0 0L8 8m2-2 2 2"
          />
        </svg>
        <p class="mb-2 text-sm text-gray-600">
          <span class="font-semibold text-green-600">Click to upload</span>
        </p>
        <p class="text-xs text-gray-500" v-if="!isVideo">
          SVG, PNG, JPG or GIF (MAX. 800x400px)
        </p>
        <p class="text-xs text-gray-500" v-else>
          WEBP, MP4 or MP3 (MAX. 10MB)
        </p>
      </div>
      <Field
        as="input"
        id="dropzone-file"
        ref="fileInput"
        type="file"
        :name="name"
        class="hidden"
        :multiple="isMultiple"
      />
    </label>
    
    <p class="text-sm text-red-500" v-if="errorMessage">
      {{ errorMessage }}
    </p>

    <div class="w-full mt-2 flex flex-wrap gap-3">
      <div
        v-if="thumbnail"
        class="relative flex flex-col items-center gap-1 pr-3 border-r border-gray-200"
      >
        <img
          :src="thumbnail"
          class="h-16 w-16 rounded-lg object-cover border-2 border-green-400"
          alt="Thumbnail"
        />
        <span class="text-xs text-gray-500">Thumbnail</span>
      </div>
      
      <div
        v-for="(previewUrl, index) in imagePreviews"
        :key="index"
        class="relative group"
      >
        <img
          :src="previewUrl"
          class="h-16 w-16 rounded-lg object-cover cursor-pointer border"
          :class="{
            'border-green-500': selectedThumbnail && selectedThumbnail.name === images[index].name,
            'border-gray-200': !(selectedThumbnail && selectedThumbnail.name === images[index].name)
          }"
          alt="Preview"
          @click="selectThumbnail(index)"
        />
        <button
          type="button"
          @click="removeImage(index)"
          class="absolute -top-2 -right-2 bg-red-500 text-white rounded-full p-0.5 text-xs opacity-0 group-hover:opacity-100 transition-opacity"
        >
          <svg
            class="w-3 h-3"
            aria-hidden="true"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 14 14"
          >
            <path
              stroke="currentColor"
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="m1 1 6 6m0 0 6 6M7 7l6-6M7 7l-6 6"
            />
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>