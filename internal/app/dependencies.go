package app

import (
	"context"
	"crypto/rsa"
	"errors"
	"finance-backend/internal/config"
	handlers "finance-backend/internal/delivery/http/handlers"
	"finance-backend/internal/domain/transaction"
	"finance-backend/internal/gateways/file_gateway"
	articleRepository "finance-backend/internal/repository/article"
	categoryRepository "finance-backend/internal/repository/category"
	transactionRepository "finance-backend/internal/repository/transaction"
	userRepository "finance-backend/internal/repository/user"
	"os"
	"strings"
	"time"

	"finance-backend/internal/usecase/article"
	"finance-backend/internal/usecase/category"
	"finance-backend/internal/usecase/user"

	"finance-backend/pkg/logger"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// AppDependencies содержит все зависимости приложения.
type AppDependencies struct {
	Config             *config.Config
	Logger             *logger.Logger
	CategoryUseCase    category.ICategoryUseCase
	ArticleUseCase     article.IArticleUseCase
	UserUseCase        user.IUserUseCase
	TransactionService transaction.Service
	AnalyticsHandler   *handlers.AnalyticsHandler
	DB                 *sqlx.DB
}

func InitDependencies() (*AppDependencies, error) {

	// 1. Загрузить конфигурацию
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	// 2. Инициализировать логгер
	log := logger.NewLogger()

	// 3. Подключение к базе данных
	db, err := sqlx.Connect(cfg.Database.Engine, cfg.Database.DSN())

	if err != nil {
		log.Error(context.TODO(), "failed to connect to database", map[string]interface{}{
			"error": err.Error(),
			"host":  cfg.Database.Host,
			"port":  cfg.Database.Port,
		})
		return nil, err
	}

	// 3.1 Подключение к S3
	sess, err := session.NewSession(&aws.Config{
		Endpoint:         aws.String(cfg.S3.Url),
		Credentials:      credentials.NewStaticCredentials(cfg.S3.S3RootUser, cfg.S3.S3RootPassword, ""),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String("us-east-1"), // может быть любое, главное заполнить
	})

	if err != nil {
		log.Fatal(context.TODO(), "Failed to create AWS session", map[string]interface{}{"error": err.Error()})
	}

	log.Info(context.TODO(), "connected_to_s3", map[string]interface{}{"endpoint": cfg.S3.Url})

	key, err := loadPrivateKeyFromEnv("JWT_PRIVATE_KEY")

	if err != nil {
		log.Fatal(context.TODO(), "Failed to get JWT private key", map[string]interface{}{"error": err.Error()})
	}
	// 4. Репозитории
	categoryRepo := categoryRepository.NewCategoryRepository(log, db)
	articleRepo := articleRepository.NewArticleRepository(log, db)
	userRepo := userRepository.NewUserRepository(db, log)
	transactionRepo := transactionRepository.NewTransactionRepository(db, log)

	// 4.1 Гейтвеи
	file_gw := file_gateway.NewS3Gateway(sess, log)

	// 5. Бизнес-логика
	categoryUseCase := category.NewCategoryUseCase(log, categoryRepo)
	articleUseCase := article.NewArticleUseCase(log, articleRepo, file_gw, cfg.ImageBucketName)
	userUseCase := user.NewUserUseCase(userRepo, key, time.Hour*24)
	transactionService := transaction.NewService(transactionRepo)

	analyticsHandler := handlers.NewAnalyticsHandler(db, log)

	return &AppDependencies{
		Config:             cfg,
		Logger:             log,
		CategoryUseCase:    categoryUseCase,
		ArticleUseCase:     articleUseCase,
		UserUseCase:        userUseCase,
		TransactionService: transactionService,
		AnalyticsHandler:   analyticsHandler,
		DB:                 db,
	}, nil
}

// CloseDependencies закрывает все ресурсы приложения.
func (d *AppDependencies) CloseDependencies() {
	if err := d.DB.Close(); err != nil {
		d.Logger.Error(context.TODO(), "failed to close database connection", map[string]interface{}{
			"error": err.Error()})
	}
	d.Logger.Info(context.TODO(), "Application dependencies closed successfully", nil)
}

func loadPrivateKeyFromEnv(envVar string) (*rsa.PrivateKey, error) {
	raw := os.Getenv(envVar)
	if raw == "" {
		return nil, errors.New("missing environment variable: " + envVar)
	}

	pemStr := strings.ReplaceAll(raw, `\n`, "\n")

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(pemStr))
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}
