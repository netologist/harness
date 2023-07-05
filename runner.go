//go:generate mockgen -source=runner.go -destination=runner_mock.go -package harness Runner

package harness

import "context"

type ExitType int

const (
	ExitTypeNormal ExitType = iota
	ExitTypeCancel
	ExitTypeSignal
)

type Runner interface {
	Name() string
	Run(ctx context.Context) error
	Shutdown(exitType ExitType)
	OnError(err error)
}
