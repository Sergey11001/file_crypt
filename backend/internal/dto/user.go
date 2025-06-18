package dto

import (
	"github.com/google/uuid"
)

type User struct {
	UUID      uuid.UUID `json:"uuid"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	PublicKey string    `json:"public_key"`
}

type SignUpInput struct {
	Name      string
	Email     string
	Password  string
	PublicKey []byte
}

type SignInInput struct {
	Email    string
	Password string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}
