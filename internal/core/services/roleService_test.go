package services

import (
	"testing"

	"github.com/AntonyIS/usafi-hub-user-service/config"
	"github.com/AntonyIS/usafi-hub-user-service/internal/adapter/logger"
	"github.com/AntonyIS/usafi-hub-user-service/internal/adapter/repository"
	"github.com/AntonyIS/usafi-hub-user-service/internal/core/domain"
)

func TestRoleService(t *testing.T) {
	logger, err := logger.NewDefaultLogger()
	if err != nil {
		panic(err)
	}

	config, err := config.NewConfig(logger)
	if err != nil {
		panic(err)
	}

	roleRepo, _ := repository.NewRolePostgresClient(*config)
	roleService := NewRoleService(roleRepo)

	t.Run("Testing CreateRole", func(t *testing.T) {
		role := domain.Role{
			Name:        "Customer",
			Description: "UsafiHub Customer",
		}

		newRole, err := roleService.CreateRole(role)
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
		roles, err := roleService.GetRoles()
		if err != nil {
			t.Errorf("error reading roles: %v", err)
		}
		roleId := roles[0].RoleId
		dbRole, err := roleService.GetRoleById(roleId)
		if err != nil {
			t.Errorf("error reading role: %v", err)
		}
		if dbRole == nil {
			t.Error("expected role found nil")
		}
	})

	t.Run("Testing GetRoles", func(t *testing.T) {
		roles, err := roleService.GetRoles()
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
		roles, err := roleService.GetRoles()
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
		err = roleService.UpdateRole(*role)

		if err != nil {
			t.Errorf("error updating role: %v", err)
		}

		dbRole, err := roleService.GetRoleById(roleId)

		if err != nil {
			t.Errorf("error reading role while updating: %v", err)
		}

		if dbRole.Name != "Administrator" {
			t.Errorf("expected role name to be 'Administrator': found %v", dbRole.Name)
		}

	})

	t.Run("Testing deleting role", func(t *testing.T) {
		roles, err := roleService.GetRoles()
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

		err = roleService.DeleteRole(roleId)
		if err != nil {
			t.Errorf("error deleting role: %v", err)
		}

		dbRole, err := roleService.GetRoleById(roleId)
		if err == nil {
			t.Errorf("error reading role: %v", err)
		}
		if dbRole != nil {
			t.Error("expected role to be nil")
		}

	})
}
