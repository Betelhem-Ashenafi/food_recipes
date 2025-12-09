// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: false }, // Disabled to fix iframe sandbox security warning
  modules: ['@nuxtjs/tailwindcss'],

  // Fix dev.json and preload issues
  experimental: {
    payloadExtraction: false,
  },

  // Apollo Client is configured via plugin (plugins/apollo.client.ts)
})
