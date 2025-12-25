// Composable to get API URL from environment variables
export const useApi = () => {
  const config = useRuntimeConfig();
  
  const apiUrl = computed(() => {
    // Check environment variable first (for Vercel)
    if (process.env.NUXT_PUBLIC_API_URL) {
      return process.env.NUXT_PUBLIC_API_URL;
    }
    // Fallback to runtime config
    if (config.public?.apiUrl) {
      return config.public.apiUrl;
    }
    // Last resort: localhost for local dev
    return 'http://localhost:8081';
  });
  
  const hasuraUrl = computed(() => {
    // Check environment variable first (for Vercel)
    if (process.env.NUXT_PUBLIC_HASURA_URL) {
      return process.env.NUXT_PUBLIC_HASURA_URL;
    }
    // Fallback to runtime config
    if (config.public?.hasuraUrl) {
      return config.public.hasuraUrl;
    }
    // Last resort: localhost for local dev
    return 'http://localhost:8080/v1/graphql';
  });
  
  return {
    apiUrl: apiUrl.value,
    hasuraUrl: hasuraUrl.value
  };
};


