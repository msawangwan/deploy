package cache

import (
	"fmt"
)

// ItemNotInCacheError ...
type ItemNotInCacheError struct {
	Item string
}

// Error ...
func (e ItemNotInCacheError) Error() string {
	return fmt.Sprintf("[error][item_not_in_cache][%s]", e.Item)
}

// NewItemNotInCacheError ...
func NewItemNotInCacheError(item string) ItemNotInCacheError {
	return ItemNotInCacheError{Item: item}
}
