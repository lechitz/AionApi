package constants

const (
	MsgFormatConString         = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC"
	MsgTryingStartsPostgresDB  = "trying starts postgres db"
	ErrorToStartsThePostgresDB = "error to starts the postgres db"
	MsgToRetrieveSQLFromGorm   = "failed to retrieve SQL DB from Gorm"
	ErrorToCloseThePostgresDB  = "error to close the postgres db"
	MsgDBConnection            = "db connection"
	MsgTryingToConnect         = "trying to connect"
	Port                       = "port"
	DBName                     = "dbname"
	Host                       = "host"
	Error                      = "error"
)
