package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/AntonyIS/usafi-hub-user-service/internal/core/domain"
	"github.com/AntonyIS/usafi-hub-user-service/internal/core/ports"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo   ports.UserRepository
	logger ports.LoggerService
	jwtKey []byte
}

func NewUserService(repo ports.UserRepository, logger ports.LoggerService, jwtKey []byte) *userService {
	service := userService{
		repo:   repo,
		logger: logger,
		jwtKey: jwtKey,
	}
	return &service
}

func (svc userService) LoginUser(email, password string) (string, error) {
	user, err := svc.GetUserByEmail(email)

	if err != nil {
		svc.logger.Error(fmt.Sprintf("Failed to get user: %v", err))
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))

	if err != nil {
		svc.logger.Error(fmt.Sprintf("Password comparison error: %v", err))
		return "", fmt.Errorf("password comparison error: %v", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.UserId,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString(svc.jwtKey)
	if err != nil {
		svc.logger.Error(err.Error())
		return "", err
	}

	return tokenString, nil
}

func (svc userService) CreateUser(user domain.User) (*domain.User, error) {
	dbUser, _ := svc.GetUserByEmail(user.Email)

	if dbUser != nil {
		svc.logger.Error("create user : user with email exists")
		return nil, errors.New("user with email exists")
	}
	user.UserId = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		svc.logger.Error(fmt.Sprintf("create user : failed to hash password: %v", err))
		return nil, fmt.Errorf("create user: failed to hash password: %v", err)
	}
	user.PasswordHash = string(bytes)
	return svc.repo.CreateUser(user)
}

func (svc userService) GetUserById(userId string) (*domain.User, error) {
	user, err := svc.repo.GetUserById(userId)
	if err != nil {
		svc.logger.Error(fmt.Sprintf("get user by id : failed to get user by id: %v", err))
		return nil, fmt.Errorf("get user by id : failed to get user by id: %v", err)
	}
	return user, nil
}

func (svc userService) GetUserByEmail(email string) (*domain.User, error) {
	user, err := svc.repo.GetUserByEmail(email)
	if err != nil {
		svc.logger.Error(fmt.Sprintf("get user by email : failed to get user by email: %v", err))
		return nil, fmt.Errorf("get user by email : failed to get user by email: %v", err)
	}
	return user, nil
}

func (svc userService) GetUsers() ([]*domain.User, error) {
	users, err := svc.repo.GetUsers()
	if err != nil {
		svc.logger.Error(fmt.Sprintf("get users: failed to get users: %v", err))
		return nil, fmt.Errorf("get users: failed to get users: %v", err)
	}
	return users, nil
}

func (svc userService) GetUsersWithRole(roleName string) ([]*domain.User, error) {
	users, err := svc.repo.GetUsersWithRole(roleName)
	if err != nil {
		svc.logger.Error(fmt.Sprintf("get users with role: failed to get users with role: %v", err))
		return nil, fmt.Errorf("get users with role: failed to get users with role: %v", err)
	}
	return users, nil
}

func (svc userService) UpdateUser(user domain.User) (*domain.User, error) {
	user.UpdatedAt = time.Now()
	dbUser, err := svc.repo.UpdateUser(user)
	if err != nil {
		svc.logger.Error(fmt.Sprintf("update user: failed to update user: %v", err))
		return nil, fmt.Errorf("update user: failed to update user: %v", err)
	}
	return dbUser, nil
}

func (svc userService) DeleteUser(userId string) error {
	err := svc.repo.DeleteUser(userId)
	if err != nil {
		svc.logger.Error(fmt.Sprintf("delete user with id user: failed to delete user with id user: %v", err))
		return fmt.Errorf("delete user with id user: failed to delete user with id user: %v", err)
	}
	return nil
}
