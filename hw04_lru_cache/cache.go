package hw04lrucache

import "sync"

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
	mutex    sync.Mutex
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
		cache.mutex.Lock()
		cache.queue.MoveToFront(val)
		cache.mutex.Unlock()
	} else {
		cache.mutex.Lock()
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
		cache.mutex.Unlock()

	}
	return ok
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	val, ok := cache.items[key]
	if ok {
		cache.mutex.Lock()
		cache.queue.MoveToFront(val)
		cache.mutex.Unlock()
		typedItem, ok := val.Value.(cacheItem)
		if ok {
			cache.mutex.Lock()
			cache.queue.PushFront(val)
			cache.queue.Remove(val)
			cache.mutex.Unlock()
			return typedItem.value, ok
		}
		return nil, false
	} else {
		return nil, false
	}
}

func (cache *lruCache) Clear() {
	cache.mutex.Lock()
	cache.queue = NewList()
	cache.items =
		make(map[Key]*ListItem, cache.capacity)
	cache.mutex.Unlock()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
