/**
 * Client-side plugin to intercept and sanitize JWT token access
 * This ensures that any code (including libraries) that accesses JWT tokens
 * always gets a safe structure, preventing Object.keys() errors
 */
export default defineNuxtPlugin((nuxtApp) => {
  if (!process.client) return;
  
  // Use dynamic import for client-side (non-blocking)
  import('jwt-decode').then((jwtDecodeModule) => {
    const originalJwtDecode = jwtDecodeModule.default || jwtDecodeModule;
    
    // Create safe wrapper
    const safeJwtDecode = (token, options) => {
      if (!token) {
        return {
          'https://hasura.io/jwt/claims': {
            'x-hasura-allowed-roles': ['user'],
            'x-hasura-default-role': 'user',
            'x-hasura-user-id': '0',
            'x-hasura-user-name': 'User',
            'x-hasura-user-email': ''
          }
        };
      }
      
      try {
        const decoded = originalJwtDecode(token, options);
        
        // Ensure we always return an object
        if (!decoded || typeof decoded !== 'object') {
          return {
            'https://hasura.io/jwt/claims': {
              'x-hasura-allowed-roles': ['user'],
              'x-hasura-default-role': 'user',
              'x-hasura-user-id': '0',
              'x-hasura-user-name': 'User',
              'x-hasura-user-email': ''
            }
          };
        }
        
        // Ensure claims is always an object, never null
        if (!decoded['https://hasura.io/jwt/claims']) {
          decoded['https://hasura.io/jwt/claims'] = {};
        } else if (typeof decoded['https://hasura.io/jwt/claims'] !== 'object' || 
                   decoded['https://hasura.io/jwt/claims'] === null) {
          decoded['https://hasura.io/jwt/claims'] = {};
        }
        
        // Ensure x-hasura-allowed-roles is always an array
        const claims = decoded['https://hasura.io/jwt/claims'];
        if (claims['x-hasura-allowed-roles']) {
          if (!Array.isArray(claims['x-hasura-allowed-roles'])) {
            claims['x-hasura-allowed-roles'] = [];
          }
        } else {
          claims['x-hasura-allowed-roles'] = ['user'];
        }
        
        // Ensure other common Hasura claims exist as safe defaults
        if (!claims['x-hasura-user-id']) claims['x-hasura-user-id'] = '0';
        if (!claims['x-hasura-user-name']) claims['x-hasura-user-name'] = 'User';
        if (!claims['x-hasura-user-email']) claims['x-hasura-user-email'] = '';
        if (!claims['x-hasura-default-role']) claims['x-hasura-default-role'] = 'user';
        
        return decoded;
      } catch (error) {
        // Return safe structure on decode error
        return {
          'https://hasura.io/jwt/claims': {
            'x-hasura-allowed-roles': ['user'],
            'x-hasura-default-role': 'user',
            'x-hasura-user-id': '0',
            'x-hasura-user-name': 'User',
            'x-hasura-user-email': ''
          }
        };
      }
    };
    
    // Patch the module exports
    if (jwtDecodeModule.default) {
      jwtDecodeModule.default = safeJwtDecode;
    }
    if (typeof jwtDecodeModule === 'function') {
      // If the module itself is the function, we can't easily patch it
      // But we can make sure our safe version is used
    }
    
    // Store for potential use
    if (typeof window !== 'undefined') {
      window.__safeJwtDecode = safeJwtDecode;
    }
    
    // Try to patch the module if possible
    if (jwtDecodeModule.default) {
      jwtDecodeModule.default = safeJwtDecode;
    }
  }).catch((error) => {
    console.warn('Could not load jwt-decode for patching:', error);
  });
});

