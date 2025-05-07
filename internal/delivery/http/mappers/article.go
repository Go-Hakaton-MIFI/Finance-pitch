package mappers

import (
	"finance-backend/internal/delivery/http/schemas"
	"finance-backend/internal/domain"
	"finance-backend/pkg/utils"
)

func MapArticleToArticleResponse(article *domain.Article) schemas.ArticleResponse {
	categories := make([]schemas.CategoryResponse, 0, len(article.Categories))
	for _, category := range article.Categories {
		categories = append(categories, MapCategoryToCategoryACommonResponse(&category))
	}
	return schemas.ArticleResponse{
		ID:          article.ID,
		Header:      article.Header,
		SubHeader:   article.SubHeader,
		Description: article.Description,
		Image:       article.Image,
		Categories:  categories,
	}
}

func MapPaginatedArticlesToPaginatedArticlesResponse(
	input utils.RestfullPaginatedEntities[domain.Article],
) utils.RestfullPaginatedEntities[schemas.ArticleResponse] {
	mappedItems := make([]schemas.ArticleResponse, len(input.Items))
	for i, article := range input.Items {
		mappedItems[i] = MapArticleToArticleResponse(&article)
	}

	return utils.RestfullPaginatedEntities[schemas.ArticleResponse]{
		Items:    mappedItems,
		Next:     input.Next,
		Previous: input.Previous,
	}
}
