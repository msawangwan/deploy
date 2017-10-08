package cache

// KVStorer ...
type KVStorer interface {
	Store(k, v string)
	Fetch(k string) (v string, e error)
	Flush() (n int, e error)
}

type nullKVStore struct{}

// Store ...
func (nu nullKVStore) Store(k, v string) {}

// Fetch ...
func (nu nullKVStore) Fetch(k string) (v string, e error) { return }

// Flush ....
func (nu nullKVStore) Flush() (n int, e error) { return }

// NullKVStore is a global
var NullKVStore = nullKVStore{}
