-- Function to calculate average rating for a recipe
CREATE OR REPLACE FUNCTION recipe_average_rating(recipe_row recipes)
RETURNS FLOAT AS $$
  SELECT COALESCE(AVG(rating), 0)
  FROM ratings
  WHERE recipe_id = recipe_row.id;
$$ LANGUAGE sql STABLE;

-- Function to count likes for a recipe
CREATE OR REPLACE FUNCTION recipe_likes_count(recipe_row recipes)
RETURNS BIGINT AS $$
  SELECT COUNT(*)
  FROM likes
  WHERE recipe_id = recipe_row.id;
$$ LANGUAGE sql STABLE;
