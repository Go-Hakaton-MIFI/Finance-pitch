package article

import (
	"context"
	"finance-backend/internal/domain"
	"finance-backend/pkg/utils"
)

type IArticleRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Article, error)

	SearchPaginated(
		ctx context.Context,
		limit int, offset int,
		search *string,
		categoriesIDs []int,
	) (utils.RestfullPaginatedEntities[domain.Article], error)

	Create(
		ctx context.Context,
		header string,
		subHeader string,
		description string,
	) (*domain.Article, error)

	Update(
		ctx context.Context,
		id int64,
		header *string,
		image *string,
		subHeader *string,
		description *string,
	) error

	Delete(ctx context.Context, id int64) error

	LinkCategories(ctx context.Context, id int64, categoriesIDs []int) error
}
