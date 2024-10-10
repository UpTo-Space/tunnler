package cache

import (
	"bytes"
	"container/list"
	"fmt"
	"testing"
)

func TestLRUCreation(t *testing.T) {
	tests := []struct {
		capacity int
	}{
		{0},
		{-1},
		{5},
		{100000},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Capacity %v", tt.capacity), func(t *testing.T) {
			lru := NewLRUCache(uint64(tt.capacity))

			if lru.capacity != uint64(tt.capacity) && lru.size != 0 {
				t.Errorf("LRU isn't initialized correctly")
			}
		})
	}
}

func TestLRUAddGet(t *testing.T) {
	lru := NewLRUCache(10)

	type kv struct {
		Key   string
		Value string
	}

	type testData struct {
		name     string
		capacity int
		KV       []kv
		Expected []kv
	}

	tests := []testData{
		testData{
			name:     "Full Capacity",
			capacity: 5,
			KV: []kv{
				kv{Key: "A", Value: "a"},
				kv{Key: "B", Value: "b"},
				kv{Key: "C", Value: "c"},
				kv{Key: "D", Value: "d"},
				kv{Key: "E", Value: "e"},
			},
			Expected: []kv{
				kv{Key: "A", Value: "a"},
				kv{Key: "B", Value: "b"},
				kv{Key: "C", Value: "c"},
				kv{Key: "D", Value: "d"},
				kv{Key: "E", Value: "e"},
			},
		},
		testData{
			capacity: 3,
			KV: []kv{
				kv{Key: "A", Value: "a"},
				kv{Key: "B", Value: "b"},
				kv{Key: "C", Value: "c"},
				kv{Key: "D", Value: "d"},
				kv{Key: "E", Value: "e"},
			},
			Expected: []kv{
				kv{Key: "C", Value: "c"},
				kv{Key: "D", Value: "d"},
				kv{Key: "E", Value: "e"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tsize := 0
			for _, v := range tt.KV {
				lru.Set(v.Key, v.Value)
				tsize += len(v.Value)
			}

			if lru.size != uint64(tsize) {
				t.Errorf("size of LRU (%v) doesn't match expected size of %v", lru.size, tsize)
			}

			for _, e := range tt.KV {
				v := lru.Get(e.Key)
				if e.Value != v {
					t.Errorf("value of LRU key %v returned %v. expected %v", e.Key, e.Value, v)
				}

				if e := lru.hash[e.Key]; e == nil {
					t.Errorf("list element with key %v shouldn't be nil", e.Value)
				} else {
					if !elementShouldBeFirst(e) {
						t.Errorf("list element with key %v should be first in the list", e.Value)
					}
				}
			}
		})
	}
}

func TestLRUCapacity(t *testing.T) {
	type testData []struct {
		key          string
		value        string
		expectedKeys []string
	}
	type testsData struct {
		data     testData
		capacity int
	}

	tests := []testsData{
		testsData{
			capacity: 3,
			data: testData{
				{"A", "a", []string{"A"}},
				{"B", "b", []string{"A", "B"}},
				{"C", "c", []string{"A", "B", "C"}},
				{"D", "d", []string{"B", "C", "D"}},
				{"E", "e", []string{"C", "D", "E"}},
			},
		},
		testsData{
			capacity: 1,
			data: testData{
				{"A", "a", []string{"A"}},
				{"B", "b", []string{"B"}},
				{"C", "c", []string{"C"}},
				{"D", "d", []string{"D"}},
				{"E", "e", []string{"E"}},
			},
		},
		testsData{
			capacity: 5,
			data: testData{
				{"A", "a", []string{"A"}},
				{"B", "b", []string{"A", "B"}},
				{"C", "c", []string{"A", "B", "C"}},
				{"D", "d", []string{"A", "B", "C", "D"}},
				{"E", "e", []string{"A", "B", "C", "D", "E"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Capacity %v", tt.capacity), func(t *testing.T) {
			lru := NewLRUCache(uint64(tt.capacity))
			for _, d := range tt.data {
				lru.Set(d.key, d.value)
				if lru.list.Len() > int(lru.capacity) {
					t.Errorf("the LRU shouldn't be bigger than each individial character added. capacity: %v, actual length: %v", lru.capacity, lru.list.Len())
				}

				for _, e := range d.expectedKeys {
					if el := lru.hash[e]; el == nil {
						t.Errorf("element should exist but doesn't: %v", e)
					}
				}
			}
		})
	}
}

func TestLRUClear(t *testing.T) {
	lru := NewLRUCache(5)

	data := []struct {
		key   string
		value string
	}{
		{"A", "a"},
		{"B", "b"},
		{"C", "c"},
		{"D", "d"},
		{"E", "e"},
	}

	for _, dt := range data {
		lru.add(dt.key, dt.value)
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

	data := []struct {
		key   string
		value string
	}{
		{"A", "a"},
		{"B", "b"},
		{"C", "c"},
		{"D", "d"},
		{"E", "e"},
	}

	for _, dt := range data {
		lru.add(dt.key, dt.value)
	}

	if lru.size <= lru.capacity {
		t.Errorf("LRU capacity should be exceeded by now. actual size: %v and capacity: %v", lru.size, lru.capacity)
	}

	if len(lru.hash) != len(data) || lru.list.Len() != len(data) {
		t.Error("LRU map and list should have the same length as the tests")
	}

	lru.ensureCapacity()

	if lru.size != lru.capacity {
		t.Error("LRU size should be limited to capacity now")
	}
}

func TestLRUCacheSerialization(t *testing.T) {
	lru := NewLRUCache(5)
	data := []struct {
		key   string
		value string
	}{
		{"A", "a"},
		{"B", "b"},
		{"C", "c"},
		{"D", "d"},
		{"E", "e"},
	}

	for _, dt := range data {
		lru.Set(dt.key, dt.value)
	}

	bb := &bytes.Buffer{}
	if err := lru.saveCache(bb); err != nil {
		t.Errorf("error in saving cache %v", err)
	}

	if b := bb.Bytes(); b == nil || len(string(b)) == 0 {
		t.Error("bytes of buffer seems to be empty")
	}

	if d, err := lru.loadFromReader(bb); err != nil {
		t.Errorf("error loading LRU from byte buffer: %v", err)
	} else {
		if len(d) != len(lru.hash) {
			t.Errorf("length of deserialized LRU is different. Deserialized %v, Expected %v", len(d), len(lru.hash))
		}
	}
}

func elementShouldBeFirst(e *list.Element) bool {
	return e.Prev() == nil
}
