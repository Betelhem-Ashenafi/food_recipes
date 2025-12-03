import { createClient, ssrExchange, fetchExchange } from '@urql/core';
import urql from '@urql/vue';
import { defineNuxtPlugin } from 'nuxt/app';

export default defineNuxtPlugin((nuxtApp) => {
  const ssr = ssrExchange({
    isClient: process.client,
  });

  const client = createClient({
    url: 'http://localhost:8080/v1/graphql',
    exchanges: [ssr, fetchExchange],
    fetchOptions: () => {
      const token = process.client ? localStorage.getItem('token') : null;
      return {
        headers: {
          ...(token ? { Authorization: `Bearer ${token}` } : {}),
        },
      };
    },
  });

  nuxtApp.vueApp.use(urql, client);
});
