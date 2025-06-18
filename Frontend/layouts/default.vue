
<template>
  <div class="min-h-screen flex flex-col">
    <!-- Navbar at the top -->
    <UINavbar />

    <!-- Content in the middle -->
    <main class="flex-grow">
      <slot></slot>
    </main>

    <!-- Footer at the bottom -->
    <UIFooter />
  </div>
</template>