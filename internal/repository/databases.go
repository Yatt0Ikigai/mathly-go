package repository

import (
	"context"
	"fmt"
	"mathly/internal/infrastructure"
	"strings"
)

type Databases interface {
	Close()
	Health(ctx context.Context) error
	DB() infrastructure.Postgres
}

type databases struct {
	db infrastructure.Postgres
}

type DatabasesConfig struct {
	Redis RedisConfig `json:"redis"`
	SQL   SQLConfig   `json:"sql"`
}

type RedisConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`

	DB       string `json:"db"`
	User     string `json:"user"`
	Password string `json:"password"`
	Schema   string `json:"schema"`
}

type SQLConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`

	DB       string `json:"db"`
	User     string `json:"user"`
	Password string `json:"password"`
	Schema   string `json:"schema"`
}

func NewDatabases(config *DatabasesConfig) (Databases, error) {
	sqlConfig := config.SQL
	var errMsgs []string

	db, err := infrastructure.NewPostgres(infrastructure.PostgresConfig{
		Path:                     fmt.Sprintf(`host=%s port=%d user=%s password=%s dbname=%s search_path=%s sslmode=disable`, sqlConfig.Host, sqlConfig.Port, sqlConfig.User, sqlConfig.Password, sqlConfig.DB, sqlConfig.Schema),
		MaxIdleConns:             10,
		MaxOpenConns:             10,
		ConnMaxLifetimeInMinutes: 120,
	})
	if err != nil {
		errMsgs = append(errMsgs, fmt.Sprintf("Cannot connect to postgres db - details: %s", err.Error()))
	}

	if len(errMsgs) > 0 {
		return nil, fmt.Errorf("errors occurred while initializing databases connections: %s", strings.Join(errMsgs, "; "))
	}

	return &databases{
		db: db,
	}, nil
}

func (d *databases) Close() {
	d.db.Close()
}

func (d *databases) Health(ctx context.Context) error {
	var errMsgs []string

	err := d.db.Health()
	if err != nil {
		errMsgs = append(errMsgs, fmt.Sprintf("Cannot connect to mysql db - details: %s", err.Error()))
	}

	if len(errMsgs) > 0 {
		return fmt.Errorf("errors occurred while initializing databases connections: %s", strings.Join(errMsgs, "; "))
	}

	return nil
}

func (d *databases) DB() infrastructure.Postgres {
	return d.db
}
