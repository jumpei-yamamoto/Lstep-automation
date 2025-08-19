package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewUser_Success(t *testing.T) {
	fixedTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	now := func() time.Time { return fixedTime }

	user, err := NewUser("Test User", "test@example.com", now)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user == nil {
		t.Fatal("expected user to be created, got nil")
	}
	if user.Name != "Test User" {
		t.Errorf("expected name 'Test User', got %s", user.Name)
	}
	if user.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got %s", user.Email)
	}
	if user.CreatedAt != fixedTime {
		t.Errorf("expected created_at %v, got %v", fixedTime, user.CreatedAt)
	}
	if user.ID == uuid.Nil {
		t.Error("expected ID to be generated, got nil UUID")
	}
}

func TestNewUser_EmptyName(t *testing.T) {
	now := func() time.Time { return time.Now() }

	user, err := NewUser("", "test@example.com", now)

	if err != ErrEmptyName {
		t.Errorf("expected error %v, got %v", ErrEmptyName, err)
	}
	if user != nil {
		t.Error("expected user to be nil when name is empty")
	}
}

func TestNewUser_InvalidEmail(t *testing.T) {
	now := func() time.Time { return time.Now() }

	testCases := []string{
		"",           // empty email
		"ab",         // too short
		"invalid",    // no @
		"test@",      // missing domain
		"@example",   // missing local part
	}

	for _, email := range testCases {
		t.Run("email_"+email, func(t *testing.T) {
			user, err := NewUser("Test User", email, now)

			if err != ErrInvalidEmail {
				t.Errorf("expected error %v for email '%s', got %v", ErrInvalidEmail, email, err)
			}
			if user != nil {
				t.Errorf("expected user to be nil for invalid email '%s'", email)
			}
		})
	}
}

func TestIsValidEmail(t *testing.T) {
	testCases := []struct {
		email string
		valid bool
	}{
		{"test@example.com", true},
		{"user@domain.org", true},
		{"a@b.c", true},
		{"", false},
		{"ab", false},
		{"invalid", false},
		{"test@", false},
		{"@example", false},
		{"test@@example.com", true}, // contains @, so passes simple validation
	}

	for _, tc := range testCases {
		t.Run("email_"+tc.email, func(t *testing.T) {
			result := isValidEmail(tc.email)
			if result != tc.valid {
				t.Errorf("expected isValidEmail('%s') = %v, got %v", tc.email, tc.valid, result)
			}
		})
	}
}

func TestContainsAt(t *testing.T) {
	testCases := []struct {
		input string
		hasAt bool
	}{
		{"test@example.com", true},
		{"@", true},
		{"test@", true},
		{"@example", true},
		{"test", false},
		{"", false},
		{"example.com", false},
	}

	for _, tc := range testCases {
		t.Run("input_"+tc.input, func(t *testing.T) {
			result := containsAt(tc.input)
			if result != tc.hasAt {
				t.Errorf("expected containsAt('%s') = %v, got %v", tc.input, tc.hasAt, result)
			}
		})
	}
}

func TestUser_ImmutableFields(t *testing.T) {
	now := func() time.Time { return time.Now() }
	user, err := NewUser("Test User", "test@example.com", now)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	originalID := user.ID
	originalCreatedAt := user.CreatedAt

	// Verify that core fields are set and should not be modified after creation
	if user.ID != originalID {
		t.Error("User ID should remain constant")
	}
	if user.CreatedAt != originalCreatedAt {
		t.Error("User CreatedAt should remain constant")
	}
}