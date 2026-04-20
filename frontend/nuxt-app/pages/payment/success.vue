<template>
  <div class="relative min-h-screen pb-20">
    <!-- Background Image with Overlay - Consistent with other pages -->
    <div class="fixed inset-0 z-0">
      <img 
        src="https://images.unsplash.com/photo-1504674900247-0877df9cc836?ixlib=rb-4.0.3&auto=format&fit=crop&w=2070&q=80" 
        alt="Elegant Food" 
        class="w-full h-full object-cover brightness-60"
      >
      <div class="absolute inset-0 bg-black/80"></div>
    </div>

    <!-- Content -->
    <div class="relative z-10 flex items-center justify-center min-h-screen px-4 py-16">
      <div class="max-w-md w-full bg-white/10 backdrop-blur-lg border border-white/20 rounded-2xl shadow-2xl p-8 text-center">
        <div class="mb-6">
          <div class="w-20 h-20 rounded-full flex items-center justify-center mx-auto mb-4 shadow-lg"
               :class="statusBadgeClass">
            <svg class="w-12 h-12 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
          </div>
          <h1 class="text-3xl font-bold text-white mb-2">{{ statusTitle }}</h1>
          <p class="text-gray-300">{{ statusDescription }}</p>
        </div>

        <div v-if="recipe" class="mb-6 p-4 bg-emerald-500/20 border border-emerald-400/50 rounded-lg">
          <h2 class="font-semibold text-white mb-2">{{ recipe.title }}</h2>
          <p class="text-sm text-emerald-200">You now have full access to this recipe!</p>
        </div>

        <div class="space-y-3">
          <button
            @click="goToRecipe"
            :disabled="!canOpenRecipe"
            class="w-full bg-gradient-to-r from-emerald-600 to-teal-600 text-white font-bold py-4 rounded-lg hover:from-emerald-500 hover:to-teal-500 transition-all shadow-lg hover:shadow-emerald-500/30 transform hover:-translate-y-0.5"
            :class="{ 'opacity-60 cursor-not-allowed hover:from-emerald-600 hover:to-teal-600 transform-none': !canOpenRecipe }"
          >
            {{ canOpenRecipe ? 'View Recipe' : 'Waiting For Confirmation' }}
          </button>
          <button
            v-if="paymentStatus === 'pending'"
            @click="recheckNow"
            class="w-full bg-amber-500/20 border border-amber-400/50 text-amber-100 font-semibold py-4 rounded-lg hover:bg-amber-500/30 transition-all"
          >
            Check Payment Again
          </button>
          <button
            @click="goToHome"
            class="w-full bg-white/10 border border-white/30 text-white font-semibold py-4 rounded-lg hover:bg-white/20 transition-all"
          >
            Back to Home
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue';
import { useApolloClient } from '@vue/apollo-composable';
import gql from 'graphql-tag';
import { useRouter, useRoute } from '#app';
import { verifyPaymentAction } from '~/utils/payment-actions';

const { client } = useApolloClient();
const recipe = ref(null);
const recipeId = ref(null);
const router = useRouter();
const route = useRoute();
const pollTimer = ref(null);
const verifyAttempts = ref(0);
const maxVerifyAttempts = 30;

const firstQueryValue = (query, key) => {
  const value = query?.[key];
  if (Array.isArray(value)) return value[0];
  return value;
};

const normalizeQueryObject = () => {
  const normalized = {};

  // Start from Nuxt route query.
  Object.entries(route.query || {}).forEach(([rawKey, rawValue]) => {
    const key = String(rawKey || '').replace(/^amp;/, '');
    if (!key) return;
    if (normalized[key] === undefined) {
      normalized[key] = Array.isArray(rawValue) ? rawValue[0] : rawValue;
    }
  });

  // Parse raw search to recover keys from URLs containing '&amp;'.
  if (process.client) {
    const rawSearch = String(window.location.search || '').replace(/&amp;/g, '&');
    const params = new URLSearchParams(rawSearch.startsWith('?') ? rawSearch.slice(1) : rawSearch);
    for (const [rawKey, rawValue] of params.entries()) {
      const key = String(rawKey || '').replace(/^amp;/, '');
      if (!key) continue;
      if (normalized[key] === undefined || normalized[key] === '') {
        normalized[key] = rawValue;
      }
    }
  }

  return normalized;
};

const normalizeStatus = (rawStatus) => {
  const status = String(rawStatus || '').toLowerCase().trim();
  if (['success', 'completed', 'paid'].includes(status)) return 'success';
  if (['pending', 'processing', 'created', 'initiated'].includes(status)) return 'pending';
  if (status === 'failed') return 'failed';
  return status || 'pending';
};

const resolvedStatus = ref(normalizeStatus(route.query.status || 'pending'));
const resolvedMessage = ref(String(route.query.message || '').trim());

