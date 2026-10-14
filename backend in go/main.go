package main

import (
	"fmt"
	"store/store/kv"
	"sync"
	"time"
)

// Benchmark: how long do many concurrent readers take?
//
// Run it once now, then refactor the store Lock to RLock
//   - With sync.Mutex:   1143.906 ms
//   - With sync.RWMutex: 11.629 ms
func main() {
	store := kv.NewStore(0)
	store.Set("name", "jane")

	const readers = 100
	const readsEach = 10

	start := time.Now()
	var wg sync.WaitGroup
	for range readers {
		wg.Go(func() {
			for range readsEach {
				store.Get("name") // this is what we want to optimise with RWMutex
				// with Mutex, it will take some time
				// with RWMutex it is quicker
			}
		})
	}
	wg.Wait()

	elapsed := time.Since(start)
	fmt.Printf("%d readers did %d concurrent reads in %.3fms\n",
		readers, readers*readsEach, float64(elapsed.Microseconds())/1000)
}
