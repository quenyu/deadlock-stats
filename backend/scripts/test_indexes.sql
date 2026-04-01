-- ============================================================================
-- Index Performance Testing Script
-- ============================================================================
-- This script helps verify that indexes are being used correctly
-- and provides performance benchmarks.
--
-- Usage:
--   psql -U postgres -d deadlock_stats -f test_indexes.sql
-- ============================================================================

\timing on
\echo '============================================================================'
\echo 'INDEX PERFORMANCE TESTING'
\echo '============================================================================'

-- ============================================================================
-- 1. PLAYER_MATCH_STATS QUERIES
-- ============================================================================

\echo ''
\echo '--- Test 1: User match history (should use idx_pms_user_created) ---'
EXPLAIN ANALYZE
SELECT * FROM player_match_stats 
WHERE user_id = (SELECT id FROM users LIMIT 1)
ORDER BY id DESC
LIMIT 20;

\echo ''
\echo '--- Test 2: Hero statistics (should use idx_pms_user_hero) ---'
EXPLAIN ANALYZE
SELECT hero_name, COUNT(*), AVG(kills), AVG(deaths), AVG(assists)
FROM player_match_stats
WHERE user_id = (SELECT id FROM users LIMIT 1)
GROUP BY hero_name;

\echo ''
\echo '--- Test 3: Win rate calculation (should use idx_pms_user_result) ---'
EXPLAIN ANALYZE
SELECT result, COUNT(*) 
FROM player_match_stats
WHERE user_id = (SELECT id FROM users LIMIT 1)
GROUP BY result;

\echo ''
\echo '--- Test 4: Match participants (should use idx_pms_match_kills) ---'
EXPLAIN ANALYZE
SELECT user_id, hero_name, kills, deaths, assists
FROM player_match_stats
WHERE match_id = (SELECT id FROM matches LIMIT 1)
ORDER BY kills DESC;

-- ============================================================================
-- 2. MATCHES QUERIES
-- ============================================================================

\echo ''
\echo '--- Test 5: Recent matches (should use idx_matches_time) ---'
EXPLAIN ANALYZE
SELECT * FROM matches
ORDER BY match_time DESC
LIMIT 100;

\echo ''
\echo '--- Test 6: Map-specific matches (should use idx_matches_map_time) ---'
EXPLAIN ANALYZE
SELECT * FROM matches
WHERE map_name = (SELECT map_name FROM matches LIMIT 1)
ORDER BY match_time DESC
LIMIT 50;

-- ============================================================================
-- 3. USERS QUERIES
-- ============================================================================

\echo ''
\echo '--- Test 7: Fuzzy nickname search (should use idx_users_nickname_gin) ---'
EXPLAIN ANALYZE
SELECT * FROM users
WHERE nickname ILIKE '%test%'
LIMIT 10;

\echo ''
\echo '--- Test 8: Exact nickname lookup (should use idx_users_nickname_lower) ---'
EXPLAIN ANALYZE
SELECT * FROM users
WHERE LOWER(nickname) = 'testuser';

\echo ''
\echo '--- Test 9: Steam ID lookup (should use unique index) ---'
EXPLAIN ANALYZE
SELECT * FROM users
WHERE steam_id = '76561198000000000';

-- ============================================================================
-- 4. PLAYER_STATS QUERIES
-- ============================================================================

\echo ''
\echo '--- Test 10: KD Ratio leaderboard (should use idx_ps_kd_ratio) ---'
EXPLAIN ANALYZE
SELECT user_id, kd_ratio, win_rate, favorite_hero
FROM player_stats
WHERE kd_ratio > 0
ORDER BY kd_ratio DESC
LIMIT 100;

\echo ''
\echo '--- Test 11: Win Rate leaderboard (should use idx_ps_win_rate) ---'
EXPLAIN ANALYZE
SELECT user_id, win_rate, kd_ratio, favorite_hero
FROM player_stats
WHERE win_rate > 0
ORDER BY win_rate DESC
LIMIT 100;

-- ============================================================================
-- 5. CROSSHAIRS QUERIES
-- ============================================================================

\echo ''
\echo '--- Test 12: Public crosshairs feed (should use idx_crosshairs_public_created) ---'
EXPLAIN ANALYZE
SELECT * FROM crosshairs
WHERE is_public = true
ORDER BY created_at DESC
LIMIT 50;

