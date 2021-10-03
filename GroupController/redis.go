package main

import (
	"crypto/tls"
	"log"
	"time"

	"github.com/go-redis/redis/v7"
)

func initRedis(host string, passwd string) *redis.Client {
	options := redis.Options{
		Addr:     host,
		Password: passwd,
	}
	options.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	client := redis.NewClient(&options)
	if _, err := client.Ping().Result(); err != nil {
		log.Printf("Redis Init exceptoin retrying %v", err)
		time.Sleep(1)
		return initRedis(host, passwd)
	}
	return client
}
