package handlers

import (
	"encoding/json"
	"finance-backend/internal/delivery/http/schemas"
	"finance-backend/internal/domain/transaction"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type TransactionHandler struct {
	validate     *validator.Validate
	transService transaction.Service
}

func NewTransactionHandler(transService transaction.Service) *TransactionHandler {
	return &TransactionHandler{
		validate:     validator.New(),
		transService: transService,
	}
}

func (h *TransactionHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	var filter schemas.TransactionFilter
	if r.Method == "POST" {
		if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
	}

	// Получаем транзакции из базы данных
	transactions, err := h.transService.GetTransactions(r.Context(), filter)
	if err != nil {
		log.Printf("Error getting transactions: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

func (h *TransactionHandler) GetPreparedTransactions(w http.ResponseWriter, r *http.Request) {
	transactions, err := h.transService.GetPreparedTransactions(r.Context())
	if err != nil {
		log.Printf("Error getting prepared transactions: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

func (h *TransactionHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.transService.GetCategories(r.Context())
	if err != nil {
		log.Printf("Error getting categories: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func (h *TransactionHandler) GetTransactionStatuses(w http.ResponseWriter, r *http.Request) {
	statuses, err := h.transService.GetTransactionStatuses(r.Context())
	if err != nil {
		log.Printf("Error getting transaction statuses: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statuses)
}

func (h *TransactionHandler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	err = h.transService.DeleteTransaction(r.Context(), id)
	if err != nil {
		log.Printf("Error deleting transaction: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Transaction %d deleted successfully", id),
	})
}

func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction schemas.Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Валидация обязательных полей
	if err := h.validate.Struct(transaction); err != nil {
		http.Error(w, "Validation failed", http.StatusBadRequest)
		return
	}

	createdTransaction, err := h.transService.CreateTransaction(r.Context(), transaction)
	if err != nil {
		log.Printf("Error creating transaction: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdTransaction)
}

func (h *TransactionHandler) CreatePreparedTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction schemas.PreparedTransaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Валидация обязательных полей
	if err := h.validate.Struct(transaction); err != nil {
		http.Error(w, "Validation failed", http.StatusBadRequest)
		return
	}

	createdTransaction, err := h.transService.CreatePreparedTransaction(r.Context(), transaction)
	if err != nil {
		log.Printf("Error creating prepared transaction: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdTransaction)
}
