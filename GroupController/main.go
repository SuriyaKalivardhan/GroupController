package main

import (
	"context"
	"log"
	"net/http"
	"os"

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

	var redisHost, redisPasswd string
	if redisHost = os.Getenv("AZUREML_OAI_REDIS_HOST"); redisHost == "" {
		log.Println("redisHost not present in env variable setting to localhost")
		redisHost = "localhost:6379"
	}

	if redisPasswd = os.Getenv("AZUREML_OAI_REDIS_KEY"); redisPasswd == "" {
		log.Println("redisPassWord not present in env variable")
		redisPasswd = ""
	}

	return initRedis(redisHost, redisPasswd)
}

func handleRedisMessages(redisClient *redis.Client, controllerChannel string, controller *Controller) {
	pubSub := redisClient.Subscribe(controllerChannel)
	for {
		select {
		case msg := <-pubSub.Channel():
			//log.Printf("RCV %v", msg.Payload)
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
