package pg

import (
	"context"
	"database/sql"
)

// gobaseDB wraps the *sql.DB provided
type gobaseDB struct {
	*sql.DB
	info InstanceInfo
}

// BeginTx begins a transaction with the database in receiver and returns a Transactor
func (i gobaseDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (Transactor, error) {
	return i.DB.BeginTx(ctx, opts)
}

// ExecContext wraps the base connector
func (i gobaseDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return i.DB.ExecContext(ctx, query, args...)
}

// QueryContext wraps the base connector
func (i gobaseDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return i.DB.QueryContext(ctx, query, args...)
}

// QueryRowContext wraps the base connector
func (i gobaseDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return i.DB.QueryRowContext(ctx, query, args...)
}

// InstanceInfo returns info about the DB
func (i gobaseDB) InstanceInfo() InstanceInfo {
	return i.info
}

// gobaseTx wraps the Transactor provided
type gobaseTx struct {
	Transactor
	info InstanceInfo
	ctx  context.Context
}

// ExecContext wraps the base connector
func (i gobaseTx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return i.Transactor.ExecContext(ctx, query, args...)
}

// QueryContext wraps the base connector
func (i gobaseTx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return i.Transactor.QueryContext(ctx, query, args...)
}

// QueryRowContext wraps the base connector
func (i gobaseTx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return i.Transactor.QueryRowContext(ctx, query, args...)
}

// Commit commits the transaction
func (i gobaseTx) Commit() error {
	return i.Transactor.Commit()
}

// Rollback aborts the transaction
func (i gobaseTx) Rollback() error {
	return i.Transactor.Rollback()
}
