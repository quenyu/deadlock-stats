package indexes

import (
	"context"
	"database/sql"

	"go.uber.org/zap"
)

type Analyzer struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewAnalyzer(db *sql.DB, logger *zap.Logger) *Analyzer {
	return &Analyzer{
		db:     db,
		logger: logger,
	}
}

type IndexStats struct {
	SchemaName  string
	TableName   string
	IndexName   string
	IndexScans  int64
	TuplesRead  int64
	TuplesFetch int64
	SizeBytes   int64
}

// GetUnusedIndexes returns indexes that are never used
func (a *Analyzer) GetUnusedIndexes(ctx context.Context) ([]IndexStats, error) {
	query := `
		SELECT
			schemaname,
			tablename,
			indexname,
			idx_scan,
			idx_tup_read,
			idx_tup_fetch,
			pg_relation_size(indexrelid::regclass) as size_bytes
		FROM pg_stat_user_indexes
		WHERE schemaname = 'public'
		  AND idx_scan = 0
		  AND indexrelid::regclass::text NOT LIKE '%_pkey'
		ORDER BY pg_relation_size(indexrelid::regclass) DESC
	`

	rows, err := a.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []IndexStats
	for rows.Next() {
		var s IndexStats
		if err := rows.Scan(
			&s.SchemaName,
			&s.TableName,
			&s.IndexName,
			&s.IndexScans,
			&s.TuplesRead,
			&s.TuplesFetch,
			&s.SizeBytes,
		); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}

	return stats, rows.Err()
}

// GetMostUsedIndexes returns most frequently used indexes
func (a *Analyzer) GetMostUsedIndexes(ctx context.Context, limit int) ([]IndexStats, error) {
	query := `
		SELECT
			schemaname,
			tablename,
			indexname,
			idx_scan,
			idx_tup_read,
			idx_tup_fetch,
			pg_relation_size(indexrelid::regclass) as size_bytes
		FROM pg_stat_user_indexes
		WHERE schemaname = 'public'
		ORDER BY idx_scan DESC
		LIMIT $1
	`

	rows, err := a.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []IndexStats
	for rows.Next() {
		var s IndexStats
		if err := rows.Scan(
			&s.SchemaName,
			&s.TableName,
			&s.IndexName,
			&s.IndexScans,
			&s.TuplesRead,
			&s.TuplesFetch,
			&s.SizeBytes,
		); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}

	return stats, rows.Err()
}

// LogUnusedIndexes logs unused indexes
func (a *Analyzer) LogUnusedIndexes(ctx context.Context) {
	indexes, err := a.GetUnusedIndexes(ctx)
	if err != nil {
		a.logger.Error("failed to get unused indexes", zap.Error(err))
		return
	}

	if len(indexes) == 0 {
		a.logger.Info("no unused indexes found")
		return
	}

	a.logger.Warn("found unused indexes",
		zap.Int("count", len(indexes)),
	)

	for _, idx := range indexes {
		a.logger.Debug("unused index",
			zap.String("table", idx.TableName),
			zap.String("index", idx.IndexName),
			zap.Int64("size_bytes", idx.SizeBytes),
		)
	}
}
