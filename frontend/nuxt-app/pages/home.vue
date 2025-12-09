<template>
  <div class="relative min-h-screen">
    <!-- Background Image with Overlay (Matching Login Page) -->
    <div 
      class="fixed inset-0 z-0 bg-cover bg-center brightness-95"
      style="background-image: url('https://images.unsplash.com/photo-1544025162-d76694265947?auto=format&fit=crop&w=2100&q=80');"
    >
      <div class="absolute inset-0 bg-gradient-to-b from-black/40 via-black/10 to-black/60"></div>
    </div>

    <!-- Content Container -->
    <div class="relative z-10">
      <!-- Hero Section with Search -->
      <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 pt-16 pb-16">
        <div class="flex flex-col items-center justify-center mb-14">
          <div class="bg-white/10 backdrop-blur-lg border border-white/20 rounded-3xl shadow-2xl px-8 py-10 w-full text-center animate-fade-in-up">
            <h1 class="text-5xl md:text-6xl lg:text-7xl font-serif font-extrabold text-white mb-4 tracking-tight drop-shadow-lg">
              <span class="text-emerald-400">Legendary</span> Recipes Await
            </h1>
            <p class="text-lg md:text-xl text-gray-200 max-w-2xl mx-auto font-light mb-2">
              Unleash your inner chef. Discover, create, and share premium recipes in a stunning, modern experience.
            </p>
            <p class="text-base text-emerald-300 font-medium mb-0">Premium. Inspiring. Unforgettable.</p>
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
              class="block w-full pl-12 pr-4 py-4 border border-white/20 rounded-xl bg-white/10 backdrop-blur-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 transition-all text-lg"
            />
          </div>
        </div>

        <!-- Categories Section -->
        <div class="mb-16">
          <div class="text-center mb-8">
            <h2 class="text-3xl font-bold text-white mb-2">Browse by Category</h2>
            <p class="text-gray-400">Explore recipes organized by cuisine and type</p>
          </div>
          
          <!-- Loading Categories -->
          <div v-if="categoriesPending" class="flex justify-center py-12">
            <div class="animate-spin rounded-full h-12 w-12 border-t-4 border-b-4 border-emerald-500"></div>
          </div>

          <!-- Categories Grid -->
          <div v-else class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6 gap-4">
            <button
              v-for="category in categories"
              :key="category.id"
              @click="selectedCategory = selectedCategory === category.id ? null : category.id"
              :class="[
                'group relative overflow-hidden rounded-xl p-6 backdrop-blur-lg border transition-all duration-300 transform hover:scale-105',
                selectedCategory === category.id
                  ? 'bg-emerald-500/30 border-emerald-400/50 shadow-lg shadow-emerald-500/20'
                  : 'bg-white/10 border-white/20 hover:bg-white/15 hover:border-white/30'
              ]"
            >
              <div class="text-center">
                <div class="text-4xl mb-2">{{ getCategoryEmoji(category.name) }}</div>
                <p class="text-white font-semibold text-sm mt-2">{{ category.name }}</p>
              </div>
              <div v-if="selectedCategory === category.id" class="absolute top-2 right-2">
                <svg class="w-5 h-5 text-emerald-400" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                </svg>
              </div>
            </button>
          </div>
        </div>

        <!-- Filters Section -->
        <div class="max-w-4xl mx-auto mb-12">
          <div class="bg-white/10 backdrop-blur-lg border border-white/20 rounded-2xl p-6">
            <div class="grid grid-cols-1 md:grid-cols-5 gap-4">
              <!-- Preparation Time Filter -->
              <div>
                <label class="block text-sm font-medium text-gray-300 mb-2">Preparation Time</label>
                <select
                  v-model="prepTimeFilter"
                  class="block w-full px-4 py-2 border border-white/20 rounded-lg bg-black/20 text-white focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500"
                >
                  <option value="">Any Time</option>
                  <option value="15">15 minutes or less</option>
                  <option value="30">30 minutes or less</option>
                  <option value="60">1 hour or less</option>
                  <option value="120">2 hours or less</option>
                </select>
              </div>

              <!-- Ingredient Filter -->
              <div>
                <label class="block text-sm font-medium text-gray-300 mb-2">Ingredient</label>
                <input
                  v-model="ingredientFilter"
                  type="text"
                  placeholder="e.g., chicken, tomato"
                  class="block w-full px-4 py-2 border border-white/20 rounded-lg bg-black/20 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500"
                />
              </div>

              <!-- Creator Filter -->
              <div>
                <label class="block text-sm font-medium text-gray-300 mb-2">Creator</label>
                <input
                  v-model="creatorFilter"
                  type="text"
                  placeholder="Chef name..."
                  class="block w-full px-4 py-2 border border-white/20 rounded-lg bg-black/20 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500"
                />
              </div>

              <!-- Sort By -->
              <div>
                <label class="block text-sm font-medium text-gray-300 mb-2">Sort By</label>
                <select
                  v-model="sortBy"
                  class="block w-full px-4 py-2 border border-white/20 rounded-lg bg-black/20 text-white focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500"
                >
                  <option value="newest">Newest First</option>
                  <option value="oldest">Oldest First</option>
                  <option value="title">Title A-Z</option>
                </select>
              </div>

              <!-- Clear Filters -->
              <div class="flex items-end">
                <button
                  @click="clearFilters"
                  class="w-full px-4 py-2 border border-white/30 rounded-lg bg-white/5 text-white hover:bg-white/10 transition-colors font-medium"
                >
                  Clear Filters
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Recipes Section -->
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pb-16">
        <div class="text-center mb-8">
          <h2 class="text-3xl font-bold text-white mb-2">
            {{ showingAll ? 'All Recipes' : 'Featured Recipes' }}
          </h2>
          <p class="text-gray-400">{{ filteredRecipes.length }} recipe{{ filteredRecipes.length !== 1 ? 's' : '' }} found</p>
        </div>

        <!-- Loading State -->
        <div v-if="pending" class="flex justify-center py-20">
          <div class="animate-spin rounded-full h-16 w-16 border-t-4 border-b-4 border-emerald-500"></div>
        </div>

        <!-- Error State -->
        <div v-else-if="error" class="text-center py-20">
          <div class="inline-block p-4 rounded-full bg-red-500/20 border border-red-500/50 text-red-400 mb-4">
            <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
            </svg>
          </div>
          <h3 class="text-lg font-medium text-white">Oops! Something went wrong.</h3>
          <p class="mt-2 text-gray-400">{{ error.message }}</p>
        </div>

        <!-- Recipe Grid -->
        <div v-else-if="displayedRecipes.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <div
            v-for="recipe in displayedRecipes"
            :key="recipe.id"
            class="group bg-white/10 backdrop-blur-lg border border-white/20 rounded-2xl overflow-hidden hover:bg-white/15 hover:border-emerald-400/50 transition-all duration-300 transform hover:-translate-y-2 hover:shadow-2xl hover:shadow-emerald-500/20 cursor-pointer"
            @click="navigateToRecipe(recipe.id)"
          >
            <!-- Image -->
            <div class="relative h-64 overflow-hidden">
              <img 
                :src="getRecipeImage(recipe)" 
                :alt="recipe.title"
                class="w-full h-full object-cover transition-transform duration-700 group-hover:scale-110"
              >
              <div class="absolute inset-0 bg-gradient-to-t from-black/80 via-transparent to-transparent"></div>
              
              <!-- Price Badge -->
              <div class="absolute top-4 right-4">
                <span v-if="recipe.price > 0" class="inline-flex items-center px-3 py-1.5 rounded-full text-xs font-bold bg-emerald-500/90 text-white backdrop-blur-sm shadow-lg">
                  💎 {{ recipe.price }} Credits
                </span>
                <span v-else class="inline-flex items-center px-3 py-1.5 rounded-full text-xs font-bold bg-green-500/90 text-white backdrop-blur-sm shadow-lg">
                  Free
                </span>
              </div>

              <!-- Preparation Time -->
              <div class="absolute bottom-4 left-4 flex items-center text-white text-sm font-medium">
                <svg class="w-5 h-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                {{ recipe.preparation_time || 'N/A' }} min
              </div>
            </div>

            <!-- Content -->
            <div class="p-6">
              <!-- User Info -->
              <div class="flex items-center mb-4">
                <div class="h-10 w-10 rounded-full bg-gradient-to-br from-emerald-400 to-teal-500 flex items-center justify-center text-white font-bold text-sm shadow-lg">
                  {{ recipe.user?.name?.charAt(0) || 'C' }}
                </div>
                <div class="ml-3">
                  <p class="text-white font-medium text-sm">{{ recipe.user?.name || 'Chef' }}</p>
                  <p class="text-gray-400 text-xs">{{ formatDate(recipe.created_at) }}</p>
                </div>
              </div>
              
              <h3 class="text-xl font-bold text-white mb-2 group-hover:text-emerald-400 transition-colors line-clamp-1">
                {{ recipe.title }}
              </h3>
              <p class="text-gray-300 text-sm line-clamp-2 mb-6 min-h-[2.5rem]">
                {{ recipe.description }}
              </p>

              <NuxtLink 
                :to="`/recipes/${recipe.id}`" 
                class="block w-full text-center py-3 rounded-xl bg-gradient-to-r from-emerald-600 to-teal-600 text-white font-semibold hover:from-emerald-500 hover:to-teal-500 transition-all duration-300 shadow-lg hover:shadow-emerald-500/30 transform hover:-translate-y-0.5"
              >
                View Recipe
              </NuxtLink>
            </div>
          </div>
        </div>

        <!-- No Results -->
        <div v-else class="text-center py-20">
          <div class="inline-block p-4 rounded-full bg-gray-500/20 border border-gray-500/50 text-gray-400 mb-4">
            <svg class="w-12 h-12" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <h3 class="text-xl font-medium text-white mb-2">No recipes found</h3>
          <p class="text-gray-400">Try adjusting your search or filters</p>
        </div>

        <!-- Show More Button -->
        <div v-if="!showingAll && filteredRecipes.length > featuredCount" class="text-center mt-12">
          <button
            @click="showingAll = true"
            class="px-8 py-4 bg-gradient-to-r from-emerald-600 to-teal-600 text-white font-bold rounded-xl hover:from-emerald-500 hover:to-teal-500 transition-all duration-300 shadow-lg hover:shadow-emerald-500/30 transform hover:-translate-y-1"
          >
            Show All Recipes ({{ filteredRecipes.length - featuredCount }} more)
          </button>
        </div>

        <div v-if="showingAll && filteredRecipes.length > featuredCount" class="text-center mt-8">
          <button
            @click="showingAll = false; scrollToTop()"
            class="px-6 py-3 bg-white/10 border border-white/20 text-white font-medium rounded-lg hover:bg-white/20 transition-all"
          >
            Show Less
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useQuery } from '@vue/apollo-composable';
import gql from 'graphql-tag';
import { defineComponent } from 'vue';

