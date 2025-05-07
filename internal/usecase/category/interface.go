package category

import (
	"context"
	"finance-backend/internal/domain"
	"finance-backend/pkg/utils"
)

type ICategoryUseCase interface {
	GetCategoryByID(ctx context.Context, id int64) (*domain.Category, error)

	SearchCategoriesPaginated(ctx context.Context, limit int, offset int, search *string) (utils.PaginatedEntities[domain.Category], error)

	SearchCategoriesFlat(ctx context.Context, search *string) ([]domain.Category, error)

	CreateCategory(ctx context.Context, name string) (*domain.Category, error)

	UpdateCategoryName(ctx context.Context, categoryID int64, categoryName string) (*domain.Category, error)

	DeleteCategory(ctx context.Context, id int64) error
}
