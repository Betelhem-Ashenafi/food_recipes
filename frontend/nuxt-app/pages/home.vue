<template>
  <div class="relative min-h-screen">
    <!-- Background Image -->
    <div class="fixed inset-0 z-0">
      <img 
        src="https://images.unsplash.com/photo-1544025162-d76694265947?ixlib=rb-4.0.3&auto=format&fit=crop&w=2070&q=80" 
        alt="Food background" 
        class="w-full h-full object-cover brightness-60"
      >
      <div class="absolute inset-0 bg-black/80"></div>
    </div>

    <!-- Content -->
    <div class="relative z-10">
      <!-- Hero Section -->
      <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 pt-16 pb-16">
        <div class="flex flex-col items-center justify-center mb-14">
          <div class="bg-white/10 backdrop-blur-lg border border-white/20 rounded-3xl shadow-2xl px-8 py-10 w-full text-center">
            <h1 class="text-5xl md:text-6xl lg:text-7xl font-serif font-extrabold text-white mb-4">
              <span class="text-emerald-400">Legendary</span> Recipes Await
            </h1>
            <p class="text-lg md:text-xl text-gray-200 max-w-2xl mx-auto font-light mb-2">
              Unleash your inner chef. Discover, create, and share premium recipes.
            </p>
            <p class="text-base text-emerald-300 font-medium">Premium. Inspiring. Unforgettable.</p>
          </div>
        </div>

        <!-- Search Bar -->
        <div class="max-w-2xl mx-auto mb-12">
          <div class="relative">
            <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
              <svg class="h-6 w-6 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Search recipes by title..."
              @keyup.enter="handleSearchSubmit"
              class="block w-full pl-12 pr-4 py-4 border border-white/20 rounded-xl bg-white/10 backdrop-blur-lg text-white placeholder-gray-400 focus:ring-2 focus:ring-emerald-500"
            />
          </div>
        </div>

        <!-- Categories -->
        <div class="mb-16">
          <div class="text-center mb-8">
            <h2 class="text-3xl font-bold text-white mb-2">Browse by Category</h2>
            <p class="text-gray-400">Explore recipes organized by cuisine and type</p>
          </div>
          
          <div v-if="categoriesPending" class="flex justify-center py-12">
            <div class="animate-spin rounded-full h-12 w-12 border-t-4 border-b-4 border-emerald-500"></div>
          </div>
          <div v-else-if="categoriesError" class="text-center py-8 text-red-300">
            Failed to load categories: {{ categoriesError.message }}
          </div>
          <div v-else-if="categories.length === 0" class="text-center py-8 text-gray-400">
            Categories are unavailable right now.
          </div>
          <div v-else class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6 gap-4">
            <button
              v-for="category in categories"
              :key="category.id"
              @click="selectedCategory = selectedCategory === category.id ? null : category.id"
              :class="[
                'group relative overflow-hidden rounded-xl border transition-all duration-300 hover:scale-105 hover:shadow-2xl',
                selectedCategory === category.id
                  ? 'ring-2 ring-emerald-400 shadow-lg shadow-emerald-500/50'
                  : 'hover:shadow-xl'
              ]"
            >
              <div class="relative h-32 md:h-40 overflow-hidden">
                <img
                  v-if="shouldShowCategoryImage(category)"
                  :src="normalizedCategoryImageUrl(category)"
                  :alt="category.name"
                  class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-110"
                  @error="markCategoryImageFailed(category.id)"
                />
                <div
                  v-else
                  class="w-full h-full bg-gradient-to-br from-emerald-500/40 via-teal-500/40 to-blue-500/40 flex items-center justify-center"
                >
                  <span class="text-white font-semibold text-sm text-center px-2">{{ category.name }}</span>
                </div>
                <div class="absolute inset-0 bg-gradient-to-t from-black/80 via-black/40 to-transparent"></div>
                <div class="absolute bottom-0 left-0 right-0 p-3 text-center">
                  <p class="text-white font-bold text-sm md:text-base drop-shadow-lg">{{ category.name }}</p>
                </div>
                <div v-if="selectedCategory === category.id" class="absolute top-2 right-2 bg-emerald-500 rounded-full p-1.5 shadow-lg">
                  <svg class="w-4 h-4 text-white" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                  </svg>
                </div>
              </div>
            </button>
          </div>
        </div>

        <!-- Filters -->
        <div class="max-w-4xl mx-auto mb-12">
          <div class="bg-white/10 backdrop-blur-lg border border-white/20 rounded-2xl p-6">
            <div class="grid grid-cols-1 md:grid-cols-5 gap-4">
              <div class="w-full">
                <label class="block text-sm font-medium text-gray-300 mb-2">Prep Time</label>
                <select v-model="prepTimeFilter" class="block w-full px-4 py-2 border border-white/20 rounded-lg bg-black/20 text-white focus:ring-2 focus:ring-emerald-500">
                  <option value="">Any Time</option>
                  <option value="15">≤15 min</option>
                  <option value="30">≤30 min</option>
                  <option value="60">≤1 hour</option>
                  <option value="120">≤2 hours</option>
                </select>
              </div>
              <div class="w-full">
                <label class="block text-sm font-medium text-gray-300 mb-2">Ingredient</label>
                <input v-model="ingredientFilter" type="text" placeholder="e.g., chicken" class="block w-full px-4 py-2 border border-white/20 rounded-lg bg-black/20 text-white placeholder-gray-400 focus:ring-2 focus:ring-emerald-500">
              </div>
              <div class="w-full">
                <label class="block text-sm font-medium text-gray-300 mb-2">Creator</label>
                <input v-model="creatorFilter" type="text" placeholder="Chef name" class="block w-full px-4 py-2 border border-white/20 rounded-lg bg-black/20 text-white placeholder-gray-400 focus:ring-2 focus:ring-emerald-500">
              </div>
              <div class="flex flex-col justify-end">
                <button @click="applyFilters" class="w-full px-4 py-2 border border-emerald-400/40 rounded-lg bg-emerald-500/20 text-white hover:bg-emerald-500/30 font-medium">
                  Apply Filter
                </button>
              </div>
              <div class="flex flex-col justify-end">
                <button @click="clearFilters" class="w-full px-4 py-2 border border-white/30 rounded-lg bg-white/5 text-white hover:bg-white/10 font-medium">
                  Clear
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Recipes -->
      <div ref="recipesSection" class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pb-16">
        <div class="text-center mb-8">
          <h2 class="text-3xl font-bold text-white mb-2">All Recipes</h2>
          <p class="text-gray-400">{{ recipes.length }} recipe{{ recipes.length !== 1 ? 's' : '' }} found</p>
        </div>

        <div v-if="pending" class="flex justify-center py-20">
          <div class="animate-spin rounded-full h-16 w-16 border-t-4 border-b-4 border-emerald-500"></div>
        </div>
        <div v-else-if="error" class="text-center py-20">
          <div class="inline-block p-4 rounded-full bg-red-500/20 border border-red-500/50 text-red-400 mb-4">⚠️</div>
          <h3 class="text-lg font-medium text-white">Something went wrong</h3>
          <p class="mt-2 text-gray-400">{{ error.message }}</p>
        </div>
        <div v-else-if="recipes.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <div
            v-for="recipe in recipes"
            :key="recipe.id"
            class="group bg-white/10 backdrop-blur-lg border border-white/20 rounded-2xl overflow-hidden hover:bg-white/15 hover:border-emerald-400/50 transition-all duration-300 hover:-translate-y-2 hover:shadow-2xl hover:shadow-emerald-500/20"
          >
            <div class="relative h-64 overflow-hidden">
                <img
                  :src="(recipe.recipe_images && recipe.recipe_images[0]?.url) || 'https://images.unsplash.com/placeholder.jpg'"
                  :alt="recipe.title"
                  class="w-full h-full object-cover transition-transform duration-700 group-hover:scale-110"
                  @error="event => event.target.src = 'https://images.unsplash.com/placeholder.jpg'"
                >
              <div class="absolute inset-0 bg-gradient-to-t from-black/80 via-transparent to-transparent"></div>
              <div class="absolute top-4 right-4">
                <span v-if="recipe.price > 0" class="inline-flex items-center px-3 py-1.5 rounded-full text-xs font-bold bg-emerald-500/90 text-white backdrop-blur-sm shadow-lg">💎 {{ recipe.price }} Credits</span>
                <span v-else class="inline-flex items-center px-3 py-1.5 rounded-full text-xs font-bold bg-green-500/90 text-white backdrop-blur-sm shadow-lg">Free</span>
              </div>
              <div class="absolute bottom-4 left-4 flex items-center text-white text-sm font-medium">
                <svg class="w-5 h-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                {{ recipe.preparation_time || 'N/A' }} min
              </div>
            </div>
            <div class="p-6">
              <div class="flex items-center mb-4">
                <div class="h-10 w-10 rounded-full bg-gradient-to-br from-emerald-400 to-teal-500 flex items-center justify-center text-white font-bold text-sm shadow-lg">
                  {{ getRecipeAuthorName(recipe).charAt(0) || 'C' }}
                </div>
                <div class="ml-3">
                  <p class="text-white font-medium text-sm">{{ getRecipeAuthorName(recipe) }}</p>
                  <p class="text-gray-400 text-xs">{{ formatDate(recipe.created_at) }}</p>
                </div>
              </div>
              <h3 class="text-xl font-bold text-white mb-2 group-hover:text-emerald-400 line-clamp-1">{{ recipe.title }}</h3>
              <p class="text-emerald-300 text-xs mb-2">{{ getRecipeCategoryName(recipe) }}</p>
              <p class="text-gray-300 text-sm line-clamp-2 mb-6">{{ recipe.description }}</p>
              <NuxtLink :to="`/recipes/${recipe.id}`" class="block w-full text-center py-3 rounded-xl bg-gradient-to-r from-emerald-600 to-teal-600 text-white font-semibold hover:from-emerald-500 hover:to-teal-500 transition-all duration-300 shadow-lg hover:shadow-emerald-500/30">
                View Recipe
              </NuxtLink>
            </div>
          </div>
        </div>
        <div v-else class="text-center py-20">
          <div class="inline-block p-4 rounded-full bg-gray-500/20 border border-gray-500/50 text-gray-400 mb-4">😕</div>
          <h3 class="text-xl font-medium text-white mb-2">No recipes found</h3>
          <p class="text-gray-400">Try adjusting your filters</p>
        </div>

        <!-- Pagination -->
        <div class="flex justify-center gap-4 mt-8">
          <button @click="goPreviousPage" :disabled="page === 1" class="px-4 py-2 bg-gray-200 rounded disabled:opacity-50">Previous</button>
          <span class="text-white">Page {{ page }}</span>
          <button @click="goNextPage" :disabled="recipes.length < limit" class="px-4 py-2 bg-gray-200 rounded disabled:opacity-50">Next</button>
        </div>
      </div>
    </div>

    <!-- Mobile Create Button -->
    <NuxtLink 
      v-if="isAuthenticated"
      to="/create"
      class="fixed bottom-6 right-6 z-50 sm:hidden bg-gradient-to-r from-emerald-600 to-teal-600 text-white rounded-full p-4 shadow-2xl hover:shadow-emerald-500/50 transition-transform hover:scale-110"
      title="Create Recipe"
    >
      <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
      </svg>
    </NuxtLink>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue';
