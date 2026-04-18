import { GraphQLClient } from 'graphql-request'

// GraphQL API endpoint - replace with your actual backend URL
const API_URL = import.meta.env.VITE_GRAPHQL_ENDPOINT || 'http://localhost:4000/graphql'

// Create GraphQL client instance
export const graphQLClient = new GraphQLClient(API_URL, {
  headers: {
    'Content-Type': 'application/json',
  },
})

// Helper function to set authentication token
export const setAuthToken = (token) => {
  if (token) {
    graphQLClient.setHeader('Authorization', `Bearer ${token}`)
  } else {
    graphQLClient.deleteHeader('Authorization')
  }
}

export default graphQLClient
