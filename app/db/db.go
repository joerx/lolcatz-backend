package db

import (
	"context"
	"database/sql"
)

// DB abstracts database specific clients
type DB interface {
	Close() error
	Ping(ctx context.Context) error
	Exec(ctx context.Context, q string, args ...interface{}) (sql.Result, error)
	Query(ctx context.Context, q string, args ...interface{}) (*sql.Rows, error)
	Prepare(ctx context.Context, q string) (*sql.Stmt, error)
}
