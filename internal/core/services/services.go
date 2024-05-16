package services

import (
	"time"

	"github.com/AntonyIS/usafi-hub-user-service/internal/core/domain"
	"github.com/AntonyIS/usafi-hub-user-service/internal/core/ports"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type BaseService struct {
	repo ports.BaseRepository
}

type UserService struct {
	repo ports.UserRepository
}

type RoleService struct {
	repo ports.RoleRepository
}

type UserRoleService struct {
	repo ports.UserRoleRepository
}

func NewBaseService(repo ports.BaseRepository) *BaseService {
	service := BaseService{
		repo: repo,
	}
	return &service
}

func NewUserService(repo ports.UserRepository) *UserService {
	service := UserService{
		repo: repo,
	}
	return &service
}

func NewRoleService(repo ports.RoleRepository) *RoleService {
	service := RoleService{
		repo: repo,
	}
	return &service
}

func NewUserRoleService(repo ports.UserRoleRepository) *UserRoleService {
	service := UserRoleService{
		repo: repo,
	}
	return &service
}

func (svc UserService) CreateUser(user domain.User) (*domain.User, error) {
	user.UserId = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	password, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = string(password)
	return svc.repo.CreateUser(user)
}

func (svc UserService) GetUserById(userId string) (*domain.User, error) {
	return svc.repo.GetUserById(userId)
}

func (svc UserService) LoginUser(email, password string) (*domain.User, error) {
	user, err := svc.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))

	if err == nil {
		return user, nil
	}
	return nil, err
}

func (svc UserService) GetUserByEmail(email string) (*domain.User, error) {
	return svc.repo.GetUserByEmail(email)
}

func (svc UserService) GetUsers() ([]*domain.User, error) {
	return svc.repo.GetUsers()
}

func (svc UserService) GetUsersWithRole(roleName string) ([]*domain.User, error) {
	return svc.repo.GetUsersWithRole(roleName)
}

func (svc UserService) UpdateUser(user domain.User) (*domain.User, error) {
	user.UpdatedAt = time.Now()
	return svc.repo.UpdateUser(user)
}

func (svc UserService) DeleteUser(userId string) error {
	return svc.repo.DeleteUser(userId)
}

func (svc RoleService) CreateRole(role domain.Role) (*domain.Role, error) {
	role.RoleId = uuid.New().String()
	return svc.repo.CreateRole(role)
}

func (svc RoleService) GetRoleById(roleId string) (*domain.Role, error) {
	return svc.repo.GetRoleById(roleId)
}

func (svc RoleService) GetRoles() ([]*domain.Role, error) {
	return svc.repo.GetRoles()
}

func (svc RoleService) UpdateRole(role domain.Role) error {
	return svc.repo.UpdateRole(role)
}

func (svc RoleService) DeleteRole(roleId string) error {
	return svc.repo.DeleteRole(roleId)
}

func (svc UserRoleService) AddUserRole(userRole domain.UserRole) error {
	return svc.repo.AddUserRole(userRole)
}

func (svc UserRoleService) RemoveUserRole(userRole domain.UserRole) error {
	return svc.repo.RemoveUserRole(userRole)
}

func (svc BaseService) DropTables() error {
	return svc.repo.DropTables()
}
