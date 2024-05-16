package cmd

import (
	"github.com/AntonyIS/usafi-hub-user-service/config"
	"github.com/AntonyIS/usafi-hub-user-service/internal/adapter/app"
	"github.com/AntonyIS/usafi-hub-user-service/internal/adapter/logger"
	"github.com/AntonyIS/usafi-hub-user-service/internal/adapter/repository"
	"github.com/AntonyIS/usafi-hub-user-service/internal/core/services"
)

func RunService() {
	logger, err := logger.NewDefaultLogger()
	if err != nil {
		panic(err)
	}

	config, err := config.NewConfig(logger)
	if err != nil {
		panic(err)
	}

	roleRepo, _ := repository.NewRolePostgresClient(*config, logger)
	userRepo, _ := repository.NewUserPostgresClient(*config, logger)
	userRoleRepo, _ := repository.NewUserRolePostgresClient(*config, logger)

	userService := services.NewUserService(userRepo)
	roleService := services.NewRoleService(roleRepo)
	userRoleService := services.NewUserRoleService(userRoleRepo)

	app.InitGinRoutes(userService, roleService, userRoleService, *config)
}
