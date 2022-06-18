package main

import (
	"net/url"
	"regexp"
	"strconv"

	"github.com/go-redis/redis/v8"
)

var redisDBPath = regexp.MustCompile(`/(?P<InitTime>[0-9]{1,2})`)

func redisNew(serverURI string) (*redis.Client, error) {
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
		Password: pass, // no password set
		DB:       db,   // use default DB
	})

	return rdb, nil

}
