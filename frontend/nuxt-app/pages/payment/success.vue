<template>
  <div class="min-h-screen bg-gradient-to-br from-emerald-50 via-teal-50 to-cyan-50 flex items-center justify-center p-4">
    <div class="max-w-md w-full bg-white rounded-2xl shadow-2xl p-8 text-center">
      <div class="mb-6">
        <div class="w-20 h-20 bg-emerald-500 rounded-full flex items-center justify-center mx-auto mb-4">
          <svg class="w-12 h-12 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
        </div>
        <h1 class="text-3xl font-bold text-gray-900 mb-2">Payment Successful!</h1>
        <p class="text-gray-600">Your recipe purchase has been confirmed.</p>
      </div>

      <div v-if="recipe" class="mb-6 p-4 bg-emerald-50 rounded-lg">
        <h2 class="font-semibold text-gray-900 mb-2">{{ recipe.title }}</h2>
        <p class="text-sm text-gray-600">You now have full access to this recipe!</p>
      </div>

      <div class="space-y-3">
        <button
          @click="goToRecipe"
          class="w-full bg-gradient-to-r from-emerald-600 to-teal-600 text-white font-semibold py-3 rounded-lg hover:from-emerald-500 hover:to-teal-500 transition-all"
        >
          View Recipe
        </button>
        <button
          @click="goToHome"
          class="w-full bg-gray-200 text-gray-700 font-semibold py-3 rounded-lg hover:bg-gray-300 transition-all"
        >
          Back to Home
        </button>
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

const verifyPayment = async () => {
  const txRef = route.query.tx_ref;
  if (!txRef) {
    console.error('No transaction reference');
    return;
  }

  try {
    const token = useCookie('auth_token').value;
    const response = await fetch(`http://localhost:8081/payment/verify?tx_ref=${txRef}`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });

    if (response.ok) {
      const data = await response.json();
      if (data.status === 'success') {
        // Extract recipe ID from tx_ref (format: tx-{recipeId}-{timestamp})
        const parts = txRef.split('-');
        if (parts.length >= 2) {
          const recipeId = parseInt(parts[1]);
          // Fetch recipe details
          // For now, just store recipe ID
          recipe.value = { id: recipeId };
        }
      }
    }
  } catch (err) {
    console.error('Error verifying payment:', err);
  }
};

const goToRecipe = () => {
  if (recipe.value?.id) {
    router.push(`/recipes/${recipe.value.id}`);
  } else {
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

