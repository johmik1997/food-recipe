export const useAuth = () => {
  const router = useRouter()
  
  const login = async (email, password) => {
    // Apollo mutation logic here
  }

  const register = async (userData) => {
    // Apollo mutation logic here  
  }

  const logout = () => {
    const token = useCookie('auth_token')
    token.value = null
    localStorage.removeItem('token') // Add this
  localStorage.removeItem('user')  
    router.push('/login')
  }

  return { login, register, logout }
}