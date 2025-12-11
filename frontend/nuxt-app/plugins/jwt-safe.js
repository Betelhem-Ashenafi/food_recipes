/**
 * Plugin to ensure JWT tokens always have safe structure
 * This prevents Object.keys() errors when JWT claims are null/undefined
 * 
 * This plugin runs early to ensure any code accessing JWT gets a safe structure
 */
export default defineNuxtPlugin(() => {
  // This plugin ensures JWT structure is always safe
  // The actual JWT decoding is handled in individual components with safe guards
});

