<template>
  <div class="max-w-xl mx-auto p-6">
    <h1 class="text-2xl font-bold mb-6 text-center">Payment Form</h1>
    <form @submit.prevent="submitPayment" class="space-y-4">
      <div>
        <label for="amount" class="block font-semibold mb-1">Amount (ETB):</label>
        <input v-model="form.amount" type="number" id="amount" required min="1"
               class="w-full border border-gray-300 rounded px-3 py-2" />
      </div>

      <div>
        <label for="currency" class="block font-semibold mb-1">Currency:</label>
        <select v-model="form.currency" id="currency" required
                class="w-full border border-gray-300 rounded px-3 py-2">
          <option value="ETB">ETB</option>
          <option value="USD">USD</option>
        </select>
      </div>

      <div>
        <label for="email" class="block font-semibold mb-1">Email:</label>
        <input v-model="form.email" type="email" id="email" required
               class="w-full border border-gray-300 rounded px-3 py-2" />
      </div>

      <div>
        <label for="firstName" class="block font-semibold mb-1">First Name:</label>
        <input v-model="form.first_name" type="text" id="firstName" required
               class="w-full border border-gray-300 rounded px-3 py-2" />
      </div>

      <div>
        <label for="lastName" class="block font-semibold mb-1">Last Name:</label>
        <input v-model="form.last_name" type="text" id="lastName" required
               class="w-full border border-gray-300 rounded px-3 py-2" />
      </div>

      <div>
        <label for="phoneNumber" class="block font-semibold mb-1">Phone Number:</label>
        <input v-model="form.phone_number" type="tel" id="phoneNumber" required
               class="w-full border border-gray-300 rounded px-3 py-2" />
      </div>

      <button type="submit" class="w-full bg-green-600 text-white py-2 rounded hover:bg-green-700">
        Pay Now
      </button>

      <p v-if="error" class="text-red-600 mt-2">{{ error }}</p>
    </form>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const form = ref({
  amount: 100,
  currency: 'ETB',
  email: 'test@example.com',
  first_name: 'John',
  last_name: 'Doe',
  phone_number: '0912345678',
  success_url: 'http://localhost:3000/success',
  error_url: 'http://localhost:3000/error',
})

const error = ref('')

const submitPayment = async () => {
  error.value = ''
  try {
    const res = await fetch('http://localhost:8081/pay', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(form.value),
    })

    const result = await res.json()
    if (result.status === 'success' && result.data?.checkout_url) {
      window.location.href = result.data.checkout_url
    } else {
      error.value = result.message || 'Failed to start payment'
    }
  } catch (err) {
    console.error(err)
    error.value = 'An error occurred while starting the payment'
  }
}
</script>
