package infrastructure

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type Postgres interface {
	Close() error
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)

	Health() error
}

type postgres struct {
	DB *sql.DB
}

type PostgresConfig struct {
	Path                     string
	MaxIdleConns             int
	MaxOpenConns             int
	ConnMaxLifetimeInMinutes int
}

func NewPostgres(config PostgresConfig) (Postgres, error) {
	db, err := sql.Open("postgres", config.Path)
	if err != nil {
		return nil, fmt.Errorf("postgres - couldn't open db: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("postgres - connection is not healthy: %w", err)
	}

	db.SetConnMaxLifetime(time.Minute * time.Duration(config.ConnMaxLifetimeInMinutes))
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)

	return postgres{
		DB: db,
	}, nil
}

func (m postgres) Close() error {
	return m.DB.Close()
}

func (m postgres) Query(query string, args ...any) (*sql.Rows, error) {
	return m.DB.Query(query, args...)
}

func (m postgres) QueryRow(query string, args ...any) *sql.Row {
	return m.DB.QueryRow(query, args...)
}

func (m postgres) Health() error {
	return m.DB.Ping()
}

func (m postgres) Exec(query string, args ...any) (sql.Result, error) {
	return m.DB.Exec(query, args...)
}
