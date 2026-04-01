package ratelimit

import (
	"context"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type Limiter struct {
	limiters map[string]*limiterEntry
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

type limiterEntry struct {
	limiter  *rate.Limiter
	lastUsed time.Time
}

func NewLimiter(rateLimit rate.Limit, burst int) *Limiter {
	return &Limiter{
		limiters: make(map[string]*limiterEntry),
		rate:     rateLimit,
		burst:    burst,
	}
}

func (l *Limiter) GetLimiter(key string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	entry, exists := l.limiters[key]
	if !exists {
		limiter := rate.NewLimiter(l.rate, l.burst)
		entry = &limiterEntry{
			limiter:  limiter,
			lastUsed: time.Now(),
		}
		l.limiters[key] = entry
	} else {
		entry.lastUsed = time.Now()
	}

	return entry.limiter
}

// Allow checks to see if the request is allowed
func (l *Limiter) Allow(key string, limit int) (allowed bool, remaining int, resetAt int64, err error) {
	limiter := l.GetLimiter(key)

	if rate.Limit(limit) != l.rate {
		limiter.SetLimit(rate.Limit(limit))
	}

	allowed = limiter.Allow()

	tokens := int(limiter.Tokens())
	if tokens < 0 {
		tokens = 0
	}

	resetAt = time.Now().Add(time.Second / time.Duration(limit)).Unix()

	return allowed, tokens, resetAt, nil
}

// Close closes the limiter (corresponds to the ILimiter interface)
func (l *Limiter) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.limiters = make(map[string]*limiterEntry)

	return nil
}

func (l *Limiter) Wait(ctx context.Context, key string) error {
	limiter := l.GetLimiter(key)
	return limiter.Wait(ctx)
}

// Cleanup removes old limiters to save memory
func (l *Limiter) Cleanup() {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Remove limiters older than 1 hour
	cutoff := time.Now().Add(-time.Hour)
	for key, entry := range l.limiters {
		if entry.lastUsed.Before(cutoff) {
			delete(l.limiters, key)
		}
	}
}
