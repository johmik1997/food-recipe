CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  email TEXT UNIQUE NOT NULL,
  password TEXT NOT NULL,
  avatar_url TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE categories (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL UNIQUE,
  image_url TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);   
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
  category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

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
);CREATE TABLE user_bookmarks (
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

CREATE INDEX idx_recipes_user_id ON recipes(user_id);
CREATE INDEX idx_recipes_category_id ON recipes(category_id);
CREATE INDEX idx_recipe_steps_recipe_id ON recipe_steps(recipe_id);
CREATE INDEX idx_user_likes_recipe_id ON user_likes(recipe_id);
-- 1. First create the function (only needed once)
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 2. Create triggers for all tables
-- Users table
CREATE TRIGGER update_users_timestamp
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION update_timestamp();

-- Recipes table
CREATE TRIGGER update_recipes_timestamp
BEFORE UPDATE ON recipes
FOR EACH ROW EXECUTE FUNCTION update_timestamp();

-- Comments table
CREATE TRIGGER update_comments_timestamp
BEFORE UPDATE ON comments
FOR EACH ROW EXECUTE FUNCTION update_timestamp();

-- Ratings table
CREATE TRIGGER update_ratings_timestamp
BEFORE UPDATE ON ratings
FOR EACH ROW EXECUTE FUNCTION update_timestamp();

-- Transactions table (if it has updated_at)
CREATE TRIGGER update_transactions_timestamp
BEFORE UPDATE ON transactions
FOR EACH ROW EXECUTE FUNCTION update_timestamp();  


ALTER TABLE users ADD COLUMN IF NOT EXISTS avatat_image_url TEXT;
 
CREATE UNIQUE INDEX one_featured_image_per_recipe
ON recipe_images(recipe_id)
WHERE (is_featured = true);
Create a materialized view for trending recipes
CREATE MATERIALIZED VIEW trending_recipes AS
SELECT 
  r.id AS recipe_id,
  r.title,
  r.featured_image_url,
  -- Calculate trending score based on:
  -- 40% likes, 30% bookmarks, 20% comments, 10% recency
  (
    (COALESCE((SELECT COUNT(*) FROM user_likes ul WHERE ul.recipe_id = r.id), 0) * 0.4) +
    (COALESCE((SELECT COUNT(*) FROM user_bookmarks ub WHERE ub.recipe_id = r.id), 0) * 0.3) +
    (COALESCE((SELECT COUNT(*) FROM comments c WHERE c.recipe_id = r.id), 0) * 0.2) +
    (1.0 / (1.0 + EXTRACT(EPOCH FROM (NOW() - r.created_at)/86400)) * 0.1
  ) AS trending_score
FROM recipes r
ORDER BY trending_score DESC
LIMIT 50;

-- Create a function to refresh the view
CREATE OR REPLACE FUNCTION refresh_trending_recipes()
RETURNS TRIGGER AS $$
BEGIN
  REFRESH MATERIALIZED VIEW CONCURRENTLY trending_recipes;
  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Set up triggers to refresh when relevant data changes
CREATE TRIGGER refresh_trending_after_like
AFTER INSERT OR UPDATE OR DELETE ON user_likes
FOR EACH STATEMENT EXECUTE FUNCTION refresh_trending_recipes();

CREATE TRIGGER refresh_trending_after_bookmark
AFTER INSERT OR UPDATE OR DELETE ON user_bookmarks
FOR EACH STATEMENT EXECUTE FUNCTION refresh_trending_recipes();

CREATE TRIGGER refresh_trending_after_comment
AFTER INSERT OR UPDATE OR DELETE ON comments
FOR EACH STATEMENT EXECUTE FUNCTION refresh_trending_recipes();

CREATE TRIGGER refresh_trending_after_recipe_change
AFTER INSERT OR UPDATE OR DELETE ON recipes
FOR EACH STATEMENT EXECUTE FUNCTION refresh_trending_recipes();