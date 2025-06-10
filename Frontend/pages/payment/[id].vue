<template>
    <div class="min-h-screen bg-gray-50">
    <Navbar />
    
  <div class="max-w-xl mx-auto p-6">
    <h1 class="text-2xl font-bold mb-6 text-center">Payment for {{ recipeTitle }}</h1>
    <form @submit.prevent="submitPayment" class="space-y-4">
      <!-- Recipe Information (display only) -->
      <div class="bg-gray-50 p-4 rounded-lg">
        <h2 class="font-semibold text-lg">{{ recipeTitle }}</h2>
        <p class="text-gray-600">Price: {{ form.amount }} ETB</p>
      </div>
      <div>
        <label for="email" class="block font-semibold mb-1">Email:</label>
        <input v-model="form.email" type="email" id="email" required
               class="w-full border border-gray-300 rounded px-3 py-2" />
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label for="first_name" class="block font-semibold mb-1">First Name:</label>
          <input v-model="form.first_name" type="text" id="firstName" required
                 class="w-full border border-gray-300 rounded px-3 py-2" />
        </div>
        <div>
          <label for="last_name" class="block font-semibold mb-1">Last Name:</label>
          <input v-model="form.last_name" type="text" id="lastName" required
                 class="w-full border border-gray-300 rounded px-3 py-2" />
        </div>
      </div>

      <div>
        <label for="phone_number" class="block font-semibold mb-1">Phone Number:</label>
        <input v-model="form.phone_number" type="tel" id="phoneNumber" required
               class="w-full border border-gray-300 rounded px-3 py-2" />
      </div>

      <button type="submit" 
              class="w-full bg-green-600 text-white py-2 rounded hover:bg-green-700 disabled:opacity-50"
              :disabled="isLoading">
        <span v-if="isLoading">Processing...</span>
        <span v-else>Pay Now ({{ form.amount }} ETB)</span>
      </button>

      <div v-if="error" class="p-4 bg-red-50 text-red-600 rounded-lg">
        {{ error }}
      </div>
    </form>
  </div>
  <Footer/>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import gql from 'graphql-tag'
import { useQuery } from '@vue/apollo-composable'

const route = useRoute()
const recipeTitle = ref('')
const isLoading = ref(false)
const error = ref('')
const userId = ref(null)

// Initialize form with recipe ID from route
const form = ref({
  recipe_id: route.params.id,
  amount: 0,
  currency: 'ETB',
  email: '',
  first_name: '',
  last_name: '',
  phone_number: '',
  success_url: '',
  error_url: '',
})

// GraphQL Query
const GET_RECIPE_DETAILS = gql`
  query GetRecipePrice($id: uuid!) {
    recipes(where: {id: {_eq: $id}}, limit: 1) {
      id
      title
      price
    }
  }
`

// Fetch recipe details using Apollo
const { onResult, onError, loading: queryLoading } = useQuery(
  GET_RECIPE_DETAILS,
  () => ({
    id: form.value.recipe_id
  }),
  {
    enabled: computed(() => !!form.value.recipe_id)
  }
)

// Handle query results
onResult((result) => {
  if (result.data?.recipes?.[0]) {
    const recipe = result.data.recipes[0]
    form.value.amount = recipe.price || 0
    recipeTitle.value = recipe.title
 setTimeout(() => {
  form.value.success_url = window.location.origin + '/payment/success' 
  }, 5000)  
    form.value.error_url   = window.location.origin + '/payment/error'
  }})
// Handle query errors
onError((err) => {
  error.value = err.message || 'Failed to load recipe details'
  console.error('Query error:', err)
})

// Watch loading state
watch(queryLoading, (loading) => {
  isLoading.value = loading
})

// Get user ID from localStorage
onMounted(() => {
  const user = JSON.parse(localStorage.getItem('user'))
  userId.value = user?.id || null
  if (!userId.value) {
    // Redirect to login or show an error message
    navigateTo('/login')
    // For this example, we'll just set an error message
    error.value = 'User not authenticated'

  }
})

// Submit payment to your backend
const submitPayment = async () => {
  try {
    isLoading.value = true
    error.value = ''

    const response = await fetch('http://localhost:3001/api/payments', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(form.value)
    })

    const result = await response.json()

    if (result.status === 'success' && result.data?.checkout_url) {
      window.location.href = result.data.checkout_url
      
    } else {
      throw new Error(result.message || 'Payment failed')
    }
  } catch (err) {
    console.error('Payment error:', err)
    error.value = err.message || 'An error occurred during payment'
  } finally {
    isLoading.value = false
  }
}
</script>