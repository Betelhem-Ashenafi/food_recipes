-- V10: Add Hasura-trackable wrapper for recipe access checks.
-- Keep can_user_access_recipe_content(...) as the source-of-truth boolean rule.
CREATE OR REPLACE FUNCTION get_accessible_recipe(
    p_user_id INT,
    p_recipe_id INT
)
RETURNS SETOF recipes
LANGUAGE sql
STABLE
AS $$
    SELECT r.*
    FROM recipes r
    WHERE r.id = p_recipe_id
      AND can_user_access_recipe_content(p_user_id, p_recipe_id);
$$;
