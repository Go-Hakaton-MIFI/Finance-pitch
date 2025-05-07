package schemas

type ArticleResponse struct {
	ID          int64              `json:"id"`
	Header      string             `json:"header"`
	SubHeader   string             `json:"sub_header"`
	Description string             `json:"description"`
	Image       *string            `json:"image"`
	Categories  []CategoryResponse `json:"categories"`
}

type CreateArticleRequest struct {
	Header      string `json:"header" validate:"required,min=2"`
	SubHeader   string `json:"sub_header" validate:"required,min=2"`
	Description string `json:"description" validate:"required,min=5"`
}

type UpdateArticleRequest struct {
	Header      *string `json:"header"`
	SubHeader   *string `json:"sub_header"`
	Description *string `json:"description"`
}

type LinkCategoryRequest struct {
	CategoriesIDs []int `json:"categories_ids" validate:"required"`
}
