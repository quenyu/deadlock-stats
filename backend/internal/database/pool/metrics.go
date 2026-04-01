package pool

import (
	"sync"
	"time"
)

type Metrics struct {
	mu                  sync.RWMutex
	OpenConnections     int
	InUse               int
	Idle                int
	WaitCount           int64
	WaitDuration        time.Duration
	MaxIdleClosed       int64
	MaxLifetimeClosed   int64
	ConnectionErrors    int64
	LastHealthCheck     time.Time
	HealthCheckDuration time.Duration
}

func NewMetrics() *Metrics {
	return &Metrics{}
}

func (m *Metrics) Update(
	open, inUse, idle int,
	waitCount, maxIdleClosed, maxLifetimeClosed int64,
	waitDuration time.Duration,
) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.OpenConnections = open
	m.InUse = inUse
	m.Idle = idle
	m.WaitCount = waitCount
	m.WaitDuration = waitDuration
	m.MaxIdleClosed = maxIdleClosed
	m.MaxLifetimeClosed = maxLifetimeClosed
}

// IncrementErrors increments connection error counter
func (m *Metrics) IncrementErrors() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ConnectionErrors++
}

// UpdateHealthCheck updates health check metrics
func (m *Metrics) UpdateHealthCheck(duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.LastHealthCheck = time.Now()
	m.HealthCheckDuration = duration
}

// Snapshot returns a copy of current metrics
func (m *Metrics) Snapshot() Metrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return Metrics{
		OpenConnections:     m.OpenConnections,
		InUse:               m.InUse,
		Idle:                m.Idle,
		WaitCount:           m.WaitCount,
		WaitDuration:        m.WaitDuration,
		MaxIdleClosed:       m.MaxIdleClosed,
		MaxLifetimeClosed:   m.MaxLifetimeClosed,
		ConnectionErrors:    m.ConnectionErrors,
		LastHealthCheck:     m.LastHealthCheck,
		HealthCheckDuration: m.HealthCheckDuration,
	}
}
