package main

import (
	"log"
	"store/store"
	"sync"
)

// Command is a parsed Redis-ish operation: opcode plus its arguments.
type Command struct {
	Op    string
	Key   string
	Value string
}

// dispatch is a tiny switch that turns a Command into a store call.
func dispatch(s store.Storer, c Command) {
	// Challenge #1: Implement the dispatcher, it receives a command and
	// parses it to call the appropriate store method.
	// ex:
	// command -> {Op: "SET", Key: "env", Value: "production"},
	// it's a SET then call s.Set("env", "production")
	switch c.Op {
	case "SET":
		log.Printf("SET key %s", c.Key)
		s.Set(c.Key, c.Value)
	case "GET":
		s.Get(c.Key)
	case "INCR":
		s.Incr(c.Key)
	}
}

// RestoreOnBoot replays a slice of commands against the store.
func RestoreOnBoot(s store.Storer, cmds []Command) {
	var wg sync.WaitGroup

	for _, c := range cmds {
		wg.Go(func() {
			dispatch(s, c)
		})
	}
	wg.Wait()
}
