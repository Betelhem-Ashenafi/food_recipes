<template>
  <div class="relative min-h-screen pb-20">
    <!-- Background -->
    <div class="fixed inset-0 z-0">
      <img 
        src="https://images.unsplash.com/photo-1504674900247-0877df9cc836?ixlib=rb-4.0.3&auto=format&fit=crop&w=2070&q=80" 
        alt="Food" 
        class="w-full h-full object-cover"
      >
      <div class="absolute inset-0 bg-gradient-to-b from-black/85 via-black/80 to-black/90"></div>
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
        <span class="block sm:inline"> {{ error.message }}</span>
      </div>
    </div>

    <!-- Recipe Content -->
    <div v-else-if="recipe" class="relative z-10 max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 pt-8">
      <!-- Hero Image -->
      <div class="mb-8 rounded-2xl overflow-hidden shadow-2xl">
        <img 
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

          <!-- Price Badge -->
          <div class="ml-4">
            <div v-if="recipe.price > 0" class="bg-emerald-500 text-white px-6 py-3 rounded-full text-xl font-bold shadow-lg">
              💎 {{ recipe.price }} Credits
            </div>
            <div v-else class="bg-green-500 text-white px-6 py-3 rounded-full text-xl font-bold shadow-lg">
              Free
            </div>
          </div>
        </div>

        <!-- Description -->
        <p class="text-gray-200 text-lg leading-relaxed mb-8">
          {{ recipe.description }}
        </p>

        <!-- Social Actions -->
        <div v-if="isAuthenticated" class="flex gap-4 mb-8 pb-8 border-b border-white/20">
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
          <span class="inline-flex items-center px-4 py-2 rounded-full text-sm font-medium bg-emerald-500/30 text-emerald-300 border border-emerald-400/50">
            {{ recipe.category?.name || 'Uncategorized' }}
          </span>
        </div>
      </div>

      <!-- Ingredients -->
      <div class="bg-white/10 backdrop-blur-lg border border-white/20 rounded-2xl p-8 mb-8">
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

      <!-- Steps -->
      <div class="bg-white/10 backdrop-blur-lg border border-white/20 rounded-2xl p-8 mb-8">
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
      <div v-if="isAuthenticated" class="bg-white/10 backdrop-blur-lg border border-white/20 rounded-2xl p-8 mb-8">
        <h2 class="text-2xl font-bold text-white mb-4">Rate this Recipe</h2>
        <div class="flex items-center gap-2">
          <button 
            v-for="star in 5" 
            :key="star"
            @click="submitRating(star)"
            :disabled="ratingLoading"
            class="transition-transform hover:scale-110"
          >
            <svg 
              class="w-10 h-10" 
              :fill="star <= userRating ? '#fbbf24' : 'none'" 
              :stroke="star <= userRating ? '#fbbf24' : '#9ca3af'" 
              viewBox="0 0 24 24"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
            </svg>
          </button>
        </div>
        <p v-if="ratingSuccess" class="text-emerald-400 text-sm mt-2">Rating submitted!</p>
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
const token = useCookie('auth_token');
const isAuthenticated = computed(() => !!token.value);

// Get user info from JWT
const userInfo = computed(() => {
  if (!token.value) return null;
  try {
    return jwtDecode(token.value);
  } catch {
    return null;
  }
});
const userName = computed(() => userInfo.value?.['https://hasura.io/jwt/claims']?.['x-hasura-user-name'] || 'User');
const userEmail = computed(() => userInfo.value?.['https://hasura.io/jwt/claims']?.['x-hasura-user-email'] || '');

// GraphQL Query for Recipe
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
      }
    }
  }
`;

// Use useQuery (reactive, no await needed)
const { result, loading: pending, error } = useQuery(query, { id: recipeId });
const recipe = computed(() => result.value?.recipes_by_pk);

// Fetch Ingredients (REST API)
const ingredientsData = ref([]);
const fetchIngredients = async () => {
  try {
    const response = await fetch(`http://localhost:8081/recipes/${recipeId}/ingredients`);
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
    const response = await fetch(`http://localhost:8081/recipes/${recipeId}/steps`);
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
  try {
    const response = await fetch(`http://localhost:8081/recipes/${recipeId}/comments`);
    if (response.ok) {
      commentsData.value = await response.json();
    }
  } catch (err) {
    console.error('Error fetching comments:', err);
  }
};

