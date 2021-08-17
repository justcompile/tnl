package data

import (
	"context"

	"github.com/go-redis/redis/v8"
)

const domainsKey = "domains"

type RedisStore struct {
	client *redis.Client
}

func (r *RedisStore) Get(ctx context.Context, key string) (string, error) {
	res := r.client.HGet(ctx, domainsKey, key)

	return res.Result()
}

func (r *RedisStore) Save(ctx context.Context, key string, value string) error {
	r.client.HSet(ctx, domainsKey, key, value)
	return nil
}

func NewRedisStore() (*RedisStore, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	cmd := rdb.Ping(context.Background())
	if err := cmd.Err(); err != nil {
		return nil, err
	}

	store := &RedisStore{
		client: rdb,
	}

	return store, nil
}
