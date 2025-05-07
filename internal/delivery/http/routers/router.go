package http

import (
	"finance-backend/internal/delivery/http/handlers"
	"finance-backend/internal/domain/transaction"
	"finance-backend/pkg/middleware"

	"github.com/gorilla/mux"
)

func NewMuxRouter(
	// categoryHandler *handlers.CategoryHandler,
	// articleHandler *handlers.ArticleHandler,
	userHandler *handlers.UserHandler,
	analyticsHandler *handlers.AnalyticsHandler,
	transactionService transaction.Service,
) *mux.Router {
	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()

	router.Use(middleware.RequestIdMiddleware)
	router.Use(middleware.LoggingProcessTimeMiddleware)
	router.Use(middleware.RecoverMiddleware)

	authRouter := router.NewRoute().Subrouter()
	authRouter.Use(middleware.JWTParserMiddleware)

	router.HandleFunc("/subject_types", userHandler.GetSubjectTypes).Methods("GET")
	router.HandleFunc("/registration", userHandler.RegisterUser).Methods("POST")
	router.HandleFunc("/login", userHandler.GetAccessToken).Methods("POST")

	// authRouter.HandleFunc("/admin/categories/{id}", categoryHandler.GetAdminCategoryById).Methods("GET")
	// authRouter.HandleFunc("/categories/{id}", categoryHandler.GetCommonCategoryById).Methods("GET")
	// authRouter.HandleFunc("/categories", categoryHandler.SearchCategoriesFlat).Methods("GET")
	// authRouter.HandleFunc("/admin/categories", categoryHandler.SearchCategoriesPaginated).Methods("GET")
	// authRouter.HandleFunc("/categories", categoryHandler.CreateCategory).Methods("POST")
	// authRouter.HandleFunc("/categories/{id}", categoryHandler.UpdateCategoryName).Methods("PUT")
	// authRouter.HandleFunc("/categories/{id}", categoryHandler.DeleteCategory).Methods("DELETE")

	// authRouter.HandleFunc("/articles/{id}", articleHandler.GetCommonArticleById).Methods("GET")
	// authRouter.HandleFunc("/articles", articleHandler.SearchArticlesPaginated).Methods("GET")
	// authRouter.HandleFunc("/articles", articleHandler.CreateArticle).Methods("POST")
	// authRouter.HandleFunc("/articles/{id}", articleHandler.UpdateArticle).Methods("PUT")
	// authRouter.HandleFunc("/articles/{id}/image", articleHandler.UploadImage).Methods("PUT")
	// authRouter.HandleFunc("/articles/{id}", articleHandler.DeleteArticle).Methods("DELETE")
	// authRouter.HandleFunc("/articles/{id}/categories", articleHandler.LinkCategories).Methods("PUT")

	router.HandleFunc("/analytics/dynamics/by-period", analyticsHandler.GetDynamicsByPeriod).Methods("POST")
	router.HandleFunc("/analytics/categories-summary", analyticsHandler.GetCategoriesSummary).Methods("POST")

	transactionHandler := handlers.NewTransactionHandler(transactionService)
	SetupRoutes(router, transactionHandler)

	return router
}

func SetupRoutes(router *mux.Router, transactionHandler *handlers.TransactionHandler) {
	// Маршруты для транзакций
	router.HandleFunc("/transactions", transactionHandler.GetTransactions).Methods("GET")
	router.HandleFunc("/transactions/filter", transactionHandler.GetTransactions).Methods("POST")
	router.HandleFunc("/transactions", transactionHandler.CreateTransaction).Methods("POST")
	router.HandleFunc("/transactions/{id}", transactionHandler.DeleteTransaction).Methods("DELETE")

	// Маршруты для подготовленных транзакций
	router.HandleFunc("/transactions/prepared", transactionHandler.GetPreparedTransactions).Methods("GET")
	router.HandleFunc("/transactions/prepared", transactionHandler.CreatePreparedTransaction).Methods("POST")

	// Маршруты для категорий и статусов
	router.HandleFunc("/categories", transactionHandler.GetCategories).Methods("GET")
	router.HandleFunc("/trans_statuses", transactionHandler.GetTransactionStatuses).Methods("GET")
}
