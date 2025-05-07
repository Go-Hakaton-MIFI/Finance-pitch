package schemas

import "time"

type CategoryResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type CategoryAdminResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UpdateOrCreateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=3"`
}
