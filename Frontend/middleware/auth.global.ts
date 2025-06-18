export default defineNuxtRouteMiddleware(async (to, from) => {
  if (import.meta.client) {
    const useAuthStore = authStore();
    await useAuthStore.initAuthState();

    // If the user is authenticated
    if (useAuthStore.$state.isAuthed) {
      useAuthStore.getRole();

      // If the user is an admin
      if (
        useAuthStore.$state.isAdmin &&
        useAuthStore.$state.role === "systemAdmin"
      ) {
        return;
      } else {
        // If the user is not an admin
        if (
          useAuthStore.isAuthPage(to.name) ||
          useAuthStore.isAdminPage(to.name)
        ) {
          if (to.name !== "unauthorized") {
            return navigateTo("/unauthorized");
          }
        }
      }
    } else {
      // If the user is not authenticated
      if (!useAuthStore.isAnonymousPage(to.name)) {
        console.log(useAuthStore.isAnonymousPage(to.name));
        console.log("current page", to.name);
        if (to.name !== "unauthorized") {
          return navigateTo("/unauthorized");
        }
      }
    }
  }
});
