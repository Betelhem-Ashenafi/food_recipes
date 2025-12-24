<template>
  <div class="relative min-h-screen pb-20">
    <!-- Background -->
    <div class="fixed inset-0 z-0">
      <img 
        src="https://images.unsplash.com/photo-1504674900247-0877df9cc836?ixlib=rb-4.0.3&auto=format&fit=crop&w=2070&q=80" 
        alt="Food" 
        class="w-full h-full object-cover brightness-60"
      >
      <div class="absolute inset-0 bg-black/80"></div>
    </div>

    <!-- Loading State -->
    <div v-if="pending" class="relative z-10 flex justify-center items-center min-h-screen">
      <div class="text-center">
        <div class="animate-spin rounded-full h-16 w-16 border-t-4 border-b-4 border-emerald-500 mx-auto"></div>
        <p class="mt-4 text-white">Loading recipe...</p>
      </div>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="relative z-10 max-w-2xl mx-auto px-4 py-20">
      <div class="bg-red-500/20 border border-red-500/50 text-red-200 px-6 py-4 rounded-lg">
        <strong class="font-bold">Error!</strong>
        <span class="block sm:inline"> {{ error?.message || error?.toString() || 'Failed to load recipe' }}</span>
      </div>
    </div>

    <!-- Recipe Content -->
    <div v-else-if="recipe" class="relative z-10 max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 pt-8">
      <!-- Image Gallery with Navigation -->
      <div class="mb-8 rounded-2xl overflow-hidden shadow-2xl relative">
        <div v-if="recipeImages.length > 0" class="relative">
          <!-- Main Image Display -->
          <img 
            :src="recipeImages[currentImageIndex].url" 
            :alt="recipe.title"
            class="w-full h-96 object-cover transition-opacity duration-300"
          >
          
          <!-- Navigation Arrows (only if more than 1 image) -->
          <div v-if="recipeImages.length > 1" class="absolute inset-0 flex items-center justify-between p-4">
            <button
              @click="previousImage"
              class="bg-black/50 hover:bg-black/70 text-white p-3 rounded-full transition-all"
              :disabled="currentImageIndex === 0"
            >
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
              </svg>
            </button>
            <button
              @click="nextImage"
              class="bg-black/50 hover:bg-black/70 text-white p-3 rounded-full transition-all"
              :disabled="currentImageIndex === recipeImages.length - 1"
            >
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
              </svg>
            </button>
          </div>
          
          <!-- Image Indicators -->
          <div v-if="recipeImages.length > 1" class="absolute bottom-4 left-1/2 transform -translate-x-1/2 flex gap-2">
            <div
              v-for="(img, index) in recipeImages"
              :key="index"
              class="h-2 rounded-full transition-all"
              :class="index === currentImageIndex ? 'w-8 bg-white' : 'w-2 bg-white/50'"
            ></div>
          </div>
          
          <!-- Featured Badge -->
        </div>
        
        <!-- Fallback if no images -->
        <img 
          v-else
          :src="recipe.thumbnail_url || 'https://images.unsplash.com/photo-1495521821757-a1efb6729352?ixlib=rb-4.0.3&auto=format&fit=crop&w=1200&q=80'" 
          :alt="recipe.title"
          class="w-full h-96 object-cover"
        >
      </div>

      <!-- Main Content -->
      <div class="bg-white/10 backdrop-blur-lg border border-white/20 rounded-2xl p-8 mb-8">
        <!-- Title & Actions -->
        <div class="flex justify-between items-start mb-6">
          <div class="flex-1">
            <h1 class="text-4xl md:text-5xl font-serif font-bold text-white mb-4">
              {{ recipe.title }}
            </h1>
            
            <!-- Meta Info -->
            <div class="flex flex-wrap items-center gap-4 text-gray-300">
              <div class="flex items-center">
                <div class="h-10 w-10 rounded-full bg-gradient-to-br from-emerald-400 to-teal-500 flex items-center justify-center text-white font-bold text-sm mr-2">
                  {{ recipe.user?.name?.charAt(0) || 'C' }}
                </div>
                <span class="font-medium">{{ recipe.user?.name || 'Chef' }}</span>
              </div>
              <span>•</span>
              <div class="flex items-center">
                <svg class="w-5 h-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                {{ recipe.preparation_time }} min
              </div>
              <span>•</span>
              <div class="flex items-center">
                <svg class="w-5 h-5 mr-1 text-yellow-400" fill="currentColor" viewBox="0 0 20 20">
                  <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                </svg>
                <span v-if="ratingData">{{ ratingData.average_rating.toFixed(1) }} ({{ ratingData.count }} ratings)</span>
                <span v-else>No ratings yet</span>
              </div>
            </div>
          </div>

          <!-- Price Badge & Actions -->
          <div class="ml-4 flex flex-col items-end gap-3">
            <div v-if="recipe.price > 0" class="bg-emerald-500 text-white px-6 py-3 rounded-full text-xl font-bold shadow-lg">
              💎 {{ recipe.price }} Credits
            </div>
            <div v-else class="bg-green-500 text-white px-6 py-3 rounded-full text-xl font-bold shadow-lg">
              Free
            </div>
            
            <!-- Edit Button (only for owner) -->
            <button
              v-if="isAuthenticated && isOwner"
              @click="router.push(`/recipes/${recipeId}/edit`)"
              class="px-4 py-2 bg-emerald-500/20 hover:bg-emerald-500/30 border border-emerald-400/50 text-emerald-300 rounded-lg transition-colors flex items-center gap-2"
              title="Edit Recipe"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
              </svg>
              Edit Recipe
            </button>
          </div>
        </div>

        <!-- Description -->
        <p class="text-gray-200 text-lg leading-relaxed mb-8">
          {{ recipe.description }}
        </p>

        <!-- Social Actions -->
        <div class="flex gap-4 mb-8 pb-8 border-b border-white/20">
          <button 
            @click="toggleLike"
            :disabled="actionLoading"
            class="flex-1 flex items-center justify-center gap-2 px-4 py-3 rounded-lg transition-colors font-medium"
            :class="isLiked ? 'bg-red-500 text-white hover:bg-red-600' : 'bg-white/10 text-gray-300 hover:bg-white/20 border border-white/20'"
          >
            <svg class="w-5 h-5" :fill="isLiked ? 'currentColor' : 'none'" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
            </svg>
            {{ isLiked ? 'Liked' : 'Like' }}
          </button>

          <button 
            @click="toggleBookmark"
            :disabled="actionLoading"
            class="flex-1 flex items-center justify-center gap-2 px-4 py-3 rounded-lg transition-colors font-medium"
            :class="isBookmarked ? 'bg-yellow-500 text-white hover:bg-yellow-600' : 'bg-white/10 text-gray-300 hover:bg-white/20 border border-white/20'"
          >
            <svg class="w-5 h-5" :fill="isBookmarked ? 'currentColor' : 'none'" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z" />
            </svg>
            {{ isBookmarked ? 'Bookmarked' : 'Bookmark' }}
          </button>

          <button 
            v-if="recipe.price > 0 && !hasPurchased"
            @click="handleBuyRecipe"
            :disabled="buying"
            class="flex-1 flex items-center justify-center gap-2 px-4 py-3 bg-gradient-to-r from-emerald-600 to-teal-600 text-white rounded-lg hover:from-emerald-500 hover:to-teal-500 transition-all font-bold shadow-lg disabled:opacity-50"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z" />
            </svg>
            {{ buying ? 'Processing...' : 'Buy Recipe' }}
          </button>
        </div>

        <!-- Category -->
        <div class="mb-6">
          <div class="inline-flex items-center gap-3 px-4 py-2 rounded-full text-sm font-medium bg-emerald-500/30 text-emerald-300 border border-emerald-400/50 backdrop-blur-sm">
            <!-- Category Image -->
            <div v-if="recipe.category?.image_url" class="w-8 h-8 rounded-full overflow-hidden border-2 border-emerald-400/50 flex-shrink-0">
              <img 
                :src="recipe.category.image_url" 
                :alt="recipe.category.name"
                class="w-full h-full object-cover"
                @error="(e) => e.target.style.display = 'none'"
              />
            </div>
            <!-- Category Name -->
            <span>{{ recipe.category?.name || 'Uncategorized' }}</span>
          </div>
        </div>
      </div>

      <!-- Ingredients - Only show if free, purchased, or owner -->
      <div v-if="recipe.price === 0 || hasPurchased || isOwner" class="bg-white/10 backdrop-blur-lg border border-white/20 rounded-2xl p-8 mb-8">
        <h2 class="text-3xl font-bold text-white mb-6 flex items-center">
          <span class="text-emerald-400 mr-2">🥕</span> Ingredients
        </h2>
        
        <div v-if="ingredientsData && ingredientsData.length > 0" class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div 
            v-for="ingredient in ingredientsData" 
            :key="ingredient.id"
            class="flex items-center p-4 bg-white/5 rounded-lg border border-white/10"
          >
            <div class="h-2 w-2 bg-emerald-400 rounded-full mr-3"></div>
            <span class="text-gray-200">{{ ingredient.quantity }} {{ ingredient.unit }} {{ ingredient.name }}</span>
          </div>
        </div>
        <p v-else class="text-gray-400">No ingredients listed</p>
      </div>

      <!-- Payment Required Message for Ingredients -->
      <div v-else class="bg-white/10 backdrop-blur-lg border border-white/20 rounded-2xl p-8 mb-8">
        <div class="text-center py-8">
          <div class="w-16 h-16 bg-emerald-500/20 rounded-full flex items-center justify-center mx-auto mb-4">
            <svg class="w-8 h-8 text-emerald-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
            </svg>
          </div>
          <h3 class="text-2xl font-bold text-white mb-2">Ingredients Locked</h3>
          <p class="text-gray-300 mb-6">Purchase this recipe to view the full ingredients list and preparation steps.</p>
          <button
            @click="handleBuyRecipe"
            :disabled="buying"
            class="px-6 py-3 bg-gradient-to-r from-emerald-600 to-teal-600 text-white rounded-lg hover:from-emerald-500 hover:to-teal-500 transition-all font-bold shadow-lg disabled:opacity-50"
          >
            {{ buying ? 'Processing...' : `Buy Recipe for ${recipe.price} Credits` }}
          </button>
        </div>
      </div>

      <!-- Steps - Only show if free, purchased, or owner -->
      <div v-if="recipe.price === 0 || hasPurchased || isOwner" class="bg-white/10 backdrop-blur-lg border border-white/20 rounded-2xl p-8 mb-8">
        <h2 class="text-3xl font-bold text-white mb-6 flex items-center">
          <span class="text-emerald-400 mr-2">📋</span> Preparation Steps
        </h2>
        
        <div v-if="stepsData && stepsData.length > 0" class="space-y-6">
          <div 
            v-for="(step, index) in stepsData" 
            :key="step.id"
            class="flex gap-4"
          >
            <div class="flex-shrink-0 w-10 h-10 bg-emerald-500 text-white rounded-full flex items-center justify-center font-bold text-lg shadow-lg">
              {{ index + 1 }}
            </div>
            <div class="flex-1 bg-white/5 rounded-lg p-4 border border-white/10">
              <p class="text-gray-200 leading-relaxed">{{ step.instruction }}</p>
            </div>
          </div>
        </div>
        <p v-else class="text-gray-400">No steps listed</p>
      </div>

      <!-- Rating Section -->
      <div class="bg-white/10 backdrop-blur-lg border border-white/20 rounded-2xl p-8 mb-8">
        <h2 class="text-2xl font-bold text-white mb-4">Rate this Recipe</h2>
        <div v-if="isAuthenticated" class="flex items-center gap-2">
          <button 
            v-for="star in 5" 
            :key="star"
            @click="submitRating(star)"
            :disabled="ratingLoading"
            class="transition-transform hover:scale-110"
            :title="userRating > 0 ? `Your rating: ${userRating} stars` : `Click to rate ${star} stars`"
          >
            <svg 
              class="w-10 h-10 transition-colors" 
              :fill="star <= currentUserRating ? '#fbbf24' : 'none'" 
              :stroke="star <= currentUserRating ? '#fbbf24' : '#9ca3af'" 
              viewBox="0 0 24 24"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
            </svg>
          </button>
          <span v-if="currentUserRating > 0" class="text-yellow-400 ml-3 text-sm font-semibold">Your rating: {{ currentUserRating }} ⭐</span>
          <span v-else class="text-gray-400 ml-3 text-sm">Click a star to rate</span>
        </div>
        <div v-else class="p-4 bg-emerald-500/20 border border-emerald-400/50 rounded-lg text-center">
          <p class="text-white">Please <NuxtLink to="/login" class="text-emerald-400 hover:text-emerald-300 font-semibold underline">log in</NuxtLink> to rate this recipe.</p>
        </div>
        <p v-if="ratingSuccess && isAuthenticated" class="text-emerald-400 text-sm mt-2">Rating submitted!</p>
      </div>

      <!-- Comments Section -->
      <div class="bg-white/10 backdrop-blur-lg border border-white/20 rounded-2xl p-8">
        <h2 class="text-3xl font-bold text-white mb-6 flex items-center">
          <span class="text-emerald-400 mr-2">💬</span> Comments ({{ commentsData?.length || 0 }})
        </h2>

        <!-- Add Comment Form -->
        <div v-if="isAuthenticated" class="mb-8">
          <textarea 
            v-model="newComment"
            rows="3"
            placeholder="Share your thoughts about this recipe..."
            class="w-full px-4 py-3 border border-white/20 rounded-lg bg-black/20 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-emerald-500"
          ></textarea>
          <button 
            @click="submitComment"
            :disabled="!newComment.trim() || commentLoading"
            class="mt-3 px-6 py-3 bg-gradient-to-r from-emerald-600 to-teal-600 text-white font-semibold rounded-lg hover:from-emerald-500 hover:to-teal-500 transition-all disabled:opacity-50"
          >
            {{ commentLoading ? 'Posting...' : 'Post Comment' }}
          </button>
        </div>
        <div v-else class="mb-8 p-4 bg-emerald-500/20 border border-emerald-400/50 rounded-lg text-center">
          <p class="text-white mb-3">Please <NuxtLink to="/login" class="text-emerald-400 hover:text-emerald-300 font-semibold underline">log in</NuxtLink> to comment on this recipe.</p>
        </div>

        <!-- Comments List -->
        <div v-if="commentsData && commentsData.length > 0" class="space-y-4">
          <div 
            v-for="comment in commentsData" 
            :key="comment.id"
            class="bg-white/5 border border-white/10 rounded-lg p-4"
          >
            <div class="flex items-center mb-2">
              <div class="h-8 w-8 rounded-full bg-gradient-to-br from-emerald-400 to-teal-500 flex items-center justify-center text-white font-bold text-xs mr-2">
                {{ comment.user_name?.charAt(0) || 'U' }}
              </div>
              <span class="text-white font-medium">{{ comment.user_name || 'User' }}</span>
              <span class="text-gray-400 text-sm ml-auto">{{ formatDate(comment.created_at) }}</span>
            </div>
            <p class="text-gray-200">{{ comment.content }}</p>
          </div>
        </div>
        <p v-else class="text-gray-400 text-center py-8">No comments yet. Be the first to comment!</p>
      </div>
    </div>

    <!-- Not Found -->
    <div v-else class="relative z-10 text-center py-20">
      <p class="text-white text-xl">Recipe not found.</p>
      <NuxtLink to="/home" class="text-emerald-400 hover:text-emerald-300 mt-4 inline-block">
        ← Back to Home
      </NuxtLink>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { useQuery } from '@vue/apollo-composable';
