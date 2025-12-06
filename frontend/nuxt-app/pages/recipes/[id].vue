<template>
  <div class="container mx-auto px-4 py-8">
    <div v-if="pending" class="text-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600 mx-auto"></div>
      <p class="mt-4 text-gray-500">Loading recipe...</p>
    </div>

    <div v-else-if="error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative">
      <strong class="font-bold">Error!</strong>
      <span class="block sm:inline"> {{ error.message }}</span>
    </div>

    <div v-else-if="data?.recipes_by_pk" class="bg-white rounded-lg shadow-lg overflow-hidden">
      <img 
        :src="data.recipes_by_pk.thumbnail_url || 'https://via.placeholder.com/800x400?text=No+Image'" 
        alt="Recipe Image" 
        class="w-full h-96 object-cover"
      >
      <div class="p-8">
        <div class="flex justify-between items-start mb-6">
          <h1 class="text-4xl font-bold text-gray-900">{{ data.recipes_by_pk.title }}</h1>
          <span v-if="data.recipes_by_pk.price > 0" class="bg-green-100 text-green-800 text-lg font-semibold px-4 py-1 rounded">
            {{ data.recipes_by_pk.price }} Credits
          </span>
          <span v-else class="bg-blue-100 text-blue-800 text-lg font-semibold px-4 py-1 rounded">
            Free
          </span>
        </div>

        <div class="flex items-center mb-8 text-gray-600">
          <span class="mr-4">By <span class="font-medium text-gray-900">{{ data.recipes_by_pk.user?.full_name || 'Unknown Chef' }}</span></span>
          <span>• {{ new Date(data.recipes_by_pk.created_at).toLocaleDateString() }}</span>
        </div>

        <p class="text-gray-700 text-lg mb-8 leading-relaxed">{{ data.recipes_by_pk.description }}</p>
        
        <div class="border-t border-gray-200 pt-8">
          <h2 class="text-2xl font-bold mb-4">Ingredients & Steps</h2>
          <div class="bg-gray-50 p-6 rounded-lg text-center">
            <p class="text-gray-600 mb-4">Unlock this recipe to see the full ingredients list and preparation steps.</p>
            <button @click="buyRecipe" class="bg-indigo-600 text-white px-8 py-3 rounded-lg text-lg font-medium hover:bg-indigo-700 transition-colors shadow-md">
              Unlock for {{ data.recipes_by_pk.price }} Credits
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <div v-else class="text-center py-12 text-gray-500">
      Recipe not found.
    </div>
  </div>
</template>

<script setup>
const route = useRoute();
const recipeId = route.params.id;

const query = gql`
  query GetRecipe($id: Int!) {
    recipes_by_pk(id: $id) {
      id
      title
      description
      price
      thumbnail_url
      created_at
      user {
        full_name
      }
    }
  }
`;

const { data, error, pending } = await useAsyncQuery(query, { id: recipeId });

const buyRecipe = () => {
  alert('Purchase functionality coming in Slice 3!');
};
</script>
