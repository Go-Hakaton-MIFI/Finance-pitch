package category

import (
	"context"
	"finance-backend/internal/domain"
	"finance-backend/pkg/logger"
	"finance-backend/pkg/utils"
	"time"

	"github.com/jmoiron/sqlx"
)

type CategoryRepository struct {
	db  *sqlx.DB
	log *logger.Logger
}

func NewCategoryRepository(logger *logger.Logger, db *sqlx.DB) *CategoryRepository {
	return &CategoryRepository{
		db:  db,
		log: logger,
	}
}

func (r *CategoryRepository) GetByID(ctx context.Context, id int64) (*domain.Category, error) {
	query := `SELECT id, name, created_at, updated_at FROM categories WHERE id = ?`
	var category domain.Category
	r.log.Info(ctx, "Getting category by id", nil)
	err := r.db.GetContext(ctx, &category, query, id)
	if err != nil {
		r.log.Error(ctx, "error performing db op", map[string]interface{}{"error": err, "category_id": id})
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) GetByName(ctx context.Context, name string) (*domain.Category, error) {
	query := `SELECT id, name, created_at, updated_at AS updated_at FROM categories WHERE name = ?`
	var category domain.Category
	r.log.Info(ctx, "Getting category by name", nil)
	err := r.db.GetContext(ctx, &category, query, name)
	if err != nil {
		r.log.Error(ctx, "error performing db op", map[string]interface{}{"error": err, "category_name": name})
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) SearchPaginated(ctx context.Context, limit int, offset int, search *string) (utils.PaginatedEntities[domain.Category], error) {
	query := `SELECT id, name, created_at, updated_at FROM categories WHERE (? IS NULL OR LOWER(name) LIKE CONCAT('%', LOWER(?), '%')) LIMIT ? OFFSET ?`
	var categories []domain.Category
	r.log.Info(ctx, "Search category paginated", map[string]interface{}{"limit": limit, "offset": offset, "search": search})
	err := r.db.SelectContext(ctx, &categories, query, search, search, limit, offset)
	if err != nil {
		r.log.Error(ctx, "error performing db op", map[string]interface{}{"error": err})
		return utils.PaginatedEntities[domain.Category]{}, err
	}

	countQuery := `SELECT COUNT(*) FROM categories WHERE (? IS NULL OR LOWER(name) LIKE CONCAT('%', LOWER(?), '%'))`
	var total int
	err = r.db.GetContext(ctx, &total, countQuery, search, search)
	if err != nil {
		r.log.Error(ctx, "error performing db op", map[string]interface{}{"error": err})
		return utils.PaginatedEntities[domain.Category]{}, err
	}

	pageCount := (total + limit - 1) / limit // Вычисляем количество страниц

	return utils.PaginatedEntities[domain.Category]{
		Items:            categories,
		Total:            total,
		PageNumber:       offset/limit + 1,
		ObjectsCount:     len(categories),
		ObjectsCounTotal: total,
		PageCount:        pageCount,
	}, nil
}

func (r *CategoryRepository) SearchFlat(ctx context.Context, search *string) ([]domain.Category, error) {
	query := `SELECT id, name, created_at, updated_at FROM categories WHERE (? IS NULL OR LOWER(name) LIKE CONCAT('%', LOWER(?), '%'))`
	var category []domain.Category
	r.log.Info(ctx, "Search categories flat", nil)
	err := r.db.SelectContext(ctx, &category, query, search, search)
	if err != nil {
		r.log.Error(ctx, "error performing db op", map[string]interface{}{"error": err})
		return []domain.Category{}, err
	}
	return category, nil
}

func (r *CategoryRepository) Create(ctx context.Context, name string) (*domain.Category, error) {
	query := `INSERT INTO categories (name) VALUES (?)`
	r.log.Info(ctx, "Create category", nil)
	result, err := r.db.ExecContext(ctx, query, name)
	if err != nil {
		r.log.Error(ctx, "error performing db op", map[string]interface{}{"error": err, "category_name": name})
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		r.log.Error(ctx, "error fetching last inserted id", map[string]interface{}{"error": err, "category_name": name})
		return nil, err
	}
	category, err := r.GetByID(ctx, id)
	return category, err
}

func (r *CategoryRepository) UpdateName(ctx context.Context, categoryID int64, categoryName string) error {
	query := `UPDATE categories SET name = ?, updated_at = ? WHERE id = ?`
	r.log.Info(ctx, "Update category by id", nil)
	_, err := r.db.ExecContext(ctx, query, categoryName, time.Now(), categoryID)
	if err != nil {
		r.log.Error(ctx, "error performing db op", map[string]interface{}{"error": err, "category_name": categoryName})
	}
	return err
}

func (r *CategoryRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM categories WHERE id = ?`
	r.log.Info(ctx, "Delete category by id", nil)
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.log.Error(ctx, "error performing db op", map[string]interface{}{"error": err, "category_id": id})
	}
	return err
}
