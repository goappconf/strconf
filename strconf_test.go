package strconf

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	c := New()

	c.Set("username", "alice", time.Minute)

	value, ok := c.Get("username")
	if !ok {
		t.Fatal("expected value")
	}

	if value != "alice" {
		t.Fatalf("expected alice, got %v", value)
	}
}
