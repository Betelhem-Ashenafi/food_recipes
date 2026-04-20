<template>
  <div class="relative min-h-screen pb-20">
    <!-- Background -->
    <div class="fixed inset-0 z-0">
      <img 
        src="https://images.unsplash.com/photo-1556910103-1c02745aae4d?ixlib=rb-4.0.3&auto=format&fit=crop&w=2070&q=80" 
        alt="Kitchen" 
        class="w-full h-full object-cover brightness-60"
      >
      <!-- Enhanced dark overlay for optimal text readability -->
      <div class="absolute inset-0 bg-black/80"></div>
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
        <h1 class="text-3xl font-bold text-white mb-6 flex items-center">
          <span class="text-emerald-400 mr-2">✏️</span>
          Edit Recipe: {{ recipe.title || 'Loading...' }}
        </h1>
        
        <form @submit.prevent="handleFormSubmit">
          <!-- Basic Information -->
          <div class="mb-8">
            <h2 class="text-2xl font-bold text-white mb-6 flex items-center">
              <span class="text-emerald-400 mr-2">📝</span> Basic Information
            </h2>
            
            <div class="mb-6">
              <label for="recipe-title" class="block text-sm font-medium text-gray-300 mb-2">Recipe Title</label>
              <input 
                id="recipe-title"
                v-model="formTitle"
                type="text" 
                required
                minlength="3"
                class="block w-full px-4 py-3 border border-white/20 rounded-lg bg-black/20 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-emerald-500"
                placeholder="Recipe title"
              />
              <p v-if="formErrors.title" class="text-red-400 text-xs mt-1">{{ formErrors.title }}</p>
            </div>

            <div class="mb-6">
              <label for="recipe-description" class="block text-sm font-medium text-gray-300 mb-2">Description</label>
              <textarea 
                id="recipe-description"
                v-model="formDescription"
                rows="4"
                required
                minlength="10"
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
                    @click="formCategoryId = cat.id"
                    :class="[
                      'group relative overflow-hidden rounded-lg border-2 transition-all duration-300 transform hover:scale-105',
                      formCategoryId === cat.id
                        ? 'border-emerald-400 ring-2 ring-emerald-400/50 shadow-lg shadow-emerald-500/30'
                        : 'border-white/20 hover:border-white/40'
                    ]"
                  >
                    <!-- Category Image -->
                    <div class="relative h-20 md:h-24 overflow-hidden">
                      <img
                        v-if="shouldShowCategoryImage(cat)"
                        :src="normalizedCategoryImageUrl(cat)"
                        :alt="cat.name"
                        class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-110"
                        @error="markCategoryImageFailed(cat.id)"
                      />
                      <!-- Fallback if no image -->
                      <div 
                        v-else
                        class="w-full h-full bg-gradient-to-br from-emerald-500/40 via-teal-500/40 to-blue-500/40 flex items-center justify-center"
                      >
                        <span class="text-white font-semibold text-sm text-center px-2">{{ cat.name }}</span>
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
                <label for="recipe-preparation-time" class="block text-sm font-medium text-gray-300 mb-2">Preparation Time (minutes)</label>
                <input 
                  id="recipe-preparation-time"
                  v-model.number="formPreparationTime"
                  type="number"
                  required
                  min="1"
                  class="block w-full px-4 py-3 border border-white/20 rounded-lg bg-black/20 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-emerald-500"
                  placeholder="e.g., 45"
                />
                <p v-if="formErrors.preparation_time" class="text-red-400 text-xs mt-1">{{ formErrors.preparation_time }}</p>
              </div>
            </div>

            <div class="mt-6">
              <label for="recipe-price" class="block text-sm font-medium text-gray-300 mb-2">Price (0 for free)</label>
              <input 
                id="recipe-price"
                v-model.number="formPrice"
                type="number"
                step="0.01"
                min="0"
                class="block w-full px-4 py-3 border border-white/20 rounded-lg bg-black/20 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-emerald-500"
                placeholder="0.00"
              />
              <p v-if="formErrors.price" class="text-red-400 text-xs mt-1">{{ formErrors.price }}</p>
            </div>
          </div>

          
          <!-- Ingredients -->
          <div class="mb-8">
            <h2 class="text-2xl font-bold text-white mb-6 flex items-center">
              <span class="text-emerald-400 mr-2">🥕</span> Ingredients
            </h2>
            <div v-for="(ing, index) in ingredients" :key="index" class="mb-4 flex gap-3">
              <input v-model="ing.name" placeholder="Name" class="flex-1 px-4 py-3 border border-white/20 rounded-lg bg-black/20 text-white" />
              <input v-model.number="ing.quantity" type="number" min="0" step="any" placeholder="Qty" class="w-24 px-4 py-3 border border-white/20 rounded-lg bg-black/20 text-white" />
              <select v-model="ing.unit_id" class="w-40 px-4 py-3 border border-white/20 rounded-lg bg-black/20 text-white">
                <option value="" disabled>Select</option>
                <option v-for="u in units" :key="u.id" :value="String(u.id)">{{ u.name }}</option>
              </select>
              <button type="button" @click="removeIngredient(index)" class="px-4 py-3 bg-red-500/20 text-red-400 rounded-lg">Remove</button>
            </div>
            <button type="button" @click="addIngredient" class="w-full px-4 py-3 border-2 border-dashed border-emerald-400/50 rounded-lg text-emerald-400">+ Add Ingredient</button>
          </div>

          <!-- Steps -->
          <div class="mb-8">
            <h2 class="text-2xl font-bold text-white mb-6 flex items-center">
              <span class="text-emerald-400 mr-2">📋</span> Preparation Steps
            </h2>
            <div v-for="(step, index) in steps" :key="index" class="mb-4">
              <textarea v-model="step.instruction" rows="3" placeholder="Step instruction" class="w-full px-4 py-3 border border-white/20 rounded-lg bg-black/20 text-white"></textarea>
              <button type="button" @click="removeStep(index)" class="mt-2 px-4 py-2 bg-red-500/20 text-red-400 rounded-lg">Remove</button>
            </div>
            <button type="button" @click="addStep" class="w-full px-4 py-3 border-2 border-dashed border-emerald-400/50 rounded-lg text-emerald-400">+ Add Step</button>
          </div>

          <div v-if="updateError" class="mb-6 p-4 rounded-lg bg-red-500/20 border border-red-500/50 text-red-200">
            <strong>Error:</strong> {{ updateError }}
          </div>
          
          <!-- Show validation errors summary -->
          <div v-if="Object.keys(formErrors).length > 0" class="mb-6 p-4 rounded-lg bg-yellow-500/20 border border-yellow-500/50 text-yellow-200">
            <strong>Please fix the following errors:</strong>
            <ul class="list-disc list-inside mt-2">
              <li v-for="(error, field) in formErrors" :key="field">{{ error }}</li>
            </ul>
          </div>

          <div class="flex gap-4">
            <button type="button" @click="$router.push(`/recipes/${recipeId}`)" class="flex-1 px-6 py-4 border border-white/30 rounded-lg text-white">Cancel</button>
            <button 
              type="submit" 
              :disabled="isSubmitting || uploadingImage"
              class="flex-1 px-6 py-4 bg-gradient-to-r from-emerald-600 to-teal-600 text-white font-bold rounded-lg disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {{ isSubmitting ? 'Updating...' : 'Update Recipe' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { jwtDecode } from 'jwt-decode';
