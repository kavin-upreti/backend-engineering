package kv

import (
	"errors"
	"reflect"
	"testing"

	"store/store"
)

func TestKeys_ReturnsAllKeysSorted(t *testing.T) {
	store := NewStore(3)
	store.Set("charlie", "3")
	store.Set("alpha", "1")
	store.Set("bravo", "2")

	got := store.Keys()

	want := []string{"alpha", "bravo", "charlie"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Keys() = %v, want %v", got, want)
	}
}

func TestKeys_EmptyStore(t *testing.T) {
	store := NewStore(3)
	store.Set("kavin", "5")
	got := store.Keys()
	if len(got) != 0 {
		t.Errorf("Keys on empty store = %v, expected empty slice", got)
	}
}

func TestSetGet_RoundTrip(t *testing.T) {
	store := NewStore(3)
	store.Set("hello", "world")

	// fail if we get an error for existing key or get a wrong value
	if val, err := store.Get("hello"); err != nil || val != "world" {
		t.Errorf("Get(%q) = %q, %v; want %q, nil", "hello", val, err, "world")
	}

	// fail when we dont get an error for missing key or dont get an empty value
	if val, err := store.Get("missing"); err == nil || val != "" {
		t.Errorf("Get(%q) = %q, %v; want %q, an error", "missing", val, err, "")
	}
}

func TestGetSet_EmptyKey(t *testing.T) {
	mystore := NewStore(2)

	if _, err := mystore.Get(""); err == nil || !errors.Is(err, store.ErrEmptyKey) {
		t.Error("Get() failed")
	}
}

// Table driven test
func TestDelete(t *testing.T) {
	tests := []struct {
		name      string
		setup     map[string]string
		deleteKey string
		wantLen   int
	}{
		{"deletes existing key", map[string]string{"a": "1", "b": "2"}, "a", 1},
		{"no-op on missing key", map[string]string{"a": "1"}, "x", 1},
		{"no-op on empty store", map[string]string{}, "x", 0},
		{"deletes last remaining key", map[string]string{"only": "value"}, "only", 0},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			store := NewStore(2)
			for k, v := range tc.setup {
				store.Set(k, v)
			}
			store.Delete(tc.deleteKey)
			if got := len(store.Keys()); got != tc.wantLen {
				t.Errorf("after Delete(%q): Len() = %d, want %d", tc.deleteKey, got, tc.wantLen)
			}
			if _, error := store.Get(tc.deleteKey); error == nil {
				t.Errorf("Key: (%q), still present after delete", tc.deleteKey)
			}
		})
	}
}
