package main

import (
	"fmt"
	"math/rand/v2"
	"time"
)

// we couldnt communicate data from the goroutines to the main code
// handling/sending data from goroutines isnt easy
// so we need channels for this
// solves the idea of sending data from goroutines to other goroutines
func main() {
	results := make(chan int, 1) // unbuffered channels, blocks program until someone reads
	// 1 means that hold 1 value(s) without anyone receiving them yet
	/*
		ch := make(chan int, 2)
		ch <- 1 // doesn't block, buffer has room
		ch <- 2 // doesn't block, buffer now full
		// ch <- 3 would block here until someone receives
		fmt.Println(<-ch, <-ch) // 1 2
	*/

	go worker(results)
	go worker(results)

	fmt.Println(<-results, <-results) // to read
	// started 2 workers so can only print twice btw, if I did it thrice then deadlock hit
	// why? because first <-results waits, then worker sends smth
	// so for third <-results, it is waiting but no worker is left to send
}

func worker(results chan<- int) { // results <-chan int means you can read only, chan<- means send from worker
	// without <- anywhere its read and write
	time.Sleep(time.Second)
	// we want to communicate this thing below to our main
	results <- rand.IntN(100)
}
