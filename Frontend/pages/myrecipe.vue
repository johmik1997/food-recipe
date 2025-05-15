<template>
    <div>
      <div class="flex flex-col justify-between mb-6 sm:flex-row sm:items-center">
        <div>
          <h1 class="text-2xl font-bold text-gray-900">My Recipes</h1>
          <p class="mt-1 text-gray-600">All recipes you've shared with the community</p>
        </div>
        <NuxtLink
          to="/dashboard/recipes/add"
          class="inline-flex items-center px-4 py-2 mt-4 text-sm font-medium text-white bg-primary-600 border border-transparent rounded-md shadow-sm sm:mt-0 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
        >
          <Icon name="mdi:plus" class="w-5 h-5 mr-2" />
          Add New Recipe
        </NuxtLink>
      </div>
      
      <div class="flex items-center mb-4 space-x-4">
        <div class="relative flex-1">
          <div class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none">
            <Icon name="mdi:magnify" class="w-5 h-5 text-gray-400" />
          </div>
          <input
            type="text"
            v-model="searchQuery"
            class="block w-full py-2 pl-10 pr-3 text-sm bg-white border border-gray-300 rounded-md focus:ring-primary-500 focus:border-primary-500"
            placeholder="Search your recipes..."
          />
        </div>
        <div class="w-48">
          <select
            v-model="selectedCategory"
            class="block w-full py-2 pl-3 pr-10 text-sm bg-white border border-gray-300 rounded-md focus:ring-primary-500 focus:border-primary-500"
          >
            <option value="">All Categories</option>
            <option v-for="category in categories" :key="category.id" :value="category.id">
              {{ category.name }}
            </option>
          </select>
        </div>
      </div>
      
      <div v-if="loading" class="flex justify-center py-12">
        <div class="w-8 h-8 border-4 border-primary-500 border-t-transparent rounded-full animate-spin"></div>
      </div>
      
      <div v-else-if="filteredRecipes.length === 0" class="py-12 text-center">
        <Icon name="mdi:book-off" class="w-12 h-12 mx-auto text-gray-400" />
        <h3 class="mt-2 text-sm font-medium text-gray-900">No recipes found</h3>
        <p class="mt-1 text-sm text-gray-500">
          Get started by creating a new recipe.
        </p>
        <div class="mt-6">
          <NuxtLink
            to="/dashboard/recipes/add"
            class="inline-flex items-center px-4 py-2 text-sm font-medium text-white bg-primary-600 border border-transparent rounded-md shadow-sm hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
          >
            <Icon name="mdi:plus" class="w-5 h-5 mr-2" />
            New Recipe
          </NuxtLink>
        </div>
      </div>
      
      <div v-else class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
        <div
          v-for="recipe in filteredRecipes"
          :key="recipe.id"
          class="overflow-hidden bg-white rounded-lg shadow"
        >
          <div class="relative">
            <img
              :src="recipe.featured_image || '/placeholder-recipe.jpg'"
              :alt="recipe.title"
              class="object-cover w-full h-48"
            />
            <div class="absolute top-2 right-2">
              <span
                class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                :class="getCategoryClass(recipe.category_id)"
              >
                {{ getCategoryName(recipe.category_id) }}
              </span>
            </div>
          </div>
          <div class="p-4">
            <div class="flex items-start justify-between">
              <div>
                <h3 class="text-lg font-medium text-gray-900">
                  {{ recipe.title }}
                </h3>
                <p class="mt-1 text-sm text-gray-500">
                  {{ formatTime(recipe.prep_time) }} • {{ recipe.servings }} servings
                </p>
              </div>
              <div class="flex items-center ml-2">
                <Icon name="mdi:heart" class="w-5 h-5 text-red-500" />
                <span class="ml-1 text-sm text-gray-500">{{ recipe.likes_count }}</span>
              </div>
            </div>
            <p class="mt-2 text-sm text-gray-600 line-clamp-2">
              {{ recipe.description }}
            </p>
            <div class="flex justify-between mt-4">
              <div class="flex space-x-2">
                <NuxtLink
                  :to="`/dashboard/recipes/edit/${recipe.id}`"
                  class="inline-flex items-center px-3 py-1 text-sm font-medium text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200"
                >
                  <Icon name="mdi:pencil" class="w-4 h-4 mr-1" />
                  Edit
                </NuxtLink>
                <button
                  @click="confirmDelete(recipe.id)"
                  class="inline-flex items-center px-3 py-1 text-sm font-medium text-red-600 bg-red-100 rounded-md hover:bg-red-200"
                >
                  <Icon name="mdi:trash-can" class="w-4 h-4 mr-1" />
                  Delete
                </button>
              </div>
              <NuxtLink
                :to="`/recipes/${recipe.slug}`"
                class="inline-flex items-center px-3 py-1 text-sm font-medium text-white bg-primary-600 rounded-md hover:bg-primary-700"
              >
                View
              </NuxtLink>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Pagination -->
      <div v-if="filteredRecipes.length > 0" class="flex items-center justify-between mt-6">
        <div class="text-sm text-gray-700">
          Showing <span class="font-medium">{{ startItem }}</span> to <span class="font-medium">{{ endItem }}</span> of <span class="font-medium">{{ totalRecipes }}</span> recipes
        </div>
        <div class="flex space-x-2">
          <button
            @click="prevPage"
            :disabled="currentPage === 1"
            class="px-3 py-1 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Previous
          </button>
          <button
            @click="nextPage"
            :disabled="currentPage * itemsPerPage >= totalRecipes"
            class="px-3 py-1 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Next
          </button>
        </div>
      </div>
      
      <!-- Delete Confirmation Modal -->
      <div v-if="showDeleteModal" class="fixed inset-0 z-50 overflow-y-auto">
        <div class="flex items-end justify-center min-h-screen px-4 pt-4 pb-20 text-center sm:block sm:p-0">
          <div class="fixed inset-0 transition-opacity" aria-hidden="true">
            <div class="absolute inset-0 bg-gray-500 opacity-75"></div>
          </div>
          <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>
          <div class="inline-block overflow-hidden text-left align-bottom transition-all transform bg-white rounded-lg shadow-xl sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
            <div class="px-4 pt-5 pb-4 bg-white sm:p-6 sm:pb-4">
              <div class="sm:flex sm:items-start">
                <div class="flex items-center justify-center flex-shrink-0 w-12 h-12 mx-auto bg-red-100 rounded-full sm:mx-0 sm:h-10 sm:w-10">
                  <Icon name="mdi:alert-circle" class="w-6 h-6 text-red-600" />
                </div>
                <div class="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left">
                  <h3 class="text-lg font-medium leading-6 text-gray-900">Delete recipe</h3>
                  <div class="mt-2">
                    <p class="text-sm text-gray-500">Are you sure you want to delete this recipe? This action cannot be undone.</p>
                  </div>
                </div>
              </div>
            </div>
            <div class="px-4 py-3 bg-gray-50 sm:px-6 sm:flex sm:flex-row-reverse">
              <button
                type="button"
                @click="deleteRecipe"
                class="inline-flex justify-center w-full px-4 py-2 text-base font-medium text-white bg-red-600 border border-transparent rounded-md shadow-sm hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 sm:ml-3 sm:w-auto sm:text-sm"
              >
                Delete
              </button>
              <button
                type="button"
                @click="showDeleteModal = false"
                class="inline-flex justify-center w-full px-4 py-2 mt-3 text-base font-medium text-gray-700 bg-white border border-gray-300 rounded-md shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
              >
                Cancel
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </template>
  
  <script setup>
  const categories = ref([
    { id: 1, name: 'Breakfast', color: 'bg-blue-100 text-blue-800' },
    { id: 2, name: 'Lunch', color: 'bg-green-100 text-green-800' },
    { id: 3, name: 'Dinner', color: 'bg-purple-100 text-purple-800' },
    { id: 4, name: 'Dessert', color: 'bg-pink-100 text-pink-800' },
    { id: 5, name: 'Appetizer', color: 'bg-yellow-100 text-yellow-800' }
  ]);
  
  const loading = ref(true);
  const recipes = ref([]);
  const searchQuery = ref('');
  const selectedCategory = ref('');
  const showDeleteModal = ref(false);
  const recipeToDelete = ref(null);
  const currentPage = ref(1);
  const itemsPerPage = 9;
  
  // Mock data - replace with actual data fetching
  onMounted(() => {
    setTimeout(() => {
      recipes.value = [
        {
          id: 1,
          title: 'Homemade Pizza',
          description: 'Delicious homemade pizza with fresh ingredients and homemade dough.',
          prep_time: 45,
          servings: 4,
          featured_image: '/pizza.jpg',
          category_id: 3,
          likes_count: 24,
          slug: 'homemade-pizza'
        },
        {
          id: 2,
          title: 'Chocolate Chip Cookies',
          description: 'Classic chocolate chip cookies that are soft in the middle and crispy on the edges.',
          prep_time: 30,
          servings: 24,
          featured_image: '/cookies.jpg',
          category_id: 4,
          likes_count: 56,
          slug: 'chocolate-chip-cookies'
        },
        {
          id: 3,
          title: 'Avocado Toast',
          description: 'Simple yet delicious avocado toast with a sprinkle of chili flakes.',
          prep_time: 10,
          servings: 2,
          featured_image: '/avocado-toast.jpg',
          category_id: 1,
          likes_count: 18,
          slug: 'avocado-toast'
        },
        {
          id: 4,
          title: 'Vegetable Stir Fry',
          description: 'Quick and healthy vegetable stir fry with a savory sauce.',
          prep_time: 20,
          servings: 2,
          featured_image: '/stir-fry.jpg',
          category_id: 2,
          likes_count: 32,
          slug: 'vegetable-stir-fry'
        },
        {
          id: 5,
          title: 'Pasta Carbonara',
          description: 'Classic Italian pasta dish with eggs, cheese, pancetta, and pepper.',
          prep_time: 25,
          servings: 4,
          featured_image: '/carbonara.jpg',
          category_id: 3,
          likes_count: 42,
          slug: 'pasta-carbonara'
        },
        {
          id: 6,
          title: 'Berry Smoothie',
          description: 'Refreshing smoothie with mixed berries, banana, and yogurt.',
          prep_time: 5,
          servings: 1,
          featured_image: '/smoothie.jpg',
          category_id: 1,
          likes_count: 15,
          slug: 'berry-smoothie'
        }
      ];
      loading.value = false;
    }, 1000);
  });
  
  const filteredRecipes = computed(() => {
    let filtered = recipes.value;
    
    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase();
      filtered = filtered.filter(recipe => 
        recipe.title.toLowerCase().includes(query) || 
        recipe.description.toLowerCase().includes(query)
      );
    }
    
    if (selectedCategory.value) {
      filtered = filtered.filter(recipe => 
        recipe.category_id === parseInt(selectedCategory.value)
      );
    }
    
    // Apply pagination
    const start = (currentPage.value - 1) * itemsPerPage;
    const end = start + itemsPerPage;
    return filtered.slice(start, end);
  });
  
  const totalRecipes = computed(() => {
    let filtered = recipes.value;
    
    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase();
      filtered = filtered.filter(recipe => 
        recipe.title.toLowerCase().includes(query) || 
        recipe.description.toLowerCase().includes(query)
      );
    }
    
    if (selectedCategory.value) {
      filtered = filtered.filter(recipe => 
        recipe.category_id === parseInt(selectedCategory.value)
      );
    }
    
    return filtered.length;
  });
  
  const startItem = computed(() => {
    return (currentPage.value - 1) * itemsPerPage + 1;
  });
  
  const endItem = computed(() => {
    const end = currentPage.value * itemsPerPage;
    return end > totalRecipes.value ? totalRecipes.value : end;
  });
  
  const getCategoryName = (categoryId) => {
    const category = categories.value.find(c => c.id === categoryId);
    return category ? category.name : 'Uncategorized';
  };
  
  const getCategoryClass = (categoryId) => {
    const category = categories.value.find(c => c.id === categoryId);
    return category ? category.color : 'bg-gray-100 text-gray-800';
  };
  
  const formatTime = (minutes) => {
    if (minutes < 60) return `${minutes} min`;
    const hours = Math.floor(minutes / 60);
    const mins = minutes % 60;
    return `${hours}h ${mins > 0 ? `${mins}m` : ''}`;
  };
  
  const confirmDelete = (recipeId) => {
    recipeToDelete.value = recipeId;
    showDeleteModal.value = true;
  };
  
  const deleteRecipe = async () => {
    try {
      // Delete recipe from backend
      // await useMutation(DELETE_RECIPE_MUTATION, {
      //   variables: { id: recipeToDelete.value }
      // });
      
      // Remove from local state
      recipes.value = recipes.value.filter(recipe => recipe.id !== recipeToDelete.value);
      showDeleteModal.value = false;
    } catch (error) {
      console.error('Error deleting recipe:', error);
    }
  };
  
  const nextPage = () => {
    if (currentPage.value * itemsPerPage < totalRecipes.value) {
      currentPage.value++;
    }
  };
  
  const prevPage = () => {
    if (currentPage.value > 1) {
      currentPage.value--;
    }
  };
  </script>








