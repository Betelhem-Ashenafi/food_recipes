// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  modules: ['@nuxtjs/tailwindcss', '@nuxtjs/apollo'],

  apollo: {
    clients: {
      default: {
        httpEndpoint: 'http://localhost:8080/v1/graphql',
        tokenStorage: 'localStorage',
        authHeader: 'Authorization',
        authType: 'Bearer',
        // @ts-ignore
        devtools: { enabled: false }
      }
    }
  }
})
