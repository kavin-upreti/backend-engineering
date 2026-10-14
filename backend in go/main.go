package main

import (
	"fmt"

	"store/store"
	"store/store/kv"
)

// in sequential prog, inserting 50 keys where 1 key takes an avg 1 sec would take 50s
// but in concurrent prog (goroutines), inserting 50 keys would take same time as the key taking longest
// to make a func concurrent, define it as "go func()" instead of "func()"
func main() {
	cmds := []Command{
		{Op: "SET", Key: "env", Value: "production"},
		{Op: "SET", Key: "version", Value: "0.0.1"},
		{Op: "SET", Key: "debug", Value: "true"},
		{Op: "GET", Key: "env"},
		{Op: "SET", Key: "region", Value: "eu-west-1"},
		{Op: "GET", Key: "version"},
	}

	s := kv.NewStore(0)
	RestoreOnBoot(s, cmds)

	fmt.Println("restored, store size:", s.Len())
	fmt.Println("keys:", s.Keys())
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

// func SetKeyWithEncryption(store store.Storer, key, val string) (string, error) {
// 	encoded := base64.StdEncoding.EncodeToString([]byte(val))
// 	if err := store.Set(key, encoded); err != nil {
// 		return "", err
// 	}
// 	return store.Get(key)
// }

// func CreateStore(maxSize int) store.Storer {
// 	plain := kv.NewStore(maxSize)
// 	logger := NewLoggingMiddleware(plain)
// 	metrics := NewMetricsMiddleware(logger)
// 	return metrics
// }

// 	store.Len()
// 	store.Get("key")

// 	return store.Get(key)
// }
