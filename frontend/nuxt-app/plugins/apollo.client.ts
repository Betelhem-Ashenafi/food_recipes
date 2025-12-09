import { ApolloClient, InMemoryCache, createHttpLink } from '@apollo/client/core';
import { setContext } from '@apollo/client/link/context';
import { provideApolloClient, DefaultApolloClient } from '@vue/apollo-composable';

export default defineNuxtPlugin((nuxtApp) => {
  // Auth link to inject JWT from cookie
  const authLink = setContext((_, { headers }) => {
    let token = '';
    if (process.client) {
      token = useCookie('auth_token').value ?? '';
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

  const apolloClient = new ApolloClient({
    link: authLink.concat(httpLink),
    cache: new InMemoryCache(),
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
