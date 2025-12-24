<template>
  <div class="relative min-h-screen pb-20">
    <!-- Background Image with Overlay - Sophisticated Food Presentation (Profile Page) -->
    <div class="fixed inset-0 z-0">
      <img 
        src="https://images.unsplash.com/photo-1495521821757-a1efb6729352?ixlib=rb-4.0.3&auto=format&fit=crop&w=2070&q=80" 
        alt="Elegant Food Presentation" 
        class="w-full h-full object-cover brightness-95"
      >
      <div class="absolute inset-0 bg-black/80"></div>
    </div>

    <!-- Content -->
    <div class="relative z-10 max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 pt-16 pb-16">
      <!-- Profile Header -->
      <div class="bg-white/10 backdrop-blur-lg border border-white/20 rounded-2xl p-8 mb-8 shadow-2xl hover:bg-white/15 transition-all duration-300">
        <div class="flex items-center gap-6">
          <div class="h-24 w-24 rounded-full bg-gradient-to-br from-emerald-400 to-teal-500 flex items-center justify-center text-white font-bold text-4xl shadow-lg">
            {{ userName.charAt(0) }}
          </div>
          <div>
            <h1 class="text-3xl font-bold text-white mb-2">{{ userName }}</h1>
            <p class="text-gray-300">{{ userEmail }}</p>
          </div>
        </div>
      </div>

      <!-- Tabs -->
      <div class="bg-white/10 backdrop-blur-lg border border-white/20 rounded-2xl overflow-hidden shadow-2xl">
        <div class="flex border-b border-white/20">
          <button
            v-for="tab in tabs"
            :key="tab.id"
            @click="activeTab = tab.id"
            :class="[
              'flex-1 px-6 py-4 font-semibold transition-colors',
              activeTab === tab.id
                ? 'bg-emerald-500/30 text-emerald-300 border-b-2 border-emerald-400'
                : 'text-gray-400 hover:text-white hover:bg-white/5'
            ]"
          >
            {{ tab.label }}
          </button>
        </div>

        <div class="p-8">
          <!-- My Recipes -->
          <div v-if="activeTab === 'recipes'">
            <h2 class="text-2xl font-bold text-white mb-6">My Recipes</h2>
            <div v-if="myRecipes.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              <div
                v-for="recipe in myRecipes"
                :key="recipe.id"
                class="bg-white/10 backdrop-blur-lg border border-white/20 rounded-2xl overflow-hidden hover:bg-white/15 hover:border-emerald-400/50 transition-all cursor-pointer group relative"
                @click="navigateToRecipe(recipe.id)"
              >
                <img :src="recipe.thumbnail_url || getDefaultImage(recipe.title)" :alt="recipe.title" class="w-full h-48 object-cover" />
                <div class="p-4">
                  <div class="flex items-center justify-between mb-2">
                    <h3 class="text-xl font-bold text-white group-hover:text-emerald-400 transition-colors">{{ recipe.title }}</h3>
                    <div class="flex gap-2">
                      <button
                        @click.stop.prevent="editRecipe(recipe.id, $event)"
                        type="button"
                        class="p-2 bg-emerald-500/20 hover:bg-emerald-500/30 rounded-lg transition-colors z-10 relative"
                        title="Edit Recipe"
                      >
                        <svg class="w-4 h-4 text-emerald-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                        </svg>
                      </button>
                      <button
                        @click.stop="deleteRecipe(recipe.id)"
                        class="p-2 bg-red-500/20 hover:bg-red-500/30 rounded-lg transition-colors"
                        title="Delete Recipe"
                      >
                        <svg class="w-4 h-4 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                        </svg>
                      </button>
                    </div>
                  </div>
                  <p class="text-gray-300 text-sm line-clamp-2 mb-2">{{ recipe.description }}</p>
                  <div class="flex items-center gap-4 text-xs text-gray-400">
                    <span>⏱️ {{ recipe.preparation_time }} min</span>
                    <span v-if="recipe.price > 0">💎 {{ recipe.price }} Credits</span>
                    <span v-else>🆓 Free</span>
                  </div>
                </div>
              </div>
            </div>
            <p v-else class="text-gray-400 text-center py-8">You haven't created any recipes yet.</p>
          </div>

          <!-- Bookmarked -->
          <div v-if="activeTab === 'bookmarked'">
            <h2 class="text-2xl font-bold text-white mb-6">Bookmarked Recipes</h2>
            <div v-if="bookmarkedRecipes.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              <div
                v-for="recipe in bookmarkedRecipes"
                :key="recipe.id"
                @click="navigateToRecipe(recipe.id)"
                class="bg-white/10 backdrop-blur-lg border border-white/20 rounded-2xl overflow-hidden hover:bg-white/15 hover:border-emerald-400/50 transition-all cursor-pointer"
              >
                <img :src="recipe.thumbnail_url || getDefaultImage(recipe.title)" :alt="recipe.title" class="w-full h-48 object-cover" />
                <div class="p-4">
                  <h3 class="text-xl font-bold text-white mb-2">{{ recipe.title }}</h3>
                  <p class="text-gray-300 text-sm line-clamp-2">{{ recipe.description }}</p>
                </div>
              </div>
            </div>
            <p v-else class="text-gray-400 text-center py-8">You haven't bookmarked any recipes yet.</p>
          </div>

          <!-- Purchased -->
          <div v-if="activeTab === 'purchased'">
            <h2 class="text-2xl font-bold text-white mb-6">Purchased Recipes</h2>
            <div v-if="purchasedRecipes && purchasedRecipes.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              <div
                v-for="recipe in purchasedRecipes"
                :key="recipe.id"
                @click="navigateToRecipe(recipe.id)"
                class="bg-white/10 backdrop-blur-lg border border-emerald-400/50 rounded-2xl overflow-hidden hover:bg-white/15 transition-all cursor-pointer"
              >
                <img :src="recipe.thumbnail_url || getDefaultImage(recipe.title)" :alt="recipe.title" class="w-full h-48 object-cover" />
                <div class="p-4">
                  <div class="flex items-center justify-between mb-2">
                    <h3 class="text-xl font-bold text-white">{{ recipe.title }}</h3>
                    <span class="text-emerald-400 text-xs font-bold">💎 Purchased</span>
                  </div>
                  <p class="text-gray-300 text-sm line-clamp-2">{{ recipe.description }}</p>
                </div>
              </div>
            </div>
            <p v-else class="text-gray-400 text-center py-8">You haven't purchased any recipes yet.</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { useQuery } from '@vue/apollo-composable';
