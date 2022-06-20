package main

import (
	"context"
	"net/url"
	"regexp"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type redisClient struct {
	listLength int64
	rdb        *redis.Client
}

const (
	redisListRecords     = "Records"
	redisListClientError = "ClientErrors"
)

var redisDBPath = regexp.MustCompile(`/(?P<InitTime>[0-9]{1,2})`)

func redisClientNew(serverURI string, listLength int64) (*redisClient, error) {
	rURI, rURIErr := url.Parse(serverURI)
	if rURIErr != nil {
		return nil, rURIErr
	}

	db := 0
	if ok := redisDBPath.MatchString(rURI.Path); ok {
		matched := redisDBPath.FindStringSubmatch(rURI.Path)
		d, dE := strconv.ParseInt(matched[1], 10, 64)
		if dE != nil {
			return nil, rURIErr
		}
		db = int(d)
	}

	pass, _ := rURI.User.Password()

	rdb := redis.NewClient(&redis.Options{
		Addr:     rURI.Host,
		Username: rURI.User.Username(),
		Password: pass,
		DB:       db,
	})

	rc := redisClient{
		listLength: listLength,
		rdb:        rdb,
	}

	return &rc, nil
}

func (rc *redisClient) countRecords() int64 {
	return rc.rdb.LLen(context.Background(), redisListRecords).Val()
}

func (rc *redisClient) pushRecord(b []byte) error {
	return rc.rdb.RPush(context.Background(), redisListRecords, b).Err()
}

func (rc *redisClient) popRecord() ([]string, error) {
	ctx := context.Background()
	results, err := rc.rdb.LRange(ctx, redisListRecords, 0, rc.listLength-1).Result()
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (rc *redisClient) popRecordSubmit() error {
	ctx := context.Background()
	return rc.rdb.LTrim(ctx, redisListRecords, rc.listLength, -1).Err()
}

func (rc *redisClient) countClientError() int64 {
	return rc.rdb.LLen(context.Background(), redisListClientError).Val()
}

func (rc *redisClient) pushClientError(b []byte) error {
	return rc.rdb.RPush(context.Background(), redisListClientError, b).Err()
}

func (rc *redisClient) popClientError() ([]string, error) {
	ctx := context.Background()
	results, err := rc.rdb.LRange(ctx, redisListClientError, 0, rc.listLength-1).Result()
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (rc *redisClient) popClientErrorSubmit() error {
	ctx := context.Background()
	return rc.rdb.LTrim(ctx, redisListClientError, rc.listLength, -1).Err()
}
