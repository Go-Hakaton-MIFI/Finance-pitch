package user

import (
	"context"
	"database/sql"
	"errors"
	"finance-backend/internal/domain"
	"finance-backend/pkg/logger"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db  *sqlx.DB
	log *logger.Logger
}

func NewUserRepository(db *sqlx.DB, logger *logger.Logger) *UserRepository {
	return &UserRepository{
		db:  db,
		log: logger,
	}
}

func (ur *UserRepository) CreateUser(ctx context.Context, data *domain.UserCreationData) error {
	tx, err := ur.db.BeginTxx(ctx, nil)
	if err != nil {
		ur.log.Error(ctx, "error starting transaction", map[string]interface{}{
			"error": err,
		})
		return domain.ErrDBConnection
	}
	defer tx.Rollback()

	var participantID int
	err = tx.GetContext(ctx, &participantID, `
        INSERT INTO Participants (part_type, part_name, part_bank, part_account, part_inn, part_phone)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING part_id
    `, data.UserType, data.Name, data.Bank, data.Account, data.INN, data.Phone)
	if err != nil {
		ur.log.Error(ctx, "error inserting participant", map[string]interface{}{
			"error": err,
			"login": data.Login,
		})
		return domain.ErrTypeInsertion
	}

	_, err = tx.ExecContext(ctx, `
        INSERT INTO Users (login_name, password, role, part_id)
        VALUES ($1, $2, 'user', $3)
    `, data.Login, data.Password, participantID)
	if err != nil {
		ur.log.Error(ctx, "error inserting user", map[string]interface{}{
			"error": err,
			"login": data.Login,
		})
		return domain.ErrTypeInsertion
	}

	if err := tx.Commit(); err != nil {
		ur.log.Error(ctx, "error committing transaction", map[string]interface{}{
			"error": err,
			"login": data.Login,
		})
		return domain.ErrDBConnection
	}

	return nil
}

func (ur *UserRepository) GetRawUserByLogin(ctx context.Context, login string) (*domain.RawUser, error) {
	var result struct {
		Login        string `db:"login_name"`
		PasswordHash string `db:"password"`
		Role         string `db:"role"`
	}

	err := ur.db.GetContext(ctx, &result, `
        SELECT login_name, password, role
        FROM Users
        WHERE login_name = $1
    `, login)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		ur.log.Error(ctx, "error querying user", map[string]interface{}{
			"error": err,
			"login": login,
		})
		return nil, domain.ErrDBConnection
	}

	return &domain.RawUser{
		User: domain.User{
			Login:   result.Login,
			IsAdmin: result.Role == "admin",
		},
		PasswordHash: result.PasswordHash,
	}, nil
}

var _ IUserRepository = (*UserRepository)(nil)
