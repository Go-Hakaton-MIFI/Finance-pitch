package delivery

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"finance-backend/internal/domain"
	"finance-backend/internal/usecase/transaction"
)

type TransactionHandler struct {
	useCase *transaction.UseCase
}

func NewTransactionHandler(useCase *transaction.UseCase) *TransactionHandler {
	return &TransactionHandler{useCase: useCase}
}

type createTransactionRequest struct {
	Amount      float64 `json:"amount"`
	Type        string  `json:"type"`
	CategoryID  int64   `json:"category_id"`
	StatusID    int64   `json:"status_id"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Получаем ID пользователя из контекста (предполагается, что middleware аутентификации уже установил его)
	userID := r.Context().Value("user_id").(int64)

	// Парсим дату
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	transaction := &domain.Transaction{
		Amount:      req.Amount,
		Type:        domain.TransactionType(req.Type),
		CategoryID:  req.CategoryID,
		StatusID:    req.StatusID,
		Description: req.Description,
		UserID:      userID,
		Date:        date,
	}

	if err := h.useCase.Create(r.Context(), transaction); err != nil {
		http.Error(w, "Failed to create transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}

func (h *TransactionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	transaction, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		if err == domain.ErrNotFound {
			http.Error(w, "Transaction not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get transaction", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(transaction)
}

func (h *TransactionHandler) ListByUserID(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)
	limit := 10
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	transactions, err := h.useCase.ListByUserID(r.Context(), userID, limit, offset)
	if err != nil {
		http.Error(w, "Failed to list transactions", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(transactions)
}

type updateTransactionRequest struct {
	Amount      float64 `json:"amount"`
	Type        string  `json:"type"`
	CategoryID  int64   `json:"category_id"`
	StatusID    int64   `json:"status_id"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

func (h *TransactionHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	var req updateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(int64)

	// Парсим дату
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	transaction := &domain.Transaction{
		ID:          id,
		Amount:      req.Amount,
		Type:        domain.TransactionType(req.Type),
		CategoryID:  req.CategoryID,
		StatusID:    req.StatusID,
		Description: req.Description,
		UserID:      userID,
		Date:        date,
	}

	if err := h.useCase.Update(r.Context(), transaction); err != nil {
		if err == domain.ErrNotFound {
			http.Error(w, "Transaction not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transaction)
}

func (h *TransactionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(int64)

	if err := h.useCase.Delete(r.Context(), id, userID); err != nil {
		if err == domain.ErrNotFound {
			http.Error(w, "Transaction not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
