package ratelimit

import (
	"context"
	"sync"
	"time"
)

type bucketState struct {
	tokens     float64
	lastRefill time.Time
}

// MemoryBucket is an in-memory token bucket rate limiter.
type MemoryBucket struct {
	mu      sync.Mutex
	buckets map[string]*bucketState
}

// NewMemoryBucket creates a new in-memory token bucket.
func NewMemoryBucket() *MemoryBucket {
	return &MemoryBucket{buckets: make(map[string]*bucketState)}
}

// Allow tries to consume cost tokens from the bucket identified by key.
// capacity is the max tokens, refillPerSec is the token refill rate.
// Returns (allowed, remaining, error).
func (b *MemoryBucket) Allow(_ context.Context, key string, capacity, cost int64, refillPerSec float64) (bool, float64, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	s, ok := b.buckets[key]
	if !ok {
		s = &bucketState{tokens: float64(capacity), lastRefill: now}
		b.buckets[key] = s
	}

	elapsed := now.Sub(s.lastRefill).Seconds()
	s.tokens += elapsed * refillPerSec
	if s.tokens > float64(capacity) {
		s.tokens = float64(capacity)
	}
	s.lastRefill = now

	if s.tokens >= float64(cost) {
		s.tokens -= float64(cost)
		return true, s.tokens, nil
	}
	return false, s.tokens, nil
}
