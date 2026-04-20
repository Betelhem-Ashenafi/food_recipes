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
          :src="getRecipeImage()" 
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
        <div class="flex flex-wrap gap-3 sm:gap-4 mb-8 pb-8 border-b border-white/20">
          <button 
            @click="toggleLike"
            :disabled="actionLoading"
            class="flex-1 min-w-[130px] sm:min-w-0 flex items-center justify-center gap-2 px-4 py-3 rounded-lg transition-colors font-medium text-sm sm:text-base"
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
            class="flex-1 min-w-[130px] sm:min-w-0 flex items-center justify-center gap-2 px-4 py-3 rounded-lg transition-colors font-medium text-sm sm:text-base"
            :class="isBookmarked ? 'bg-yellow-500 text-white hover:bg-yellow-600' : 'bg-white/10 text-gray-300 hover:bg-white/20 border border-white/20'"
          >
            <svg class="w-5 h-5" :fill="isBookmarked ? 'currentColor' : 'none'" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z" />
            </svg>
            {{ isBookmarked ? 'Bookmarked' : 'Bookmark' }}
          </button>

          <button 
            v-if="Number(recipe.price) > 0 && !hasPurchased && !isOwner"
            @click="handleBuyRecipe"
            :disabled="buying"
            class="w-full sm:flex-1 sm:w-auto flex items-center justify-center gap-2 px-4 py-3 bg-gradient-to-r from-emerald-600 to-teal-600 text-white rounded-lg hover:from-emerald-500 hover:to-teal-500 transition-all font-bold shadow-lg disabled:opacity-50"
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
            <div v-if="shouldShowRecipeCategoryImage" class="w-8 h-8 rounded-full overflow-hidden border-2 border-emerald-400/50 flex-shrink-0">
              <img
                :src="normalizedRecipeCategoryImageUrl"
                :alt="recipeCategory.name"
                class="w-full h-full object-cover"
                @error="recipeCategoryImageFailed = true"
              />
            </div>
            <div v-else class="w-8 h-8 rounded-full overflow-hidden border-2 border-emerald-400/50 flex-shrink-0 bg-emerald-500/40 flex items-center justify-center text-white text-xs font-bold">
              {{ recipeCategoryInitial }}
            </div>
            <!-- Category Name -->
            <span>{{ recipeCategory?.name || 'Uncategorized' }}</span>
          </div>
        </div>
      </div>

      <!-- Ingredients - Only show if free, purchased, or owner -->
      <div v-if="canViewContent" class="bg-white/10 backdrop-blur-lg border border-white/20 rounded-2xl p-8 mb-8">
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
      <div v-if="canViewContent" class="bg-white/10 backdrop-blur-lg border border-white/20 rounded-2xl p-8 mb-8">
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
        <div v-if="canInteract">
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
          <p v-if="!isAuthenticated" class="text-white">Please log in to rate this recipe.</p>
          <p v-else class="text-white">You must purchase this recipe to rate it.</p>
        </div>
        <p v-if="ratingSuccess && canInteract" class="text-emerald-400 text-sm mt-2">Rating submitted!</p>
      </div>

      <!-- Comments Section -->
      <div class="bg-white/10 backdrop-blur-lg border border-white/20 rounded-2xl p-8">
        <h2 class="text-3xl font-bold text-white mb-6 flex items-center">
          <span class="text-emerald-400 mr-2">💬</span> Comments ({{ commentsData?.length || 0 }})
        </h2>

        <!-- Add Comment Form -->
        <div v-if="canInteract" class="mb-8">
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
          <p v-if="!isAuthenticated" class="text-white mb-3">Please log in to comment on this recipe.</p>
          <p v-else class="text-white mb-3">You must purchase this recipe to comment.</p>
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
import { ref, computed, watch } from 'vue';
import { useApolloClient, useQuery } from '@vue/apollo-composable';
import { jwtDecode } from 'jwt-decode';
import {
  RECIPE_QUERY,
  GET_ACCESSIBLE_RECIPE_QUERY,
  CATEGORY_QUERY,
  CATEGORIES_QUERY,
  RECIPE_IMAGES_QUERY,
  RECIPE_INGREDIENTS_QUERY,
  RECIPE_STEPS_QUERY,
  RECIPE_COMMENTS_QUERY,
  RECIPE_COMMENTS_QUERY_CAMEL,
  RECIPE_RATING_QUERY,
  UNITS_QUERY,
  CHECK_PURCHASES_RATINGS_QUERY,
  CHECK_SOCIAL_QUERY,
  LIKE_RECIPE_MUTATION,
  UNLIKE_RECIPE_MUTATION,
  BOOKMARK_RECIPE_MUTATION,
  UNBOOKMARK_RECIPE_MUTATION,
  RATE_RECIPE_MUTATION,
  POST_COMMENT_MUTATION,
  POST_COMMENT_MUTATION_CAMEL,
  UPDATE_COMMENT_MUTATION,
  DELETE_COMMENT_MUTATION,
  UPDATE_RECIPE_INLINE_MUTATION
} from '~/utils/recipe-detail.gql';
import { initializePaymentAction, verifyPaymentAction } from '~/utils/payment-actions';

