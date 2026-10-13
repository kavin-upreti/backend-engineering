package main

import (
	"encoding/base64"
	"fmt"
	"store/store"
	"store/store/kv"
	"store/store/ttl"
	"time"
)

func main() {
	kvstore := kv.NewStore(20)
	encrypted, err := SetKeyWithEncryption(kvstore, "a", "hey you cant see me lol")

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(encrypted)

	ttlStore := ttl.NewTTLStore(time.Second * 2)
	encryptedTtlValue, err := SetKeyWithEncryption(ttlStore, "a", "hey you cant see me lol")

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(encryptedTtlValue)

}

// storer is an interface so both normal store and ttlstore have it
func SetKeyWithEncryption(store store.Storer, key, val string) (string, error) {
	encoded := base64.StdEncoding.EncodeToString([]byte(val))
	if err := store.Set(key, encoded); err != nil {
		return "", err
	}
	return store.Get(key)
}
