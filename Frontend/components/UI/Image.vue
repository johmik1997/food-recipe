<script setup>
import { Field } from "vee-validate";
import { ref, defineProps, defineEmits, onMounted } from "vue";

const emit = defineEmits(["image-changed"]);
const imagePreviews = ref([]);
const fileInput = ref(null);
const videoFileName = ref("");

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
});

const handleFileChange = (event) => {
  const files = event.target.files;
  if (!files || files.length === 0) return;

  imagePreviews.value = [];

  if (props.isVideo) {
    videoFileName.value = files[0].name;
  } else {
    for (let i = 0; i < files.length; i++) {
      const file = files[i];
      if (!file.type.startsWith("image/")) continue;

      const reader = new FileReader();
      reader.onload = () => {
        imagePreviews.value.push(reader.result);
      };
      reader.readAsDataURL(file);
    }
  }
  emit("image-changed", event);
};
</script>
<template>
  <fieldset class="fieldset">
    <legend class="fieldset-legend">
      Pick a file
    </legend>
    <Field
      as="input"
      type="file"
      class="file-input"
      id="dropzone-file"
      ref="fileInput"
      :name="name"
      :multiple="isMultiple"
      @change="handleFileChange"
    />
    <label class="fieldset-label">Max size 2MB</label>

    <p class="error-message" v-if="errorMessage">
      {{ errorMessage }}
    </p>

    <div class="preview-container">
      <div v-for="(previewUrl, index) in imagePreviews" :key="index">
        <img :src="previewUrl" class="preview-image" alt="Preview" />
      </div>
      <div v-if="isVideo && videoFileName">{{ videoFileName }}</div>
    </div>
  </fieldset>
</template>

<style scoped>
.fieldset {
  border: 2px dashed #e5e7eb;
  padding: 1.5rem;
  border-radius: 0.5rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.75rem;
  background-color: white;
  transition: all 0.3s ease;
}

.fieldset:hover {
  border-color: #4ade80;
}

.fieldset-legend {
  font-weight: 600;
  font-size: 0.875rem;
  color: #374151;
  padding: 0 0.5rem;
}

.file-input {
  display: block;
  cursor: pointer;
  padding: 0.5rem 1rem;
  border-radius: 0.375rem;
  background-color: #f0fdf4;
  color: #065f46;
  border: 1px solid #d1fae5;
  transition: all 0.2s ease;
}

.file-input:hover {
  background-color: #dcfce7;
  border-color: #a7f3d0;
}

.file-input:focus {
  outline: none;
  box-shadow: 0 0 0 2px rgba(74, 222, 128, 0.5);
}

.fieldset-label {
  font-size: 0.75rem;
  color: #6b7280;
}

.error-message {
  font-size: 0.75rem;
  color: #ef4444;
  margin-top: 0.5rem;
}

.preview-container {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-top: 1rem;
}

.preview-image {
  width: 3.75rem;
  height: 3.75rem;
  object-fit: cover;
  border-radius: 9999px;
  border: 2px solid #d1fae5;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
}
</style>