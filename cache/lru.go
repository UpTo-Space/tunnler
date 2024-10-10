package cache

import (
	"container/list"
	"encoding/gob"
	"io"
	"os"
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
	Key      string
	Value    string
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
		return e.Value.(*entry).Value
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

func (c *LRUCache) Entries() []entry {
	entries := []entry{}
	for _, it := range c.hash {
		entries = append(entries, *it.Value.(*entry))
	}
	return entries
}

func (c *LRUCache) Save(cacheFile string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	file, err := os.OpenFile(cacheFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	return c.saveCache(file)
}

func (c *LRUCache) Load(cacheFile string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	file, err := os.OpenFile(cacheFile, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}

	entries, err := c.loadFromReader(file)
	if err != nil {
		return err
	}

	if len(entries) > 0 {
		c.setCapacity(len(entries))
	}

	for _, e := range entries {
		c.Set(e.Key, e.Value)
	}

	return nil
}

func (c *LRUCache) saveCache(w io.Writer) error {
	entries := c.Entries()

	enc := gob.NewEncoder(w)
	return enc.Encode(entries)
}

func (c *LRUCache) loadFromReader(r io.Reader) ([]entry, error) {
	var entries []entry

	dec := gob.NewDecoder(r)
	if err := dec.Decode(&entries); err != nil {
		return nil, err
	}
	return entries, nil
}

func (c *LRUCache) ensureCapacity() {
	for c.size > c.capacity {
		lastElem := c.list.Back()
		lastValue := lastElem.Value.(*entry)
		c.list.Remove(lastElem)
		delete(c.hash, lastValue.Key)
		c.size -= uint64(lastValue.size)
	}
}

func (c *LRUCache) setCapacity(capacity int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.capacity = uint64(capacity)
	c.ensureCapacity()
}

func (c *LRUCache) add(key string, value string) {
	e := &entry{
		Key:      key,
		Value:    value,
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
	e.Value.(*entry).Value = value
	e.Value.(*entry).size = size
	c.size += uint64(sizeDiff)
}

func (c *LRUCache) moveToFront(e *list.Element) {
	e.Value.(*entry).lastUsed = time.Now().UTC()
	c.list.MoveToFront(e)
}
