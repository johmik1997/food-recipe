<template>
    <footer class="bg-gray-900 text-white pt-16 pb-8">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <!-- Main Footer Content -->
        <div class="grid grid-cols-1 md:grid-cols-4 gap-12">
          
          <!-- Brand Column -->
          <div class="space-y-4">
            <NuxtLink to="/" class="flex items-center gap-2">
              <span class="text-2xl font-bold text-green-400">🍲 Kushna</span>
            </NuxtLink>
            <p class="text-gray-400">
              Discover, create, and share delicious recipes with our global food community.
            </p>
            <div class="flex space-x-4">
              <a href="#" class="text-gray-400 hover:text-green-400 transition-colors">
                <span class="sr-only">Facebook</span>
                <svg class="h-6 w-6" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M22 12c0-5.523-4.477-10-10-10S2 6.477 2 12c0 4.991 3.657 9.128 8.438 9.878v-6.987h-2.54V12h2.54V9.797c0-2.506 1.492-3.89 3.777-3.89 1.094 0 2.238.195 2.238.195v2.46h-1.26c-1.243 0-1.63.771-1.63 1.562V12h2.773l-.443 2.89h-2.33v6.988C18.343 21.128 22 16.991 22 12z"/>
                </svg>
              </a>
              <a href="#" class="text-gray-400 hover:text-green-400 transition-colors">
                <span class="sr-only">Instagram</span>
                <svg class="h-6 w-6" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M12 2.163c3.204 0 3.584.012 4.85.07 3.252.148 4.771 1.691 4.919 4.919.058 1.265.069 1.645.069 4.849 0 3.205-.012 3.584-.069 4.849-.149 3.225-1.664 4.771-4.919 4.919-1.266.058-1.644.07-4.85.07-3.204 0-3.584-.012-4.849-.07-3.26-.149-4.771-1.699-4.919-4.92-.058-1.265-.07-1.644-.07-4.849 0-3.204.013-3.583.07-4.849.149-3.227 1.664-4.771 4.919-4.919 1.266-.057 1.645-.069 4.849-.069zM12 0C8.741 0 8.333.014 7.053.072 2.695.272.273 2.69.073 7.052.014 8.333 0 8.741 0 12c0 3.259.014 3.668.072 4.948.2 4.358 2.618 6.78 6.98 6.98C8.333 23.986 8.741 24 12 24c3.259 0 3.668-.014 4.948-.072 4.354-.2 6.782-2.618 6.979-6.98.059-1.28.073-1.689.073-4.948 0-3.259-.014-3.667-.072-4.947-.196-4.354-2.617-6.78-6.979-6.98C15.668.014 15.259 0 12 0zm0 5.838a6.162 6.162 0 100 12.324 6.162 6.162 0 000-12.324zM12 16a4 4 0 110-8 4 4 0 010 8zm6.406-11.845a1.44 1.44 0 100 2.881 1.44 1.44 0 000-2.881z"/>
                </svg>
              </a>
              <a href="#" class="text-gray-400 hover:text-green-400 transition-colors">
                <span class="sr-only">Pinterest</span>
                <svg class="h-6 w-6" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M12 0C5.373 0 0 5.372 0 12c0 5.084 3.163 9.426 7.627 11.174-.105-.949-.2-2.405.042-3.441.218-.937 1.407-5.965 1.407-5.965s-.359-.719-.359-1.782c0-1.668.967-2.914 2.171-2.914 1.023 0 1.518.769 1.518 1.69 0 1.029-.655 2.568-.994 3.995-.283 1.194.599 2.169 1.777 2.169 2.133 0 3.772-2.249 3.772-5.495 0-2.873-2.064-4.882-5.012-4.882-3.414 0-5.418 2.561-5.418 5.207 0 1.031.397 2.138.893 2.738a.36.36 0 01.083.345l-.333 1.36c-.053.22-.174.267-.402.161-1.499-.698-2.436-2.889-2.436-4.649 0-3.785 2.75-7.262 7.929-7.262 4.163 0 7.398 2.967 7.398 6.931 0 4.136-2.607 7.464-6.227 7.464-1.216 0-2.359-.631-2.75-1.378l-.748 2.853c-.271 1.043-1.002 2.35-1.492 3.146C9.57 23.812 10.763 24 12 24c6.627 0 12-5.373 12-12 0-6.628-5.373-12-12-12z"/>
                </svg>
              </a>
            </div>
          </div>
  
          <!-- Quick Links -->
          <div class="space-y-4">
            <h3 class="text-lg font-semibold text-green-400">Quick Links</h3>
            <ul class="space-y-2">
              <li><NuxtLink to="/recipes" class="text-gray-400 hover:text-green-400 transition-colors">Browse Recipes</NuxtLink></li>
              <li><NuxtLink to="/categories" class="text-gray-400 hover:text-green-400 transition-colors">Categories</NuxtLink></li>
              <li><NuxtLink to="/popular" class="text-gray-400 hover:text-green-400 transition-colors">Popular Recipes</NuxtLink></li>
              <li><NuxtLink to="/seasonal" class="text-gray-400 hover:text-green-400 transition-colors">Seasonal Picks</NuxtLink></li>
            </ul>
          </div>
  
          <!-- Help & Legal -->
          <div class="space-y-4">
            <h3 class="text-lg font-semibold text-green-400">Help & Legal</h3>
            <ul class="space-y-2">
              <li><NuxtLink to="/faq" class="text-gray-400 hover:text-green-400 transition-colors">FAQ</NuxtLink></li>
              <li><NuxtLink to="/contact" class="text-gray-400 hover:text-green-400 transition-colors">Contact Us</NuxtLink></li>
              <li><NuxtLink to="/privacy" class="text-gray-400 hover:text-green-400 transition-colors">Privacy Policy</NuxtLink></li>
              <li><NuxtLink to="/terms" class="text-gray-400 hover:text-green-400 transition-colors">Terms of Service</NuxtLink></li>
            </ul>
          </div>
  
          <!-- Newsletter -->
          <div class="space-y-4">
            <h3 class="text-lg font-semibold text-green-400">Stay Updated</h3>
            <p class="text-gray-400">Subscribe to our newsletter for new recipes and tips.</p>
            <form class="flex flex-col space-y-3">
              <input 
                type="email" 
                placeholder="Your email" 
                class="px-4 py-2 rounded bg-gray-800 border border-gray-700 focus:outline-none focus:ring-1 focus:ring-green-400"
                required
              >
              <button 
                type="submit" 
                class="px-4 py-2 bg-green-600 hover:bg-green-700 rounded text-white transition-colors"
              >
                Subscribe
              </button>
            </form>
          </div>
        </div>
  
        <!-- Copyright -->
        <div class="mt-16 pt-8 border-t border-gray-800 text-center text-gray-500 text-sm">
          <p>© {{ new Date().getFullYear() }} Kushna. All rights reserved.</p>
          <p class="mt-2">Made with ❤️ for food lovers worldwide</p>
        </div>
      </div>
    </Footer>
  </template>
  
  <script setup>
  // No script needed for basic footer
  </script>