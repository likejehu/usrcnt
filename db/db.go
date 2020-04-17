package db

import (
	"errors"

	"github.com/go-redis/redis/v7"
)

// Error404 is 404 err for store, when key is not found
var Error404 = errors.New("key not found")

// RedisStore is an implementation of the Storer Interface backed by Redis
type RedisStore struct {
	client *redis.Client
}

// Cache is instance of redis storage
var Cache = NewRedisStore("localhost:6379")

// NewRedisStore returns a new Store Instance
func NewRedisStore(address string) *RedisStore {
	var store = &RedisStore{}
	store.client = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return store
}
