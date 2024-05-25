package app

import (
	"fmt"
	"log"

	"github.com/AntonyIS/usafi-hub-user-service/config"
	"github.com/AntonyIS/usafi-hub-user-service/internal/core/ports"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitGinRoutes(userService ports.UserService, roleService ports.RoleService, userRoleService ports.UserRoleService, config config.Config, logger ports.LoggerService) {
	gin.SetMode(gin.DebugMode)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	handler := NewGinHandler(
		userService,
		roleService,
		userRoleService,
	)

	homeRoutes := router.Group("/")
	userRoutes := router.Group("/users/v1")
	roleRoutes := router.Group("/roles/v1")
	userRoleRoutes := router.Group("/user_roles/v1")
	authRoutes := router.Group("/auth/v1")

	middleware := NewMiddleware(userService, logger, config.SECRET_KEY)

	homeRoutes.Use(middleware.AuthorizeToken)
	userRoutes.Use(middleware.AuthorizeToken)
	roleRoutes.Use(middleware.AuthorizeToken)
	userRoleRoutes.Use(middleware.AuthorizeToken)

	{
		homeRoutes.GET("/", handler.Home)
		homeRoutes.GET("/health-check", handler.Healthcheck)
	}

	{
		userRoutes.POST("/", handler.CreateUser)
		userRoutes.POST("/get", handler.GetUserByEmail)
		userRoutes.GET("/:user_id", handler.GetUserById)
		userRoutes.GET("/", handler.GetUsers)
		userRoutes.GET("/roles/:role_name", handler.GetUsersWithRole)
		userRoutes.PUT(":user_id", handler.UpdateUser)
		userRoutes.DELETE("/:user_id", handler.DeleteUser)
	}
	{
		roleRoutes.POST("/", handler.CreateRole)
		roleRoutes.GET("/:role_id", handler.GetRoleById)
		roleRoutes.GET("/", handler.GetRoles)
		roleRoutes.PUT(":role_id", handler.UpdateRole)
		roleRoutes.DELETE("/:role_id", handler.DeleteRole)
	}

	{
		userRoleRoutes.POST("/", handler.AddUserRole)
		userRoleRoutes.GET("/:user_role_id", handler.RemoveUserRole)
	}

	{
		authRoutes.POST("/signup", handler.SignupUser)
		authRoutes.POST("/login", handler.LoginUser)
		authRoutes.POST("/forgot-password", handler.ForgotPassword)
	}
	log.Printf("Server running on port 0.0.0.0:%s", config.SERVER_PORT)
	router.Run(fmt.Sprintf(":%s", config.SERVER_PORT))
}
