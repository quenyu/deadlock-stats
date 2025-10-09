package ratelimit

import "errors"

var (
	// ErrInvalidConfig indicates invalid rate limiter configuration
	ErrInvalidConfig = errors.New("invalid rate limiter configuration")

	// ErrRateLimitExceeded indicates rate limit has been exceeded
	ErrRateLimitExceeded = errors.New("rate limit exceeded")

	// ErrRedisConnection indicates Redis connection failure
	ErrRedisConnection = errors.New("redis connection failed")

	// ErrLimiterClosed indicates operation on closed limiter
	ErrLimiterClosed = errors.New("limiter is closed")
)
