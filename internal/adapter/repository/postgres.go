package repository

import (
	"database/sql"
	"fmt"

	"github.com/AntonyIS/usafi-hub-user-service/config"
	"github.com/AntonyIS/usafi-hub-user-service/internal/core/domain"
	_ "github.com/lib/pq"
)

type postgresClient struct {
	db *sql.DB
	// logger              ports.LoggerService
	usersTablename      string
	rolesTablename      string
	rolesUsersTablename string
	tablenames          []string
}

func NewBasePostgresClient(config config.Config) (*postgresClient, error) {
	dbname := config.POSTGRES_DB
	tablenames := []string{config.USER_ROLE_TABLE, config.ROLE_TABLE, config.USER_TABLE, "roles"}
	user := config.POSTGRES_USER
	password := config.POSTGRES_PASSWORD
	port := config.POSTGRES_PORT
	host := config.POSTGRES_HOST

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbname, password)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &postgresClient{
		db:                  db,
		usersTablename:      config.USER_TABLE,
		rolesTablename:      config.ROLE_TABLE,
		rolesUsersTablename: config.USER_ROLE_TABLE,
		tablenames:          tablenames,
	}, nil
}

func NewUserPostgresClient(config config.Config) (*postgresClient, error) {
	dbname := config.POSTGRES_DB
	tablename := config.USER_TABLE
	user := config.POSTGRES_USER
	password := config.POSTGRES_PASSWORD
	port := config.POSTGRES_PORT
	host := config.POSTGRES_HOST

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbname, password)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	queryString := fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS %s (
            user_id VARCHAR(255) PRIMARY KEY UNIQUE,
            username VARCHAR(255) NOT NULL,
            password_hash VARCHAR(255) NOT NULL,
            email VARCHAR(255) UNIQUE NOT NULL,
            fullname VARCHAR(255) NOT NULL,
            phone_number VARCHAR(255),
            avatar VARCHAR(255),
            address VARCHAR(255),
            created_at TIMESTAMP,
            updated_at TIMESTAMP
        )
    `, tablename)

	_, err = db.Exec(queryString)
	if err != nil {
		// logger.Error(fmt.Sprintf("Failed to create user table: %v", err))
		return nil, err
	}
	// logger.Info("Connected to the database successfully")
	return &postgresClient{
		db:                  db,
		usersTablename:      config.USER_TABLE,
		rolesTablename:      config.ROLE_TABLE,
		rolesUsersTablename: config.USER_ROLE_TABLE,
		tablenames:          []string{},
		// logger:              logger,
	}, nil
}

func NewRolePostgresClient(config config.Config) (*postgresClient, error) {
	dbname := config.POSTGRES_DB
	tablename := config.ROLE_TABLE
	user := config.POSTGRES_USER
	password := config.POSTGRES_PASSWORD
	port := config.POSTGRES_PORT
	host := config.POSTGRES_HOST

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbname, password)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	roleQueryString := fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS %s (
            role_id VARCHAR(255) PRIMARY KEY UNIQUE,
            name VARCHAR(255) UNIQUE NOT NULL,
            description VARCHAR(255)
        )
    `, tablename)

	_, err = db.Exec(roleQueryString)
	if err != nil {
		// logger.Error(fmt.Sprintf("Failed to create role table: %v", err))
		return nil, err
	}
	// logger.Info("Connected to the database successfully")
	return &postgresClient{
		db:                  db,
		usersTablename:      config.USER_TABLE,
		rolesTablename:      config.ROLE_TABLE,
		rolesUsersTablename: config.USER_ROLE_TABLE,
		tablenames:          []string{},
		// logger:              logger,
	}, nil
}

