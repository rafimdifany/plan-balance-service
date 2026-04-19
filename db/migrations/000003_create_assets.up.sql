CREATE TABLE IF NOT EXISTS assets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('BANK', 'CASH', 'EWALLET')),
    icon VARCHAR(50) NOT NULL DEFAULT 'account_balance_wallet',
    color VARCHAR(7) NOT NULL DEFAULT '#6366F1',
    balance DECIMAL(15, 2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP DEFAULT NULL
);

-- Index for searching user's assets
CREATE INDEX IF NOT EXISTS idx_assets_user_id ON assets(user_id);

-- Unique constraint: a user cannot have two active assets with same name.
CREATE UNIQUE INDEX IF NOT EXISTS unique_active_asset_name_per_user 
ON assets (user_id, name) 
WHERE deleted_at IS NULL;
