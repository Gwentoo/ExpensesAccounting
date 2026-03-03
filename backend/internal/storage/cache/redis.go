package cache

import "github.com/go-redis/redis/v8"

type DB struct {
	Client *redis.Client
}

func NewRedisDB(addr, password string, db int) *DB {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &DB{Client: client}
}
