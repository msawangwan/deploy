package cache_dep

import (
	"testing"

	"github.com/msawangwan/ci.io/lib/symbol"
)

var c Cache

func TestCreateACache(t *testing.T) {
	c = newBasicCache()

	if c.Contains(0) {
		t.Fatalf("shouldn't contain item: %s", symbol.FailMark)
	}

	t.Logf("create a cache %s", symbol.PassMark)
}

func TestCacheItem(t *testing.T) {
	c = newBasicCache()

	c.Store("some_value")
	c.Store("another_value")

	t.Logf("cached items: %+v", c)
	t.Logf("able to store items: %s", symbol.PassMark)
}

func TestFetchItemFromStore(t *testing.T) {
	c = newBasicCache()

	vs := []string{
		"value1",
		"value2",
		"value3",
	}

	ks := []int{}

	for _, item := range vs {
		k := c.Store(item)
		ks = append(ks, k)
	}

	t.Logf("%v", ks)
}
