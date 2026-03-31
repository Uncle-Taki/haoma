package user

import (
	"testing"
)

func TestNewUser(t *testing.T) {
	u, err := NewUser("Alice", "alice@example.com", "secret123", "555-1234", "123 Main St")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if u.Name != "Alice" {
		t.Errorf("Expected name Alice, got %s", u.Name)
	}
	if u.Email != "alice@example.com" {
		t.Errorf("Expected email alice@example.com, got %s", u.Email)
	}
	if u.Phone != "555-1234" {
		t.Errorf("Expected phone 555-1234, got %s", u.Phone)
	}
	if u.Address != "123 Main St" {
		t.Errorf("Expected address 123 Main St, got %s", u.Address)
	}
	if u.PasswordHash == "" {
		t.Error("Expected password hash to be set")
	}
	if u.PasswordHash == "secret123" {
		t.Error("Password should be hashed, not stored in plain text")
	}
	if u.ID.String() == "" {
		t.Error("Expected UUID to be generated")
	}
}

func TestUser_ValidatePassword(t *testing.T) {
	u, err := NewUser("Bob", "bob@example.com", "mypassword", "", "")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !u.ValidatePassword("mypassword") {
		t.Error("Expected valid password to pass validation")
	}
	if u.ValidatePassword("wrongpassword") {
		t.Error("Expected invalid password to fail validation")
	}
}