const route = useRoute();
const router = useRouter();
const recipeId = Number.parseInt(route.params.id, 10) || 0;

const token = useCookie('auth_token');
const isAuthenticated = computed(() => !!token.value);
const fromPayment = computed(() => route.query.fromPayment === 'true' || route.query.payment === 'success');
const paymentTxRef = computed(() => route.query.tx_ref || route.query.txRef || route.query.txref || route.query['amp;tx_ref']);

const { client } = useApolloClient();

const { result, loading: pending, error, refetch: refetchRecipe } = useQuery(
  RECIPE_QUERY,
  () => ({ id: recipeId }),
  {
    skip: () => !recipeId,
    errorPolicy: 'all',
    fetchPolicy: 'cache-and-network'
  }
);

const recipe = computed(() => result.value?.recipes_by_pk || null);

const categoryId = computed(() => {
  const raw = recipe.value?.category_id;
  const parsed = typeof raw === 'string' ? Number.parseInt(raw, 10) : raw;
  return Number.isFinite(parsed) && parsed > 0 ? parsed : null;
});

const { result: categoryResult } = useQuery(
  CATEGORY_QUERY,
  () => ({ id: categoryId.value ?? -1 }),
  {
    skip: () => !categoryId.value,
    errorPolicy: 'all',
    fetchPolicy: 'cache-and-network'
  }
);

const { result: categoriesResult } = useQuery(
  CATEGORIES_QUERY,
  () => ({}),
  { errorPolicy: 'all', fetchPolicy: 'cache-and-network' }
);

const recipeCategory = computed(() => categoryResult.value?.categories_by_pk || null);
const allCategories = computed(() => categoriesResult.value?.categories || []);
const recipeCategoryImageFailed = ref(false);

