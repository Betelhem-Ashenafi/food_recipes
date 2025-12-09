<template>
  <div class="relative min-h-screen pb-20">
    <!-- Background Image with Overlay - Inspiring Kitchen Scene (Create Page) -->
    <div class="fixed inset-0 z-0">
      <img 
        src="https://images.unsplash.com/photo-1556910103-1c02745aae4d?ixlib=rb-4.0.3&auto=format&fit=crop&w=2070&q=80" 
        alt="Elegant Kitchen" 
        class="w-full h-full object-cover brightness-95"
      >
      <div class="absolute inset-0 bg-gradient-to-b from-black/40 via-black/10 to-black/60"></div>
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
              <label class="block text-sm font-medium text-gray-300 mb-2">Recipe Title</label>
              <Field 
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
              <label class="block text-sm font-medium text-gray-300 mb-2">Description</label>
              <Field 
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
                <label class="block text-sm font-medium text-gray-300 mb-2">Category</label>
                <Field 
                  name="category_id" 
                  as="select"
                  class="block w-full px-4 py-3 border border-white/20 rounded-lg bg-black/20 text-white focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500"
                  :class="{ 'border-red-500 focus:ring-red-500': errors.category_id }"
                >
                  <option value="">Select a category</option>
                  <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
                </Field>
                <ErrorMessage name="category_id" class="text-red-400 text-xs mt-1" />
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-300 mb-2">Preparation Time (minutes)</label>
                <Field 
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
              <label class="block text-sm font-medium text-gray-300 mb-2">Price (0 for free)</label>
              <Field 
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

          <!-- Thumbnail Image -->
          <div class="mb-8">
            <h2 class="text-2xl font-bold text-white mb-6 flex items-center">
              <span class="text-emerald-400 mr-2">📸</span> Featured Image
            </h2>
            
            <div class="border-2 border-dashed border-white/30 rounded-lg p-8 text-center hover:border-emerald-400/50 transition-colors">
              <input 
                type="file" 
                @change="handleImageUpload" 
                accept="image/*"
                class="hidden" 
                ref="imageInput"
              />
              
              <div v-if="!thumbnailUrl" @click="$refs.imageInput.click()" class="cursor-pointer">
                <svg class="mx-auto h-12 w-12 text-gray-400" stroke="currentColor" fill="none" viewBox="0 0 48 48">
                  <path d="M28 8H12a4 4 0 00-4 4v20m32-12v8m0 0v8a4 4 0 01-4 4H12a4 4 0 01-4-4v-4m32-4l-3.172-3.172a4 4 0 00-5.656 0L28 28M8 32l9.172-9.172a4 4 0 015.656 0L28 28m0 0l4 4m4-24h8m-4-4v8m-12 4h.02" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
                </svg>
                <p class="mt-2 text-sm text-gray-300">Click to upload featured image</p>
                <p class="text-xs text-gray-400 mt-1">PNG, JPG, GIF up to 10MB</p>
              </div>

              <div v-else class="relative">
                <img :src="thumbnailUrl" alt="Preview" class="max-h-64 mx-auto rounded-lg" />
                <button 
                  type="button"
                  @click="thumbnailUrl = ''; $refs.imageInput.value = ''"
                  class="absolute top-2 right-2 bg-red-500 text-white p-2 rounded-full hover:bg-red-600"
                >
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>
            </div>
            <p v-if="uploadingImage" class="text-emerald-400 text-sm mt-2">Uploading image...</p>
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
const thumbnailUrl = ref('');
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
const fetchCategories = async () => {
  try {
    const { data } = await useFetch('http://localhost:8081/categories');
    categories.value = data.value || [];
  } catch (err) {
    console.error('Error fetching categories:', err);
  }
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
  category_id: yup.number().required().label('Category'),
  preparation_time: yup.number().required().min(1).label('Preparation Time'),
  price: yup.number().min(0).label('Price'),
});

// Handle Image Upload
const handleImageUpload = async (event) => {
  const file = event.target.files[0];
  if (!file) return;

  uploadingImage.value = true;
  const formData = new FormData();
  formData.append('file', file);

  try {
    const response = await fetch('http://localhost:8081/upload', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token.value}`
      },
      body: formData
    });

    const data = await response.json();
    thumbnailUrl.value = data.url;
  } catch (err) {
    createError.value = 'Failed to upload image';
    console.error(err);
  } finally {
    uploadingImage.value = false;
  }
};

// Handle Form Submit
const handleCreateRecipe = async (values) => {
  createError.value = '';

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

  // Validate image
  if (!thumbnailUrl.value) {
    createError.value = 'Please upload a featured image';
    return;
  }

  const recipeData = {
    category_id: parseInt(values.category_id),
    title: values.title,
    description: values.description,
    preparation_time: parseInt(values.preparation_time),
    price: parseFloat(values.price) || 0,
    thumbnail_url: thumbnailUrl.value,
    ingredients: validIngredients,
    steps: validSteps
  };

  try {
    const response = await fetch('http://localhost:8081/recipes', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token.value}`
      },
      body: JSON.stringify(recipeData)
    });

    if (!response.ok) {
      const errorData = await response.text();
      throw new Error(errorData || 'Failed to create recipe');
    }

    const data = await response.json();
    
    // Success - redirect to recipe detail page
    router.push(`/recipes/${data.id}`);
  } catch (err) {
    createError.value = err.message || 'An error occurred while creating the recipe';
    console.error(err);
  }
};
</script>
