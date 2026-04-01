-- ============================================================================
-- Query Performance Analysis Helper
-- ============================================================================
-- This script provides utilities for analyzing query performance and
-- identifying slow queries that might benefit from indexes.
-- ============================================================================

-- Enable query execution time tracking
\timing on

-- ============================================================================
-- CONFIGURATION
-- ============================================================================

-- Show query plans for slow queries
SET enable_seqscan = on;  -- Allow sequential scans
SET work_mem = '16MB';     -- Memory for sorts and hashes
SET random_page_cost = 1.1; -- SSD optimization

-- ============================================================================
-- QUERY PLAN ANALYSIS HELPERS
-- ============================================================================

\echo '============================================================================'
\echo 'QUERY PLAN COST ANALYSIS'
\echo '============================================================================'

-- Function to compare query performance with/without specific index
CREATE OR REPLACE FUNCTION compare_with_without_index(
    query_text TEXT,
    index_name TEXT
) RETURNS TABLE (
    scenario TEXT,
    execution_time_ms NUMERIC,
    total_cost NUMERIC
) AS $$
DECLARE
    start_time TIMESTAMP;
    end_time TIMESTAMP;
    plan_json JSONB;
BEGIN
    -- With index
    EXECUTE 'SET enable_indexscan = on';
    start_time := clock_timestamp();
    EXECUTE query_text;
    end_time := clock_timestamp();
    
    RETURN QUERY SELECT 
        'with_index'::TEXT,
        EXTRACT(MILLISECONDS FROM (end_time - start_time)),
        0::NUMERIC;
    
    -- Without index
    EXECUTE 'SET enable_indexscan = off';
    start_time := clock_timestamp();
    EXECUTE query_text;
    end_time := clock_timestamp();
    
    RETURN QUERY SELECT 
        'without_index'::TEXT,
        EXTRACT(MILLISECONDS FROM (end_time - start_time)),
        0::NUMERIC;
    
    EXECUTE 'SET enable_indexscan = on';
END;
$$ LANGUAGE plpgsql;

-- ============================================================================
-- MISSING INDEX DETECTION
-- ============================================================================

\echo ''
\echo '--- Detect missing indexes on foreign keys ---'
SELECT
    tc.table_name,
    kcu.column_name,
    ccu.table_name AS foreign_table_name,
    ccu.column_name AS foreign_column_name
FROM information_schema.table_constraints AS tc
JOIN information_schema.key_column_usage AS kcu
    ON tc.constraint_name = kcu.constraint_name
    AND tc.table_schema = kcu.table_schema
JOIN information_schema.constraint_column_usage AS ccu
    ON ccu.constraint_name = tc.constraint_name
    AND ccu.table_schema = tc.table_schema
WHERE tc.constraint_type = 'FOREIGN KEY'
    AND tc.table_schema = 'public'
    AND NOT EXISTS (
        SELECT 1
        FROM pg_indexes
        WHERE schemaname = 'public'
            AND tablename = tc.table_name
            AND indexdef LIKE '%' || kcu.column_name || '%'
    );

-- ============================================================================
-- SEQUENTIAL SCAN DETECTION
-- ============================================================================

\echo ''
\echo '--- Tables with high sequential scan ratio ---'
SELECT
    schemaname,
    tablename,
    seq_scan,
    seq_tup_read,
    idx_scan,
    idx_tup_fetch,
    CASE 
        WHEN seq_scan + idx_scan > 0 THEN
            ROUND(100.0 * seq_scan / (seq_scan + idx_scan), 2)
        ELSE 0
    END AS seq_scan_pct,
    pg_size_pretty(pg_relation_size(schemaname||'.'||tablename)) as table_size
FROM pg_stat_user_tables
WHERE schemaname = 'public'
    AND seq_scan + idx_scan > 0
ORDER BY seq_scan DESC;

-- ============================================================================
-- CACHE HIT RATIO
-- ============================================================================

\echo ''
\echo '--- Table cache hit ratio (should be > 99%) ---'
SELECT
    schemaname,
    tablename,
    heap_blks_read,
    heap_blks_hit,
    CASE 
        WHEN heap_blks_read + heap_blks_hit > 0 THEN
            ROUND(100.0 * heap_blks_hit / (heap_blks_read + heap_blks_hit), 2)
        ELSE 100
    END AS cache_hit_ratio
