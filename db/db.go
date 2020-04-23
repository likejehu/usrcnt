package db

import (
	"log"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
)

// Error404 is 404 err for store, when key is not found
var Error404 = errors.New("key not found")

// RedisCache is instance of redis storage
var RedisCache = NewRedisCache("redis://localhost:6379")

// RedisStore is an implementation of the Storer Interface backed by Redis
type RedisStore struct {
	Client redis.Conn
}

// NewRedisCache returns a new Store Instance
func NewRedisCache(address string) *RedisStore {
	var cache = &RedisStore{}
	// Initialize the redis connection to a redis instance running on your local machine
	conn, err := redis.DialURL(address)
	if err != nil {
		log.Fatal(err)
	}
	cache.Client = conn
	return cache
}

// Get fetches the specified key from the underlying store
func (r *RedisStore) Get(key string) (int, error) {

	val, err := redis.Int(r.Client.Do("GET", key))
	return val, err
}

// Set updates the specified key in the underlying store
func (r *RedisStore) Set(key string, value string) error {
	_, err := r.Client.Do("SETEX", key, "120", value)
	return errors.Wrap(err, "error: settin  with SETEX")
}

// SETNXToZero sets key (if it not exists) to zero
func (r *RedisStore) SETNXToZero(key string) error {
	_, err := r.Client.Do("SETNX", key, 0)
	return errors.Wrap(err, "error: settin  with SETNXtoZero")
}

// Increment adds +1 to count
func (r *RedisStore) Increment(key string) (int, error) {
	res, err := redis.Int(r.Client.Do("INCR", key))
	return res, errors.Wrap(err, "error: settin  with SETNXtoZero")
}

// Exists checks if the specified key is in use
func (r *RedisStore) Exists(key string) (int, error) {
	e, err := redis.Int(r.Client.Do("EXISTS", key))
	if err != nil {
		return e, errors.Wrap(err, "error while checkin if token exists")
	}

	return e, nil
}
