-- StackBill Initial Schema
-- Creates all tables for v1.0.0

BEGIN;

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password VARCHAR(255) NOT NULL,
    nickname VARCHAR(50) DEFAULT '',
    avatar VARCHAR(500) DEFAULT ''
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON users (username) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users (email) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users (deleted_at);

-- Categories table
CREATE TABLE IF NOT EXISTS categories (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    user_id BIGINT NOT NULL,
    name VARCHAR(50) NOT NULL,
    type VARCHAR(20) NOT NULL,
    color VARCHAR(20) DEFAULT '',
    icon VARCHAR(50) DEFAULT '',
    sort_order INT DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_categories_user_id ON categories (user_id);
CREATE INDEX IF NOT EXISTS idx_categories_deleted_at ON categories (deleted_at);

-- Subscriptions table
CREATE TABLE IF NOT EXISTS subscriptions (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    user_id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500) DEFAULT '',
    category_id BIGINT DEFAULT 0,
    amount DECIMAL(10,2) NOT NULL DEFAULT 0,
    currency VARCHAR(10) DEFAULT 'USD',
    billing_cycle VARCHAR(20) NOT NULL,
    billing_interval INT DEFAULT 1,
    next_payment_date TIMESTAMPTZ,
    start_date TIMESTAMPTZ,
    payment_method VARCHAR(50) DEFAULT '',
    auto_renew BOOLEAN DEFAULT TRUE,
    status VARCHAR(20) DEFAULT 'active',
    website_url VARCHAR(500) DEFAULT '',
    remark VARCHAR(500) DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id ON subscriptions (user_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_category_id ON subscriptions (category_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_deleted_at ON subscriptions (deleted_at);

-- Assets table
CREATE TABLE IF NOT EXISTS assets (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    user_id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL,
    asset_type VARCHAR(30) NOT NULL,
    provider VARCHAR(100) DEFAULT '',
    identifier VARCHAR(200) DEFAULT '',
    url VARCHAR(500) DEFAULT '',
    expire_date TIMESTAMPTZ,
    cost_amount DECIMAL(10,2) DEFAULT 0,
    cost_currency VARCHAR(10) DEFAULT 'USD',
    billing_cycle VARCHAR(20) DEFAULT '',
    status VARCHAR(20) DEFAULT 'active',
    description VARCHAR(500) DEFAULT '',
    remark VARCHAR(500) DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_assets_user_id ON assets (user_id);
CREATE INDEX IF NOT EXISTS idx_assets_deleted_at ON assets (deleted_at);

-- Reminder read tracking table
CREATE TABLE IF NOT EXISTS reminder_reads (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    user_id BIGINT NOT NULL,
    target_type VARCHAR(30) NOT NULL,
    target_id BIGINT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_reminder_read ON reminder_reads (user_id, target_type, target_id);

-- Reminder dismissed tracking table
CREATE TABLE IF NOT EXISTS reminder_dismisseds (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    user_id BIGINT NOT NULL,
    target_type VARCHAR(30) NOT NULL,
    target_id BIGINT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_reminder_dismissed ON reminder_dismisseds (user_id, target_type, target_id);

COMMIT;
