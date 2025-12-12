import { ApolloClient, InMemoryCache, createHttpLink } from '@apollo/client/core';
import { setContext } from '@apollo/client/link/context';
import { provideApolloClient, DefaultApolloClient } from '@vue/apollo-composable';

export default defineNuxtPlugin((nuxtApp) => {
  // Auth link to inject JWT from cookie
  const authLink = setContext((_, { headers }) => {
    let token = '';
    if (process.client) {
      const cookie = useCookie('auth_token');
      token = cookie.value ?? '';
    }
    const newHeaders: Record<string, string> = { ...headers } as Record<string, string>;
    // Only include Authorization when a token is present
    if (token) newHeaders.Authorization = `Bearer ${token}`;
    // For local development: include Hasura admin secret so frontend can access Hasura
    // WARNING: sending admin secret from the browser is insecure. Use only for local testing.
    newHeaders['x-hasura-admin-secret'] = 'myhasurasecret';
    return { headers: newHeaders };
  });

  const httpLink = createHttpLink({
    uri: 'http://localhost:8080/v1/graphql',
  });

  // Create cache with safe handling for null/undefined
  const cache = new InMemoryCache({
    // Add type policies to handle null/undefined safely
    typePolicies: {
      Query: {
        fields: {
          recipes: {
            merge(existing = [], incoming) {
              // Ensure we always return an array, never null/undefined
              if (!incoming || !Array.isArray(incoming)) {
                return existing || [];
              }
              return incoming;
            },
            read(existing) {
              // Ensure read always returns an array
              return existing || [];
            }
          },
          categories: {
            merge(existing = [], incoming) {
              if (!incoming || !Array.isArray(incoming)) {
                return existing || [];
              }
              return incoming;
            },
            read(existing) {
              return existing || [];
            }
          },
          recipes_by_pk: {
            merge(existing, incoming) {
              // Return incoming if it exists, otherwise existing, but never null
              return incoming || existing || null;
            },
            read(existing) {
              return existing || null;
            }
          }
        }
      }
    },
    // Add result caching policy to prevent null/undefined issues
    resultCaching: true
  });

  const apolloClient = new ApolloClient({
    link: authLink.concat(httpLink),
    cache: cache,
    // Add default options to handle errors gracefully
    defaultOptions: {
      watchQuery: {
        errorPolicy: 'all',
        fetchPolicy: 'cache-and-network',
        // Return partial data even if there are errors
        returnPartialData: true
      },
      query: {
        errorPolicy: 'all',
        fetchPolicy: 'cache-first',
        returnPartialData: true
      },
      mutate: {
        errorPolicy: 'all'
      }
    }
  });

  // Register as the default Apollo client for the Vue app so `useQuery`/`useMutation` work
  if (nuxtApp.vueApp && nuxtApp.vueApp.provide) {
    nuxtApp.vueApp.provide(DefaultApolloClient, apolloClient);
  }
  // Also make available via Nuxt injection
  nuxtApp.provide('apollo', apolloClient);
  // Also provide with $ prefix so `useNuxtApp().$apollo` is always available
  nuxtApp.provide('$apollo', apolloClient);
});
