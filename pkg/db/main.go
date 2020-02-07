package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Config holds the database config
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	UseTLS   bool // ignored for now
}

// Client is a generic database client
type Client struct {
	db *sql.DB
}

// NewClient creates a new client already connected to the database
func NewClient(cf Config) (*Client, error) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cf.Host, cf.Port, cf.User, cf.Password, cf.Name,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Client{db}, nil
}

// Close closes the underlying database connection
func (c *Client) Close() {
	c.db.Close()
}

// Ping the database, making sure it's still available establishing a connection if needed
func (c *Client) Ping(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := c.db.PingContext(ctx); err != nil {
		return err
	}

	return nil
}
