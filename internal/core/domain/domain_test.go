package domain

import (
	"testing"
	"time"
)

func TestUserDomain(t *testing.T) {
	// Test user creation
	t.Run("Test user creation", func(t *testing.T) {
		user := User{
			UserId:       "1",
			Username:     "john_doe",
			PasswordHash: "hashed_password",
			Email:        "john.doe@example.com",
			FullName:     "John Doe",
			PhoneNumber:  "1234567890",
			Avatar:       "avatar_url",
			Address:      "123 Main St",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		if user.UserId != "1" {
			t.Errorf("expected UserId to be '1', got %s", user.UserId)
		}
	})
}

func TestRoleDomain(t *testing.T) {
	// Test role creation
	t.Run("Test role creation", func(t *testing.T) {
		role := Role{
			RoleId:"1",
			Name:"Admin",
			Description: "Admininstrator role",
		}

		if role.Name != "Admin" {
			t.Errorf("expected Name to be 'Admin', got %s", role.Name)
		}
	})
}

func TestUserRoleDomain(t *testing.T) {
	// Test UserRole creation
	t.Run("Test UserRole creation", func(t *testing.T) {
		userRole := UserRole{
			RoleId: "1",
			UserId: "1",
		}

		if userRole.RoleId != "1" {
			t.Errorf("expected RoleId to be '1', got %s", userRole.RoleId)
		}

		if userRole.UserId != "1" {
			t.Errorf("expected UserId to be '1', got %s", userRole.UserId)
		}
	})
}
