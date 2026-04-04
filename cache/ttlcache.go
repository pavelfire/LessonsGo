package cache

import (
	"errors"
	"sync"
	"time"
)

var ErrNonPositiveTTL = errors.New("ttl must be positive")

type entry[V any] struct {
	value     V
	expiresAt time.Time
}

// TTLCache is an in-memory key-value store with per-entry expiration.
// Expired entries are removed on Get and by a background sweeper.
type TTLCache[K comparable, V any] struct {
	mu              sync.RWMutex
	items           map[K]entry[V]
	cleanupInterval time.Duration
	stopCh          chan struct{}
	stopOnce        sync.Once
}

// NewTTLCache creates a cache. cleanupInterval controls how often a
// background pass removes expired keys (0 disables sweeping).
func NewTTLCache[K comparable, V any](cleanupInterval time.Duration) *TTLCache[K, V] {
	c := &TTLCache[K, V]{
		items:           make(map[K]entry[V]),
		cleanupInterval: cleanupInterval,
		stopCh:          make(chan struct{}),
	}
	if cleanupInterval > 0 {
		go c.sweepLoop()
	}
	return c
}

// Set stores value under key until now+ttl.
func (c *TTLCache[K, V]) Set(key K, value V, ttl time.Duration) error {
	if ttl <= 0 {
		return ErrNonPositiveTTL
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = entry[V]{value: value, expiresAt: time.Now().Add(ttl)}
	return nil
}

// Get returns the value if present and not expired.
func (c *TTLCache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.items[key]
	if !ok {
		var zero V
		return zero, false
	}
	if time.Now().After(e.expiresAt) {
		delete(c.items, key)
		var zero V
		return zero, false
	}
	return e.value, true
}

// Delete removes key regardless of expiry.
func (c *TTLCache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Len returns the number of keys including not-yet-swept expired ones.
func (c *TTLCache[K, V]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

// Stop ends the background sweeper. Safe to call multiple times.
func (c *TTLCache[K, V]) Stop() {
	c.stopOnce.Do(func() { close(c.stopCh) })
}

func (c *TTLCache[K, V]) sweepLoop() {
	t := time.NewTicker(c.cleanupInterval)
	defer t.Stop()
	for {
		select {
		case <-c.stopCh:
			return
		case <-t.C:
			c.deleteExpired()
		}
	}
}

func (c *TTLCache[K, V]) deleteExpired() {
	now := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, e := range c.items {
		if now.After(e.expiresAt) {
			delete(c.items, k)
		}
	}
}
