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
	
	// Workflow related errors
	ErrWorkflowNameEmpty         = errors.New("workflow name is empty")
	ErrWorkflowInactive          = errors.New("workflow is inactive")
	ErrStepNotFound              = errors.New("step not found")
	ErrDuplicateStepOrder        = errors.New("duplicate step order")
	ErrInvalidWorkflowTransition = errors.New("invalid workflow status transition")
	ErrCannotModifyActiveWorkflow = errors.New("cannot modify active workflow")
	ErrStepsRequired             = errors.New("workflow must have at least one step")
	
	// Step related errors
	ErrInvalidStepType = errors.New("invalid step type")
	ErrEmptyStepName   = errors.New("step name is empty")
	ErrInvalidOrder    = errors.New("invalid step order")
	ErrInvalidConfig   = errors.New("invalid step configuration")
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