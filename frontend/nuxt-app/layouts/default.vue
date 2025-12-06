<template>
  <div class="min-h-screen bg-gray-50 font-sans selection:bg-indigo-100 selection:text-indigo-700">
    <nav class="sticky top-0 z-50 bg-white/90 backdrop-blur-md border-b border-gray-200 transition-all duration-300">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex">
            <div class="flex-shrink-0 flex items-center">
              <NuxtLink to="/home" class="flex items-center gap-2 group">
                <span class="text-2xl">🍳</span>
                <span class="text-2xl font-extrabold bg-clip-text text-transparent bg-gradient-to-r from-indigo-600 to-purple-600 group-hover:from-indigo-500 group-hover:to-purple-500 transition-all">
                  Chef's Table
                </span>
              </NuxtLink>
            </div>
            <div class="hidden sm:ml-10 sm:flex sm:space-x-8">
              <NuxtLink to="/home" class="border-transparent text-gray-500 hover:text-indigo-600 hover:border-indigo-600 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium transition-colors">
                Discover
              </NuxtLink>
              <NuxtLink v-if="isAuthenticated" to="/create" class="border-transparent text-gray-500 hover:text-indigo-600 hover:border-indigo-600 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium transition-colors">
                Create Recipe
              </NuxtLink>
            </div>
          </div>
          <div class="flex items-center gap-4">
            <div v-if="!isAuthenticated" class="flex items-center space-x-4">
              <NuxtLink to="/login" class="text-gray-500 hover:text-indigo-600 font-medium text-sm transition-colors">Log in</NuxtLink>
              <NuxtLink to="/register" class="bg-indigo-600 text-white hover:bg-indigo-700 px-4 py-2 rounded-full text-sm font-medium shadow-md hover:shadow-lg transition-all transform hover:-translate-y-0.5">
                Sign up
              </NuxtLink>
            </div>
            <div v-else class="flex items-center space-x-4">
              <div class="hidden md:block text-right">
                <p class="text-xs text-gray-400 uppercase tracking-wider font-semibold">Welcome back</p>
                <p class="text-sm font-medium text-gray-900">Chef</p>
              </div>
              <button @click="handleLogout" class="text-gray-400 hover:text-red-600 p-2 rounded-full hover:bg-red-50 transition-all" title="Logout">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>
    </nav>

    <main>
      <slot />
    </main>
  </div>
</template>

<script setup>
const { onLogout } = useApollo();
const router = useRouter();
const token = useCookie('apollo:default.token');

// Computed property to make it reactive
const isAuthenticated = computed(() => !!token.value);

const handleLogout = async () => {
  await onLogout();
  // Force refresh or redirect to ensure state is cleared
  router.push('/login');
};
</script>
