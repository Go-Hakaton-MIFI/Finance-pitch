package schemas

import "time"

type Transaction struct {
	ID            int       `json:"id"`
	UserType      string    `json:"user_type"` // ФЛ или ЮЛ
	DateTime      time.Time `json:"date_time"` // Дата и время операции
	TransType     string    `json:"trans_type"`
	Amount        float64   `json:"amount"`      // Сумма (точность до 5 знаков)
	CategoryID    int       `json:"category_id"` // ID категории
	StatusID      int       `json:"status_id"`
	SenderBank    string    `json:"sender_bank"`    // Банк отправителя
	ReceiverINN   string    `json:"receiver_inn"`   // ИНН получателя
	ReceiverPhone string    `json:"receiver_phone"` // Телефон получателя
	Comment       string    `json:"comment"`        // Комментарий к операции
	CategoryName  string    `json:"category_name"`
	StatusName    string    `json:"status_name"`
}

type TransactionFilter struct {
	UserType      string    `json:"user_type"`
	TransType     string    `json:"trans_type"`
	SenderBank    string    `json:"sender_bank"`
	ReceiverINN   string    `json:"receiver_inn"`
	ReceiverPhone string    `json:"receiver_phone"`
	DateFrom      time.Time `json:"date_from"`
	DateTo        time.Time `json:"date_to"`
	CategoryID    int       `json:"category_id"`
	StatusID      int       `json:"status_id"`
}

type PreparedTransaction struct {
	ID            int       `json:"id"`
	UserType      string    `json:"user_type"`
	DateTime      time.Time `json:"date_time"`
	TransType     string    `json:"trans_type"`
	Amount        float64   `json:"amount"`
	CategoryID    int       `json:"category_id"`
	StatusID      int       `json:"status_id"`
	SenderBank    string    `json:"sender_bank"`
	ReceiverINN   string    `json:"receiver_inn"`
	ReceiverPhone string    `json:"receiver_phone"`
	Comment       string    `json:"comment"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"` // credit или debit
}

type TransactionStatus struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Структуры для аналитики
type DynamicsResponse struct {
	PeriodStart string  `json:"period_start"`
	Count       int     `json:"count"`
	Amount      float64 `json:"amount"`
}

type CategorySummaryResponse struct {
	Category string  `json:"category"`
	Count    int     `json:"count"`
	Amount   float64 `json:"amount"`
}

type BankSummaryResponse struct {
	Bank   string  `json:"bank"`
	Count  int     `json:"count"`
	Amount float64 `json:"amount"`
}

type StatusSummaryResponse struct {
	Status string  `json:"status"`
	Count  int     `json:"count"`
	Amount float64 `json:"amount"`
}
