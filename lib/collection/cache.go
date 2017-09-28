package collection

type Cache interface {
	Valid() bool
	Contains(k int) bool
	Store(v interface{}) int
	Fetch(k int) (v interface{}, e error)
}

type BasicCache struct {
	store map[int]string
	key   int
}

func newBasicCache() Cache {
	return &BasicCache{
		store: make(map[int]string),
		key:   0,
	}
}

func (c *BasicCache) Valid() bool { return true }

func (c *BasicCache) Contains(k int) bool {
	if _, contains := c.store[k]; contains {
		return true
	}
	return false
}

func (c *BasicCache) Store(v interface{}) int {
	c.key++
	c.store[c.key] = v.(string)
	return c.key
}

func (c *BasicCache) Fetch(k int) (v interface{}, e error) { return }
