package postgres

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/jackc/pgx/v5"
)

type User struct {
	ID        int64
	Email     string
	UserName  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (db *DB) CreateUser(ctx context.Context, email, username, password string) (int64, error) {
	id, _ := rand.Int(rand.Reader, big.NewInt(1000000000))

	_, err := db.Pool.Exec(
		ctx,
		"INSERT INTO users (user_id, email, username, password) VALUES ($1, $2, $3, $4)",
		id, email, username, password,
	)
	if err != nil {
		return 0, err
	}

	return id.Int64(), nil
}

func (db *DB) IsUserExists(ctx context.Context, email string) (bool, error) {
	var count int

	query := "SELECT COUNT(*) FROM users WHERE email = $1"

	err := db.Pool.QueryRow(ctx, query, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return count > 0, nil
}

func (db *DB) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User

	query := "SELECT * FROM users WHERE email = $1"

	err := db.Pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.UserName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, fmt.Errorf("user is not exists")
		}

		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

func (db *DB) GetUserByID(ctx context.Context, id int64) (*User, error) {
	var user User

	query := "SELECT * FROM users WHERE user_id = $1"

	err := db.Pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.UserName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, fmt.Errorf("user is not exists")
		}

		return nil, fmt.Errorf("failed to get user by id: %v", err)
	}

	return &user, nil
}

func (db *DB) UpdatePassword(ctx context.Context, id int64, newPass string) error {
	query := "UPDATE users SET password = $1, updated_at = $2 WHERE user_id = $3"
	_, err := db.Pool.Exec(
		ctx,
		query,
		newPass,
		time.Now(),
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}
