package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisLimiter struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisLimiter(client *redis.Client, ttl time.Duration) *RedisLimiter {
	return &RedisLimiter{
		client: client,
		ttl:    ttl,
	}
}

// Allow checks if request is allowed using token bucket algorithm
func (rl *RedisLimiter) Allow(key string, limit int) (allowed bool, remaining int, resetAt int64, err error) {
	ctx := context.Background()
	redisKey := fmt.Sprintf("ratelimit:%s", key)

	// Lua script for atomic token bucket implementation
	script := `
		local key = KEYS[1]
		local limit = tonumber(ARGV[1])
		local ttl = tonumber(ARGV[2])
		local now = tonumber(ARGV[3])
		
		-- Get current state
		local state = redis.call('HMGET', key, 'tokens', 'last_refill')
		local tokens = tonumber(state[1])
		local last_refill = tonumber(state[2])
		
		-- Initialize if not exists
		if tokens == nil then
			tokens = limit
			last_refill = now
		end
		
		-- Calculate tokens to add based on time passed
		local time_passed = now - last_refill
		local tokens_to_add = math.floor(time_passed * limit)
		
		if tokens_to_add > 0 then
			tokens = math.min(limit, tokens + tokens_to_add)
			last_refill = now
		end
		
		-- Try to consume a token
		local allowed = 0
		if tokens > 0 then
			tokens = tokens - 1
			allowed = 1
		end
		
		-- Update state
		redis.call('HMSET', key, 'tokens', tokens, 'last_refill', last_refill)
		redis.call('EXPIRE', key, ttl)
		
		-- Calculate reset time (when next token will be available)
		local reset_at = last_refill + 1
		
		return {allowed, tokens, reset_at}
	`

	now := time.Now().Unix()
	result, err := rl.client.Eval(ctx, script, []string{redisKey}, limit, int(rl.ttl.Seconds()), now).Result()
	if err != nil {
		return false, 0, 0, fmt.Errorf("%w: %v", ErrRedisConnection, err)
	}

	// Parse result
	values, ok := result.([]interface{})
	if !ok || len(values) != 3 {
		return false, 0, 0, fmt.Errorf("unexpected redis response")
	}

	allowedInt, _ := values[0].(int64)
	remainingInt, _ := values[1].(int64)
	resetAtInt, _ := values[2].(int64)

	return allowedInt == 1, int(remainingInt), resetAtInt, nil
}

// Close closes the Redis limiter
func (rl *RedisLimiter) Close() error {
	// Redis client is managed externally, don't close it
	return nil
}

// Ping checks if Redis connection is alive
func (rl *RedisLimiter) Ping(ctx context.Context) error {
	return rl.client.Ping(ctx).Err()
}

// Reset resets rate limit for a specific key
func (rl *RedisLimiter) Reset(ctx context.Context, key string) error {
	redisKey := fmt.Sprintf("ratelimit:%s", key)
	return rl.client.Del(ctx, redisKey).Err()
}

// GetState returns current state for a key
func (rl *RedisLimiter) GetState(ctx context.Context, key string) (tokens int, lastRefill int64, err error) {
	redisKey := fmt.Sprintf("ratelimit:%s", key)

	result, err := rl.client.HMGet(ctx, redisKey, "tokens", "last_refill").Result()
	if err != nil {
		return 0, 0, err
	}

	if len(result) != 2 {
		return 0, 0, fmt.Errorf("unexpected redis response")
	}

	if result[0] != nil {
		if tokensStr, ok := result[0].(string); ok {
			fmt.Sscanf(tokensStr, "%d", &tokens)
		}
	}

	if result[1] != nil {
		if lastRefillStr, ok := result[1].(string); ok {
			fmt.Sscanf(lastRefillStr, "%d", &lastRefill)
		}
	}

	return tokens, lastRefill, nil
}
