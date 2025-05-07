package app

import (
	"context"
	handlers "finance-backend/internal/delivery/http/handlers"
	routers "finance-backend/internal/delivery/http/routers"
	"finance-backend/pkg/logger"
	"fmt"
	"net/http"
)

func StartHTTPServer(deps *AppDependencies) *http.Server {
	stdLogger := logger.NewStdLogger(deps.Logger)

	server := &http.Server{
		Addr: fmt.Sprintf("%s:%s", deps.Config.Server.Address, deps.Config.Server.Port),
		Handler: routers.NewMuxRouter(
			// handlers.NewCategoryHandler(*deps.Logger, deps.CategoryUseCase),
			// handlers.NewArticleHandler(*deps.Logger, deps.ArticleUseCase),
			handlers.NewUserHandler(stdLogger, deps.UserUseCase),
			deps.AnalyticsHandler,
			deps.TransactionService,
		),
	}

	go func() {
		deps.Logger.Info(context.Background(), "Starting HTTP server on", map[string]interface{}{"address": deps.Config.Server.Address, "port": deps.Config.Server.Port})
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			deps.Logger.Fatal(context.Background(), "Failed to start server", map[string]interface{}{"error": err.Error(), "addr": server.Addr})
		}
	}()

	return server
}
