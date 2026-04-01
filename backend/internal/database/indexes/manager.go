package indexes

import (
	"context"
	"database/sql"
	"fmt"

	"go.uber.org/zap"
)

type Manager struct {
	db       *sql.DB
	analyzer *Analyzer
	logger   *zap.Logger
}

func NewManager(db *sql.DB, logger *zap.Logger) *Manager {
	return &Manager{
		db:       db,
		analyzer: NewAnalyzer(db, logger),
		logger:   logger,
	}
}

func (m *Manager) Analyzer() *Analyzer {
	return m.analyzer
}

// VacuumAnalyze runs VACUUM ANALYZE on specified tables
func (m *Manager) VacuumAnalyze(ctx context.Context, tables ...string) error {
	if len(tables) == 0 {
		if _, err := m.db.ExecContext(ctx, "VACUUM ANALYZE"); err != nil {
			return fmt.Errorf("vacuum analyze failed: %w", err)
		}
		m.logger.Info("vacuum analyze completed for all tables")
		return nil
	}

	// Vacuum specific tables
	for _, table := range tables {
		query := fmt.Sprintf("VACUUM ANALYZE %s", table)
		if _, err := m.db.ExecContext(ctx, query); err != nil {
			return fmt.Errorf("vacuum analyze %s failed: %w", table, err)
		}
		m.logger.Debug("vacuum analyze completed", zap.String("table", table))
	}

	m.logger.Info("vacuum analyze completed", zap.Int("tables", len(tables)))
	return nil
}

// ReindexTable rebuilds all indexes on a table
func (m *Manager) ReindexTable(ctx context.Context, table string) error {
	query := fmt.Sprintf("REINDEX TABLE %s", table)
	if _, err := m.db.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("reindex table %s failed: %w", table, err)
	}

	m.logger.Info("reindex completed", zap.String("table", table))
	return nil
}

// ReindexIndex rebuilds specific index
func (m *Manager) ReindexIndex(ctx context.Context, index string) error {
	query := fmt.Sprintf("REINDEX INDEX %s", index)
	if _, err := m.db.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("reindex %s failed: %w", index, err)
	}

	m.logger.Info("reindex completed", zap.String("index", index))
	return nil
}

// DropIndex drops an index
func (m *Manager) DropIndex(ctx context.Context, index string) error {
	query := fmt.Sprintf("DROP INDEX IF EXISTS %s", index)
	if _, err := m.db.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("drop index %s failed: %w", index, err)
	}

	m.logger.Info("index dropped", zap.String("index", index))
	return nil
}
