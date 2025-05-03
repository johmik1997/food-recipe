<template>
  <div class="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-50 flex items-center justify-center p-4">
    <div class="w-full max-w-md">
      <div class="text-center mb-8">
        <NuxtImg src="/logo.png" alt="Logo" class="mx-auto h-16 w-auto mb-4" />
        <h2 class="text-3xl font-bold text-indigo-900 font-serif">Sign In</h2>
      </div>

      <Form @submit="handleLogin" class="space-y-6" v-slot="{ errors }">
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

        <button
          type="submit"
          :disabled="loading"
          class="w-full flex justify-center items-center py-3 px-4 border border-transparent rounded-lg shadow-sm text-lg font-medium text-white bg-gradient-to-r from-indigo-500 to-blue-500 hover:from-indigo-600 hover:to-blue-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-all duration-300 disabled:opacity-70"
        >
          <span v-if="!loading">Sign In</span>
          <span v-else>Signing in...</span>
        </button>
        <div class="bg-indigo-50 px-8 py-4 rounded-b-2xl text-center">
          <p class="text-sm text-indigo-700">
            Don't have an account? 
            <NuxtLink to="/register" class="font-medium text-indigo-600 hover:text-indigo-500">
              Sign up
            </NuxtLink>
          </p>
        </div>
      </Form>
    </div>
  </div>
</template>

<script>
import { defineComponent } from 'vue';
import gql from 'graphql-tag';
const LOGIN_MUTATION = gql`
  mutation Login($email: String!, $password: String!) {
    login(email: $email, password: $password) {
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
    async handleLogin({ email, password }) {
      this.loading = true;

      try {
        const response = await this.$apollo.mutate({
          mutation: LOGIN_MUTATION,
          variables: { email, password },
        });

        const { data } = response;
        if (data.login) {
          console.log('Login successful', data);
          alert('Login successful');
          // Store token in localStorage
          localStorage.setItem('token', data.login.token);
          localStorage.setItem('user', JSON.stringify(data.login.user));
          // Redirect to the dashboard or home page
          this.$router.push('/dashboard');
        } else {
          alert('Invalid credentials');
        }
      } catch (error) {
        console.error('Login error:', error);
        alert('An error occurred while logging in');
      } finally {
        this.loading = false;
      }
    },
  },
});
</script>
