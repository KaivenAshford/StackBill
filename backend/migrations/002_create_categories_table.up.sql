CREATE TABLE categories (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    user_id BIGINT NOT NULL,
    name VARCHAR(50) NOT NULL,
    type VARCHAR(20) NOT NULL,
    color VARCHAR(20) DEFAULT '',
    icon VARCHAR(50) DEFAULT '',
    sort_order INT DEFAULT 0
);

CREATE INDEX idx_categories_user_id ON categories (user_id);
CREATE INDEX idx_categories_deleted_at ON categories (deleted_at);