// GraphQL Query for Recipes (using Vue Apollo with Hasura)
const recipesQuery = gql`
  query GetRecipes {
    recipes(order_by: { created_at: desc }) {
      id
      title
      description
      thumbnail_url
      price
      preparation_time
      created_at
      category_id
      user_id
      user {
        id
        name
      }
      category {
        id
        name
      }
    }
  }
`;

// GraphQL Query for Categories (using Vue Apollo with Hasura)
const categoriesQuery = gql`
  query GetCategories {
    categories(order_by: { name: asc }) {
      id
      name
      image_url
    }
  }
`;

// Fetch Recipes using Vue Apollo
const { result: recipesData, error: recipesError, loading: recipesPending } = useQuery(recipesQuery);

// Fetch Categories using Vue Apollo
const { result: categoriesData, loading: categoriesPending } = useQuery(categoriesQuery);

// Computed properties
const recipes = computed(() => recipesData.value?.recipes || []);
const categories = computed(() => categoriesData.value?.categories || []);
const pending = computed(() => recipesPending.value);
const error = computed(() => recipesError.value);

// Filters and Display State
const searchQuery = ref('');
const selectedCategory = ref(null);
const prepTimeFilter = ref('');
const ingredientFilter = ref('');
const creatorFilter = ref('');
const sortBy = ref('newest');
const showingAll = ref(false);
const featuredCount = 6;
const showFeatured = ref(true);
const showTrending = ref(true);

