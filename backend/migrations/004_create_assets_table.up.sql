CREATE TABLE assets (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    user_id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL,
    asset_type VARCHAR(30) NOT NULL,
    provider VARCHAR(100) DEFAULT '',
    identifier VARCHAR(200) DEFAULT '',
    url VARCHAR(500) DEFAULT '',
    expire_date TIMESTAMP WITH TIME ZONE,
    cost_amount DECIMAL(10,2) DEFAULT 0,
    cost_currency VARCHAR(10) DEFAULT 'USD',
    billing_cycle VARCHAR(20) DEFAULT '',
    status VARCHAR(20) DEFAULT 'active',
    description VARCHAR(500) DEFAULT '',
    remark VARCHAR(500) DEFAULT ''
);

CREATE INDEX idx_assets_user_id ON assets (user_id);
CREATE INDEX idx_assets_deleted_at ON assets (deleted_at);