import { useQuery } from '@vue/apollo-composable';
import gql from 'graphql-tag';

// Authentication
const token = useCookie('auth_token');
const isAuthenticated = computed(() => !!token.value);

// Filters
const searchQuery = ref('');
const selectedCategory = ref(null);
const prepTimeFilter = ref('');
const ingredientFilter = ref('');
const creatorFilter = ref('');
const page = ref(1);
const limit = ref(6);
const recipesSection = ref(null);
const failedCategoryImageIds = ref(new Set());

// GraphQL Queries
const categoriesQuery = gql`
  query GetCategories {
    categories(order_by: { name: asc }) {
      id
      name
      image_url
    }
  }
`;

const recipesQuery = gql`
  query GetRecipes(
    $limit: Int, 
    $offset: Int,
    $where: recipes_bool_exp
  ) {
    recipes(
      limit: $limit, 
      offset: $offset,
      where: $where,
      order_by: { created_at: desc }
    ) {
      id
      title
      description
      category_id
      price
      preparation_time
      created_at
      user_id
      user { id name }
      recipe_images(order_by: [{ is_featured: desc }, { id: asc }], limit: 1) { url }
    }
  }
`;

const usersByIdsQuery = gql`
  query GetUsersByIds($userIds: [Int!]) {
    users(where: { id: { _in: $userIds } }) {
      id
      name
    }
  }
`;

