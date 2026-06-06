-- Remove subscription_id from assets table

BEGIN;

DROP INDEX IF EXISTS idx_assets_subscription_id;
ALTER TABLE assets DROP COLUMN IF EXISTS subscription_id;

COMMIT;
