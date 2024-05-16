package services

import (
	"testing"
	"time"

	"github.com/AntonyIS/usafi-hub-user-service/config"
	"github.com/AntonyIS/usafi-hub-user-service/internal/adapter/logger"
	"github.com/AntonyIS/usafi-hub-user-service/internal/adapter/repository"
	"github.com/AntonyIS/usafi-hub-user-service/internal/core/domain"
)

type testingServices struct {
	baseService     BaseService
	userService     UserService
	roleService     RoleService
	userRoleService UserRoleService
}

func initServices() testingServices {
	logger, err := logger.NewDefaultLogger()
	if err != nil {
		panic(err)
	}

	config, err := config.NewConfig(logger)
	if err != nil {
		panic(err)
	}

	baseRepo, _ := repository.NewBasePostgresClient(*config, logger)
	roleRepo, _ := repository.NewRolePostgresClient(*config, logger)
	userRepo, _ := repository.NewUserPostgresClient(*config, logger)
	userRoleRepo, _ := repository.NewUserRolePostgresClient(*config, logger)

	baseService := NewBaseService(baseRepo)
	userService := NewUserService(userRepo)
	roleService := NewRoleService(roleRepo)
	userRoleService := NewUserRoleService(userRoleRepo)

	return testingServices{
		baseService:     *baseService,
		userService:     *userService,
		roleService:     *roleService,
		userRoleService: *userRoleService,
	}

}

func TestUserManagementService(t *testing.T) {
	services := initServices()

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
		newUser, err := services.userService.CreateUser(user)
		if err != nil {
			t.Errorf("error adding user: %v", err)
		}

		if newUser.Email != "john.doe@example.com" {
			t.Errorf("expected at email john.doe@example.com, got : %v", newUser.Email)
		}
	})

	t.Run("Testing GetUsers", func(t *testing.T) {
		users, err := services.userService.GetUsers()
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
		dbusers, err := services.userService.GetUsers()
		if err != nil {
			t.Errorf("error getting users when testing GetUserById: %v", err)
		}

		dbUser := dbusers[0]
		userID := dbUser.UserId

		user, err := services.userService.GetUserById(userID)

		if err != nil {
			t.Errorf("error getting user with ID %s: %v", userID, err)
		}

		if user == nil {
			t.Errorf("received nil user with ID %s", userID)
		}
	})

	t.Run("Testing GetUserByEmail", func(t *testing.T) {
		email := "john.doe@example.com"
		dbuser, err := services.userService.GetUserByEmail(email)
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
		newUser, err := services.userService.CreateUser(user)
		if err != nil {
			t.Fatalf("error adding user: %v", err)
		}
		newEmail := "j.doe@example.com"
		avatarURL := "https://img.freepik.com/free-psd/3d-illustration-person-with-sunglasses_23-2149436188.jpg?size=338&ext=jpg&ga=GA1.1.1369675164.1715385600&semt=ais_user"
		newUser.Avatar = avatarURL
		newUser.Email = newEmail
		updatedUser, err := services.userService.UpdateUser(*newUser)
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
		newUser, err := services.userService.CreateUser(user)
		if err != nil {
			t.Fatalf("error adding user: %v", err)
		}

		err = services.userService.DeleteUser(newUser.UserId)
		if err != nil {
			t.Fatalf("error deleting user: %v", err)
		}

		deletedUser, err := services.userService.GetUserById(newUser.UserId)
		if deletedUser != nil || err == nil {
			t.Error("expected user to be deleted, but user still exists")
		}
	})

	t.Run("Testing GetUsersWithRole", func(t *testing.T) {
		_, err := services.userService.GetUsers()
		if err != nil {
			t.Error("error getting users")
		}

		role := domain.Role{
			Name:        "Admin",
			Description: "UsafiHub Administrator",
		}

		newRole, err := services.roleService.CreateRole(role)
		if err != nil {
			t.Errorf("error adding role: %v", err)
		}

		users, err := services.userService.GetUsers()
		if err != nil {
			t.Errorf("error reading users while testing GetUsersWithRole: %v", err)
		}

		userRole := domain.UserRole{
			UserId: users[0].UserId,
			RoleId: newRole.RoleId,
		}

		err = services.userRoleService.AddUserRole(userRole)

		if err != nil {
			t.Errorf("error adding user role %v: %v", userRole, err)
		}
	})
	// t.Run("Testing Deleting tables", func(t *testing.T) {
	// 	err := services.baseService.DropTables()
	// 	if err != nil {
	// 		t.Errorf("error deleting tables: %v", err)
	// 	}
	// })
}

func TestRoleManagementService(t *testing.T) {
	services := initServices()
	t.Run("Testing CreateRole", func(t *testing.T) {
		role := domain.Role{
			Name:        "Customer",
			Description: "UsafiHub Customer",
		}

		newRole, err := services.roleService.CreateRole(role)
		if err != nil {
			t.Errorf("error adding role: %v", err)
		}

		if newRole.Name != role.Name {
			t.Errorf("expected role name : %s found %s ", role.Name, newRole.Name)
		}

		if newRole.Description != role.Description {
			t.Errorf("expected role name : %s found %s ", role.Description, newRole.Description)
		}
	})

	t.Run("Testing GetRoleById", func(t *testing.T) {
		roles, err := services.roleService.GetRoles()
		if err != nil {
			t.Errorf("error reading roles: %v", err)
		}
		roleId := roles[0].RoleId
		dbRole, err := services.roleService.GetRoleById(roleId)
		if err != nil {
			t.Errorf("error reading role: %v", err)
		}
		if dbRole == nil {
			t.Error("expected role found nil")
		}
	})

	t.Run("Testing GetRoles", func(t *testing.T) {
		roles, err := services.roleService.GetRoles()
		if err != nil {
			t.Errorf("error reading roles: %v", err)
		}

		if roles == nil {
			t.Error("received nil roles")
		}

		if len(roles) < 1 {
			t.Error("expected at least 1 role, got 0")
		}
	})

	t.Run("Testing UpdateRole", func(t *testing.T) {
		roles, err := services.roleService.GetRoles()
		if err != nil {
			t.Errorf("error reading roles: %v", err)
		}

		if roles == nil {
			t.Error("received nil roles")
		}

		if len(roles) < 1 {
			t.Error("expected at least 1 role, got 0")
		}
		role := roles[0]
		roleId := role.RoleId
		role.Name = "Administrator"
		err = services.roleService.UpdateRole(*role)

		if err != nil {
			t.Errorf("error updating role: %v", err)
		}

		dbRole, err := services.roleService.GetRoleById(roleId)

		if err != nil {
			t.Errorf("error reading role while updating: %v", err)
		}

		if dbRole.Name != "Administrator" {
			t.Errorf("expected role name to be 'Administrator': found %v", dbRole.Name)
		}

	})

	t.Run("Testing deleting role", func(t *testing.T) {
		roles, err := services.roleService.GetRoles()
		if err != nil {
			t.Errorf("error reading roles: %v", err)
		}

		if roles == nil {
			t.Error("received nil roles")
		}

		if len(roles) < 1 {
			t.Error("expected at least 1 role, got 0")
		}
		role := roles[0]
		roleId := role.RoleId

		err = services.roleService.DeleteRole(roleId)
		if err != nil {
			t.Errorf("error deleting role: %v", err)
		}

		dbRole, err := services.roleService.GetRoleById(roleId)
		if err == nil {
			t.Errorf("error reading role: %v", err)
		}
		if dbRole != nil {
			t.Error("expected role to be nil")
		}

	})
}
