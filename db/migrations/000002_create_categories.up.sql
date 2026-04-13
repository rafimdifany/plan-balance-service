CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('INCOME', 'EXPENSE')),
    icon VARCHAR(50) NOT NULL DEFAULT 'default',
    color VARCHAR(7) NOT NULL DEFAULT '#6366F1',
    budget_limit DECIMAL(15, 2) DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP DEFAULT NULL
);

-- Index for searching user's categories
CREATE INDEX IF NOT EXISTS idx_categories_user_id ON categories(user_id);

-- Unique constraint: a user cannot have two active categories with same name and type.
-- We use a partial unique index to allow re-using the same name after a category is soft-deleted.
CREATE UNIQUE INDEX IF NOT EXISTS unique_active_category_name_per_user 
ON categories (user_id, name, type) 
WHERE deleted_at IS NULL;
