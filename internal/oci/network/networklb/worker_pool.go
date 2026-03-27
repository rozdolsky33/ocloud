package networklb

import (
	"context"

	"golang.org/x/sync/errgroup"
)

const (
	defaultWorkerCount = 12
)

type Work func() error

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
