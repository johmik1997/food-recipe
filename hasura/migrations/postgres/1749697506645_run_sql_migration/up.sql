CREATE OR REPLACE FUNCTION public.set_current_timestamp_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 4. Create recipes table
CREATE TABLE public.recipes (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
  category_id INTEGER REFERENCES public.categories(id) ON DELETE SET NULL,

  title VARCHAR NOT NULL,
  description TEXT,
  prep_time INTEGER NOT NULL,
  cook_time INTEGER NOT NULL,
  total_time INTEGER GENERATED ALWAYS AS (prep_time + cook_time) STORED,

  servings INTEGER,
  average_rating NUMERIC DEFAULT 0,
  featured_image TEXT,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 5. Trigger to auto-update updated_at
CREATE TRIGGER set_public_recipes_updated_at
BEFORE UPDATE ON public.recipes
FOR EACH ROW
EXECUTE FUNCTION public.set_current_timestamp_updated_at();
