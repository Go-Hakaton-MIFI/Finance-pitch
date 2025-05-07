package transaction

import "context"

type Repository interface {
	GetTransactions(ctx context.Context, filter *TransactionFilter) ([]Transaction, error)
	GetPreparedTransactions(ctx context.Context) ([]PreparedTransaction, error)
	GetCategories(ctx context.Context) ([]Category, error)
	GetTransactionStatuses(ctx context.Context) ([]TransactionStatus, error)
	DeleteTransaction(ctx context.Context, id int) error
	CreateTransaction(ctx context.Context, transaction *Transaction) error
	CreatePreparedTransaction(ctx context.Context, transaction *PreparedTransaction) error
}
