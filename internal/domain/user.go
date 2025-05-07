package domain

type UserType string

const (
	UserTypeFL UserType = "ФЛ"
	UserTypeUL UserType = "ЮЛ"
)

func (u UserType) IsValid() bool {
	return u == UserTypeFL || u == UserTypeUL
}

func (u UserType) String() string {
	return string(u)
}

func (UserType) Values() []UserType {
	return []UserType{UserTypeFL, UserTypeUL}
}

type AccessToken string

func (t AccessToken) String() string {
	return string(t)
}

type UserCreationData struct {
	UserType UserType
	Login    string
	Name     string
	Password string
	Bank     string
	Account  string
	INN      string
	Phone    string
}

type User struct {
	Login   string
	IsAdmin bool
}

type RawUser struct {
	User
	PasswordHash string
}
