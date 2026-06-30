package db

import (
	"context"
	"logarda/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client
var errorQueue string
var metricQueue string
var predictedMetricQueue string
var RedisNil = redis.Nil

func CreateRedisClient() {
	// create new redis client and connect to redis server
	Redis = redis.NewClient(&redis.Options{Addr: config.REDIS_URL, ReadTimeout: -1, PoolSize: 10, ConnMaxIdleTime: 1 * time.Minute})

	// init queue keys
	errorQueue = "error_messages"
	metricQueue = "metrics"
	predictedMetricQueue = "stg_predicted_metrics"
}

// use brop in production for workers (blocking pop - waits until queue has messages to pop, wont return nil)
// @ sub-goroutine
func ConsumeErrorEvents() (string, error) {
	ctx := context.Background()
	result, err := Redis.BRPop(ctx, 30*time.Second, errorQueue).Result()

	if err != nil {
		return "", err
	}
	// get the error event (json string) and send it to LLM api to unmarshal and make inference
	errorEvent := result[1] // value

	return errorEvent, nil
}

func ConsumeMetricEvents() (string, error) {
	ctx := context.Background()
	result, err := Redis.BRPop(ctx, 30*time.Second, metricQueue).Result()

	if err != nil {
		return "", err
	}
	// get the error event (json string) and send it to LLM api to unmarshal and make inference
	metricEvent := result[1] // value

	return metricEvent, nil
}

func ConsumePredictedMetricEvents() (string, error) {
	ctx := context.Background()
	result, err := Redis.BRPop(ctx, 30*time.Second, predictedMetricQueue).Result()

	if err != nil {
		return "", err
	}
	// get the metric event (json string) and send it to LLM api to unmarshal and make inference
	predictedMetricEvent := result[1] // value

	return predictedMetricEvent, nil
}
