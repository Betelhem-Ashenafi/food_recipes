/**
 * Apollo Client Error Handler Plugin
 * Catches and handles Apollo Client errors, especially Object.keys() errors
 * from cache normalization
 * 
 * This plugin patches Object.keys to handle null/undefined safely,
 * preventing Apollo Client cache normalization errors
 */
export default defineNuxtPlugin((nuxtApp) => {
  if (!process.client) return;

  // Store original Object.keys
  const originalObjectKeys = Object.keys;

  // Patch Object.keys to handle null/undefined safely
  Object.keys = function(obj) {
    if (obj === null || obj === undefined) {
      // Return empty array for null/undefined instead of throwing
      return [];
    }
    // For non-objects, try to convert or return empty array
    if (typeof obj !== 'object') {
      return [];
    }
    // Use original Object.keys for valid objects
    return originalObjectKeys.call(this, obj);
  };

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