import gql from 'graphql-tag';
import { jwtDecode } from 'jwt-decode';

const route = useRoute();
const router = useRouter();
const recipeId = parseInt(route.params.id);
// Check if coming from payment success page
const fromPayment = computed(() => route.query.fromPayment === 'true' || route.query.payment === 'success');
const token = useCookie('auth_token');
const isAuthenticated = computed(() => !!token.value);

// Get user info from JWT - always return safe object structure
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
  const claims = userInfo.value?.['https://hasura.io/jwt/claims'];
  if (!claims || typeof claims !== 'object' || claims === null) return 'User';
  return claims['x-hasura-user-name'] || 'User';
});
const userEmail = computed(() => {
  const claims = userInfo.value?.['https://hasura.io/jwt/claims'];
  if (!claims || typeof claims !== 'object' || claims === null) return '';
  return claims['x-hasura-user-email'] || '';
});
const userId = computed(() => {
  const claims = userInfo.value?.['https://hasura.io/jwt/claims'];
  if (!claims || typeof claims !== 'object' || claims === null) return null;
  const id = claims['x-hasura-user-id'];
  return id ? parseInt(id) : null;
});
const isOwner = computed(() => {
  if (!isAuthenticated.value || !userId.value || !recipe.value) return false;
  // Ensure both are numbers for comparison
  const recipeUserId = typeof recipe.value.user_id === 'string' ? parseInt(recipe.value.user_id) : recipe.value.user_id;
  const currentUserId = typeof userId.value === 'string' ? parseInt(userId.value) : userId.value;
  return recipeUserId === currentUserId;
});

