<template>
  <div class="min-h-screen font-sans">
    <nav class="sticky top-0 z-50 bg-black/40 backdrop-blur-lg border-b border-white/10 transition-all duration-300">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex">
            <div class="flex-shrink-0 flex items-center">
              <NuxtLink to="/home" class="flex items-center gap-2 group">
                <span class="text-2xl">👨‍🍳</span>
                <span class="text-xl sm:text-2xl font-serif font-bold bg-clip-text text-transparent bg-gradient-to-r from-emerald-400 to-teal-400 group-hover:from-emerald-300 group-hover:to-teal-300 transition-all">
                  Chef's Table
                </span>
              </NuxtLink>
            </div>
            <!-- Desktop Navigation -->
            <div class="hidden sm:ml-10 sm:flex sm:space-x-8">
              <NuxtLink to="/home" class="border-transparent text-gray-300 hover:text-emerald-400 hover:border-emerald-400 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium transition-colors">
                Discover
              </NuxtLink>
              <NuxtLink v-if="isAuthenticated" to="/create" class="border-transparent text-gray-300 hover:text-emerald-400 hover:border-emerald-400 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium transition-colors">
                Create Recipe
              </NuxtLink>
            </div>
          </div>
          <div class="flex items-center gap-4">
            <!-- Mobile Menu Button -->
            <button 
              v-if="isAuthenticated"
              @click="mobileMenuOpen = !mobileMenuOpen"
              class="sm:hidden text-gray-300 hover:text-emerald-400 p-2 rounded-md transition-colors"
              aria-label="Toggle menu"
            >
              <svg v-if="!mobileMenuOpen" xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
              </svg>
              <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
            
            <!-- Desktop Auth Buttons -->
            <div v-if="!isAuthenticated" class="hidden sm:flex items-center space-x-4">
              <NuxtLink to="/login" class="text-gray-300 hover:text-emerald-400 font-medium text-sm transition-colors">Log in</NuxtLink>
              <NuxtLink to="/register" class="bg-gradient-to-r from-emerald-600 to-teal-600 text-white hover:from-emerald-500 hover:to-teal-500 px-4 py-2 rounded-full text-sm font-medium shadow-lg hover:shadow-emerald-500/30 transition-all transform hover:-translate-y-0.5">
                Sign up
              </NuxtLink>
            </div>
            <div v-else class="hidden sm:flex items-center space-x-4">
              <NuxtLink to="/profile" class="text-gray-300 hover:text-emerald-400 font-medium text-sm transition-colors flex items-center gap-2">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                </svg>
                Profile
              </NuxtLink>
              <button @click="handleLogout" class="text-gray-400 hover:text-red-400 p-2 rounded-full hover:bg-red-500/10 transition-all" title="Logout">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                </svg>
              </button>
            </div>
            
            <!-- Mobile Auth Buttons (when not authenticated) -->
            <div v-if="!isAuthenticated" class="sm:hidden flex items-center space-x-2">
              <NuxtLink to="/login" class="text-gray-300 hover:text-emerald-400 font-medium text-xs transition-colors">Log in</NuxtLink>
              <NuxtLink to="/register" class="bg-gradient-to-r from-emerald-600 to-teal-600 text-white hover:from-emerald-500 hover:to-teal-500 px-3 py-1.5 rounded-full text-xs font-medium shadow-lg">
                Sign up
              </NuxtLink>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Mobile Menu Dropdown -->
      <div 
        v-if="mobileMenuOpen && isAuthenticated" 
        class="sm:hidden border-t border-white/10 bg-black/60 backdrop-blur-lg"
      >
        <div class="px-4 pt-2 pb-4 space-y-1">
          <NuxtLink 
            to="/home" 
            @click="mobileMenuOpen = false"
            class="block px-3 py-2 text-gray-300 hover:text-emerald-400 hover:bg-white/5 rounded-md text-base font-medium transition-colors"
          >
            Discover
          </NuxtLink>
          <NuxtLink 
            to="/create" 
            @click="mobileMenuOpen = false"
            class="block px-3 py-2 text-emerald-400 hover:text-emerald-300 hover:bg-emerald-500/10 rounded-md text-base font-medium transition-colors flex items-center gap-2"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            Create Recipe
          </NuxtLink>
          <NuxtLink 
            to="/profile" 
            @click="mobileMenuOpen = false"
            class="block px-3 py-2 text-gray-300 hover:text-emerald-400 hover:bg-white/5 rounded-md text-base font-medium transition-colors flex items-center gap-2"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
            </svg>
            Profile
          </NuxtLink>
          <button 
            @click="handleLogout"
            class="w-full text-left px-3 py-2 text-red-400 hover:text-red-300 hover:bg-red-500/10 rounded-md text-base font-medium transition-colors flex items-center gap-2"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
            </svg>
            Logout
          </button>
        </div>
      </div>
    </nav>

    <main>
      <slot />
    </main>
  </div>
</template>

<script setup>
import { ref } from 'vue';

const router = useRouter();
const token = useCookie('auth_token');

// Mobile menu state
const mobileMenuOpen = ref(false);

// Computed property to make it reactive
const isAuthenticated = computed(() => !!token.value);

const handleLogout = async () => {
  // Clear the auth token
  token.value = null;
  // Close mobile menu
  mobileMenuOpen.value = false;
  // Force refresh or redirect to ensure state is cleared
  router.push('/login');
};
</script>
