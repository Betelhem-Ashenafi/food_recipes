<template>
  <div class="relative min-h-screen w-full overflow-hidden flex items-center justify-center px-4">
    <!-- Background Image with Overlay (Unified with Home/Index) -->
    <div class="absolute inset-0 z-0">
      <img 
        src="https://images.unsplash.com/photo-1504674900247-0877df9cc836?ixlib=rb-4.0.3&auto=format&fit=crop&w=2070&q=80" 
        alt="Elegant Food" 
        class="w-full h-full object-cover brightness-60"
      >
      <div class="absolute inset-0 bg-black/80"></div>
    </div>

    <!-- Glassmorphism Card -->
    <div class="relative z-10 w-full max-w-md transform transition-all duration-500 hover:scale-[1.01]">
      <div class="bg-white/10 backdrop-blur-lg border border-white/20 rounded-2xl shadow-2xl p-8 sm:p-10">
        
        <!-- Header -->
        <div class="text-center mb-10">
          <div class="text-5xl mb-4">👨‍🍳</div>
          <h2 class="text-3xl font-serif font-bold text-white tracking-wide">
            Welcome Back
          </h2>
          <p class="mt-2 text-gray-300 text-sm">
            Sign in to continue your culinary journey
          </p>
        </div>
        
        <Form @submit="handleLogin" :validation-schema="schema" v-slot="{ errors, isSubmitting }" class="space-y-6">
          
          <!-- Email Field -->
          <div>
            <label for="email" class="sr-only">Email address</label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <svg class="h-5 w-5 text-gray-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                  <path d="M2.003 5.884L10 9.882l7.997-3.998A2 2 0 0016 4H4a2 2 0 00-1.997 1.884z" />
                  <path d="M18 8.118l-8 4-8-4V14a2 2 0 002 2h12a2 2 0 002-2V8.118z" />
                </svg>
              </div>
              <Field 
                name="email" 
                type="email" 
                class="block w-full pl-10 pr-3 py-3 border border-white/10 rounded-lg leading-5 bg-black/20 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 sm:text-sm transition-colors" 
                placeholder="Email address" 
                :class="{ 'border-red-500 focus:ring-red-500': errors.email }"
              />
            </div>
            <ErrorMessage name="email" class="text-red-400 text-xs mt-1 ml-1" />
          </div>

          <!-- Password Field -->
          <div>
            <label for="password" class="sr-only">Password</label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <svg class="h-5 w-5 text-gray-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z" clip-rule="evenodd" />
                </svg>
              </div>
              <Field 
                name="password" 
                :type="showPassword ? 'text' : 'password'" 
                class="block w-full pl-10 pr-10 py-3 border border-white/10 rounded-lg leading-5 bg-black/20 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 sm:text-sm transition-colors" 
                placeholder="Password" 
                :class="{ 'border-red-500 focus:ring-red-500': errors.password }"
              />
              <button
                type="button"
                @click="showPassword = !showPassword"
                class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-300 focus:outline-none transition-colors"
                tabindex="-1"
              >
                <!-- Eye icon (when password is hidden) -->
                <svg v-if="!showPassword" class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                  <path d="M10 12a2 2 0 100-4 2 2 0 000 4z" />
                  <path fill-rule="evenodd" d="M.458 10C1.732 5.943 5.522 3 10 3s8.268 2.943 9.542 7c-1.274 4.057-5.064 7-9.542 7S1.732 14.057.458 10zM14 10a4 4 0 11-8 0 4 4 0 018 0z" clip-rule="evenodd" />
                </svg>
                <!-- Eye-slash icon (when password is visible) -->
                <svg v-else class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M3.707 2.293a1 1 0 00-1.414 1.414l14 14a1 1 0 001.414-1.414l-1.473-1.473A10.014 10.014 0 0019.542 10C18.268 5.943 14.478 3 10 3a9.994 9.994 0 00-4.838 1.232L3.707 2.293zM14.95 6.05a4 4 0 00-5.9 5.9l1.519-1.52A2.5 2.5 0 0110 8.5c.36 0 .69.08.99.22l1.519-1.52zM2.853 4.853l1.414 1.414A9.975 9.975 0 00.458 10c1.274 4.057 5.064 7 9.542 7 2.29 0 4.408-.738 6.131-1.997l1.414 1.414a1 1 0 001.414-1.414l-14-14a1 1 0 00-1.414 1.414zm2.94 2.94l1.415 1.415A2.5 2.5 0 005 10a4 4 0 004 4 2.5 2.5 0 001.792-.77l1.415 1.415A9.975 9.975 0 0110 13c-4.478 0-8.268-2.943-9.542-7a9.97 9.97 0 015.293-4.207z" clip-rule="evenodd" />
                </svg>
              </button>
            </div>
            <ErrorMessage name="password" class="text-red-400 text-xs mt-1 ml-1" />
          </div>

          <!-- Error Message -->
          <div v-if="loginError" class="p-3 rounded-lg bg-red-500/20 border border-red-500/50 text-red-200 text-sm text-center">
            {{ loginError }}
          </div>

          <!-- Submit Button -->
          <div>
            <button 
              type="submit" 
              :disabled="isSubmitting"
              class="group relative w-full flex justify-center py-3 px-4 border border-transparent text-sm font-bold rounded-lg text-white bg-gradient-to-r from-emerald-600 to-teal-600 hover:from-emerald-500 hover:to-teal-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-emerald-500 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-300 shadow-lg hover:shadow-emerald-500/30 transform hover:-translate-y-0.5"
            >
              <span class="absolute left-0 inset-y-0 flex items-center pl-3">
                <svg v-if="!isSubmitting" class="h-5 w-5 text-emerald-200 group-hover:text-emerald-100" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                  <path fill-rule="evenodd" d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z" clip-rule="evenodd" />
                </svg>
                <svg v-else class="animate-spin h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
              </span>
              {{ isSubmitting ? 'Signing in...' : 'Sign In' }}
            </button>
          </div>
        </Form>

        <!-- Footer Links -->
        <div class="mt-8 text-center">
          <p class="text-sm text-gray-300">
            Don't have an account?
            <NuxtLink to="/register" class="font-medium text-emerald-400 hover:text-emerald-300 transition-colors underline decoration-emerald-400/30 hover:decoration-emerald-300">
              Create one now
            </NuxtLink>
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { Form, Field, ErrorMessage } from 'vee-validate';
import * as yup from 'yup';
import { gql } from '@apollo/client/core';

