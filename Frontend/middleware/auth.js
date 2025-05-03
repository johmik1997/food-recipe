export default defineNuxtRouteMiddleware(async (to) => {
  const auth = useAuth()
  
  // Skip middleware for auth pages
  if (to.path === '/login' || to.path === '/register') {
    if (auth.isAuthenticated) {
      return navigateTo('/dashboard')
    }
    return
  }

  // Check authentication
  if (!auth.isAuthenticated) {
    return navigateTo('/login')
  }


})