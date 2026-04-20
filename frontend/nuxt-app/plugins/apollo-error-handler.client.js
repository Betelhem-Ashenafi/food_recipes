/**
 * Lightweight runtime error logger for known client-side issues.
 *
 * Avoid mutating global built-ins (like Object.keys), which can introduce
 * hard-to-debug side effects in framework/runtime code.
 */
export default defineNuxtPlugin((nuxtApp) => {
  if (!process.client) return;

  // Surface a concise hint when this specific error appears.
  window.addEventListener('unhandledrejection', (event) => {
    const error = event.reason;
    if (
      error &&
      error.message &&
      error.message.includes('Cannot convert undefined or null to object') &&
      error.stack &&
      error.stack.includes('checkInGroup')
    ) {
      console.warn('[Client] checkInGroup null-object error detected.');
    }
  });
});

