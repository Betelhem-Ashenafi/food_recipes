<template>
  <div>
    <h1 class="text-3xl font-bold mb-6">Delicious Recipes</h1>
    <div v-if="fetching" class="text-center">Loading...</div>
    <div v-else-if="error" class="text-red-500">Error: {{ error.message }}</div>
    <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div v-for="recipe in data.recipes" :key="recipe.id" class="bg-white rounded-lg shadow-md overflow-hidden">
        <img :src="recipe.thumbnail_url || 'https://via.placeholder.com/300'" alt="Recipe Image" class="w-full h-48 object-cover">
        <div class="p-4">
          <h2 class="text-xl font-semibold mb-2">{{ recipe.title }}</h2>
          <p class="text-gray-600 mb-4">{{ recipe.description }}</p>
          <NuxtLink :to="`/recipes/${recipe.id}`" class="bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600">View Recipe</NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useQuery } from '@urql/vue';

const { data, fetching, error } = useQuery({
  query: `
    query GetRecipes {
      recipes {
        id
        title
        description
        thumbnail_url
        price
      }
    }
  `
});
</script>
