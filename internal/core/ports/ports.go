package ports

import (
	"github.com/AntonyIS/usafi-hub-user-service/internal/core/domain"
)

type UserService interface {
	CreateUser(user domain.User) (*domain.User, error)
	GetUsersWithRole(roleName string) ([]*domain.User, error)
	GetUserById(userId string) (*domain.User, error)
	GetUsers() ([]*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
	UpdateUser(user domain.User) (*domain.User, error)
	DeleteUser(userId string) error
	LoginUser(email, password string) (*domain.User, error)
}

type RoleService interface {
	CreateRole(role domain.Role) (*domain.Role, error)
	GetRoleById(roleId string) (*domain.Role, error)
	GetRoles() ([]*domain.Role, error)
	UpdateRole(role domain.Role) error
	DeleteRole(roleId string) error
}

type UserRoleService interface {
	AddUserRole(userRole domain.UserRole) error
	RemoveUserRole(userRole domain.UserRole) error
}

type UserRepository interface {
	CreateUser(user domain.User) (*domain.User, error)
	GetUsersWithRole(roleName string) ([]*domain.User, error)
	GetUserById(userId string) (*domain.User, error)
	GetUsers() ([]*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
	UpdateUser(user domain.User) (*domain.User, error)
	DeleteUser(userId string) error
}

type RoleRepository interface {
	CreateRole(role domain.Role) (*domain.Role, error)
	GetRoleById(roleId string) (*domain.Role, error)
	GetRoles() ([]*domain.Role, error)
	UpdateRole(role domain.Role) error
	DeleteRole(roleId string) error
}

type UserRoleRepository interface {
	AddUserRole(userRole domain.UserRole) error
	RemoveUserRole(userRole domain.UserRole) error
}

type LoggerService interface {
	Info(message string)
	Warning(message string)
	Error(message string)
}

type BaseRepository interface {
	DropTables() error
}

type BaseService interface {
	DropTables() error
}
