package harness

type options struct {
	runners     []Runner
	onError     func(error)
	onCompleted func()
}

func defaultOptions() *options {
	return &options{
		runners:     []Runner{},
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
