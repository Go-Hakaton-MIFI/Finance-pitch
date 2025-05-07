package schemas

import "finance-backend/internal/domain"

type UserRegistrationSchema struct {
	UserType domain.UserType `json:"userType" validate:"required,oneof=ФЛ ЮЛ"`
	Login    string          `json:"loginName" validate:"required,min=3,max=50"`
	Name     string          `json:"partName" validate:"required,min=3,max=50"`
	Password string          `json:"password" validate:"required,min=6"`
	Bank     string          `json:"bank" validate:"required"`
	Account  string          `json:"account" validate:"required,len=20,numeric"`
	INN      string          `json:"inn" validate:"required,len=11,numeric"`
	Phone    string          `json:"phone" validate:"required,e164"`
}

func (us *UserRegistrationSchema) ToDomainEntity() *domain.UserCreationData {
	return &domain.UserCreationData{
		UserType: us.UserType,
		Login:    us.Login,
		Name:     us.Name,
		Password: us.Password,
		Bank:     us.Bank,
		Account:  us.Account,
		INN:      us.INN,
		Phone:    us.Phone,
	}
}

type UserLoginSchema struct {
	Login    string `json:"loginName" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6"`
}
