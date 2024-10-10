package cache

import (
	"container/list"
	"testing"
)

func TestLRUCreation(t *testing.T) {
	lru := NewLRUCache(100)

	if lru.capacity != 100 && lru.size != 0 {
		t.Errorf("LRU isn't initialized correctly")
	}
}

func TestLRUAddGet(t *testing.T) {
	lru := NewLRUCache(10)

	tests := []struct {
		key   string
		value string
	}{
		{"A", "a"},
		{"B", "b"},
		{"longer", "long"},
	}

	tsize := 0
	for _, tt := range tests {
		lru.Set(tt.key, tt.value)
		tsize += len(tt.value)
		if lru.size != uint64(tsize) {
			t.Errorf("size of LRU (%v) doesn't match expected size of %v", lru.size, tsize)
		}

		v := lru.Get(tt.key)
		if tt.value != v {
			t.Errorf("value of LRU key %v returned %v. expected %v", tt.key, tt.value, v)
		}

		if e := lru.hash[tt.key]; e == nil {
			t.Errorf("list element with key %v shouldn't be nil", tt.key)
		} else {
			if !elementShouldBeFirst(e) {
				t.Errorf("list element with key %v should be first in the list", tt.key)
			}
		}
	}
}

func TestLRUCapacity(t *testing.T) {
	lru := NewLRUCache(3)

	tests := []struct {
		key          string
		value        string
		expectedKeys []string
	}{
		{"A", "a", []string{"A"}},
		{"B", "b", []string{"A", "B"}},
		{"C", "c", []string{"A", "B", "C"}},
		{"D", "d", []string{"B", "C", "D"}},
		{"E", "e", []string{"C", "D", "E"}},
	}

	for _, tt := range tests {
		lru.Set(tt.key, tt.value)
		if lru.list.Len() > int(lru.capacity) {
			t.Errorf("the LRU shouldn't be bigger than each individial character added. capacity: %v, actual length: %v", lru.capacity, lru.list.Len())
		}

		for _, e := range tt.expectedKeys {
			if el := lru.hash[e]; el == nil {
				t.Errorf("element should exist but doesn't: %v", e)
			}
		}
	}
}

func TestLRUClear(t *testing.T) {
	lru := NewLRUCache(5)

	tests := []struct {
		key   string
		value string
	}{
		{"A", "a"},
		{"B", "b"},
		{"C", "c"},
		{"D", "d"},
		{"E", "e"},
	}

	for _, tt := range tests {
		lru.add(tt.key, tt.value)
	}

	lru.Clear()

	if lru.list.Front() != nil {
		t.Error("an empty LRU shouldn't have a front element")
	}

	if len(lru.hash) > 0 {
		t.Error("an empty LRU shouldn't have any element in it's hash map")
	}
}

func TestLRUCapacityClear(t *testing.T) {
	lru := NewLRUCache(2)

	tests := []struct {
		key   string
		value string
	}{
		{"A", "a"},
		{"B", "b"},
		{"C", "c"},
		{"D", "d"},
		{"E", "e"},
	}

	for _, tt := range tests {
		lru.add(tt.key, tt.value)
	}

	if lru.size <= lru.capacity {
		t.Errorf("LRU capacity should be exceeded by now. actual size: %v and capacity: %v", lru.size, lru.capacity)
	}

	if len(lru.hash) != len(tests) || lru.list.Len() != len(tests) {
		t.Error("LRU map and list should have the same length as the tests")
	}

	lru.ensureCapacity()

	if lru.size != lru.capacity {
		t.Error("LRU size should be limited to capacity now")
	}
}

func elementShouldBeFirst(e *list.Element) bool {
	return e.Prev() == nil
}