func NewUserRolePostgresClient(config config.Config) (*postgresClient, error) {
	dbname := config.POSTGRES_DB
	tablename := config.USER_ROLE_TABLE
	user := config.POSTGRES_USER
	password := config.POSTGRES_PASSWORD
	port := config.POSTGRES_PORT
	host := config.POSTGRES_HOST

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbname, password)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	userRoleQueryString := fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS %s (
            user_id VARCHAR(255),
            role_id VARCHAR(255),
            PRIMARY KEY (user_id, role_id),
            CONSTRAINT fk_user_roles_user_id FOREIGN KEY (user_id) REFERENCES %s(user_id),
            CONSTRAINT fk_user_roles_role_id FOREIGN KEY (role_id) REFERENCES %s(role_id)
        )
    `, config.USER_ROLE_TABLE, config.USER_TABLE, config.ROLE_TABLE)

	_, err = db.Exec(userRoleQueryString)
	if err != nil {
		return nil, err
	}

	unique_user_role_name := fmt.Sprintf(`unique_%s`, config.USER_ROLE_TABLE)
	dropQueryString := fmt.Sprintf(`
    	ALTER TABLE %s
    	DROP CONSTRAINT IF EXISTS %s;
	`, tablename, unique_user_role_name)
	_, err = db.Exec(dropQueryString)
	if err != nil {
		return nil, err
	}

	alterQueryString := fmt.Sprintf(`
        ALTER TABLE %s
        ADD CONSTRAINT %s UNIQUE (user_id, role_id);
    `, tablename, unique_user_role_name)

	_, err = db.Exec(alterQueryString)
	if err != nil {
		return nil, err
	}
	return &postgresClient{
		db:                  db,
		usersTablename:      config.USER_TABLE,
		rolesTablename:      config.ROLE_TABLE,
		rolesUsersTablename: config.USER_ROLE_TABLE,
		tablenames:          []string{},
	}, nil
}

func (svc postgresClient) CreateUser(user domain.User) (*domain.User, error) {

	query := fmt.Sprintf(`
        INSERT INTO %s (user_id, username, password_hash, email, fullname, phone_number, avatar, address, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    `, svc.usersTablename)
	_, err := svc.db.Exec(query,
		user.UserId,
		user.Username,
		user.PasswordHash,
		user.Email,
		user.FullName,
		user.PhoneNumber,
		user.Avatar,
		user.Address,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return svc.GetUserById(user.UserId)
}

func (svc postgresClient) GetUserById(userId string) (*domain.User, error) {
	query := fmt.Sprintf(`
        SELECT user_id, username, password_hash, email, fullname, phone_number, avatar, address, created_at, updated_at
        FROM %s
        WHERE user_id = $1
    `, svc.usersTablename)
	row := svc.db.QueryRow(query, userId)
	user := &domain.User{}
	err := row.Scan(&user.UserId, &user.Username, &user.PasswordHash, &user.Email, &user.FullName, &user.PhoneNumber, &user.Avatar, &user.Address, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (svc postgresClient) GetUserByEmail(email string) (*domain.User, error) {
	query := fmt.Sprintf(`
        SELECT user_id, username, password_hash, email, fullname, phone_number, avatar, address, created_at, updated_at
        FROM %s
        WHERE email = $1
    `, svc.usersTablename)
	row := svc.db.QueryRow(query, email)
	user := &domain.User{}
	err := row.Scan(&user.UserId, &user.Username, &user.PasswordHash, &user.Email, &user.FullName, &user.PhoneNumber, &user.Avatar, &user.Address, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (svc postgresClient) GetUsers() ([]*domain.User, error) {
	query := fmt.Sprintf(`
        SELECT user_id, username, password_hash, email, fullname, phone_number, avatar, address, created_at, updated_at
        FROM %s
    `, svc.usersTablename)
	rows, err := svc.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*domain.User{}
	for rows.Next() {
		user := &domain.User{}
		err := rows.Scan(&user.UserId, &user.Username, &user.PasswordHash, &user.Email, &user.FullName, &user.PhoneNumber, &user.Avatar, &user.Address, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (svc postgresClient) GetUsersWithRole(roleName string) ([]*domain.User, error) {
	query := fmt.Sprintf(`
        SELECT u.user_id, u.username, u.password_hash, u.email, u.fullname, u.phone_number, u.avatar, u.address, u.created_at, u.updated_at
        FROM %s u
        JOIN %s ur ON u.user_id = ur.user_id
        JOIN %s r ON ur.role_id = r.role_id
        WHERE r.name = $1
    `, svc.usersTablename, svc.rolesUsersTablename, svc.rolesTablename)
	rows, err := svc.db.Query(query, roleName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*domain.User{}
	for rows.Next() {
		user := &domain.User{}
		err := rows.Scan(&user.UserId, &user.Username, &user.PasswordHash, &user.Email, &user.FullName, &user.PhoneNumber, &user.Avatar, &user.Address, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (svc postgresClient) UpdateUser(user domain.User) (*domain.User, error) {
	query := fmt.Sprintf(`
        UPDATE %s
        SET username=$2, password_hash=$3, email=$4, fullname=$5, phone_number=$6, avatar=$7, address=$8, updated_at=$9
        WHERE user_id=$1
    `, svc.usersTablename)
	_, err := svc.db.Exec(query, user.UserId, user.Username, user.PasswordHash, user.Email, user.FullName, user.PhoneNumber, user.Avatar, user.Address, user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return svc.GetUserById(user.UserId)
}

func (svc postgresClient) DeleteUser(userId string) error {
	query := fmt.Sprintf(`
        DELETE FROM %s
        WHERE user_id=$1
    `, svc.usersTablename)
	_, err := svc.db.Exec(query, userId)
	if err != nil {
		return err
	}
	return nil
}

func (svc postgresClient) CreateRole(role domain.Role) (*domain.Role, error) {
	roles, err := svc.GetRoles()
	if err != nil {
		return nil, err
	}
	for _, roleItem := range roles {
		if roleItem.Name == role.Name && roleItem.Description == role.Description {
			return roleItem, nil
		}
	}
	query := fmt.Sprintf(`
        INSERT INTO %s (role_id, name, description)
        VALUES ($1, $2, $3)
	`, svc.rolesTablename)

	_, err = svc.db.Exec(query, role.RoleId, role.Name, role.Description)
	if err != nil {
		return nil, err
	}
	return svc.GetRoleById(role.RoleId)
}

func (svc postgresClient) GetRoleById(roleId string) (*domain.Role, error) {
	query := fmt.Sprintf(`
        SELECT role_id, name, description
        FROM %s
        WHERE role_id = $1
	`, svc.rolesTablename)
	row := svc.db.QueryRow(query, roleId)
	role := &domain.Role{}
	err := row.Scan(&role.RoleId, &role.Name, &role.Description)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (svc postgresClient) GetRoles() ([]*domain.Role, error) {
	query := fmt.Sprintf(`
        SELECT role_id, name, description
        FROM %s
    `, svc.rolesTablename)
	rows, err := svc.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := []*domain.Role{}
	for rows.Next() {
		role := &domain.Role{}
		err := rows.Scan(&role.RoleId, &role.Name, &role.Description)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return roles, nil
}

func (svc postgresClient) UpdateRole(role domain.Role) error {
	query := fmt.Sprintf(`
        UPDATE %s
        SET name=$2, description=$3
        WHERE role_id=$1
	`, svc.rolesTablename)
	_, err := svc.db.Exec(query, role.RoleId, role.Name, role.Description)
	if err != nil {
		return err
	}
	return nil
}

func (svc postgresClient) DeleteRole(roleId string) error {
	query := fmt.Sprintf(`
        DELETE FROM %s
        WHERE role_id=$1
	`, svc.rolesTablename)
	_, err := svc.db.Exec(query, roleId)
	if err != nil {
		return err
	}
	return nil
}

func (svc postgresClient) AddUserRole(userRole domain.UserRole) error {
	query := fmt.Sprintf(`
        INSERT INTO %s (user_id, role_id)
        VALUES ($1, $2)
   	`, svc.rolesUsersTablename)

	_, err := svc.db.Exec(query, userRole.UserId, userRole.RoleId)
	if err != nil {
		return err
	}

	return nil
}

func (svc postgresClient) RemoveUserRole(userRole domain.UserRole) error {
	query := fmt.Sprintf(`
        DELETE FROM %s
        WHERE user_id=$1 AND role_id=$2
   	`, svc.rolesUsersTablename)
	_, err := svc.db.Exec(query, userRole.UserId, userRole.RoleId)
	if err != nil {
		return err
	}
	return nil
}

func (svc postgresClient) DropTables() error {
	for _, tablename := range svc.tablenames {
		_, err := svc.db.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, tablename))
		if err != nil {
			return err
		}
	}
	return nil

}
