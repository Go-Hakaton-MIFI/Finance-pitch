package handlers

import (
	"encoding/json"
	"errors"
	"finance-backend/internal/delivery/http/schemas"
	"finance-backend/internal/domain"
	"fmt"
	"log"
	"net/http"
	"strings"

	uc "finance-backend/internal/usecase/user"

	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	logger      *log.Logger
	userUseCase uc.IUserUseCase
}

func NewUserHandler(logger *log.Logger, userUseCase uc.IUserUseCase) *UserHandler {
	return &UserHandler{
		logger:      logger,
		userUseCase: userUseCase,
	}
}

func (h *UserHandler) GetSubjectTypes(w http.ResponseWriter, r *http.Request) {
	types := []map[string]string{
		{
			"subjectType": "INDIVIDUAL",
			"subjectName": "Физическое лицо",
		},
		{
			"subjectType": "LEGAL",
			"subjectName": "Юридическое лицо",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(types)
}

func (uh *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var requestEntity schemas.UserRegistrationSchema
	if err := json.NewDecoder(r.Body).Decode(&requestEntity); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"msg": "JSON decode error", "error": err.Error()})
		return
	}

	validate := validator.New()
	if err := validate.Struct(requestEntity); err != nil {
		errorMap := make(map[string]string)
		for _, verr := range err.(validator.ValidationErrors) {
			errorMap[strings.ToLower(verr.Field())] = fmt.Sprintf(
				"Field validation for '%s' failed on the '%s' tag", strings.ToLower(verr.Field()), verr.Tag(),
			)
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorMap)
		return
	}

	token, err := uh.userUseCase.RegisterUser(r.Context(), requestEntity.ToDomainEntity())

	if err != nil {
		var de *domain.DomainError
		if errors.As(err, &de) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": de.Message})
			return
		}
		// Обработка других ошибок
		uh.logger.Printf("Error registering user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token.String()})
}

func (uh *UserHandler) GetAccessToken(w http.ResponseWriter, r *http.Request) {
	var requestEntity schemas.UserLoginSchema

	if err := json.NewDecoder(r.Body).Decode(&requestEntity); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"msg": "JSON decode error", "error": err.Error()})
		return
	}

	validate := validator.New()
	if err := validate.Struct(requestEntity); err != nil {
		errorMap := make(map[string]string)
		for _, verr := range err.(validator.ValidationErrors) {
			errorMap[strings.ToLower(verr.Field())] = fmt.Sprintf(
				"Field validation for '%s' failed on the '%s' tag", strings.ToLower(verr.Field()), verr.Tag(),
			)
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorMap)
		return
	}

	token, err := uh.userUseCase.GetAccessToken(r.Context(), requestEntity.Login, requestEntity.Password)

	if err != nil {
		var de *domain.DomainError
		if errors.As(err, &de) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": de.Message})
			return
		}
		// Обработка других ошибок
		uh.logger.Printf("Error getting access token: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token.String()})
}
