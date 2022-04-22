package hw05parallelexecution

import (
	"errors"
	"sync"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrZeroWorkers         = errors.New("need at least one worker")
	ErrLimitErrors         = errors.New("m should be more zero")
)

type Task func() error

type errorCounter struct {
	mu    sync.Mutex
	count int
}

func (c *errorCounter) Add() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

func (c *errorCounter) Get() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrLimitErrors
	}

	if n <= 0 {
		return ErrZeroWorkers
	}

	ch := make(chan Task)
	wg := sync.WaitGroup{}
	wg.Add(n)
	errCounter := errorCounter{}

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for t := range ch {
				if t() != nil {
					errCounter.Add()
				}
			}
		}()
	}

	for _, t := range tasks {
		if errCounter.Get() < m {
			ch <- t
		}
	}
	close(ch)
	wg.Wait()

	if errCounter.Get() >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
