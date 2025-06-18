<script setup>
import { onMounted, ref } from "vue";
const recipesStore = useRecipeStore();
const categories = ref([]);

onMounted(async () => {
  try {
    await recipesStore.getCategories();
    categories.value = recipesStore.categories;
  } catch (error) {
    console.error("Failed to load categories", error);
  }
});
</script>

<template>
  <div>Hello {{}}</div>
</template>
