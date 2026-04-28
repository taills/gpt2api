package lock

import (
	"context"
	"errors"
	"sync"
	"time"
)

// ErrNotAcquired is returned when the lock cannot be acquired.
var ErrNotAcquired = errors.New("lock: not acquired")

type entry struct {
	token   string
	expires time.Time
}

// MemoryLock is a thread-safe in-memory distributed lock for single-process use.
type MemoryLock struct {
	mu      sync.Mutex
	entries map[string]entry
}

// NewMemoryLock creates a new in-memory lock.
func NewMemoryLock() *MemoryLock {
	return &MemoryLock{entries: make(map[string]entry)}
}

// Acquire tries to acquire the lock for key with the given token and TTL.
// Returns ErrNotAcquired if the key is already locked.
func (l *MemoryLock) Acquire(_ context.Context, key, token string, ttl time.Duration) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if e, ok := l.entries[key]; ok && time.Now().Before(e.expires) {
		return ErrNotAcquired
	}
	l.entries[key] = entry{token: token, expires: time.Now().Add(ttl)}
	return nil
}

// Release releases the lock only if the token matches the current holder.
func (l *MemoryLock) Release(_ context.Context, key, token string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if e, ok := l.entries[key]; ok && e.token == token {
		delete(l.entries, key)
	}
	return nil
}

// Refresh extends the TTL of a lock only if the token matches the current holder.
func (l *MemoryLock) Refresh(_ context.Context, key, token string, ttl time.Duration) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if e, ok := l.entries[key]; ok && e.token == token && time.Now().Before(e.expires) {
		l.entries[key] = entry{token: token, expires: time.Now().Add(ttl)}
	}
	return nil
}
