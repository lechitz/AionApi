// Package constants contains constants used for database operations.
package constants

// MsgFormatConString is the template for the database connection string.
const MsgFormatConString = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC"

// MsgTryingStartsDB indicates an attempt to establish a database connection.
const MsgTryingStartsDB = "trying to establish Database connection"

// ErrorToStartDB indicates a failure to start the PostgreSQL database.
const ErrorToStartDB = "error to start the postgres db"

// MsgToRetrieveSQLFromGorm indicates a failure to retrieve the SQL DB from Gorm.
const MsgToRetrieveSQLFromGorm = "failed to retrieve SQL DB from Gorm"

// ErrorToCloseDB indicates a failure to close the PostgreSQL database.
const ErrorToCloseDB = "error to close the postgres db"

// MsgDBConnection indicates the initialization of a database connection.
const MsgDBConnection = "initializing Database connection"

// FailedToPingDB indicates a failed ping to the database.
const FailedToPingDB = "failed to ping DB"

// ErrDBConnectionAttempt indicates a failed database connection attempt.
const ErrDBConnectionAttempt = "connection attempt failed"

// MsgPostgresConnectionClosed indicates that the PostgreSQL connection was closed successfully.
const MsgPostgresConnectionClosed = "Database connection closed successfully"

// Port is the key for the database port value.
const Port = "port"

// DBName is the key for the database name value.
const DBName = "dbname"

// Host is the key for the database host value.
const Host = "host"

// Error is the key for error messages in database context.
const Error = "error"

// Try is a generic key for "try" attempts in database context.
const Try = "try"
