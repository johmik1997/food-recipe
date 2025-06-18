<template>
  <div class="flex flex-col container mx-auto">
    <AdminUsers />
  </div>
</template>
<script setup>
useSeoMeta({
  title: "recipe-app | Admin-dashboard",
  description: "my cool recipe app.",
});
definePageMeta({
  layout: "admin",
});
</script>
