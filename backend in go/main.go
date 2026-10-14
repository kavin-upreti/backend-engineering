package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("firing 3 concurrent requests...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4) // try *2, will show deadline exceeded
	defer cancel()

	// WITHOUT CONTEXT AT ALL:
	// the issue would be that this is fine if all 3 requests take the same time
	// but if one of them takes longer then it will wait forever until it returns
	// it wont timeout
	// and we want it to timeout
	// so that is where we use context
	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			makeRequest(ctx, id)
		}(i)
	}
	wg.Wait()
	fmt.Println("all done")
}

func makeRequest(ctx context.Context, id int) {
	fmt.Printf("[req %d] making the request...\n", id)

	done := make(chan struct{})
	go func() {
		time.Sleep(3 * time.Second) // simulates a slow network call
		close(done)
	}()

	// race condition: done flag or the context timeout?
	select {
	case <-done:
		fmt.Printf("[req %d] request done!\n", id)
	case <-ctx.Done():
		fmt.Printf("[req %d] has timed-out: %v!\n", id, ctx.Err())
	}
}
