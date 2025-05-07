package article

var (
	SearchArticlesWithCategoriesQuery = `
		SELECT 
			articles.id, header, sub_header, description, image, 
			articles.created_at, articles.updated_at, 
			categories.name as category_name, categories.id as category_id
		FROM articles
		LEFT JOIN categoriesArticles on categoriesArticles.article_id = articles.id
		LEFT JOIN categories on categoriesArticles.category_id = categories.id
		WHERE 
			(? IS NULL OR (
				LOWER(header) LIKE CONCAT('%', LOWER(?), '%') or 
				LOWER(sub_header) LIKE CONCAT('%', LOWER(?), '%') or 
				LOWER(description) LIKE CONCAT('%', LOWER(?), '%') or 
			))
			AND categoriesArticles.category_id IN (?)
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	SearchArticlesQuery = `
		SELECT 
			articles.id, header, sub_header, description, image, 
			articles.created_at, articles.updated_at, 
			categories.name as category_name, categories.id as category_id
		FROM articles
		LEFT JOIN categoriesArticles on categoriesArticles.article_id = articles.id
		LEFT JOIN categories on categoriesArticles.category_id = categories.id
		WHERE 
			? IS NULL OR (
				LOWER(header) LIKE CONCAT('%', LOWER(?), '%') or 
				LOWER(sub_header) LIKE CONCAT('%', LOWER(?), '%') or 
				LOWER(description) LIKE CONCAT('%', LOWER(?), '%') or 
			)
		ORDER BY articles.id DESC
		LIMIT ? OFFSET ?
	`

	CountSearchArticlesWithCategoriesQuery = `
		SELECT 
			COUNT(*)
		FROM articles
		LEFT JOIN categoriesArticles on categoriesArticles.article_id = articles.id
		LEFT JOIN categories on categoriesArticles.category_id = categories.id
		WHERE 
			? IS NULL OR (
				LOWER(header) LIKE CONCAT('%', LOWER(?), '%') or 
				LOWER(sub_header) LIKE CONCAT('%', LOWER(?), '%') or 
				LOWER(description) LIKE CONCAT('%', LOWER(?), '%') or 
			)
			AND categories.id IN (?)
	`

	CountSearchArticlesQuery = `
		SELECT 
			COUNT(*)
		FROM articles
		LEFT JOIN categoriesArticles on categoriesArticles.article_id = articles.id
		LEFT JOIN categories on categoriesArticles.category_id = categories.id
		WHERE 
			? IS NULL OR (
				LOWER(header) LIKE CONCAT('%', LOWER(?), '%') or 
				LOWER(sub_header) LIKE CONCAT('%', LOWER(?), '%') or 
				LOWER(description) LIKE CONCAT('%', LOWER(?), '%') or 
			)
	`
)