// Fetch categories
const { result: categoriesData, loading: categoriesPending, error: categoriesError } = useQuery(
  categoriesQuery,
  {},
  { errorPolicy: 'all', fetchPolicy: 'cache-and-network' }
);
const categories = computed(() => categoriesData.value?.categories || []);

// Build filter variables for recipes
const hasValue = (value) => value !== null && value !== undefined && value !== '';

const filterVariables = computed(() => {
  const title = searchQuery.value?.trim() || '';
  const creator = creatorFilter.value?.trim() || '';
  const ingredient = ingredientFilter.value?.trim() || '';

  const andFilters = [];

  if (title !== '') {
    andFilters.push({ title: { _ilike: `%${title}%` } });
  }

  if (hasValue(selectedCategory.value)) {
    andFilters.push({ category_id: { _eq: selectedCategory.value } });
  }

  if (hasValue(prepTimeFilter.value)) {
    andFilters.push({ preparation_time: { _lte: parseInt(prepTimeFilter.value) } });
  }

  if (creator !== '') {
    andFilters.push({ user: { name: { _ilike: `%${creator}%` } } });
  }

  if (ingredient !== '') {
    andFilters.push({ recipe_ingredients: { name: { _ilike: `%${ingredient}%` } } });
  }

  return {
    limit: limit.value,
    offset: (page.value - 1) * limit.value,
    where: andFilters.length > 0 ? { _and: andFilters } : {},
  };
});