// GraphQL Query for Recipe (using Vue Apollo with Hasura)
const query = gql`
  query GetRecipe($id: Int!) {
    recipes_by_pk(id: $id) {
      id
      title
      description
      price
      thumbnail_url
      created_at
      preparation_time
      user_id
      category_id
      user {
        id
        name
      }
      category {
        id
        name
        image_url
      }
    }
  }
`;

// Use useQuery with proper error handling
// Skip query if recipeId is invalid
const { result, loading: pending, error } = useQuery(
  query, 
  () => ({ id: recipeId }),
  { 
    skip: () => !recipeId || recipeId === 0,
    errorPolicy: 'all',
    fetchPolicy: 'cache-and-network',
    returnPartialData: true
  }
);

const recipe = computed(() => {
  if (!result.value || typeof result.value !== 'object') return null;
  if (!result.value.recipes_by_pk) return null;
  return result.value.recipes_by_pk;
});

// Fetch Recipe Images (REST API)
const recipeImages = ref([]);
const currentImageIndex = ref(0);

// Get API URL from runtime config
const config = useRuntimeConfig();
const getApiUrl = () => config.public?.apiUrl || 'http://localhost:8081';

const fetchRecipeImages = async () => {
  try {
    // Fetch images from recipe_images table
    const response = await fetch(`${getApiUrl()}/recipes/${recipeId}/images`);
    if (response.ok) {
      const images = await response.json();
      if (images && images.length > 0) {
        recipeImages.value = images;
        // Set current index to featured image if available
        const featuredIndex = recipeImages.value.findIndex(img => img.is_featured);
        if (featuredIndex >= 0) {
          currentImageIndex.value = featuredIndex;
        } else {
          currentImageIndex.value = 0;
        }
      } else {
        // No images in recipe_images table, use thumbnail as fallback
        if (recipe.value?.thumbnail_url) {
          recipeImages.value = [{ url: recipe.value.thumbnail_url, is_featured: true, id: 0, recipe_id: recipeId }];
          currentImageIndex.value = 0;
        }
      }
    } else {
      // If endpoint fails, use thumbnail as fallback
      if (recipe.value?.thumbnail_url) {
        recipeImages.value = [{ url: recipe.value.thumbnail_url, is_featured: true, id: 0, recipe_id: recipeId }];
        currentImageIndex.value = 0;
      }
    }
  } catch (err) {
    console.error('Error fetching images:', err);
    // Fallback to thumbnail
    if (recipe.value?.thumbnail_url) {
      recipeImages.value = [{ url: recipe.value.thumbnail_url, is_featured: true, id: 0, recipe_id: recipeId }];
      currentImageIndex.value = 0;
    }
  }
};

