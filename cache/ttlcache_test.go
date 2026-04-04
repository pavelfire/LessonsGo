package cache

import (
	"testing"
	"time"
)

func TestTTLCache_GetAfterExpiry(t *testing.T) {
	c := NewTTLCache[string, int](0)
	defer c.Stop()

	if err := c.Set("a", 1, 50*time.Millisecond); err != nil {
		t.Fatal(err)
	}
	v, ok := c.Get("a")
	if !ok || v != 1 {
		t.Fatalf("expected 1, got %v, ok=%v", v, ok)
	}
	time.Sleep(60 * time.Millisecond)
	_, ok = c.Get("a")
	if ok {
		t.Fatal("expected miss after ttl")
	}
}

func TestTTLCache_SetInvalidTTL(t *testing.T) {
	c := NewTTLCache[string, int](0)
	defer c.Stop()
	if err := c.Set("a", 1, 0); err != ErrNonPositiveTTL {
		t.Fatalf("expected ErrNonPositiveTTL, got %v", err)
	}
}

func TestTTLCache_Sweeper(t *testing.T) {
	c := NewTTLCache[string, int](20 * time.Millisecond)
	defer c.Stop()

	_ = c.Set("x", 42, 30*time.Millisecond)
	time.Sleep(80 * time.Millisecond)
	if c.Len() != 0 {
		t.Fatalf("sweeper should remove expired keys, len=%d", c.Len())
	}
}
