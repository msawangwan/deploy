package main

import (
	"testing"

	"github.com/msawangwan/ci.io/lib/cache"
)

var (
	mock_wsCache cache.KVStorer
)

func TestCreateTempWorkspace(t *testing.T) {
	mock_wsCache = cache.NullKVStore
}
