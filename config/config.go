package config

import (
	"os"

	"github.com/AntonyIS/usafi-hub-user-service/internal/core/ports"
	"github.com/joho/godotenv"
)

type Config struct {
	ENV               string
	SECRET_KEY        string
	SERVER_PORT       string
	POSTGRES_DB       string
	POSTGRES_HOST     string
	POSTGRES_PORT     string
	POSTGRES_USER     string
	POSTGRES_PASSWORD string
	USER_TABLE        string
	ROLE_TABLE        string
	USER_ROLE_TABLE   string
	DEBUG             bool
	TEST              bool
}

func NewConfig(logger ports.LoggerService) (*Config, error) {
	ENV := os.Getenv("ENV")

	switch ENV {
	case "development":
		err := godotenv.Load(".env")
		if err != nil {
			logger.Error(err.Error())
			return nil, err
		}

	}

	var (
		SECRET_KEY        = os.Getenv("SECRET_KEY")
		SERVER_PORT       = "5000"
		POSTGRES_DB       = "usafihub"
		POSTGRES_HOST     = "postgres"
		POSTGRES_PORT     = "5432"
		POSTGRES_USER     = "postgres"
		POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
		USER_TABLE        = ""
		ROLE_TABLE        = ""
		USER_ROLE_TABLE   = ""
		DEBUG             = false
		TEST              = false
	)

	switch ENV {
	case "production":
		TEST = false
		DEBUG = false

	case "production_test":
		TEST = true
		DEBUG = true
		USER_TABLE = "Prod_Test_Users"
		ROLE_TABLE = "Prod_Test_Roles"
		USER_ROLE_TABLE = "Prod_Test_UserRoles"

	case "development":
		TEST = true
		DEBUG = true
		POSTGRES_HOST = "localhost"
		USER_TABLE = "Dev_Users"
		ROLE_TABLE = "Dev_Roles"
		USER_ROLE_TABLE = "Dev_UserRoles"

	case "development_test":
		TEST = true
		DEBUG = true
		SECRET_KEY = "testsecret"
		POSTGRES_PASSWORD = "pass1234"
		POSTGRES_HOST = "localhost"
		USER_TABLE = "Test_Users"
		ROLE_TABLE = "Test_Roles"
		USER_ROLE_TABLE = "Test_UserRoles"

	case "docker":
		TEST = true
		DEBUG = true
		USER_TABLE = "Docker_Users"
		ROLE_TABLE = "Docker_Roles"
		USER_ROLE_TABLE = "Docker_UserRoles"

	case "docker_test":
		TEST = true
		DEBUG = true
		USER_TABLE = "Docker_Test_Users"
		ROLE_TABLE = "Docker_Test_Roles"
		USER_ROLE_TABLE = "Docker_Test_UserRoles"
	}

	config := Config{
		ENV:               ENV,
		SECRET_KEY:        SECRET_KEY,
		SERVER_PORT:       SERVER_PORT,
		POSTGRES_DB:       POSTGRES_DB,
		POSTGRES_HOST:     POSTGRES_HOST,
		POSTGRES_PORT:     POSTGRES_PORT,
		POSTGRES_USER:     POSTGRES_USER,
		POSTGRES_PASSWORD: POSTGRES_PASSWORD,
		USER_TABLE:        USER_TABLE,
		ROLE_TABLE:        ROLE_TABLE,
		USER_ROLE_TABLE:   USER_ROLE_TABLE,
		DEBUG:             DEBUG,
		TEST:              TEST,
	}

	return &config, nil
}
