package main

import (
	"context"
	"finance-backend/internal/app"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	deps, err := app.InitDependencies()
	if err != nil {
		deps.Logger.Fatal(context.TODO(), "failed to initialize dependencies", map[string]interface{}{"error": err.Error()})
	}
	defer deps.CloseDependencies()

	server := app.StartHTTPServer(deps)
	defer server.Shutdown(context.Background())

	waitForShutdown(deps)
}

func waitForShutdown(deps *app.AppDependencies) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	deps.Logger.Info(context.TODO(), "Shutting down gracefully...", map[string]interface{}{})
}
