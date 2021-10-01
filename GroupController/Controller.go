package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v7"
)

func NewController(client *redis.Client, rootContext context.Context) (*Controller, context.CancelFunc) {
	controllerContext, cancelFunc := context.WithCancel(rootContext)
	c := &Controller{client, controllerContext, make(map[string]*worker)}
	go c.watchonWorkersListener()
	return c, cancelFunc
}

type Controller struct {
	redisClient *redis.Client
	ctx         context.Context
	workers     map[string]*worker
}

func (c *Controller) handleRedisMessage(redisPayload string) {
	var message BootChannelMessage
	err := json.Unmarshal([]byte(redisPayload), &message)
	if err != nil {
		log.Printf("Unexpected REDIS error message %v", err)
		return
	}

	if c.workers[message.Id] != nil {
		c.workers[message.Id].Update(&message)
	} else {
		log.Printf("Adding new worker %v", message.Id)
		_, cancel := context.WithCancel(c.ctx)
		worker := NewWorker(message, cancel)
		c.workers[worker.id] = worker
	}
}

func (c *Controller) watchonWorkersListener() {
	for {
		select {
		case <-c.ctx.Done():
			log.Println("context closed for Controller in WatchWorkerListern, returning")
			return
		case <-time.After(5 * time.Second):
			var channels []string
			localwokers := make(map[string]*worker, len(c.workers))
			for _, worker := range c.workers {
				channels = append(channels, worker.listenerChannel)
				localwokers[worker.listenerChannel] = worker
			}

			if len(channels) == 0 {
				continue
			}
			result, err := c.redisClient.PubSubNumSub(channels...).Result()
			if err != nil {
				log.Printf("Unexpected error while Wathicng on listeners %v", err)
				continue
			}

			for channel, count := range result {
				if localwokers[channel] == nil {
					continue
				}

				if count == 0 {
					log.Printf("No subsribre for %v, removing %s", channel, localwokers[channel].id)
					delete(c.workers, localwokers[channel].id)
					localwokers[channel].close()
				}
			}
		}
	}
}

func (c *Controller) handleHealthcheck(writer http.ResponseWriter, request *http.Request) {
	log.Printf("Recived the request: %v", request.Method)
	writer.Write([]byte("Okay"))
}

func (c *Controller) handleAllocate(writer http.ResponseWriter, request *http.Request) {
	log.Println("Recieved the request")
	var allocateRequest AllocateRequest
	err := json.NewDecoder(request.Body).Decode(&allocateRequest)
	if err != nil {
		log.Printf("Unexpected response while parsing the request %v", err)
	}

	result := "FAILURE"
	if c.Allocate(&allocateRequest) {
		result = "SUCCCESS"
	}

	writer.Write([]byte(result))
}

func (c *Controller) Allocate(request *AllocateRequest) bool {
	var freeWorkers []*worker
	var currentWorkers []*worker

	for _, worker := range c.workers {
		if worker.client == "" {
			freeWorkers = append(freeWorkers, worker)
		} else if worker.client == request.Id {
			currentWorkers = append(currentWorkers, worker)
		}
	}

	diff := min(request.Target-len(currentWorkers), len(freeWorkers))

	if diff >= 0 {
		newWorkers := freeWorkers[0:diff]
		for _, worker := range newWorkers {
			if worker.Register(request, c.redisClient) {
				currentWorkers = append(currentWorkers, worker)
			}
		}
		return len(currentWorkers) == request.Target

	} else if diff < 0 {
		removed := 0
		removeWorkers := currentWorkers[0 : diff*-1]
		for _, worker := range removeWorkers {
			if worker.UnRegister(request, c.redisClient) {
				removed++
			}
		}
		return (len(currentWorkers) - removed) == request.Target
	}

	return true
}
