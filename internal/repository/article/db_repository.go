package article

import (
	"context"
	"database/sql"
	"errors"
	"finance-backend/internal/domain"
	"finance-backend/pkg/logger"
	"finance-backend/pkg/utils"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type ArticleRepository struct {
	db  *sqlx.DB
	log *logger.Logger
}

func NewArticleRepository(logger *logger.Logger, db *sqlx.DB) *ArticleRepository {
	return &ArticleRepository{
		db:  db,
		log: logger,
	}
}

func (r *ArticleRepository) GetByID(ctx context.Context, id int64) (*domain.Article, error) {
	query_article := `
		SELECT 
			id, header, sub_header, description, image, 
			articles.created_at, articles.updated_at
		FROM articles
		WHERE 
			id = ?
	`
	var article domain.Article
	r.log.Info(ctx, "Getting article by id", map[string]interface{}{"article_id": id})
	err := r.db.GetContext(ctx, &article, query_article, id)
	if err != nil {
		r.log.Error(ctx, "error performing db op", map[string]interface{}{"error": err, "article_id": id})
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrArticleNotFound
		}
		return nil, err
	}

	query_categories := `
		SELECT 
			categories.id, categories.name
		FROM articles
		INNER JOIN categoriesArticles on categoriesArticles.article_id = articles.id
		INNER JOIN categories on categoriesArticles.category_id = categories.id
		WHERE 
			articles.id = ?
	`

	r.log.Info(ctx, "Getting article categories", map[string]interface{}{"article_id": id})
	var categories []domain.Category

	err = r.db.SelectContext(ctx, &categories, query_categories, id)
	if err != nil {
		r.log.Error(ctx, "error performing db op", map[string]interface{}{"error": err, "article_id": id})
		return nil, err
	}

	article.Categories = categories

	return &article, nil
}

func (r *ArticleRepository) SearchPaginated(
	ctx context.Context,
	limit int, offset int,
	search *string,
	categoriesIDs []int,
) (utils.RestfullPaginatedEntities[domain.Article], error) {

	var queryRawArticles string
	var countQuery string
	var params []interface{}
	var paramsCount []interface{}

	if len(categoriesIDs) > 0 {
		params = []interface{}{
			search,
			search,
			search,
			search,
			search,
			search,
			categoriesIDs,
			limit,
			offset,
		}

		paramsCount = []interface{}{
			search,
			search,
			search,
			search,
			search,
			search,
			categoriesIDs,
		}

		queryRawArticles, params, _ = sqlx.In(SearchArticlesWithCategoriesQuery, params...)
		countQuery, paramsCount, _ = sqlx.In(CountSearchArticlesWithCategoriesQuery, paramsCount...)
	} else {
		params = []interface{}{
			search,
			search,
			search,
			search,
			search,
			search,
			limit,
			offset,
		}

		paramsCount = []interface{}{
			search,
			search,
			search,
			search,
			search,
			search,
		}

		queryRawArticles = SearchArticlesQuery
		countQuery = CountSearchArticlesQuery
	}

	var rawArticles []RawArticle
	r.log.Info(ctx, "Search articles paginated", map[string]interface{}{"limit": limit, "offset": offset, "search": search})
	err := r.db.SelectContext(
		ctx,
		&rawArticles,
		queryRawArticles,
		params...,
	)
	if err != nil {
		r.log.Error(ctx, "error performing db op", map[string]interface{}{"error": err})
		return utils.RestfullPaginatedEntities[domain.Article]{}, err
	}

	var total int
	err = r.db.GetContext(ctx, &total, countQuery, paramsCount...)
	if err != nil {
		r.log.Error(ctx, "error performing db op", map[string]interface{}{"error": err})
		return utils.RestfullPaginatedEntities[domain.Article]{}, err
	}

	var articles = make([]domain.Article, 0, total)
	var articlesMap = make(map[int64]domain.Article)

	for _, item := range rawArticles {

		if value, ok := articlesMap[item.ID]; ok && item.CategoryId != nil && item.CategoryName != nil {
			value.Categories = append(value.Categories, domain.Category{ID: *item.CategoryId, Name: *item.CategoryName})
		} else {
			article := domain.Article{
				ID:          item.ID,
				Header:      item.Header,
				SubHeader:   item.SubHeader,
				Description: item.Description,
				Image:       item.Image,
				CreatedAt:   item.CreatedAt,
				UpdatedAt:   item.UpdatedAt,
			}

			if item.CategoryId != nil && item.CategoryName != nil {
				article.Categories = []domain.Category{{ID: *item.CategoryId, Name: *item.CategoryName}}
			}
			articlesMap[item.ID] = article
		}
	}

	for _, article := range articlesMap {
		articles = append(articles, article)
	}

	// TODO: utils
	var previous *string
	var next *string

	if total > limit+offset {
		nextOffset := offset + limit
		n := fmt.Sprintf("limit=%d&offset=%d", limit, nextOffset)
		next = &n
	}

	if offset > 0 {
		prevOffset := offset - limit
		if prevOffset < 0 {
			prevOffset = 0
		}
		if offset > total {
			prevOffset = total - limit
		}
		p := fmt.Sprintf("limit=%d&offset=%d", limit, prevOffset)
		previous = &p
	}

	return utils.RestfullPaginatedEntities[domain.Article]{
		Items:    articles,
		Next:     next,
		Previous: previous,
	}, nil
}

