package my_errgroup

import (
	"context"
	"sync"
)

type MyGroup struct {
	cancel context.CancelFunc
	ctx    context.Context
	wg     sync.WaitGroup
	mu     sync.Mutex
	err    error
}

func WithContext(ctx context.Context) (*MyGroup, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &MyGroup{cancel: cancel, ctx: ctx}, ctx
}

func (g *MyGroup) Go(f func() error) {
	g.add(f)
}

func (g *MyGroup) add(f func() error) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()

		err := f()
		if err != nil {
			g.mu.Lock()
			if g.err == nil {
				g.err = err
			}
			g.mu.Unlock()
			g.cancel()
		}
	}()
}

func (g *MyGroup) Wait() error {
	g.wg.Wait()
	return g.err
}
