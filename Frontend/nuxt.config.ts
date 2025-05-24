import { defineNuxtConfig } from 'nuxt/config'
import tailwindcss from "@tailwindcss/vite";

export default defineNuxtConfig({
  modules: ['@nuxtjs/apollo', '@nuxt/image', '@vee-validate/nuxt','nuxt-icon'],
  plugins: [
    './plugins/apollo','./plugins/vee-validate'
  ],
  build: {
    transpile: ['@apollo/client', '@vue/apollo-composable', '@vee-validate/rules', '@vee-validate','ts-invariant/process'],
  },
  runtimeConfig: {
    public: {
      apollo: {
        clients: {
          default: {
            httpEndpoint: 'http://localhost:8080/v1/graphql' // your Go backend
          }
        }
      },
      veeValidate: {
        autoImports: true,
        componentNames: {
          Form: 'VeeForm',
          Field: 'VeeField',
          FieldArray: 'VeeFieldArray',
          ErrorMessage: 'VeeErrorMessage',
        },
      },
    }
  },

  css: ['~/assets/css/main.css'],
  devtools: { enabled: true },

  vite: {
    plugins: [tailwindcss()]
  },

  compatibilityDate: '2025-04-22'
})