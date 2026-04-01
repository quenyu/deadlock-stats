-- Удаляем индексы производительности

-- Удаляем индексы для player_match_stats
DROP INDEX IF EXISTS idx_pms_user_id;
DROP INDEX IF EXISTS idx_pms_match_id;
DROP INDEX IF EXISTS idx_pms_hero_name;
DROP INDEX IF EXISTS idx_pms_match_time;

-- Удаляем индексы для matches
DROP INDEX IF EXISTS idx_matches_time;
DROP INDEX IF EXISTS idx_matches_duration;

-- Удаляем индексы для users
DROP INDEX IF EXISTS idx_users_nickname;
DROP INDEX IF EXISTS idx_users_steam_id;

-- Удаляем индексы для crosshairs
DROP INDEX IF EXISTS idx_crosshairs_user_id;
DROP INDEX IF EXISTS idx_crosshairs_created_at;
DROP INDEX IF EXISTS idx_crosshairs_likes;

-- Удаляем индексы для votes
DROP INDEX IF EXISTS idx_votes_user_id;
DROP INDEX IF EXISTS idx_votes_crosshair_id;

-- Удаляем индексы для comments
DROP INDEX IF EXISTS idx_comments_crosshair_id;
DROP INDEX IF EXISTS idx_comments_user_id;
DROP INDEX IF EXISTS idx_comments_created_at;

-- Удаляем индексы для builds
DROP INDEX IF EXISTS idx_builds_user_id;
DROP INDEX IF EXISTS idx_builds_hero_name;
DROP INDEX IF EXISTS idx_builds_created_at;
DROP INDEX IF EXISTS idx_builds_rating;

-- Удаляем составные индексы
DROP INDEX IF EXISTS idx_pms_user_hero;
DROP INDEX IF EXISTS idx_pms_user_time;
DROP INDEX IF EXISTS idx_crosshairs_user_created;