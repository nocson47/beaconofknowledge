package redisadapters

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/nocson47/beaconofknowledge/config"
)

func ConnectRedis(cfg *config.Configuration) *redis.Client {
	addr := fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort)
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
}

func PingRedis(client *redis.Client) error {
	ctx := context.Background()
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		return err
	}
	fmt.Printf("Redis ping response: %s\n", pong)
	return nil
}
