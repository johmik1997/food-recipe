<script setup>
// import { userStore } from "~/stores/user";
import { useToast } from "vue-toast-notification";
import { ref, computed, onMounted } from "vue";

const useUserStore = userStore();
const toast = useToast();
const emit = defineEmits(["refetch-deleted-users"]);
const currentPage = ref(1);
const itemsPerPage = 5;
const userToDelete = ref(null);
const showDeleteConfirmation = ref(false);
const showDeleteModal = ref(false);
const deleteModalRef = ref(null);

const openDeleteModal = (user) => {
  userToDelete.value = user;
  deleteModalRef.value.showModal();
};

const closeDeleteModal = () => {
  deleteModalRef.value.close();
};
// Fetch users on mount
onMounted(async () => {
  try {
    await useUserStore.getUsers();
  } catch (error) {
    toast.error("Failed to load users");
    console.error("Failed to load users", error);
  }
});

// Function to delete a user
const handleDeleteUser = async (userId) => {
  console.log("Deleting user with ID:", userId);
  try {
    const success = await useUserStore.deleteUser(userId);
    if (success) {
      toast.success("user deleted successfully!");
      userToDelete.value = null;
      await useUserStore.getUsers();
      showDeleteModal.value = false;
    } else {
      toast.error("Failed to delete user");
    }
  } catch (err) {
    toast.error("Failed to delete user");
    console.error(err);
  }
};

// Pagination logic
const paginatedUsers = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage;
  const end = start + itemsPerPage;
  return useUserStore.users.slice(start, end);
});

const totalPages = computed(() =>
  Math.ceil(useUserStore.users.length / itemsPerPage)
);

const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++;
  }
};

const prevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--;
  }
};
</script>

<template>
  <div class="p-6 bg-gray-100 mb-20 dark:bg-[#41d137]">
    <div class="bg-white dark:bg-[#41d137] shadow-md rounded-lg p-6">
      <h1 class="text-2xl font-bold text-gray-700 dark:text-gray-100 mb-4">
        Users
      </h1>
      <div class="overflow-x-auto">
        <table class="table-auto w-full border-collapse border border-gray-300">
          <!-- Table Head -->
          <thead
            class="bg-gray-200 dark:bg-[#41d137] text-black dark:text-black uppercase text-sm leading-normal"
          >
            <tr>
              <th class="py-3 px-6 text-left">
                <label>
                  <input type="checkbox" class="checkbox" />
                </label>
              </th>
              <th class="py-3 px-6 text-left">Id</th>
              <th class="py-3 px-6 text-left">Name</th>
              <th class="py-3 px-6 text-left">Phone</th>
              <th class="py-3 px-6 text-left">Role</th>
              <th class="py-3 px-6 text-left">isVerified</th>
              <th class="py-3 px-6 text-left">Joined-in</th>
              <th class="py-3 px-10 text-start">Actions</th>
            </tr>
          </thead>
          <!-- Table Body -->
          <tbody class="text-black dark:text-gray-900 text-sm font-light">
            <tr
              v-for="user in paginatedUsers"
              :key="user.id"
              class="border-b border-gray-400 dark:border-gray-700 hover:bg-gray-100 dark:hover:bg-gray-800"
            >
              <td class="py-3 px-6 text-left">
                <label>
                  <input type="checkbox" class="checkbox" />
                </label>
              </td>
              <td class="py-3 px-6 text-left text-black dark:text-gray-100">
                <div>
                  {{ user.id }}
                </div>
              </td>
              <td class="py-3 px-6">
                <div class="flex items-center gap-3">
                  <div class="avatar">
                    <div class="mask mask-squircle h-12 w-12">
                      <img
                        :src="user.profile || '/images/user1.png'"
                        alt="User Avatar"
                      />
                    </div>
                  </div>
                  <div>
                    <div class="font-bold text-black dark:text-gray-100">
                      {{ user.username }}
                    </div>
                    <div class="text-sm text-black dark:text-gray-100">
                      {{ user.email }}
                    </div>
                  </div>
                </div>
              </td>
              <td class="py-3 px-6 text-black dark:text-gray-100">
                {{ user.phone || "N/A" }}
              </td>
              <td class="py-3 px-6 text-black dark:text-gray-100">
                {{ user.role || "N/A" }}
              </td>
              <td class="py-3 px-6 font-bold">
                <span
                  :class="
                    user.is_email_verified
                      ? 'text-green-500 bg-green-100 px-4 py-1 rounded-lg text-xs '
                      : 'text-red-500 bg-red-100  px-4 py-1 rounded-lg text-xs'
                  "
                >
                  {{ user.is_email_verified ? "True" : "False" }}
                </span>
              </td>

              <td class="py-3 px-6">
                <span
                  class="bg-blue-100 text-blue-600 px-2 py-1 rounded-lg text-xs"
                >
                  {{
                    user.created_at
                      ? new Date(user.created_at).toLocaleDateString()
                      : "N/A"
                  }}
                </span>
              </td>
              <td class="py-3 px-6 text-center">
                <!-- Action Buttons -->
                <div class="flex gap-2">
                  <NuxtLink
                    :to="`/admin/users/${user.id}`"
                    class="btn btn-info btn-xs"
                  >
                    Details
                  </NuxtLink>
                  <!-- update user -->
                  <NuxtLink
                    :to="`/admin/update/${user.id}`"
                    class="btn btn-warning btn-xs"
                  >
                    update
                  </NuxtLink>
                  <!-- <button class="btn btn-warning btn-xs">Update</button> -->
                  <button
                    class="btn btn-error btn-xs"
                    @click="openDeleteModal(user)"
                  >
                    Delete
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <!-- Pagination -->
      <div class="flex justify-between items-center mt-4">
        <button
          class="btn btn-sm btn-primary dark:text-gray-100"
          @click="prevPage"
          :disabled="currentPage === 1"
        >
          Previous
        </button>
        <span>Page {{ currentPage }} of {{ totalPages }}</span>
        <button
          class="btn btn-sm btn-primary dark:text-gray-100"
          @click="nextPage"
          :disabled="currentPage === totalPages"
        >
          Next
        </button>
      </div>
      <!-- Confirmation Modal -->
      <dialog ref="deleteModalRef" class="modal">
        <div class="modal-box">
          <form method="dialog">
            <button
              class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2"
              @click="closeDeleteModal"
            >
              x
            </button>
          </form>
          <h3 class="text-lg font-bold">Confirm Deletion</h3>
          <p class="py-4">
            Are you sure you want to delete
            <strong>{{ userToDelete?.username }}?</strong>
          </p>

          <div class="flex justify-end gap-4">
            <button
              class="btn bg-gray-100 hover:bg-gray-300"
              @click="closeDeleteModal"
            >
              cancel
            </button>
            <button
              class="btn bg-red-50 text-black hover:bg-red-500"
              @click="handleDeleteUser(userToDelete?.id)"
            >
              Confirm
            </button>
          </div>
        </div>
      </dialog>
    </div>
  </div>
</template>