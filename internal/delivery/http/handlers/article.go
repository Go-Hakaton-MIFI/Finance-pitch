package handlers

import (
	"encoding/json"
	"errors"
	"finance-backend/internal/delivery/http/mappers"
	"finance-backend/internal/delivery/http/schemas"
	"finance-backend/internal/domain"
	"finance-backend/internal/usecase/article"
	"finance-backend/pkg/logger"
	"finance-backend/pkg/utils"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type ArticleHandler struct {
	ArticleUseCase article.IArticleUseCase
	log            logger.Logger
}

func NewArticleHandler(logger logger.Logger, articleUseCase article.IArticleUseCase) *ArticleHandler {
	return &ArticleHandler{
		ArticleUseCase: articleUseCase,
		log:            logger,
	}
}

func (h *ArticleHandler) GetCommonArticleById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	article, err := h.ArticleUseCase.GetArticleByID(r.Context(), id)

	if err != nil {
		var de *domain.DomainError
		if errors.As(err, &de) {
			w.WriteHeader(http.StatusNotFound)

			json.NewEncoder(w).Encode(map[string]string{"error": de.Message})
			return
		}
	}
	json.NewEncoder(w).Encode(mappers.MapArticleToArticleResponse(article))
}

func (h *ArticleHandler) SearchArticlesPaginated(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	limit, _ := strconv.Atoi(utils.GetOrDefault(queryParams, "limit", "10"))
	offset, _ := strconv.Atoi(utils.GetOrDefault(queryParams, "offset", "0"))
	search := utils.GetOrNil(queryParams, "search")
	rawCategoriesIds := queryParams["category_id"]
	var categoriesIds []int
	if rawCategoriesIds != nil {
		categoriesIds = make([]int, 0, len(rawCategoriesIds))
		for _, i := range rawCategoriesIds {
			category_id, _ := strconv.Atoi(i)
			categoriesIds = append(categoriesIds, category_id)
		}
	}

	articles, _ := h.ArticleUseCase.SearchArticlesPaginated(r.Context(), limit, offset, search, categoriesIds)

	json.NewEncoder(w).Encode(mappers.MapPaginatedArticlesToPaginatedArticlesResponse(articles))
}

func (h *ArticleHandler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	if !utils.CheckIsAdmin(r.Context()) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"msg": "admin role required"})
		return
	}

	var requestEntity schemas.CreateArticleRequest
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

	article, err := h.ArticleUseCase.CreateArticle(
		r.Context(),
		requestEntity.Header,
		requestEntity.SubHeader,
		requestEntity.Description,
	)
	if err != nil {
		var de *domain.DomainError
		if errors.As(err, &de) {
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(map[string]string{"error": de.Message})
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(mappers.MapArticleToArticleResponse(article))
}

func (h *ArticleHandler) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	if !utils.CheckIsAdmin(r.Context()) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"msg": "admin role required"})
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	var requestEntity schemas.UpdateArticleRequest

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

	article, err := h.ArticleUseCase.UpdateArticle(
		r.Context(),
		id,
		requestEntity.Header,
		requestEntity.SubHeader,
		requestEntity.Description,
	)
	if err != nil {
		var de *domain.DomainError
		if errors.As(err, &de) {
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(map[string]string{"error": de.Message})
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mappers.MapArticleToArticleResponse(article))

}

func (h *ArticleHandler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	if !utils.CheckIsAdmin(r.Context()) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"msg": "admin role required"})
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	err := h.ArticleUseCase.DeleteArticle(r.Context(), id)
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

func (h *ArticleHandler) LinkCategories(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	var requestEntity schemas.LinkCategoryRequest

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

	article, err := h.ArticleUseCase.LinkCategories(
		r.Context(),
		id,
		requestEntity.CategoriesIDs,
	)
	if err != nil {
		var de *domain.DomainError
		if errors.As(err, &de) {
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(map[string]string{"error": de.Message})
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mappers.MapArticleToArticleResponse(article))

}

func (h *ArticleHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	if !utils.CheckIsAdmin(r.Context()) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"msg": "admin role required"})
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	contentType := fileHeader.Header.Get("Content-Type")
	size := fileHeader.Size

	article, err := h.ArticleUseCase.LinkImage(r.Context(), id, file, size, contentType)

	if err != nil {
		var de *domain.DomainError
		if errors.As(err, &de) {
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(map[string]string{"error": de.Message})
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mappers.MapArticleToArticleResponse(article))
}
