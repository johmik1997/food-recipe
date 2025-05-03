import { ApolloClient, InMemoryCache, HttpLink } from '@apollo/client/core'
import { DefaultApolloClient } from '@vue/apollo-composable'

export default defineNuxtPlugin((nuxtApp) => {
  // Configure the HTTP connection to your GraphQL API
  const httpLink = new HttpLink({
    uri: 'http://localhost:8080/v1/graphql',
    headers: {
      'x-hasura-admin-secret': 'myadminsecretkey' // Add this if needed
    }
  })

  // Create the Apollo client
  const apolloClient = new ApolloClient({
    link: httpLink,
    cache: new InMemoryCache(),
    defaultOptions: {
      query: {
        fetchPolicy: 'no-cache',
      }
    }
  })

  // Provide the Apollo client to your Vue application
  nuxtApp.vueApp.provide(DefaultApolloClient, apolloClient)

  // Also provide it via Nuxt's provide
  return {
    provide: {
      apollo: apolloClient
    }
  }
})