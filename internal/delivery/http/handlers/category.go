package handlers

import (
	"encoding/json"
	"errors"
	"finance-backend/internal/delivery/http/mappers"
	"finance-backend/internal/delivery/http/schemas"
	"finance-backend/internal/domain"
	"finance-backend/internal/usecase/category"
	"finance-backend/pkg/logger"
	"finance-backend/pkg/utils"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/gorilla/mux"
)

type CategoryHandler struct {
	CategoryUseCase category.ICategoryUseCase
	log             logger.Logger
}

func NewCategoryHandler(logger logger.Logger, categoryUseCase category.ICategoryUseCase) *CategoryHandler {
	return &CategoryHandler{
		CategoryUseCase: categoryUseCase,
		log:             logger,
	}
}

func (h *CategoryHandler) GetAdminCategoryById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	category, err := h.CategoryUseCase.GetCategoryByID(r.Context(), id)

	if err != nil {
		var de *domain.DomainError
		if errors.As(err, &de) {
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(map[string]string{"error": de.Message})
			return
		}
	}
	json.NewEncoder(w).Encode(mappers.MapCategoryToCategoryAdminResponse(category))
}

func (h *CategoryHandler) GetCommonCategoryById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	category, err := h.CategoryUseCase.GetCategoryByID(r.Context(), id)

	if err != nil {
		var de *domain.DomainError
		if errors.As(err, &de) {
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(map[string]string{"error": de.Message})
			return
		}
	}
	json.NewEncoder(w).Encode(mappers.MapCategoryToCategoryACommonResponse(category))
}

func (h *CategoryHandler) SearchCategoriesPaginated(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	limit, _ := strconv.Atoi(utils.GetOrDefault(queryParams, "limit", "10"))
	offset, _ := strconv.Atoi(utils.GetOrDefault(queryParams, "limit", "0"))
	search := utils.GetOrNil(queryParams, "search")

	categories, _ := h.CategoryUseCase.SearchCategoriesPaginated(r.Context(), limit, offset, search)

	json.NewEncoder(w).Encode(mappers.MapPaginatedCategoriesToAdminResponse(categories))
}

func (h *CategoryHandler) SearchCategoriesFlat(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	search := utils.GetOrNil(queryParams, "search")

	categories, _ := h.CategoryUseCase.SearchCategoriesFlat(r.Context(), search)

	json.NewEncoder(w).Encode(mappers.MapCategoriesToCategoriesRepsonse(categories))
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	if !utils.CheckIsAdmin(r.Context()) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"msg": "admin role required"})
		return
	}

	var requestEntity schemas.UpdateOrCreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&requestEntity); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"msg": "JSON decode error", "error": err.Error()})
		return
	}

	validate := validator.New()
	if err := validate.Struct(requestEntity); err != nil {
		errorMap := make(map[string]string)
		for _, verr := range err.(validator.ValidationErrors) {
			errorMap[strings.ToLower(verr.Field())] = fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", strings.ToLower(verr.Field()), verr.Tag())
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorMap)
		return
	}

	category, err := h.CategoryUseCase.CreateCategory(r.Context(), requestEntity.Name)
	if err != nil {
		var de *domain.DomainError
		if errors.As(err, &de) {
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(map[string]string{"error": de.Message})
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(mappers.MapCategoryToCategoryAdminResponse(category))
}

func (h *CategoryHandler) UpdateCategoryName(w http.ResponseWriter, r *http.Request) {
	if !utils.CheckIsAdmin(r.Context()) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"msg": "admin role required"})
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	var requestEntity schemas.UpdateOrCreateCategoryRequest

	if err := json.NewDecoder(r.Body).Decode(&requestEntity); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"msg": "JSON decode error", "error": err.Error()})
		return
	}

	validate := validator.New()
	if err := validate.Struct(requestEntity); err != nil {
		errorMap := make(map[string]string)
		for _, verr := range err.(validator.ValidationErrors) {
			errorMap[strings.ToLower(verr.Field())] = fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", strings.ToLower(verr.Field()), verr.Tag())
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorMap)
		return
	}

	category, err := h.CategoryUseCase.UpdateCategoryName(r.Context(), id, requestEntity.Name)
	if err != nil {
		var de *domain.DomainError
		if errors.As(err, &de) {
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(map[string]string{"error": de.Message})
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mappers.MapCategoryToCategoryAdminResponse(category))

}

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	if !utils.CheckIsAdmin(r.Context()) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"msg": "admin role required"})
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	err := h.CategoryUseCase.DeleteCategory(r.Context(), id)
	if err != nil {
		var de *domain.DomainError
		if errors.As(err, &de) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": de.Message})
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}
