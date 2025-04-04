package pg

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	pkgerrors "github.com/pkg/errors"
)

// NewPool opens a new DB connection pool, pings it and returns the pool
// connection pool definition: http://go-database-sql.org/connection-pool.html
func NewPool(
	url string,
	poolMaxOpenConns,
	poolMaxIdleConns int,
) (BeginnerExecutor, error) {
	fmt.Println("url: ", url)
	log.Println("Initializing Postgres")

	conn, err := sql.Open("postgres", url)
	if err != nil {
		return nil, pkgerrors.WithStack(fmt.Errorf("opening DB failed. err: %w", err))
	}

	conn.SetConnMaxLifetime(29 * time.Minute) // Azure's default is 30 mins.
	conn.SetMaxOpenConns(poolMaxOpenConns)
	conn.SetMaxIdleConns(poolMaxIdleConns)

	log.Println("Pinging DB...")
	if err = conn.Ping(); err != nil {
		return nil, pkgerrors.WithStack(fmt.Errorf("unable to ping DB. err: %w", err))
	}

	log.Println("DB ping successful")

	log.Println("Postgres initialized")

	return &gobaseDB{
		DB: conn,
	}, nil
}
