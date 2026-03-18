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

-- Vault invites table (pending invitations for users who don't have an account yet)
CREATE TABLE IF NOT EXISTS vault_invites (
    id          SERIAL PRIMARY KEY,
    vault_id    INTEGER NOT NULL REFERENCES vaults(id) ON DELETE CASCADE,
    invited_by  INTEGER NOT NULL REFERENCES users(id),
    email       TEXT NOT NULL,
    role        INTEGER NOT NULL CHECK (role BETWEEN 1 AND 4),
    path        VARCHAR(500) NOT NULL DEFAULT '/',
    token       TEXT NOT NULL UNIQUE,
    expires_at  TIMESTAMP NOT NULL,
    accepted_at TIMESTAMP,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_vault_invites_token ON vault_invites(token);
CREATE INDEX IF NOT EXISTS idx_vault_invites_email ON vault_invites(email);

-- User password reset requests
CREATE TABLE IF NOT EXISTS password_resets (
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token      TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    used_at    TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_password_resets_token ON password_resets(token);
