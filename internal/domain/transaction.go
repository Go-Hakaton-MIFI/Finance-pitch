package domain

import "time"

type TransactionType string

const (
	Income  TransactionType = "income"
	Expense TransactionType = "expense"
)

type Transaction struct {
	ID          int64           `db:"id"`
	Amount      float64         `db:"amount"`
	Type        TransactionType `db:"type"`
	CategoryID  int64           `db:"category_id"`
	StatusID    int64           `db:"status_id"`
	Description string          `db:"description"`
	UserID      int64           `db:"user_id"`
	Date        time.Time       `db:"date"`
	CreatedAt   time.Time       `db:"created_at"`
	UpdatedAt   time.Time       `db:"updated_at"`
}
