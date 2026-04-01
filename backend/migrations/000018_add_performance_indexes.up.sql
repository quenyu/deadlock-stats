-- Добавляем индексы для оптимизации производительности

-- Индексы для player_match_stats
CREATE INDEX IF NOT EXISTS idx_pms_user_id ON player_match_stats(user_id);
CREATE INDEX IF NOT EXISTS idx_pms_match_id ON player_match_stats(match_id);
CREATE INDEX IF NOT EXISTS idx_pms_hero_name ON player_match_stats(hero_name);
CREATE INDEX IF NOT EXISTS idx_pms_match_time ON player_match_stats(match_time DESC);

-- Индексы для matches
CREATE INDEX IF NOT EXISTS idx_matches_time ON matches(match_time DESC);
CREATE INDEX IF NOT EXISTS idx_matches_duration ON matches(duration);

-- Индексы для users (для поиска)
CREATE INDEX IF NOT EXISTS idx_users_nickname ON users USING gin(nickname gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_users_steam_id ON users(steam_id);

-- Индексы для crosshairs
CREATE INDEX IF NOT EXISTS idx_crosshairs_user_id ON crosshairs(user_id);
CREATE INDEX IF NOT EXISTS idx_crosshairs_created_at ON crosshairs(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_crosshairs_likes ON crosshairs(likes DESC);

-- Индексы для votes
CREATE INDEX IF NOT EXISTS idx_votes_user_id ON votes(user_id);
CREATE INDEX IF NOT EXISTS idx_votes_crosshair_id ON votes(crosshair_id);

-- Индексы для comments
CREATE INDEX IF NOT EXISTS idx_comments_crosshair_id ON comments(crosshair_id);
CREATE INDEX IF NOT EXISTS idx_comments_user_id ON comments(user_id);
CREATE INDEX IF NOT EXISTS idx_comments_created_at ON comments(created_at DESC);

-- Индексы для builds (если таблица существует)
CREATE INDEX IF NOT EXISTS idx_builds_user_id ON builds(user_id);
CREATE INDEX IF NOT EXISTS idx_builds_hero_name ON builds(hero_name);
CREATE INDEX IF NOT EXISTS idx_builds_created_at ON builds(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_builds_rating ON builds(rating DESC);

-- Составные индексы для частых запросов
CREATE INDEX IF NOT EXISTS idx_pms_user_hero ON player_match_stats(user_id, hero_name);
CREATE INDEX IF NOT EXISTS idx_pms_user_time ON player_match_stats(user_id, match_time DESC);
CREATE INDEX IF NOT EXISTS idx_crosshairs_user_created ON crosshairs(user_id, created_at DESC);