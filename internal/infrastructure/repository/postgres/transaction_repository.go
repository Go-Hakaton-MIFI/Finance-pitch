package postgres

import (
	"context"
	"database/sql"
	"finance-backend/internal/domain/transaction"
	"time"
)

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) transaction.Repository {
	return &transactionRepository{
		db: db,
	}
}

func (r *transactionRepository) GetTransactions(ctx context.Context, filter *transaction.TransactionFilter) ([]transaction.Transaction, error) {
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
			s.name as status_name
		FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		LEFT JOIN transaction_statuses s ON t.status_id = s.id
		WHERE ($1 = '' OR t.user_type = $1)
		AND ($2 = '' OR t.trans_type = $2)
		AND ($3 = '' OR t.sender_bank ILIKE '%' || $3 || '%')
		AND ($4 = '' OR t.receiver_inn ILIKE '%' || $4 || '%')
		AND ($5 = '' OR t.receiver_phone ILIKE '%' || $5 || '%')
		AND ($6 = 0 OR t.category_id = $6)
		AND ($7 = 0 OR t.status_id = $7)
		AND ($8::timestamp IS NULL OR t.date_time >= $8)
		AND ($9::timestamp IS NULL OR t.date_time <= $9)
		ORDER BY t.date_time DESC
	`

	rows, err := r.db.QueryContext(ctx, query,
		filter.UserType,
		filter.TransType,
		filter.SenderBank,
		filter.ReceiverINN,
		filter.ReceiverPhone,
		filter.CategoryID,
		filter.StatusID,
		filter.DateFrom,
		filter.DateTo,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []transaction.Transaction
	for rows.Next() {
		var t transaction.Transaction
		err := rows.Scan(
			&t.ID,
			&t.UserType,
			&t.DateTime,
			&t.TransType,
			&t.Amount,
			&t.CategoryID,
			&t.StatusID,
			&t.SenderBank,
			&t.ReceiverINN,
			&t.ReceiverPhone,
			&t.Comment,
			&t.CategoryName,
			&t.StatusName,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}

func (r *transactionRepository) GetPreparedTransactions(ctx context.Context) ([]transaction.PreparedTransaction, error) {
	query := `
		SELECT id, user_type, date_time, trans_type, amount, category_id, status_id,
			   sender_bank, receiver_inn, receiver_phone, comment
		FROM prepared_transactions
		ORDER BY date_time DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []transaction.PreparedTransaction
	for rows.Next() {
		var t transaction.PreparedTransaction
		err := rows.Scan(
			&t.ID,
			&t.UserType,
			&t.DateTime,
			&t.TransType,
			&t.Amount,
			&t.CategoryID,
			&t.StatusID,
			&t.SenderBank,
			&t.ReceiverINN,
			&t.ReceiverPhone,
			&t.Comment,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}

func (r *transactionRepository) GetCategories(ctx context.Context) ([]transaction.Category, error) {
	query := `SELECT id, name, type FROM categories ORDER BY name`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []transaction.Category
	for rows.Next() {
		var c transaction.Category
		err := rows.Scan(&c.ID, &c.Name, &c.Type)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}

func (r *transactionRepository) GetTransactionStatuses(ctx context.Context) ([]transaction.TransactionStatus, error) {
	query := `SELECT id, name, description FROM transaction_statuses ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statuses []transaction.TransactionStatus
	for rows.Next() {
		var s transaction.TransactionStatus
		err := rows.Scan(&s.ID, &s.Name, &s.Description)
		if err != nil {
			return nil, err
		}
		statuses = append(statuses, s)
	}

	return statuses, nil
}

func (r *transactionRepository) DeleteTransaction(ctx context.Context, id int) error {
	query := `UPDATE transactions SET status_id = 6 WHERE id = $1` // 6 = "Платеж удален"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *transactionRepository) CreateTransaction(ctx context.Context, t *transaction.Transaction) error {
	query := `
		INSERT INTO transactions (
			user_type, date_time, trans_type, amount, category_id, status_id,
			sender_bank, receiver_inn, receiver_phone, comment
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`

	if t.DateTime.IsZero() {
		t.DateTime = time.Now()
	}
	if t.StatusID == 0 {
		t.StatusID = 1 // 1 = "Новая"
	}

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

	return err
}

func (r *transactionRepository) CreatePreparedTransaction(ctx context.Context, t *transaction.PreparedTransaction) error {
	query := `
		INSERT INTO prepared_transactions (
			user_type, date_time, trans_type, amount, category_id, status_id,
			sender_bank, receiver_inn, receiver_phone, comment
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`

	if t.DateTime.IsZero() {
		t.DateTime = time.Now()
	}
	if t.StatusID == 0 {
		t.StatusID = 1 // 1 = "Новая"
	}

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

	return err
}
