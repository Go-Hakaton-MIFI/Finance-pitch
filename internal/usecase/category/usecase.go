package category

import (
	"context"

	"finance-backend/internal/domain"
	"finance-backend/internal/repository/category"
	"finance-backend/pkg/logger"
	"finance-backend/pkg/utils"
)

type CategoryUseCase struct {
	repo category.ICategoryRepository
	log  *logger.Logger
}

func NewCategoryUseCase(logger *logger.Logger, repo category.ICategoryRepository) *CategoryUseCase {
	return &CategoryUseCase{
		log:  logger,
		repo: repo,
	}
}

func (uc *CategoryUseCase) GetCategoryByID(ctx context.Context, id int64) (*domain.Category, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *CategoryUseCase) SearchCategoriesPaginated(ctx context.Context, limit int, offset int, search *string) (utils.PaginatedEntities[domain.Category], error) {
	return uc.repo.SearchPaginated(ctx, limit, offset, search)
}

func (uc *CategoryUseCase) SearchCategoriesFlat(ctx context.Context, search *string) ([]domain.Category, error) {
	return uc.repo.SearchFlat(ctx, search)
}

func (uc *CategoryUseCase) CreateCategory(ctx context.Context, name string) (*domain.Category, error) {
	existingCategory, _ := uc.repo.GetByName(ctx, name)

	if existingCategory != nil {
		return nil, domain.ErrCategoryExists
	}

	return uc.repo.Create(ctx, name)
}

func (uc *CategoryUseCase) UpdateCategoryName(ctx context.Context, categoryID int64, categoryName string) (*domain.Category, error) {
	category, err := uc.repo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, domain.ErrCategoryNotFound
	}

	if category.Name == categoryName {
		return category, nil
	}

	existingCategory, _ := uc.repo.GetByName(ctx, categoryName)

	if existingCategory != nil {
		return nil, domain.ErrCategoryExists
	}

	uc.repo.UpdateName(ctx, categoryID, categoryName)

	return uc.repo.GetByID(ctx, categoryID)
}

func (uc *CategoryUseCase) DeleteCategory(ctx context.Context, id int64) error {
	category, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if category == nil {
		return domain.ErrCategoryNotFound
	}

	return uc.repo.Delete(ctx, id)
}
