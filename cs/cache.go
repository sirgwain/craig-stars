package cs

// simple generic cache with a load function for lazy loading the cache
type cache[K comparable, V any] struct {
	data map[K]V
	load func(k K) V
}

func newCache[K comparable, V any](loader func(K) V) cache[K, V] {
	return cache[K, V]{
		data: map[K]V{},
		load: loader,
	}
}

// get retrieves a value from the cache, or loads it if it's not there
func (c cache[K, V]) get(key K) V {
	if val, ok := c.data[key]; ok {
		return val
	} else {
		val = c.load(key)
		c.data[key] = val
		return val
	}
}

// A partCache is a cache storing HullComponents for each combination of TechTag and HullSlotType.
type partCache struct {
	cache[partCacheKey, *TechHullComponent]
}

type partCacheKey struct {
	slotType HullSlotType
	techTag  TechTag
}

// create a new partCache using the specified loadFunc.
func newPartCache(loadFunc func(hst HullSlotType, tag TechTag) *TechHullComponent) partCache {
	// Wrap the provided function to allow using a single struct as a key
	wrappedLoadFunc := func(key partCacheKey) *TechHullComponent {
		return loadFunc(key.slotType, key.techTag)
	}

	cache := newCache(wrappedLoadFunc)
	pc := partCache{cache}

	return pc
}

// get retrieves a part from this partCache, or loads it if not present.
func (pc partCache) get(hst HullSlotType, tag TechTag) *TechHullComponent {
	return pc.cache.get(partCacheKey{slotType: hst, techTag: tag})
}
