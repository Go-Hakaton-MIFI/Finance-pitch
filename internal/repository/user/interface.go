package user

import (
	"context"
	"finance-backend/internal/domain"
)

type IUserRepository interface {
	GetRawUserByLogin(ctx context.Context, login string) (*domain.RawUser, error)
	CreateUser(ctx context.Context, creationData *domain.UserCreationData) error
}