-- 1. First create users table
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  email TEXT UNIQUE NOT NULL,
  password TEXT NOT NULL,
  avatar_image_url TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 2. Create categories table
CREATE TABLE categories (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL UNIQUE,
  image_url TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 3. Create recipes table
CREATE TABLE recipes (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  title TEXT NOT NULL,
  description TEXT,
  prep_time INT, -- in minutes
  cook_time INT, -- in minutes
  total_time INT GENERATED ALWAYS AS (prep_time + cook_time) STORED,
  servings INT,
  featured_image_url TEXT,
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 4. Create recipe_images table (must exist before creating its constraints)
CREATE TABLE recipe_images (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  recipe_id UUID REFERENCES recipes(id) ON DELETE CASCADE,
  image_url TEXT NOT NULL,
  is_featured BOOLEAN NOT NULL DEFAULT false,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE TABLE recipe_categories (
  recipe_id UUID REFERENCES recipes(id) ON DELETE CASCADE,
  category_id UUID REFERENCES categories(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (recipe_id, category_id)
);

-- Then create indexes for the junction table
CREATE INDEX idx_recipe_categories_recipe_id ON recipe_categories(recipe_id);
CREATE INDEX idx_recipe_categories_category_id ON recipe_categories(category_id);

CREATE UNIQUE INDEX one_featured_image_per_recipe
ON recipe_images(recipe_id)
WHERE (is_featured = true);

-- 6. Create remaining tables
CREATE TABLE recipe_steps (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  recipe_id UUID REFERENCES recipes(id) ON DELETE CASCADE,
  step_number INT NOT NULL,
  instruction TEXT NOT NULL,
  image_url TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(recipe_id, step_number)
);

CREATE TABLE recipe_ingredients (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  recipe_id UUID REFERENCES recipes(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  quantity DECIMAL(10,2),
  unit TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE user_likes (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  recipe_id UUID REFERENCES recipes(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(user_id, recipe_id)
);

CREATE TABLE user_bookmarks (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  recipe_id UUID REFERENCES recipes(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(user_id, recipe_id)
);

CREATE TABLE comments (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  recipe_id UUID REFERENCES recipes(id) ON DELETE CASCADE,
  content TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE ratings (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  recipe_id UUID REFERENCES recipes(id) ON DELETE CASCADE,
  value INT NOT NULL CHECK (value BETWEEN 1 AND 5),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(user_id, recipe_id)
);

CREATE TABLE transactions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE SET NULL,
  recipe_id UUID REFERENCES recipes(id) ON DELETE SET NULL,
  amount DECIMAL(10,2) NOT NULL,
  status TEXT NOT NULL,
  chapa_reference TEXT UNIQUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_recipes_user_id ON recipes(user_id);
CREATE INDEX idx_recipe_steps_recipe_id ON recipe_steps(recipe_id);
CREATE INDEX idx_user_likes_recipe_id ON user_likes(recipe_id);
CREATE INDEX idx_recipe_ingredients_recipe_id ON recipe_ingredients(recipe_id);
CREATE INDEX idx_comments_recipe_id ON comments(recipe_id);
CREATE INDEX idx_comments_user_id ON comments(user_id);
CREATE INDEX idx_ratings_recipe_id ON ratings(recipe_id);
CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE INDEX idx_transactions_recipe_id ON transactions(recipe_id);


CREATE INDEX idx_recipe_categories_recipe_id ON recipe_categories(recipe_id);
CREATE INDEX idx_recipe_categories_category_id ON recipe_categories(category_id);
CREATE UNIQUE INDEX one_featured_image_per_recipe
ON recipe_images(recipe_id)
WHERE (is_featured = true);
ALTER TABLE recipes ADD COLUMN trending_score DECIMAL(10,2) DEFAULT 0;
CREATE INDEX idx_recipes_trending_score ON recipes(trending_score DESC);

-- Create timestamp update function and triggers
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_timestamp
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER update_recipes_timestamp
BEFORE UPDATE ON recipes
FOR EACH ROW EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER update_comments_timestamp
BEFORE UPDATE ON comments
FOR EACH ROW EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER update_ratings_timestamp
BEFORE UPDATE ON ratings
FOR EACH ROW EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER update_transactions_timestamp
BEFORE UPDATE ON transactions
FOR EACH ROW EXECUTE FUNCTION update_timestamp();