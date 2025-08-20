package pokecache

import (
	"fmt"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	cache := NewCache(5 * time.Millisecond)
	expected := []byte("Hello")
	cache.Add("key", expected)
	actual, ok := cache.Get("key")
	if !ok {
		t.Errorf("Cache entry was not added, or was cleaned immediately")
		t.Fail()
		return
	}

	if len(actual) != len(expected) {
		t.Errorf("Cache value doesn't match the added value\n Expected: %v\n Actual %v", expected, actual)
		t.Fail()
		return
	}
	for i := range actual {
		if actual[i] != expected[i] {
			t.Errorf("Cache value doesn't match the added value\n Expected: %v\n Actual %v", expected, actual)
			t.Fail()
			return
		}
	}
	time.Sleep(50 * time.Millisecond)
	_, ok = cache.Get("key")
	if ok {
		t.Errorf("Failed to clean cache after timeout")
		t.Fail()
		return
	}
}

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond
	cache := NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}
