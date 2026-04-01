package pool

import "errors"

var (
	// ErrInvalidPoolConfig indicates invalid pool configuration
	ErrInvalidPoolConfig = errors.New("invalid pool configuration")

	// ErrConnectionFailed indicates database connection failure
	ErrConnectionFailed = errors.New("database connection failed")

	// ErrHealthCheckFailed indicates health check failure
	ErrHealthCheckFailed = errors.New("health check failed")

	// ErrPoolClosed indicates operation on closed pool
	ErrPoolClosed = errors.New("pool is closed")

	// ErrTimeout indicates operation timeout
	ErrTimeout = errors.New("operation timeout")
)
