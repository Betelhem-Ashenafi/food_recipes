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
          :src="getRecipeImage(recipe)" 
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
                <span v-if="ratingData">{{ ratingData.average_rating.toFixed(1) }}</span>
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
              @click="startInlineEdit"
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

        <div v-if="isInlineEditing && isOwner" class="mb-8 p-5 bg-black/25 border border-emerald-400/40 rounded-xl space-y-4">
          <h3 class="text-xl font-bold text-emerald-300">Edit Recipe Details</h3>

          <div>
            <label class="block text-sm text-gray-300 mb-1">Title</label>
            <input v-model="editForm.title" class="w-full px-4 py-2 border border-white/20 rounded-lg bg-black/20 text-white" />
          </div>

          <div>
            <label class="block text-sm text-gray-300 mb-1">Description</label>
            <textarea v-model="editForm.description" rows="4" class="w-full px-4 py-2 border border-white/20 rounded-lg bg-black/20 text-white"></textarea>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="block text-sm text-gray-300 mb-1">Preparation Time (min)</label>
              <input v-model.number="editForm.preparation_time" type="number" min="1" class="w-full px-4 py-2 border border-white/20 rounded-lg bg-black/20 text-white" />
            </div>
            <div>
              <label class="block text-sm text-gray-300 mb-1">Price</label>
              <input v-model.number="editForm.price" type="number" min="0" step="0.01" class="w-full px-4 py-2 border border-white/20 rounded-lg bg-black/20 text-white" />
            </div>
          </div>

          <div>
            <label class="block text-sm text-gray-300 mb-1">Category</label>
            <div class="grid grid-cols-2 md:grid-cols-3 gap-2">
              <button
                v-for="cat in allCategories"
                :key="cat.id"
                type="button"
                @click="editForm.category_id = cat.id"
                :class="[
                  'px-3 py-2 rounded-lg border text-sm transition-colors',
                  editForm.category_id === cat.id
                    ? 'bg-emerald-500/30 border-emerald-400 text-emerald-200'
                    : 'bg-black/20 border-white/20 text-gray-200 hover:border-white/40'
                ]"
              >
                {{ cat.name }}
              </button>
            </div>
          </div>

          <div>
            <h4 class="text-lg font-semibold text-white mb-2">Ingredients</h4>
            <div v-for="(ingredient, index) in editIngredients" :key="index" class="flex gap-2 mb-2">
              <input v-model="ingredient.name" placeholder="Name" class="flex-1 px-3 py-2 border border-white/20 rounded-lg bg-black/20 text-white" />
              <input v-model="ingredient.quantity" placeholder="Qty" class="w-24 px-3 py-2 border border-white/20 rounded-lg bg-black/20 text-white" />
              <select v-model="ingredient.unit_id" class="w-40 px-3 py-2 border border-white/20 rounded-lg bg-black/20 text-white">
                <option value="" disabled>Select unit</option>
                <option v-for="u in units" :key="u.id" :value="String(u.id)">{{ u.name }}</option>
              </select>
              <button @click="removeEditIngredient(index)" type="button" class="px-3 py-2 bg-red-500/20 text-red-300 rounded-lg">Remove</button>
            </div>
            <button @click="addEditIngredient" type="button" class="px-4 py-2 border border-dashed border-emerald-400/50 text-emerald-300 rounded-lg">+ Add Ingredient</button>
          </div>

          <div>
            <h4 class="text-lg font-semibold text-white mb-2">Preparation Steps</h4>
            <div v-for="(step, index) in editSteps" :key="index" class="mb-2">
              <textarea v-model="step.instruction" rows="2" placeholder="Step instruction" class="w-full px-3 py-2 border border-white/20 rounded-lg bg-black/20 text-white"></textarea>
              <button @click="removeEditStep(index)" type="button" class="mt-1 px-3 py-1 bg-red-500/20 text-red-300 rounded-lg">Remove</button>
            </div>
            <button @click="addEditStep" type="button" class="px-4 py-2 border border-dashed border-emerald-400/50 text-emerald-300 rounded-lg">+ Add Step</button>
          </div>

          <p v-if="inlineEditError" class="text-red-300 text-sm">{{ inlineEditError }}</p>

          <div class="flex gap-3 justify-end">
            <button @click="cancelInlineEdit" type="button" class="px-4 py-2 border border-white/30 text-white rounded-lg">Cancel</button>
            <button @click="saveInlineEdit" type="button" :disabled="inlineEditLoading" class="px-4 py-2 bg-emerald-600 text-white rounded-lg disabled:opacity-50">
              {{ inlineEditLoading ? 'Saving...' : 'Save Changes' }}
            </button>
          </div>
        </div>

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
            {{ isLiked ? 'Liked' : 'Like' }} ({{ likeCount }})
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
            v-if="recipe.price > 0 && !hasPurchased && !isOwner"
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
            <div v-if="recipeCategory?.image_url" class="w-8 h-8 rounded-full overflow-hidden border-2 border-emerald-400/50 flex-shrink-0">
              <img 
                :src="recipeCategory.image_url" 
                :alt="recipeCategory.name"
                class="w-full h-full object-cover"
                @error="(e) => e.target.style.display = 'none'"
              />
            </div>
            <!-- Category Name -->
            <span>{{ recipeCategory?.name || 'Uncategorized' }}</span>
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
            <span class="text-gray-200">{{ ingredient.quantity }} {{ ingredient.unit?.name || '' }} {{ ingredient.name }}</span>
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
        <div v-if="isAuthenticated && (recipe.price === 0 || hasPurchased || isOwner)">
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
          <p class="text-white">You must purchase this recipe to rate it.</p>
        </div>
        <p v-if="ratingSuccess && isAuthenticated && (recipe.price === 0 || hasPurchased || isOwner)" class="text-emerald-400 text-sm mt-2">Rating submitted!</p>
      </div>

      <!-- Comments Section -->
      <div class="bg-white/10 backdrop-blur-lg border border-white/20 rounded-2xl p-8">
        <h2 class="text-3xl font-bold text-white mb-6 flex items-center">
          <span class="text-emerald-400 mr-2">💬</span> Comments ({{ commentsData?.length || 0 }})
        </h2>

        <!-- Add Comment Form -->
        <div v-if="isAuthenticated && (recipe.price === 0 || hasPurchased || isOwner)" class="mb-8">
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
          <p class="text-white mb-3">You must purchase this recipe to comment.</p>
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
              <div class="ml-auto flex items-center gap-2">
                <button
                  v-if="isCommentOwner(comment) && editingCommentId !== comment.id"
                  @click="startEditComment(comment)"
                  class="text-emerald-300 hover:text-emerald-200 text-xs"
                >
                  Edit
                </button>
                <button
                  v-if="isCommentOwner(comment)"
                  @click="deleteComment(comment)"
                  :disabled="deletingCommentId === comment.id"
                  class="text-red-300 hover:text-red-200 text-xs disabled:opacity-50"
                >
                  {{ deletingCommentId === comment.id ? 'Deleting...' : 'Delete' }}
                </button>
                <span class="text-gray-400 text-sm">{{ formatDate(comment.created_at) }}</span>
              </div>
            </div>
            <div v-if="editingCommentId === comment.id" class="space-y-2">
              <textarea
                v-model="editingCommentText"
                rows="3"
                class="w-full px-3 py-2 border border-white/20 rounded-lg bg-black/20 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-emerald-500"
              ></textarea>
              <div class="flex gap-2 justify-end">
                <button
                  @click="cancelEditComment"
                  class="px-3 py-1 border border-white/30 rounded text-white text-xs"
                >
                  Cancel
                </button>
                <button
                  @click="saveEditedComment(comment)"
                  :disabled="!editingCommentText.trim() || commentLoading"
                  class="px-3 py-1 bg-emerald-500 text-white rounded text-xs disabled:opacity-50"
                >
                  {{ commentLoading ? 'Saving...' : 'Save' }}
                </button>
              </div>
            </div>
            <p v-else class="text-gray-200">{{ comment.content }}</p>
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
import { useApolloClient, useQuery } from '@vue/apollo-composable';
import gql from 'graphql-tag';
import { jwtDecode } from 'jwt-decode';

