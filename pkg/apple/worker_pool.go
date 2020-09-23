package apple

import (
	"context"
	"golang.org/x/sync/semaphore"
	"runtime"
)

var (
	maxWorkers = runtime.GOMAXPROCS(0)
	sem        = semaphore.NewWeighted(int64(maxWorkers))
	out        = make([]int, 32)
)

func (v *Verifier) Start() error {
	defer v.logger.Sync()
	sugar := v.logger.Sugar()

	ctx := context.Background()

	subsCh := make(chan Subscription)

	// Retrieve subscriptions in a separate goroutine
	// so that the channel won't be blocked here.
	go func() {
		err := v.LoadSubs(subsCh)
		if err != nil {
			sugar.Error(err)
		}
	}()

	// Compute the output using up to maxWorkers goroutines at a time.
	for i := range out {
		// When maxWorkers goroutines are in flight, Acquire blocks until one of the
		// workers finishes.
		if err := sem.Acquire(ctx, 1); err != nil {
			sugar.Errorf("Failed to acquire semaphore: %v", err)
			break
		}

		go func(i int) {
			sugar.Infof("Start worker %d", i)
			defer sem.Release(1)

			for subs := range subsCh {
				//err := v.Verify(subs)
				//if err != nil {
				//	sugar.Error(err)
				//}

				sugar.Infof("Will verify %v", subs)
			}
		}(i)
	}

	// Acquire all of the tokens to wait for any remaining workers to finish.
	//
	// If you are already waiting for the workers by some other means (such as an
	// errgroup.Group), you can omit this final Acquire call.
	if err := sem.Acquire(ctx, int64(maxWorkers)); err != nil {
		sugar.Infof("Failed to acquire semaphore: %v", err)
		return nil
	}

	return nil
}
