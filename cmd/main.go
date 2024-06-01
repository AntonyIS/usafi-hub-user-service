package cmd

import (
	"fmt"

	"github.com/AntonyIS/usafi-hub-user-service/config"
	"github.com/AntonyIS/usafi-hub-user-service/internal/adapter/app"
	"github.com/AntonyIS/usafi-hub-user-service/internal/adapter/logger"
	"github.com/AntonyIS/usafi-hub-user-service/internal/adapter/repository"
	"github.com/AntonyIS/usafi-hub-user-service/internal/core/services"
)

func RunService() {
	logger, err := logger.NewDefaultLogger()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to load logger: %v", err))
		panic(err)
	}

	logger.Info("Loaded logger successfully")

	config, err := config.NewConfig(logger)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to load config: %v", err))
		panic(err)
	}
	logger.Info("Loaded configurations successfully...")
	roleRepo, _ := repository.NewRolePostgresClient(*config)
	userRepo, _ := repository.NewUserPostgresClient(*config)
	userRoleRepo, _ := repository.NewUserRolePostgresClient(*config)

	logger.Info("Service repository running successfully...")

	userService := services.NewUserService(userRepo, logger, []byte(config.SECRET_KEY))
	roleService := services.NewRoleService(roleRepo)
	userRoleService := services.NewUserRoleService(userRoleRepo)

	logger.Info("Services running successfully...")
	app.InitGinRoutes(userService, roleService, userRoleService, *config, logger)
}