const setStatus = (nextStatus, message = '') => {
  const next = normalizeStatus(nextStatus);
  // Once confirmed successful, do not downgrade to pending/failed.
  if (resolvedStatus.value === 'success' && next !== 'success') {
    return;
  }
  resolvedStatus.value = next || resolvedStatus.value;
  if (String(message || '').trim()) {
    resolvedMessage.value = String(message).trim();
  }
};

const updateUrlStatus = async (nextStatus, message = '') => {
  const current = { ...normalizeQueryObject() };
  current.status = normalizeStatus(nextStatus || current.status || 'pending');
  if (String(message || '').trim()) {
    current.message = String(message).trim();
  }
  try {
    await router.replace({ path: route.path, query: current });
  } catch {
    // URL sync is best-effort only.
  }
};

const stopPolling = () => {
  if (pollTimer.value) {
    clearTimeout(pollTimer.value);
    pollTimer.value = null;
  }
};

const scheduleNextVerify = (fn, delayMs = 4000) => {
  stopPolling();
  pollTimer.value = setTimeout(async () => {
    await fn();
  }, delayMs);
};

const paymentStatus = computed(() => resolvedStatus.value || 'pending');
const canOpenRecipe = computed(() => paymentStatus.value === 'success');
const statusTitle = computed(() => {
  if (paymentStatus.value === 'pending') return 'Payment Pending';
  if (paymentStatus.value === 'failed') return 'Payment Failed';
  return 'Payment Successful!';
});
const statusDescription = computed(() => {
  if (resolvedMessage.value) return resolvedMessage.value;
  if (paymentStatus.value === 'pending') return 'Your payment is being processed. Please wait and refresh later.';
  if (paymentStatus.value === 'failed') return 'Payment was not completed. You can return and resume your payment.';
  return 'Your recipe purchase has been confirmed.';
});
const statusBadgeClass = computed(() => {
  if (paymentStatus.value === 'pending') return 'bg-amber-500';
  if (paymentStatus.value === 'failed') return 'bg-red-500';
  return 'bg-emerald-500';
});

const RECIPE_TITLE_QUERY = gql`
  query GetRecipeTitle($id: Int!) {
    recipes_by_pk(id: $id) {
      id
      title
    }
  }
`;

