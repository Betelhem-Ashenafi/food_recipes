<template>
  <div 
    @click="navigateToHome" 
    class="relative min-h-screen w-full overflow-hidden cursor-pointer group"
  >
    <!-- Background Image with Overlay -->
    <div class="absolute inset-0">
      <img 
        src="https://images.unsplash.com/photo-1504674900247-0877df9cc836?ixlib=rb-4.0.3&auto=format&fit=crop&w=2070&q=80" 
        alt="Elegant Food" 
        class="w-full h-full object-cover brightness-95 transition-transform duration-[10s] ease-out group-hover:scale-105"
      >
      <!-- Unified Gradient Overlay -->
      <div class="absolute inset-0 bg-gradient-to-b from-black/40 via-black/10 to-black/60"></div>
    </div>

    <!-- Content -->
    <div class="relative z-10 flex flex-col items-center justify-center min-h-screen text-center px-4">
      <!-- Logo / Icon -->
      <div class="mb-8 opacity-0 animate-fade-in-down" style="animation-delay: 0.2s; animation-fill-mode: forwards;">
        <span class="text-6xl">👨‍🍳</span>
      </div>

      <!-- Main Title -->
      <h1 class="font-serif text-5xl md:text-7xl font-bold text-white mb-4 tracking-wide opacity-0 animate-fade-in-up" style="animation-delay: 0.5s; animation-fill-mode: forwards;">
        <span class="text-emerald-400">Chef's</span> Table
      </h1>

      <!-- Subtitle -->
      <p class="text-gray-300 text-lg md:text-xl max-w-lg font-light tracking-wider mb-12 opacity-0 animate-fade-in-up" style="animation-delay: 0.8s; animation-fill-mode: forwards;">
        Experience the art of culinary excellence. <br>
        Discover, Create, and Share.
      </p>

      <!-- Call to Action Indicator -->
      <div class="opacity-0 animate-pulse" style="animation-delay: 1.5s; animation-fill-mode: forwards;">
        <span class="text-emerald-400 text-sm uppercase tracking-[0.3em] border-b border-emerald-400 pb-1">
          Tap anywhere to begin
        </span>
      </div>
    </div>
  </div>
</template>

<script setup>
const router = useRouter();

const navigateToHome = () => {
  router.push('/home');
};

// Define layout as 'none' or 'blank' if you don't want the navbar on this splash screen
definePageMeta({
  layout: 'blank'
});
</script>

<style scoped>
@keyframes fade-in-down {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes fade-in-up {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate-fade-in-down {
  animation: fade-in-down 1s ease-out;
}

.animate-fade-in-up {
  animation: fade-in-up 1s ease-out;
}
</style>
