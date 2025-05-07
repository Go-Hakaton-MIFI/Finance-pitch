package repository

import (
	"context"
	"database/sql"
	"time"

	"finance-backend/internal/domain"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, transaction *domain.Transaction) error {
	query := `
		INSERT INTO transactions (amount, type, category_id, status_id, description, user_id, date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	now := time.Now()
	transaction.CreatedAt = now
	transaction.UpdatedAt = now

	return r.db.QueryRowContext(
		ctx,
		query,
		transaction.Amount,
		transaction.Type,
		transaction.CategoryID,
		transaction.StatusID,
		transaction.Description,
		transaction.UserID,
		transaction.Date,
		transaction.CreatedAt,
		transaction.UpdatedAt,
	).Scan(&transaction.ID)
}

func (r *TransactionRepository) GetByID(ctx context.Context, id int64) (*domain.Transaction, error) {
	query := `
		SELECT id, amount, type, category_id, status_id, description, user_id, date, created_at, updated_at
		FROM transactions
		WHERE id = $1`

	transaction := &domain.Transaction{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&transaction.ID,
		&transaction.Amount,
		&transaction.Type,
		&transaction.CategoryID,
		&transaction.StatusID,
		&transaction.Description,
		&transaction.UserID,
		&transaction.Date,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (r *TransactionRepository) ListByUserID(ctx context.Context, userID int64, limit, offset int) ([]*domain.Transaction, error) {
	query := `
		SELECT id, amount, type, category_id, status_id, description, user_id, date, created_at, updated_at
		FROM transactions
		WHERE user_id = $1
		ORDER BY date DESC, created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*domain.Transaction
	for rows.Next() {
		transaction := &domain.Transaction{}
		err := rows.Scan(
			&transaction.ID,
			&transaction.Amount,
			&transaction.Type,
			&transaction.CategoryID,
			&transaction.StatusID,
			&transaction.Description,
			&transaction.UserID,
			&transaction.Date,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *TransactionRepository) Update(ctx context.Context, transaction *domain.Transaction) error {
	query := `
		UPDATE transactions
		SET amount = $1, type = $2, category_id = $3, status_id = $4, description = $5, date = $6, updated_at = $7
		WHERE id = $8 AND user_id = $9`

	transaction.UpdatedAt = time.Now()
	result, err := r.db.ExecContext(
		ctx,
		query,
		transaction.Amount,
		transaction.Type,
		transaction.CategoryID,
		transaction.StatusID,
		transaction.Description,
		transaction.Date,
		transaction.UpdatedAt,
		transaction.ID,
		transaction.UserID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *TransactionRepository) Delete(ctx context.Context, id, userID int64) error {
	query := `DELETE FROM transactions WHERE id = $1 AND user_id = $2`

	result, err := r.db.ExecContext(ctx, query, id, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}
