package handlers

import (
	"encoding/json"
	"finance-backend/internal/delivery/http/schemas"
	"net/http"

	"finance-backend/pkg/logger"

	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
)

type AnalyticsHandler struct {
	db       *sqlx.DB
	logger   *logger.Logger
	validate *validator.Validate
}

func NewAnalyticsHandler(db *sqlx.DB, logger *logger.Logger) *AnalyticsHandler {
	return &AnalyticsHandler{
		db:       db,
		logger:   logger,
		validate: validator.New(),
	}
}

func (h *AnalyticsHandler) GetDynamicsByPeriod(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	if period == "" {
		period = "month"
	}

	var request schemas.DynamicsByPeriodRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	if err := h.validate.Struct(request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	var interval string
	switch period {
	case "week":
		interval = "1 day"
	case "month":
		interval = "1 day"
	case "quarter":
		interval = "1 week"
	case "year":
		interval = "1 month"
	default:
		interval = "1 day"
	}

	query := `
		WITH date_series AS (
			SELECT generate_series(
				$1::timestamp with time zone,
				$2::timestamp with time zone,
				$3::interval
			)::timestamp with time zone as date
		)
		SELECT 
			to_char(ds.date, 'YYYY-MM-DD') as date,
			COUNT(t.id) as count,
			COALESCE(SUM(
				CASE 
					WHEN t.trans_type = 'credit' THEN t.amount
					WHEN t.trans_type = 'debit' THEN -t.amount
					ELSE 0
				END
			), 0) as value
		FROM date_series ds
		LEFT JOIN transactions t ON DATE(t.date_time) = DATE(ds.date)
		GROUP BY ds.date
		ORDER BY ds.date
	`

	h.logger.Info(r.Context(), "Executing dynamics query", map[string]interface{}{
		"query":  query,
		"params": []interface{}{request.Date.From, request.Date.To, interval},
	})

	rows, err := h.db.QueryContext(r.Context(), query, request.Date.From, request.Date.To, interval)
	if err != nil {
		h.logger.Error(r.Context(), "error getting dynamics", map[string]interface{}{"error": err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}
	defer rows.Close()

	var response schemas.DynamicsByPeriodResponse
	for rows.Next() {
		var item struct {
			Date  string  `db:"date"`
			Count int     `db:"count"`
			Value float64 `db:"value"`
		}
		if err := rows.Scan(&item.Date, &item.Count, &item.Value); err != nil {
			h.logger.Error(r.Context(), "error scanning row", map[string]interface{}{"error": err.Error()})
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
			return
		}
		response.Data = append(response.Data, struct {
			Date  string  `json:"date"`
			Value float64 `json:"value"`
		}{
			Date:  item.Date,
			Value: item.Value,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error(r.Context(), "error encoding response", map[string]interface{}{"error": err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}
}

func (h *AnalyticsHandler) GetCategoriesSummary(w http.ResponseWriter, r *http.Request) {
	transType := r.URL.Query().Get("trans_type")
	if transType == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "trans_type parameter is required"})
		return
	}

	var request schemas.CategoriesSummaryRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	if err := h.validate.Struct(request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	query := `
		SELECT 
			COALESCE(c.name, 'Без категории') as category,
			COALESCE(SUM(t.amount), 0) as value
		FROM categories c
		LEFT JOIN transactions t ON c.id = t.category_id 
			AND t.trans_type = $1
			AND t.date_time >= $2::timestamp with time zone
			AND t.date_time <= $3::timestamp with time zone
		WHERE c.type = $1 OR c.type IS NULL
		GROUP BY c.name
		ORDER BY value DESC
	`

	h.logger.Info(r.Context(), "Executing categories summary query", map[string]interface{}{
		"query":  query,
		"params": []interface{}{transType, request.Date.From, request.Date.To},
	})

	rows, err := h.db.QueryContext(r.Context(), query, transType, request.Date.From, request.Date.To)
	if err != nil {
		h.logger.Error(r.Context(), "error getting categories summary", map[string]interface{}{"error": err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}
	defer rows.Close()

	var response schemas.CategoriesSummaryResponse
	for rows.Next() {
		var item struct {
			Category string  `db:"category"`
			Value    float64 `db:"value"`
		}
		if err := rows.Scan(&item.Category, &item.Value); err != nil {
			h.logger.Error(r.Context(), "error scanning row", map[string]interface{}{"error": err.Error()})
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
			return
		}
		response.Data = append(response.Data, struct {
			Category string  `json:"category"`
			Value    float64 `json:"value"`
		}{
			Category: item.Category,
			Value:    item.Value,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error(r.Context(), "error encoding response", map[string]interface{}{"error": err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}
}
