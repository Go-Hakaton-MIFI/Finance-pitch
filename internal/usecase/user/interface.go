package user

import (
	"context"
	"finance-backend/internal/domain"
)

type IUserUseCase interface {
	RegisterUser(ctx context.Context, userCreationData *domain.UserCreationData) (*domain.AccessToken, error)
	GetAccessToken(ctx context.Context, login, password string) (*domain.AccessToken, error)
	GetSubjectTypes() ([]map[string]string, error)
}
