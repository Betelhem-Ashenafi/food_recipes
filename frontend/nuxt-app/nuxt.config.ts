// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: false }, // Disabled to fix iframe sandbox security warning
  modules: ['@nuxtjs/tailwindcss'],

  // Fix dev.json and preload issues
  experimental: {
    payloadExtraction: false,
  },

  // Environment variables
  runtimeConfig: {
    public: {
      apiUrl: process.env.NUXT_PUBLIC_API_URL || 'http://localhost:8081',
      hasuraUrl: process.env.NUXT_PUBLIC_HASURA_URL || 'http://localhost:8080/v1/graphql',
      hasuraAdminSecret: process.env.NUXT_PUBLIC_HASURA_ADMIN_SECRET || 'myhasurasecret'
    }
  },

  // Apollo Client is configured via plugin (plugins/apollo.client.ts)
})
