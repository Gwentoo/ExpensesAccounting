package cache

import (
	"context"
	"encoding/json"
	"time"
)

type NewPass struct {
	Email string `json:"email"`
	Token string `json:"code"`
}

func (r *DB) SavePendingPass(ctx context.Context, u *NewPass) error {
	data, err := json.Marshal(u)
	if err != nil {
		return err
	}

	return r.Client.Set(ctx, "newPass:"+u.Email, data, 10*time.Minute).Err()
}

func (r *DB) GetPendingPass(ctx context.Context, email string) (*NewPass, error) {
	val, err := r.Client.Get(ctx, "newPass:"+email).Result()
	if err != nil {
		return nil, err
	}

	var u NewPass
	if err := json.Unmarshal([]byte(val), &u); err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *DB) DeletePendingPass(ctx context.Context, email string) error {
	return r.Client.Del(ctx, "newPass:"+email).Err()
}
