-- Drop and recreate foreign key constraints with ON DELETE CASCADE for all tables referencing recipes.id

-- bookmarks table
ALTER TABLE bookmarks
DROP CONSTRAINT IF EXISTS bookmarks_recipe_id_fkey;
ALTER TABLE bookmarks
ADD CONSTRAINT bookmarks_recipe_id_fkey
FOREIGN KEY (recipe_id)
REFERENCES recipes(id)
ON DELETE CASCADE;

-- comments table
ALTER TABLE comments
DROP CONSTRAINT IF EXISTS comments_recipe_id_fkey;
ALTER TABLE comments
ADD CONSTRAINT comments_recipe_id_fkey
FOREIGN KEY (recipe_id)
REFERENCES recipes(id)
ON DELETE CASCADE;

-- likes table
ALTER TABLE likes
DROP CONSTRAINT IF EXISTS likes_recipe_id_fkey;
ALTER TABLE likes
ADD CONSTRAINT likes_recipe_id_fkey
FOREIGN KEY (recipe_id)
REFERENCES recipes(id)
ON DELETE CASCADE;

-- purchases table
ALTER TABLE purchases
DROP CONSTRAINT IF EXISTS purchases_recipe_id_fkey;
ALTER TABLE purchases
ADD CONSTRAINT purchases_recipe_id_fkey
FOREIGN KEY (recipe_id)
REFERENCES recipes(id)
ON DELETE CASCADE;

-- ratings table
ALTER TABLE ratings
DROP CONSTRAINT IF EXISTS ratings_recipe_id_fkey;
ALTER TABLE ratings
ADD CONSTRAINT ratings_recipe_id_fkey
FOREIGN KEY (recipe_id)
REFERENCES recipes(id)
ON DELETE CASCADE;

-- recipe_images table
ALTER TABLE recipe_images
DROP CONSTRAINT IF EXISTS recipe_images_recipe_id_fkey;
ALTER TABLE recipe_images
ADD CONSTRAINT recipe_images_recipe_id_fkey
FOREIGN KEY (recipe_id)
REFERENCES recipes(id)
ON DELETE CASCADE;

-- recipe_ingredients table
ALTER TABLE recipe_ingredients
DROP CONSTRAINT IF EXISTS recipe_ingredients_recipe_id_fkey;
ALTER TABLE recipe_ingredients
ADD CONSTRAINT recipe_ingredients_recipe_id_fkey
FOREIGN KEY (recipe_id)
REFERENCES recipes(id)
ON DELETE CASCADE;

-- recipe_steps table
ALTER TABLE recipe_steps
DROP CONSTRAINT IF EXISTS recipe_steps_recipe_id_fkey;
ALTER TABLE recipe_steps
ADD CONSTRAINT recipe_steps_recipe_id_fkey
FOREIGN KEY (recipe_id)
REFERENCES recipes(id)
ON DELETE CASCADE;
-- Migration: Add cascade delete to recipe_ingredients.recipe_id foreign key
ALTER TABLE recipe_ingredients
DROP CONSTRAINT recipe_ingredients_recipe_id_fkey,
ADD CONSTRAINT recipe_ingredients_recipe_id_fkey
  FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE;
