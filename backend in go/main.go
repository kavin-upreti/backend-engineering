package main

import (
	"fmt"
	"sync"
	"time"
)

// in sequential prog, inserting 50 keys where 1 key takes an avg 1 sec would take 50s
// but in concurrent prog (goroutines), inserting 50 keys would take same time as the key taking longest
// to make a func concurrent, define it as "go func()" instead of "func()"
var (
	counter int
	mu      sync.Mutex // mutual exclusion
)

func increment() {
	mu.Lock()
	defer mu.Unlock() // just execute at the end, write it here so you dont forget
	v := counter
	time.Sleep(time.Microsecond) // simulate a db query
	counter = v + 1

	// mu.Unlock()
}

func main() {
	expected := 1000
	var wg sync.WaitGroup
	for i := 0; i < expected; i++ {
		wg.Add(1)
		go func() {
			increment()
			defer wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("Counter: ", counter)
	fmt.Println("Expected: ", expected)

	// WITHOUT MU AND MU.LOCK() AND DEFER MU.UNLOCK()
	// this would have very varied values of counter and expected
	// this is because of race conditions: multiple threads are trying to incrmeent the same var at the same time
	// and so a lot of it is overriden:
	// thread 1 tried to get counter from 5 to 6 and did it
	// but thread 2 tried it at the same time too
	// and so instead of it going from 6 to 7 due to thread 2, it will just stay from 5 to 6
	// go can solve this though using .lock
	// so you hold the lock during the entire duration of you writing the value
	// so every other thread waits in line

	// but here since we have used those things we get the desired output (1000 and 1000 both)
}
