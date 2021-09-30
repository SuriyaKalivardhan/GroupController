package main

import (
	"fmt"
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
	end := make(chan string)
	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(5 * time.Second)
			result := fmt.Sprintf("This is Complete %v", i)
			end <- result
		}
	}()

	for {
		select {
		case result := <-end:
			log.Printf("This is complete %v", result)
			//return
		case <-time.After(1 * time.Second):
			keys, err := client.Keys("*").Result()
			if err != nil {
				log.Printf("Unexpected error while Printing keys %v", err)
			} else {
				for _, key := range keys {
					val, err := client.Get(key).Result()
					if err != nil {
						log.Printf("Unexpected error while Printing keys %v", err)
					} else {
						log.Printf("KEY: %v, VALUE: %v", key, val)
					}
				}
			}
		}
	}
}
