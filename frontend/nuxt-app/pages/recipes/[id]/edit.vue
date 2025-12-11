<template>
  <div class="relative min-h-screen pb-20">
    <!-- Background -->
    <div class="fixed inset-0 z-0">
      <img 
        src="https://images.unsplash.com/photo-1556910103-1c02745aae4d?ixlib=rb-4.0.3&auto=format&fit=crop&w=2070&q=80" 
        alt="Kitchen" 
        class="w-full h-full object-cover brightness-95"
      >
      <div class="absolute inset-0 bg-gradient-to-b from-black/40 via-black/10 to-black/60"></div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="relative z-10 flex justify-center items-center min-h-screen">
      <div class="text-center">
        <div class="animate-spin rounded-full h-16 w-16 border-t-4 border-b-4 border-emerald-500 mx-auto"></div>
        <p class="mt-4 text-white">Loading recipe...</p>
      </div>
    </div>

    <!-- Error State -->
    <div v-else-if="!loading && !recipe" class="relative z-10 max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 pt-16 pb-16">
      <div class="bg-red-500/20 border border-red-500/50 text-red-200 px-6 py-4 rounded-lg">
        <strong class="font-bold">Error!</strong>
        <span class="block sm:inline"> {{ updateError || 'Recipe not found' }}</span>
        <button @click="$router.push('/profile')" class="mt-4 px-4 py-2 bg-emerald-500 text-white rounded-lg">Go to Profile</button>
      </div>
    </div>

    <!-- Content -->
    <div v-else-if="recipe" class="relative z-10 max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 pt-16 pb-16">
      <div class="bg-white/10 backdrop-blur-lg border border-white/20 rounded-2xl shadow-2xl p-8">
        <h1 class="text-3xl font-bold text-white mb-6">Edit Recipe: {{ recipe.title || 'Loading...' }}</h1>
        
        <!-- Debug Info (remove this after testing) -->
        <div v-if="recipe" class="mb-4 p-4 bg-black/20 rounded-lg text-xs text-gray-300">
          <p>Recipe ID: {{ recipeId }}</p>
          <p>Title: {{ formTitle || '(empty)' }}</p>
          <p>Category ID: {{ formCategoryId }}</p>
          <p>Has Recipe: {{ !!recipe }}</p>
        </div>
        
        <form @submit.prevent="handleFormSubmit">
          <!-- Basic Information -->
          <div class="mb-8">
            <h2 class="text-2xl font-bold text-white mb-6">Basic Information</h2>
            
            <div class="mb-6">
              <label class="block text-sm font-medium text-gray-300 mb-2">Recipe Title</label>
              <input 
                :value="formTitle"
                @input="(e) => { formTitle.value = e.target.value; }"
                type="text" 
                required
                minlength="3"
                key="title-input"
                class="block w-full px-4 py-3 border border-white/20 rounded-lg bg-black/20 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-emerald-500"
                placeholder="Recipe title"
              />
              <p v-if="formErrors.title" class="text-red-400 text-xs mt-1">{{ formErrors.title }}</p>
            </div>

            <div class="mb-6">
              <label class="block text-sm font-medium text-gray-300 mb-2">Description</label>
              <textarea 
                :value="formDescription"
                @input="(e) => { formDescription.value = e.target.value; }"
                rows="4"
                required
                minlength="10"
                key="description-input"
                class="block w-full px-4 py-3 border border-white/20 rounded-lg bg-black/20 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-emerald-500"
                placeholder="Recipe description"
              ></textarea>
              <p v-if="formErrors.description" class="text-red-400 text-xs mt-1">{{ formErrors.description }}</p>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label class="block text-sm font-medium text-gray-300 mb-3">Category</label>
                
                <!-- Visual Category Selector -->
                <div class="grid grid-cols-2 md:grid-cols-3 gap-3">
                  <button
                    v-for="cat in categories"
                    :key="cat.id"
                    type="button"
                    @click="formCategoryId.value = cat.id"
                    :class="[
                      'group relative overflow-hidden rounded-lg border-2 transition-all duration-300 transform hover:scale-105',
                      formCategoryId.value === cat.id
                        ? 'border-emerald-400 ring-2 ring-emerald-400/50 shadow-lg shadow-emerald-500/30'
                        : 'border-white/20 hover:border-white/40'
                    ]"
                  >
                    <!-- Category Image -->
                    <div class="relative h-20 md:h-24 overflow-hidden">
                      <img 
                        v-if="cat.image_url"
                        :src="cat.image_url" 
                        :alt="cat.name"
                        class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-110"
                        @error="(e) => e.target.style.display = 'none'"
                      />
                      <!-- Fallback gradient if no image -->
                      <div 
                        v-else
                        class="w-full h-full bg-gradient-to-br from-emerald-500/40 via-teal-500/40 to-blue-500/40 flex items-center justify-center"
                      >
                        <span class="text-3xl">{{ getCategoryEmoji(cat.name) }}</span>
                      </div>
                      
                      <!-- Overlay gradient -->
                      <div class="absolute inset-0 bg-gradient-to-t from-black/70 via-black/30 to-transparent"></div>
                      
                      <!-- Category Name -->
                      <div class="absolute bottom-0 left-0 right-0 p-2 text-center">
                        <p class="text-white font-semibold text-xs drop-shadow-lg">{{ cat.name }}</p>
                      </div>
                      
                      <!-- Selected Checkmark -->
                      <div v-if="formCategoryId === cat.id" class="absolute top-1 right-1">
                        <div class="bg-emerald-500 rounded-full p-1 shadow-lg">
                          <svg class="w-3 h-3 text-white" fill="currentColor" viewBox="0 0 20 20">
                            <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                          </svg>
                        </div>
                      </div>
                    </div>
                  </button>
                </div>
                
                <p v-if="formErrors.category_id" class="text-red-400 text-xs mt-1">{{ formErrors.category_id }}</p>
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-300 mb-2">Preparation Time (minutes)</label>
                <input 
                  :value="formPreparationTime"
                  @input="(e) => { formPreparationTime.value = parseInt(e.target.value) || 0; }"
                  type="number"
                  required
                  min="1"
                  key="prep-time-input"
                  class="block w-full px-4 py-3 border border-white/20 rounded-lg bg-black/20 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-emerald-500"
                  placeholder="e.g., 45"
                />
                <p v-if="formErrors.preparation_time" class="text-red-400 text-xs mt-1">{{ formErrors.preparation_time }}</p>
              </div>
            </div>

            <div class="mt-6">
              <label class="block text-sm font-medium text-gray-300 mb-2">Price (0 for free)</label>
              <input 
                :value="formPrice"
                @input="(e) => { formPrice.value = parseFloat(e.target.value) || 0; }"
                type="number"
                step="0.01"
                min="0"
                key="price-input"
                class="block w-full px-4 py-3 border border-white/20 rounded-lg bg-black/20 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-emerald-500"
                placeholder="0.00"
              />
              <p v-if="formErrors.price" class="text-red-400 text-xs mt-1">{{ formErrors.price }}</p>
            </div>
          </div>

          <!-- Multiple Images with Featured Selection -->
          <div class="mb-8">
            <h2 class="text-2xl font-bold text-white mb-6">Recipe Images</h2>
            <p class="text-gray-300 text-sm mb-4">Upload multiple images and select one as featured (thumbnail)</p>
            
            <!-- Image Upload Area -->
            <div class="border-2 border-dashed border-white/30 rounded-lg p-6 text-center hover:border-emerald-400/50 transition-colors mb-4">
              <input 
                type="file" 
                @change="handleMultipleImageUpload" 
                accept="image/*"
                multiple
                class="hidden" 
                ref="imageInput"
              />
              
              <div @click="$refs.imageInput.click()" class="cursor-pointer">
                <svg class="mx-auto h-12 w-12 text-gray-400" stroke="currentColor" fill="none" viewBox="0 0 48 48">
                  <path d="M28 8H12a4 4 0 00-4 4v20m32-12v8m0 0v8a4 4 0 01-4 4H12a4 4 0 01-4-4v-4m32-4l-3.172-3.172a4 4 0 00-5.656 0L28 28M8 32l9.172-9.172a4 4 0 015.656 0L28 28m0 0l4 4m4-24h8m-4-4v8m-12 4h.02" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
                </svg>
                <p class="mt-2 text-sm text-gray-300">Click to upload images</p>
                <p class="text-xs text-gray-400 mt-1">PNG, JPG, GIF up to 10MB each (multiple files allowed)</p>
              </div>
            </div>
            <p v-if="uploadingImage" class="text-emerald-400 text-sm mb-4">Uploading images...</p>

            <!-- Uploaded Images Grid -->
            <div v-if="uploadedImages.length > 0" class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
              <div 
                v-for="(img, index) in uploadedImages" 
                :key="index"
                class="relative group"
              >
                <div class="relative">
                  <img :src="img.url" :alt="`Image ${index + 1}`" class="w-full h-32 object-cover rounded-lg border-2" :class="img.isFeatured ? 'border-emerald-500' : 'border-white/20'" />
                  
                  <!-- Featured Badge -->
                  <div v-if="img.isFeatured" class="absolute top-2 left-2 bg-emerald-500 text-white px-2 py-1 rounded text-xs font-bold">
                    ⭐ Featured
                  </div>
                  
                  <!-- Actions -->
                  <div class="absolute inset-0 bg-black/60 opacity-0 group-hover:opacity-100 transition-opacity rounded-lg flex items-center justify-center gap-2">
                    <button
                      v-if="!img.isFeatured"
                      @click="setFeaturedImage(index)"
                      class="px-3 py-1 bg-emerald-500 text-white rounded text-xs hover:bg-emerald-600"
                      title="Set as featured"
                    >
                      Set Featured
                    </button>
                    <button
                      @click="removeImage(index)"
                      class="px-3 py-1 bg-red-500 text-white rounded text-xs hover:bg-red-600"
                      title="Remove image"
                    >
                      Remove
                    </button>
                  </div>
                </div>
              </div>
            </div>
            
            <p v-if="uploadedImages.length === 0 && !uploadingImage" class="text-gray-400 text-sm text-center py-4">
              No images uploaded yet. Upload at least one image and select it as featured.
            </p>
          </div>

          <!-- Ingredients -->
          <div class="mb-8">
            <h2 class="text-2xl font-bold text-white mb-6">Ingredients</h2>
            <div v-for="(ing, index) in ingredients" :key="index" class="mb-4 flex gap-3">
              <input v-model="ing.name" placeholder="Name" class="flex-1 px-4 py-3 border border-white/20 rounded-lg bg-black/20 text-white" />
              <input v-model="ing.quantity" placeholder="Qty" class="w-24 px-4 py-3 border border-white/20 rounded-lg bg-black/20 text-white" />
              <input v-model="ing.unit" placeholder="Unit" class="w-32 px-4 py-3 border border-white/20 rounded-lg bg-black/20 text-white" />
              <button type="button" @click="removeIngredient(index)" class="px-4 py-3 bg-red-500/20 text-red-400 rounded-lg">Remove</button>
            </div>
            <button type="button" @click="addIngredient" class="w-full px-4 py-3 border-2 border-dashed border-emerald-400/50 rounded-lg text-emerald-400">+ Add Ingredient</button>
          </div>

          <!-- Steps -->
          <div class="mb-8">
            <h2 class="text-2xl font-bold text-white mb-6">Steps</h2>
            <div v-for="(step, index) in steps" :key="index" class="mb-4">
              <textarea v-model="step.instruction" rows="3" placeholder="Step instruction" class="w-full px-4 py-3 border border-white/20 rounded-lg bg-black/20 text-white"></textarea>
              <button type="button" @click="removeStep(index)" class="mt-2 px-4 py-2 bg-red-500/20 text-red-400 rounded-lg">Remove</button>
            </div>
            <button type="button" @click="addStep" class="w-full px-4 py-3 border-2 border-dashed border-emerald-400/50 rounded-lg text-emerald-400">+ Add Step</button>
          </div>

          <div v-if="updateError" class="mb-6 p-4 rounded-lg bg-red-500/20 border border-red-500/50 text-red-200">
            {{ updateError }}
          </div>

          <div class="flex gap-4">
            <button type="button" @click="$router.push(`/recipes/${recipeId}`)" class="flex-1 px-6 py-4 border border-white/30 rounded-lg text-white">Cancel</button>
            <button type="submit" :disabled="isSubmitting" class="flex-1 px-6 py-4 bg-gradient-to-r from-emerald-600 to-teal-600 text-white font-bold rounded-lg">
              {{ isSubmitting ? 'Updating...' : 'Update Recipe' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';

const route = useRoute();
const router = useRouter();
const recipeId = parseInt(route.params.id);
const token = useCookie('auth_token');
const loading = ref(true);
const recipe = ref(null);
const categories = ref([]);
const ingredients = ref([]);
const steps = ref([]);
const updateError = ref('');
const imageInput = ref(null);
const uploadedImages = ref([]);
const uploadingImage = ref(false);
const isSubmitting = ref(false);
const formErrors = ref({});

// Form values - use individual refs for maximum compatibility
const formTitle = ref('');
const formDescription = ref('');
const formCategoryId = ref(0);
const formPreparationTime = ref(0);
const formPrice = ref(0);

// Redirect if not logged in
onMounted(async () => {
  console.log('[EDIT] Page mounted, recipeId:', recipeId);
  
  if (!token.value) {
    console.log('[EDIT] No token, redirecting to login');
    router.push('/login');
    return;
  }
  
  try {
    console.log('[EDIT] Fetching recipe from:', `http://localhost:8081/recipes/${recipeId}`);
    // Fetch recipe
    const recipeRes = await fetch(`http://localhost:8081/recipes/${recipeId}`);
    console.log('[EDIT] Recipe response status:', recipeRes.status);
    
    if (!recipeRes.ok) {
      const errorText = await recipeRes.text();
      console.error('[EDIT] Recipe fetch failed:', errorText);
      throw new Error(`Recipe not found: ${errorText}`);
    }
    
    recipe.value = await recipeRes.json();
    console.log('[EDIT] Recipe loaded:', recipe.value);
    
    // Set form values - assign to individual refs
    if (recipe.value) {
      // Wait for next tick to ensure DOM is ready
      await nextTick();
      
      // Assign to individual refs
      formTitle.value = recipe.value.title || '';
      formDescription.value = recipe.value.description || '';
      formCategoryId.value = recipe.value.category_id || 0;
      formPreparationTime.value = recipe.value.preparation_time || 0;
      formPrice.value = recipe.value.price || 0;
      
      console.log('[EDIT] Form values set:');
      console.log('[EDIT] Title:', formTitle.value);
      console.log('[EDIT] Description:', formDescription.value);
    } else {
      console.error('[EDIT] Recipe value is null/undefined');
    }

    // Fetch ingredients
    const ingRes = await fetch(`http://localhost:8081/recipes/${recipeId}/ingredients`);
    if (ingRes.ok) {
      const ingData = await ingRes.json();
      // Ensure ingredients have the right structure
      ingredients.value = ingData.map(ing => ({
        name: ing.name || '',
        quantity: ing.quantity || '',
        unit: ing.unit || ''
      }));
      // If no ingredients, add one empty row
      if (ingredients.value.length === 0) {
        ingredients.value = [{ name: '', quantity: '', unit: '' }];
      }
      console.log('[EDIT] Ingredients loaded:', ingredients.value);
    }

    // Fetch steps
    const stepsRes = await fetch(`http://localhost:8081/recipes/${recipeId}/steps`);
    if (stepsRes.ok) {
      const stepsData = await stepsRes.json();
      // Ensure steps have the right structure
      steps.value = stepsData.map(step => ({
        instruction: step.instruction || '',
        image_url: step.image_url || ''
      }));
      // If no steps, add one empty row
      if (steps.value.length === 0) {
        steps.value = [{ instruction: '', image_url: '' }];
      }
      console.log('[EDIT] Steps loaded:', steps.value);
    }

    // Fetch existing images
    const imagesRes = await fetch(`http://localhost:8081/recipes/${recipeId}/images`);
    if (imagesRes.ok) {
      const images = await imagesRes.json();
      uploadedImages.value = images.map(img => ({ url: img.url, isFeatured: img.is_featured }));
    }

    // Fetch categories
    const catRes = await fetch('http://localhost:8081/categories');
    if (catRes.ok) {
      categories.value = await catRes.json();
      console.log('[EDIT] Categories loaded:', categories.value.length);
    }
  } catch (err) {
    console.error('[EDIT] Error loading data:', err);
    updateError.value = err.message || 'Failed to load recipe data';
    alert(`Error: ${err.message}`);
  } finally {
    loading.value = false;
    console.log('[EDIT] Loading complete');
  }
});

const addIngredient = () => ingredients.value.push({ name: '', quantity: '', unit: '' });
const removeIngredient = (index) => {
  if (ingredients.value.length > 1) ingredients.value.splice(index, 1);
};
const addStep = () => steps.value.push({ instruction: '', image_url: '' });
const removeStep = (index) => {
  if (steps.value.length > 1) steps.value.splice(index, 1);
};

// Get Category Emoji (fallback when image is not available)
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

// Handle Multiple Image Upload
const handleMultipleImageUpload = async (event) => {
  const files = Array.from(event.target.files);
  if (files.length === 0) return;

  uploadingImage.value = true;
  updateError.value = '';

  try {
    const uploadPromises = files.map(async (file) => {
      const formData = new FormData();
      formData.append('file', file);
      
      const response = await fetch('http://localhost:8081/upload', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token.value}`
        },
        body: formData
      });

      if (!response.ok) {
        throw new Error(`Failed to upload ${file.name}`);
      }

      const data = await response.json();
      return { url: data.url, isFeatured: false };
    });

    const uploaded = await Promise.all(uploadPromises);
    uploadedImages.value.push(...uploaded);
    
    // If this is the first image, set it as featured automatically
    if (uploadedImages.value.length === uploaded.length && uploadedImages.value.length > 0) {
      uploadedImages.value[0].isFeatured = true;
    }
  } catch (err) {
    updateError.value = err.message || 'Failed to upload images';
    console.error(err);
  } finally {
    uploadingImage.value = false;
    if (event.target) {
      event.target.value = '';
    }
  }
};

// Set Featured Image
const setFeaturedImage = (index) => {
  uploadedImages.value.forEach(img => img.isFeatured = false);
  uploadedImages.value[index].isFeatured = true;
};

// Remove Image
const removeImage = (index) => {
  const removed = uploadedImages.value.splice(index, 1)[0];
  if (removed.isFeatured && uploadedImages.value.length > 0) {
    uploadedImages.value[0].isFeatured = true;
  }
};

// Form validation
const validateForm = () => {
  formErrors.value = {};
  let isValid = true;
  
  if (!formTitle.value || formTitle.value.trim().length < 3) {
    formErrors.value.title = 'Title must be at least 3 characters';
    isValid = false;
  }
  
  if (!formDescription.value || formDescription.value.trim().length < 10) {
    formErrors.value.description = 'Description must be at least 10 characters';
    isValid = false;
  }
  
  if (!formCategoryId.value || formCategoryId.value === 0) {
    formErrors.value.category_id = 'Please select a category';
    isValid = false;
  }
  
  if (!formPreparationTime.value || formPreparationTime.value < 1) {
    formErrors.value.preparation_time = 'Preparation time must be at least 1 minute';
    isValid = false;
  }
  
  if (formPrice.value < 0) {
    formErrors.value.price = 'Price cannot be negative';
    isValid = false;
  }
  
  return isValid;
};

// Handle form submit
const handleFormSubmit = async () => {
  updateError.value = '';
  formErrors.value = {};
  
  if (!validateForm()) {
    return;
  }
  
  // Create form data object from individual refs
  const formData = {
    title: formTitle.value,
    description: formDescription.value,
    category_id: formCategoryId.value,
    preparation_time: formPreparationTime.value,
    price: formPrice.value
  };
  
  await handleUpdateRecipe(formData);
};

const handleUpdateRecipe = async (values) => {
  updateError.value = '';
  isSubmitting.value = true;
  
  console.log('[EDIT] Form submitted with values:', values);
  
  // Validate ingredients
  const validIngredients = ingredients.value.filter(ing => ing.name && ing.name.trim() !== '');
  if (validIngredients.length === 0) {
    updateError.value = 'Please add at least one ingredient';
    return;
  }

  // Validate steps
  const validSteps = steps.value.filter(step => step.instruction && step.instruction.trim() !== '');
  if (validSteps.length === 0) {
    updateError.value = 'Please add at least one preparation step';
    return;
  }

  // Validate images (allow keeping existing images if no new ones uploaded)
  if (uploadedImages.value.length === 0) {
    // If editing and no images uploaded, use existing thumbnail
    if (recipe.value?.thumbnail_url) {
      uploadedImages.value = [{ url: recipe.value.thumbnail_url, isFeatured: true }];
    } else {
      updateError.value = 'Please upload at least one image';
      return;
    }
  }

  // Find featured image
  const featuredImage = uploadedImages.value.find(img => img.isFeatured);
  if (!featuredImage) {
    // Auto-set first image as featured if none selected
    uploadedImages.value[0].isFeatured = true;
  }
  
  // Format ingredients to match backend model (name, quantity, unit)
  const formattedIngredients = validIngredients.map(ing => ({
    name: ing.name || '',
    quantity: ing.quantity || '',
    unit: ing.unit || ''
  }));
  
  // Format steps to match backend model (instruction, image_url)
  const formattedSteps = validSteps.map(step => ({
    instruction: step.instruction || '',
    image_url: step.image_url || ''
  }));
  
  // Use values from form submission
  const recipeData = {
    category_id: parseInt(values.category_id),
    title: values.title,
    description: values.description,
    preparation_time: parseInt(values.preparation_time),
    price: parseFloat(values.price) || 0,
    thumbnail_url: featuredImage.url,
    ingredients: formattedIngredients,
    steps: formattedSteps,
    images: uploadedImages.value.map(img => img.url)
  };
  
  console.log('[EDIT] Formatted recipe data to send:', recipeData);

  try {
    console.log('[EDIT] Updating recipe:', recipeId, recipeData);
    const response = await fetch(`http://localhost:8081/recipes/${recipeId}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token.value}`
      },
      body: JSON.stringify(recipeData)
    });

    if (!response.ok) {
      let errorData = '';
      try {
        errorData = await response.text();
        const errorJson = JSON.parse(errorData);
        throw new Error(errorJson.error || errorData || 'Failed to update recipe');
      } catch (parseError) {
        throw new Error(errorData || 'Failed to update recipe');
      }
    }

    // Try to parse response as JSON, but handle empty responses
    let data = null;
    const responseText = await response.text();
    if (responseText) {
      try {
        data = JSON.parse(responseText);
      } catch (e) {
        // Response is not JSON, that's okay
        console.log('[EDIT] Response is not JSON, assuming success');
      }
    }
    
    console.log('[EDIT] Recipe updated successfully:', data || 'No response data');
    alert('Recipe updated successfully!');
    router.push(`/recipes/${recipeId}`);
  } catch (err) {
    updateError.value = err.message || 'An error occurred while updating the recipe';
    console.error('[EDIT] Exception:', err);
  } finally {
    isSubmitting.value = false;
  }
};
</script>

