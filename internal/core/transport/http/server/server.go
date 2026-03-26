package core_http_server

import (
	"context"
	"errors"
	"fmt"
	core_logger "github.com/severholod/amazing-todo/internal/core/logger"
	core_http_middleware "github.com/severholod/amazing-todo/internal/core/transport/http/middleware"
	"go.uber.org/zap"
	"net/http"
)

type HTTPServer struct {
	mux        *http.ServeMux
	config     Config
	log        *core_logger.Logger
	middleware []core_http_middleware.Middleware
}

func NewHTTPServer(
	config Config,
	log *core_logger.Logger,
	middleware ...core_http_middleware.Middleware,
) *HTTPServer {
	return &HTTPServer{
		mux:        http.NewServeMux(),
		config:     config,
		log:        log,
		middleware: middleware,
	}
}

func (s *HTTPServer) RegisterApiRoutes(routers ...*APIVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersion)

		s.mux.Handle(prefix+"/", http.StripPrefix(prefix, router))
	}
}

func (s *HTTPServer) Run(ctx context.Context) error {
	mux := core_http_middleware.ChainMiddleware(s.mux, s.middleware...)
	server := &http.Server{
		Addr:    s.config.Address,
		Handler: mux,
	}

	errChan := make(chan error, 1)
	go func() {
		defer close(errChan)

		s.log.Warn("Starting HTTP server", zap.String("address", s.config.Address))

		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		if err != nil {
			return fmt.Errorf("http server error: %w", err)
		}
	case <-ctx.Done():
		s.log.Info("Shutting down HTTP server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)

		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()
			return fmt.Errorf("http server shutdown error: %w", err)
		}
		s.log.Warn("HTTP server shutdown completed")
	}
	return nil
}
