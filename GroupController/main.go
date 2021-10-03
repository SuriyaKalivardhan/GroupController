package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-redis/redis/v7"
)

const PoolManagerControlChannel = "poolmanager_control"

func main() {
	log.Println("Init")
	context := context.Background()
	redisClient := getRedisClient()

	controller, _ := NewController(redisClient, context)
	go handleRedisMessages(redisClient, PoolManagerControlChannel, controller)
	go initHttpServer(controller)
	select {}
}

func getRedisClient() *redis.Client {
	redisHost := "redis-poolmanager-0.redis.cache.windows.net:6380"
	redisPasswd := "KEfR4SJAdSokMnp1Hm1G5jZsLEc+WN+PeRCjYwDD9r0="
	return initRedis(redisHost, redisPasswd)
}

func handleRedisMessages(redisClient *redis.Client, controllerChannel string, controller *Controller) {
	pubSub := redisClient.Subscribe(controllerChannel)
	for {
		select {
		case msg := <-pubSub.Channel():
			log.Printf("Received nessage %v", msg.Payload)
			controller.handleRedisMessage(msg.Payload)
		case <-controller.ctx.Done():
			log.Println("Manager no more active, Not processing REDIS message")
			return
		}
	}
}

func initHttpServer(controller *Controller) {
	http.HandleFunc("/healthcheck", controller.handleHealthcheck)
	http.HandleFunc("/reconcile", controller.handleAllocate)
	http.ListenAndServe(":5001", nil)
}
