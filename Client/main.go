package main

import (
	"context"
	"log"

	"github.com/go-redis/redis/v7"
)

const ControllerBootChannel = "ContosoController.v1"

func main() {
	log.Println("Init")
	context := context.Background()
	redisClient := getRedisClient()

	controller, _ := NewWorkerInteractor(redisClient, context)
	go handleRedisMessages(redisClient, ControllerBootChannel, controller)
	select {}
}

func handleRedisMessages(redisClient *redis.Client, controllerChannel string, interactor *WorkerInteractor) {
	pubSub := redisClient.Subscribe(controllerChannel)
	for {
		select {
		case msg := <-pubSub.Channel():
			log.Printf("Received nessage %v", msg.Payload)
			interactor.handleRedisMessage(msg.Payload)
		case <-interactor.ctx.Done():
			log.Println("Interactor no more active, Not processing REDIS message")
			return
		}
	}
}