const nextImage = () => {
  if (currentImageIndex.value < recipeImages.value.length - 1) {
    currentImageIndex.value++;
  }
};

const previousImage = () => {
  if (currentImageIndex.value > 0) {
    currentImageIndex.value--;
  }
};

// Fetch Ingredients (REST API)
const ingredientsData = ref([]);
const fetchIngredients = async () => {
  try {
    const response = await fetch(`${getApiUrl()}/recipes/${recipeId}/ingredients`);
    if (response.ok) {
      ingredientsData.value = await response.json();
    }
  } catch (err) {
    console.error('Error fetching ingredients:', err);
  }
};

// Fetch Steps (REST API)
const stepsData = ref([]);
const fetchSteps = async () => {
  try {
    const response = await fetch(`${getApiUrl()}/recipes/${recipeId}/steps`);
    if (response.ok) {
      const data = await response.json();
      stepsData.value = data.sort((a, b) => a.step_number - b.step_number);
    }
  } catch (err) {
    console.error('Error fetching steps:', err);
  }
};

// Fetch Comments (REST API)
const commentsData = ref([]);
const fetchComments = async () => {
  if (!recipeId) {
    console.warn('[COMMENTS] No recipe ID, skipping fetch');
    return;
  }
  try {
    console.log(`[COMMENTS] Fetching comments for recipe ${recipeId}`);
    const response = await fetch(`${getApiUrl()}/recipes/${recipeId}/comments`);
    if (response.ok) {
      const data = await response.json();
      commentsData.value = data || [];
      console.log(`[COMMENTS] Loaded ${commentsData.value.length} comments`);
    } else {
      console.error(`[COMMENTS] Failed to fetch: ${response.status} ${response.statusText}`);
      commentsData.value = [];
    }
  } catch (err) {
    console.error('[COMMENTS] Error fetching comments:', err);
    commentsData.value = [];
  }
};