FROM pg_statio_user_tables
WHERE schemaname = 'public'
ORDER BY cache_hit_ratio ASC;

-- ============================================================================
-- BLOAT DETECTION
-- ============================================================================

\echo ''
\echo '--- Index bloat detection ---'
SELECT
    schemaname,
    tablename,
    indexname,
    pg_size_pretty(pg_relation_size(indexrelid::regclass)) as index_size,
    idx_scan,
    idx_tup_read,
    idx_tup_fetch,
    CASE 
        WHEN idx_tup_read > 0 THEN
            ROUND(100.0 * idx_tup_fetch / idx_tup_read, 2)
        ELSE 0
    END AS fetch_ratio
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
    AND pg_relation_size(indexrelid::regclass) > 1000000  -- > 1MB
ORDER BY pg_relation_size(indexrelid::regclass) DESC;

-- ============================================================================
-- SLOW QUERY SIMULATION
-- ============================================================================

\echo ''
\echo '============================================================================'
\echo 'SLOW QUERY PATTERNS'
\echo '============================================================================'

-- Pattern 1: Full table scan on large table
\echo ''
\echo '--- Pattern 1: Unindexed WHERE clause (should be slow without index) ---'
EXPLAIN (ANALYZE, BUFFERS, FORMAT TEXT)
SELECT * FROM player_match_stats
WHERE kills > 10
LIMIT 100;

-- Pattern 2: Unindexed ORDER BY
\echo ''
\echo '--- Pattern 2: Unindexed ORDER BY (check for Sort operation) ---'
EXPLAIN (ANALYZE, BUFFERS, FORMAT TEXT)
SELECT * FROM matches
ORDER BY duration_minutes DESC
LIMIT 50;

-- Pattern 3: JOIN without indexes
\echo ''
\echo '--- Pattern 3: JOIN performance (should use hash/merge join) ---'
EXPLAIN (ANALYZE, BUFFERS, FORMAT TEXT)
SELECT u.nickname, COUNT(pms.id)
FROM users u
LEFT JOIN player_match_stats pms ON u.id = pms.user_id
GROUP BY u.id, u.nickname
LIMIT 100;

-- ============================================================================
-- MAINTENANCE RECOMMENDATIONS
-- ============================================================================

\echo ''
\echo '============================================================================'
\echo 'MAINTENANCE RECOMMENDATIONS'
\echo '============================================================================'

-- Recommend VACUUM ANALYZE for stale statistics
\echo ''
\echo '--- Tables needing VACUUM ANALYZE ---'
SELECT
    schemaname,
    tablename,
    n_tup_ins + n_tup_upd + n_tup_del as total_changes,
    n_live_tup,
    n_dead_tup,
    CASE 
        WHEN n_live_tup > 0 THEN
            ROUND(100.0 * n_dead_tup / (n_live_tup + n_dead_tup), 2)
        ELSE 0
    END AS dead_tuple_pct,
    last_vacuum,
    last_autovacuum,
    last_analyze,
    last_autoanalyze
FROM pg_stat_user_tables
WHERE schemaname = 'public'
ORDER BY n_dead_tup DESC;

-- Recommend REINDEX for bloated indexes
\echo ''
\echo '--- Indexes that might benefit from REINDEX ---'
SELECT
    schemaname,
    tablename,
    indexname,
    pg_size_pretty(pg_relation_size(indexrelid::regclass)) as index_size,
    idx_scan,
    CASE 
        WHEN idx_scan = 0 THEN 'CONSIDER DROPPING (never used)'
        WHEN idx_scan < 100 THEN 'RARELY USED'
        ELSE 'FREQUENTLY USED'
    END AS usage_status
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
    AND pg_relation_size(indexrelid::regclass) > 5000000  -- > 5MB
ORDER BY idx_scan ASC, pg_relation_size(indexrelid::regclass) DESC;

\echo ''
\echo '============================================================================'
\echo 'ANALYSIS COMPLETE'
\echo '============================================================================'

