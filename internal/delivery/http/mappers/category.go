package mappers

import (
	"finance-backend/internal/delivery/http/schemas"
	"finance-backend/internal/domain"
	"finance-backend/pkg/utils"
)

func MapCategoryToCategoryAdminResponse(category *domain.Category) schemas.CategoryAdminResponse {
	return schemas.CategoryAdminResponse{
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
}

func MapCategoryToCategoryACommonResponse(category *domain.Category) schemas.CategoryResponse {
	return schemas.CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
	}
}

func MapPaginatedCategoriesToAdminResponse(
	input utils.PaginatedEntities[domain.Category],
) utils.PaginatedEntities[schemas.CategoryAdminResponse] {
	mappedItems := make([]schemas.CategoryAdminResponse, len(input.Items))
	for i, category := range input.Items {
		mappedItems[i] = schemas.CategoryAdminResponse{
			ID:        category.ID,
			Name:      category.Name,
			CreatedAt: category.CreatedAt,
			UpdatedAt: category.UpdatedAt,
		}
	}

	return utils.PaginatedEntities[schemas.CategoryAdminResponse]{
		Items:            mappedItems,
		Total:            input.Total,
		PageNumber:       input.PageNumber,
		ObjectsCount:     input.ObjectsCount,
		ObjectsCounTotal: input.ObjectsCounTotal,
		PageCount:        input.PageCount,
	}
}

func MapCategoriesToCategoriesRepsonse(categories []domain.Category) []schemas.CategoryResponse {
	mappedItems := make([]schemas.CategoryResponse, len(categories))
	for i, category := range categories {
		mappedItems[i] = schemas.CategoryResponse{
			ID:   category.ID,
			Name: category.Name,
		}
	}
	return mappedItems
}
