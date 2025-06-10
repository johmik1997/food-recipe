// plugins/apollo.js
import {
  ApolloClient,
  InMemoryCache,
  createHttpLink,
  ApolloLink
} from '@apollo/client/core'
import { DefaultApolloClient } from '@vue/apollo-composable'
import { defineNuxtPlugin, useRuntimeConfig } from '#app'
import { useCookie } from '#imports'
// plugins/apollo.js
export default defineNuxtPlugin((nuxtApp) => {
  const authToken = useCookie('auth_token')

  const httpLink = createHttpLink({
    uri: 'http://localhost:8080/v1/graphql',
  })

  const authMiddleware = new ApolloLink((operation, forward) => {
    const token = authToken.value || localStorage.getItem('token')
    
    operation.setContext(({ headers = {} }) => ({
      headers: {
      ...headers,
      'content-type': 'application/json',
      'x-hasura-admin-secret': 'myadminsecretkey',
      'authorization': `Bearer ${token}`,
      }
    }))
    
    return forward(operation)
  })

  const client = new ApolloClient({
    link: ApolloLink.from([authMiddleware, httpLink]),
    cache: new InMemoryCache(),
    defaultOptions: {
      query: {
        fetchPolicy: 'no-cache',
      },
      mutate: {
        fetchPolicy: 'no-cache',
      }
    }
  })

  nuxtApp.vueApp.provide(DefaultApolloClient, client)
})