func (r *ArticleRepository) Create(
	ctx context.Context,
	header string,
	subHeader string,
	description string,
) (*domain.Article, error) {
	r.log.Info(ctx, "Create article", nil)

	query := `
		INSERT INTO articles 
			(header, sub_header, description) 
		VALUES 
			(?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query, header, subHeader, description)
	if err != nil {
		r.log.Error(ctx, "error performing db op", map[string]interface{}{"error": err})
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		r.log.Error(ctx, "error fetching last inserted id", map[string]interface{}{"error": err})
		return nil, err
	}

	article, err := r.GetByID(ctx, id)
	return article, err
}

func (r *ArticleRepository) Update(
	ctx context.Context,
	id int64,
	header *string,
	image *string,
	subHeader *string,
	description *string,
) error {
	r.log.Info(ctx, "Update article by id", map[string]interface{}{"article_id": id})

	// TODO: util?
	params := make([]interface{}, 0, 6)

	var queryBuilder strings.Builder
	queryBuilder.WriteString("UPDATE articles SET")

	if header != nil {
		queryBuilder.WriteString(" header = ?,")
		params = append(params, *header)
	}
	if image != nil {
		queryBuilder.WriteString(" image = ?,")
		params = append(params, *image)
	}

	if subHeader != nil {
		queryBuilder.WriteString(" sub_header = ?,")
		params = append(params, *subHeader)
	}
	if description != nil {
		queryBuilder.WriteString(" description = ?,")
		params = append(params, *description)
	}

	if len(params) == 0 {
		return nil
	}

	query := strings.TrimSuffix(queryBuilder.String(), ",")
	query += " WHERE id = ?"
	params = append(params, id)

	_, err := r.db.ExecContext(ctx, query, params...)

	if err != nil {
		r.log.Error(ctx, "error performing db op", map[string]interface{}{"error": err, "article_id": id})
		return err
	}

	return err
}

func (r *ArticleRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM articles WHERE id = ?`
	r.log.Info(ctx, "Delete article by id", map[string]interface{}{"article_id": id})
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.log.Error(ctx, "error performing db op", map[string]interface{}{"error": err, "article_id": id})
	}
	return err
}

func (r *ArticleRepository) LinkCategories(ctx context.Context, id int64, categoriesIDs []int) error {
	r.log.Info(ctx, "Update article categories", map[string]interface{}{"article_id": id})

	tx, err := r.db.Beginx()

	if err != nil {
		r.log.Error(ctx, "error creating transaction", map[string]interface{}{"error": err, "article_id": id})
		return err
	}

	defer tx.Rollback()

	deleteQuery := `DELETE FROM categoriesArticles WHERE article_id = ?`
	_, err = tx.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		r.log.Error(ctx, "error performing db op", map[string]interface{}{"error": err, "article_id": id})
		return err
	}

	var articlesCategoriesMap []map[string]interface{}

	for _, i := range categoriesIDs {
		articlesCategoriesMap = append(articlesCategoriesMap, map[string]interface{}{"article_id": id, "category_id": i})
	}

	_, err = tx.NamedExec(`
		INSERT INTO categoriesArticles 
			(article_id, category_id)
		VALUES 
			(:article_id, :category_id)
	`, articlesCategoriesMap)

	if err != nil {
		r.log.Error(ctx, "error performing db op", map[string]interface{}{"error": err, "article_id": id})
		return err
	}

	if err = tx.Commit(); err != nil {
		r.log.Error(ctx, "error commiting transaction", map[string]interface{}{"error": err})
		return err
	}

	return err
}

var _ IArticleRepository = (*ArticleRepository)(nil)
