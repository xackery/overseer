package signal

import (
	"context"
	"sync"
)

var (
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
)

func init() {
	ctx, cancel = context.WithCancel(context.Background())
}

// Cancel cancels the context.
func Cancel() {
	cancel()
}

// Ctx returns the context.
func Ctx() context.Context {
	return ctx
}

// WaitWorker waits for all workers to finish.
func WaitWorker() {
	wg.Wait()
}

// AddWorker adds a worker.
func AddWorker() {
	wg.Add(1)
}

// FinishWorker finishes a worker.
func FinishWorker() {
	wg.Done()
}
