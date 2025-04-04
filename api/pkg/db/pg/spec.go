package pg

import (
	"context"
	"database/sql"
)

// Executor can perform SQL queries.
type Executor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// ContextExecutor can perform SQL queries with context
type ContextExecutor interface {
	Executor

	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// Beginner allows creation of context aware transactions with options.
type Beginner interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (Transactor, error)
}

// BeginnerExecutor can context-aware perform SQL queries and
// create context-aware transactions with options
type BeginnerExecutor interface {
	Beginner
	ContextExecutor

	Close() error
	InstanceInfo() InstanceInfo
}

// Transactor is an interface for an sql.Tx
type Transactor interface {
	Commit() error
	Rollback() error
	ContextExecutor
}

// InstanceInfo holds info about the DB
type InstanceInfo struct {
	dbName string
}
