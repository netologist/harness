package harness

import (
	"os"
	"syscall"
)

type options struct {
	runners     []Runner
	signals     []os.Signal
	onError     func(error)
	onCompleted func()
}

func defaultOptions() *options {
	return &options{
		runners:     []Runner{},
		signals:     []os.Signal{os.Interrupt, syscall.SIGINT, syscall.SIGTERM},
		onError:     func(err error) {},
		onCompleted: func() {},
	}
}

type Option func(o *options)

func Register(runner ...Runner) Option {
	return func(o *options) {
		o.runners = append(o.runners, runner...)
	}
}

func SetSignal(signal ...os.Signal) Option {
	return func(o *options) {
		o.signals = signal
	}
}

func OnCompleted(fnOnCompleted func()) Option {
	return func(o *options) {
		o.onCompleted = fnOnCompleted
	}
}

func OnError(fnOnError func(error)) Option {
	return func(o *options) {
		o.onError = fnOnError
	}
}
