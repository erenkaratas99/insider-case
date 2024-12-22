package pkg

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"

	"time"
)

func NewRedisClient(connectionUri string) (*redis.Client, error) {
	rc := redis.NewClient(&redis.Options{
		Addr:     connectionUri,
		Password: "",
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rc.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return rc, nil
}
