-- V9: Add Postgres events, trigger-driven outbox, helper functions,
-- and generated properties for Hasura integration.

-- 1) Event outbox table (for Postgres -> Hasura event trigger webhook pattern)
CREATE TABLE IF NOT EXISTS payment_events (
    id BIGSERIAL PRIMARY KEY,
    purchase_id INT NOT NULL REFERENCES purchases(id) ON DELETE CASCADE,
    tx_ref VARCHAR(255) NOT NULL,
    user_id INT NOT NULL,
    recipe_id INT NOT NULL,
    old_status VARCHAR(32),
    new_status VARCHAR(32) NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_payment_events_purchase_id ON payment_events(purchase_id);
CREATE INDEX IF NOT EXISTS idx_payment_events_created_at ON payment_events(created_at DESC);

-- 2) Trigger function to capture purchase status transitions as events
CREATE OR REPLACE FUNCTION emit_payment_event()
RETURNS TRIGGER AS $$
DECLARE
    old_status_normalized TEXT;
    new_status_normalized TEXT;
BEGIN
    old_status_normalized := CASE WHEN TG_OP = 'INSERT' THEN NULL ELSE LOWER(COALESCE(OLD.status, '')) END;
    new_status_normalized := LOWER(COALESCE(NEW.status, ''));

    -- Emit event for insert, or when status actually changes
    IF TG_OP = 'INSERT' OR old_status_normalized IS DISTINCT FROM new_status_normalized THEN
        INSERT INTO payment_events (
            purchase_id,
            tx_ref,
            user_id,
            recipe_id,
            old_status,
            new_status,
            payload
        ) VALUES (
            NEW.id,
            COALESCE(NEW.chapa_tx_ref, ''),
            NEW.user_id,
            NEW.recipe_id,
            CASE WHEN TG_OP = 'INSERT' THEN NULL ELSE old_status_normalized END,
            new_status_normalized,
            jsonb_build_object(
                'purchase_id', NEW.id,
                'tx_ref', COALESCE(NEW.chapa_tx_ref, ''),
                'user_id', NEW.user_id,
                'recipe_id', NEW.recipe_id,
                'amount', NEW.amount,
                'currency', NEW.currency,
                'old_status', CASE WHEN TG_OP = 'INSERT' THEN NULL ELSE old_status_normalized END,
                'new_status', new_status_normalized,
                'event_time', CURRENT_TIMESTAMP
            )
        );

        -- Optional Postgres event stream for listeners
        PERFORM pg_notify(
            'payment_status_changed',
            json_build_object(
                'purchase_id', NEW.id,
                'tx_ref', COALESCE(NEW.chapa_tx_ref, ''),
                'recipe_id', NEW.recipe_id,
                'user_id', NEW.user_id,
                'status', new_status_normalized
            )::text
        );
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_emit_payment_event ON purchases;
DROP TRIGGER IF EXISTS trg_emit_payment_event ON purchases;
CREATE TRIGGER trg_emit_payment_event
AFTER INSERT OR UPDATE ON purchases
FOR EACH ROW
EXECUTE FUNCTION emit_payment_event();

-- 3) Postgres helper functions (usable via Hasura tracked SQL functions)
-- Returns whether a user can access a recipe's protected content.
CREATE OR REPLACE FUNCTION can_user_access_recipe_content(p_user_id INT, p_recipe_id INT)
RETURNS BOOLEAN AS $$
    SELECT EXISTS (
        SELECT 1
        FROM recipes r
        WHERE r.id = p_recipe_id
          AND (
              COALESCE(r.price, 0) <= 0
              OR r.user_id = p_user_id
              OR EXISTS (
                  SELECT 1
                  FROM purchases p
                  WHERE p.user_id = p_user_id
                    AND p.recipe_id = p_recipe_id
                    AND LOWER(COALESCE(p.status, '')) = 'success'
              )
          )
    );
$$ LANGUAGE sql STABLE;

-- Computed-field style function: successful purchase count per recipe
CREATE OR REPLACE FUNCTION recipe_successful_purchases_count(recipe_row recipes)
RETURNS BIGINT AS $$
    SELECT COUNT(*)
    FROM purchases p
    WHERE p.recipe_id = recipe_row.id
      AND LOWER(COALESCE(p.status, '')) = 'success';
$$ LANGUAGE sql STABLE;

-- 4) Generated properties (stored generated columns)
ALTER TABLE IF EXISTS recipes
ADD COLUMN IF NOT EXISTS is_paid BOOLEAN GENERATED ALWAYS AS (COALESCE(price, 0) > 0) STORED;

ALTER TABLE IF EXISTS purchases
ADD COLUMN IF NOT EXISTS is_success BOOLEAN GENERATED ALWAYS AS (LOWER(COALESCE(status, '')) = 'success') STORED;
