<template>
  <div class="relative min-h-screen pb-20">
    <!-- Background Image with Overlay - Inspiring Kitchen Scene (Create Page) -->
    <div class="fixed inset-0 z-0">
      <img 
        src="https://images.unsplash.com/photo-1556910103-1c02745aae4d?ixlib=rb-4.0.3&auto=format&fit=crop&w=2070&q=80" 
        alt="Elegant Kitchen" 
        class="w-full h-full object-cover brightness-95"
      >
      <div class="absolute inset-0 bg-black/80"></div>
    </div>

    <!-- Content -->
    <div class="relative z-10 max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 pt-16 pb-16">
      <!-- Header -->
      <div class="text-center mb-10 pt-8">
        <div class="bg-white/10 backdrop-blur-lg border border-white/20 rounded-3xl shadow-2xl px-8 py-10 mb-8">
          <h1 class="text-4xl md:text-5xl lg:text-6xl font-serif font-extrabold text-white mb-4 tracking-tight drop-shadow-lg">
            Share Your <span class="text-emerald-400">Culinary</span> Creation
          </h1>
          <p class="text-lg md:text-xl text-gray-200 max-w-2xl mx-auto font-light">
            Create a new recipe and inspire food lovers around the world
          </p>
        </div>
      </div>

      <!-- Form Card -->
      <div class="bg-white/10 backdrop-blur-lg border border-white/20 rounded-2xl shadow-2xl p-8 hover:bg-white/15 transition-all duration-300">
        <Form @submit="handleCreateRecipe" :validation-schema="schema" v-slot="{ errors, isSubmitting }">
          
          <!-- Basic Information -->
          <div class="mb-8">
            <h2 class="text-2xl font-bold text-white mb-6 flex items-center">
              <span class="text-emerald-400 mr-2">📝</span> Basic Information
            </h2>
            
            <!-- Title -->
            <div class="mb-6">
              <label for="create-recipe-title" class="block text-sm font-medium text-gray-300 mb-2">Recipe Title</label>
              <Field 
                id="create-recipe-title"
                name="title" 
                type="text" 
                class="block w-full px-4 py-3 border border-white/20 rounded-lg bg-black/20 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500"
                placeholder="e.g., Grandma's Chocolate Cake"
                :class="{ 'border-red-500 focus:ring-red-500': errors.title }"
              />
              <ErrorMessage name="title" class="text-red-400 text-xs mt-1" />
            </div>

            <!-- Description -->
            <div class="mb-6">
              <label for="create-recipe-description" class="block text-sm font-medium text-gray-300 mb-2">Description</label>
              <Field 
                id="create-recipe-description"
                name="description" 
                as="textarea"
                rows="4"
                class="block w-full px-4 py-3 border border-white/20 rounded-lg bg-black/20 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500"
                placeholder="Tell us about your recipe..."
                :class="{ 'border-red-500 focus:ring-red-500': errors.description }"
              />
              <ErrorMessage name="description" class="text-red-400 text-xs mt-1" />
            </div>

            <!-- Category & Time -->
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label class="block text-sm font-medium text-gray-300 mb-3">Category</label>
                
                <!-- Visual Category Selector -->
                <div class="grid grid-cols-2 md:grid-cols-3 gap-3 mb-3">
                  <button
                    v-for="cat in categories"
                    :key="cat.id"
                    type="button"
                    @click="selectedCategoryId = cat.id"
                    :class="[
                      'group relative overflow-hidden rounded-lg border-2 transition-all duration-300 transform hover:scale-105',
                      selectedCategoryId === cat.id
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
                      <div v-if="selectedCategoryId === cat.id" class="absolute top-1 right-1">
                        <div class="bg-emerald-500 rounded-full p-1 shadow-lg">
                          <svg class="w-3 h-3 text-white" fill="currentColor" viewBox="0 0 20 20">
                            <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                          </svg>
                        </div>
                      </div>
                    </div>
                  </button>
                </div>
                
                <!-- Hidden Field for Form Validation - sync with selectedCategoryId -->
                <Field 
                  name="category_id" 
                  :model-value="selectedCategoryId"
                  @update:model-value="selectedCategoryId = $event"
                  type="hidden"
                />
                <ErrorMessage name="category_id" class="text-red-400 text-xs mt-1" />
                <p v-if="!selectedCategoryId" class="text-yellow-400 text-xs mt-1">⚠️ Please select a category</p>
              </div>

              <div>
                <label for="create-recipe-preparation-time" class="block text-sm font-medium text-gray-300 mb-2">Preparation Time (minutes)</label>
                <Field 
                  id="create-recipe-preparation-time"
                  name="preparation_time" 
                  type="number"
                  class="block w-full px-4 py-3 border border-white/20 rounded-lg bg-black/20 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500"
                  placeholder="e.g., 45"
                  :class="{ 'border-red-500 focus:ring-red-500': errors.preparation_time }"
                />
                <ErrorMessage name="preparation_time" class="text-red-400 text-xs mt-1" />
              </div>
            </div>

            <!-- Price -->
            <div class="mt-6">
              <label for="create-recipe-price" class="block text-sm font-medium text-gray-300 mb-2">Price (0 for free)</label>
              <Field 
                id="create-recipe-price"
                name="price" 
                type="number"
                step="0.01"
                class="block w-full px-4 py-3 border border-white/20 rounded-lg bg-black/20 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500"
                placeholder="0.00"
                :class="{ 'border-red-500 focus:ring-red-500': errors.price }"
              />
              <ErrorMessage name="price" class="text-red-400 text-xs mt-1" />
            </div>
          </div>

          <!-- Multiple Images with Featured Selection -->
          <div class="mb-8">
            <h2 class="text-2xl font-bold text-white mb-6 flex items-center">
              <span class="text-emerald-400 mr-2">📸</span> Recipe Images
            </h2>
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
            <h2 class="text-2xl font-bold text-white mb-6 flex items-center">
              <span class="text-emerald-400 mr-2">🥕</span> Ingredients
            </h2>
            
            <div v-for="(ingredient, index) in ingredients" :key="index" class="mb-4 flex gap-3">
              <input 
                v-model="ingredient.name"
                type="text"
                placeholder="Ingredient name"
                class="flex-1 px-4 py-3 border border-white/20 rounded-lg bg-black/20 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-emerald-500"
              />
              <input 
                v-model="ingredient.quantity"
                type="text"
                placeholder="Qty"
                class="w-24 px-4 py-3 border border-white/20 rounded-lg bg-black/20 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-emerald-500"
              />
              <input 
                v-model="ingredient.unit"
                type="text"
                placeholder="Unit"
                class="w-32 px-4 py-3 border border-white/20 rounded-lg bg-black/20 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-emerald-500"
              />
              <button 
                type="button"
                @click="removeIngredient(index)"
                class="px-4 py-3 bg-red-500/20 text-red-400 border border-red-500/50 rounded-lg hover:bg-red-500/30 transition-colors"
              >
                Remove
              </button>
            </div>

            <button 
              type="button"
              @click="addIngredient"
              class="w-full px-4 py-3 border-2 border-dashed border-emerald-400/50 rounded-lg text-emerald-400 hover:bg-emerald-500/10 transition-colors font-medium"
            >
              + Add Ingredient
            </button>
          </div>

          <!-- Steps -->
          <div class="mb-8">
            <h2 class="text-2xl font-bold text-white mb-6 flex items-center">
              <span class="text-emerald-400 mr-2">📋</span> Preparation Steps
            </h2>
            
            <div v-for="(step, index) in steps" :key="index" class="mb-4">
              <div class="flex items-start gap-3">
                <div class="flex-shrink-0 w-8 h-8 bg-emerald-500 text-white rounded-full flex items-center justify-center font-bold text-sm">
                  {{ index + 1 }}
                </div>
                <textarea 
                  v-model="step.instruction"
                  rows="3"
                  placeholder="Describe this step..."
                  class="flex-1 px-4 py-3 border border-white/20 rounded-lg bg-black/20 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-emerald-500"
                ></textarea>
                <button 
                  type="button"
                  @click="removeStep(index)"
                  class="px-4 py-3 bg-red-500/20 text-red-400 border border-red-500/50 rounded-lg hover:bg-red-500/30 transition-colors"
                >
                  Remove
                </button>
              </div>
            </div>

            <button 
              type="button"
              @click="addStep"
              class="w-full px-4 py-3 border-2 border-dashed border-emerald-400/50 rounded-lg text-emerald-400 hover:bg-emerald-500/10 transition-colors font-medium"
            >
              + Add Step
            </button>
          </div>

          <!-- Error Message -->
          <div v-if="createError" class="mb-6 p-4 rounded-lg bg-red-500/20 border border-red-500/50 text-red-200 text-sm">
            {{ createError }}
          </div>

          <!-- Submit Button -->
          <div class="flex gap-4">
            <NuxtLink 
              to="/home"
              class="flex-1 px-6 py-4 border border-white/30 rounded-lg text-white hover:bg-white/10 transition-colors font-semibold text-center"
            >
              Cancel
            </NuxtLink>
            <button 
              type="submit" 
              :disabled="isSubmitting || uploadingImage"
              class="flex-1 px-6 py-4 bg-gradient-to-r from-emerald-600 to-teal-600 text-white font-bold rounded-lg hover:from-emerald-500 hover:to-teal-500 focus:outline-none focus:ring-2 focus:ring-emerald-500 disabled:opacity-50 disabled:cursor-not-allowed transition-all shadow-lg hover:shadow-emerald-500/30 transform hover:-translate-y-0.5"
            >
              {{ isSubmitting ? 'Creating Recipe...' : 'Create Recipe' }}
            </button>
          </div>
        </Form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { Form, Field, ErrorMessage } from 'vee-validate';
