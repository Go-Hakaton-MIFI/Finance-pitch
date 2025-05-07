package transaction

import (
	"errors"
	"time"
)

var (
	ErrTransactionNotFound = errors.New("transaction not found")
)

type Transaction struct {
	ID                int       `db:"id"`
	UserType          string    `db:"user_type"`
	DateTime          time.Time `db:"date_time"`
	TransType         string    `db:"trans_type"`
	Amount            float64   `db:"amount"`
	CategoryID        int       `db:"category_id"`
	StatusID          int       `db:"status_id"`
	SenderBank        string    `db:"sender_bank"`
	ReceiverINN       string    `db:"receiver_inn"`
	ReceiverPhone     string    `db:"receiver_phone"`
	Comment           string    `db:"comment"`
	CategoryName      string    `db:"category_name"`
	CategoryType      string    `db:"category_type"`
	StatusName        string    `db:"status_name"`
	StatusDescription string    `db:"status_description"`
}

type PreparedTransaction struct {
	ID                int       `db:"id"`
	UserType          string    `db:"user_type"`
	DateTime          time.Time `db:"date_time"`
	TransType         string    `db:"trans_type"`
	Amount            float64   `db:"amount"`
	CategoryID        int       `db:"category_id"`
	StatusID          int       `db:"status_id"`
	SenderBank        string    `db:"sender_bank"`
	ReceiverINN       string    `db:"receiver_inn"`
	ReceiverPhone     string    `db:"receiver_phone"`
	Comment           string    `db:"comment"`
	CategoryName      string    `db:"category_name"`
	CategoryType      string    `db:"category_type"`
	StatusName        string    `db:"status_name"`
	StatusDescription string    `db:"status_description"`
}

type Category struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
	Type string `db:"type"`
}

type TransactionStatus struct {
	ID          int    `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

type TransactionFilter struct {
	UserType      string
	TransType     string
	SenderBank    string
	ReceiverINN   string
	ReceiverPhone string
	CategoryID    int
	StatusID      int
	DateFrom      time.Time
	DateTo        time.Time
}
