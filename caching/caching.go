package caching

// InMemoryCache uses built-in map type as a storage for cache.
type InMemoryCache struct {
	cache map[string]struct{}
}

// NewInMemoryCache returns initialized in-memory cache.
func NewInMemoryCache() InMemoryCache {
	return InMemoryCache{
		cache: make(map[string]struct{}),
	}
}

// Visit returns true if provided string was visited before.
func (c InMemoryCache) Visit(s string) bool {
	if _, ok := c.cache[s]; ok {
		return true
	}
	c.cache[s] = struct{}{}
	return false
}