// Fetch Rating (REST API)
const ratingData = ref(null);
const fetchRating = async () => {
  if (!recipeId) {
    console.warn('[RATING] No recipe ID, skipping fetch');
    return;
  }
  try {
    console.log(`[RATING] Fetching rating for recipe ${recipeId}`);
    const response = await fetch(`${getApiUrl()}/recipes/${recipeId}/rate`);
    if (response.ok) {
      ratingData.value = await response.json();
      console.log(`[RATING] Loaded rating:`, ratingData.value);
    } else {
      console.error(`[RATING] Failed to fetch: ${response.status} ${response.statusText}`);
      ratingData.value = null;
    }
  } catch (err) {
    console.error('[RATING] Error fetching rating:', err);
    ratingData.value = null;
  }
};

// Social Features State
const isLiked = ref(false);
const isBookmarked = ref(false);
const actionLoading = ref(false);
const hasPurchased = ref(false);

// Toggle Like
const toggleLike = async () => {
  if (!token.value) {
    router.push('/login');
    return;
  }
  actionLoading.value = true;
  
  try {
    const method = isLiked.value ? 'DELETE' : 'POST';
    console.log(`[LIKE] ${method} /recipes/${recipeId}/like`);
      const response = await fetch(`${getApiUrl()}/recipes/${recipeId}/like`, {
      method,
      headers: { 
        'Authorization': `Bearer ${token.value}`,
        'Content-Type': 'application/json'
      }
    });
    
    console.log(`[LIKE] Response status: ${response.status}`);
    
    if (response.ok) {
      const data = await response.json().catch(() => ({}));
      console.log('[LIKE] Success:', data);
      isLiked.value = !isLiked.value;
      // Refresh like count if needed
      await fetchRating();
    } else {
      const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
      console.error('[LIKE] Error response:', errorData);
      alert(`Failed to ${isLiked.value ? 'unlike' : 'like'} recipe: ${errorData.error || 'Please try again'}`);
    }
  } catch (err) {
    console.error('[LIKE] Exception:', err);
    alert(`Error: ${err.message || 'Failed to like recipe. Please check if backend is running.'}`);
  } finally {
    actionLoading.value = false;
  }
};

// Toggle Bookmark
const toggleBookmark = async () => {
  if (!token.value) {
    router.push('/login');
    return;
  }
  actionLoading.value = true;
  
  try {
    const method = isBookmarked.value ? 'DELETE' : 'POST';
    console.log(`[BOOKMARK] ${method} /recipes/${recipeId}/bookmark`);
      const response = await fetch(`${getApiUrl()}/recipes/${recipeId}/bookmark`, {
      method,
      headers: { 
        'Authorization': `Bearer ${token.value}`,
        'Content-Type': 'application/json'
      }
    });
    
    console.log(`[BOOKMARK] Response status: ${response.status}`);
    
    if (response.ok) {
      const data = await response.json().catch(() => ({}));
      console.log('[BOOKMARK] Success:', data);
      isBookmarked.value = !isBookmarked.value;
    } else {
      const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
      console.error('[BOOKMARK] Error response:', errorData);
      alert(`Failed to ${isBookmarked.value ? 'unbookmark' : 'bookmark'} recipe: ${errorData.error || 'Please try again'}`);
    }
  } catch (err) {
    console.error('[BOOKMARK] Exception:', err);
    alert(`Error: ${err.message || 'Failed to bookmark recipe. Please check if backend is running.'}`);
  } finally {
    actionLoading.value = false;
  }
};

