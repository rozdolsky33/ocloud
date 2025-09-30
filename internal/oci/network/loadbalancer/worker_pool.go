package loadbalancer

import (
	"context"

	"golang.org/x/sync/errgroup"
)

const (
	defaultWorkerCount = 12
)

// Work represents a unit of work executed by the worker pool
// It returns an error to allow early cancellation via errgroup.
type Work func() error

// runWithWorkers executes jobs from the channel using n workers and stops on the first error or context cancel.
func runWithWorkers(ctx context.Context, n int, jobs <-chan Work) error {
	g, ctx := errgroup.WithContext(ctx)
	for i := 0; i < n; i++ {
		g.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case w, ok := <-jobs:
					if !ok {
						return nil
					}
					if err := w(); err != nil {
						return err
					}
				}
			}
		})
	}
	return g.Wait()
}
