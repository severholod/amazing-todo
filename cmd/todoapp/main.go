package main

import (
	"context"
	"fmt"
	core_logger "github.com/severholod/amazing-todo/internal/core/logger"
	core_postgres_pool "github.com/severholod/amazing-todo/internal/core/repository/postgres/pool"
	core_http_middleware "github.com/severholod/amazing-todo/internal/core/transport/http/middleware"
	core_http_server "github.com/severholod/amazing-todo/internal/core/transport/http/server"
	users_postgres_repository "github.com/severholod/amazing-todo/internal/features/users/repository/postgres"
	users_service "github.com/severholod/amazing-todo/internal/features/users/service"
	users_transport_http "github.com/severholod/amazing-todo/internal/features/users/transport/http"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger, logErr := core_logger.NewLogger(core_logger.NewConfigMust())
	if logErr != nil {
		fmt.Println("Failed to initialize logger:", logErr)
		os.Exit(1)
	}
	defer logger.Close()
	logger.Debug("Init database connection pool")

	pool, err := core_postgres_pool.NewConnectionPool(ctx, core_postgres_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("Failed to initialize postgres connection pool:", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("Init feature", zap.String("feature", "users"))

	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportHttp := users_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("Init HTTP Server")
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Panic(),
		core_http_middleware.Trace(),
	)

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersTransportHttp.Routes()...)
	httpServer.RegisterApiRoutes(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP Server run error", zap.Error(err))
	}
}