// Fetch Rating (REST API)
const ratingData = ref(null);
const fetchRating = async () => {
  try {
    const response = await fetch(`http://localhost:8081/recipes/${recipeId}/rate`);
    if (response.ok) {
      ratingData.value = await response.json();
    }
  } catch (err) {
    console.error('Error fetching rating:', err);
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
    alert('Please login to like recipes');
    return;
  }
  actionLoading.value = true;
  
  try {
    const method = isLiked.value ? 'DELETE' : 'POST';
    console.log(`[LIKE] ${method} /recipes/${recipeId}/like`);
    const response = await fetch(`http://localhost:8081/recipes/${recipeId}/like`, {
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
    alert('Please login to bookmark recipes');
    return;
  }
  actionLoading.value = true;
  
  try {
    const method = isBookmarked.value ? 'DELETE' : 'POST';
    console.log(`[BOOKMARK] ${method} /recipes/${recipeId}/bookmark`);
    const response = await fetch(`http://localhost:8081/recipes/${recipeId}/bookmark`, {
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

const submitRating = async (rating) => {
  if (!token.value) return;
  ratingLoading.value = true;
  ratingSuccess.value = false;
  
  try {
    const response = await fetch(`http://localhost:8081/recipes/${recipeId}/rate`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token.value}`
      },
      body: JSON.stringify({ rating })
    });
    
    if (response.ok) {
      userRating.value = rating;
      ratingSuccess.value = true;
      await fetchRating();
      setTimeout(() => { ratingSuccess.value = false; }, 3000);
    }
  } catch (err) {
    console.error('Error submitting rating:', err);
  } finally {
    ratingLoading.value = false;
  }
};

// Comments
const newComment = ref('');
const commentLoading = ref(false);

const submitComment = async () => {
  if (!token.value) {
    alert('Please login to comment');
    return;
  }
  if (!newComment.value.trim()) {
    alert('Please enter a comment');
    return;
  }
  commentLoading.value = true;
  
  try {
    console.log(`[COMMENT] POST /recipes/${recipeId}/comments`);
    const response = await fetch(`http://localhost:8081/recipes/${recipeId}/comments`, {
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
  if (!token.value || !recipe.value) return;
  buying.value = true;
  
  try {
    const response = await fetch('http://localhost:8081/payment/initialize', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token.value}`
      },
      body: JSON.stringify({
        amount: recipe.value.price.toString(),
        email: userEmail.value || 'user@example.com',
        first_name: userName.value?.split(' ')[0] || 'User',
        last_name: userName.value?.split(' ').slice(1).join(' ') || 'Name',
        recipe_id: recipeId
      })
    });
    
    if (response.ok) {
      const data = await response.json();
      window.location.href = data.checkout_url;
    }
  } catch (err) {
    console.error('Error initializing payment:', err);
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

// Check if user liked/bookmarked this recipe
const checkUserInteractions = async () => {
  if (!token.value || !recipe.value) return;
  
  try {
    // Check like
    const likeCheck = await fetch(`http://localhost:8081/recipes/${recipeId}/like/check`, {
      headers: { 'Authorization': `Bearer ${token.value}` }
    });
    if (likeCheck.ok) {
      const likeData = await likeCheck.json();
      isLiked.value = likeData.liked || false;
    }
    
    // Check bookmark
    const bookmarkCheck = await fetch(`http://localhost:8081/recipes/${recipeId}/bookmark/check`, {
      headers: { 'Authorization': `Bearer ${token.value}` }
    });
    if (bookmarkCheck.ok) {
      const bookmarkData = await bookmarkCheck.json();
      isBookmarked.value = bookmarkData.bookmarked || false;
    }
    
    // Check purchase
    const purchaseCheck = await fetch(`http://localhost:8081/recipes/${recipeId}/purchase/check`, {
      headers: { 'Authorization': `Bearer ${token.value}` }
    });
    if (purchaseCheck.ok) {
      const purchaseData = await purchaseCheck.json();
      hasPurchased.value = purchaseData.purchased || false;
    }
  } catch (err) {
    console.error('Error checking interactions:', err);
  }
};

// Watch for recipe changes and check interactions
watch(() => recipe.value, async (newRecipe) => {
  if (newRecipe && token.value) {
    await checkUserInteractions();
  }
}, { immediate: true });

// Load data on mount
onMounted(async () => {
  await fetchIngredients();
  await fetchSteps();
  await fetchComments();
  await fetchRating();
  if (token.value && recipe.value) {
    await checkUserInteractions();
  }
});
</script>
