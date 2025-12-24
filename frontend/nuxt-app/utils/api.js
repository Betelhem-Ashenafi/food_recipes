// Centralized API configuration
export const getApiUrl = () => {
  // In browser/client-side
  if (process.client) {
    return process.env.NUXT_PUBLIC_API_URL || 'http://localhost:8081';
  }
  // Server-side
  const config = useRuntimeConfig();
  return config.public.apiUrl || 'http://localhost:8081';
};

export const getHasuraUrl = () => {
  if (process.client) {
    return process.env.NUXT_PUBLIC_HASURA_URL || 'http://localhost:8080/v1/graphql';
  }
  const config = useRuntimeConfig();
  return config.public.hasuraUrl || 'http://localhost:8080/v1/graphql';
};

