package testutil

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"omg/api/pkg/db/pg"
	"omg/api/pkg/env"
)

// appDB caches a pg connection for reuse
var appDB *sql.DB

// WithTxDB provides callback with a `pg.BeginnerExecutor` for running pg related tests
// where the `pg.BeginnerExecutor` is actually powered by a pg transaction
// and will be rolled back (so no data is actually written into pg)
func WithTxDB(t *testing.T, callback func(pg.BeginnerExecutor)) {
	if appDB == nil {
		var err error
		appDB, err = sql.Open("postgres", env.GetAndValidateF("DB_URL"))
		appDB.SetMaxOpenConns(50)
		appDB.SetConnMaxLifetime(30 * time.Minute)
		require.NoError(t, err)
	}

	tx, err := appDB.BeginTx(context.Background(), nil)
	require.NoError(t, err)

	// explicitly hardcoded to rollback via `defer` because
	// in case `callback` does a t.Fatal we must still run tx.Rollback
	defer tx.Rollback()

	callback(&txDB{Tx: tx})
}

type txStatus int

// Constants of rollback status
const (
	Started txStatus = iota
	Committed
	RolledBack
)

// txDB wraps `sql.Tx` in our tests to pretend to be `sql.DB`
type txDB struct {
	*sql.Tx

	Status txStatus
}

// Commit implements Transactor
func (tx *txDB) Commit() error {
	if tx.Status == Committed || tx.Status == RolledBack {
		return sql.ErrTxDone
	}

	tx.Status = Committed

	return nil
}

// Rollback implements Transactor
// After a call to Commit or Rollback, all operations on the
// transaction fail with ErrTxDone.
func (tx *txDB) Rollback() error {
	if tx.Status == Committed || tx.Status == RolledBack {
		return sql.ErrTxDone
	}

	tx.Status = RolledBack

	return nil
}

// Close is a test mock for pg.Close() which closes the pg connection
func (*txDB) Close() error {
	return nil
}

// BeginTx implements Beginner
func (tx *txDB) BeginTx(context.Context, *sql.TxOptions) (pg.Transactor, error) {
	tx.Status = Started

	return FakeTx{tx.Tx}, nil
}

func (*txDB) InstanceInfo() pg.InstanceInfo {
	return pg.InstanceInfo{}
}

// FakeTx is a test mock for the Transactor interface
type FakeTx struct {
	*sql.Tx
}

// Commit prevents a FakeTx from committing the parent tx
func (FakeTx) Commit() error {
	return nil
}

// Rollback prevents a FakeTx from rolling back the parent tx
func (FakeTx) Rollback() error {
	return nil
}
