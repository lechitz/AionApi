package constants

const (
	MsgFormatConString          = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC"
	MsgTryingStartsDB           = "trying to establish Database connection"
	ErrorToStartDB              = "error to starts the postgres db"
	MsgToRetrieveSQLFromGorm    = "failed to retrieve SQL DB from Gorm"
	ErrorToCloseDB              = "error to close the postgres db"
	MsgDBConnection             = "initializing Database connection"
	FailedToPingDB              = "failed to ping DB"
	ErrDBConnectionAttempt      = "connection attempt failed"
	MsgPostgresConnectionClosed = "Database connection closed successfully"
	Port                        = "port"
	DBName                      = "dbname"
	Host                        = "host"
	Error                       = "error"
	Try                         = "try"
)
