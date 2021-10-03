package main

import (
	"encoding/json"
	"log"

	"github.com/go-redis/redis/v7"
)

type worker struct {
	workerId        string
	listenerChannel string
	controllerId    string
	close           func()
}

func NewWorker(message WorkerMessage, close func()) *worker {
	return &worker{
		workerId:        message.WorkerId,
		listenerChannel: message.ListenerChannel,
		controllerId:    message.Controller,
		close:           close,
	}
}

func (w *worker) Update(message *WorkerMessage) {
	if message.WorkerId != w.workerId {
		log.Printf("Mismatch boot assignment %v for %v", message.WorkerId, w.workerId)
		return
	}
	if message.ListenerChannel != w.listenerChannel {
		log.Printf("Mismatch boot assignment %v for %v", message.WorkerId, w.workerId)
		w.listenerChannel = message.ListenerChannel
	}
	if message.Controller != w.controllerId {
		log.Printf("Change in client from  %v to %v, reassigning", message.Controller, w.controllerId)
		w.controllerId = message.Controller
	}
}

func (w *worker) Register(request *AllocateRequest, redisClient *redis.Client) bool {
	if w.controllerId != "" {
		log.Println("Unexpected allocation while this worker is assigned")
		return false
	}

	controlMessage := ControlMessage{
		Method:        "Bind",
		ControllerId:  request.ControllerId,
		RedisAddress:  request.RedisAddress,
		RedisPassword: request.RedisPassword,
		RedisUseSSL:   request.RedisUseSSL,
	}
	redisMessage, err := json.Marshal(&controlMessage)

	if err != nil {
		log.Printf("Unexcepted exception while serializing %v", err)
	}

	//below two lines are needed, just for test purpose
	// clientsRedisFromWorker := initRedis(request.RedisHost, request.RedisPassword)
	// clientsRedisFromWorker.Publish(request.RegisterChannel, redisMessage)

	redisClient.Publish(w.listenerChannel, redisMessage)
	w.controllerId = request.ControllerId
	return true
}

func (w *worker) UnRegister(request *AllocateRequest, redisClient *redis.Client) bool {
	if w.controllerId != request.ControllerId {
		log.Println("Unexpected allocation while this worker is assigned")
		return false
	}

	controlMessage := ControlMessage{
		Method:        "UnBind",
		ControllerId:  request.ControllerId,
		RedisAddress:  request.RedisAddress,
		RedisPassword: request.RedisPassword,
		RedisUseSSL:   request.RedisUseSSL,
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