import gql from 'graphql-tag';
import { jwtDecode } from 'jwt-decode';
import { useRouter } from 'vue-router';

const router = useRouter();

// Get user from JWT - always return safe object structure
const token = useCookie('auth_token');
const userInfo = computed(() => {
  if (!token.value) return {};
  try {
    const decoded = jwtDecode(token.value);
    // Ensure we always return an object with safe structure
    if (!decoded || typeof decoded !== 'object') return {};
    // Ensure claims is always an object, not null
    if (decoded['https://hasura.io/jwt/claims']) {
      if (typeof decoded['https://hasura.io/jwt/claims'] !== 'object' || decoded['https://hasura.io/jwt/claims'] === null) {
        decoded['https://hasura.io/jwt/claims'] = {};
      }
    } else {
      decoded['https://hasura.io/jwt/claims'] = {};
    }
    return decoded;
  } catch {
    return {};
  }
});

const userName = computed(() => {
  // First try localStorage (most reliable, stored during login)
  if (process.client) {
    const storedName = localStorage.getItem('user_name');
    if (storedName && typeof storedName === 'string' && storedName.trim() !== '') {
      return storedName;
    }
  }
  
  // Then try to get name from top-level JWT claims
  if (userInfo.value?.name && typeof userInfo.value.name === 'string' && userInfo.value.name.trim() !== '') {
    return userInfo.value.name;
  }
  
  // Then try Hasura claims
  const claims = userInfo.value?.['https://hasura.io/jwt/claims'];
  if (claims && typeof claims === 'object' && claims !== null) {
    if (claims['x-hasura-user-name'] && typeof claims['x-hasura-user-name'] === 'string' && claims['x-hasura-user-name'].trim() !== '') {
      return claims['x-hasura-user-name'];
    }
  }
  
  // Last resort fallback
  return 'User';
});
const userEmail = computed(() => {
  const claims = userInfo.value?.['https://hasura.io/jwt/claims'];
  if (!claims || typeof claims !== 'object' || claims === null) return '';
  return claims['x-hasura-user-email'] || '';
});
const userId = computed(() => {
  const claims = userInfo.value?.['https://hasura.io/jwt/claims'];
  if (!claims || typeof claims !== 'object' || claims === null) return 0;
  return parseInt(claims['x-hasura-user-id']) || 0;
});

// Tabs
const activeTab = ref('recipes');
const tabs = [
  { id: 'recipes', label: 'My Recipes' },
  { id: 'bookmarked', label: 'Bookmarked' },
  { id: 'purchased', label: 'Purchased' }
];

// GraphQL Query for User's Recipes (using Vue Apollo with Hasura)
const myRecipesQuery = gql`
  query GetUserRecipes($userId: Int!) {
    recipes(where: { user_id: { _eq: $userId } }, order_by: { created_at: desc }) {
      id
      title
      description
      thumbnail_url
      price
      preparation_time
      created_at
    }
  }
`;

// Skip query if userId is invalid to prevent Apollo errors
const { result: recipesResult, refetch: refetchMyRecipes } = useQuery(
  myRecipesQuery, 
  () => ({ userId: userId.value }),
  { 
    skip: () => !userId.value || userId.value === 0,
    errorPolicy: 'all',
    fetchPolicy: 'cache-and-network',
    returnPartialData: true
  }
);