const verifyPayment = async () => {
  const normalizedQuery = normalizeQueryObject();

  // Debug: Log all query params to see what Chapa is actually sending
  console.log('Payment success page - All query params:', route.query);
  console.log('Payment success page - Full URL:', window.location.href);
  console.log('Payment success page - Normalized query params:', normalizedQuery);
  
  const txRef =
    firstQueryValue(normalizedQuery, 'tx_ref') ||
    firstQueryValue(normalizedQuery, 'txRef') ||
    firstQueryValue(normalizedQuery, 'txref');
  let extractedRecipeId = null;

  // Try to get recipe ID from tx_ref (format: tx-{recipeId}-{timestamp})
  if (txRef) {
    console.log('Found tx_ref:', txRef);
    const parts = txRef.split('-');
    console.log('tx_ref parts:', parts);
    if (parts.length >= 2) {
      extractedRecipeId = parseInt(parts[1]);
      console.log('Extracted recipe ID from tx_ref:', extractedRecipeId);
    }
  }

  // Also check if recipe_id is in query params (fallback)
  if (!extractedRecipeId && firstQueryValue(normalizedQuery, 'recipe_id')) {
    extractedRecipeId = parseInt(firstQueryValue(normalizedQuery, 'recipe_id'));
    console.log('Found recipe_id in query params:', extractedRecipeId);
  }

  // Check if recipe_id is in the URL hash or other places
  if (!extractedRecipeId) {
    // Try to extract from URL hash if Chapa puts it there
    const hash = window.location.hash;
    if (hash) {
      const hashParams = new URLSearchParams(hash.substring(1));
      const hashRecipeId = hashParams.get('recipe_id');
      if (hashRecipeId) {
        extractedRecipeId = parseInt(hashRecipeId);
        console.log('Found recipe_id in hash:', extractedRecipeId);
      }
    }
  }

  // Fallback: Check sessionStorage (stored when user clicks "Buy Recipe")
  if (!extractedRecipeId && process.client) {
    const storedRecipeId = sessionStorage.getItem('pending_payment_recipe_id');
    if (storedRecipeId) {
      extractedRecipeId = parseInt(storedRecipeId);
      console.log('Found recipe_id in sessionStorage:', extractedRecipeId);
      // Clear it after use
      sessionStorage.removeItem('pending_payment_recipe_id');
    }
  }

  if (!extractedRecipeId) {
    // No recipe ID found - user might have navigated directly
    console.warn('No recipe ID found in payment success page');
    console.warn('Available query params:', Object.keys(normalizedQuery));
    resolvedStatus.value = 'failed';
    if (!resolvedMessage.value) setStatus('failed', 'Could not determine the recipe for this payment.');
    return;
  }

  // Rewrite malformed amp-prefixed query params to canonical keys on first load.
  await updateUrlStatus(paymentStatus.value, resolvedMessage.value);

  // Store recipe ID so we can always navigate to it
  recipeId.value = extractedRecipeId;

  // Verify payment - try with tx_ref first, then with recipe_id
  const verifyPaymentWithRetry = async () => {
    if (verifyAttempts.value >= maxVerifyAttempts) {
      setStatus('pending', 'Payment is still pending. Please keep this page open and try again shortly.');
      await updateUrlStatus('pending', 'Payment is still pending. Please keep this page open and try again shortly.');
      return false;
    }
    verifyAttempts.value += 1;
    
    // First try with tx_ref if available
    if (txRef) {
      try {
        console.log('💳 Verifying payment with tx_ref:', txRef);

        const data = await verifyPaymentAction(client, { tx_ref: txRef });
        if (data?.status === 'success') {
          console.log('✅ Payment verified successfully with tx_ref! Purchase recorded.');
          setStatus('success', String(data?.message || 'Your recipe purchase has been confirmed.'));
          await updateUrlStatus('success', String(data?.message || 'Your recipe purchase has been confirmed.'));
          stopPolling();
          return true;
        } else if (data?.status === 'pending') {
          setStatus('pending', String(data?.message || 'Payment is still processing.'));
          await updateUrlStatus('pending', String(data?.message || 'Payment is still processing.'));
        } else {
          setStatus('failed', String(data?.message || 'Payment was not completed.'));
          await updateUrlStatus('failed', String(data?.message || 'Payment was not completed.'));
          stopPolling();
          console.warn('⚠️ Payment verification returned non-success status:', data);
          return false;
        }
      } catch (err) {
        console.warn('⚠️ Error verifying payment with tx_ref:', err);
      }
    }
    
    // If tx_ref verification failed or no tx_ref, try with recipe_id
    if (extractedRecipeId) {
      try {
        console.log(`💳 Verifying payment with recipe_id: ${extractedRecipeId} (attempt ${verifyAttempts.value}/${maxVerifyAttempts})`);

        const data = await verifyPaymentAction(client, { recipe_id: extractedRecipeId });
        if (data?.status === 'success') {
          console.log('✅ Payment verified successfully with recipe_id! Purchase confirmed.');
          setStatus('success', String(data?.message || 'Your recipe purchase has been confirmed.'));
          await updateUrlStatus('success', String(data?.message || 'Your recipe purchase has been confirmed.'));
          stopPolling();
          return true;
        } else if (data?.status === 'pending') {
          setStatus('pending', String(data?.message || 'Payment is still processing.'));
          await updateUrlStatus('pending', String(data?.message || 'Payment is still processing.'));
          console.log(`⏳ Payment still processing, retrying on same URL... (${verifyAttempts.value}/${maxVerifyAttempts})`);
          const delay = Math.min(2000 + verifyAttempts.value * 300, 7000);
          scheduleNextVerify(verifyPaymentWithRetry, delay);
          return false;
        } else {
          setStatus('failed', String(data?.message || 'Payment was not completed.'));
          await updateUrlStatus('failed', String(data?.message || 'Payment was not completed.'));
          stopPolling();
          console.warn('⚠️ Payment verification returned:', data);
          return false;
        }
      } catch (err) {
        setStatus('pending', 'Verification is retrying. Please keep this page open.');
        await updateUrlStatus('pending', 'Verification is retrying. Please keep this page open.');
        console.warn('⚠️ Error verifying payment with recipe_id:', err);
        scheduleNextVerify(verifyPaymentWithRetry, 5000);
        return false;
      }
    }
    
    return false;
  };
  
  // Start verification
  await verifyPaymentWithRetry();

  // Fetch recipe details to show title
  try {
    const result = await client.query({
      query: RECIPE_TITLE_QUERY,
      variables: { id: extractedRecipeId },
      fetchPolicy: 'network-only'
    });
    recipe.value = result?.data?.recipes_by_pk || { id: extractedRecipeId, title: 'Recipe' };
  } catch (err) {
    console.error('Error fetching recipe:', err);
    recipe.value = { id: extractedRecipeId, title: 'Recipe' };
  }
};

const goToRecipe = () => {
  if (!canOpenRecipe.value) {
    return;
  }
  // Always use recipeId.value if available, fallback to recipe.value?.id
  const idToUse = recipeId.value || recipe.value?.id;
  if (idToUse) {
    console.log('Navigating to recipe:', idToUse);
    // Add query param to indicate coming from payment
    router.push(`/recipes/${idToUse}?fromPayment=true&payment=success`);
  } else {
    console.warn('No recipe ID available, navigating to home');
    router.push('/home');
  }
};

const goToHome = () => {
  router.push('/home');
};

const recheckNow = async () => {
  verifyAttempts.value = 0;
  stopPolling();
  await verifyPayment();
};

onMounted(async () => {
  try {
    await verifyPayment();
  } catch (err) {
    console.error('Error in payment success page:', err);
  }
});

onBeforeUnmount(() => {
  stopPolling();
});
</script>


