package main

import (
	"log"
	"store/store"
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
	}
}

// RestoreOnBoot replays a slice of commands against the store.
func RestoreOnBoot(s store.Storer, cmds []Command) {
	for _, c := range cmds {
		// without go before dispatch, it would take 4 seconds to happen
		// because of the time.Sleep in set
		// but with go, it takes 1, however the log prints dont happen
		// because main got to the end before the goroutines were even scheduled
		// so they never reached the log.Printf
		go dispatch(s, c)
	}

	// if I add a time.Sleep(time.Second * 2) here, it will reach the log printf
	// however in that case, all 4 goroutines reached m[key] = value at the same time
	// so rather than corrupting the map, go crashes the program
	// time.Sleep(time.Second * 2):
	/*
		2026/10/14 10:19:04 SET key debug
		2026/10/14 10:19:04 SET key env
		2026/10/14 10:19:04 SET key version
		2026/10/14 10:19:04 SET key region
		fatal error: concurrent map writes

		goroutine 21 [running]:
		internal/runtime/maps.fatal({0x100665ba5?, 0x10061a6b8?})
				/usr/local/go/src/runtime/panic.go:1181 +0x20
		store/store/kv.(*Store).Set(0x49d0a6290040, {0x100662d91, 0x7}, {0x100662b13, 0x5})
				/Users/kavin/Desktop/Learning/Backend Engineering/backend in go/store/kv/kv.go:110 +0x13c
		main.dispatch({0x100741110, 0x49d0a6290040}, {{0x1006629a4, 0x3}, {0x100662d91, 0x7}, {0x100662b13, 0x5}})
				/Users/kavin/Desktop/Learning/Backend Engineering/backend in go/command.go:26 +0x110
		created by main.RestoreOnBoot in goroutine 1
				/Users/kavin/Desktop/Learning/Backend Engineering/backend in go/command.go:40 +0x58
		exit status 2
	*/
}
