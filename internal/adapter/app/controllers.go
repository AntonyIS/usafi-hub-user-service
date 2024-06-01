package app

import (
	"net/http"

	"github.com/AntonyIS/usafi-hub-user-service/internal/core/domain"
	"github.com/AntonyIS/usafi-hub-user-service/internal/core/ports"
	"github.com/gin-gonic/gin"
)

type GinHandler interface {
	Home(ctx *gin.Context)
	Healthcheck(ctx *gin.Context)
	CreateUser(ctx *gin.Context)
	GetUsers(ctx *gin.Context)
	GetUsersWithRole(ctx *gin.Context)
	GetUserById(ctx *gin.Context)
	GetUserByEmail(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
	CreateRole(ctx *gin.Context)
	GetRoleById(ctx *gin.Context)
	GetRoles(ctx *gin.Context)
	UpdateRole(ctx *gin.Context)
	DeleteRole(ctx *gin.Context)
	AddUserRole(ctx *gin.Context)
	RemoveUserRole(ctx *gin.Context)
	SignupUser(ctx *gin.Context)
	LoginUser(ctx *gin.Context)
	ForgotPassword(ctx *gin.Context)
	GenerateToken(ctx *gin.Context)
}

type handler struct {
	userService     ports.UserService
	roleService     ports.RoleService
	userRoleService ports.UserRoleService
}

func NewGinHandler(userService ports.UserService, roleService ports.RoleService, userRoleService ports.UserRoleService) GinHandler {
	routerHandler := handler{
		userService:     userService,
		roleService:     roleService,
		userRoleService: userRoleService,
	}
	return routerHandler
}

func (h handler) Home(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"responseMessage": "UsafiHub User Service",
		"responseCode":    http.StatusOK,
	})
}

func (h handler) Healthcheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"responseMessage": "UsafiHub User Service Health Check",
		"responseCode":    http.StatusOK,
	})
}

func (h handler) CreateUser(ctx *gin.Context) {
	var user domain.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"responseMessage": err.Error(),
			"responseCode":    http.StatusBadRequest,
		})
		return
	}

	dbUser, err := h.userService.CreateUser(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"responseMessage": err.Error(),
			"responseCode":    http.StatusInternalServerError,
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"responseMessage": "User created successfully",
		"responseCode":    http.StatusCreated,
		"data":            dbUser,
	})
}

func (h handler) GetUsers(ctx *gin.Context) {
	users, err := h.userService.GetUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"responseMessage": err.Error(),
			"responseCode":    http.StatusInternalServerError,
		})
		return
	}

	if len(users) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"responseMessage": "No users found",
			"responseCode":    http.StatusOK,
			"response":        users,
		})
	} else {
		ctx.JSON(http.StatusOK, users)
	}
}

func (h handler) GetUsersWithRole(ctx *gin.Context) {
	roleName := ctx.Param("role_name")

	users, err := h.userService.GetUsersWithRole(roleName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"responseMessage": err.Error(),
			"responseCode":    http.StatusInternalServerError,
		})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (h handler) GetUserById(ctx *gin.Context) {
	userId := ctx.Param("user_id")
	user, err := h.userService.GetUserById(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"responseMessage": err.Error(),
			"responseCode":    http.StatusInternalServerError,
		})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (h handler) GetUserByEmail(ctx *gin.Context) {
	var user struct {
		Email string
	}
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"responseMessage": err.Error(),
			"responseCode":    http.StatusBadRequest,
		})
		return
	}

	dbUser, err := h.userService.GetUserByEmail(user.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"responseMessage": err.Error(),
			"responseCode":    http.StatusInternalServerError,
		})
		return
	}

	ctx.JSON(http.StatusOK, dbUser)
}

func (h handler) UpdateUser(ctx *gin.Context) {
	var user domain.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"responseMessage": err.Error(),
			"responseCode":    http.StatusBadRequest,
		})
		return
	}

	dbUser, err := h.userService.UpdateUser(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"responseMessage": err.Error(),
			"responseCode":    http.StatusInternalServerError,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"responseMessage": "User updated successfully",
		"responseCode":    http.StatusOK,
		"data":            dbUser,
	})
}

func (h handler) DeleteUser(ctx *gin.Context) {
	userId := ctx.Param("userId")
	err := h.userService.DeleteUser(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"responseMessage": err.Error(),
			"responseCode":    http.StatusInternalServerError,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"responseMessage": "User deleted successfully",
		"responseCode":    http.StatusOK,
	})
}

func (h handler) CreateRole(ctx *gin.Context) {
	var role domain.Role
	if err := ctx.ShouldBindJSON(&role); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"responseMessage": err.Error(),
			"responseCode":    http.StatusBadRequest,
		})
		return
	}

	newRole, err := h.roleService.CreateRole(role)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"responseMessage": err.Error(),
			"responseCode":    http.StatusInternalServerError,
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"responseMessage": "Role created successfully",
		"responseCode":    http.StatusCreated,
		"data":            newRole,
	})
}

