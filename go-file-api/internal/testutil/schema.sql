-- Users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Vaults table
CREATE TABLE IF NOT EXISTS vaults (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Vault users table (junction table with roles and paths)
CREATE TABLE IF NOT EXISTS vault_users (
    id SERIAL PRIMARY KEY,
    vault_id INTEGER NOT NULL REFERENCES vaults(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    path VARCHAR(500) NOT NULL DEFAULT '/',
    role INTEGER NOT NULL CHECK (role BETWEEN 1 AND 4),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(vault_id, user_id, path)
);

-- Indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_vault_users_vault_id ON vault_users(vault_id);
CREATE INDEX IF NOT EXISTS idx_vault_users_user_id ON vault_users(user_id);
CREATE INDEX IF NOT EXISTS idx_vault_users_vault_user ON vault_users(vault_id, user_id);