const route = useRoute();
const router = useRouter();
const recipeId = parseInt(route.params.id);

const config = useRuntimeConfig();
const getApiUrl = () => config.public?.apiUrl || 'http://localhost:8081';
// Check if coming from payment success page
const fromPayment = computed(() => route.query.fromPayment === 'true' || route.query.payment === 'success');
const paymentTxRef = computed(() => route.query.tx_ref || route.query.txRef || route.query.txref || route.query['amp;tx_ref']);
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
  if (claims && typeof claims === 'object' && claims !== null) {
    const id = claims['x-hasura-user-id'];
    if (id) return parseInt(id);
  }
  if (process.client) {
    const storedId = localStorage.getItem('user_id');
    if (storedId) return parseInt(storedId);
  }
  return null;
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
      created_at
      preparation_time
      user_id
      category_id
      user {
        name
      }
    }
  }
`

const CATEGORY_QUERY = gql`
  query GetCategory($id: Int!) {
    categories_by_pk(id: $id) {
      id
      name
      image_url
    }
  }
`;

const CATEGORIES_QUERY = gql`
  query GetCategories {
    categories(order_by: { name: asc }) {
      id
      name
      image_url
    }
  }
`;

const RECIPE_IMAGES_QUERY = gql`
  query GetRecipeImages($id: Int!) {
    recipe_images(where: { recipe_id: { _eq: $id } }, order_by: { id: asc }) {
      id
      recipe_id
      url
      is_featured
    }
  }