func (h handler) GetRoleById(ctx *gin.Context) {
	roleID := ctx.Param("role_id")
	role, err := h.roleService.GetRoleById(roleID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"responseMessage": err.Error(),
			"responseCode":    http.StatusInternalServerError,
		})
		return
	}

	if role == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"responseMessage": "Role not found",
			"responseCode":    http.StatusNotFound,
		})
		return
	}

	ctx.JSON(http.StatusOK, role)
}

func (h handler) GetRoles(ctx *gin.Context) {
	roles, err := h.roleService.GetRoles()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"responseMessage": err.Error(),
			"responseCode":    http.StatusInternalServerError,
		})
		return
	}
	if len(roles) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"responseMessage": "No roles found",
			"responseCode":    http.StatusOK,
			"response":        roles,
		})

	} else {
		ctx.JSON(http.StatusOK, roles)
	}
}

func (h handler) UpdateRole(ctx *gin.Context) {
	roleID := ctx.Param("role_id")
	var role domain.Role
	if err := ctx.ShouldBindJSON(&role); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"responseMessage": err.Error(),
			"responseCode":    http.StatusBadRequest,
		})
		return
	}

	role.RoleId = roleID
	if err := h.roleService.UpdateRole(role); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"responseMessage": err.Error(),
			"responseCode":    http.StatusInternalServerError,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"responseMessage": "Role updated successfully",
		"responseCode":    http.StatusOK,
	})
}

func (h handler) DeleteRole(ctx *gin.Context) {
	roleID := ctx.Param("role_id")
	role, err := h.roleService.GetRoleById(roleID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"responseMessage": err.Error(),
			"responseCode":    http.StatusNotFound,
		})
		return
	}

	if role == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"responseMessage": err.Error(),
			"responseCode":    http.StatusNotFound,
		})
		return
	}
	if err := h.roleService.DeleteRole(roleID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"responseMessage": err.Error(),
			"responseCode":    http.StatusInternalServerError,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"responseMessage": "Role deleted successfully",
		"responseCode":    http.StatusOK,
	})
}

func (h handler) AddUserRole(ctx *gin.Context) {
	var userRole domain.UserRole
	if err := ctx.ShouldBindJSON(&userRole); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"responseMessage": err.Error(),
			"responseCode":    http.StatusBadRequest,
		})
		return
	}

	err := h.userRoleService.AddUserRole(userRole)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"responseMessage": err.Error(),
			"responseCode":    http.StatusInternalServerError,
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"responseMessage": "User role created successfully",
		"responseCode":    http.StatusCreated,
	})
}

func (h handler) RemoveUserRole(ctx *gin.Context) {
	var userRole domain.UserRole
	if err := ctx.ShouldBindJSON(&userRole); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"responseMessage": err.Error(),
			"responseCode":    http.StatusBadRequest,
		})
		return
	}

	err := h.userRoleService.RemoveUserRole(userRole)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"responseMessage": err.Error(),
			"responseCode":    http.StatusInternalServerError,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"responseMessage": "User role deleted successfully",
		"responseCode":    http.StatusOK,
	})
}

func (h handler) SignupUser(ctx *gin.Context) {
	var user domain.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"responseMessage": err.Error(),
			"responseCode":    http.StatusBadRequest,
		})
		return
	}

	dbUser, err := h.userService.CreateUser(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"responseMessage": err.Error(),
			"responseCode":    http.StatusInternalServerError,
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"responseMessage": "User created successfully",
		"responseCode":    http.StatusCreated,
		"data":            dbUser,
	})
}

func (h handler) LoginUser(ctx *gin.Context) {
	type User struct {
		Email        string `json:"email"`
		PasswordHash string `json:"password"`
	}
	var user User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"responseMessage": err.Error(),
			"responseCode":    http.StatusBadRequest,
		})
		return
	}

	token, err := h.userService.LoginUser(user.Email, user.PasswordHash)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"responseMessage": "Invalid email or password",
				"responseCode":    http.StatusUnauthorized,
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"responseMessage": err.Error(),
			"responseCode":    http.StatusInternalServerError,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token": token,
	})
}

func (h handler) ForgotPassword(ctx *gin.Context) {
	var user domain.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"responseMessage": err.Error(),
			"responseCode":    http.StatusBadRequest,
		})
		return
	}

	dbUser, err := h.userService.CreateUser(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"responseMessage": err.Error(),
			"responseCode":    http.StatusInternalServerError,
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"responseMessage": "User created successfully",
		"responseCode":    http.StatusCreated,
		"data":            dbUser,
	})
}

func (h handler) GenerateToken(ctx *gin.Context) {
	type User struct {
		Email        string `json:"email"`
		PasswordHash string `json:"password_hash"`
	}
	var user User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"responseMessage": err.Error(),
			"responseCode":    http.StatusBadRequest,
		})
		return
	}

	token, err := h.userService.LoginUser(user.Email, user.PasswordHash)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"responseMessage": "Invalid email or password",
				"responseCode":    http.StatusUnauthorized,
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"responseMessage": err.Error(),
			"responseCode":    http.StatusInternalServerError,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