// Rating
const userRating = ref(0);
const ratingLoading = ref(false);
const ratingSuccess = ref(false);

// Computed to ensure reactivity for star display
const currentUserRating = computed(() => {
  const rating = Number(userRating.value) || 0;
  console.log('[RATING] Current user rating computed:', rating);
  return rating;
});

const submitRating = async (rating) => {
  if (!token.value) {
    router.push('/login');
    return;
  }
  
  if (rating < 1 || rating > 5) {
    alert('Rating must be between 1 and 5');
    return;
  }
  
  ratingLoading.value = true;
  ratingSuccess.value = false;
  
  try {
    console.log(`[RATING] Submitting rating ${rating} for recipe ${recipeId}`);
    const response = await fetch(`${getApiUrl()}/recipes/${recipeId}/rate`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token.value}`
      },
      body: JSON.stringify({ rating })
    });
    
    console.log(`[RATING] Response status: ${response.status}`);
    
    if (response.ok) {
      const data = await response.json().catch(() => ({}));
      console.log(`[RATING] Success:`, data);
      userRating.value = parseInt(rating);
      console.log(`[RATING] userRating.value set to: ${userRating.value}`);
      ratingSuccess.value = true;
      await fetchRating(); // Refresh average rating
      setTimeout(() => { ratingSuccess.value = false; }, 3000);
      alert('Rating submitted successfully!');
    } else {
      const errorData = await response.json().catch(() => ({ error: 'Rating failed' }));
      console.error('[RATING] Error response:', errorData);
      alert(`Failed to submit rating: ${errorData.error || 'Please try again'}`);
    }
  } catch (err) {
    console.error('[RATING] Exception:', err);
    alert('Failed to submit rating. Please check if backend is running.');
  } finally {
    ratingLoading.value = false;
  }
};

// Comments
const newComment = ref('');
const commentLoading = ref(false);

const submitComment = async () => {
  if (!token.value) {
    router.push('/login');
    return;
  }
  if (!newComment.value.trim()) {
    alert('Please enter a comment');
    return;
  }
  commentLoading.value = true;
  
  try {
    console.log(`[COMMENT] POST /recipes/${recipeId}/comments`);
      const response = await fetch(`${getApiUrl()}/recipes/${recipeId}/comments`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token.value}`
      },
      body: JSON.stringify({ content: newComment.value.trim() })
    });
    
    console.log(`[COMMENT] Response status: ${response.status}`);
    
    if (response.ok) {
      const data = await response.json().catch(() => ({}));
      console.log('[COMMENT] Success:', data);
      newComment.value = '';
      await fetchComments(); // Refresh comments list
    } else {
      const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
      console.error('[COMMENT] Error response:', errorData);
      alert(`Failed to post comment: ${errorData.error || 'Please try again'}`);
    }
  } catch (err) {
    console.error('[COMMENT] Exception:', err);
    alert(`Error: ${err.message || 'Failed to post comment. Please check if backend is running.'}`);
  } finally {
    commentLoading.value = false;
  }
};

