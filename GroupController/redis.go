package main

import (
	"github.com/go-redis/redis/v7"
)

func initRedis(host string, passwd string) *redis.Client {
	options := redis.Options{
		Addr:     host,
		Password: passwd,
	}
	client := redis.NewClient(&options)
	return client
}
