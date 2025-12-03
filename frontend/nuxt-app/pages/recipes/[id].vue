<template>
  <div v-if="fetching" class="text-center">Loading...</div>
  <div v-else-if="error" class="text-red-500">Error: {{ error.message }}</div>
  <div v-else class="bg-white rounded-lg shadow-md overflow-hidden">
    <img :src="data.recipes_by_pk.thumbnail_url || 'https://via.placeholder.com/800x400'" alt="Recipe Image" class="w-full h-64 object-cover">
    <div class="p-8">
      <h1 class="text-3xl font-bold mb-4">{{ data.recipes_by_pk.title }}</h1>
      <p class="text-gray-600 mb-6">{{ data.recipes_by_pk.description }}</p>
      <div class="text-2xl font-bold text-green-600 mb-6">${{ data.recipes_by_pk.price }}</div>
      
      <button @click="buyRecipe" class="bg-green-500 text-white px-6 py-3 rounded-lg text-lg hover:bg-green-600">
        Buy Recipe
      </button>
    </div>
  </div>
</template>

<script setup>
import { useRoute } from 'vue-router';
import { useQuery } from '@urql/vue';

const route = useRoute();
const recipeId = route.params.id;

const { data, fetching, error } = useQuery({
  query: `
    query GetRecipe($id: Int!) {
      recipes_by_pk(id: $id) {
        id
        title
        description
        thumbnail_url
        price
      }
    }
  `,
  variables: { id: parseInt(recipeId) }
});

const buyRecipe = async () => {
  try {
    const token = localStorage.getItem('token');
    if (!token) {
      alert('Please login to purchase');
      return;
    }

    const response = await fetch('http://localhost:8081/payment/initialize', {
      method: 'POST',
      headers: { 
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify({ 
        amount: data.value.recipes_by_pk.price,
        email: "user@example.com", // Should come from user profile
        first_name: "User",
        last_name: "Name",
        tx_ref: `tx-${Date.now()}`
      })
    });

    const result = await response.json();
    if (result.checkout_url) {
      window.location.href = result.checkout_url;
    } else {
      alert('Payment initialization failed');
    }
  } catch (e) {
    console.error(e);
    alert('Error processing payment');
  }
};
</script>
