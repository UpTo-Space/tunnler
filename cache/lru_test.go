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

func elementShouldBeFirst(e *list.Element) bool {
	return e.Prev() == nil
}
