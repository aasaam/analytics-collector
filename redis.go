package main

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

const redisKeyRecords = "RECORDS"

func redisGetClient(connectionString string) (*redis.Client, error) {
	opt, err := redis.ParseURL(connectionString)

	if err != nil {
		return nil, err
	}

	red := redis.NewClient(opt)

	pong, pongErr := red.Ping(context.Background()).Result()

	if pongErr != nil {
		return nil, pongErr
	}

	if pong == "PONG" {
		return red, nil
	}

	return nil, errors.New("cannot connect to redis")
}
