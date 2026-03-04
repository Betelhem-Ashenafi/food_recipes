<template>
  <div class="relative min-h-screen w-full overflow-hidden flex items-center justify-center px-4">
    <!-- Background image with dark overlay -->
    <div class="absolute inset-0 z-0">
      <img 
        src="https://images.unsplash.com/photo-1504674900247-0877df9cc836?auto=format&fit=crop&w=2070&q=80" 
        alt="Food background" 
        class="w-full h-full object-cover brightness-60"
      >
      <div class="absolute inset-0 bg-black/80"></div>
    </div>

    <!-- Login card -->
    <div class="relative z-10 w-full max-w-md">
      <div class="bg-white/10 backdrop-blur-lg border border-white/20 rounded-2xl shadow-2xl p-8 sm:p-10">
        <!-- Header -->
        <div class="text-center mb-10">
          <div class="text-4xl font-serif text-white mb-3">👨‍🍳 Chef</div>
          <h2 class="text-3xl font-serif font-bold text-white">Welcome Back</h2>
          <p class="mt-2 text-gray-300 text-sm">Sign in to continue your culinary journey</p>
        </div>

        <!-- Login form with Vee-Validate -->
        <Form @submit="handleLogin" :validation-schema="schema" v-slot="{ errors, isSubmitting }" class="space-y-6">
          <!-- Email field -->
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

          <!-- Password field -->
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
                <svg v-if="!showPassword" class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                  <path d="M10 12a2 2 0 100-4 2 2 0 000 4z" />
                  <path fill-rule="evenodd" d="M.458 10C1.732 5.943 5.522 3 10 3s8.268 2.943 9.542 7c-1.274 4.057-5.064 7-9.542 7S1.732 14.057.458 10zM14 10a4 4 0 11-8 0 4 4 0 018 0z" clip-rule="evenodd" />
                </svg>
                <svg v-else class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M3.707 2.293a1 1 0 00-1.414 1.414l14 14a1 1 0 001.414-1.414l-1.473-1.473A10.014 10.014 0 0019.542 10C18.268 5.943 14.478 3 10 3a9.994 9.994 0 00-4.838 1.232L3.707 2.293zM14.95 6.05a4 4 0 00-5.9 5.9l1.519-1.52A2.5 2.5 0 0110 8.5c.36 0 .69.08.99.22l1.519-1.52zM2.853 4.853l1.414 1.414A9.975 9.975 0 00.458 10c1.274 4.057 5.064 7 9.542 7 2.29 0 4.408-.738 6.131-1.997l1.414 1.414a1 1 0 001.414-1.414l-14-14a1 1 0 00-1.414 1.414zm2.94 2.94l1.415 1.415A2.5 2.5 0 005 10a4 4 0 004 4 2.5 2.5 0 001.792-.77l1.415 1.415A9.975 9.975 0 0110 13c-4.478 0-8.268-2.943-9.542-7a9.97 9.97 0 015.293-4.207z" clip-rule="evenodd" />
                </svg>
              </button>
            </div>
            <ErrorMessage name="password" class="text-red-400 text-xs mt-1 ml-1" />
          </div>

          <!-- Error message display -->
          <div v-if="loginError" class="p-3 rounded-lg bg-red-500/20 border border-red-500/50 text-red-200 text-sm text-center">
            {{ loginError }}
          </div>

          <!-- Submit button -->
          <button 
            type="submit" 
            :disabled="isSubmitting"
            class="w-full py-3 px-4 bg-gradient-to-r from-emerald-600 to-teal-600 text-white font-bold rounded-lg hover:from-emerald-500 hover:to-teal-500 disabled:opacity-50 transition-all"
          >
            <span v-if="!isSubmitting">Sign In</span>
            <span v-else class="flex items-center justify-center">
              <svg class="animate-spin h-5 w-5 mr-2 text-white" viewBox="0 0 24 24">...</svg>
              Signing in...
            </span>
          </button>
        </Form>

        <!-- Link to register -->
        <div class="mt-8 text-center">
          <p class="text-sm text-gray-300">
            Don't have an account?
            <NuxtLink to="/register" class="text-emerald-400 hover:text-emerald-300 underline">
              Create one now
            </NuxtLink>
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { Form, Field, ErrorMessage } from 'vee-validate';
import * as yup from 'yup';
import { useMutation } from '@vue/apollo-composable';
import { gql } from '@apollo/client/core';
import { useRouter } from 'vue-router';
import { useCookie } from '#app';

const router = useRouter();
const loginError = ref('');
const showPassword = ref(false);

// Redirect to home if already logged in
onMounted(() => {
  const token = useCookie('auth_token');
  if (token.value) router.push('/home');
});

// Use blank layout (no navbar)
definePageMeta({ layout: 'blank' });

// Validation rules
const schema = yup.object({
  email: yup.string().required().email(),
  password: yup.string().required().min(6),
});

// GraphQL login mutation (via Hasura action)
const LOGIN_MUTATION = gql`
  mutation Login($email: String!, $password: String!) {
    login(email: $email, password: $password) {
      token
      user_id
      name
      email
    }
  }
`;

// Apollo mutation hook
const { mutate: login, onDone, onError } = useMutation(LOGIN_MUTATION);

// Form submission
const handleLogin = (values) => {
  loginError.value = '';
  login({ email: values.email, password: values.password });
};

// On successful login
onDone((result) => {
  const { token, user_id, name, email } = result.data.login;
  // Store token in cookie
  const cookie = useCookie('auth_token');
  cookie.value = token;
  // Save user info locally
  localStorage.setItem('user_id', user_id);
  localStorage.setItem('user_name', name);
  localStorage.setItem('user_email', email);
  // Go to home page
  router.push('/home');
});

// On error
onError((err) => {
  loginError.value = err.graphQLErrors?.[0]?.message || 'Login failed';
});
</script>