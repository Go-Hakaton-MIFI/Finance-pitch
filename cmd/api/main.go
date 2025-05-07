package main

import (
	"crypto/rand"
	"crypto/rsa"
	"finance-backend/internal/app"
	"finance-backend/internal/delivery/http/handlers"
	approuters "finance-backend/internal/delivery/http/routers"
	userrepo "finance-backend/internal/repository/user"
	userusecase "finance-backend/internal/usecase/user"
	"log"
	"net/http"
	"time"
)

func main() {
	logger := log.New(log.Writer(), "API: ", log.LstdFlags)

	// Инициализация зависимостей
	deps, err := app.InitDependencies()
	if err != nil {
		logger.Fatalf("Failed to initialize dependencies: %v", err)
	}
	defer deps.CloseDependencies()

	// Генерация ключа для JWT
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		logger.Fatalf("Failed to generate private key: %v", err)
	}

	// Инициализация репозиториев
	userRepo := userrepo.NewUserRepository(deps.DB, deps.Logger)

	// Инициализация use cases
	transactionService := deps.TransactionService
	userUseCase := userusecase.NewUserUseCase(userRepo, privateKey, 24*time.Hour)

	// Инициализация обработчиков
	userHandler := handlers.NewUserHandler(logger, userUseCase)
	analyticsHandler := handlers.NewAnalyticsHandler(deps.DB, deps.Logger)

	// Настройка маршрутизации
	router := approuters.NewMuxRouter(userHandler, analyticsHandler, transactionService)

	// Запуск сервера
	logger.Println("Server starting on :8089")
	if err := http.ListenAndServe(":8089", router); err != nil {
		logger.Fatalf("Server failed to start: %v", err)
	}
}
