package db

import (
	"context"
	"logarda/internal/config"

	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client
var errorQueue string
var metricQueue string

func CreateRedisClient() {
	// create new redis client and connect to redis server
	Redis = redis.NewClient(&redis.Options{Addr: config.REDIS_URL})

	// init queue keys
	errorQueue = "error_messages"
	metricQueue = "metrics"
}

// use brop in production for workers (blocking pop - waits until queue has messages to pop, wont return nil)
// @ sub-goroutine
func ConsumeErrorEvents() (string, error) {
	ctx := context.Background()
	result, err := Redis.BRPop(ctx, 0, errorQueue).Result()

	if err != nil {
		return "", err
	}
	// get the error event (json string) and send it to LLM api to unmarshal and make inference
	errorEvent := result[1] // value

	return errorEvent, nil
}

func ConsumeMetricEvents() (string, error) {
	ctx := context.Background()
	result, err := Redis.BRPop(ctx, 0, metricQueue).Result()

	if err != nil {
		return "", err
	}
	// get the error event (json string) and send it to LLM api to unmarshal and make inference
	metricEvent := result[1] // value

	return metricEvent, nil
}