const myRecipes = computed(() => {
  if (!recipesResult.value || typeof recipesResult.value !== 'object') return [];
  if (!recipesResult.value.recipes || !Array.isArray(recipesResult.value.recipes)) return [];
  return recipesResult.value.recipes;
});

// Function to refetch recipes
const fetchMyRecipes = async () => {
  if (refetchMyRecipes) {
    await refetchMyRecipes();
  }
};

// Fetch Bookmarked Recipes
const bookmarkedRecipes = ref([]);
const fetchBookmarkedRecipes = async () => {
  if (!token.value || !userId.value) {
    bookmarkedRecipes.value = [];
    return;
  }
  try {
    const response = await fetch(`http://localhost:8081/users/${userId.value}/bookmarks`, {
      headers: { 'Authorization': `Bearer ${token.value}` }
    });
    if (response.ok) {
      const data = await response.json();
      bookmarkedRecipes.value = data || [];
    } else if (response.status === 401) {
      console.warn('⚠️ Token invalid for bookmarks - please log in again');
      token.value = null;
      bookmarkedRecipes.value = [];
    } else {
      bookmarkedRecipes.value = [];
    }
  } catch (err) {
    console.error('Error fetching bookmarks:', err);
    bookmarkedRecipes.value = [];
  }
};

// Fetch Purchased Recipes
const purchasedRecipes = ref([]);
const fetchPurchasedRecipes = async () => {
  if (!token.value || !userId.value) {
    purchasedRecipes.value = [];
    return;
  }
  try {
    const response = await fetch(`http://localhost:8081/users/${userId.value}/purchases`, {
      headers: { 'Authorization': `Bearer ${token.value}` }
    });
    if (response.ok) {
      const data = await response.json();
      purchasedRecipes.value = data || [];
    } else if (response.status === 401) {
      console.warn('⚠️ Token invalid for purchases - please log in again');
      token.value = null;
      purchasedRecipes.value = [];
    } else {
      purchasedRecipes.value = [];
    }
  } catch (err) {
    console.error('Error fetching purchases:', err);
    purchasedRecipes.value = [];
  }
};

// Navigate to recipe
const navigateToRecipe = (recipeId) => {
  router.push(`/recipes/${recipeId}`);
};

// Get default image
const getDefaultImage = (title) => {
  const t = title?.toLowerCase() || '';
  if (t.includes('avocado') || t.includes('salad')) {
    return 'https://images.unsplash.com/photo-1546069901-ba9599a7e63c?auto=format&fit=crop&w=800&q=80';
  }
  if (t.includes('chocolate') || t.includes('cake')) {
    return 'https://images.unsplash.com/photo-1578985545062-69928b1d9587?auto=format&fit=crop&w=800&q=80';
  }
  return 'https://images.unsplash.com/photo-1495521821757-a1efb6729352?auto=format&fit=crop&w=800&q=80';
};

// Edit Recipe
const editRecipe = (recipeId, event) => {
  if (event) {
    event.preventDefault();
    event.stopPropagation();
  }
  console.log('[PROFILE] Edit button clicked for recipe:', recipeId);
  router.push(`/recipes/${recipeId}/edit`);
};

// Delete Recipe
const deleteRecipe = async (recipeId) => {
  if (!confirm('Are you sure you want to delete this recipe? This action cannot be undone.')) {
    return;
  }
  
  if (!token.value) {
    router.push('/login');
    return;
  }
  
  try {
    console.log(`[DELETE] Deleting recipe ${recipeId}`);
    const response = await fetch(`http://localhost:8081/recipes/${recipeId}`, {
      method: 'DELETE',
      headers: { 
        'Authorization': `Bearer ${token.value}`,
        'Content-Type': 'application/json'
      }
    });
    
    console.log(`[DELETE] Response status: ${response.status}`);
    
    if (response.ok) {
      console.log(`[DELETE] Recipe ${recipeId} deleted successfully`);
      // Refetch recipes to update the list
      await fetchMyRecipes();
      alert('Recipe deleted successfully!');
    } else {
      const errorData = await response.json().catch(() => ({}));
      console.error('[DELETE] Error response:', errorData);
      alert(errorData.error || 'Failed to delete recipe');
    }
  } catch (err) {
    console.error('[DELETE] Exception:', err);
    alert('Failed to delete recipe. Please check if backend is running.');
  }
};

// Watch for tab changes to refresh data
watch(activeTab, (newTab) => {
  if (newTab === 'bookmarked') {
    fetchBookmarkedRecipes();
  } else if (newTab === 'purchased') {
    fetchPurchasedRecipes();
  }
});

// Load data on mount
onMounted(() => {
  // Redirect if not logged in
  if (!token.value) {
    router.push('/login');
    return;
  }
  
  fetchMyRecipes();
  fetchBookmarkedRecipes();
  fetchPurchasedRecipes();
});
</script>

