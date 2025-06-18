<template>
  <div class="sticky top-0 z-10 bg-white shadow-sm">
    <div class="container mx-auto px-4">
      <nav class="flex items-center justify-between h-16">
        <!-- Logo -->
        <div class="flex-1">
          <img src="/images/logo1.png" alt="logo" class="h-10" />
        </div>

        <!-- Search Bar -->
        <input
          type="text"
          placeholder="What are we cooking today?"
          class="input input-bordered w-[120px] md:w-[300px] rounded-full border-green-300 focus:border-green-500 focus:ring-green-200 mx-4"
        />

        <!-- Navigation Items -->
        <div v-if="user" class="flex items-center space-x-4">
          <!-- Recipes Dropdown -->
          <div class="relative">
            <Menu as="div" class="relative inline-block text-left">
              <div>
                <MenuButton
                  class="inline-flex w-full justify-center rounded-md bg-green-600 px-4 py-2 text-sm font-medium text-white hover:bg-green-700 focus:outline-none focus-visible:ring-2 focus-visible:ring-green-500 focus-visible:ring-opacity-75 transition-colors duration-200"
                >
                  Recipes
                  <ChevronDownIcon
                    class="-mr-1 ml-2 h-5 w-5 text-green-100"
                    aria-hidden="true"
                  />
                </MenuButton>
              </div>

              <transition
                enter-active-class="transition duration-100 ease-out"
                enter-from-class="transform scale-95 opacity-0"
                enter-to-class="transform scale-100 opacity-100"
                leave-active-class="transition duration-75 ease-in"
                leave-from-class="transform scale-100 opacity-100"
                leave-to-class="transform scale-95 opacity-0"
              >
                <MenuItems
                  class="absolute right-0 mt-2 w-56 origin-top-right divide-y divide-green-100 rounded-md bg-white shadow-lg ring-1 ring-green-200 focus:outline-none"
                >
                  <div class="px-1 py-1">
                    <MenuItem v-slot="{ active }">
                      <button
                        :class="[
                          active ? 'bg-green-100 text-green-800' : 'text-green-700',
                          'group flex w-full items-center rounded-md px-2 py-2 text-sm',
                        ]"
                      >
                        <PencilIcon class="mr-2 h-4 w-4 text-green-500" />
                        Edit Recipe
                      </button>
                    </MenuItem>
                    <MenuItem v-slot="{ active }">
                      <button
                        :class="[
                          active ? 'bg-green-100 text-green-800' : 'text-green-700',
                          'group flex w-full items-center rounded-md px-2 py-2 text-sm',
                        ]"
                      >
                        <DocumentDuplicateIcon class="mr-2 h-4 w-4 text-green-500" />
                        Duplicate
                      </button>
                    </MenuItem>
                  </div>
                  <div class="px-1 py-1">
                    <MenuItem v-slot="{ active }">
                      <button
                        :class="[
                          active ? 'bg-green-100 text-green-800' : 'text-green-700',
                          'group flex w-full items-center rounded-md px-2 py-2 text-sm',
                        ]"
                      >
                        <ArchiveIcon class="mr-2 h-4 w-4 text-green-500" />
                        Archive
                      </button>
                    </MenuItem>
                  </div>
                  <div class="px-1 py-1">
                    <MenuItem v-slot="{ active }">
                      <button
                        :class="[
                          active ? 'bg-green-100 text-green-800' : 'text-green-700',
                          'group flex w-full items-center rounded-md px-2 py-2 text-sm',
                        ]"
                      >
                        <TrashIcon class="mr-2 h-4 w-4 text-green-500" />
                        Delete
                      </button>
                    </MenuItem>
                  </div>
                </MenuItems>
              </transition>
            </Menu>
          </div>

          <!-- Add other navigation items here -->
        </div>
      </nav>
    </div>
  </div>
</template>

<style>
/* Custom styles for the navbar */
nav {
  transition: all 0.3s ease;
}

/* Input focus style */
.input:focus {
  box-shadow: 0 0 0 2px rgba(5, 150, 105, 0.2);
}

/* Dropdown transition */
.dropdown-enter-active,
.dropdown-leave-active {
  transition: all 0.2s ease;
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

/* Menu item hover effect */
.menu-item:hover {
  background-color: #f0fdf4;
}

/* Active link styling */
.router-link-active {
  @apply text-green-700 font-semibold border-b-2 border-green-500;
}
</style>