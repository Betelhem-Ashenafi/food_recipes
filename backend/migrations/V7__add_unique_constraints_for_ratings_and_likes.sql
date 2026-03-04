-- Remove duplicate likes per user/recipe pair
DELETE FROM likes a
USING likes b
WHERE a.recipe_id = b.recipe_id
  AND a.user_id = b.user_id
  AND a.ctid < b.ctid;

-- Remove duplicate ratings per user/recipe pair
DELETE FROM ratings a
USING ratings b
WHERE a.recipe_id = b.recipe_id
  AND a.user_id = b.user_id
  AND a.ctid < b.ctid;

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM pg_constraint
    WHERE conname = 'likes_recipe_user_unique'
  ) THEN
    ALTER TABLE likes
    ADD CONSTRAINT likes_recipe_user_unique UNIQUE (recipe_id, user_id);
  END IF;
END $$;

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM pg_constraint
    WHERE conname = 'ratings_recipe_user_unique'
  ) THEN
    ALTER TABLE ratings
    ADD CONSTRAINT ratings_recipe_user_unique UNIQUE (recipe_id, user_id);
  END IF;
END $$;
