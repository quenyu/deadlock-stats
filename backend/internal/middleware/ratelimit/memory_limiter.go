package ratelimit

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// MemoryLimiter implements in-memory rate limiting using token bucket
type MemoryLimiter struct {
	limiters map[string]*limiterEntry
	mu       sync.RWMutex
	burst    int
	ttl      time.Duration
	stopChan chan struct{}
}

type MemoryLimiterEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// NewMemoryLimiter creates a new in-memory rate limiter
func NewMemoryLimiter(burst int, ttl time.Duration) *MemoryLimiter {
	ml := &MemoryLimiter{
		limiters: make(map[string]*limiterEntry),
		burst:    burst,
		ttl:      ttl,
		stopChan: make(chan struct{}),
	}

	go ml.cleanup()

	return ml
}

// Allow checks if request is allowed
func (ml *MemoryLimiter) Allow(key string, limit int) (allowed bool, remaining int, resetAt int64, err error) {
	limiter := ml.getLimiter(key, limit)

	allowed = limiter.Allow()

	tokens := int(limiter.Tokens())
	if tokens < 0 {
		tokens = 0
	}

	resetAt = time.Now().Add(time.Second / time.Duration(limit)).Unix()

	return allowed, tokens, resetAt, nil
}

// getLimiter returns or creates a rate limiter for the given key
func (ml *MemoryLimiter) getLimiter(key string, limit int) *rate.Limiter {
	ml.mu.Lock()
	defer ml.mu.Unlock()

	entry, exists := ml.limiters[key]
	if !exists {
		limiter := rate.NewLimiter(rate.Limit(limit), ml.burst)
		ml.limiters[key] = &limiterEntry{
			limiter:  limiter,
			lastUsed: time.Now(),
		}
		return limiter
	}

	entry.lastUsed = time.Now()

	entry.limiter.SetLimit(rate.Limit(limit))

	return entry.limiter
}

// cleanup removes expired limiters
func (ml *MemoryLimiter) cleanup() {
	ticker := time.NewTicker(ml.ttl)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ml.removeExpired()
		case <-ml.stopChan:
			return
		}
	}
}

// removeExpired removes limiters that haven't been used recently
func (ml *MemoryLimiter) removeExpired() {
	ml.mu.Lock()
	defer ml.mu.Unlock()

	now := time.Now()
	for key, entry := range ml.limiters {
		if now.Sub(entry.lastUsed) > ml.ttl {
			delete(ml.limiters, key)
		}
	}
}

// Close stops the cleanup goroutine
func (ml *MemoryLimiter) Close() error {
	close(ml.stopChan)

	ml.mu.Lock()
	defer ml.mu.Unlock()
	ml.limiters = make(map[string]*limiterEntry)

	return nil
}

// Count returns the number of active limiters
func (ml *MemoryLimiter) Count() int {
	ml.mu.RLock()
	defer ml.mu.RUnlock()
	return len(ml.limiters)
}
