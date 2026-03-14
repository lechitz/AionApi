// Package postgres provides a GORM-based implementation of the DB interface.
package postgres

import (
	"context"

	"github.com/lechitz/AionApi/internal/platform/ports/output/db"
	"gorm.io/gorm"
)

// gormDB wraps *gorm.DB and implements the db.DB interface.
// This adapter allows the application to remain database-agnostic while using GORM internally.
type gormDB struct {
	db *gorm.DB
}

// NewDBAdapter wraps a *gorm.DB instance and returns it as the db.DB interface.
// This follows the Adapter pattern, isolating GORM details from the application layer.
func NewDBAdapter(gormConn *gorm.DB) db.DB {
	if gormConn == nil {
		panic("gorm.DB cannot be nil - ensure database connection is established")
	}
	return &gormDB{db: gormConn}
}

// WithContext returns a new DB instance bound to the given context.
func (g *gormDB) WithContext(ctx context.Context) db.DB {
	return &gormDB{db: g.db.WithContext(ctx)}
}

// Create inserts a new record into the database.
func (g *gormDB) Create(value interface{}) db.DB {
	return &gormDB{db: g.db.Create(value)}
}

// Save updates all fields of an existing record or creates it if it doesn't exist.
func (g *gormDB) Save(value interface{}) db.DB {
	return &gormDB{db: g.db.Save(value)}
}

// Updates specific fields of a record.
func (g *gormDB) Updates(updates interface{}) db.DB {
	return &gormDB{db: g.db.Updates(updates)}
}

// Update updates a single column.
func (g *gormDB) Update(column string, value interface{}) db.DB {
	return &gormDB{db: g.db.Update(column, value)}
}

// Delete performs a soft delete (if the model has a DeletedAt field) or hard delete.
func (g *gormDB) Delete(value interface{}) db.DB {
	return &gormDB{db: g.db.Delete(value)}
}

// First finds the first record ordered by primary key.
func (g *gormDB) First(dest interface{}, conds ...interface{}) db.DB {
	return &gormDB{db: g.db.First(dest, conds...)}
}

// Find retrieves all records that match the given conditions.
func (g *gormDB) Find(dest interface{}, conds ...interface{}) db.DB {
	return &gormDB{db: g.db.Find(dest, conds...)}
}

// Model specifies the model to use for the query.
func (g *gormDB) Model(value interface{}) db.DB {
	return &gormDB{db: g.db.Model(value)}
}

// Select specifies fields to retrieve from the database.
func (g *gormDB) Select(query interface{}, args ...interface{}) db.DB {
	return &gormDB{db: g.db.Select(query, args...)}
}

// Where adds a WHERE clause to the query.
func (g *gormDB) Where(query interface{}, args ...interface{}) db.DB {
	return &gormDB{db: g.db.Where(query, args...)}
}

// Order specifies the order when retrieving records.
func (g *gormDB) Order(value interface{}) db.DB {
	return &gormDB{db: g.db.Order(value)}
}

// Limit specifies the maximum number of records to retrieve.
func (g *gormDB) Limit(limit int) db.DB {
	return &gormDB{db: g.db.Limit(limit)}
}

// Offset specifies the number of records to skip before starting to return records.
func (g *gormDB) Offset(offset int) db.DB {
	return &gormDB{db: g.db.Offset(offset)}
}

// Count counts the number of records that match the query.
func (g *gormDB) Count(count *int64) db.DB {
	return &gormDB{db: g.db.Count(count)}
}

// Scan scans the result into the dest variable.
func (g *gormDB) Scan(dest interface{}) db.DB {
	return &gormDB{db: g.db.Scan(dest)}
}

// Exec executes a raw SQL query (without returning rows).
func (g *gormDB) Exec(sql string, values ...interface{}) db.DB {
	return &gormDB{db: g.db.Exec(sql, values...)}
}

// Raw executes a raw SQL query and returns rows (for SELECT statements).
// This allows complex queries like full-text search, JOINs, and GROUP BY.
func (g *gormDB) Raw(sql string, values ...interface{}) db.DB {
	return &gormDB{db: g.db.Raw(sql, values...)}
}

// Error returns the error from the last operation (if any).
func (g *gormDB) Error() error {
	return g.db.Error
}

// RowsAffected returns the number of rows affected by the last operation.
func (g *gormDB) RowsAffected() int64 {
	return g.db.RowsAffected
}

// Begin starts a new transaction.
func (g *gormDB) Begin() db.DB {
	return &gormDB{db: g.db.Begin()}
}

// Commit commits the current transaction.
func (g *gormDB) Commit() db.DB {
	return &gormDB{db: g.db.Commit()}
}

// Rollback rolls back the current transaction.
func (g *gormDB) Rollback() db.DB {
	return &gormDB{db: g.db.Rollback()}
}

// Transaction executes a function within a transaction.
func (g *gormDB) Transaction(fc func(tx db.DB) error) error {
	return g.db.Transaction(func(tx *gorm.DB) error {
		return fc(&gormDB{db: tx})
	})
}
