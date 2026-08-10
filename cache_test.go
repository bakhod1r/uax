package uax

import (
	"sync"
	"testing"
)

func TestLRUCache(t *testing.T) {
	c := NewCache(2)

	cli1 := &Client{Raw: "ua1"}
	cli2 := &Client{Raw: "ua2"}
	cli3 := &Client{Raw: "ua3"}

	c.Set("ua1", cli1)
	c.Set("ua2", cli2)

	if _, ok := c.Get("ua1"); !ok {
		t.Error("expected ua1 to be in cache")
	}

	c.Set("ua3", cli3) // Should evict ua2 because ua1 was recently accessed

	if _, ok := c.Get("ua2"); ok {
		t.Error("expected ua2 to be evicted")
	}

	if _, ok := c.Get("ua3"); !ok {
		t.Error("expected ua3 to be in cache")
	}
}

func TestLRUCache_Concurrent(t *testing.T) {
	c := NewCache(100)
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			c.Set("ua", &Client{Raw: "ua"})
			c.Get("ua")
		}(i)
	}

	wg.Wait()
}

func TestCacheDefaultSize(t *testing.T) {
	c := NewCache(0)
	if c.maxEntries != 1000 {
		t.Errorf("expected default maxEntries to be 1000, got %d", c.maxEntries)
	}
}
