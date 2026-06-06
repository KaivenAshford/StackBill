-- Add subscription_id to assets table for asset-subscription association

BEGIN;

ALTER TABLE assets ADD COLUMN IF NOT EXISTS subscription_id BIGINT DEFAULT 0;
CREATE INDEX IF NOT EXISTS idx_assets_subscription_id ON assets (subscription_id);

COMMIT;
