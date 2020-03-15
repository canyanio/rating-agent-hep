package state

import (
	"context"

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