// Buy Recipe
const buying = ref(false);
const handleBuyRecipe = async () => {
  if (!token.value) {
    router.push('/login');
    return;
  }
  if (!recipe.value) return;
  buying.value = true;
  
  try {
      const response = await fetch(`${getApiUrl()}/payment/initialize`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token.value}`
      },
      body: JSON.stringify({
        amount: recipe.value.price.toString(),
        email: userEmail.value || 'testuser@gmail.com',
        first_name: userName.value?.split(' ')[0] || 'Test',
        last_name: userName.value?.split(' ').slice(1).join(' ') || 'User',
        recipe_id: recipeId
      })
    });
    
    if (response.ok) {
      const data = await response.json();
      // Store recipe ID in sessionStorage as backup for payment success page
      if (process.client) {
        sessionStorage.setItem('pending_payment_recipe_id', recipeId.toString());
      }
      window.location.href = data.checkout_url;
    } else {
      const errorData = await response.json().catch(() => ({ error: 'Payment initialization failed' }));
      alert(`Payment Error: ${errorData.error || 'Failed to initialize payment. Please try again.'}`);
      console.error('Payment initialization error:', errorData);
    }
  } catch (err) {
    console.error('Error initializing payment:', err);
    alert('Failed to initialize payment. Please check your connection and try again.');
  } finally {
    buying.value = false;
  }
};

// Get real image based on recipe title
const getRecipeImage = (recipe) => {
  if (recipe?.thumbnail_url) return recipe.thumbnail_url;
  
  const title = recipe?.title?.toLowerCase() || '';
  
  if (title.includes('avocado') || title.includes('salad')) {
    return 'https://images.unsplash.com/photo-1546069901-ba9599a7e63c?auto=format&fit=crop&w=1200&q=80';
  }
  if (title.includes('chocolate') || title.includes('cake')) {
    return 'https://images.unsplash.com/photo-1578985545062-69928b1d9587?auto=format&fit=crop&w=1200&q=80';
  }
  if (title.includes('pasta')) {
    return 'https://images.unsplash.com/photo-1621996346565-e3dbc646d9a9?auto=format&fit=crop&w=1200&q=80';
  }
  if (title.includes('pizza')) {
    return 'https://images.unsplash.com/photo-1513104890138-7c749659a591?auto=format&fit=crop&w=1200&q=80';
  }
  if (title.includes('burger')) {
    return 'https://images.unsplash.com/photo-1568901346375-23c9450c58cd?auto=format&fit=crop&w=1200&q=80';
  }
  if (title.includes('sushi')) {
    return 'https://images.unsplash.com/photo-1579584425555-c3ce17fd4351?auto=format&fit=crop&w=1200&q=80';
  }
  
  return 'https://images.unsplash.com/photo-1495521821757-a1efb6729352?auto=format&fit=crop&w=1200&q=80';
};

// Format Date
const formatDate = (dateString) => {
  if (!dateString) return '';
  const date = new Date(dateString);
  return date.toLocaleDateString();
};

// Helper function to check if token is expired
const isTokenExpired = () => {
  if (!token.value) return true;
  try {
    const decoded = jwtDecode(token.value);
    if (decoded.exp && decoded.exp < Date.now() / 1000) {
      return true;
    }
    return false;
  } catch (err) {
    return true;
  }
};

// Check if user liked/bookmarked/purchased this recipe
const checkUserInteractions = async () => {
  if (!token.value || !recipe.value) {
    console.log('Skipping interaction check - no token or recipe');
    return;
  }
  
  // Check if token is expired
  if (isTokenExpired()) {
    console.warn('⚠️ Token is expired - please log in again');
    // Clear expired token
    token.value = null;
    return;
  }
  
  try {
    // Check like
    try {
      const likeCheck = await fetch(`${getApiUrl()}/recipes/${recipeId}/like/check`, {
        headers: { 'Authorization': `Bearer ${token.value}` }
      });
      if (likeCheck.ok) {
        const likeData = await likeCheck.json();
        isLiked.value = likeData.liked || false;
      } else if (likeCheck.status === 401) {
        console.warn('⚠️ Token invalid for like check - please log in again');
        token.value = null;
      }
    } catch (err) {
      console.warn('Error checking like:', err);
    }
    
    // Check bookmark
    try {
      const bookmarkCheck = await fetch(`${getApiUrl()}/recipes/${recipeId}/bookmark/check`, {
        headers: { 'Authorization': `Bearer ${token.value}` }
      });
      if (bookmarkCheck.ok) {
        const bookmarkData = await bookmarkCheck.json();
        isBookmarked.value = bookmarkData.bookmarked || false;
      } else if (bookmarkCheck.status === 401) {
        console.warn('⚠️ Token invalid for bookmark check - please log in again');
        token.value = null;
      }
    } catch (err) {
      console.warn('Error checking bookmark:', err);
    }
    
    // Check purchase - CRITICAL for paid recipes
    try {
      const purchaseCheck = await fetch(`${getApiUrl()}/recipes/${recipeId}/purchase/check`, {
        headers: { 'Authorization': `Bearer ${token.value}` }
      });
      
      if (purchaseCheck.ok) {
        const purchaseData = await purchaseCheck.json();
        const wasPurchased = hasPurchased.value;
        hasPurchased.value = purchaseData.purchased || false;
        
        console.log('🔍 Purchase check result:', {
          purchased: hasPurchased.value,
          wasPurchased: wasPurchased,
          recipeId: recipeId,
          price: recipe.value.price,
          isOwner: isOwner.value
        });
        
        // If purchase status changed to true, trigger content load
        if (hasPurchased.value && !wasPurchased) {
          console.log('✅ Purchase status changed to TRUE - will load content');
          // Trigger content load
          await loadContentIfAllowed();
        }
      } else if (purchaseCheck.status === 401) {
        console.warn('⚠️ Token invalid for purchase check - please log in again');
        const errorText = await purchaseCheck.text().catch(() => 'Unknown error');
        console.warn('Error details:', errorText);
        token.value = null;
      } else {
        console.warn('⚠️ Purchase check failed:', purchaseCheck.status, purchaseCheck.statusText);
        const errorText = await purchaseCheck.text().catch(() => 'Unknown error');
        console.warn('Error details:', errorText);
      }
    } catch (err) {
      console.error('❌ Error checking purchase:', err);
    }
    
    // Check user's rating
    try {
      const ratingCheck = await fetch(`${getApiUrl()}/recipes/${recipeId}/rate/check`, {
        headers: { 'Authorization': `Bearer ${token.value}` }
      });
      if (ratingCheck.ok) {
        const ratingData = await ratingCheck.json();
        console.log('⭐ Rating check response:', ratingData);
        if (ratingData.rated && ratingData.rating && ratingData.rating > 0) {
          userRating.value = parseInt(ratingData.rating);
          console.log(`⭐ User's previous rating set to: ${userRating.value} stars`);
        } else {
          userRating.value = 0;
          console.log('⭐ User has not rated this recipe yet');
        }
      } else if (ratingCheck.status === 401) {
        console.warn('⚠️ Token invalid for rating check - please log in again');
        token.value = null;
        userRating.value = 0;
      } else {
        console.warn('⚠️ Rating check failed:', ratingCheck.status);
        userRating.value = 0;
      }
    } catch (err) {
      console.warn('Error checking user rating:', err);
      userRating.value = 0;
    }
  } catch (err) {
    console.error('Error checking interactions:', err);
  }
};

