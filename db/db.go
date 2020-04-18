package db

import (
	"errors"

	"github.com/gomodule/redigo/redis"
)

// Error404 is 404 err for store, when key is not found
var Error404 = errors.New("key not found")

// RedisStore is an implementation of the Storer Interface backed by Redis
type RedisStore struct {
	client *redis.Conn
}

// NewRedisStore returns a new Store Instance
func NewRedisStore(address string) *redis.Conn {
	var cache redis.Conn
	// Initialize the redis connection to a redis instance running on your local machine
	conn, err := redis.DialURL(address)
	if err != nil {
		panic(err)
	}
	cache = conn
	return &cache
}
