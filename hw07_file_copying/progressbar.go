package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type ProgressBar struct {
	wg      sync.WaitGroup
	current int64
	count   int64
	stack   chan int64
}

// NewProgressBar inits Bar count process.
func NewProgressBar(count int64) *ProgressBar {
	pb := new(ProgressBar)
	pb.stack = make(chan int64)

	pb.count = count
	pb.wg.Add(1)

	go func(pb *ProgressBar) {
		defer pb.wg.Done()
		for current := range pb.stack {
			if pb.count > 0 {
				v := 100 * current / pb.count
				fmt.Printf("\r%v%%", v)
			} else {
				v := current
				fmt.Printf("\r%v", v)
			}
		}
	}(pb)

	return pb
}

// Add increases bar counter on delta value.
func (pb *ProgressBar) Add(delta int64) {
	atomic.AddInt64(&pb.current, delta)
	pb.stack <- pb.current
}

// Finish closes counting process.
func (pb *ProgressBar) Finish() {
	close(pb.stack)
	pb.wg.Wait()
	fmt.Println()
}
