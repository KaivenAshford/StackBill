-- Remove updated_at from reminder tracking tables

BEGIN;

ALTER TABLE reminder_reads DROP COLUMN IF EXISTS updated_at;
ALTER TABLE reminder_dismisseds DROP COLUMN IF EXISTS updated_at;

COMMIT;