`;

const RECIPE_INGREDIENTS_QUERY = gql`
  query GetRecipeIngredients($id: Int!) {
    recipe_ingredients(where: { recipe_id: { _eq: $id } }, order_by: { id: asc }) {
      id
      name
      quantity
      unit_id
      unit {
        id
        name
      }
    }
  }
`;

const UNITS_QUERY = gql`
  query GetUnits {
    units(order_by: { id: asc }) {
      id
      name
    }
  }
`;

const RECIPE_STEPS_QUERY = gql`
  query GetRecipeSteps($id: Int!) {
    recipe_steps(where: { recipe_id: { _eq: $id } }, order_by: { step_number: asc }) {
      id
      step_number
      instruction
    }
  }
`;

const RECIPE_COMMENTS_QUERY = gql`
  query GetRecipeComments($id: Int!) {
    comments(where: { recipe_id: { _eq: $id } }, order_by: { created_at: desc }) {
      id
      user_id
      content
      created_at
      user {
        name
      }
    }
  }
`;

const RECIPE_COMMENTS_QUERY_CAMEL = gql`
  query GetRecipeCommentsCamel($id: Int!) {
    comments(where: { recipeId: { _eq: $id } }, order_by: { createdAt: desc }) {
      id
      userId
      content
      createdAt
      user {
        name
      }
    }
  }
`;

const RECIPE_RATING_QUERY = gql`
  query GetRecipeRating($id: Int!) {
    recipes_by_pk(id: $id) {
      id
      recipe_average_rating
      recipe_likes_count
    }
  }
`;

const CHECK_PURCHASES_RATINGS_QUERY = gql`
  query CheckPurchasesRatings($recipeId: Int!, $userId: Int!) {
    purchases(where: { recipe_id: { _eq: $recipeId }, user_id: { _eq: $userId } }) {
      id
      status
    }
    ratings(where: { recipe_id: { _eq: $recipeId }, user_id: { _eq: $userId } }) {
      rating
    }
  }
`;

const USER_PURCHASES_QUERY = gql`
  query GetUserPurchases($userId: Int!) {
    purchases(where: { user_id: { _eq: $userId } }) {
      id
      recipe_id
      status
    }
  }
`;

const CHECK_SOCIAL_QUERY = gql`
  query CheckSocial($recipeId: Int!, $userId: Int!) {
    likes(where: { recipe_id: { _eq: $recipeId }, user_id: { _eq: $userId } }) {
      recipe_id
    }
    bookmarks(where: { recipe_id: { _eq: $recipeId }, user_id: { _eq: $userId } }) {
      recipe_id
    }
  }
`;

const LIKE_RECIPE_MUTATION = gql`
  mutation LikeRecipe($recipeId: Int!) {
    insert_likes_one(object: { recipe_id: $recipeId }) {
      recipe_id
    }
  }
`;

const UNLIKE_RECIPE_MUTATION = gql`
  mutation UnlikeRecipe($recipeId: Int!, $userId: Int!) {
    delete_likes(where: { recipe_id: { _eq: $recipeId }, user_id: { _eq: $userId } }) {
      affected_rows
    }
  }
