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
          <div class="w-20 h-20 bg-emerald-500 rounded-full flex items-center justify-center mx-auto mb-4 shadow-lg">
            <svg class="w-12 h-12 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
          </div>
          <h1 class="text-3xl font-bold text-white mb-2">Payment Successful!</h1>
          <p class="text-gray-300">Your recipe purchase has been confirmed.</p>
        </div>

        <div v-if="recipe" class="mb-6 p-4 bg-emerald-500/20 border border-emerald-400/50 rounded-lg">
          <h2 class="font-semibold text-white mb-2">{{ recipe.title }}</h2>
          <p class="text-sm text-emerald-200">You now have full access to this recipe!</p>
        </div>

        <div class="space-y-3">
          <button
            @click="goToRecipe"
            class="w-full bg-gradient-to-r from-emerald-600 to-teal-600 text-white font-bold py-4 rounded-lg hover:from-emerald-500 hover:to-teal-500 transition-all shadow-lg hover:shadow-emerald-500/30 transform hover:-translate-y-0.5"
          >
            View Recipe
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
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';

const route = useRoute();
const router = useRouter();
const recipe = ref(null);
const recipeId = ref(null);

const verifyPayment = async () => {
  // Debug: Log all query params to see what Chapa is actually sending
  console.log('Payment success page - All query params:', route.query);
  console.log('Payment success page - Full URL:', window.location.href);
  
  const txRef = route.query.tx_ref || route.query.txRef || route.query.txref;
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
  if (!extractedRecipeId && route.query.recipe_id) {
    extractedRecipeId = parseInt(route.query.recipe_id);
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
    console.warn('Available query params:', Object.keys(route.query));
    return;
  }

  // Store recipe ID so we can always navigate to it
  recipeId.value = extractedRecipeId;

  // Verify payment - try with tx_ref first, then with recipe_id
  let retryCount = 0;
  const maxRetries = 5;
  
  const verifyPaymentWithRetry = async () => {
    if (retryCount >= maxRetries) {
      console.warn('⚠️ Max retries reached for payment verification');
      return false;
    }
    
    const token = useCookie('auth_token').value;
    
    // First try with tx_ref if available
    if (txRef) {
      try {
        console.log('💳 Verifying payment with tx_ref:', txRef);
        
        const response = await fetch(`http://localhost:8081/payment/verify?tx_ref=${txRef}`, {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        });

        if (response.ok) {
          const data = await response.json();
          if (data.status === 'success') {
            console.log('✅ Payment verified successfully with tx_ref! Purchase recorded.');
            return true;
          } else {
            console.warn('⚠️ Payment verification returned non-success status:', data);
          }
        } else {
          console.warn('⚠️ Payment verification with tx_ref failed:', response.status);
        }
      } catch (err) {
        console.warn('⚠️ Error verifying payment with tx_ref:', err);
      }
    }
    
    // If tx_ref verification failed or no tx_ref, try with recipe_id
    if (extractedRecipeId) {
      try {
        console.log(`💳 Verifying payment with recipe_id: ${extractedRecipeId} (attempt ${retryCount + 1}/${maxRetries})`);
        
        const response = await fetch(`http://localhost:8081/payment/verify?recipe_id=${extractedRecipeId}`, {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        });

        if (response.ok) {
          const data = await response.json();
          if (data.status === 'success') {
            console.log('✅ Payment verified successfully with recipe_id! Purchase confirmed.');
            return true;
          } else if (data.status === 'pending') {
            retryCount++;
            console.log(`⏳ Payment still processing, will retry in 3 seconds... (${retryCount}/${maxRetries})`);
            // Retry after delay with exponential backoff
            const delay = Math.min(3000 * retryCount, 10000); // Max 10 seconds
            setTimeout(async () => {
              await verifyPaymentWithRetry();
            }, delay);
            return false;
          } else {
            console.warn('⚠️ Payment verification returned:', data);
          }
        } else {
          const errorText = await response.text().catch(() => 'Unknown error');
          console.warn('⚠️ Payment verification with recipe_id failed:', response.status, errorText);
        }
      } catch (err) {
        console.warn('⚠️ Error verifying payment with recipe_id:', err);
      }
    }
    
    return false;
  };
  
  // Start verification
  await verifyPaymentWithRetry();

  // Fetch recipe details to show title
  try {
    const recipeResponse = await fetch(`http://localhost:8081/recipes/${extractedRecipeId}`);
    if (recipeResponse.ok) {
      recipe.value = await recipeResponse.json();
      // Ensure recipe has id property
      if (!recipe.value.id) {
        recipe.value.id = extractedRecipeId;
      }
    } else {
      // If recipe fetch fails, just store ID
      recipe.value = { id: extractedRecipeId, title: 'Recipe' };
    }
  } catch (err) {
    console.error('Error fetching recipe:', err);
    recipe.value = { id: extractedRecipeId, title: 'Recipe' };
  }
};

const goToRecipe = () => {
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

onMounted(() => {
  verifyPayment();
});
</script>


