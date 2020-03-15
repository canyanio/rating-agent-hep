package state

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/pkg/errors"
)

// RedisManager is the Redis state manager
type RedisManager struct {
	client        *redis.Client
	redisAddress  string
	redisPassword string
	redisDb       int
}

// NewRedisManager returns a new Redis Manager objecft
func NewRedisManager(redisAddress, redisPassword string, redisDb int) *RedisManager {
	return &RedisManager{
		redisAddress:  redisAddress,
		redisPassword: redisPassword,
		redisDb:       redisDb,
	}
}

// Connect connects to Redis
func (m *RedisManager) Connect(context context.Context) error {
	client := redis.NewClient(&redis.Options{
		Addr:     m.redisAddress,
		Password: m.redisPassword, // no password set
		DB:       m.redisDb,       // use default DB
	})
	_, err := client.Ping().Result()
	if err != nil {
		return errors.Wrap(err, "unable to connect to redis")
	}
	m.client = client
	return nil
}

// Close disconnects from Redis
func (m *RedisManager) Close(context context.Context) error {
	if m.client == nil {
		return errors.New("not connected")
	}
	err := m.client.Close()
	m.client = nil
	return err
}

// Set updates the data associated with a key
func (m *RedisManager) Set(context context.Context, key string, data interface{}, ttl int) error {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "unable to marshal request to JSON")
	}
	err = m.client.Set(key, dataJSON, time.Duration(ttl)*time.Second).Err()
	return err
}

// Get retrives the data associated with a key
func (m *RedisManager) Get(context context.Context, key string, destination interface{}) error {
	dataJSON, err := m.client.Get(key).Result()
	if err == redis.Nil {
		return nil
	} else if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(dataJSON), destination)
	if err != nil {
		return errors.Wrap(err, "unable to marshal request to JSON")
	}
	return nil
}

// Delete deletes a key and its associated data
func (m *RedisManager) Delete(context context.Context, key string) error {
	err := m.client.Del(key).Err()
	return err
}

// Delete deletes a key and its associated data
func (m *RedisManager) flushAll(context context.Context) {
	m.client.FlushAll()
}
