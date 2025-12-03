<template>
  <div class="min-h-screen bg-gray-100">
    <nav class="bg-white shadow-md p-4">
      <div class="container mx-auto flex justify-between items-center">
        <NuxtLink to="/" class="text-xl font-bold text-green-600">Food Recipes</NuxtLink>
        <div>
          <NuxtLink to="/" class="mr-4 text-gray-600 hover:text-green-600">Home</NuxtLink>
          <NuxtLink v-if="!isAuthenticated" to="/login" class="text-gray-600 hover:text-green-600">Login</NuxtLink>
          <button v-else @click="logout" class="text-gray-600 hover:text-green-600">Logout</button>
        </div>
      </div>
    </nav>
    <main class="container mx-auto p-4">
      <slot />
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';

const isAuthenticated = ref(false);
const router = useRouter();

onMounted(() => {
  isAuthenticated.value = !!localStorage.getItem('token');
});

const logout = () => {
  localStorage.removeItem('token');
  isAuthenticated.value = false;
  router.push('/login');
};
</script>
