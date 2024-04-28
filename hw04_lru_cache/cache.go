package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool // Добавить значение в кэш по ключу.
	Get(key Key) (interface{}, bool)     // Получить значение из кэша по ключу.
	Clear()                              // Очистить кэш.
}

type lruCache struct {
	capacity int               // ёмкость (количество сохраняемых в кэше элементов)
	queue    List              // очередь [последних используемых элементов] на основе двусвязного списка
	items    map[Key]*ListItem // словарь, отображающий ключ на элемент очереди
	mx       sync.Mutex
}

type entry struct {
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

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mx.Lock()
	defer c.mx.Unlock()

	item, ok := c.items[key]
	if ok {
		entryFromItem(item).value = value
		c.queue.MoveToFront(item)
		return true
	}

	if c.queue.Len() == c.capacity {
		backItem := c.queue.Back()
		c.queue.Remove(backItem)
		delete(c.items, entryFromItem(backItem).key)
	}

	createdItem := c.queue.PushFront(&entry{key: key, value: value})
	c.items[key] = createdItem

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mx.Lock()
	defer c.mx.Unlock()

	item, ok := c.items[key]
	if !ok {
		return nil, false
	}

	c.queue.MoveToFront(item)

	return entryFromItem(item).value, true
}

func (c *lruCache) Clear() {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func entryFromItem(item *ListItem) *entry {
	return item.Value.(*entry)
}
