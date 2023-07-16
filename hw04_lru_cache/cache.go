package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool // Добавить значение в кэш по ключу.
	Get(key Key) (interface{}, bool)     // Получить значение из кэша по ключу
	Clear()                              // Очистить кэш
}

var _ Cache = (*lruCache)(nil)

type lruCache struct {
	sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type pair struct {
	key   Key
	value interface{}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.Lock()
	defer c.Unlock()

	if item, ok := c.items[key]; ok {
		c.queue.MoveToFront(item)
		item.Value = pair{key: key, value: value}
		return true
	}

	newPair := pair{key: key, value: value}
	item := c.queue.PushFront(newPair)
	c.items[key] = item
	if c.queue.Len() > c.capacity {
		lastElem := c.queue.Back()
		delete(c.items, lastElem.Value.(pair).key)
		c.queue.Remove(lastElem)
	}

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.Lock()
	defer c.Unlock()

	if item, ok := c.items[key]; ok {
		c.queue.MoveToFront(item)
		return item.Value.(pair).value, true
	}

	return nil, false
}

func (c *lruCache) Clear() {
	c.Lock()
	defer c.Unlock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
