package collection

type Cache interface {
	Valid() bool
	Contains(k string) bool
	Store(k string, v interface{}) bool
	Fetch(k string) (v interface{}, e error)
}

type BasicCache map[string]string

func newBasicCache() Cache { return BasicCache(make(map[string]string)) }

func (c BasicCache) Valid() bool                             { return true }
func (c BasicCache) Contains(k string) bool                  { return true }
func (c BasicCache) Store(k string, v interface{}) bool      { return true }
func (c BasicCache) Fetch(k string) (v interface{}, e error) { return }
