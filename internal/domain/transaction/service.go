package transaction

import (
	"context"
	"finance-backend/internal/delivery/http/schemas"
)

type Service interface {
	GetTransactions(ctx context.Context, filter schemas.TransactionFilter) ([]schemas.Transaction, error)
	GetPreparedTransactions(ctx context.Context) ([]schemas.PreparedTransaction, error)
	GetCategories(ctx context.Context) ([]schemas.Category, error)
	GetTransactionStatuses(ctx context.Context) ([]schemas.TransactionStatus, error)
	DeleteTransaction(ctx context.Context, id int64) error
	CreateTransaction(ctx context.Context, transaction schemas.Transaction) (schemas.Transaction, error)
	CreatePreparedTransaction(ctx context.Context, transaction schemas.PreparedTransaction) (schemas.PreparedTransaction, error)
}
