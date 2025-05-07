package delivery

import "fmt"

type DelivaryError struct {
	Code    string
	Message string
}

func (e *DelivaryError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

var (
	ErrJwtTokenParsing = &DelivaryError{
		Code:    "JWT_TOKEN_ERROR",
		Message: "Ошибка парсинга jwt токена",
	}
)
