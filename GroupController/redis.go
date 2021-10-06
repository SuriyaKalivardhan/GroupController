package main

import (
	"crypto/tls"
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis/v7"
)

func initRedis(host string, passwd string) *redis.Client {
	options := redis.Options{
		Addr:     host,
		Password: passwd,
	}

	if !strings.Contains(host, "localhost") {
		options.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}

	client := redis.NewClient(&options)
	if _, err := client.Ping().Result(); err != nil {
		log.Printf("Redis Init exception retrying %v", err)
		time.Sleep(1 * time.Second)
		return initRedis(host, passwd) //UGLY Recursion, except stack overflow crash and container restarts
	}
	log.Println("Successfull redis Init")
	return client
}