import * as yup from 'yup';

const router = useRouter();
const createError = ref('');
const imageInput = ref(null);
const uploadedImages = ref([]); // Array of {url: string, isFeatured: boolean}
const uploadingImage = ref(false);

// Redirect if not logged in
const token = useCookie('auth_token');
onMounted(() => {
  if (!token.value) {
    router.push('/login');
  }
  fetchCategories();
});

// Fetch Categories
const categories = ref([]);
const selectedCategoryId = ref(null);
const fetchCategories = async () => {
  try {
    const config = useRuntimeConfig();
    const apiUrl = config.public.apiUrl || 'http://localhost:8081';
    const data = await $fetch(`${apiUrl}/categories`);
    categories.value = data || [];
  } catch (err) {
    console.error('Error fetching categories:', err);
  }
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

// Dynamic Ingredients
const ingredients = ref([
  { name: '', quantity: '', unit: '' }
]);

const addIngredient = () => {
  ingredients.value.push({ name: '', quantity: '', unit: '' });
};

const removeIngredient = (index) => {
  if (ingredients.value.length > 1) {
    ingredients.value.splice(index, 1);
  }
};

// Dynamic Steps
const steps = ref([
  { instruction: '', image_url: '' }
]);

const addStep = () => {
  steps.value.push({ instruction: '', image_url: '' });
};

const removeStep = (index) => {
  if (steps.value.length > 1) {
    steps.value.splice(index, 1);
  }
};

// Validation Schema
const schema = yup.object({
  title: yup.string().required().min(3).label('Title'),
  description: yup.string().required().min(10).label('Description'),
  category_id: yup.mixed().test('required', 'Category is required', function(value) {
    // Check both the form value and the selectedCategoryId
    const catId = value || selectedCategoryId.value;
    return catId !== null && catId !== undefined && catId !== '' && catId !== 0;
  }).label('Category'),
  preparation_time: yup.number().required().min(1).label('Preparation Time'),
  price: yup.number().min(0).label('Price'),
});

// Handle Multiple Image Upload
const handleMultipleImageUpload = async (event) => {
  const files = Array.from(event.target.files);
  if (files.length === 0) return;

  uploadingImage.value = true;
  createError.value = '';

  try {
    // Upload all files
    const uploadPromises = files.map(async (file) => {
      const formData = new FormData();
      formData.append('file', file);
      
      const config = useRuntimeConfig();
      const apiUrl = config.public.apiUrl || 'http://localhost:8081';
      const response = await fetch(`${apiUrl}/upload`, {
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
    
    // Add to uploadedImages array
    uploadedImages.value.push(...uploaded);
    
    // If this is the first image, set it as featured automatically
    if (uploadedImages.value.length === uploaded.length) {
      uploadedImages.value[0].isFeatured = true;
    }
  } catch (err) {
    createError.value = err.message || 'Failed to upload images';
    console.error(err);
  } finally {
    uploadingImage.value = false;
    // Reset input
    if (event.target) {
      event.target.value = '';
    }
  }
};

// Set Featured Image
const setFeaturedImage = (index) => {
  // Unset all featured
  uploadedImages.value.forEach(img => img.isFeatured = false);
  // Set selected as featured
  uploadedImages.value[index].isFeatured = true;
};

// Remove Image
const removeImage = (index) => {
  const removed = uploadedImages.value.splice(index, 1)[0];
  // If removed was featured and there are other images, set first as featured
  if (removed.isFeatured && uploadedImages.value.length > 0) {
    uploadedImages.value[0].isFeatured = true;
  }
};

// Handle Form Submit
const handleCreateRecipe = async (values) => {
  createError.value = '';

  // Use selectedCategoryId if category_id from form is missing
  const categoryId = values.category_id || selectedCategoryId.value;
  
  if (!categoryId || categoryId === 0) {
    createError.value = 'Please select a category';
    return;
  }

  // Validate ingredients
  const validIngredients = ingredients.value.filter(ing => ing.name.trim() !== '');
  if (validIngredients.length === 0) {
    createError.value = 'Please add at least one ingredient';
    return;
  }

  // Validate steps
  const validSteps = steps.value.filter(step => step.instruction.trim() !== '');
  if (validSteps.length === 0) {
    createError.value = 'Please add at least one preparation step';
    return;
  }

  // Validate images
  if (uploadedImages.value.length === 0) {
    createError.value = 'Please upload at least one image';
    return;
  }

  // Find featured image
  const featuredImage = uploadedImages.value.find(img => img.isFeatured);
  if (!featuredImage) {
    createError.value = 'Please select a featured image';
    return;
  }

  const recipeData = {
    category_id: parseInt(categoryId),
    title: values.title,
    description: values.description,
    preparation_time: parseInt(values.preparation_time),
    price: parseFloat(values.price) || 0,
    thumbnail_url: featuredImage.url, // Use featured image as thumbnail
    ingredients: validIngredients,
    steps: validSteps,
    images: uploadedImages.value.map(img => img.url) // All images for later upload
  };

  try {
    console.log('[CREATE] Submitting recipe:', recipeData);
    console.log('[CREATE] Token present:', !!token.value);
    
    const config = useRuntimeConfig();
    const apiUrl = config.public.apiUrl || 'http://localhost:8081';
    const response = await fetch(`${apiUrl}/recipes`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token.value}`
      },
      body: JSON.stringify(recipeData)
    });

    console.log('[CREATE] Response status:', response.status);

    if (!response.ok) {
      const errorData = await response.text();
      console.error('[CREATE] Error response:', errorData);
      throw new Error(errorData || 'Failed to create recipe');
    }

    const data = await response.json();
    console.log('[CREATE] Recipe created successfully:', data);
    
    if (!data || !data.id) {
      throw new Error('Recipe created but no ID returned');
    }
    
    // Show success message
    alert('Recipe created successfully! Redirecting to your recipe...');
    
    // Success - redirect to recipe detail page
    await router.push(`/recipes/${data.id}`);
  } catch (err) {
    createError.value = err.message || 'An error occurred while creating the recipe';
    console.error('[CREATE] Exception:', err);
    console.error('[CREATE] Full error:', err);
    
    // Don't navigate on error - stay on create page to show error
    alert('Failed to create recipe: ' + err.message);
  }
};
</script>