// Load content if user has access (free, purchased, or owner)
const loadContentIfAllowed = async () => {
  if (!recipe.value) {
    console.log('Cannot load content - no recipe');
    return;
  }
  
  const canViewContent = recipe.value.price === 0 || hasPurchased.value || isOwner.value;
  
  console.log('📋 Checking content access:', {
    price: recipe.value.price,
    hasPurchased: hasPurchased.value,
    isOwner: isOwner.value,
    userId: userId.value,
    recipeUserId: recipe.value.user_id,
    canViewContent: canViewContent
  });
  
  if (canViewContent) {
    console.log('✅ User has access - loading ingredients and steps...');
    try {
      await Promise.all([
        fetchIngredients(),
        fetchSteps()
      ]);
      console.log('✅ Content loaded successfully');
    } catch (err) {
      console.error('❌ Error loading content:', err);
    }
  } else {
    console.log('❌ User does not have access to content');
  }
};

// Watch for recipe changes and check interactions
watch(() => recipe.value, async (newRecipe) => {
  if (newRecipe) {
    console.log('📝 Recipe loaded:', newRecipe.title);
    await fetchRecipeImages();
    
    // Fetch comments and ratings (these are public, always fetch)
    await Promise.all([
      fetchComments(),
      fetchRating()
    ]);
    
    // Check user interactions first (including purchase status)
    if (token.value) {
      await checkUserInteractions();
    }
    
    // Load content after checking interactions
    await loadContentIfAllowed();
    
    // Also check after a short delay to catch reactive updates
    setTimeout(() => {
      loadContentIfAllowed();
    }, 500);
  }
}, { immediate: false });

// Watch for purchase status changes - CRITICAL for showing content after payment
watch(() => hasPurchased.value, async (purchased, oldValue) => {
  console.log('💰 Purchase status changed:', { 
    purchased, 
    oldValue, 
    recipeId: recipeId,
    recipePrice: recipe.value?.price 
  });
  
  if (purchased && recipe.value) {
    // User just purchased, load content immediately
    console.log('✅ Purchase confirmed! Loading full recipe content...');
    await loadContentIfAllowed();
  }
}, { immediate: false });

// Watch for owner status changes
watch(() => isOwner.value, async (owner) => {
  console.log('👤 Owner status changed:', { owner, recipeId: recipeId });
  
  if (owner && recipe.value) {
    // Owner always has access, load content
    console.log('✅ User is owner - loading content...');
    await loadContentIfAllowed();
  }
});

// Load data on mount
onMounted(async () => {
  console.log('🚀 Component mounted, recipe:', recipe.value?.title || 'not loaded yet');
  
  // Comments and ratings are public - fetch them immediately regardless of recipe state
  console.log('📥 Fetching public data (comments, ratings)...');
  await Promise.all([
    fetchComments(),
    fetchRating()
  ]);
  
  // Wait for recipe to load from Apollo query
  if (recipe.value) {
    // If coming from payment, check purchase status with retries
    if (fromPayment.value && token.value) {
      console.log('💳 Coming from payment page - checking purchase status with retries...');
      
      // Check purchase status multiple times with delays to ensure backend has processed
      const checkPurchaseWithRetry = async (attempt = 1, maxAttempts = 5) => {
        console.log(`🔄 Purchase check attempt ${attempt}/${maxAttempts}...`);
        await checkUserInteractions();
        
        // If still not purchased and we have more attempts, retry
        if (!hasPurchased.value && attempt < maxAttempts) {
          const delay = attempt * 1000; // Increasing delay: 1s, 2s, 3s, 4s
          console.log(`⏳ Retrying purchase check in ${delay}ms...`);
          setTimeout(() => checkPurchaseWithRetry(attempt + 1, maxAttempts), delay);
        } else if (hasPurchased.value) {
          console.log('✅ Purchase confirmed after payment!');
        } else {
          console.warn('⚠️ Purchase not confirmed after all retries');
        }
      };
      
      // Start checking immediately
      await checkPurchaseWithRetry();
    } else if (token.value) {
      // Normal flow - check all interactions
      console.log('🔍 Checking user interactions...');
      await checkUserInteractions();
    }
    
    // Load content based on access
    await loadContentIfAllowed();
    
    // Also check after delays to catch reactive updates
    setTimeout(() => loadContentIfAllowed(), 500);
    setTimeout(() => loadContentIfAllowed(), 1500);
    
    // If coming from payment, check again after longer delay
    if (fromPayment.value) {
      setTimeout(() => {
        console.log('💳 Final purchase check after payment...');
        checkUserInteractions().then(() => loadContentIfAllowed());
      }, 3000);
    }
    
    // Comments and ratings are always visible
    await fetchComments();
    await fetchRating();
  } else {
    console.log('⏳ Waiting for recipe to load...');
  }
});
</script>
