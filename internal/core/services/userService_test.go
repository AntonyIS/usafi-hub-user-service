package services

import (
	"testing"
	"time"

	"github.com/AntonyIS/usafi-hub-user-service/config"
	"github.com/AntonyIS/usafi-hub-user-service/internal/adapter/logger"
	"github.com/AntonyIS/usafi-hub-user-service/internal/adapter/repository"
	"github.com/AntonyIS/usafi-hub-user-service/internal/core/domain"
)

func TestUserService(t *testing.T) {
	logger, err := logger.NewDefaultLogger()
	if err != nil {
		panic(err)
	}

	config, err := config.NewConfig(logger)
	if err != nil {
		panic(err)
	}

	repo, _ := repository.NewUserPostgresClient(*config)
	roleRepo, _ := repository.NewRolePostgresClient(*config)
	userRoleRepo, _ := repository.NewUserRolePostgresClient(*config)
	baseRepo, _ := repository.NewBasePostgresClient(*config)

	userService := NewUserService(repo, logger, []byte(config.SECRET_KEY))
	roleService := NewRoleService(roleRepo)
	userRoleService := NewUserRoleService(userRoleRepo)
	baseService := NewBaseService(baseRepo)

	t.Run("Testing CreateUser", func(t *testing.T) {
		user := domain.User{
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
		newUser, err := userService.CreateUser(user)
		if err != nil {
			t.Errorf("error adding user: %v", err)
		}

		if newUser.Email != "john.doe@example.com" {
			t.Errorf("expected at email john.doe@example.com, got : %v", newUser.Email)
		}
	})

	t.Run("Testing GetUsers", func(t *testing.T) {
		users, err := userService.GetUsers()
		if err != nil {
			t.Errorf("error getting users: %v", err)
		}

		if users == nil {
			t.Error("received nil users")
		}

		if len(users) < 1 {
			t.Error("expected at least 1 user, got 0")
		}

	})

	t.Run("Testing GetUserById", func(t *testing.T) {
		dbusers, err := userService.GetUsers()
		if err != nil {
			t.Errorf("error getting users when testing GetUserById: %v", err)
		}

		dbUser := dbusers[0]
		userID := dbUser.UserId

		user, err := userService.GetUserById(userID)

		if err != nil {
			t.Errorf("error getting user with ID %s: %v", userID, err)
		}

		if user == nil {
			t.Errorf("received nil user with ID %s", userID)
		}
	})

	t.Run("Testing GetUserByEmail", func(t *testing.T) {
		email := "john.doe@example.com"
		dbuser, err := userService.GetUserByEmail(email)
		if err != nil {
			t.Errorf("error getting users when testing GetUserById: %v", err)
		}

		if dbuser == nil {
			t.Errorf("received nil user with email %s", email)
		}

		if dbuser.Email != email {
			t.Errorf("expected email to be %s: found %s", email, dbuser.Email)
		}

	})

	t.Run("Testing UpdateUser", func(t *testing.T) {
		user := domain.User{
			Username:     "mary_doe",
			PasswordHash: "hashed_password",
			Email:        "mary.doe@example.com",
			FullName:     "mary Doe",
			PhoneNumber:  "0987654321",
			Avatar:       "avatar_url",
			Address:      "500 Main St",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		newUser, err := userService.CreateUser(user)
		if err != nil {
			t.Fatalf("error adding user: %v", err)
		}
		newEmail := "j.doe@example.com"
		avatarURL := "https://img.freepik.com/free-psd/3d-illustration-person-with-sunglasses_23-2149436188.jpg?size=338&ext=jpg&ga=GA1.1.1369675164.1715385600&semt=ais_user"
		newUser.Avatar = avatarURL
		newUser.Email = newEmail
		updatedUser, err := userService.UpdateUser(*newUser)
		if err != nil {
			t.Fatalf("error updating user: %v", err)
		}

		if updatedUser.Email != newEmail {
			t.Errorf("expected email %s, got %s", newEmail, updatedUser.Email)
		}

		if updatedUser.Avatar != avatarURL {
			t.Errorf("expected email %s, got %s", avatarURL, updatedUser.Avatar)
		}
	})

	t.Run("Testing DeleteUser", func(t *testing.T) {
		user := domain.User{
			Username:     "joe_doe",
			PasswordHash: "hashed_password",
			Email:        "joe.doe@example.com",
			FullName:     "joe Doe",
			PhoneNumber:  "0567654321",
			Avatar:       "avatar_url",
			Address:      "500 Main St",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		newUser, err := userService.CreateUser(user)
		if err != nil {
			t.Fatalf("error adding user: %v", err)
		}

		err = userService.DeleteUser(newUser.UserId)
		if err != nil {
			t.Fatalf("error deleting user: %v", err)
		}

		deletedUser, err := userService.GetUserById(newUser.UserId)
		if deletedUser != nil || err == nil {
			t.Error("expected user to be deleted, but user still exists")
		}
	})

	t.Run("Testing GetUsersWithRole", func(t *testing.T) {
		_, err := userService.GetUsers()
		if err != nil {
			t.Error("error getting users")
		}

		role := domain.Role{
			Name:        "Admin",
			Description: "UsafiHub Administrator",
		}

		newRole, err := roleService.CreateRole(role)
		if err != nil {
			t.Errorf("error adding role: %v", err)
		}

		users, err := userService.GetUsers()
		if err != nil {
			t.Errorf("error reading users while testing GetUsersWithRole: %v", err)
		}

		userRole := domain.UserRole{
			UserId: users[0].UserId,
			RoleId: newRole.RoleId,
		}

		err = userRoleService.AddUserRole(userRole)

		if err != nil {
			t.Errorf("error adding user role %v: %v", userRole, err)
		}
	})

	t.Run("Testing Deleting tables", func(t *testing.T) {
		err := baseService.DropTables()
		if err != nil {
			t.Errorf("error deleting tables: %v", err)
		}
	})

}
