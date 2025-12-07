package users

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB *pgxpool.Pool
}

func (r *Repository) FindByEmail(email string) (*User, error) {
	row := r.DB.QueryRow(context.Background(),
		`SELECT id, email, password_hash FROM users WHERE email = $1`, email)

	u := User{}
	err := row.Scan(&u.Id, &u.Email, &u.PasswordHash)
	return &u, err
}

func (r *Repository) Create(email, passwordHash string) (int, error) {
	row := r.DB.QueryRow(context.Background(),
		`INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id`,
		email, passwordHash)

	var id int
	return id, row.Scan(&id)
}
