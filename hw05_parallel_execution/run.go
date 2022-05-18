package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type errCnt struct {
	counter int32
	max     int32
}

func (e *errCnt) addCnt() {
	atomic.AddInt32(&e.counter, 1)
}

func (e *errCnt) isExceeded() bool {
	if e.max < 0 {
		return false
	}
	if atomic.CompareAndSwapInt32(&e.counter, e.max, e.max) {
		return true
	}
	return false
}

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	wgWorker := &sync.WaitGroup{}
	err := &errCnt{max: int32(m)}

	stack := make(chan Task)

	for i := 0; i < n; i++ {
		wgWorker.Add(1)
		go worker(wgWorker, stack, err)
	}

	for _, t := range tasks {
		if err.isExceeded() {
			break
		}
		stack <- t
	}
	close(stack)

	wgWorker.Wait()

	if err.isExceeded() {
		return ErrErrorsLimitExceeded
	}

	return nil
}

// worker is the base goroutine launches tasks.
func worker(wg *sync.WaitGroup, stack chan Task, err *errCnt) {
	defer wg.Done()

	if err.isExceeded() {
		return
	}

	for task := range stack {
		e := task()
		if e != nil {
			if err.isExceeded() {
				return
			}
			err.addCnt()
		}
	}
}
