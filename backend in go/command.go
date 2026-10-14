package main

import (
	"fmt"
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
func dispatch(s store.Storer, c Command) error {
	// Challenge #1: Implement the dispatcher, it receives a command and
	// parses it to call the appropriate store method.
	// ex:
	// command -> {Op: "SET", Key: "env", Value: "production"},
	// it's a SET then call s.Set("env", "production")
	switch c.Op {
	case "SET":
		return s.Set(c.Key, c.Value)
	case "INCR":
		_, err := s.Incr(c.Key)
		return err
	default:
		return fmt.Errorf("unknown operation %q encountered", c.Op)
	}
}

// RestoreOnBoot replays a slice of commands against the store.
func RestoreOnBoot(s store.Storer, cmds []Command) {
	var wg sync.WaitGroup
	errs := make(chan error, len(cmds))
	for _, c := range cmds {
		wg.Go(func() {
			errs <- dispatch(s, c)
		})
	}
	go func() {
		wg.Wait()
		close(errs)
	}()
	// if we didnt close the channel, the loop below (err := range errs) would go forever; deadlock
	// we want to close it after all of the goroutines are done with their stuff

	// lets log the errors we get from a channel for now
	for err := range errs {
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}

}