const normalizeImageUrl = (url) => {
  const raw = String(url || '').trim();
  if (!raw) return '';
  if (/^https?:\/\//i.test(raw)) return raw;
  if (/^\/\//.test(raw)) return `https:${raw}`;
  return `https://${raw}`;
};

const normalizedRecipeCategoryImageUrl = computed(() => normalizeImageUrl(recipeCategory.value?.image_url));

const shouldShowRecipeCategoryImage = computed(() => {
  return !!normalizedRecipeCategoryImageUrl.value && !recipeCategoryImageFailed.value;
});

const recipeCategoryInitial = computed(() => {
  const name = String(recipeCategory.value?.name || '').trim();
  return name ? name.charAt(0).toUpperCase() : 'C';
});

const getClaims = () => {
  if (!token.value) return {};
  try {
    const decoded = jwtDecode(token.value);
    return decoded?.['https://hasura.io/jwt/claims'] || {};
  } catch {
    return {};
  }
};

const userId = computed(() => {
  const claimId = getClaims()['x-hasura-user-id'];
  if (claimId) return Number.parseInt(claimId, 10);
  if (process.client) {
    const storedId = localStorage.getItem('user_id');
    if (storedId) return Number.parseInt(storedId, 10);
  }
  return null;
});

const userName = computed(() => getClaims()['x-hasura-user-name'] || 'User');
const userEmail = computed(() => getClaims()['x-hasura-user-email'] || '');

const isOwner = computed(() => {
  if (!recipe.value || !userId.value) return false;
  return Number(recipe.value.user_id) === Number(userId.value);
});

const dbHasAccess = ref(false);

const canViewContent = computed(() => {
  if (!recipe.value) return false;
  return Number(recipe.value.price) === 0 || hasPurchased.value || isOwner.value || dbHasAccess.value;
});

const canInteract = computed(() => isAuthenticated.value && canViewContent.value);

const units = ref([]);
const recipeImages = ref([]);
const currentImageIndex = ref(0);
const ingredientsData = ref([]);
const stepsData = ref([]);
const commentsData = ref([]);
const ratingData = ref(null);
const likeCount = ref(0);

const isLiked = ref(false);
const isBookmarked = ref(false);
const hasPurchased = ref(false);

const actionLoading = ref(false);
const buying = ref(false);

const userRating = ref(0);
const currentUserRating = computed(() => Number(userRating.value) || 0);
const ratingLoading = ref(false);
const ratingSuccess = ref(false);

const newComment = ref('');
const commentLoading = ref(false);
const editingCommentId = ref(null);
const editingCommentText = ref('');
const deletingCommentId = ref(null);

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
const editIngredients = ref([{ name: '', quantity: '', unit_id: '' }]);
const editSteps = ref([{ instruction: '' }]);
const hasAutoOpenedInlineEdit = ref(false);

const createIngredient = () => ({ name: '', quantity: '', unit_id: '' });
const createStep = () => ({ instruction: '' });

const requireAuth = () => {
  if (!token.value || !userId.value) {
    router.push('/login');
    return false;
  }
  return true;
};

const fetchList = async (query, targetRef, options = {}) => {
  try {
    const resultQuery = await client.query({
      query,
      variables: options.variables || {},
      fetchPolicy: 'network-only'
    });
    const data = resultQuery?.data?.[options.key] || [];
    targetRef.value = options.transform ? options.transform(data) : data;
  } catch {
    targetRef.value = [];
  }
};

const fetchUnits = () => fetchList(UNITS_QUERY, units, { key: 'units' });

const fetchRecipeImages = async () => {
  await fetchList(RECIPE_IMAGES_QUERY, recipeImages, { key: 'recipe_images', variables: { id: recipeId } });
  const featuredIndex = recipeImages.value.findIndex((img) => img.is_featured);
  currentImageIndex.value = featuredIndex >= 0 ? featuredIndex : 0;
};

const fetchIngredients = () => fetchList(RECIPE_INGREDIENTS_QUERY, ingredientsData, {
  key: 'recipe_ingredients',
  variables: { id: recipeId }
});

const fetchSteps = () => fetchList(RECIPE_STEPS_QUERY, stepsData, {
  key: 'recipe_steps',
  variables: { id: recipeId },
  transform: (items) => [...items].sort((a, b) => a.step_number - b.step_number)
});

const fetchComments = async () => {
  if (!recipeId) return;
  try {
    const resultQuery = await client.query({
      query: RECIPE_COMMENTS_QUERY,
      variables: { id: recipeId },
      fetchPolicy: 'network-only'
    });
    const comments = resultQuery?.data?.comments || [];
    commentsData.value = comments.map((comment) => ({
      id: comment.id,
      user_id: comment.user_id,
      user_name: comment.user?.name || 'User',
      content: comment.content,
      created_at: comment.created_at
    }));
  } catch {
    try {
      const fallback = await client.query({
        query: RECIPE_COMMENTS_QUERY_CAMEL,
        variables: { id: recipeId },
        fetchPolicy: 'network-only'
      });
      const comments = fallback?.data?.comments || [];
      commentsData.value = comments.map((comment) => ({
        id: comment.id,
        user_id: comment.userId,
        user_name: comment.user?.name || 'User',
        content: comment.content,
        created_at: comment.createdAt
      }));
    } catch {
      commentsData.value = [];
    }
  }
};

const fetchRating = async () => {
  try {
    const resultQuery = await client.query({
      query: RECIPE_RATING_QUERY,
      variables: { id: recipeId },
      fetchPolicy: 'network-only'
    });
    const data = resultQuery?.data?.recipes_by_pk;
    ratingData.value = { average_rating: Number(data?.recipe_average_rating || 0) };
    likeCount.value = Number(data?.recipe_likes_count || 0);
  } catch {
    ratingData.value = null;
    likeCount.value = 0;
  }
};

const loadProtectedContent = async () => {
  if (!canViewContent.value) {
    ingredientsData.value = [];
    stepsData.value = [];
    return;
  }
  await Promise.all([fetchIngredients(), fetchSteps()]);
};

const verifyPurchaseFromPayment = async () => {
  if (!fromPayment.value || !token.value || !recipeId || !userId.value) return;
  try {
    const input = paymentTxRef.value
      ? { tx_ref: paymentTxRef.value }
      : { recipe_id: recipeId };

    const data = await verifyPaymentAction(client, input);
    if (data?.status === 'success') {
      hasPurchased.value = true;
    }
  } catch {
  }
};

const checkUserInteractions = async () => {
  if (!token.value || !userId.value || !recipeId) return;
  try {
    const [purchaseResult, socialResult] = await Promise.all([
      client.query({
        query: CHECK_PURCHASES_RATINGS_QUERY,
        variables: { recipeId, userId: userId.value },
        fetchPolicy: 'network-only'
      }),
      client.query({
        query: CHECK_SOCIAL_QUERY,
        variables: { recipeId, userId: userId.value },
        fetchPolicy: 'network-only'
      }).catch(() => ({ data: {} }))
    ]);

    const purchases = purchaseResult?.data?.purchases || [];
    const ratings = purchaseResult?.data?.ratings || [];
    const likes = socialResult?.data?.likes || [];
    const bookmarks = socialResult?.data?.bookmarks || [];

    hasPurchased.value = purchases.length > 0;
    isLiked.value = likes.length > 0;
    isBookmarked.value = bookmarks.length > 0;
    userRating.value = ratings[0]?.rating ? Number.parseInt(ratings[0].rating, 10) : 0;
  } catch {
  }
};

const checkAccessFromDb = async () => {
  if (!token.value || !userId.value || !recipeId) {
    dbHasAccess.value = false;
    return;
  }

  try {
    const resultQuery = await client.query({
      query: GET_ACCESSIBLE_RECIPE_QUERY,
      variables: { userId: userId.value, recipeId },
      fetchPolicy: 'network-only'
    });
    const rows = resultQuery?.data?.get_accessible_recipe || [];
    dbHasAccess.value = rows.length > 0;
  } catch {
    dbHasAccess.value = false;
  }
};

const nextImage = () => {
  if (currentImageIndex.value < recipeImages.value.length - 1) {
    currentImageIndex.value += 1;
  }
};

const previousImage = () => {
  if (currentImageIndex.value > 0) {
    currentImageIndex.value -= 1;
  }
};

const toggleLike = async () => {
  if (!requireAuth()) return;
  actionLoading.value = true;
  try {
    if (isLiked.value) {
      await client.mutate({
        mutation: UNLIKE_RECIPE_MUTATION,
        variables: { recipeId, userId: userId.value }
      });
      isLiked.value = false;
    } else {
      await client.mutate({
        mutation: LIKE_RECIPE_MUTATION,
        variables: { recipeId }
      });
      isLiked.value = true;
    }
    await fetchRating();
  } finally {
    actionLoading.value = false;
  }
};

const toggleBookmark = async () => {
  if (!requireAuth()) return;
  actionLoading.value = true;
  try {
    if (isBookmarked.value) {
      await client.mutate({
        mutation: UNBOOKMARK_RECIPE_MUTATION,
        variables: { recipeId, userId: userId.value }
      });
      isBookmarked.value = false;
    } else {
      await client.mutate({
        mutation: BOOKMARK_RECIPE_MUTATION,
        variables: { recipeId }
      });
      isBookmarked.value = true;
    }
  } finally {
    actionLoading.value = false;
  }
};

const submitRating = async (rating) => {
  if (!requireAuth() || rating < 1 || rating > 5) return;
  ratingLoading.value = true;
  ratingSuccess.value = false;
  try {
    await client.mutate({
      mutation: RATE_RECIPE_MUTATION,
      variables: { recipeId, rating }
    });
    userRating.value = Number(rating);
    ratingSuccess.value = true;
    await fetchRating();
    setTimeout(() => {
      ratingSuccess.value = false;
    }, 3000);
  } finally {
    ratingLoading.value = false;
  }
};

const isCommentOwner = (comment) => Number(comment?.user_id) === Number(userId.value);

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
  if (!requireAuth() || !isCommentOwner(comment) || !editingCommentText.value.trim()) return;
  commentLoading.value = true;
  try {
    await client.mutate({
      mutation: UPDATE_COMMENT_MUTATION,
      variables: { id: comment.id, content: editingCommentText.value.trim() }
    });
    cancelEditComment();
    await fetchComments();
  } finally {
    commentLoading.value = false;
  }
};

const deleteComment = async (comment) => {
  if (!requireAuth() || !isCommentOwner(comment)) return;
  if (!confirm('Delete this comment?')) return;
  deletingCommentId.value = comment.id;
  try {
    await client.mutate({
      mutation: DELETE_COMMENT_MUTATION,
      variables: { id: comment.id }
    });
    if (editingCommentId.value === comment.id) {
      cancelEditComment();
    }
    await fetchComments();
  } finally {
    deletingCommentId.value = null;
  }
};

const submitComment = async () => {
  if (!requireAuth() || !newComment.value.trim()) return;
  commentLoading.value = true;
  const content = newComment.value.trim();
  try {
    try {
      await client.mutate({
        mutation: POST_COMMENT_MUTATION,
        variables: { recipeId, content }
      });
    } catch {
      await client.mutate({
        mutation: POST_COMMENT_MUTATION_CAMEL,
        variables: { recipeId, userId: userId.value, content }
      });
    }
    newComment.value = '';
    await fetchComments();
  } finally {
    commentLoading.value = false;
  }
};

const handleBuyRecipe = async () => {
  if (!requireAuth() || !recipe.value) return;
  buying.value = true;
  try {
    const normalizedUserName = (userName.value || '').trim() || 'Customer';
    const data = await initializePaymentAction(client, {
      amount: String(recipe.value.price),
      email: userEmail.value || 'testuser@gmail.com',
      user_name: normalizedUserName,
      recipe_id: recipeId
    });

    if (data?.status === 'success' && !data?.checkout_url) {
      hasPurchased.value = true;
      await loadProtectedContent();
      const q = new URLSearchParams({
        recipe_id: String(recipeId),
        status: 'success',
        message: String(data?.message || 'Payment already confirmed')
      });
      if (data?.tx_ref) {
        q.set('tx_ref', String(data.tx_ref));
      }
      await navigateTo(`/payment/success?${q.toString()}`);
      return;
    }
    if (process.client) {
      sessionStorage.setItem('pending_payment_recipe_id', String(recipeId));
    }
    if (data?.checkout_url) {
      window.location.href = data.checkout_url;
    }
  } catch (error) {
    console.error('Failed to initialize payment:', error);
    alert('Failed to initialize payment. Please try again.');
  } finally {
    buying.value = false;
  }
};

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

  editIngredients.value = (ingredientsData.value || []).map((item) => ({
    name: item.name || '',
    quantity: item.quantity || '',
    unit_id: String(item.unit_id || item.unit?.id || '')
  }));
  if (editIngredients.value.length === 0) {
    editIngredients.value = [createIngredient()];
  }

  editSteps.value = (stepsData.value || []).map((step) => ({ instruction: step.instruction || '' }));
  if (editSteps.value.length === 0) {
    editSteps.value = [createStep()];
  }

  inlineEditError.value = '';
  isInlineEditing.value = true;
};

