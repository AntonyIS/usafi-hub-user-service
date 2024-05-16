package domain

import (
	"time"
)

type User struct {
	UserId       string    `json:"user_id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"password_hash"`
	Email        string    `json:"email"`
	FullName     string    `json:"fullname"`
	PhoneNumber  string    `json:"phone_number"`
	Avatar       string    `json:"avatar"`
	Address      string    `json:"address"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Role struct {
	RoleId      string `json:"role_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UserRole struct {
	UserId string `json:"user_id"`
	RoleId string `json:"role_id"`
}
