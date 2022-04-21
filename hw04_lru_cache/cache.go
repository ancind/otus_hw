package hw04lrucache

import "errors"

type Key string

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

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (lc *lruCache) Set(k Key, v interface{}) bool {
	item := &cacheItem{key: k, value: v}
	if listItem, ok := lc.items[k]; ok {
		listItem.Value = item
		lc.queue.MoveToFront(listItem)

		return true
	}

	pushedItem := lc.queue.PushFront(item)
	if lc.queue.Len() > lc.capacity {
		convertedCacheItem, err := convertToCacheItem(lc.queue.Back().Value)
		if err != nil {
			return false
		}

		lc.queue.Remove(lc.queue.Back())
		delete(lc.items, convertedCacheItem.key)
	}

	lc.items[k] = pushedItem

	return false
}

func (lc *lruCache) Get(k Key) (interface{}, bool) {
	if listItem, ok := lc.items[k]; ok {
		lc.queue.MoveToFront(listItem)

		convertedCacheItem, err := convertToCacheItem(listItem.Value)
		if err != nil {
			return nil, false
		}

		return convertedCacheItem.value, true
	}

	return nil, false
}

func (lc *lruCache) Clear() {
	lc.queue = NewList()
	for item := range lc.items {
		delete(lc.items, item)
	}
}

func convertToCacheItem(value interface{}) (*cacheItem, error) {
	cItem, ok := value.(*cacheItem)
	if !ok {
		return nil, errors.New("first element is number")
	}

	return cItem, nil
}
