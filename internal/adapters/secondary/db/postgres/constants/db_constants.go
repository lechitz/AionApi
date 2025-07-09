// Package constants contains constants used for database operations.
package constants

// MsgFormatConString is the template for the database connection string.
const MsgFormatConString = "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s"

// MsgTryingStartsDB indicates an attempt to establish a database connection.
const MsgTryingStartsDB = "trying to establish Database connection"

// ErrorToStartDB indicates a failure to start the PostgresSQL database.
const ErrorToStartDB = "error to start the postgres db"

// MsgToRetrieveSQLFromGorm indicates a failure to retrieve the SQL DB from Gorm.
const MsgToRetrieveSQLFromGorm = "failed to retrieve SQL DB from Gorm"

// ErrorToCloseDB indicates a failure to close the PostgresSQL database.
const ErrorToCloseDB = "error to close the postgres db"

// MsgDBConnection indicates the initialization of a database connection.
const MsgDBConnection = "initializing Database connection"

// FailedToPingDB indicates a failed ping to the database.
const FailedToPingDB = "failed to ping DB"

// ErrDBConnectionAttempt indicates a failed database connection attempt.
const ErrDBConnectionAttempt = "connection attempt failed"

// MsgPostgresConnectionClosed indicates that the PostgresSQL connection was closed successfully.
const MsgPostgresConnectionClosed = "Database connection closed successfully"
