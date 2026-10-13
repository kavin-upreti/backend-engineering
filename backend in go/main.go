package main

import (
	"encoding/base64"
	"fmt"

	"store/store"
	"store/store/kv"
)

func main() {
	kvStore := kv.NewStore(20)
	PopulateDefaults(kvStore)

	newstore := kvStore.Clone()
	newstore.Set("env", "production")

	fmt.Println(*kvStore)
	fmt.Println(*newstore)
	// this updated both to production, even kvstore and newstore
	// because the clone returns a pointer as it is a map
	// and a map is default a pointer

}

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

func SetKeyWithEncryption(store store.Storer, key, val string) (string, error) {
	encoded := base64.StdEncoding.EncodeToString([]byte(val))
	if err := store.Set(key, encoded); err != nil {
		return "", err
	}
	return store.Get(key)
}

// func CreateStore(maxSize int) store.Storer {
// 	plain := kv.NewStore(maxSize)
// 	logger := NewLoggingMiddleware(plain)
// 	metrics := NewMetricsMiddleware(logger)
// 	return metrics
// }

// func SetKeyWithEncryption(store store.Storer, key, val string) (string, error) {
// 	encoded := base64.StdEncoding.EncodeToString([]byte(val))
// 	if err := store.Set(key, encoded); err != nil {
// 		return "", err
// 	}

// 	store.Len()
// 	store.Get("key")

// 	return store.Get(key)
// }
