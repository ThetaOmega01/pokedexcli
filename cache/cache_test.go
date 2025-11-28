package cache

import (
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	cache := NewCache(5 * time.Second)

	// Test adding and getting an entry
	key := "test-key"
	val := []byte("test-value")

	cache.Add(key, val)

	retrieved, ok := cache.Get(key)
	if !ok {
		t.Errorf("Expected to find key %s in cache", key)
	}

	if string(retrieved) != string(val) {
		t.Errorf("Expected value %s, got %s", string(val), string(retrieved))
	}
}

func TestGetNotFound(t *testing.T) {
	cache := NewCache(5 * time.Second)

	// Test getting a non-existent entry
	_, ok := cache.Get("non-existent-key")
	if ok {
		t.Error("Expected to not find non-existent key in cache")
	}
}

func TestReapLoop(t *testing.T) {
	interval := 10 * time.Millisecond
	cache := NewCache(interval)

	// Add an entry
	key := "test-key"
	val := []byte("test-value")
	cache.Add(key, val)

	// Verify it exists
	_, ok := cache.Get(key)
	if !ok {
		t.Error("Expected to find key in cache immediately after adding")
	}

	// Wait for longer than the interval
	time.Sleep(interval + 5*time.Millisecond)

	// The entry should be reaped
	_, ok = cache.Get(key)
	if ok {
		t.Error("Expected entry to be reaped after interval")
	}
}

func TestConcurrentAccess(t *testing.T) {
	cache := NewCache(5 * time.Second)

	// Test concurrent writes and reads
	done := make(chan bool)

	// Writer goroutine
	go func() {
		for range 100 {
			cache.Add("key", []byte("value"))
		}
		done <- true
	}()

	// Reader goroutine
	go func() {
		for range 100 {
			cache.Get("key")
		}
		done <- true
	}()

	// Wait for both goroutines to complete
	<-done
	<-done

	// If we get here without a panic, the test passes
}
