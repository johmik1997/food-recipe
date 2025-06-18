<script setup>
import { ref, onMounted } from "vue";
import { useRoute } from "vue-router";
import { useToast } from "vue-toast-notification";

useSeoMeta({
  title: "recipe-app | User Details Page",
  description: "The project app meta.",
});

const route = useRoute();
const toast = useToast();
const userId = parseInt(route.params.id);
const useUserStore = userStore();
const userDetails = ref(null);
const isLoading = ref(false);

// Fetch user details
const fetchUserDetails = async () => {
  isLoading.value = true;
  try {
    await useUserStore.getUser(userId);
    userDetails.value = useUserStore.user;
    console.log("user details", userDetails.value);
  } catch (error) {
    toast.error("Failed to load user details");
    console.error(error);
  } finally {
    isLoading.value = false;
  }
};

// Fetch data on component mount
onMounted(fetchUserDetails);
</script>

<template>
  <div class="p-6 min-h-80vh bg-white">
    <!-- Loading Spinner -->
    <div v-if="isLoading" class="flex justify-center items-center min-h-screen">
      <span class="loading loading-spinner loading-lg text-green-500"></span>
    </div>

    <!-- User Details -->
    <div
      v-else-if="userDetails"
      class="bg-white shadow-lg rounded-xl p-6 max-w-4xl mx-auto"
    >
      <!-- User Basic Info -->
      <div class="flex flex-col md:flex-row items-center gap-6 mb-8">
        <img
          :src="userDetails.profile || '/images/user1.png'"
          alt="User Avatar"
          class="h-24 w-24 md:h-32 md:w-32 rounded-full border-4 border-green-500 shadow-md"
        />
        <div class="text-center md:text-left">
          <h1 class="text-3xl font-bold text-gray-800">
            {{ userDetails.username }}
          </h1>
          <p class="text-gray-600">
            {{ userDetails.email }}
          </p>
          <p class="text-gray-600">
            Phone: <span>{{ userDetails.phone || "N/A" }}</span>
          </p>
          <p
            :class="
              userDetails.is_email_verified ? 'text-green-500' : 'text-red-500'
            "
            class="font-semibold"
          >
            Email Verified: {{ userDetails.is_email_verified ? "Yes" : "No" }}
          </p>
          <p class="text-gray-600">
            Role: <span class="font-semibold">{{ userDetails.role }}</span>
          </p>
        </div>
      </div>
    </div>

    <!-- No User Details -->
    <div v-else class="text-center text-gray-500">
      <p>No user details found.</p>
    </div>
  </div>
</template>