import { useApolloClient } from '@vue/apollo-composable';
import gql from 'graphql-tag';

const route = useRoute();
const router = useRouter();
const recipeId = parseInt(route.params.id);
const token = useCookie('auth_token');
const loading = ref(true);
const recipe = ref(null);
const categories = ref([]);
const units = ref([]);
const ingredients = ref([]);
const steps = ref([]);
const updateError = ref('');
const imageInput = ref(null);
// Edit image logic removed
const uploadedImages = ref([]);
const uploadingImage = ref(false);
const isSubmitting = ref(false);
const formErrors = ref({});
const failedCategoryImageIds = ref(new Set());
const { client } = useApolloClient();

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

const RECIPE_QUERY = gql`
  query GetRecipe($id: Int!) {
    recipes_by_pk(id: $id) {
      id
      title
      description
      category_id
      preparation_time
      price
      user_id
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

const RECIPE_STEPS_QUERY = gql`
  query GetRecipeSteps($id: Int!) {
    recipe_steps(where: { recipe_id: { _eq: $id } }, order_by: { step_number: asc }) {
      id
      step_number
      instruction
      image_url
    }
  }
`;

const RECIPE_IMAGES_QUERY = gql`
  query GetRecipeImages($id: Int!) {
    recipe_images(where: { recipe_id: { _eq: $id } }, order_by: { id: asc }) {
      id
      url
      is_featured
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

const UNITS_QUERY = gql`
  query GetUnits {
    units(order_by: { id: asc }) {
      id
      name
    }
  }
`;

const UPDATE_RECIPE_MUTATION = gql`
  mutation UpdateRecipe(
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

const UPLOAD_FILE_MUTATION = gql`
  mutation UploadFile($file: UploadInput!) {
    uploadFile(file: $file) {
      url
    }
  }
`;

// Form values - use individual refs for maximum compatibility
const formTitle = ref('');
const formDescription = ref('');
const formCategoryId = ref(0);
const formPreparationTime = ref(0);
const formPrice = ref(0);
const createIngredient = () => ({ name: '', quantity: '', unit_id: '' });
const createStep = () => ({ instruction: '', image_url: '' });

// Get current user ID from token
const getCurrentUserId = () => {
  if (!token.value) return null;
  try {
    const decoded = jwtDecode(token.value);
    const claims = decoded?.['https://hasura.io/jwt/claims'];
    if (claims && claims['x-hasura-user-id']) {
      return parseInt(claims['x-hasura-user-id']);
    }
    // Fallback to direct user_id claim
    if (decoded?.user_id) {
      return typeof decoded.user_id === 'string' ? parseInt(decoded.user_id) : decoded.user_id;
    }
    return null;
  } catch {
    return null;
  }
};

// Redirect if not logged in
onMounted(async () => {
  if (!token.value) {
    router.push('/login');
    return;
  }

  const currentUserId = getCurrentUserId();

  try {
    const [recipeResult, ingredientsResult, stepsResult, imagesResult, categoriesResult, unitsResult] = await Promise.all([
      client.query({ query: RECIPE_QUERY, variables: { id: recipeId }, fetchPolicy: 'network-only' }),
      client.query({ query: RECIPE_INGREDIENTS_QUERY, variables: { id: recipeId }, fetchPolicy: 'network-only' }),
      client.query({ query: RECIPE_STEPS_QUERY, variables: { id: recipeId }, fetchPolicy: 'network-only' }),
      client.query({ query: RECIPE_IMAGES_QUERY, variables: { id: recipeId }, fetchPolicy: 'network-only' }),
      client.query({ query: CATEGORIES_QUERY, fetchPolicy: 'network-only' }),
      client.query({ query: UNITS_QUERY, fetchPolicy: 'network-only' })
    ]);

    recipe.value = recipeResult?.data?.recipes_by_pk || null;
    if (!recipe.value) {
      throw new Error('Recipe not found');
    }

    // Check ownership
    const recipeUserId = recipe.value?.user_id;
    if (recipeUserId && currentUserId && recipeUserId !== currentUserId) {
      updateError.value = 'You do not have permission to edit this recipe';
      router.push(`/recipes/${recipeId}`);
      return;
    }

    if (recipe.value) {
      formTitle.value = recipe.value.title || '';
      formDescription.value = recipe.value.description || '';
      formCategoryId.value = recipe.value.category_id || 0;
      formPreparationTime.value = recipe.value.preparation_time || 0;
      formPrice.value = recipe.value.price || 0;
    }

    const ingData = ingredientsResult?.data?.recipe_ingredients || [];
    ingredients.value = ingData.map(ing => ({
      name: ing.name || '',
      quantity: ing.quantity || '',
      unit_id: ing.unit_id ? String(ing.unit_id) : ''
    }));
    if (ingredients.value.length === 0) {
      ingredients.value = [createIngredient()];
    }

    const stepsData = stepsResult?.data?.recipe_steps || [];
    steps.value = stepsData.map(step => ({
      instruction: step.instruction || '',
      image_url: step.image_url || ''
    }));
    if (steps.value.length === 0) {
      steps.value = [createStep()];
    }

    categories.value = categoriesResult?.data?.categories || [];
    units.value = unitsResult?.data?.units || [];
  } catch (err) {
    updateError.value = err.message || 'Failed to load recipe data';
  } finally {
    loading.value = false;
  }
});

const addIngredient = () => ingredients.value.push(createIngredient());
const removeIngredient = (index) => {
  if (ingredients.value.length > 1) ingredients.value.splice(index, 1);
};
const addStep = () => steps.value.push(createStep());
const removeStep = (index) => {
  if (steps.value.length > 1) steps.value.splice(index, 1);
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
const handleFormSubmit = async (event) => {
  if (event) {
    event.preventDefault();
    event.stopPropagation();
  }

  updateError.value = '';
  formErrors.value = {};

  const isValid = validateForm();
  if (!isValid) {
    updateError.value = 'Please fix the errors above before submitting.';
    return;
  }

  updateError.value = '';
  isSubmitting.value = true;

  // Validate ingredients
  const validIngredients = ingredients.value.filter(ing => ing.name && ing.name.trim() !== '');
  if (validIngredients.length === 0) {
    updateError.value = 'Please add at least one ingredient';
    isSubmitting.value = false;
    return;
  }

  const hasMissingUnit = validIngredients.some((ing) => !ing.unit_id);
  if (hasMissingUnit) {
    updateError.value = 'Please select a unit for every ingredient';
    isSubmitting.value = false;
    return;
  }

  const hasInvalidQuantity = validIngredients.some((ing) => {
    if (ing.quantity === '' || ing.quantity === null || ing.quantity === undefined) return true;
    const quantityValue = Number(ing.quantity);
    return Number.isNaN(quantityValue) || quantityValue <= 0;
  });
  if (hasInvalidQuantity) {
    updateError.value = 'Please enter a valid quantity (number > 0) for every ingredient';
    isSubmitting.value = false;
    return;
  }

  // Validate steps
  const validSteps = steps.value.filter(step => step.instruction && step.instruction.trim() !== '');
  if (validSteps.length === 0) {
    updateError.value = 'Please add at least one preparation step';
    isSubmitting.value = false;
    return;
  }

  // Format ingredients to match schema (name, quantity, unit_id)
  const formattedIngredients = validIngredients.map(ing => ({
    name: ing.name || '',
    quantity: ing.quantity || '',
    unit_id: parseInt(ing.unit_id)
  }));
  
  
  
  // Use values from form submission
  const recipeData = {
    category_id: parseInt(formCategoryId.value),
    title: formTitle.value,
    description: formDescription.value,
    preparation_time: parseInt(formPreparationTime.value),
    price: parseFloat(formPrice.value) || 0,
    ingredients: formattedIngredients,
    steps: formattedSteps,
    images: uploadedImages.value.map(img => img.url)
  };

  try {
    const recipeSet = {
      category_id: recipeData.category_id,
      title: recipeData.title,
      description: recipeData.description,
      preparation_time: recipeData.preparation_time,
      price: recipeData.price
    };

    const ingredientsPayload = formattedIngredients.map((ing) => ({
      recipe_id: recipeId,
      name: ing.name,
      quantity: ing.quantity,
      unit_id: ing.unit_id
    }));

    const stepsPayload = formattedSteps.map((step, index) => ({
      recipe_id: recipeId,
      step_number: index + 1,
      instruction: step.instruction,
      image_url: step.image_url || ''
    }));

    const imagesPayload = uploadedImages.value.map((img) => ({
      recipe_id: recipeId,
      url: img.url,
      is_featured: img.isFeatured
    }));

    await client.mutate({
      mutation: UPDATE_RECIPE_MUTATION,
      variables: {
        recipeId,
        recipe: recipeSet,
        ingredients: ingredientsPayload,
        steps: stepsPayload,
        images: imagesPayload
      }
    });

    await router.push(`/recipes/${recipeId}`);
  } catch (err) {
    updateError.value = err.message || 'An error occurred while updating the recipe';
  } finally {
    isSubmitting.value = false;
  }
};
</script>