\echo ''
\echo '--- Test 13: Popular crosshairs (should use idx_crosshairs_public_views) ---'
EXPLAIN ANALYZE
SELECT * FROM crosshairs
WHERE is_public = true
ORDER BY view_count DESC
LIMIT 50;

\echo ''
\echo '--- Test 14: User crosshairs (should use idx_crosshairs_author_created) ---'
EXPLAIN ANALYZE
SELECT * FROM crosshairs
WHERE author_id = (SELECT id FROM users LIMIT 1)
ORDER BY created_at DESC;

-- ============================================================================
-- 6. BUILDS QUERIES
-- ============================================================================

\echo ''
\echo '--- Test 15: Public builds (should use idx_builds_public_created) ---'
EXPLAIN ANALYZE
SELECT * FROM builds
WHERE is_public = true
ORDER BY created_at DESC
LIMIT 50;

\echo ''
\echo '--- Test 16: Version-specific builds (should use idx_builds_version_created) ---'
EXPLAIN ANALYZE
SELECT * FROM builds
WHERE game_version IS NOT NULL
ORDER BY created_at DESC
LIMIT 50;

-- ============================================================================
-- 7. COMPLEX JOINS
-- ============================================================================

\echo ''
\echo '--- Test 17: User stats with recent matches (multiple indexes) ---'
EXPLAIN ANALYZE
SELECT 
    u.nickname,
    ps.kd_ratio,
    ps.win_rate,
    COUNT(pms.id) as match_count
FROM users u
LEFT JOIN player_stats ps ON u.id = ps.user_id
LEFT JOIN player_match_stats pms ON u.id = pms.user_id
WHERE u.id = (SELECT id FROM users LIMIT 1)
GROUP BY u.id, u.nickname, ps.kd_ratio, ps.win_rate;

\echo ''
\echo '--- Test 18: Match details with participants (should use multiple indexes) ---'
EXPLAIN ANALYZE
SELECT 
    m.id,
    m.map_name,
    m.duration_minutes,
    m.match_time,
    COUNT(pms.id) as player_count,
    AVG(pms.kills) as avg_kills
FROM matches m
LEFT JOIN player_match_stats pms ON m.id = pms.match_id
GROUP BY m.id
ORDER BY m.match_time DESC
LIMIT 20;

-- ============================================================================
-- INDEX USAGE STATISTICS
-- ============================================================================

\echo ''
\echo '============================================================================'
\echo 'INDEX USAGE STATISTICS'
\echo '============================================================================'
\echo ''

-- Show index sizes
\echo '--- Index Sizes ---'
SELECT
    schemaname,
    tablename,
    indexname,
    pg_size_pretty(pg_relation_size(indexrelid::regclass)) as index_size
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
ORDER BY pg_relation_size(indexrelid::regclass) DESC;

\echo ''
\echo '--- Index Usage (requires statistics collection) ---'
SELECT
    schemaname,
    tablename,
    indexname,
    idx_scan as index_scans,
    idx_tup_read as tuples_read,
    idx_tup_fetch as tuples_fetched
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
ORDER BY idx_scan DESC;

\echo ''
\echo '--- Unused Indexes (idx_scan = 0) ---'
SELECT
    schemaname,
    tablename,
    indexname,
    pg_size_pretty(pg_relation_size(indexrelid::regclass)) as index_size
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
  AND idx_scan = 0
  AND indexrelid::regclass::text NOT LIKE '%_pkey'
ORDER BY pg_relation_size(indexrelid::regclass) DESC;

\echo ''
\echo '--- Table Sizes ---'
SELECT
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as total_size,
    pg_size_pretty(pg_relation_size(schemaname||'.'||tablename)) as table_size,
    pg_size_pretty(pg_indexes_size(schemaname||'.'||tablename)) as indexes_size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

\echo ''
\echo '============================================================================'
\echo 'TESTING COMPLETE'
\echo '============================================================================'
\echo ''
\echo 'How to interpret results:'
\echo '  - Look for "Index Scan" in EXPLAIN ANALYZE output (good)'
\echo '  - "Seq Scan" means table scan without index (bad for large tables)'
\echo '  - Execution time should be < 10ms for indexed queries'
\echo '  - Check "Index Usage Statistics" for unused indexes'
\echo ''

