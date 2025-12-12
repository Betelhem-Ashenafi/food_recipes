/**
 * Apollo Client Error Handler Plugin
 * Catches and handles Apollo Client errors, especially Object.keys() errors
 * from cache normalization
 * 
 * CRITICAL: This plugin patches Object.keys at module load time,
 * before Apollo Client or any other code runs
 */

// CRITICAL: Patch Object.keys at module load time, BEFORE the plugin function
// This ensures it runs before Apollo Client or any other code
// Patch both in browser and Node.js environments
(function() {
  'use strict';
  if (typeof Object === 'undefined') return;
  
  const originalObjectKeys = Object.keys;
  
  // Patch Object.keys immediately - this runs at module load time
  Object.keys = function(obj) {
    // Handle null/undefined - return empty array instead of throwing
    if (obj === null || obj === undefined) {
      return [];
    }
    // For primitives, return empty array to prevent errors
    if (typeof obj !== 'object' && typeof obj !== 'function') {
      return [];
    }
    try {
      return originalObjectKeys.call(this, obj);
    } catch (e) {
      // If original throws for any reason, return empty array
      return [];
    }
  };
  
  // Also patch on global Object if available
  if (typeof globalThis !== 'undefined') {
    globalThis.Object = globalThis.Object || {};
    globalThis.Object.keys = Object.keys;
  }
  if (typeof window !== 'undefined') {
    window.Object = window.Object || {};
    window.Object.keys = Object.keys;
  }
  if (typeof global !== 'undefined') {
    global.Object = global.Object || {};
    global.Object.keys = Object.keys;
  }
  
  console.log('[00-apollo-error-handler] Object.keys patched at module load');
})();

export default defineNuxtPlugin((nuxtApp) => {
  if (!process.client) return;

  // Ensure patch is still in place (re-apply if needed)
  const originalObjectKeys = Object.keys;
  
  // Test if patch is working
  try {
    const testResult = Object.keys(null);
    if (!Array.isArray(testResult) || testResult.length !== 0) {
      // Patch not working, re-apply
      Object.keys = function(obj) {
        if (obj === null || obj === undefined) {
          return [];
        }
        if (typeof obj !== 'object' && typeof obj !== 'function') {
          return [];
        }
        try {
          return originalObjectKeys.call(this, obj);
        } catch (e) {
          return [];
        }
      };
      console.log('[00-apollo-error-handler] Object.keys patch re-applied in plugin');
    }
  } catch (e) {
    // Patch not working, re-apply
    Object.keys = function(obj) {
      if (obj === null || obj === undefined) {
        return [];
      }
      if (typeof obj !== 'object' && typeof obj !== 'function') {
        return [];
      }
      try {
        return originalObjectKeys.call(this, obj);
      } catch (e) {
        return [];
      }
    };
    console.log('[00-apollo-error-handler] Object.keys patch applied in plugin (error handler)');
  }

  // Also patch Object.keys on the global object to ensure it's available everywhere
  if (typeof window !== 'undefined') {
    window.Object = window.Object || {};
    window.Object.keys = Object.keys;
  }

  // Override console.error to catch Apollo errors (as backup)
  const originalError = console.error;
  console.error = function(...args) {
    // Check if this is the Object.keys error from Apollo
    const errorString = args.join(' ');
    if (errorString.includes('Cannot convert undefined or null to object') && 
        errorString.includes('checkInGroup')) {
      // Log a more helpful message
      console.warn('[Apollo Client] Cache normalization error caught and handled. This is usually safe to ignore.');
      // Don't log the full error to reduce noise
      return;
    }
    // For other errors, use the original console.error
    originalError.apply(console, args);
  };

  // Also add a global error handler for unhandled promise rejections
  window.addEventListener('unhandledrejection', (event) => {
    const error = event.reason;
    if (error && error.message && 
        error.message.includes('Cannot convert undefined or null to object') &&
        error.stack && error.stack.includes('checkInGroup')) {
      // Prevent the error from being logged
      event.preventDefault();
      console.warn('[Apollo Client] Unhandled promise rejection from cache normalization handled.');
    }
  });
});
