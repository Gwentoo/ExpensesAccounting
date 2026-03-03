package cache

import (
	"context"
	"encoding/json"
	"time"
)

type PendingUser struct {
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Username string    `json:"username"`
	Code     string    `json:"code"`
	Expires  time.Time `json:"expires"`
}

func (r *DB) SavePendingUser(ctx context.Context, u *PendingUser, ttl time.Duration) error {
	data, err := json.Marshal(u)
	if err != nil {
		return err
	}

	return r.Client.Set(ctx, "pending:"+u.Email, data, ttl).Err()
}

func (r *DB) GetPendingUser(ctx context.Context, email string) (*PendingUser, error) {
	val, err := r.Client.Get(ctx, "pending:"+email).Result()
	if err != nil {
		return nil, err
	}

	var u PendingUser
	if err := json.Unmarshal([]byte(val), &u); err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *DB) DeletePendingUser(ctx context.Context, email string) error {
	return r.Client.Del(ctx, "pending:"+email).Err()
}
