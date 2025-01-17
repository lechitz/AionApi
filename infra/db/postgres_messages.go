package db

const (
	msgFormatConString         = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable"
	msgTryingStartsPostgresDB  = "trying starts postgres db"
	errorToStartsThePostgresDB = "error to starts the postgres db"
	msgToRetrieveSQLFromGorm   = "failed to retrieve SQL DB from Gorm"
	errorToCloseThePostgresDB  = "error to close the postgres db"
	msgDBConnection            = "db connection"
	msgTryingToConnect         = "trying to connect"
)
