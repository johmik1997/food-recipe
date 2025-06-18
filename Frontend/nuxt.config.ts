import { defineNuxtConfig } from "nuxt/config";

export default defineNuxtConfig({
  compatibilityDate: "2024-11-01",
  devtools: { enabled: true },
  css: ["~/assets/css/main.css"],
  modules: [
    "@nuxtjs/apollo",
    "@pinia/nuxt",
    "@nuxtjs/color-mode",
    "@nuxtjs/tailwindcss",

    [
      "@vee-validate/nuxt",
      {
        autoImports: true,
      },
    ],
    [
      "@vee-validate/nuxt",
      {
        autoImports: true,
        componentNames: {
          Form: "VeeForm",
          Field: "VeeField",
          FieldArray: "VeeFieldArray",
          ErrorMessage: "VeeErrorMessage",
        },
      },
    ],
    "@pinia/nuxt",
  ],
  colorMode: {
    classSuffix: "",
  },
  postcss: {
    plugins: {
      tailwindcss: {},
      autoprefixer: {},
    },
  },
  pinia: {
    storesDirs: ["./stores/**", "./custom-folder/stores/**"],
  },
  
  apollo: {
    autoImports: true,
    clients: {
      default: {
        httpEndpoint: "http://localhost:8084/v1/graphql",
        tokenStorage: "localStorage",
        authHeader: "Authorization",
        authType: "Bearer",
        tokenName: "authToken",
        httpLinkOptions: {
          // headers: {
          //   "x-hasura-admin-secret": "myadminsecretkey",
          // },
          // headers: {
          //   Authorization: localStorage.getItem("authToken")
          //     ? `Bearer ${localStorage.getItem("authToken")}`
          //     : "",
          //   "x-hasura-role": localStorage.getItem("authRole") || "anonymous",
          // },
        },

        defaultOptions: {
          query: {
            fetchPolicy: "network-only",
          },
          watchQuery: {
            fetchPolicy: "network-only",
          },
        },
      },
    },
  },

  build: {
    transpile: [/@nuxtjs[\\/]composition-api/],
  },
});