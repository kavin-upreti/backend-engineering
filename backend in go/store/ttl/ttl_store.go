package ttl

import (
	"fmt"
	"time"

	"store/store"
)

// =============================================================================
// TTLStore — a second engine that expires keys automatically
// =============================================================================
// TTLStore wraps a map but adds an expiry time to every value.
// Get returns ErrKeyNotFound for expired keys, exactly like a missing key.

type TTLStore struct {
	data map[string]ttlEntry
	ttl  time.Duration
}

// key "a" -> 5 seconds

type ttlEntry struct {
	value     string
	expiresAt time.Time
}

func NewTTLStore(ttl time.Duration) *TTLStore {
	return &TTLStore{
		data: make(map[string]ttlEntry),
		ttl:  ttl,
	}
}

func (t *TTLStore) Set(key string, value string) error {
	if key == "" {
		return store.ErrEmptyKey
	}

	// Challenge: Insert the ttlEntry (calculate the time to expiry based on the ttl duration)
	exptime := time.Now().Add(t.ttl)
	t.data[key] = ttlEntry{
		value:     value,
		expiresAt: exptime,
	}
	return nil
}

func (t *TTLStore) Get(key string) (string, error) {
	if key == "" {
		return "", store.ErrEmptyKey
	}

	entry, ok := t.data[key]

	if !ok || time.Now().After(entry.expiresAt) {
		// lazy eviction - no background job scanning for expired keys
		// deleted on demand.
		delete(t.data, key)

		return "", fmt.Errorf("key %q does not exist", key)
	}

	return entry.value, nil
}

func (t *TTLStore) Delete(key string) { delete(t.data, key) }
func (t *TTLStore) Len() int          { return len(t.data) }
