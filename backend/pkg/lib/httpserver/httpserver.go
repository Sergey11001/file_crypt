package httpserver

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"univer/pkg/lib/log"
)

type Config struct {
	Addr            string        `default:"localhost:8000"`
	ShutdownTimeout time.Duration `default:"60s"`
}

type HTTPServer struct {
	logger          Logger
	addr            string
	handler         http.Handler
	shutdownTimeout time.Duration
}

func (s *HTTPServer) Name() string {
	return "http server"
}

func New(config Config, logger Logger, root http.Handler) (*HTTPServer, error) {
	if logger == nil {
		panic("http server: nil logger")
	}
	if root == nil {
		panic("http server: nil root")
	}

	return &HTTPServer{
		logger:          logger,
		addr:            config.Addr,
		handler:         root,
		shutdownTimeout: config.ShutdownTimeout,
	}, nil
}

func (s *HTTPServer) Run(ctx context.Context) {
	if ctx == nil {
		panic("http server: nil context")
	}

	driver := &http.Server{
		Addr:              s.addr,
		Handler:           s.handler,
		ReadHeaderTimeout: 0,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	stop := make(chan struct{})
	go func() {
		defer close(stop)

		err := driver.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("http server: listening failed", log.ErrorAttr(err))
		}
	}()

	select {
	case <-ctx.Done():
		s.logger.Debug("http server: context done", log.ErrorAttr(ctx.Err()))

		shutdownCtx := context.Background()
		if s.shutdownTimeout > 0 {
			var cancel context.CancelFunc
			shutdownCtx, cancel = context.WithTimeout(shutdownCtx, s.shutdownTimeout)
			defer cancel()
		}

		err := driver.Shutdown(shutdownCtx)
		if err != nil {
			s.logger.Error("http server: shutdown failed", log.ErrorAttr(err))
		}

		<-stop
	case <-stop:
	}
}
