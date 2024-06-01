package services

import (
	"github.com/AntonyIS/usafi-hub-user-service/internal/core/domain"
	"github.com/AntonyIS/usafi-hub-user-service/internal/core/ports"
	"github.com/google/uuid"
)

type roleService struct {
	repo ports.RoleRepository
}

func NewRoleService(repo ports.RoleRepository) *roleService {
	service := roleService{
		repo: repo,
	}
	return &service
}

func (svc roleService) CreateRole(role domain.Role) (*domain.Role, error) {
	role.RoleId = uuid.New().String()
	return svc.repo.CreateRole(role)
}

func (svc roleService) GetRoleById(roleId string) (*domain.Role, error) {
	return svc.repo.GetRoleById(roleId)
}

func (svc roleService) GetRoles() ([]*domain.Role, error) {
	return svc.repo.GetRoles()
}

func (svc roleService) UpdateRole(role domain.Role) error {
	return svc.repo.UpdateRole(role)
}

func (svc roleService) DeleteRole(roleId string) error {
	return svc.repo.DeleteRole(roleId)
}
