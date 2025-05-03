export default defineNuxtPlugin((nuxtApp) => {
    const token = useCookie('auth-token', {
      httpOnly: true,
      secure: process.env.NODE_ENV === 'production',
      sameSite: 'strict',
      maxAge: 60 * 60 * 24 // 1 day
    })
  
    const setAuthToken = (newToken) => {
      token.value = newToken
    }
  
    const getAuthToken = () => token.value
  
    const logout = () => {
      token.value = null
      navigateTo('/login')
    }
  
    return {
      provide: {
        auth: {
          setAuthToken,
          getAuthToken,
          logout
        }
      }
    }
  })