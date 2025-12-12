import { computed } from 'vue';
import { jwtDecode } from 'jwt-decode';

/**
 * Safe JWT decoding that always returns a valid object structure
 * This prevents Object.keys() errors when JWT claims are null/undefined
 */
export const useAuth = () => {
  const token = useCookie('auth_token');
  
  const userInfo = computed(() => {
    if (!token.value) {
      return {
        'https://hasura.io/jwt/claims': {}
      };
    }
    
    try {
      const decoded = jwtDecode(token.value);
      
      // Ensure we always return an object
      if (!decoded || typeof decoded !== 'object') {
        return {
          'https://hasura.io/jwt/claims': {}
        };
      }
      
      // Ensure claims is always an object, never null
      if (!decoded['https://hasura.io/jwt/claims']) {
        decoded['https://hasura.io/jwt/claims'] = {};
      } else if (typeof decoded['https://hasura.io/jwt/claims'] !== 'object' || decoded['https://hasura.io/jwt/claims'] === null) {
        decoded['https://hasura.io/jwt/claims'] = {};
      }
      
      // Ensure x-hasura-allowed-roles is always an array if present
      const claims = decoded['https://hasura.io/jwt/claims'];
      if (claims['x-hasura-allowed-roles']) {
        if (!Array.isArray(claims['x-hasura-allowed-roles'])) {
          claims['x-hasura-allowed-roles'] = [];
        }
      } else {
        claims['x-hasura-allowed-roles'] = [];
      }
      
      return decoded;
    } catch (error) {
      console.error('JWT decode error:', error);
      return {
        'https://hasura.io/jwt/claims': {}
      };
    }
  });
  
  const isAuthenticated = computed(() => !!token.value);
  
  const userName = computed(() => {
    const claims = userInfo.value?.['https://hasura.io/jwt/claims'];
    if (!claims || typeof claims !== 'object' || claims === null) return 'User';
    return claims['x-hasura-user-name'] || 'User';
  });
  
  const userEmail = computed(() => {
    const claims = userInfo.value?.['https://hasura.io/jwt/claims'];
    if (!claims || typeof claims !== 'object' || claims === null) return '';
    return claims['x-hasura-user-email'] || '';
  });
  
  const userId = computed(() => {
    const claims = userInfo.value?.['https://hasura.io/jwt/claims'];
    if (!claims || typeof claims !== 'object' || claims === null) return 0;
    return parseInt(claims['x-hasura-user-id']) || 0;
  });
  
  const userRoles = computed(() => {
    const claims = userInfo.value?.['https://hasura.io/jwt/claims'];
    if (!claims || typeof claims !== 'object' || claims === null) return [];
    const roles = claims['x-hasura-allowed-roles'];
    return Array.isArray(roles) ? roles : [];
  });
  
  return {
    token,
    userInfo,
    isAuthenticated,
    userName,
    userEmail,
    userId,
    userRoles
  };
};



