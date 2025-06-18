package runner

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Runner struct {
	logger    Logger
	ctx       context.Context //nolint:containedctx
	breakOnce sync.Once
	breakChan chan string // module name
}

func NewRunner(ctx context.Context, logger Logger) (*Runner, error) {
	if ctx == nil {
		panic("runner: nil context")
	}
	if logger == nil {
		panic("runner: nil logger")
	}

	return &Runner{
		logger:    logger,
		ctx:       ctx,
		breakChan: make(chan string, 1),
	}, nil
}

func (r *Runner) RunModule(module Module) func() {
	if module == nil {
		panic("runner: nil module")
	}

	moduleName := module.Name()

	ctx, cancel := context.WithCancel(r.ctx)
	done := make(chan struct{})

	r.logger.Debug(moduleName + " running")

	go func() {
		defer r.breakOnce.Do(func() {
			r.breakChan <- moduleName
		})
		defer close(done)

		module.Run(ctx)
	}()

	var isTerminated bool

	return func() {
		if isTerminated {
			return
		}
		isTerminated = true

		cancel()
		<-done

		r.logger.Debug(moduleName + " terminated")
	}
}

func (r *Runner) Listen() error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigChan:
		r.logger.Debug("signal received: " + sig.String())

		return nil
	case moduleName := <-r.breakChan:
		return fmt.Errorf("%w: %s", errors.New("app: unexpected break"), moduleName)
	}
}
