// This is a Nuxt plugin – it runs when your app starts.
import { createHttpLink, InMemoryCache, ApolloClient } from '@apollo/client/core';
import { setContext } from '@apollo/client/link/context';
import { DefaultApolloClient } from '@vue/apollo-composable';
export default defineNuxtPlugin((nuxtApp) => {
  // -------------------- AUTH: Attach JWT token to every request --------------------
  // Think of this as a "stamp" that adds the user's ID card to every envelope.
  const authLink = setContext((_, prevContext) => {
    let token = '';
    if (process.client) {
      const cookie = useCookie('auth_token');
      token = cookie.value ?? '';
    }
    const headers = prevContext.headers || {};
    const newHeaders = { ...headers };
    if (token) newHeaders.Authorization = `Bearer ${token}`;
    return { headers: newHeaders };
  });

  // -------------------- WHERE TO SEND REQUESTS --------------------
  // This is the address of your Hasura GraphQL API.
  const config = useRuntimeConfig();
  const hasuraUrl = config.public?.hasuraUrl || 'http://localhost:8080/v1/graphql';

  // This creates the network link – the actual path the messenger walks.
  const httpLink = createHttpLink({ uri: hasuraUrl });

  // -------------------- CACHE: Remember data so you don't ask twice --------------------
  // Apollo keeps a "memory" of previous answers. This cache configuration tells it how to handle special cases.
  const cache = new InMemoryCache({
    // Here we're just making sure that lists (like recipes) are always treated as arrays, even if empty.
    typePolicies: {
      Query: {
        fields: {
          recipes: {
            merge(existing = [], incoming) {
              return incoming || existing; // keep existing if incoming is nothing
            },
            read(existing) {
              return existing || []; // always return an array (never null)
            }
          },
          categories: {
            merge(existing = [], incoming) {
              return incoming || existing;
            },
            read(existing) {
              return existing || [];
            }
          },
          // For a single recipe, just return it or null.
          recipes_by_pk: {
            merge(existing, incoming) {
              return incoming || existing || null;
            },
            read(existing) {
              return existing || null;
            }
          }
        }
      }
    },
    resultCaching: true // enable caching
  });

  // -------------------- PUT IT ALL TOGETHER --------------------
  // Build the messenger: first attach auth, then send request, with a cache.
  const apolloClient = new ApolloClient({
    link: authLink.concat(httpLink), // chain: add token, then send
    cache: cache,
    defaultOptions: {
      watchQuery: {
        errorPolicy: 'all',        // don't crash on errors, just tell me
        fetchPolicy: 'cache-and-network', // use cache first, then update from network
        returnPartialData: true     // give me whatever data you have, even if some fields error
      },
      query: {
        errorPolicy: 'all',
        fetchPolicy: 'cache-first' // for one‑off queries, use cache first
      },
      mutate: {
        errorPolicy: 'all'
      }
    }
  });

  // -------------------- MAKE IT AVAILABLE TO VUE COMPONENTS --------------------
  // Now any component can use `useQuery`, `useMutation` and they'll automatically use this client.
  nuxtApp.vueApp.provide(DefaultApolloClient, apolloClient);
  nuxtApp.provide('apollo', apolloClient); // optional, but nice to have
});