const router = useRouter();
const loginError = ref('');
const showPassword = ref(false);

// Redirect if already logged in
onMounted(() => {
  const token = useCookie('auth_token');
  if (token.value) {
    router.push('/home');
  }
});

// Define layout as 'blank' to hide the default navbar
definePageMeta({
  layout: 'blank'
});

// 1. Define the Validation Rules (Schema)
const schema = yup.object({
  email: yup.string().required().email().label('Email Address'),
  password: yup.string().required().min(6).label('Password'),
});

// 2. Handle the Form Submit (REAL backend integration via REST)
const handleLogin = async (values) => {
  loginError.value = '';
  
  try {
    // Call REAL backend login endpoint that queries REAL database
    const config = useRuntimeConfig();
    const apiUrl = config.public.apiUrl || 'http://localhost:8081';
    console.log('[LOGIN] Using API URL:', apiUrl); // Debug log
    const data = await $fetch(`${apiUrl}/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        email: values.email,
        password: values.password
      })
    });

    if (data?.error) {
      loginError.value = data.error;
      return;
    }

    if (data?.token) {
      // Store REAL JWT token in cookie for subsequent requests
      const cookie = useCookie('auth_token');
      cookie.value = data.token;
      
      // Store user info in localStorage
      if (data?.user) {
        localStorage.setItem('user_id', data.user.id);
        localStorage.setItem('user_name', data.user.name);
        localStorage.setItem('user_email', data.user.email);
      }
      
      console.log('Login successful - token stored:', data.token.substring(0, 20) + '...');
      router.push('/home');
    } else {
      loginError.value = 'Login failed - no token received';
    }
    
  } catch (err) {
    console.error('Login error:', err);
    if (err.data?.error) {
      loginError.value = err.data.error;
    } else if (err.statusCode === 401) {
      loginError.value = 'Invalid email or password';
    } else {
      loginError.value = 'Login failed - please try again';
    }
  }
};
</script>
