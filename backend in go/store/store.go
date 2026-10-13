package store

import (
	"errors"
)

var ErrKeyDNE = errors.New("Key does not exist")
var ErrEmptyKey = errors.New("Blank key is unacceptable")
var ErrMaxCapacity = errors.New("Maximum capacity of store has been reached")

type Storer interface {
	Get(key string) (string, error)
	Set(key, value string) error
	// Keys() []string
	Delete(key string)
	Len() int
}
