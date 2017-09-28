package collection

import (
	"testing"

	"github.com/msawangwan/ci.io/lib/symbol"
)

type testCache map[string]string

func newCache() Cache {
	return testCache(make(map[string]string))
}

func (c testCache) Contains(k string) bool {
	if _, ok := c[label]; ok {
		return true
	}

	return false
}

func (c testCache) Store(k string, v interface{}) bool {
	c[k] = v.(string)
	return true
}

func (c testCache) Fetch(k string) (v interface{}, e error) { return }
func (c testCache) Valid() bool                             { return true }

var c Cache

func TestCreateACache(t *testing.T) {
	c = newCache()

	if c.Contains("some_item") {
		t.Fatalf("shouldn't contain item: %s", symbol.FailMark)
	}

	t.Logf("create a cache %s", symbol.PassMark)
}

func TestCacheItem(t *testing.T) {
	c = newCache()
	c.Store("some_key", "some_value")
	c.Store("another_key", "another_value")

	t.Logf("cached items: %+v", c)
	t.Logf("able to store items: %s", symbol.PassMark)
}

type kvpair struct{ k, v string }

func TestFetchItemFromStore(t *testing.T) {
	c = newCache()

	kvpairs := []kvpair{
		kvpair{"key1", "value1"},
		kvpair{"key2", "value2"},
		kvpair{"key3", "value3"},
	}

	for _, item := range kvpairs {
		t.Logf("%s, %s", item.k, item.v)
	}
}
