package transaction

import (
	"context"
	"finance-backend/internal/domain/transaction"
	"finance-backend/pkg/logger"
	"strconv"

	"github.com/jmoiron/sqlx"
)

type TransactionRepository struct {
	db     *sqlx.DB
	logger *logger.Logger
}

func NewTransactionRepository(db *sqlx.DB, logger *logger.Logger) *TransactionRepository {
	return &TransactionRepository{
		db:     db,
		logger: logger,
	}
}

func (r *TransactionRepository) GetTransactions(ctx context.Context, filter *transaction.TransactionFilter) ([]transaction.Transaction, error) {
	query := `
		SELECT 
			transactions.id,
			transactions.user_type,
			transactions.date_time,
			transactions.trans_type,
			transactions.amount,
			transactions.category_id,
			transactions.status_id,
			transactions.sender_bank,
			transactions.receiver_inn,
			transactions.receiver_phone,
			transactions.comment,
			c.name as category_name,
			c.type as category_type,
			s.name as status_name,
			s.description as status_description
		FROM transactions
		LEFT JOIN categories c ON transactions.category_id = c.id
		LEFT JOIN transaction_statuses s ON transactions.status_id = s.id
		WHERE 1=1
	`

	args := []interface{}{}
	if filter != nil {
		if filter.UserType != "" {
			query += " AND transactions.user_type = $" + strconv.Itoa(len(args)+1)
			args = append(args, filter.UserType)
		}
		if filter.TransType != "" {
			query += " AND transactions.trans_type = $" + strconv.Itoa(len(args)+1)
			args = append(args, filter.TransType)
		}
		if filter.SenderBank != "" {
			query += " AND transactions.sender_bank = $" + strconv.Itoa(len(args)+1)
			args = append(args, filter.SenderBank)
		}
		if filter.ReceiverINN != "" {
			query += " AND transactions.receiver_inn = $" + strconv.Itoa(len(args)+1)
			args = append(args, filter.ReceiverINN)
		}
		if filter.ReceiverPhone != "" {
			query += " AND transactions.receiver_phone = $" + strconv.Itoa(len(args)+1)
			args = append(args, filter.ReceiverPhone)
		}
		if filter.CategoryID != 0 {
			query += " AND transactions.category_id = $" + strconv.Itoa(len(args)+1)
			args = append(args, filter.CategoryID)
		}
		if filter.StatusID != 0 {
			query += " AND transactions.status_id = $" + strconv.Itoa(len(args)+1)
			args = append(args, filter.StatusID)
		}
		if !filter.DateFrom.IsZero() {
			query += " AND transactions.date_time >= $" + strconv.Itoa(len(args)+1)
			args = append(args, filter.DateFrom)
		}
		if !filter.DateTo.IsZero() {
			query += " AND transactions.date_time <= $" + strconv.Itoa(len(args)+1)
			args = append(args, filter.DateTo)
		}
	}

	query += " ORDER BY transactions.date_time DESC"

	var transactions []transaction.Transaction
	if err := r.db.SelectContext(ctx, &transactions, query, args...); err != nil {
		r.logger.Error(ctx, "error getting transactions", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	return transactions, nil
}

func (r *TransactionRepository) GetPreparedTransactions(ctx context.Context) ([]transaction.PreparedTransaction, error) {
	query := `
		SELECT 
			t.id,
			t.user_type,
			t.date_time,
			t.trans_type,
			t.amount,
			t.category_id,
			t.status_id,
			t.sender_bank,
			t.receiver_inn,
			t.receiver_phone,
			t.comment,
			c.name as category_name,
			c.type as category_type,
			s.name as status_name,
			s.description as status_description
		FROM prepared_transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		LEFT JOIN trans_statuses s ON t.status_id = s.id
		ORDER BY t.date_time DESC
	`

	var transactions []transaction.PreparedTransaction
	if err := r.db.SelectContext(ctx, &transactions, query); err != nil {
		r.logger.Error(ctx, "error getting prepared transactions", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	return transactions, nil
}

func (r *TransactionRepository) GetCategories(ctx context.Context) ([]transaction.Category, error) {
	query := `
		SELECT 
			id,
			name,
			COALESCE(type, '') as type
		FROM categories
		ORDER BY name
	`

	var categories []transaction.Category
	if err := r.db.SelectContext(ctx, &categories, query); err != nil {
		r.logger.Error(ctx, "error getting categories", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	return categories, nil
}

func (r *TransactionRepository) GetTransactionStatuses(ctx context.Context) ([]transaction.TransactionStatus, error) {
	query := `
		SELECT 
			id,
			name,
			description
		FROM transaction_statuses
		ORDER BY id
	`

	var statuses []transaction.TransactionStatus
	if err := r.db.SelectContext(ctx, &statuses, query); err != nil {
		r.logger.Error(ctx, "error getting transaction statuses", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	return statuses, nil
}

func (r *TransactionRepository) DeleteTransaction(ctx context.Context, id int) error {
	query := "DELETE FROM transactions WHERE id = $1"
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error(ctx, "error deleting transaction", map[string]interface{}{"error": err.Error(), "id": id})
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		r.logger.Error(ctx, "error getting rows affected", map[string]interface{}{"error": err.Error(), "id": id})
		return err
	}

	if rows == 0 {
		r.logger.Error(ctx, "transaction not found", map[string]interface{}{"id": id})
		return transaction.ErrTransactionNotFound
	}

	return nil
}

func (r *TransactionRepository) CreateTransaction(ctx context.Context, t *transaction.Transaction) error {
	query := `
		INSERT INTO transactions (
			user_type,
			date_time,
			trans_type,
			amount,
			category_id,
			status_id,
			sender_bank,
			receiver_inn,
			receiver_phone,
			comment
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		) RETURNING id
	`

	err := r.db.QueryRowContext(ctx, query,
		t.UserType,
		t.DateTime,
		t.TransType,
		t.Amount,
		t.CategoryID,
		t.StatusID,
		t.SenderBank,
		t.ReceiverINN,
		t.ReceiverPhone,
		t.Comment,
	).Scan(&t.ID)

	if err != nil {
		r.logger.Error(ctx, "error creating transaction", map[string]interface{}{"error": err.Error()})
		return err
	}

	return nil
}

func (r *TransactionRepository) CreatePreparedTransaction(ctx context.Context, t *transaction.PreparedTransaction) error {
	query := `
		INSERT INTO prepared_transactions (
			user_type,
			date_time,
			trans_type,
			amount,
			category_id,
			status_id,
			sender_bank,
			receiver_inn,
			receiver_phone,
			comment
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		) RETURNING id
	`

	err := r.db.QueryRowContext(ctx, query,
		t.UserType,
		t.DateTime,
		t.TransType,
		t.Amount,
		t.CategoryID,
		t.StatusID,
		t.SenderBank,
		t.ReceiverINN,
		t.ReceiverPhone,
		t.Comment,
	).Scan(&t.ID)

	if err != nil {
		r.logger.Error(ctx, "error creating prepared transaction", map[string]interface{}{"error": err.Error()})
		return err
	}

	return nil
}

// Проверка соответствия интерфейсу
var _ transaction.Repository = (*TransactionRepository)(nil)
