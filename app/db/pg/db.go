package pg

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/joerx/lolcatz-backend/db"
)

type pgdb struct {
	db *sql.DB
}

// Config holds configuration values for postgres database connections
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

// NewWithSchema creates a new database client and loads the initial database schema
func NewWithSchema(ctx context.Context, cf Config) (db.DB, error) {
	pgdb, err := NewClient(ctx, cf)
	if err != nil {
		return nil, err
	}

	if err := loadSchema(pgdb); err != nil {
		return nil, err
	}

	log.Println("Database connection initialized")
	return pgdb, nil
}

// NewClient creates and connects a new postgres database client
func NewClient(ctx context.Context, cf Config) (db.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cf.Host, cf.Port, cf.User, cf.Password, cf.Name,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	return pgdb{db}, nil
}

func (p pgdb) Close() error {
	return p.db.Close()
}

func (p pgdb) Ping(ctx context.Context) error {
	return p.db.PingContext(ctx)
}

func (p pgdb) Exec(ctx context.Context, q string, args ...interface{}) (sql.Result, error) {
	return p.db.ExecContext(ctx, q, args...)
}

func (p pgdb) Query(ctx context.Context, q string, args ...interface{}) (*sql.Rows, error) {
	return p.db.QueryContext(ctx, q, args...)
}
