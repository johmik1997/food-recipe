<template>
  <div class="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-50 flex items-center justify-center p-4">
    <div class="w-full max-w-md">
      <div class="text-center mb-8">
        <NuxtImg src="/logo.png" alt="Logo" class="mx-auto h-16 w-auto mb-4" />
        <h2 class="text-3xl font-bold text-indigo-900 font-serif">Create Account</h2>
      </div>

      <Form @submit="handleRegister" class="space-y-6" v-slot="{ errors }">
        <div>
          <label for="name" class="block text-sm font-medium text-indigo-800">Full Name</label>
          <div class="mt-1 relative rounded-md shadow-sm">
            <Field
              id="name"
              name="name"
              type="text"
              rules="required|min:3"
              :class="`block w-full pl-3 pr-3 py-3 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 placeholder-indigo-400 ${
                errors.name ? 'border-red-300' : 'border-indigo-300'
              }`"
              placeholder="John Doe"
            />
            <ErrorMessage name="name" class="mt-2 text-sm text-rose-600" />
          </div>
        </div>

        <div>
          <label for="email" class="block text-sm font-medium text-indigo-800">Email</label>
          <div class="mt-1 relative rounded-md shadow-sm">
            <Field
              id="email"
              name="email"
              type="email"
              rules="required|email"
              :class="`block w-full pl-3 pr-3 py-3 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 placeholder-indigo-400 ${
                errors.email ? 'border-red-300' : 'border-indigo-300'
              }`"
              placeholder="your@email.com"
            />
            <ErrorMessage name="email" class="mt-2 text-sm text-rose-600" />
          </div>
        </div>

        <div>
          <label for="password" class="block text-sm font-medium text-indigo-800">Password</label>
          <div class="mt-1 relative rounded-md shadow-sm">
            <Field
              id="password"
              name="password"
              type="password"
              rules="required|min:8"
              :class="`block w-full pl-3 pr-3 py-3 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 placeholder-indigo-400 ${
                errors.password ? 'border-red-300' : 'border-indigo-300'
              }`"
              placeholder="••••••••"
            />
            <ErrorMessage name="password" class="mt-2 text-sm text-rose-600" />
          </div>
        </div>

        <div>
          <label for="confirmPassword" class="block text-sm font-medium text-indigo-800">Confirm Password</label>
          <div class="mt-1 relative rounded-md shadow-sm">
            <Field
              id="confirmPassword"
              name="confirmPassword"
              type="password"
              rules="required|confirmed:@password"
              :class="`block w-full pl-3 pr-3 py-3 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 placeholder-indigo-400 ${
                errors.confirmPassword ? 'border-red-300' : 'border-indigo-300'
              }`"
              placeholder="••••••••"
            />
            <ErrorMessage name="confirmPassword" class="mt-2 text-sm text-rose-600" />
          </div>
        </div>

        <button
          type="submit"
          :disabled="loading"
          class="w-full flex justify-center items-center py-3 px-4 border border-transparent rounded-lg shadow-sm text-lg font-medium text-white bg-gradient-to-r from-indigo-500 to-blue-500 hover:from-indigo-600 hover:to-blue-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-all duration-300 disabled:opacity-70"
        >
          <span v-if="!loading">Register</span>
          <span v-else>Creating account...</span>
        </button>
        <div class="bg-indigo-50 px-8 py-4 rounded-b-2xl text-center">
          <p class="text-sm text-indigo-700">
            Already have an account? 
            <NuxtLink to="/login" class="font-medium text-indigo-600 hover:text-indigo-500">
              Sign in
            </NuxtLink>
          </p>
        </div>
      </Form>
    </div>
  </div>
</template>

<script>
import { defineComponent } from 'vue';
//import { useMutation } from '@apollo/client/core';
import gql from 'graphql-tag';

const REGISTER_MUTATION = gql`
  mutation Register($email: String!, $name: String!, $password: String!) {
    register(input: { email: $email, name: $name, password: $password }) {
      token
      user {
        email
        id
      }
    }
  }
`;

export default defineComponent({
  data() {
    return {
      loading: false,
    };
  },
  methods: {
    async handleRegister({ name, email, password, confirmPassword }) {
      if (password !== confirmPassword) {
        alert('Passwords do not match!');
        return;
      }

      this.loading = true;

      try {
        const response = await this.$apollo.mutate({
          mutation: REGISTER_MUTATION,
          variables: { email, name, password },
        });

        const { data } = response;
        if (data.register) {
          console.log('Registration successful', data);
          alert('Account created successfully');
          // Optionally, store token and redirect
        } else {
          alert('Failed to create account');
        }
      } catch (error) {
        console.error('Registration error:', error);
        alert('An error occurred while creating the account');
      } finally {
        this.loading = false;
      }
    },
  },
});
</script>


