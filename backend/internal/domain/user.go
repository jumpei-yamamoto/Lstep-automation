package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrEmailAlreadyUsed = errors.New("email already used")
	ErrInvalidEmail     = errors.New("invalid email")
	ErrEmptyName        = errors.New("name is empty")
)

type User struct {
	ID        uuid.UUID
	Name      string
	Email     string
	CreatedAt time.Time
}

func NewUser(name, email string, now func() time.Time) (*User, error) {
	if name == "" {
		return nil, ErrEmptyName
	}
	if !isValidEmail(email) {
		return nil, ErrInvalidEmail
	}
	return &User{
		ID:        uuid.New(),
		Name:      name,
		Email:     email,
		CreatedAt: now(),
	}, nil
}

func isValidEmail(s string) bool {
	return len(s) >= 3 && len(s) <= 255 && containsAt(s)
}

func containsAt(s string) bool {
	for _, ch := range s {
		if ch == '@' {
			return true
		}
	}
	return false
}