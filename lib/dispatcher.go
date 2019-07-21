package lib

import (
	"context"
	"sync"
	"time"
)

type DispatchOptions struct {
	Config    *Config
	Ctx       context.Context
	Interval  time.Duration
	Namespace string
}

type Dispatcher struct {
	opts *DispatchOptions
}

func (d *Dispatcher) Poll() error {
	timer := time.NewTicker(d.opts.Interval)
	defer timer.Stop()

	for {
		select {
		case <-d.opts.Ctx.Done():
			return nil
		case <-timer.C:
			var wg sync.WaitGroup

			errors := make(chan error, len(d.opts.Config.Accounts))

			for _, account := range d.opts.Config.Accounts {
				wg.Add(1)
				go d.startWorker(&wg, account, errors)
			}

			wg.Wait()

			close(errors)

			for err := range errors {
				if err != nil {
					return err
				}
			}
		}
	}
}

func (d *Dispatcher) startWorker(wg *sync.WaitGroup, account *account, results chan error) {
	defer wg.Done()

	worker := newWorker(account)

	results <- worker.do()
}

func NewDispatcher(opts *DispatchOptions) *Dispatcher {
	return &Dispatcher{
		opts,
	}
}
