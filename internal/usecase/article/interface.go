package article

import (
	"context"
	"finance-backend/internal/domain"
	"finance-backend/pkg/utils"
	"io"
)

type IArticleUseCase interface {
	GetArticleByID(ctx context.Context, id int64) (*domain.Article, error)

	SearchArticlesPaginated(
		ctx context.Context,
		limit int,
		offset int,
		search *string,
		categoriesIDs []int,
	) (utils.RestfullPaginatedEntities[domain.Article], error)

	CreateArticle(
		ctx context.Context,
		header string,
		subHeader string,
		description string,
	) (*domain.Article, error)

	UpdateArticle(
		ctx context.Context,
		id int64,
		header *string,
		subHeader *string,
		description *string,
	) (*domain.Article, error)

	DeleteArticle(ctx context.Context, id int64) error

	LinkCategories(ctx context.Context, id int64, categoriesIDs []int) (*domain.Article, error)

	LinkImage(ctx context.Context, id int64, body io.ReadSeeker, size int64, contentType string) (*domain.Article, error)
}
