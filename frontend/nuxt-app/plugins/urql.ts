import { createClient, ssrExchange, fetchExchange } from '@urql/core';
import { defineNuxtPlugin } from 'nuxt/app';
import { ref } from 'vue';

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

  nuxtApp.vueApp.provide('$urql', client);
  nuxtApp.provide('urql', client);
});
