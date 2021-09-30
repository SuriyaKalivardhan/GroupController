package main

import (
	"log"
	"time"

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

func monitorRedisKeys(client *redis.Client) {
	for {
		keys, err := client.Keys("*").Result()
		if err != nil {
			log.Printf("Unexpected error while Printing keys %v", err)
		} else {
			for _, key := range keys {
				log.Printf("KEY is %v", key)
				val, err := client.Get(key).Result()
				log.Printf("KEY is %v", key)
				if err != nil {
					log.Printf("Unexpected error while Printing keys %v", err)
				} else {
					log.Printf("VALUE IS %v", val)
				}
			}
		}
		time.Sleep(5 * time.Second)
	}
}
