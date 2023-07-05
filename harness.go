package harness

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

func New(options ...Option) Handler {
	o := defaultOptions()
	for _, opt := range options {
		opt(o)
	}

	return &handler{
		options: o,
	}
}

type Handler interface {
	Start(ctx context.Context)
}

type handler struct {
	options *options
}

func (h *handler) Start(ctx context.Context) {
	exitTypeCh := make(chan ExitType, 1)
	ctx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()

	// Use a wait group to track the completion of all runners
	var wg sync.WaitGroup
	wg.Add(len(h.options.runners))

	// Start each runner in a separate goroutine
	go func() {
		for _, runner := range h.options.runners {
			go func(r Runner, fnCancel context.CancelFunc) {

				defer wg.Done()
				defer func(r Runner, fnCancel context.CancelFunc) {
					if rec := recover(); rec != nil {
						err := fmt.Errorf("recovered %+v", rec)
						h.options.onError(err)
						r.OnError(err)
						fnCancel()
					}
				}(r, fnCancel)

				err := r.Run(ctx)

				if err != nil {
					h.options.onError(err)
					r.OnError(err)
					fnCancel()
				}
			}(runner, cancelFunc)
		}
		wg.Wait()
		time.Sleep(100 * time.Microsecond)
		exitTypeCh <- ExitTypeNormal
	}()

	go func() {
		// Wait for termination signal or cancellation
		signalCh := make(chan os.Signal, 1)
		// defer close(signalCh)
		signal.Notify(signalCh, h.options.signals...)

		// Wait for all runners to complete
		select {
		case <-signalCh:
			exitTypeCh <- ExitTypeSignal
		case <-ctx.Done():
			exitTypeCh <- ExitTypeCancel
		}
	}()

	exitType := <-exitTypeCh
	h.gracefulShutdown(exitType)
}

func (h *handler) gracefulShutdown(exitType ExitType) {
	var wg sync.WaitGroup
	for _, runner := range h.options.runners {
		wg.Add(1)
		go func(r Runner) {

			defer wg.Done()
			defer func(r Runner) {
				if rec := recover(); rec != nil {
					err := fmt.Errorf("recovered %+v", rec)
					h.options.onError(err)
					r.OnError(err)
				}
			}(r)
			r.Shutdown(exitType)
		}(runner)
	}
	wg.Wait()
	h.options.onCompleted()
}
