import { useMutation } from '@vue/apollo-composable'
import gql from 'graphql-tag'

export const useAuth = () => {
  const loading = ref(false)
  const error = ref(null)
  const token = useCookie('auth-token')
  const router = useRouter()

  const register = async (userData) => {
    loading.value = true
    error.value = null
    
    try {
      // Use the mutation with proper context
      const { mutate } = useMutation(gql`
        mutation Register($input: RegisterInput!) {
          register(input: $input) {
            token
            user {
              id
              email
            }
          }
        }
      `)

      const result = await mutate({ input: userData })
      
      if (result?.data?.register?.token) {
        token.value = result.data.register.token
        router.push('/')
      } else {
        throw new Error('Registration failed')
      }
      
      return result
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  return { register, loading, error }
}