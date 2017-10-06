package main

import (
	"testing"

	"github.com/msawangwan/ci.io/lib/cache"
)

var (
	wsCacheMock cache.KVStorer
)

func TestCreateTempWorkspace(t *testing.T) {
	wsCacheMock = cache.NullKVStore
}
