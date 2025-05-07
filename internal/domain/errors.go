package domain

import (
	"errors"
	"fmt"
)

type DomainError struct {
	Code    string
	Message string
}

func (e *DomainError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

var (
	ErrCategoryNotFound = &DomainError{
		Code:    "CATEGORY_NOT_FOUND",
		Message: "Категория не найдена",
	}

	ErrCategoryExists = &DomainError{
		Code:    "CATEGORY_EXISTS",
		Message: "Категория с таким названием уже существует",
	}

	ErrArticleNotFound = &DomainError{
		Code:    "ARTICLE_NOT_FOUND",
		Message: "Статья не найдена",
	}

	ErrDBConnection = &DomainError{
		Code:    "DB_CONNECTION_ERROR",
		Message: "Не удается подключиться к БД",
	}

	ErrTypeInsertion = &DomainError{
		Code:    "INSERTION_ERROR",
		Message: "Ошибка при выполнении операции",
	}

	ErrWrongLoginOrPassword = &DomainError{
		Code:    "CREDS_INVALID_ERROR",
		Message: "Неверный логин или пароль",
	}

	ErrUserAlreadyExists = &DomainError{
		Code:    "USER_EXISTS_ERROR",
		Message: "Пользователь уже зарегистрирован",
	}

	ErrS3Connection = &DomainError{
		Code:    "S3_CONNECTION_ERROR",
		Message: "Не удается подключиться к файловому хранилищу",
	}

	ErrNotFound = errors.New("entity not found")
)
