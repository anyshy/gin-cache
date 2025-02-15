package persist

import (
	"errors"
	"time"

	"github.com/go-redis/redis"
)

// RedisStore store http response in redis
type RedisStore struct {
	RedisClient *redis.Client
}

// NewRedisStore create a redis memory store with redis client
func NewRedisStore(redisClient *redis.Client) *RedisStore {
	return &RedisStore{
		RedisClient: redisClient,
	}
}

// Set put key value pair to redis, and expire after expireDuration
func (store *RedisStore) Set(key string, value interface{}, expire time.Duration) error {
	payload, err := Serialize(value)
	if err != nil {
		return err
	}
	return store.RedisClient.Set(key, payload, expire).Err()
}

// Delete remove key in redis, do nothing if key doesn't exist
func (store *RedisStore) Delete(key string) error {
	return store.RedisClient.Del(key).Err()
}

// Get retrieves an item from redis, if key doesn't exist, return ErrCacheMiss
func (store *RedisStore) Get(key string, value interface{}) error {
	payload, err := store.RedisClient.Get(key).Bytes()

	if errors.Is(err, redis.Nil) {
		return ErrCacheMiss
	}

	if err != nil {
		return err
	}
	return Deserialize(payload, value)
}
