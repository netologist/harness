package harness_test

import (
	"context"
	"errors"
	"syscall"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/netologist/harness"
	"github.com/stretchr/testify/assert"
)

func TestNewHarness(t *testing.T) {
	tests := []struct {
		name             string
		runFunc          func(context.Context) error
		shutdownFunc     func(harness.ExitType)
		signalFunc       func()
		expectedErr      error
		expectedExitType harness.ExitType
	}{
		{
			name: "Successful Run and Shutdown",
			runFunc: func(ctx context.Context) error {
				// Simulate some work
				return nil
			},
			signalFunc: func() {},
			shutdownFunc: func(exitType harness.ExitType) {
				assert.Equal(t, harness.ExitTypeNormal, exitType)
			},
			expectedErr: nil,
		},
		{
			name: "Handle When Run Returns Error",
			runFunc: func(ctx context.Context) error {
				return errors.New("run error")
			},
			signalFunc: func() {},
			shutdownFunc: func(exitType harness.ExitType) {
				assert.Equal(t, harness.ExitTypeCancel, exitType)
			},
			expectedErr: errors.New("run error"),
		},
		{
			name: "Handle Run When Throw Panic",
			runFunc: func(ctx context.Context) error {
				panic("run panic error")
			},
			signalFunc: func() {},
			shutdownFunc: func(exitType harness.ExitType) {
				assert.Equal(t, harness.ExitTypeCancel, exitType)
			},
			expectedErr: errors.New("recovered run panic error"),
		},
		{
			name: "Handle Shutdown When Throw Panic",
			runFunc: func(ctx context.Context) error {
				// Simulate successful execution
				return nil
			},
			signalFunc: func() {},
			shutdownFunc: func(exitType harness.ExitType) {
				assert.Equal(t, harness.ExitTypeNormal, exitType)
				panic("shutdown panic")
			},
			expectedErr: errors.New("recovered shutdown panic"),
		},
		{
			name: "Handle Terminate When Send Signal",
			runFunc: func(ctx context.Context) error {
				time.Sleep(2 * time.Second)
				return nil
			},
			signalFunc: func() {
				go func() {
					time.Sleep(100 * time.Microsecond)
					assert.Nil(t, syscall.Kill(syscall.Getpid(), syscall.SIGTERM))
				}()
			},
			shutdownFunc: func(exitType harness.ExitType) {
				assert.Equal(t, harness.ExitTypeSignal, exitType)
			},
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var actualErr error
			var isGracefullyShutdown bool

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRunner := harness.NewMockRunner(ctrl)
			mockRunner.EXPECT().Run(gomock.Any()).DoAndReturn(test.runFunc).Times(1)
			mockRunner.EXPECT().Shutdown(gomock.Any()).DoAndReturn(test.shutdownFunc).Times(1)
			if test.expectedErr != nil {
				mockRunner.EXPECT().OnError(gomock.Any()).Times(1)
			} else {
				mockRunner.EXPECT().OnError(gomock.Any()).Times(0)
			}

			test.signalFunc()

			harness.New(
				harness.Register(mockRunner),
				harness.OnError(func(err error) {
					actualErr = err
				}),
				harness.OnCompleted(func() {
					isGracefullyShutdown = true
				}),
			).Start(context.Background())

			// Check the results
			if test.expectedErr != nil {
				assert.Equal(t, test.expectedErr, actualErr)
			} else {
				assert.Empty(t, actualErr)
			}
			assert.True(t, isGracefullyShutdown)
		})
	}
}
