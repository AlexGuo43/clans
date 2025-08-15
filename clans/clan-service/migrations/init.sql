CREATE TABLE IF NOT EXISTS clans (
    id SERIAL PRIMARY KEY,
    name VARCHAR(20) UNIQUE NOT NULL,
    display_name VARCHAR(50) NOT NULL,
    description TEXT,
    owner_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    is_public BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS clan_memberships (
    id SERIAL PRIMARY KEY,
    clan_id INTEGER NOT NULL REFERENCES clans(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL DEFAULT 'member',
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(clan_id, user_id)
);

CREATE INDEX idx_clans_name ON clans(name);
CREATE INDEX idx_clans_owner_id ON clans(owner_id);
CREATE INDEX idx_clans_is_public ON clans(is_public);
CREATE INDEX idx_clans_created_at ON clans(created_at);

CREATE INDEX idx_clan_memberships_clan_id ON clan_memberships(clan_id);
CREATE INDEX idx_clan_memberships_user_id ON clan_memberships(user_id);
CREATE INDEX idx_clan_memberships_role ON clan_memberships(role);
CREATE INDEX idx_clan_memberships_joined_at ON clan_memberships(joined_at);