package cache

import (
	"container/list"
	"sync"
	"time"
)

type LRUCache struct {
	list     *list.List
	hash     map[string]*list.Element
	mu       sync.Mutex
	capacity uint64
	size     uint64
}

type entry struct {
	key      string
	value    string
	size     int
	lastUsed time.Time
}

func NewLRUCache(capacity uint64) *LRUCache {
	cache := &LRUCache{
		capacity: capacity,
		list:     list.New(),
		hash:     make(map[string]*list.Element),
		size:     0,
	}

	return cache
}

func (c *LRUCache) Set(key string, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if e := c.hash[key]; e != nil {
		c.update(e, value)
	} else {
		c.add(key, value)
	}
	c.ensureCapacity()
}

func (c *LRUCache) Get(key string) string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if e := c.hash[key]; e != nil {
		c.moveToFront(e)
		return e.Value.(*entry).value
	}

	return ""
}

func (c *LRUCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if e := c.hash[key]; e != nil {
		c.capacity -= uint64(e.Value.(*entry).size)
		c.list.Remove(e)
		delete(c.hash, key)
	}
}

func (c *LRUCache) Clear() {
	c.list = list.New().Init()
	c.hash = make(map[string]*list.Element)
	c.size = 0
}

func (c *LRUCache) ensureCapacity() {
	for c.size > c.capacity {
		lastElem := c.list.Back()
		lastValue := lastElem.Value.(*entry)
		c.list.Remove(lastElem)
		delete(c.hash, lastValue.key)
		c.size -= uint64(lastValue.size)
	}
}

func (c *LRUCache) add(key string, value string) {
	e := &entry{
		key:      key,
		value:    value,
		size:     len(value),
		lastUsed: time.Now().UTC(),
	}
	listElement := c.list.PushFront(e)
	c.hash[key] = listElement
	c.size += uint64(e.size)
}

func (c *LRUCache) update(e *list.Element, value string) {
	size := len(value)
	sizeDiff := size - e.Value.(*entry).size
	e.Value.(*entry).value = value
	e.Value.(*entry).size = size
	c.size += uint64(sizeDiff)
}

func (c *LRUCache) moveToFront(e *list.Element) {
	e.Value.(*entry).lastUsed = time.Now().UTC()
	c.list.MoveToFront(e)
}
