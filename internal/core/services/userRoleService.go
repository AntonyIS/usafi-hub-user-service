package services

import (
	"github.com/AntonyIS/usafi-hub-user-service/internal/core/domain"
	"github.com/AntonyIS/usafi-hub-user-service/internal/core/ports"
)

type userRoleService struct {
	repo ports.UserRoleRepository
}

func NewUserRoleService(repo ports.UserRoleRepository) *userRoleService {
	service := userRoleService{
		repo: repo,
	}
	return &service
}

func (svc userRoleService) AddUserRole(userRole domain.UserRole) error {
	return svc.repo.AddUserRole(userRole)
}

func (svc userRoleService) RemoveUserRole(userRole domain.UserRole) error {
	return svc.repo.RemoveUserRole(userRole)
}