const cancelInlineEdit = () => {
  isInlineEditing.value = false;
  inlineEditError.value = '';
};

const addEditIngredient = () => {
  editIngredients.value.push(createIngredient());
};

const removeEditIngredient = (index) => {
  if (editIngredients.value.length > 1) {
    editIngredients.value.splice(index, 1);
  }
};

const addEditStep = () => {
  editSteps.value.push(createStep());
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

  const validIngredients = editIngredients.value.filter((item) => item.name?.trim());
  if (validIngredients.length === 0 || validIngredients.some((item) => !item.unit_id)) {
    inlineEditError.value = 'Please add valid ingredients and select units.';
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
          preparation_time: Number.parseInt(editForm.value.preparation_time, 10),
          price: Number.parseFloat(editForm.value.price) || 0,
          category_id: Number.parseInt(editForm.value.category_id, 10)
        },
        ingredients: validIngredients.map((item) => ({
          recipe_id: recipeId,
          name: item.name.trim(),
          quantity: item.quantity || '',
          unit_id: Number.parseInt(item.unit_id, 10)
        })),
        steps: validSteps.map((step, index) => ({
          recipe_id: recipeId,
          step_number: index + 1,
          instruction: step.instruction.trim()
        }))
      }
    });

    await refetchRecipe({ id: recipeId });
    await Promise.all([fetchRating(), fetchIngredients(), fetchSteps()]);
    isInlineEditing.value = false;
  } catch (err) {
    inlineEditError.value = err?.message || 'Failed to update recipe.';
  } finally {
    inlineEditLoading.value = false;
  }
};

