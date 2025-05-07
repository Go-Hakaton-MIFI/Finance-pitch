package transaction

import (
	"context"
	"finance-backend/internal/delivery/http/schemas"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetTransactions(ctx context.Context, filter schemas.TransactionFilter) ([]schemas.Transaction, error) {
	domainFilter := &TransactionFilter{
		UserType:      filter.UserType,
		TransType:     filter.TransType,
		SenderBank:    filter.SenderBank,
		ReceiverINN:   filter.ReceiverINN,
		ReceiverPhone: filter.ReceiverPhone,
		DateFrom:      filter.DateFrom,
		DateTo:        filter.DateTo,
		CategoryID:    filter.CategoryID,
		StatusID:      filter.StatusID,
	}

	transactions, err := s.repo.GetTransactions(ctx, domainFilter)
	if err != nil {
		return nil, err
	}

	result := make([]schemas.Transaction, len(transactions))
	for i, t := range transactions {
		result[i] = schemas.Transaction{
			ID:            t.ID,
			UserType:      t.UserType,
			DateTime:      t.DateTime,
			TransType:     t.TransType,
			Amount:        t.Amount,
			CategoryID:    t.CategoryID,
			StatusID:      t.StatusID,
			SenderBank:    t.SenderBank,
			ReceiverINN:   t.ReceiverINN,
			ReceiverPhone: t.ReceiverPhone,
			Comment:       t.Comment,
			CategoryName:  t.CategoryName,
			StatusName:    t.StatusName,
		}
	}
	return result, nil
}

func (s *service) GetPreparedTransactions(ctx context.Context) ([]schemas.PreparedTransaction, error) {
	transactions, err := s.repo.GetPreparedTransactions(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]schemas.PreparedTransaction, len(transactions))
	for i, t := range transactions {
		result[i] = schemas.PreparedTransaction{
			ID:            t.ID,
			UserType:      t.UserType,
			DateTime:      t.DateTime,
			TransType:     t.TransType,
			Amount:        t.Amount,
			CategoryID:    t.CategoryID,
			StatusID:      t.StatusID,
			SenderBank:    t.SenderBank,
			ReceiverINN:   t.ReceiverINN,
			ReceiverPhone: t.ReceiverPhone,
			Comment:       t.Comment,
		}
	}
	return result, nil
}

func (s *service) GetCategories(ctx context.Context) ([]schemas.Category, error) {
	categories, err := s.repo.GetCategories(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]schemas.Category, len(categories))
	for i, c := range categories {
		result[i] = schemas.Category{
			ID:   c.ID,
			Name: c.Name,
			Type: c.Type,
		}
	}
	return result, nil
}

func (s *service) GetTransactionStatuses(ctx context.Context) ([]schemas.TransactionStatus, error) {
	statuses, err := s.repo.GetTransactionStatuses(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]schemas.TransactionStatus, len(statuses))
	for i, st := range statuses {
		result[i] = schemas.TransactionStatus{
			ID:   st.ID,
			Name: st.Name,
		}
	}
	return result, nil
}

func (s *service) DeleteTransaction(ctx context.Context, id int64) error {
	return s.repo.DeleteTransaction(ctx, int(id))
}

func (s *service) CreateTransaction(ctx context.Context, transaction schemas.Transaction) (schemas.Transaction, error) {
	domainTransaction := &Transaction{
		ID:            transaction.ID,
		UserType:      transaction.UserType,
		DateTime:      transaction.DateTime,
		TransType:     transaction.TransType,
		Amount:        transaction.Amount,
		CategoryID:    transaction.CategoryID,
		StatusID:      transaction.StatusID,
		SenderBank:    transaction.SenderBank,
		ReceiverINN:   transaction.ReceiverINN,
		ReceiverPhone: transaction.ReceiverPhone,
		Comment:       transaction.Comment,
	}

	err := s.repo.CreateTransaction(ctx, domainTransaction)
	if err != nil {
		return schemas.Transaction{}, err
	}

	return transaction, nil
}

func (s *service) CreatePreparedTransaction(ctx context.Context, transaction schemas.PreparedTransaction) (schemas.PreparedTransaction, error) {
	domainTransaction := &PreparedTransaction{
		ID:            transaction.ID,
		UserType:      transaction.UserType,
		DateTime:      transaction.DateTime,
		TransType:     transaction.TransType,
		Amount:        transaction.Amount,
		CategoryID:    transaction.CategoryID,
		StatusID:      transaction.StatusID,
		SenderBank:    transaction.SenderBank,
		ReceiverINN:   transaction.ReceiverINN,
		ReceiverPhone: transaction.ReceiverPhone,
		Comment:       transaction.Comment,
	}

	err := s.repo.CreatePreparedTransaction(ctx, domainTransaction)
	if err != nil {
		return schemas.PreparedTransaction{}, err
	}

	return transaction, nil
}
