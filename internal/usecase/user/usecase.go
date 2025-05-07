package user

import (
	"context"
	"crypto/rsa"
	"finance-backend/internal/domain"
	repo "finance-backend/internal/repository/user"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	repo       repo.IUserRepository
	privateKey *rsa.PrivateKey
	tokenTTL   time.Duration
}

func NewUserUseCase(repo repo.IUserRepository, privateKey *rsa.PrivateKey, tokenTTL time.Duration) *UserUseCase {
	return &UserUseCase{
		repo:       repo,
		privateKey: privateKey,
		tokenTTL:   tokenTTL,
	}
}

func (u *UserUseCase) GetSubjectTypes() ([]map[string]string, error) {
	return []map[string]string{
		{
			"id":   "individual",
			"name": "Физическое лицо",
		},
		{
			"id":   "legal",
			"name": "Юридическое лицо",
		},
	}, nil
}

func (u *UserUseCase) RegisterUser(ctx context.Context, data *domain.UserCreationData) (*domain.AccessToken, error) {
	existingUser, err := u.repo.GetRawUserByLogin(ctx, data.Login)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	data.Password = string(hashedPassword)
	if err := u.repo.CreateUser(ctx, data); err != nil {
		return nil, err
	}

	return u.issueToken(data.Login, false)
}

func (u *UserUseCase) GetAccessToken(ctx context.Context, login, password string) (*domain.AccessToken, error) {
	rawUser, err := u.repo.GetRawUserByLogin(ctx, login)
	if err != nil {
		return nil, err
	}
	if rawUser == nil {
		return nil, domain.ErrWrongLoginOrPassword
	}

	if err := bcrypt.CompareHashAndPassword([]byte(rawUser.PasswordHash), []byte(password)); err != nil {
		return nil, domain.ErrWrongLoginOrPassword
	}

	return u.issueToken(rawUser.Login, rawUser.IsAdmin)
}

func (u *UserUseCase) issueToken(login string, isAdmin bool) (*domain.AccessToken, error) {
	claims := jwt.MapClaims{
		"sub":      login,
		"is_admin": isAdmin,
		"exp":      time.Now().Add(u.tokenTTL).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signed, err := token.SignedString(u.privateKey)
	if err != nil {
		return nil, err
	}

	t := domain.AccessToken(signed)

	return &t, nil
}

var _ IUserUseCase = (*UserUseCase)(nil)