`;

const BOOKMARK_RECIPE_MUTATION = gql`
  mutation BookmarkRecipe($recipeId: Int!) {
    insert_bookmarks_one(object: { recipe_id: $recipeId }) {
      recipe_id
    }
  }
`;

const UNBOOKMARK_RECIPE_MUTATION = gql`
  mutation UnbookmarkRecipe($recipeId: Int!, $userId: Int!) {
    delete_bookmarks(where: { recipe_id: { _eq: $recipeId }, user_id: { _eq: $userId } }) {
      affected_rows
    }
  }
`;

const RATE_RECIPE_MUTATION = gql`
  mutation RateRecipe($recipeId: Int!, $rating: Int!) {
    insert_ratings_one(
      object: { recipe_id: $recipeId, rating: $rating }
      on_conflict: { constraint: ratings_recipe_user_unique, update_columns: [rating] }
    ) {
      rating
    }
  }
`;

const POST_COMMENT_MUTATION = gql`
  mutation PostComment($recipeId: Int!, $content: String!) {
    insert_comments_one(object: { recipe_id: $recipeId, content: $content }) {
      id
    }
  }
`;

const UPDATE_COMMENT_MUTATION = gql`
  mutation UpdateComment($id: Int!, $content: String!) {
    update_comments(where: { id: { _eq: $id } }, _set: { content: $content }) {
      affected_rows
    }
  }
`;

const DELETE_COMMENT_MUTATION = gql`
  mutation DeleteComment($id: Int!) {
    delete_comments(where: { id: { _eq: $id } }) {
      affected_rows
    }
  }
`;

const UPDATE_RECIPE_INLINE_MUTATION = gql`
  mutation UpdateRecipeInline(
    $recipeId: Int!
    $recipe: recipes_set_input!
    $ingredients: [recipe_ingredients_insert_input!]!
    $steps: [recipe_steps_insert_input!]!
  ) {
    update_recipes_by_pk(pk_columns: { id: $recipeId }, _set: $recipe) {
      id
    }
    delete_recipe_ingredients(where: { recipe_id: { _eq: $recipeId } }) {
      affected_rows
    }
    insert_recipe_ingredients(objects: $ingredients) {
      affected_rows
    }
    delete_recipe_steps(where: { recipe_id: { _eq: $recipeId } }) {
      affected_rows
    }
    insert_recipe_steps(objects: $steps) {
      affected_rows
    }
  }
`;

const POST_COMMENT_MUTATION_CAMEL = gql`
  mutation PostCommentCamel($recipeId: Int!, $userId: Int!, $content: String!) {
    insert_comments_one(object: { recipeId: $recipeId, userId: $userId, content: $content }) {
      id
    }
  }
