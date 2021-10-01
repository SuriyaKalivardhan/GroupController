package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-redis/redis/v7"
)

const ControllerBootChannel = "ControllerBootChannel.v1"

func main() {
	log.Println("Init")
	context := context.Background()
	redisClient := getRedisClient()

	controller, _ := NewController(redisClient, context)
	go handleRedisMessages(redisClient, ControllerBootChannel, controller)
	go initHttpServer(controller)
	select {}
}

func getRedisClient() *redis.Client {
	redisHost := "localhost:6379"
	redisPasswd := ""
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
			log.Println("Controller no more active, Not processing REDIS message")
			return
		}
	}
}

func initHttpServer(controller *Controller) {
	http.HandleFunc("/healthcheck", controller.handleHealthcheck)
	http.HandleFunc("/allocate", controller.handleAllocate)
	http.ListenAndServe(":5001", nil)
}
