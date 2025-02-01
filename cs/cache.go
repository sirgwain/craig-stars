package cs

// simple generic cache with a load function for lazy loading the cache
type cache[K comparable, V any] struct {
	data map[K]V
	load func(k K) V
}

func newCache[K comparable, V any](loader func(k K) V) cache[K, V] {
	return cache[K, V]{
		data: map[K]V{},
		load: loader,
	}
}

// get gets a value from the cache, or loads it if it's not there
func (c *cache[K, V]) get(key K) V {
	if val, ok := c.data[key]; ok {
		return val
	} else {
		val = c.load(key)
		c.data[key] = val
		return val
	}
}

