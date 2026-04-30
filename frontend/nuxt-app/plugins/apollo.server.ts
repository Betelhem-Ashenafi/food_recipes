import { createHttpLink, InMemoryCache, ApolloClient } from '@apollo/client/core';
import { setContext } from '@apollo/client/link/context';
import { DefaultApolloClient } from '@vue/apollo-composable';

export default defineNuxtPlugin((nuxtApp) => {
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

  const config = useRuntimeConfig();
  const hasuraUrl = config.public?.hasuraUrl || 'http://localhost:8080/v1/graphql';

  const httpLink = createHttpLink({ uri: hasuraUrl });

  const cache = new InMemoryCache({
    typePolicies: {
      Query: {
        fields: {
          recipes: {
            merge(existing = [], incoming) {
              return incoming || existing;
            },
            read(existing) {
              return existing || [];
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
    resultCaching: true
  });

  const apolloClient = new ApolloClient({
    link: authLink.concat(httpLink),
    cache,
    defaultOptions: {
      watchQuery: {
        errorPolicy: 'all',
        fetchPolicy: 'cache-and-network',
        returnPartialData: true
      },
      query: {
        errorPolicy: 'all',
        fetchPolicy: 'cache-first'
      },
      mutate: {
        errorPolicy: 'all'
      }
    }
  });

  nuxtApp.vueApp.provide(DefaultApolloClient, apolloClient);
  nuxtApp.provide('apollo', apolloClient);
});