// Get Category Emoji
const getCategoryEmoji = (name) => {
  const emojiMap = {
    'Italian': '🍝',
    'Mexican': '🌮',
    'Asian': '🍜',
    'Dessert': '🍰',
    'Breakfast': '🥞',
    'Lunch': '🥗',
    'Dinner': '🍽️',
    'Vegetarian': '🥬',
    'Vegan': '🌱',
    'Seafood': '🐟',
    'Pasta': '🍝',
    'Pizza': '🍕',
    'Salad': '🥗',
    'Soup': '🍲',
    'Beverage': '🥤',
  };
  return emojiMap[name] || '🍳';
};

// Format Date
const formatDate = (dateString) => {
  if (!dateString) return '';
  const date = new Date(dateString);
  const now = new Date();
  const diffTime = Math.abs(now - date);
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
  
  if (diffDays === 0) return 'Today';
  if (diffDays === 1) return 'Yesterday';
  if (diffDays < 7) return `${diffDays} days ago`;
  if (diffDays < 30) return `${Math.floor(diffDays / 7)} weeks ago`;
  return date.toLocaleDateString();
};

// Filtered Recipes
const filteredRecipes = computed(() => {
  if (!recipes.value || recipes.value.length === 0) return [];

  let filtered = [...recipes.value];

  // Search by title
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase();
    filtered = filtered.filter(recipe => 
      recipe.title.toLowerCase().includes(query)
    );
  }

  // Filter by category
  if (selectedCategory.value) {
    filtered = filtered.filter(recipe => 
      recipe.category_id === selectedCategory.value
    );
  }

  // Filter by preparation time
  if (prepTimeFilter.value) {
    const maxTime = parseInt(prepTimeFilter.value);
    filtered = filtered.filter(recipe => 
      recipe.preparation_time && recipe.preparation_time <= maxTime
    );
  }

  // Filter by ingredient (search in title/description)
  if (ingredientFilter.value) {
    const ingredient = ingredientFilter.value.toLowerCase();
    filtered = filtered.filter(recipe => 
      recipe.title.toLowerCase().includes(ingredient) ||
      recipe.description?.toLowerCase().includes(ingredient)
    );
  }

  // Filter by creator
  if (creatorFilter.value) {
    filtered = filtered.filter(recipe => 
      recipe.user?.name?.toLowerCase().includes(creatorFilter.value.toLowerCase())
    );
  }

  // Sort
  if (sortBy.value === 'newest') {
    filtered.sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
  } else if (sortBy.value === 'oldest') {
    filtered.sort((a, b) => new Date(a.created_at) - new Date(b.created_at));
  } else if (sortBy.value === 'title') {
    filtered.sort((a, b) => a.title.localeCompare(b.title));
  }

  return filtered;
});

