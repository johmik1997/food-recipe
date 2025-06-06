import { gql } from '@apollo/client'

export const LOGIN_MUTATION = gql`
  mutation Login($username: String!, $password: String!) {
    login(username: $username, password: $password) {
      token
      user {
        id
        username
        isAdmin
      }
    }
  }
`

export const REGISTER_MUTATION = gql`
  mutation Register($username: String!, $password: String!) {
    register(username: $username, password: $password) {
      id
      username
      isAdmin
      createdAt
    }
  }
`