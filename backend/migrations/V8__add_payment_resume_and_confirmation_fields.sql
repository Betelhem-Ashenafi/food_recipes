ALTER TABLE IF EXISTS purchases
ADD COLUMN IF NOT EXISTS checkout_url TEXT;


CREATE INDEX IF NOT EXISTS idx_purchases_user_recipe_status
ON purchases(user_id, recipe_id, status, created_at DESC);