// Fetch recipes
const { result: recipesData, error: recipesError, loading: recipesPending, refetch } = useQuery(
  recipesQuery,
  () => filterVariables.value,
  { errorPolicy: 'all', fetchPolicy: 'cache-and-network' }
);
const recipesRaw = computed(() => recipesData.value?.recipes || []);

const recipeUserIds = computed(() => {
  const uniqueIds = new Set(
    recipesRaw.value
      .map((recipe) => recipe.user_id)
      .filter((userId) => userId !== null && userId !== undefined)
  );
  return Array.from(uniqueIds);
});

const usersVariables = computed(() => ({
  userIds: recipeUserIds.value,
}));

const { result: usersData } = useQuery(
  usersByIdsQuery,
  usersVariables,
  {
    errorPolicy: 'all',
    fetchPolicy: 'cache-and-network',
    enabled: computed(() => recipeUserIds.value.length > 0),
  }
);

const usersById = computed(() => {
  return (usersData.value?.users || []).reduce((lookup, user) => {
    lookup[user.id] = user.name;
    return lookup;
  }, {});
});

const recipes = computed(() => {
  return recipesRaw.value.map((recipe) => {
    const fallbackName = usersById.value[recipe.user_id];
    if (!recipe.user && fallbackName) {
      return {
        ...recipe,
        user: {
          id: recipe.user_id,
          name: fallbackName,
        },
      };
    }
    return recipe;
  });
});

// Combined loading/error
const pending = computed(() => recipesPending.value || categoriesPending.value);
const error = computed(() => recipesError.value || categoriesError.value);

const refetchWithCurrentVariables = () => {
  return refetch({ ...filterVariables.value });
};

const goToPage = (targetPage) => {
  const nextPage = Math.max(1, targetPage);
  if (nextPage === page.value) return;
  page.value = nextPage;
  refetchWithCurrentVariables();
};

const goPreviousPage = () => {
  goToPage(page.value - 1);
};

const goNextPage = () => {
  goToPage(page.value + 1);
};

const refreshRecipes = (resetPage = false) => {
  if (resetPage && page.value !== 1) {
    page.value = 1;
  }
  refetchWithCurrentVariables();
};

// Reset to first page and refetch when filters change
watch([searchQuery, selectedCategory, ingredientFilter, creatorFilter, prepTimeFilter], () => {
  refreshRecipes(true);
});

// Helper functions
const formatDate = (dateString) => {
  if (!dateString) return '';
  // Only format date on client to avoid SSR/client mismatch
  if (typeof window === 'undefined') return dateString; // SSR: return raw string or placeholder
  const date = new Date(dateString);
  if (Number.isNaN(date.getTime())) return '';
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  });
};

const handleImageError = (e) => e.target.style.display = 'none';

const normalizeImageUrl = (url) => {
  const raw = String(url || '').trim();
  if (!raw) return '';
  if (/^https?:\/\//i.test(raw)) return raw;
  if (/^\/\//.test(raw)) return `https:${raw}`;
  return `https://${raw}`;
};

const normalizedCategoryImageUrl = (category) => {
  return normalizeImageUrl(category?.image_url);
};

const shouldShowCategoryImage = (category) => {
  const categoryId = Number(category?.id);
  const hasImage = normalizedCategoryImageUrl(category) !== '';
  return hasImage && !failedCategoryImageIds.value.has(categoryId);
};

const markCategoryImageFailed = (categoryId) => {
  const next = new Set(failedCategoryImageIds.value);
  next.add(Number(categoryId));
  failedCategoryImageIds.value = next;
};

const getNonEmptyText = (...values) => {
  for (const value of values) {
    if (typeof value === 'string' && value.trim() !== '') {
      return value;
    }
  }
  return '';
};

const getRecipeAuthorName = (recipe) => {
  return getNonEmptyText(usersById.value[recipe.user_id], recipe.user?.name) || 'Unknown User';
};

const getRecipeCategoryName = (recipe) => {
  const category = categories.value.find((item) => item.id === recipe.category_id);
  return category?.name || 'Unknown Category';
};

const clearFilters = () => {
  searchQuery.value = '';
  selectedCategory.value = null;
  prepTimeFilter.value = '';
  ingredientFilter.value = '';
  creatorFilter.value = '';
  refreshRecipes(true);
  scrollToRecipes();
};

const scrollToRecipes = () => {
  if (process.client && recipesSection.value) {
    recipesSection.value.scrollIntoView({ behavior: 'smooth', block: 'start' });
  }
};

const applyFilters = () => {
  refreshRecipes(true);
  scrollToRecipes();
};

const handleSearchSubmit = () => {
  applyFilters();
};
</script>

<style scoped>
.line-clamp-1 {
  display: -webkit-box;
  -webkit-line-clamp: 1;
  line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>