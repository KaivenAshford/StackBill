-- Rollback StackBill Initial Schema
-- Drops all tables from v1.0.0

BEGIN;

DROP TABLE IF EXISTS reminder_dismisseds;
DROP TABLE IF EXISTS reminder_reads;
DROP TABLE IF EXISTS assets;
DROP TABLE IF EXISTS subscriptions;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS users;

COMMIT;
