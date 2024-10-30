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
	cache.mutex.Lock()
	val, ok := cache.items[key]
	cache.mutex.Unlock()
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
			cache.mutex.Unlock()
			return false
		}
		newItem := cache.queue.PushFront(newCacheItem)
		cache.items[key] = newItem
		cache.mutex.Unlock()
	}
	return ok
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	cache.mutex.Lock()
	val, ok := cache.items[key]
	if ok {
		cache.queue.MoveToFront(val)
		typedItem, ok := val.Value.(cacheItem)
		if ok {
			cache.queue.PushFront(val)
			cache.queue.Remove(val)

			cache.mutex.Unlock()
			return typedItem.value, ok
		}
		cache.mutex.Unlock()
		return nil, false
	}
	cache.mutex.Unlock()
	return nil, false
}

func (cache *lruCache) Clear() {
	cache.mutex.Lock()
	cache.queue = NewList()
	cache.items = make(map[Key]*ListItem, cache.capacity)
	cache.mutex.Unlock()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		mutex:    sync.Mutex{},
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
