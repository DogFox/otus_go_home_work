package hw04lrucache

type Key string

type cacheItem struct {
	key   Key
	value interface{}
}

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (cache *lruCache) Set(key Key, value interface{}) bool {
	// ключ может быть, тогда обновим

	newCacheItem := cacheItem{key: key, value: value}
	val, ok := cache.items[key]
	if ok {
		val.Value = newCacheItem
		cache.queue.MoveToFront(val)
	} else {
		if cache.queue.Len() == cache.capacity {
			lastItem := cache.queue.Back()

			typedItem, ok := lastItem.Value.(cacheItem)
			if ok {
				delete(cache.items, typedItem.key)
				cache.queue.Remove(lastItem)
			}
			return false
		}
		newItem := cache.queue.PushFront(newCacheItem)
		cache.items[key] = newItem
	}
	return ok
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	val, ok := cache.items[key]
	if ok {
		cache.queue.MoveToFront(val)
		typedItem, ok := val.Value.(cacheItem)
		if ok {
			cache.queue.PushFront(val)
			cache.queue.Remove(val)
			return typedItem.value, ok
		}
		return nil, false
	} else {
		return nil, false
	}
}

func (cache *lruCache) Clear() {
	cache.queue = NewList()
	cache.items =
		make(map[Key]*ListItem, cache.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