// Displayed Recipes (featured or all)
const displayedRecipes = computed(() => {
  return showingAll.value ? filteredRecipes.value : filteredRecipes.value.slice(0, featuredCount);
});

// Navigate to Recipe
const navigateToRecipe = (id) => {
  navigateTo(`/recipes/${id}`);
};

// Scroll to Top
const scrollToTop = () => {
  window.scrollTo({ top: 0, behavior: 'smooth' });
};

// Trending Recipes (based on likes - mock for now, would need likes count from backend)
const trendingRecipes = computed(() => {
  return recipes.value.slice(0, 6);
});

// Featured Recipes
const featuredRecipes = computed(() => {
  return recipes.value.slice(0, 6);
});

// Get real image based on recipe title
const getRecipeImage = (recipe) => {
  if (recipe.thumbnail_url) return recipe.thumbnail_url;
  
  const title = recipe.title.toLowerCase();
  
  // Map recipe keywords to real images
  if (title.includes('avocado') || title.includes('salad')) {
    return 'https://images.unsplash.com/photo-1546069901-ba9599a7e63c?auto=format&fit=crop&w=800&q=80';
  }
  if (title.includes('chocolate') || title.includes('cake')) {
    return 'https://images.unsplash.com/photo-1578985545062-69928b1d9587?auto=format&fit=crop&w=800&q=80';
  }
  if (title.includes('pasta')) {
    return 'https://images.unsplash.com/photo-1621996346565-e3dbc646d9a9?auto=format&fit=crop&w=800&q=80';
  }
  if (title.includes('pizza')) {
    return 'https://images.unsplash.com/photo-1513104890138-7c749659a591?auto=format&fit=crop&w=800&q=80';
  }
  if (title.includes('burger')) {
    return 'https://images.unsplash.com/photo-1568901346375-23c9450c58cd?auto=format&fit=crop&w=800&q=80';
  }
  if (title.includes('sushi')) {
    return 'https://images.unsplash.com/photo-1579584425555-c3ce17fd4351?auto=format&fit=crop&w=800&q=80';
  }
  if (title.includes('soup')) {
    return 'https://images.unsplash.com/photo-1547592166-23ac45744acd?auto=format&fit=crop&w=800&q=80';
  }
  
  // Default food image
  return 'https://images.unsplash.com/photo-1495521821757-a1efb6729352?auto=format&fit=crop&w=800&q=80';
};

// Clear Filters
const clearFilters = () => {
  searchQuery.value = '';
  selectedCategory.value = null;
  prepTimeFilter.value = '';
  ingredientFilter.value = '';
  creatorFilter.value = '';
  sortBy.value = 'newest';
  showingAll.value = false;
};
</script>

<style scoped>
.line-clamp-1 {
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>