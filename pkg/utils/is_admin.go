package utils

import (
	"context"
	"finance-backend/internal/domain"
)

func CheckIsAdmin(context context.Context) bool {
	user, ok := context.Value(ContextKeyUser).(domain.User)
	return ok && user.IsAdmin
}
