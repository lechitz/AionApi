package contextkeys

import "time"

// ContextKey constants

const (
	Users       = "users"
	User        = "user"
	UserID      = "id"
	Name        = "name"
	Username    = "username"
	Email       = "email"
	Password    = "password"
	CreatedAt   = "created_at"
	UpdatedAt   = "updated_at"
	DeletedAt   = "deleted_at"
	AuthToken   = "auth_token"
	Error       = "error"
	Domain      = "localhost"
	Token       = "token"
	Path        = "/"
	Context     = "context"
	ContextPath = "contextPath"
	Host        = "host"
	Port        = "port"
	DBName      = "dbname"
	Config      = "config"

	Setting = "setting"
)

// Service constants

const (
	ExpTimeToken = 1 * time.Hour
	Exp          = "exp"
)
