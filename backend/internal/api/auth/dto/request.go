package dto

import "errors"

type CredentialsRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r CredentialsRequest) Validate() error {
	return nil
}

type AuthResponse struct {
	Token string `json:"token"`
	User  any    `json:"user"`
}

var (
	ErrInvalidEmail    = errors.New("invalid email")
	ErrPasswordTooShort = errors.New("password too short")
)
