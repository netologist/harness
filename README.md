## Harness Package
The `harness` package provides functionality to manage and handle multiple runners concurrently. It allows graceful shutdown of these runners in response to termination signals or cancellation events.

### Usage
To use the `harness` package, follow the steps below:

1. Import the `harness` package:

```go
import "github.com/netologist/harness"
```
2. Create a new handler using the New function, passing one or more Runner instances:

```go
runners := []harness.Runner{
    // Initialize your runners here
}

handler := harness.New(runners...)
```

3. Start the handler by calling the Start method, passing a context:

```go
ctx := context.Background()
handler.Start(ctx)
```

4. Graceful shutdown:

   - If a termination signal (e.g., SIGINT or SIGTERM) is received, the handler will initiate a graceful shutdown by calling the Shutdown method on each runner.
   - If cancellation is triggered on the provided context, the handler will also initiate a graceful shutdown.

### Example Use
Here's an example use case to illustrate how the `harness` package can be used:

```go
package main

import (
	"context"
	"fmt"
	"github.com/netologist/harness"
	"os"
	"os/signal"
	"syscall"
)

type TestRunner struct {
}
func (r *TestRunner) Name() string {
	return "test runner"
}
func (r *TestRunner) Run(ctx context.Context) error {
	return nil
}
func (r *TestRunner) Shutdown(exitType harness.ExitType) {
	log.Printf("NAME: '%s' - EXIT_TYPE: %d", r.Name(), exitType)
}

func (r *TestRunner) OnError(err error) {
	log.Printf("NAME: '%s' - ERROR: %+v", r.Name(), err)
}

func main() {
    ctx := context.Background()

	// Create a new TestRunner instance
	testRunner := &TestRunner{
		// Initialize your runner
		// ...
	}

	// Create the handler with the runner
	harness.New(
		harness.Register(testRunner),
		harness.OnError(func(err error) {
			log.Printf("error: %+v", err)
		}),
		harness.OnCompleted(func() {
			log.Printf("successfully completed")
		}),
		harnes.SetSignal(os.Interrupt, syscall.SIGINT, syscall.SIGTERM), // if you want customise signals
	).Start(context.Background())
}
```

In this example, we create a custom TestRunner struct that implements the Runner interface required by the `harness` package. We then create a handler with the testRunner instance and start it in a separate goroutine. We handle termination signals and cancellation requests, triggering the corresponding actions to gracefully shut down the runners. Finally, we wait for the handler to complete and perform any necessary cleanup or exit operations.

Feel free to customize the example and adapt it to your specific use case.

Please note that this is a simplified example, and you will need to implement the Runner interface methods and define your custom logic within the TestRunner struct based on your requirements.

I hope this helps! Let me know if you have any further questions.