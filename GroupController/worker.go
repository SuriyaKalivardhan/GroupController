package main

import (
	"fmt"
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

	clientsRedisFromWorker := initRedis(request.RedisHost, request.RedisPassword)
	clientsRedisFromWorker.Publish(request.RegisterChannel, fmt.Sprintf("Publishing Register from %v targeting %v", w.id, request.Id))

	redisClient.Publish(w.listenerChannel, fmt.Sprintf("Publishing Register from %v targeting %v", w.id, request.Id))
	w.client = request.Id
	return true
}

func (w *worker) UnRegister(request *AllocateRequest, redisClient *redis.Client) bool {
	if w.client != request.Id {
		log.Println("Unexpected allocation while this worker is assigned")
		return false
	}

	clientsRedisFromWorker := initRedis(request.RedisHost, request.RedisPassword)
	clientsRedisFromWorker.Publish(request.RegisterChannel, fmt.Sprintf("Publishing UnRegister from %v targeting %v", w.id, request.Id))
	w.client = ""

	redisClient.Publish(w.listenerChannel, fmt.Sprintf("Publishing UnRegister from %v targeting %v", w.id, request.Id))
	return true
}
