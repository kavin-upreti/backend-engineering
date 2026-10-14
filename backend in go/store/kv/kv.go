package kv

import (
	"fmt"
	"sort"
	"time"

	"store/store"
)

type Store struct {
	data    map[string]string
	maxSize int // 0 = unlimited
}

func (s *Store) Clone() *Store {
	// return &Store{
	// 	data:    s.data,
	// 	maxSize: s.maxSize,
	// }
	// if i returned this, updating the clone would update the original too as data = s.data is a map
	// and map is address by default

	cp := &Store{
		data:    make(map[string]string, len(s.data)),
		maxSize: s.maxSize,
	}

	for k, v := range s.data {
		cp.data[k] = v
	}

	return cp

	/*
		kvStore := kv.NewStore(20)
		PopulateDefaults(kvStore)

		newstore := kvStore.Clone()
		newstore.Set("env", "production")

		fmt.Println(*kvStore)
		fmt.Println(*newstore)
		// this would update both to production, kvstore and newstore
		// because the clone returns a pointer as it is a map
		// and a map is default a pointer
		// that is why we updated our Clone func

		this was in main.go earlier

		func PopulateDefaults(s store.Storer) error {
			defaults := map[string]string{
				"env":     "development",
				"version": "0.0.1",
				"debug":   "true",
			}
			for k, v := range defaults {
				if err := s.Set(k, v); err != nil {
					return fmt.Errorf("PopulateDefaults: %w", err)
				}
			}
			return nil
		}
	*/
}

// returns every key in the store sorted alphabetically
func (s *Store) Keys() []string {
	// first we wrote the test in store_test.go
	// then we do it here
	keys := make([]string, 0, len(s.data))
	for key := range s.data {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func NewStore(maxSize int) *Store {
	return &Store{
		data:    make(map[string]string),
		maxSize: maxSize,
	}
}

func (s *Store) Get(key string) (string, error) {
	if key == "" {
		return "", store.ErrEmptyKey
	}

	val, ok := s.data[key]

	if !ok {
		return "", fmt.Errorf("Key %q does not exist", key)
	}
	return val, nil
}

func (s *Store) Set(key, value string) error {
	if key == "" {
		return store.ErrEmptyKey
	}
	// this is for the goroutines example, uncomment to test
	time.Sleep(time.Second)
	_, exists := s.data[key]
	if s.maxSize > 0 && s.Len() >= s.maxSize && !exists {
		return fmt.Errorf("Set(%q): %w", key, store.ErrMaxCapacity)
	}

	s.data[key] = value
	return nil
}

func (s *Store) Delete(key string) {
	delete(s.data, key)
}

func (s *Store) Len() int {
	return len(s.data)
}
