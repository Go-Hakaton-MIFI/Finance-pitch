package category

import (
	"context"
	"finance-backend/internal/domain"
	"finance-backend/pkg/utils"
)

type ICategoryRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Category, error)

	GetByName(ctx context.Context, name string) (*domain.Category, error)

	SearchPaginated(ctx context.Context, limit int, offset int, search *string) (utils.PaginatedEntities[domain.Category], error)

	SearchFlat(ctx context.Context, search *string) ([]domain.Category, error)

	Create(ctx context.Context, name string) (*domain.Category, error)

	UpdateName(ctx context.Context, categoryID int64, categoryName string) error

	Delete(ctx context.Context, id int64) error
}
