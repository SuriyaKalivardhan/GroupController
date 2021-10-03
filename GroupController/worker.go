package main

import (
	"encoding/json"
	"log"

	"github.com/go-redis/redis/v7"
)

type worker struct {
	id              string
	listenerChannel string
	client          string
	close           func()
}

func NewWorker(message BootChannelMessage, close func()) *worker {
	return &worker{
		id:              message.Id,
		listenerChannel: message.ListenerChannel,
		client:          message.BindedClient,
		close:           close,
	}
}

func (w *worker) Update(message *BootChannelMessage) {
	if message.Id != w.id {
		log.Printf("Mismatch boot assignment %v for %v", message.Id, w.id)
		return
	}
	if message.ListenerChannel != w.listenerChannel {
		log.Printf("Mismatch boot assignment %v for %v", message.Id, w.id)
		w.listenerChannel = message.ListenerChannel
	}
	if message.BindedClient != w.client {
		log.Printf("Change in client from  %v to %v, reassigning", w.client, message.BindedClient)
		w.client = message.BindedClient
	}
}

func (w *worker) Register(request *AllocateRequest, redisClient *redis.Client) bool {
	if w.client != "" {
		log.Println("Unexpected allocation while this worker is assigned")
		return false
	}

	controlMessage := GroupControllerControlMessage{
		Method:          "Register",
		ClientId:        request.Id,
		RedisHost:       request.RedisHost,
		RedisPort:       request.RedisPort,
		RedisPassword:   request.RedisPassword,
		RegisterChannel: request.RegisterChannel,
	}
	redisMessage, err := json.Marshal(&controlMessage)

	if err != nil {
		log.Printf("Unexcepted exception while serializing %v", err)
	}

	//below two lines are needed, just for test purpose
	// clientsRedisFromWorker := initRedis(request.RedisHost, request.RedisPassword)
	// clientsRedisFromWorker.Publish(request.RegisterChannel, redisMessage)

	redisClient.Publish(w.listenerChannel, redisMessage)
	w.client = request.Id
	return true
}

func (w *worker) UnRegister(request *AllocateRequest, redisClient *redis.Client) bool {
	if w.client != request.Id {
		log.Println("Unexpected allocation while this worker is assigned")
		return false
	}

	controlMessage := GroupControllerControlMessage{
		Method:          "UnRegister",
		ClientId:        request.Id,
		RedisHost:       request.RedisHost,
		RedisPort:       request.RedisPort,
		RedisPassword:   request.RedisPassword,
		RegisterChannel: request.RegisterChannel,
	}

	redisMessage, err := json.Marshal(&controlMessage)

	if err != nil {
		log.Printf("Unexcepted exception while serializing %v", err)
	}

	//Below Three lines are not needed, just for debugging purpose
	// clientsRedisFromWorker := initRedis(request.RedisHost, request.RedisPassword)
	// clientsRedisFromWorker.Publish(request.RegisterChannel, redisMessage)
	// w.client = ""

	redisClient.Publish(w.listenerChannel, redisMessage)
	return true
}
