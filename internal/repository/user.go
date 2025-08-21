package repository

import (
	"database/sql"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"fmt"
	"mathly/internal/infrastructure"
	"mathly/internal/models"
)

type User interface {
	GetByID(id uuid.UUID) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Insert(user *models.User) (models.User, error)
}

type user struct {
	db infrastructure.Postgres
}

func newUser(db infrastructure.Postgres) User {
	return &user{db}
}

func (u *user) GetByEmail(email string) (*models.User, error) {
	var user models.User

	row := u.db.QueryRow(`SELECT id, email, nickname, password_hash, created_at, updated_at FROM users WHERE email = $1 ;`, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Nickname,
		&user.Hash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &user, fmt.Errorf("failure while looking up user with email %s : %v", email, err)
	}

	return &user, nil
}

func (u *user) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User

	row := u.db.QueryRow(`SELECT id, email, nickname, password_hash, created_at, updated_at FROM users WHERE id = $1 ;`, id)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Nickname,
		&user.Hash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &user, fmt.Errorf("failure while looking up user with id %s : %v", id, err)
	}

	return &user, nil
}

func (u *user) Insert(user *models.User) (models.User, error) {
	_, err := u.db.Query(`
	INSERT INTO users (
		id, email, nickname, password_hash, created_at, updated_at
	) VALUES (
		$1, $2, $3, $4, $5, $6
	);`,
		user.ID.String(),
		user.Email,
		user.Nickname,
		user.Hash,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return *user, fmt.Errorf("failure while inserting user: %s", err)
	}

	return *user, nil
}
