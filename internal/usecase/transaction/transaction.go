package transaction

import (
	"context"
	"finance-backend/internal/domain"
	"finance-backend/internal/repository"
)

type UseCase struct {
	repo *repository.TransactionRepository
}

func NewUseCase(repo *repository.TransactionRepository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) Create(ctx context.Context, transaction *domain.Transaction) error {
	return uc.repo.Create(ctx, transaction)
}

func (uc *UseCase) GetByID(ctx context.Context, id int64) (*domain.Transaction, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *UseCase) ListByUserID(ctx context.Context, userID int64, limit, offset int) ([]*domain.Transaction, error) {
	return uc.repo.ListByUserID(ctx, userID, limit, offset)
}

func (uc *UseCase) Update(ctx context.Context, transaction *domain.Transaction) error {
	return uc.repo.Update(ctx, transaction)
}

func (uc *UseCase) Delete(ctx context.Context, id, userID int64) error {
	return uc.repo.Delete(ctx, id, userID)
}
