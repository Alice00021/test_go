package entity

type User struct {
	Entity
	Name        string
	Surname     string
	Username    string
	Password    string
	Email       string
	IsVerified  bool
	VerifyToken *string
	FilePath    *string
	Rating      float32
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type CreateUserInput struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UpdateUserInput struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	Surname  string  `json:"surname"`
	Username string  `json:"username"`
	Rating   float32 `json:"rating"`
}
type ChangePasswordInput struct {
	ID              int64  `json:"id"`
	OldPassword     string `json:"oldPassword"`
	NewPassword     string `json:"newPassword"`
	ConfirmPassword string `json:"confirmPassword"`
}

func NewUser(name, surname, username, password, email string) *User {
	return &User{
		Name:       name,
		Surname:    surname,
		Username:   username,
		Password:   password,
		Email:      email,
		IsVerified: false,
	}
}

type FilterUserInput struct {
	IsVerified *bool `json:"is_verified"`
}