const getRecipeImage = () => {
  if (recipeImages.value.length > 0) {
    const featured = recipeImages.value.find((img) => img.is_featured);
    return featured?.url || recipeImages.value[0]?.url;
  }
  return 'https://images.unsplash.com/photo-1495521821757-a1efb6729352?auto=format&fit=crop&w=1200&q=80';
};

const formatDate = (dateString) => {
  if (!dateString) return '';
  return new Date(dateString).toLocaleDateString();
};

watch(
  () => recipe.value,
  async (loadedRecipe) => {
    if (!loadedRecipe) return;
    await Promise.all([fetchRecipeImages(), fetchComments(), fetchRating()]);
    await verifyPurchaseFromPayment();
    await checkUserInteractions();
    await checkAccessFromDb();
    await loadProtectedContent();

    if (!hasAutoOpenedInlineEdit.value && isOwner.value && String(route.query.edit || '') === '1') {
      hasAutoOpenedInlineEdit.value = true;
      await startInlineEdit();
    }
  },
  { immediate: true }
);

watch(
  () => recipeCategory.value?.image_url,
  () => {
    // Reset failure flag when category/image changes.
    recipeCategoryImageFailed.value = false;
  }
);

watch(
  () => canViewContent.value,
  async (allowed) => {
    if (allowed) {
      await loadProtectedContent();
    }
  }
);
</script>
 expalin what is the main role of this code