`;

// Use useQuery with proper error handling
// Skip query if recipeId is invalid
const { result, loading: pending, error, refetch: refetchRecipe } = useQuery(
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

const categoryId = computed(() => {
  const raw = recipe.value?.category_id;
  const parsed = typeof raw === 'string' ? parseInt(raw, 10) : raw;
  return Number.isFinite(parsed) && parsed > 0 ? parsed : null;
});

const categoryQueryVars = computed(() => ({ id: categoryId.value ?? -1 }));

const { result: categoryResult } = useQuery(
  CATEGORY_QUERY,
  () => categoryQueryVars.value,
  {
    skip: () => !categoryId.value,
    errorPolicy: 'all',
    fetchPolicy: 'cache-and-network',
    returnPartialData: true
  }
);

const { result: categoriesResult } = useQuery(
  CATEGORIES_QUERY,
  () => ({}),
  {
    errorPolicy: 'all',
    fetchPolicy: 'cache-and-network',
    returnPartialData: true
  }
);

const recipeCategory = computed(() => {
  const category = categoryResult.value?.categories_by_pk;
  return category || null;
});

const allCategories = computed(() => categoriesResult.value?.categories || []);

const { client } = useApolloClient();

const units = ref([]);
const fetchUnits = async () => {
  try {
    const result = await client.query({
      query: UNITS_QUERY,
      fetchPolicy: 'network-only'
    });
    units.value = result?.data?.units || [];
  } catch (err) {
    console.error('Error fetching units:', err);
    units.value = [];
  }
};

// Fetch Recipe Images (REST API)
const recipeImages = ref([]);
const currentImageIndex = ref(0);

const fetchRecipeImages = async () => {
  try {
    const result = await client.query({
      query: RECIPE_IMAGES_QUERY,
      variables: { id: recipeId },
      fetchPolicy: 'network-only'
    });
    const images = result?.data?.recipe_images || [];
    if (images.length > 0) {
      recipeImages.value = images;
      const featuredIndex = recipeImages.value.findIndex(img => img.is_featured);
      currentImageIndex.value = featuredIndex >= 0 ? featuredIndex : 0;
      return;
    }
  } catch (err) {
    console.error('Error fetching images:', err);
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

// Fetch Ingredients (GraphQL)
const ingredientsData = ref([]);
const fetchIngredients = async () => {
  try {
    const result = await client.query({
      query: RECIPE_INGREDIENTS_QUERY,
      variables: { id: recipeId },
      fetchPolicy: 'network-only'
    });
    ingredientsData.value = result?.data?.recipe_ingredients || [];
  } catch (err) {
    console.error('Error fetching ingredients:', err);
    ingredientsData.value = [];
  }
};

const isInlineEditing = ref(false);
const inlineEditLoading = ref(false);
const inlineEditError = ref('');
const editForm = ref({
  title: '',
  description: '',
  preparation_time: 1,
  price: 0,
  category_id: null
});
const editIngredients = ref([]);
const editSteps = ref([]);

const startInlineEdit = async () => {
  if (!isOwner.value || !recipe.value) return;
  if (units.value.length === 0) {
    await fetchUnits();
  }

  editForm.value = {
    title: recipe.value.title || '',
    description: recipe.value.description || '',
    preparation_time: Number(recipe.value.preparation_time || 1),
    price: Number(recipe.value.price || 0),
    category_id: recipe.value.category_id || null
  };

  editIngredients.value = (ingredientsData.value || []).map((ingredient) => ({
    name: ingredient.name || '',
    quantity: ingredient.quantity || '',
    unit_id: ingredient.unit_id ? String(ingredient.unit_id) : String(ingredient.unit?.id || '')
  }));
  if (editIngredients.value.length === 0) {
    editIngredients.value = [{ name: '', quantity: '', unit_id: '' }];
  }

  editSteps.value = (stepsData.value || []).map((step) => ({ instruction: step.instruction || '' }));
  if (editSteps.value.length === 0) {
    editSteps.value = [{ instruction: '' }];
  }

  inlineEditError.value = '';
  isInlineEditing.value = true;
};

const cancelInlineEdit = () => {
  isInlineEditing.value = false;
  inlineEditError.value = '';
};

const addEditIngredient = () => {
  editIngredients.value.push({ name: '', quantity: '', unit_id: '' });
};

const removeEditIngredient = (index) => {
  if (editIngredients.value.length > 1) {
    editIngredients.value.splice(index, 1);
  }
};

const addEditStep = () => {
  editSteps.value.push({ instruction: '' });
};

const removeEditStep = (index) => {
  if (editSteps.value.length > 1) {
    editSteps.value.splice(index, 1);
  }
};

const saveInlineEdit = async () => {
  if (!isOwner.value || !recipe.value) return;

  inlineEditError.value = '';

  if (!editForm.value.title?.trim() || editForm.value.title.trim().length < 3) {
    inlineEditError.value = 'Title must be at least 3 characters.';
    return;
  }
  if (!editForm.value.description?.trim() || editForm.value.description.trim().length < 10) {
    inlineEditError.value = 'Description must be at least 10 characters.';
    return;
  }
  if (Number(editForm.value.preparation_time) < 1) {
    inlineEditError.value = 'Preparation time must be at least 1 minute.';
    return;
  }
  if (Number(editForm.value.price) < 0) {
    inlineEditError.value = 'Price cannot be negative.';
    return;
  }
  if (!editForm.value.category_id) {
    inlineEditError.value = 'Please select a category.';
    return;
  }

  const validIngredients = editIngredients.value.filter((ingredient) => ingredient.name?.trim());
  if (validIngredients.length === 0) {
    inlineEditError.value = 'Please add at least one ingredient.';
    return;
  }
  if (validIngredients.some((ingredient) => !ingredient.unit_id)) {
    inlineEditError.value = 'Please select a unit for every ingredient.';
    return;
  }

  const validSteps = editSteps.value.filter((step) => step.instruction?.trim());
  if (validSteps.length === 0) {
    inlineEditError.value = 'Please add at least one preparation step.';
    return;
  }

  inlineEditLoading.value = true;
  try {
    await client.mutate({
      mutation: UPDATE_RECIPE_INLINE_MUTATION,
      variables: {
        recipeId,
        recipe: {
          title: editForm.value.title.trim(),
          description: editForm.value.description.trim(),
          preparation_time: parseInt(editForm.value.preparation_time, 10),
          price: parseFloat(editForm.value.price) || 0,
          category_id: parseInt(editForm.value.category_id, 10)
        },
        ingredients: validIngredients.map((ingredient) => ({
          recipe_id: recipeId,
          name: ingredient.name.trim(),
          quantity: ingredient.quantity || '',
          unit_id: parseInt(ingredient.unit_id, 10)
        })),
        steps: validSteps.map((step, index) => ({
          recipe_id: recipeId,
          step_number: index + 1,
          instruction: step.instruction.trim()
        }))
      }
    });

    await refetchRecipe({ id: recipeId });
    await Promise.all([fetchIngredients(), fetchSteps(), fetchRating()]);
    isInlineEditing.value = false;
  } catch (err) {
    console.error('[INLINE_EDIT] Exception:', err);
    inlineEditError.value = err.message || 'Failed to update recipe.';
  } finally {
    inlineEditLoading.value = false;
  }
};

// Fetch Steps (GraphQL)
const stepsData = ref([]);
const fetchSteps = async () => {
  try {
    const result = await client.query({
      query: RECIPE_STEPS_QUERY,
      variables: { id: recipeId },
      fetchPolicy: 'network-only'
    });
    const data = result?.data?.recipe_steps || [];
    const sortedSteps = [...data].sort((a, b) => a.step_number - b.step_number);
    stepsData.value = sortedSteps;
  } catch (err) {
    console.error('Error fetching steps:', err);
    stepsData.value = [];
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
    let result;
    try {
      result = await client.query({
        query: RECIPE_COMMENTS_QUERY,
        variables: { id: recipeId },
        fetchPolicy: 'network-only'
      });
    } catch (err) {
      const message = err?.message || '';
      if (message.includes('recipe_id') || message.includes('created_at') || message.includes('user_id')) {
        result = await client.query({
          query: RECIPE_COMMENTS_QUERY_CAMEL,
          variables: { id: recipeId },
          fetchPolicy: 'network-only'
        });
      } else {
        throw err;
      }
    }
    const comments = result?.data?.comments || [];
    commentsData.value = comments.map((comment) => ({
      id: comment.id,
      user_id: comment.user_id ?? comment.userId,
      user_name: comment.user?.name || 'User',
      content: comment.content,
      created_at: comment.created_at ?? comment.createdAt
    }));
    console.log(`[COMMENTS] Loaded ${commentsData.value.length} comments`);
  } catch (err) {
    console.error('[COMMENTS] Error fetching comments:', err);
    commentsData.value = [];
  }
};

// Fetch Rating (REST API)
const ratingData = ref(null);
const likeCount = ref(0);
const fetchRating = async () => {
  if (!recipeId) {
    console.warn('[RATING] No recipe ID, skipping fetch');
    return;
  }
  try {
    console.log(`[RATING] Fetching rating for recipe ${recipeId}`);
    const result = await client.query({
      query: RECIPE_RATING_QUERY,
      variables: { id: recipeId },
      fetchPolicy: 'network-only'
    });

    const recipeStats = result?.data?.recipes_by_pk;
    const averageRating = Number(recipeStats?.recipe_average_rating ?? 0);
    likeCount.value = Number(recipeStats?.recipe_likes_count ?? 0);
    ratingData.value = { average_rating: averageRating };
    console.log(`[RATING] Loaded rating:`, ratingData.value);
  } catch (err) {
    console.error('[RATING] Error fetching rating/likes stats:', err);
    ratingData.value = null;
    likeCount.value = 0;
  }
};

// Social Features State
const isLiked = ref(false);
const isBookmarked = ref(false);
const actionLoading = ref(false);
const hasPurchased = ref(false);

const verifyPurchaseFromPayment = async () => {
  if (!fromPayment.value || !token.value || !recipeId) return false;
  try {
    const queryParam = paymentTxRef.value ? `tx_ref=${paymentTxRef.value}` : `recipe_id=${recipeId}`;
    const response = await fetch(`${getApiUrl()}/payment/verify?${queryParam}`, {
      headers: {
        Authorization: `Bearer ${token.value}`
      }
    });
    if (!response.ok) {
      return false;
    }
    const data = await response.json();
    if (data.status === 'success') {
      hasPurchased.value = true;
      return true;
    }
  } catch (err) {
    console.error('Error verifying purchase from payment:', err);
  }
  return false;
};

// Toggle Like
const toggleLike = async () => {
  if (!token.value) {
    router.push('/login');
    return;
  }
  actionLoading.value = true;
  
  try {
    if (!userId.value) {
      router.push('/login');
      return;
    }
    if (isLiked.value) {
      await client.mutate({
        mutation: UNLIKE_RECIPE_MUTATION,
        variables: { recipeId, userId: userId.value }
      });
      isLiked.value = false;
    } else {
      await client.mutate({
        mutation: LIKE_RECIPE_MUTATION,
        variables: { recipeId, userId: userId.value }
      });
      isLiked.value = true;
    }
    await fetchRating();
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
    if (!userId.value) {
      router.push('/login');
      return;
    }
    if (isBookmarked.value) {
      await client.mutate({
        mutation: UNBOOKMARK_RECIPE_MUTATION,
        variables: { recipeId, userId: userId.value }
      });
      isBookmarked.value = false;
    } else {
      await client.mutate({
        mutation: BOOKMARK_RECIPE_MUTATION,
        variables: { recipeId, userId: userId.value }
      });
      isBookmarked.value = true;
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
    if (!userId.value) {
      router.push('/login');
      return;
    }
    await client.mutate({
      mutation: RATE_RECIPE_MUTATION,
      variables: { recipeId, rating }
    });
    userRating.value = parseInt(rating);
    console.log(`[RATING] userRating.value set to: ${userRating.value}`);
    ratingSuccess.value = true;
    await fetchRating();
    setTimeout(() => { ratingSuccess.value = false; }, 3000);
    alert('Rating submitted successfully!');
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
const editingCommentId = ref(null);
const editingCommentText = ref('');
const deletingCommentId = ref(null);

const isCommentOwner = (comment) => {
  if (!isAuthenticated.value || !userId.value || !comment) return false;
  const commentUserId = typeof comment.user_id === 'string' ? parseInt(comment.user_id, 10) : comment.user_id;
  const currentUserId = typeof userId.value === 'string' ? parseInt(userId.value, 10) : userId.value;
  return Number(commentUserId) === Number(currentUserId);
};

const startEditComment = (comment) => {
  if (!isCommentOwner(comment)) return;
  editingCommentId.value = comment.id;
  editingCommentText.value = comment.content || '';
};

const cancelEditComment = () => {
  editingCommentId.value = null;
  editingCommentText.value = '';
};

const saveEditedComment = async (comment) => {
  if (!token.value) {
    router.push('/login');
    return;
  }
  if (!isCommentOwner(comment)) {
    alert('You can only edit your own comments.');
    return;
  }
  if (!editingCommentText.value.trim()) {
    alert('Comment cannot be empty.');
    return;
  }

  commentLoading.value = true;
  try {
    const result = await client.mutate({
      mutation: UPDATE_COMMENT_MUTATION,
      variables: { id: comment.id, content: editingCommentText.value.trim() }
    });
    const affectedRows = result?.data?.update_comments?.affected_rows ?? 0;
    if (affectedRows < 1) {
      throw new Error('Comment update was not permitted.');
    }
    cancelEditComment();
    await fetchComments();
  } catch (err) {
    console.error('[COMMENT] Edit exception:', err);
    alert(`Error: ${err.message || 'Failed to edit comment.'}`);
  } finally {
    commentLoading.value = false;
  }
};

const deleteComment = async (comment) => {
  if (!token.value) {
    router.push('/login');
    return;
  }
  if (!isCommentOwner(comment)) {
    alert('You can only delete your own comments.');
    return;
  }
  if (!confirm('Delete this comment?')) {
    return;
  }

  deletingCommentId.value = comment.id;
  try {
    const result = await client.mutate({
      mutation: DELETE_COMMENT_MUTATION,
      variables: { id: comment.id }
    });
    const affectedRows = result?.data?.delete_comments?.affected_rows ?? 0;
    if (affectedRows < 1) {
      throw new Error('Comment delete was not permitted.');
    }
    if (editingCommentId.value === comment.id) {
      cancelEditComment();
    }
    await fetchComments();
  } catch (err) {
    console.error('[COMMENT] Delete exception:', err);
    alert(`Error: ${err.message || 'Failed to delete comment.'}`);
  } finally {
    deletingCommentId.value = null;
  }
};

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
    if (!userId.value) {
      router.push('/login');
      return;
    }
    try {
      await client.mutate({
        mutation: POST_COMMENT_MUTATION,
        variables: { recipeId, userId: userId.value, content: newComment.value.trim() }
      });
    } catch (err) {
      const message = err?.message || '';
      if (message.includes('recipe_id') || message.includes('user_id')) {
        await client.mutate({
          mutation: POST_COMMENT_MUTATION_CAMEL,
          variables: { recipeId, userId: userId.value, content: newComment.value.trim() }
        });
      } else {
        throw err;
      }
    }
    newComment.value = '';
    await fetchComments();
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
      const rawText = await response.text().catch(() => 'Payment initialization failed');
      let errorMessage = 'Failed to initialize payment. Please try again.';
      try {
        const errorData = JSON.parse(rawText);
        errorMessage = errorData.error || errorData.message || errorMessage;
        console.error('Payment initialization error:', errorData);
      } catch {
        errorMessage = rawText || errorMessage;
        console.error('Payment initialization error:', rawText);
      }
      alert(`Payment Error: ${errorMessage}`);
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
  if (recipeImages.value?.length > 0) {
    const featured = recipeImages.value.find((img) => img.is_featured);
    return featured?.url || recipeImages.value[0]?.url;
  }
  
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
    if (!userId.value) {
      router.push('/login');
      return;
    }
    const result = await client.query({
      query: CHECK_PURCHASES_RATINGS_QUERY,
      variables: { recipeId, userId: userId.value },
      fetchPolicy: 'network-only'
    });
    let purchases = result?.data?.purchases || [];
    const ratings = result?.data?.ratings || [];

    if (purchases.length === 0) {
      try {
        const userPurchasesResult = await client.query({
          query: USER_PURCHASES_QUERY,
          variables: { userId: userId.value },
          fetchPolicy: 'network-only'
        });
        const allPurchases = userPurchasesResult?.data?.purchases || [];
        purchases = allPurchases.filter((p) => p?.recipe_id === recipeId);
      } catch (err) {
        console.warn('User purchases fallback query failed:', err);
      }
    }

    try {
      const socialResult = await client.query({
        query: CHECK_SOCIAL_QUERY,
        variables: { recipeId, userId: userId.value },
        fetchPolicy: 'network-only'
      });
      const likes = socialResult?.data?.likes || [];
      const bookmarks = socialResult?.data?.bookmarks || [];
      isLiked.value = likes.length > 0;
      isBookmarked.value = bookmarks.length > 0;
    } catch (err) {
      console.warn('Social query failed (likes/bookmarks not tracked):', err);
      isLiked.value = false;
      isBookmarked.value = false;
    }

    const wasPurchased = hasPurchased.value;
    const isSuccessfulPurchase = (status) => {
      const normalized = String(status || '').toLowerCase().trim();
      return normalized === 'success' || normalized === 'paid' || normalized === 'completed';
    };
    const hasAnyPurchase = purchases.length > 0;
    hasPurchased.value = purchases.some((p) => isSuccessfulPurchase(p?.status)) || hasAnyPurchase;

    console.log('🔍 Purchase check result:', {
      purchased: hasPurchased.value,
      wasPurchased: wasPurchased,
      recipeId: recipeId,
      price: recipe.value.price,
      isOwner: isOwner.value
    });

    if (hasPurchased.value && !wasPurchased) {
      console.log('✅ Purchase status changed to TRUE - will load content');
      await loadContentIfAllowed();
    }

    if (ratings.length > 0 && ratings[0].rating) {
      userRating.value = parseInt(ratings[0].rating);
      console.log(`⭐ User's previous rating set to: ${userRating.value} stars`);
    } else {
      userRating.value = 0;
      console.log('⭐ User has not rated this recipe yet');
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
    
    // If returning from payment, attempt verification before checking interactions
    await verifyPurchaseFromPayment();

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
  await fetchUnits();
  
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
