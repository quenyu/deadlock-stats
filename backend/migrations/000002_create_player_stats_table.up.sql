CREATE TABLE IF NOT EXISTS player_stats (
    user_id UUID PRIMARY KEY,
    kd_ratio REAL NOT NULL DEFAULT 0,
    win_rate REAL NOT NULL DEFAULT 0,
    avg_matches_per_day REAL NOT NULL DEFAULT 0,
    favorite_hero VARCHAR(255),
    last_updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);
