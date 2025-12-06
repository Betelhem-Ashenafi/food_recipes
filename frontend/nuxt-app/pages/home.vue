<template>
  <div>
    <!-- Hero Section -->
    <div class="relative bg-gray-900 overflow-hidden">
      <div class="absolute inset-0">
        <img class="w-full h-full object-cover opacity-40" src="https://images.unsplash.com/photo-1543353071-087f9bcbd111?ixlib=rb-4.0.3&auto=format&fit=crop&w=2070&q=80" alt="Food Background">
        <div class="absolute inset-0 bg-gradient-to-b from-transparent to-gray-900"></div>
      </div>
      <div class="relative max-w-7xl mx-auto py-24 px-4 sm:py-32 sm:px-6 lg:px-8">
        <h1 class="text-4xl font-extrabold tracking-tight text-white sm:text-5xl lg:text-6xl">
          Taste the Extraordinary
        </h1>
        <p class="mt-6 text-xl text-gray-300 max-w-3xl">
          Discover thousands of premium recipes from top chefs around the world. 
          Join our community to share your own culinary masterpieces and unlock exclusive content.
        </p>
        <div class="mt-10 flex space-x-4">
          <NuxtLink to="/register" class="inline-flex items-center px-6 py-3 border border-transparent text-base font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 transition-colors duration-300 shadow-lg hover:shadow-indigo-500/50">
            Get Started Free
          </NuxtLink>
          <a href="#recipes" class="inline-flex items-center px-6 py-3 border border-gray-300 text-base font-medium rounded-md text-gray-200 hover:bg-gray-800 transition-colors duration-300">
            Browse Recipes
          </a>
        </div>
      </div>
    </div>

    <!-- Recipe Feed Section -->
    <div id="recipes" class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16 bg-gray-50">
      <div class="text-center mb-12">
        <h2 class="text-3xl font-extrabold text-gray-900 sm:text-4xl">Trending Now</h2>
        <p class="mt-4 text-lg text-gray-500">The most popular dishes from our community this week.</p>
      </div>

      <!-- Loading State -->
      <div v-if="pending" class="flex justify-center py-20">
        <div class="animate-spin rounded-full h-16 w-16 border-t-4 border-b-4 border-indigo-600"></div>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="text-center py-20">
        <div class="inline-block p-4 rounded-full bg-red-100 text-red-600 mb-4">
          <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
        </div>
        <h3 class="text-lg font-medium text-gray-900">Oops! Something went wrong.</h3>
        <p class="mt-2 text-gray-500">{{ error.message }}</p>
      </div>

      <!-- Recipe Grid -->
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-10">
        <div v-for="recipe in data?.recipes" :key="recipe.id" class="group bg-white rounded-2xl shadow-sm hover:shadow-2xl transition-all duration-300 transform hover:-translate-y-1 overflow-hidden border border-gray-100">
          <!-- Image -->
          <div class="relative h-64 overflow-hidden">
            <img 
              :src="recipe.thumbnail_url || 'https://images.unsplash.com/photo-1495521821757-a1efb6729352?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80'" 
              alt="Recipe Image" 
              class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-110"
            >
            <div class="absolute top-4 right-4">
              <span v-if="recipe.price > 0" class="inline-flex items-center px-3 py-1 rounded-full text-xs font-bold bg-white/90 text-indigo-600 backdrop-blur-sm shadow-sm">
                💎 {{ recipe.price }} Credits
              </span>
              <span v-else class="inline-flex items-center px-3 py-1 rounded-full text-xs font-bold bg-green-500/90 text-white backdrop-blur-sm shadow-sm">
                Free
              </span>
            </div>
          </div>

          <!-- Content -->
          <div class="p-6">
            <div class="flex items-center mb-4">
              <div class="h-8 w-8 rounded-full bg-indigo-100 flex items-center justify-center text-indigo-600 font-bold text-xs">
                {{ recipe.user?.full_name?.charAt(0) || 'C' }}
              </div>
              <span class="ml-2 text-sm text-gray-600 font-medium">{{ recipe.user?.full_name || 'Chef' }}</span>
            </div>
            
            <h3 class="text-xl font-bold text-gray-900 mb-2 group-hover:text-indigo-600 transition-colors">
              {{ recipe.title }}
            </h3>
            <p class="text-gray-500 text-sm line-clamp-2 mb-6">
              {{ recipe.description }}
            </p>

            <NuxtLink :to="`/recipes/${recipe.id}`" class="block w-full text-center py-3 rounded-xl bg-gray-50 text-gray-900 font-semibold hover:bg-indigo-600 hover:text-white transition-colors duration-300">
              View Recipe
            </NuxtLink>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
// GraphQL Query
const query = gql`
  query GetRecipes {
    recipes(order_by: {created_at: desc}) {
      id
      title
      description
      thumbnail_url
      price
      user {
        full_name
      }
    }
  }
`;

// Fetch Data
const { data, error, pending } = await useAsyncQuery(query);
</script>
